package rabbit_mq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type JwtUser struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Channel struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type GetMessage struct {
	Content   string  `json:"content"`
	MessageId string  `json:"messageId"`
	Channel   Channel `json:"channel"`
	SendTo    JwtUser `json:"sendTo"`
}

type Server struct {
	*logrus.Logger
	producer sarama.SyncProducer

	topic             string
	AMQPConnectionURL string /*"amqp://guest:guest@localhost:5672/"*/
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", err, msg)
	}
}

func NewServer(
	logger *logrus.Logger,
	brokersList string, /*The Kafka brokers to connect to, as a comma separated list*/
	topic string, /*The Kafka topic to use*/
	url string, /*Dial accepts a string in the AMQP URI format and returns a new Connection over TCP using PlainAuth.*/
	queueName string,
) (*Server, error) {
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}

	//producer
	// conf := sarama.NewConfig()
	// conf.Producer.RequiredAcks = sarama.WaitForAll
	// conf.Producer.Return.Successes = true
	// producer, err := sarama.NewSyncProducer(strings.Split(brokersList, ","), conf)
	// if err != nil {
	// 	panic(err)
	// }

	// Rabbit MQ consumer
	conn, err := amqp.Dial(url)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	/*
	   A connection is a TCP connection from the client to server.
	   A connection is not cheap to create.
	   A channel serves as the communication protocol over the connection.
	   Channels are quite cheap.
	   We should aim to limit connections to a minimum number while establishing
	   as many channels as we need, on top of these connections.
	*/
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	//we can  now start talking to rabbit mq.
	//we need tall  the server about the queue we're interested in.
	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure Qos")

	//start consuming messages
	//we will get a go channel. We can range over this channel to get messages.
	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")
	stopchan := make(chan bool)
	go func() {
		log.Printf("Consumer ready ,PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message : %s", d.Body)
			getMsg := &GetMessage{}
			err := json.Unmarshal(d.Body, getMsg)
			if err != nil {
				log.Printf("Msg : %v ", getMsg)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	<-stopchan

	server := &Server{
		Logger: logger,
		//producer: producer,
		topic: topic,
	}
	return server, nil
}

func (s *Server) ProduceToKafka(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return errors.Wrapf(err, "fail to stored message in topic")
	}
	s.Logger.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
