package grpc

import (
	"log"

	"sm-connector-be/svc/grpc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd is the exported command.
var Cmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start a Grpc API server ",
	Args:  cobra.NoArgs,
	Run:   run,
}

var (
	addr     string
	grpcCert string
	grpcKey  string
	brokers  string
	topic    string
)

func init() {
	Cmd.Flags().StringVarP(&addr, "address", "a", ":50052", "Address to listen on")
	Cmd.Flags().StringVar(&grpcCert, "grpc-cert", "", "cert file for gRPC TLS authentication")
	Cmd.Flags().StringVar(&grpcKey, "grpc-key", "", "key file for gRPC TLS authentication")
	Cmd.Flags().StringVar(&brokers, "brokers", ":9092", "The Kafka brokers to connect to, as a comma separated list")
	Cmd.Flags().StringVar(&topic, "topic", "default_topic", "The Kafka topic to use")
}

func run(cmd *cobra.Command, _ []string) {
	logger := logrus.New()

	server, err := grpc.NewServer(
		logger,
		addr,
		grpcCert,
		grpcKey,
		brokers,
		topic,
	)
	if err != nil {
		logger.Fatalln("Couldn't start server:", err)
	}

	log.Fatalln(server.Run())
}
