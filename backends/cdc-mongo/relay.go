package cdc_mongo

import (
	"context"
	"github.com/batchcorp/plumber/api"
	"github.com/batchcorp/plumber/backends/cdc-mongo/types"
	"github.com/batchcorp/plumber/cli"
	"github.com/batchcorp/plumber/relay"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Relayer struct {
	Options        *cli.Options
	RelayCh        chan interface{}
	log            *logrus.Entry
	Service        *mongo.Client
	DefaultContext context.Context
}

func Relay(opts *cli.Options) error {

	// Create new relayer instance (+ validate token & gRPC address)
	relayCfg := &relay.Config{
		Token:       opts.RelayToken,
		GRPCAddress: opts.RelayGRPCAddress,
		NumWorkers:  opts.RelayNumWorkers,
		Timeout:     opts.RelayGRPCTimeout,
		RelayCh:     make(chan interface{}, 1),
		DisableTLS:  opts.RelayGRPCDisableTLS,
		Type:        opts.RelayType,
	}

	grpcRelayer, err := relay.New(relayCfg)
	if err != nil {
		return errors.Wrap(err, "unable to create new gRPC relayer")
	}

	// Create new service
	client, err := NewService(opts)
	if err != nil {
		return errors.Wrap(err, "unable to create mongo connection")
	}

	// Launch HTTP server
	go func() {
		if err := api.Start(opts.RelayHTTPListenAddress, opts.Version); err != nil {
			logrus.Fatalf("unable to start API server: %s", err)
		}
	}()

	// Launch gRPC Relayer
	if err := grpcRelayer.StartWorkers(); err != nil {
		return errors.Wrap(err, "unable to start gRPC relay workers")
	}

	r := &Relayer{
		Options: opts,
		RelayCh: relayCfg.RelayCh,
		Service: client,
		log:     logrus.WithField("pkg", "cdc-mongo/relay.go"),
	}

	return r.Relay()
}

func (r *Relayer) Relay() error {
	ctx := context.TODO()

	var err error
	var cs *mongo.ChangeStream

	database := r.Service.Database(r.Options.CDCMongo.Database)

	if r.Options.CDCMongo.Collection == "" {
		cs, err = database.Watch(ctx, mongo.Pipeline{})
	} else {
		coll := database.Collection(r.Options.CDCMongo.Collection)
		cs, err = coll.Watch(ctx, mongo.Pipeline{})
	}

	if err != nil {
		return errors.Wrap(err, "could not begin change stream")
	}

	defer cs.Close(ctx)

	for {
		if !cs.Next(ctx) {
			r.log.Errorf("unable to read message from mongo: %s", cs.Err())
		}

		r.RelayCh <- &types.RelayMessage{
			Value: cs.Current,
		}
	}

	return nil
}
