package rest

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//sql connector server
type Server struct {
	*logrus.Logger
	*gin.Engine
	producer sarama.SyncProducer
}

func NewServer(
	logger *logrus.Logger,
	addr string,
	brokersList string, /*The Kafka brokers to connect to, as a comma separated list*/
	userName string, /*The SASL username*/
	passwd string, /*The SASL password*/
	algorithm string, /*The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism*/
	topic string, /*The Kafka topic to use*/
	certFile string, /*The optional certificate file for client authentication*/
	keyFile string, /*The optional key file for client authentication*/
	caFile string, /*The optional certificate authority file for TLS client authentication*/
	verifySSL bool, /*Optional verify ssl certificates chain*/
	useTLS bool, /*Use TLS to communicate with the cluster*/
	mode string, /*Mode to run in: \"produce\" to produce, \"consume\" to consume*/
	logMsg bool, /*True to log consumed messages to console*/
) (*Server, error) {

	// Logger
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}

	//producer
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Retry.Max = maxRetry
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(strings.Split(brokersList, ","), conf)
	if err != nil {
		panic(err)
	}

	// Gin
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Mode, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	pprof.Register(r)

	server := &Server{
		Logger:   logger,
		Engine:   r,
		producer: producer,
	}

	authorized := r.Group("/api/v1")
	{
		authorized.POST("/send-msg", server.SendMsgToBroker)
	}

	return server, nil
}
