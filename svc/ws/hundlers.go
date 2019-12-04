package ws

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

func (s *Server) PingTest(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		AbortWithError(c, 500, `P/R/0001`, fmt.Sprintf("Failed to set websocket upgrade: %+v", err), err, nil)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func (s *Server) EchoMsg(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		AbortWithError(c, 500, `P/R/0001`, fmt.Sprintf("Failed to set websocket upgrade: %+v", err), err, nil)
		return
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			AbortWithError(c, 500, `P/R/0001`, "Failed to read message: ", err, nil)
			return
		}
		// Print the message to the console
		s.Logger.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			AbortWithError(c, 500, `P/R/0001`, "Failed to write message: ", err, nil)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				AbortWithError(c, 500, `P/R/0001`, "Failed to close connection: ", err, nil)
				return
			}
		}()
	}
}

func (s *Server) Notifications(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		AbortWithError(c, 500, `P/R/0001`, fmt.Sprintf("Failed to set websocket upgrade: %+v", err), err, nil)
		return
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			AbortWithError(c, 500, `P/R/0001`, "Failed to read message: ", err, nil)
			return
		}
		// Print the message to the console
		s.Logger.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		topicMsg := &sarama.ProducerMessage{
			Topic: s.topic,
			Value: sarama.StringEncoder(string(msg)),
		}

		partition, offset, err := s.producer.SendMessage(topicMsg)
		if err != nil {
			AbortWithError(c, 500, `P/R/0001`, `Failed to store your data`, err, nil)
			return
		}
		s.Logger.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", s.topic, partition, offset)

		// Write message back to browser
		if err = conn.WriteMessage(msgType, []byte(fmt.Sprintf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", s.topic, partition, offset))); err != nil {
			AbortWithError(c, 500, `P/R/0001`, "Failed to write message: ", err, nil)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				AbortWithError(c, 500, `P/R/0001`, "Failed to close connection: ", err, nil)
				return
			}
		}()
	}
}
