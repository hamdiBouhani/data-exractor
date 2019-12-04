package consumer

import (
	"sm-connector-be/svc/consumer"

	"github.com/spf13/cobra"
)

// cobra command to start kafka
var Cmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start a kafka connector",
	Args:  cobra.NoArgs,
	Run:   run,
}

var (
	addr      string
	brokers   string
	userName  string
	passwd    string
	algorithm string
	topic     string
	certFile  string
	keyFile   string
	caFile    string
	verifySSL bool
	useTLS    bool
	mode      string
	logMsg    bool
)

func init() {
	Cmd.Flags().StringVarP(&addr, "address", "a", ":8080", "Address to listen on")
	Cmd.Flags().StringVar(&brokers, "brokers", ":9092", "The Kafka brokers to connect to, as a comma separated list")
	Cmd.Flags().StringVar(&userName, "username", "", "The SASL username")
	Cmd.Flags().StringVar(&passwd, "passwd", "", "The SASL password")
	Cmd.Flags().StringVar(&algorithm, "algorithm", "", "The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism")
	Cmd.Flags().StringVar(&topic, "topic", "default_topic", "The Kafka topic to use")
	Cmd.Flags().StringVar(&certFile, "certificate", "", "The optional certificate file for client authentication")
	Cmd.Flags().StringVar(&keyFile, "key", "", "The optional key file for client authentication")
	Cmd.Flags().StringVar(&caFile, "ca", "", "The optional certificate authority file for TLS client authentication")
	Cmd.Flags().BoolVar(&verifySSL, "verify", false, "Optional verify ssl certificates chain")
	Cmd.Flags().BoolVar(&useTLS, "tls", false, "Use TLS to communicate with the cluster")
	Cmd.Flags().StringVar(&mode, "mode", "produce", "Mode to run in: \"produce\" to produce, \"consume\" to consume")
	Cmd.Flags().BoolVar(&logMsg, "logmsg", false, "True to log consumed messages to console")
}

func run(cmd *cobra.Command, _ []string) {
	consumer.NewConsumer(addr, brokers, userName, passwd, algorithm, topic, certFile, keyFile, caFile)
}
