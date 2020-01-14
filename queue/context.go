package queue

import (
	"fmt"
	"github.com/streadway/amqp"
)
// struct for context in callable function that defined by yourself.
type Context struct {
	QueueObj amqp.Delivery
	Channel  *amqp.Channel
}

func (c *Context) Ack(msg amqp.Publishing) string {
	if c.QueueObj.ReplyTo != "" {
		c.NextTo(c.QueueObj.Exchange, c.QueueObj.ReplyTo, msg)
	}
	c.QueueObj.Ack(false)
	return "Ack Success"
}
func (c *Context) Nack() {
	c.QueueObj.Nack(false, false)
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
