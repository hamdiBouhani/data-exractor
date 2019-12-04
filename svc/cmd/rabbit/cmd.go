package rabbit

import (
	"sm-connector-be/svc/rabbit_mq"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//docker run -d --hostname localhost --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

var (
	brokers   string /*The Kafka brokers to connect to, as a comma separated list*/
	topic     string /*The Kafka topic to use*/
	url       string /*Dial accepts a string in the AMQP URI format and returns a new Connection over TCP using PlainAuth.*/
	queueName string /*queue name*/
)

// cobra command to start kafka
var Cmd = &cobra.Command{
	Use:   "start-amqp",
	Short: "Start a rabbit mq consomer",
	Args:  cobra.NoArgs,
	Run:   run,
}

func init() {
	Cmd.Flags().StringVar(&brokers, "brokers", ":9092", "The Kafka brokers to connect to, as a comma separated list")
	Cmd.Flags().StringVar(&topic, "topic", "default_topic", "The Kafka topic to use")
	Cmd.Flags().StringVar(&url, "mp-url", "amqp://guest:guest@localhost:5672/", "Dial accepts a string in the AMQP URI format and returns a new Connection")
	Cmd.Flags().StringVar(&queueName, "queue-name", "test", "queue name")
}

func run(cmd *cobra.Command, _ []string) {
	logger := logrus.New()

	_, err := rabbit_mq.NewServer(
		logger,
		brokers,
		topic,
		url, /*Dial accepts a string in the AMQP URI format and returns a new Connection over TCP using PlainAuth.*/
		queueName,
	)
	if err != nil {
		logger.Fatalln("Couldn't start websocket server:", err)
	}

}
