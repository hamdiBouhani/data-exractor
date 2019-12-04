package consumer

import (
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	kafkaclient "github.com/uber-go/kafka-client"
	"github.com/uber-go/kafka-client/kafka"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func NewConsumer(
	addr string,
	brokersList string, /*The Kafka brokers to connect to, as a comma separated list*/
	userName string, /*The SASL username*/
	passwd string, /*The SASL password*/
	algorithm string, /*The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism*/
	topic string, /*The Kafka topic to use*/
	certFile string, /*The optional certificate file for client authentication*/
	keyFile string, /*The optional key file for client authentication*/
	caFile string, /*The optional certificate authority file for TLS client authentication*/
) {

	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// mapping from cluster name to list of broker ip addresses
	brokers := map[string][]string{
		"sample_cluster":     strings.Split(brokersList, ","),
		"sample_dlq_cluster": strings.Split(brokersList, ","),
	}
	// mapping from topic name to cluster that has that topic
	topicClusterAssignment := map[string][]string{
		"sample_topic": []string{"sample_cluster"},
	}

	// First create the kafkaclient, its the entry point for creating consumers or producers
	// It takes as input a name resolver that knows how to map topic names to broker ip addrs
	client := kafkaclient.New(kafka.NewStaticNameResolver(topicClusterAssignment, brokers), zap.NewNop(), tally.NoopScope)

	// Next, setup the consumer config for consuming from a set of topics
	config := &kafka.ConsumerConfig{
		TopicList: kafka.ConsumerTopicList{
			kafka.ConsumerTopic{ // Consumer Topic is a combination of topic + dead-letter-queue
				Topic: kafka.Topic{ // Each topic is a tuple of (name, clusterName)
					Name:    topic,
					Cluster: "sample_cluster",
				},
				DLQ: kafka.Topic{
					Name:    topic,
					Cluster: "sample_dlq_cluster",
				},
			},
		},
		GroupName:   "sample_consumer",
		Concurrency: 100, // number of go routines processing messages in parallel
	}
	config.Offsets.Initial.Offset = sarama.OffsetOldest

	// Create the consumer through the previously created client
	consumer, err := client.NewConsumer(config)
	if err != nil {
		panic(err)
	}

	// Finally, start consuming
	if err := consumer.Start(); err != nil {
		panic(err)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				return // channel closed
			}
			logger.Printf(string(msg.Value()))
		case <-sigCh:
			consumer.Stop()
			<-consumer.Closed()
			return
		}
	}

}
