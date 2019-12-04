package ws

import (
	"strings"

	"sm-connector-be/common"

	"github.com/Shopify/sarama"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*logrus.Logger
	*gin.Engine
	producer sarama.SyncProducer

	topic string
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
) (*Server, error) {
	// Logger
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}

	//producer
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
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
		topic:    topic,
	}

	//WebSocket client
	r.LoadHTMLFiles(common.GetWebSocketClient("ping.html"), common.GetWebSocketClient("client.html"), common.GetWebSocketClient("test.html"))
	r.GET("/ping", func(c *gin.Context) {
		c.HTML(200, "ping.html", nil)
	})
	r.GET("/client", func(c *gin.Context) {
		c.HTML(200, "client.html", nil)
	})

	r.GET("/test", func(c *gin.Context) {
		c.HTML(200, "test.html", nil)
	})

	websocket := r.Group("/ws")
	websocket.Use(func(c *gin.Context) {})
	{
		websocket.GET("/ping", server.PingTest)
		websocket.GET("/echo", server.EchoMsg)
		websocket.GET("/noticfications", server.Notifications)
	}

	return server, nil
}
