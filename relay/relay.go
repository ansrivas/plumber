package relay

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/batchcorp/schemas/build/go/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	sqsTypes "github.com/batchcorp/plumber/backends/aws-sqs/types"
	azureTypes "github.com/batchcorp/plumber/backends/azure/types"
	mongoTypes "github.com/batchcorp/plumber/backends/cdc-mongo/types"
	postgresTypes "github.com/batchcorp/plumber/backends/cdc-postgres/types"
	gcpTypes "github.com/batchcorp/plumber/backends/gcp-pubsub/types"
	kafkaTypes "github.com/batchcorp/plumber/backends/kafka/types"
	rabbitTypes "github.com/batchcorp/plumber/backends/rabbitmq/types"
	redisTypes "github.com/batchcorp/plumber/backends/rpubsub/types"
	rstreamsTypes "github.com/batchcorp/plumber/backends/rstreams/types"
	"github.com/batchcorp/plumber/stats"
)

const (
	DefaultNumWorkers = 10

	QueueFlushInterval = 10 * time.Second
	DefaultBatchSize   = 100 // number of messages to batch

	MaxGRPCRetries = 5

	// Maximum message size for GRPC client in bytes
	MaxGRPCMessageSize = 1024 * 1024 * 100 // 100MB
	GRPCRetrySleep     = time.Second * 5
)

type Relay struct {
	Config *Config
	log    *logrus.Entry
}

type Config struct {
	Token       string
	GRPCAddress string
	NumWorkers  int
	BatchSize   int
	RelayCh     chan interface{}
	DisableTLS  bool
	Timeout     time.Duration // general grpc timeout (used for all grpc calls)
	Type        string
}

func New(relayCfg *Config) (*Relay, error) {
	if err := validateConfig(relayCfg); err != nil {
		return nil, errors.Wrap(err, "unable to complete relay config validation")
	}

	// Verify grpc connection & token
	if err := TestConnection(relayCfg); err != nil {
		return nil, errors.Wrap(err, "unable to complete connection test")
	}

	// JSON formatter for log output if not running in a TTY - colors are fun!
	if !terminal.IsTerminal(int(os.Stderr.Fd())) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	return &Relay{
		Config: relayCfg,
		log:    logrus.WithField("pkg", "relay"),
	}, nil
}

func validateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("Relay config cannot be nil")
	}

	if cfg.Token == "" {
		return errors.New("Token cannot be empty")
	}

	if cfg.GRPCAddress == "" {
		return errors.New("GRPCAddress cannot be empty")
	}

	if cfg.RelayCh == nil {
		return errors.New("RelayCh cannot be nil")
	}

	if cfg.NumWorkers <= 0 {
		logrus.Warningf("NumWorkers cannot be <= 0 - setting to default '%d'", DefaultNumWorkers)
		cfg.NumWorkers = DefaultNumWorkers
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = DefaultBatchSize
	}

	return nil
}

func TestConnection(cfg *Config) error {
	conn, ctx, err := NewConnection(cfg.GRPCAddress, cfg.Token, cfg.Timeout, cfg.DisableTLS, false)
	if err != nil {
		return errors.Wrap(err, "unable to create new connection")
	}

	// Call the Test method to verify connectivity
	c := services.NewGRPCCollectorClient(conn)

	if _, err := c.Test(ctx, &services.TestRequest{}); err != nil {
		return errors.Wrap(err, "unable to complete Test request")
	}

	return nil
}

func NewConnection(address, token string, timeout time.Duration, disableTLS, noCtx bool) (*grpc.ClientConn, context.Context, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if !disableTLS {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{
				InsecureSkipVerify: true,
			},
		)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	dialContext, _ := context.WithTimeout(context.Background(), timeout)

	conn, err := grpc.DialContext(dialContext, address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to grpc address '%s': %s", address, err)
	}

	var ctx context.Context

	if !noCtx {
		ctx, _ = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx = context.Background()
	}

	md := metadata.Pairs("batch-token", token)
	outCtx := metadata.NewOutgoingContext(ctx, md)

	return conn, outCtx, nil
}

func (r *Relay) StartWorkers() error {
	for i := 0; i != r.Config.NumWorkers; i++ {
		r.log.WithField("workerId", i).Debug("starting worker")

		conn, ctx, err := NewConnection(r.Config.GRPCAddress, r.Config.Token, r.Config.Timeout, r.Config.DisableTLS, true)
		if err != nil {
			return fmt.Errorf("unable to create new gRPC connection for worker %d: %s", i, err)
		}

		go r.Run(i, conn, ctx)
	}

	return nil
}

