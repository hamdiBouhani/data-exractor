package rest

import (
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

type MsgReq struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func (s *Server) SendMsgToBroker(c *gin.Context) {
	var req MsgReq
	if err := c.ShouldBind(&req); err != nil {
		BindJsonErr(c, err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: req.Topic,
		Value: sarama.StringEncoder(req.Message),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		AbortWithError(c, 500, `P/R/0001`, `Failed to store your data`, err, nil)
		return
	}
	s.Logger.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", req.Topic, partition, offset)

}
