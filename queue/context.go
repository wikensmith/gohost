package queue

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/wikensmith/gohost/golog"
)

// struct for context in callable function that defined by yourself.

type Services struct {
}

func (s *Services) NewLogger() *golog.Log {
	return new(golog.Log).New()
}

type Context struct {
	QueueObj amqp.Delivery
	Channel  *amqp.Channel
	Services Services
}

func (c *Context) Ack(msg amqp.Publishing) string {
	if c.QueueObj.ReplyTo != "" {
		c.NextTo(c.QueueObj.Exchange, c.QueueObj.ReplyTo, msg)
	}
	err := c.QueueObj.Ack(false)
	if err != nil {
		fmt.Println("Ack false", err)
	}
	return "Ack Success"
}
func (c *Context) Nack() {
	err := c.QueueObj.Nack(false, false)
	if err != nil {
		fmt.Println("Nack false", err)
	}
}
func (c *Context) Text() []byte {
	return c.QueueObj.Body
}
func (c *Context) NextTo(exchangeName string, routingKey string, msg amqp.Publishing) {
	err := c.Channel.Publish(exchangeName, routingKey, false, false, msg)
	if err != nil {
		fmt.Println("MQ 消息发送失败")
	} else {
		fmt.Println("MQ 消息发送成功")
	}
}
