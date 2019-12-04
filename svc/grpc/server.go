package grpc

import (
	"log"
	"net"
	"strings"

	"sm-connector-be/pb"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	Address string
	*logrus.Logger
	producer sarama.SyncProducer
}

func NewServer(
	logger *logrus.Logger,
	listenAddr string,
	grpcTLSCertPath string,
	grpcTLSKeyPath string,
	brokers string,
	topic string,
) (*Server, error) {
	// Set up logger.
	if logger == nil {
		logger = logrus.New()
	}

	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Retry.Max = maxRetry
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), conf)
	if err != nil {
		panic(err)
	}

	return &Server{
		Address:  listenAddr,
		Logger:   logger,
		producer: producer,
	}, nil
}

// Run the server.
func (s *Server) Run() error {
	s.Logger.Infoln("Starting Grpc API server listening on:", s.Address)
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Printf("Failed to listen : %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMessageServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
		return err
	}
	return nil
}

func (s *Server) Produce(topic, message string) error {
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