func (r *Relay) Run(id int, conn *grpc.ClientConn, ctx context.Context) {
	llog := r.log.WithField("relayId", id)

	llog.Debug("Relayer started")

	queue := make([]interface{}, 0)

	// This functions as an escape-vale -- if we are pumping messages *REALLY*
	// fast - we will hit max queue size; if we are pumping messages slowly,
	// the ticker will be hit and the queue will be flushed, regardless of size.
	flushTicker := time.NewTicker(QueueFlushInterval)

	stats.Incr(r.Config.Type+"-relay-producer", 0)

	for {
		select {
		case msg := <-r.Config.RelayCh:
			queue = append(queue, msg)

			// Max queue size reached
			if len(queue) >= r.Config.BatchSize {
				r.log.Debugf("%d: max queue size reached - flushing!", id)

				r.flush(ctx, conn, queue...)

				// Reset queue; since we passed by variadic, the underlying slice can be updated
				queue = make([]interface{}, 0)

				// Reset ticker (so time-based flush doesn't occur)
				flushTicker.Reset(QueueFlushInterval)
			}
		case <-flushTicker.C:
			if len(queue) != 0 {
				r.log.Debugf("%d: flush ticker hit and queue not empty - flushing!", id)

				r.flush(ctx, conn, queue...)

				// Reset queue; same as above - safe to delete queue contents
				queue = make([]interface{}, 0)
			}
		}
	}
}

func (r *Relay) flush(ctx context.Context, conn *grpc.ClientConn, messages ...interface{}) {
	if len(messages) < 1 {
		r.log.Error("asked to flush empty message queue - bug?")
		return
	}

	// We only care about the first message since plumber can only be using
	// one message bus type at a time

	var err error

	switch v := messages[0].(type) {
	case *sqsTypes.RelayMessage:
		r.log.Debugf("flushing %d sqs message(s)", len(messages))
		err = r.handleSQS(ctx, conn, messages)
	case *rabbitTypes.RelayMessage:
		r.log.Debugf("flushing %d rabbit message(s)", len(messages))
		err = r.handleRabbit(ctx, conn, messages)
	case *kafkaTypes.RelayMessage:
		r.log.Debugf("flushing %d kafka message(s)", len(messages))
		err = r.handleKafka(ctx, conn, messages)
	case *azureTypes.RelayMessage:
		r.log.Debugf("flushing %d azure message(s)", len(messages))
		err = r.handleAzure(ctx, conn, messages)
	case *gcpTypes.RelayMessage:
		r.log.Debugf("flushing %d gcp message(s)", len(messages))
		err = r.handleGCP(ctx, conn, messages)
	case *mongoTypes.RelayMessage:
		r.log.Debugf("flushing %d mongo message(s)", len(messages))
		err = r.handleCDCMongo(ctx, conn, messages)
	case *redisTypes.RelayMessage:
		r.log.Debugf("flushing %d redis-pubsub message(s)", len(messages))
		err = r.handleRedisPubSub(ctx, conn, messages)
	case *rstreamsTypes.RelayMessage:
		r.log.Debugf("flushing %d redis-streams message(s)", len(messages))
		err = r.handleRedisStreams(ctx, conn, messages)
	case *postgresTypes.RelayMessage:
		r.log.Debugf("flushing %d cdc-postgres message(s)", len(messages))
		err = r.handleCdcPostgres(ctx, conn, messages)
	default:
		r.log.WithField("type", v).Error("received unknown message type - skipping")
		return
	}

	if err != nil {
		r.log.WithField("err", err).Error("unable to handle message")
		return
	}

	stats.Incr(r.Config.Type+"-relay-producer", len(messages))
}

// CallWithRetry will retry a GRPC call until it succeeds or reaches a maximum number of retries defined by MaxGRPCRetries
func (r *Relay) CallWithRetry(ctx context.Context, method string, publish func(ctx context.Context) error) error {
	var err error

	for i := 1; i <= MaxGRPCRetries; i++ {
		err = publish(ctx)
		if err != nil {
			r.log.Debugf("unable to complete %s call [retry %d/%d]", method, i, 5)
			time.Sleep(GRPCRetrySleep)
			continue
		}
		r.log.Debugf("successfully handled %s message", strings.Replace(method, "Add", "", 1))
		return nil
	}

	return fmt.Errorf("unable to complete %s call [reached max retries (%d)]: %s", method, MaxGRPCRetries, err)
}
