package cli

import (
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultKafkaConnectTimeout = "10s"
	DefaultKafkaGroupId        = "plumber"

	// Optimize for immediate output
	DefaultKafkaMaxWait          = "1s"
	DefaultKafkaMinBytes         = "1"
	DefaultKafkaMaxBytes         = "1"
	DefaultKafkaQueueCapacity    = "1"
	DefaultKafkaRebalanceTimeout = "0"

	// Optimize for throughput + reliability
	DefaultKafkaRelayMaxWait          = "5s"
	DefaultKafkaRelayMinBytes         = "1048576" // 1MB
	DefaultKafkaRelayMaxBytes         = "1048576" // 1MB
	DefaultKafkaRelayQueueCapacity    = "1000"
	DefaultKafkaRelayRebalanceTimeout = "5s"
	DefaultKafkaRelayCommitInterval   = "5s"
	DefaultKafkaReadOffset            = "0"
)

type KafkaOptions struct {
	// Shared
	Address     string
	Topic       string
	Timeout     time.Duration
	InsecureTLS bool
	MutualTLS   bool
	CaCert      string
	ClientCert  string
	ClientKey   string

	Username           string
	Password           string
	AuthenticationType string

	// Read
	UseConsumerGroup bool
	GroupID          string
	ReadOffset       int64 // If UseConsumerGroup is true, ReadOffset will NOT be used
	MaxWait          time.Duration
	MinBytes         int
	MaxBytes         int
	QueueCapacity    int
	RebalanceTimeout time.Duration
	CommitInterval   time.Duration

	// Write
	WriteKey string
}

func HandleKafkaFlags(readCmd, writeCmd, relayCmd *kingpin.CmdClause, opts *Options) {
	rc := readCmd.Command("kafka", "Kafka message system")

	addSharedKafkaFlags(rc, opts)
	addReadKafkaFlags(rc, opts)

	// Kafka write cmd
	wc := writeCmd.Command("kafka", "Kafka message system")

	addSharedKafkaFlags(wc, opts)
	addWriteKafkaFlags(wc, opts)

	// If PLUMBER_RELAY_TYPE is set, use env vars, otherwise use CLI flags
	relayType := os.Getenv("PLUMBER_RELAY_TYPE")

	var rec *kingpin.CmdClause

	if relayType != "" {
		rec = relayCmd
	} else {
		rec = relayCmd.Command("kafka", "Kafka message system")
	}

	addSharedKafkaFlags(rec, opts)
	addReadKafkaFlags(rec, opts)
}

func addSharedKafkaFlags(cmd *kingpin.CmdClause, opts *Options) {
	cmd.Flag("address", "Destination host address").
		Default("localhost:9092").
		Envar("PLUMBER_RELAY_KAFKA_ADDRESS").
		StringVar(&opts.Kafka.Address)
	cmd.Flag("topic", "Topic to read message(s) from").
		Required().
		Envar("PLUMBER_RELAY_KAFKA_TOPIC").
		StringVar(&opts.Kafka.Topic)
	cmd.Flag("timeout", "Connect timeout").Default(DefaultKafkaConnectTimeout).
		Envar("PLUMBER_RELAY_KAFKA_TIMEOUT").
		DurationVar(&opts.Kafka.Timeout)
	cmd.Flag("insecure-tls", "Use insecure TLS").
		Envar("PLUMBER_RELAY_KAFKA_INSECURE_TLS").
		BoolVar(&opts.Kafka.InsecureTLS)
	cmd.Flag("mtls", "Use mutual TLS").
		Envar("PLUMBER_RELAY_KAFKA_MUTUAL_TLS").
		BoolVar(&opts.Kafka.MutualTLS)
	cmd.Flag("ca-cert", "Path to ca cert").
		Envar("PLUMBER_RELAY_KAFKA_CA_CERT").
		StringVar(&opts.Kafka.CaCert)
	cmd.Flag("client-cert", "Path to client cert").
		Envar("PLUMBER_RELAY_KAFKA_CLIENT_CERT").
		StringVar(&opts.Kafka.ClientCert)
	cmd.Flag("client-key", "Path to client key").
		Envar("PLUMBER_RELAY_KAFKA_CLIENT_KEY").
		StringVar(&opts.Kafka.ClientKey)
	cmd.Flag("username", "SASL Username").
		Envar("PLUMBER_RELAY_KAFKA_USERNAME").
		StringVar(&opts.Kafka.Username)
	cmd.Flag("password", "SASL Password. If omitted, you will be prompted for the password").
		Envar("PLUMBER_RELAY_KAFKA_PASSWORD").
		StringVar(&opts.Kafka.Password)
	cmd.Flag("auth-type", "SASL Authentication type (plain or scram)").
		Default("scram").
		Envar("PLUMBER_RELAY_KAFKA_SASL_TYPE").
		StringVar(&opts.Kafka.AuthenticationType)
}

