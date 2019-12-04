package ws

import (
	"log"

	"sm-connector-be/svc/ws"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cobra command ro start kafka
var Cmd = &cobra.Command{
	Use:   "ws-start",
	Short: "Start a web sockets service",
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
)

func init() {
	Cmd.Flags().StringVarP(&addr, "address", "a", ":8089", "Address to listen on")
	Cmd.Flags().StringVar(&brokers, "brokers", ":9092", "The Kafka brokers to connect to, as a comma separated list")
	Cmd.Flags().StringVar(&userName, "username", "", "The SASL username")
	Cmd.Flags().StringVar(&passwd, "passwd", "", "The SASL password")
	Cmd.Flags().StringVar(&algorithm, "algorithm", "", "The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism")
	Cmd.Flags().StringVar(&topic, "topic", "default_topic", "The Kafka topic to use")
	Cmd.Flags().StringVar(&certFile, "certificate", "", "The optional certificate file for client authentication")
	Cmd.Flags().StringVar(&keyFile, "key", "", "The optional key file for client authentication")
	Cmd.Flags().StringVar(&caFile, "ca", "", "The optional certificate authority file for TLS client authentication")
}

func run(cmd *cobra.Command, _ []string) {
	logger := logrus.New()

	server, err := ws.NewServer(
		logger,
		addr,
		brokers,
		userName,
		passwd,
		algorithm,
		topic,
		certFile,
		keyFile,
		caFile,
	)
	if err != nil {
		logger.Fatalln("Couldn't start websocket server:", err)
	}

	log.Fatalln(server.Run(addr))
}