func addReadKafkaFlags(cmd *kingpin.CmdClause, opts *Options) {
	defaultMaxWait := DefaultKafkaMaxWait
	defaultMaxBytes := DefaultKafkaMaxBytes
	defaultMinBytes := DefaultKafkaMinBytes
	defaultQueueCapacity := DefaultKafkaQueueCapacity
	defaultRebalanceTimeout := DefaultKafkaRebalanceTimeout

	if cmd.FullCommand() == "relay" {
		defaultMaxWait = DefaultKafkaRelayMaxWait
		defaultMaxBytes = DefaultKafkaRelayMaxBytes
		defaultMinBytes = DefaultKafkaRelayMinBytes
		defaultQueueCapacity = DefaultKafkaRelayQueueCapacity
		defaultRebalanceTimeout = DefaultKafkaRelayRebalanceTimeout
	}

	cmd.Flag("use-consumer-group", "Whether plumber should use a consumer group").
		Envar("PLUMBER_RELAY_KAFKA_USE_CONSUMER_GROUP").
		Default("true").
		BoolVar(&opts.Kafka.UseConsumerGroup)
	cmd.Flag("group-id", "Specify a specific group-id to use when reading from kafka").
		Envar("PLUMBER_RELAY_KAFKA_GROUP_ID").
		Default(DefaultKafkaGroupId).
		StringVar(&opts.Kafka.GroupID)
	cmd.Flag("read-offset", "Specify what offset the consumer should read from (only works if '--use-consumer-group' is false)").
		Envar("PLUMBER_RELAY_KAFKA_READ_OFFSET").
		Default(DefaultKafkaReadOffset).
		Int64Var(&opts.Kafka.ReadOffset)
	cmd.Flag("max-wait", "How long to wait for new data when reading batches of messages").
		Envar("PLUMBER_RELAY_KAFKA_MAX_WAIT").
		Default(defaultMaxWait).
		DurationVar(&opts.Kafka.MaxWait)
	cmd.Flag("min-bytes", "Minimum number of bytes to fetch in a single kafka request (throughput optimization)").
		Envar("PLUMBER_RELAY_KAFKA_MIN_BYTES").
		Default(defaultMinBytes).
		IntVar(&opts.Kafka.MinBytes)
	cmd.Flag("max-bytes", "Maximum number of bytes to fetch in a single kafka request (throughput optimization)").
		Envar("PLUMBER_RELAY_KAFKA_MAX_BYTES").
		Default(defaultMaxBytes).
		IntVar(&opts.Kafka.MaxBytes)
	cmd.Flag("queue-capacity", "Internal queue capacity").
		Envar("PLUMBER_RELAY_KAFKA_QUEUE_CAPACITY").
		Default(defaultQueueCapacity).
		IntVar(&opts.Kafka.QueueCapacity)
	cmd.Flag("rebalance-timeout", "How long a coordinator will wait for member joins as part of a rebalance").
		Envar("PLUMBER_RELAY_KAFKA_REBALANCE_TIMEOUT").
		Default(defaultRebalanceTimeout).
		DurationVar(&opts.Kafka.RebalanceTimeout)
	cmd.Flag("commit-interval", "How often to commit offsets to broker (0 = synchronous)").
		Envar("PLUMBER_RELAY_KAFKA_COMMIT_INTERVAL").
		Default(DefaultKafkaRelayCommitInterval).
		DurationVar(&opts.Kafka.CommitInterval)
}

func addWriteKafkaFlags(cmd *kingpin.CmdClause, opts *Options) {
	cmd.Flag("key", "Key to write to kafka (not required)").StringVar(&opts.Kafka.WriteKey)
}
