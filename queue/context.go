package queue

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	QueueObj  amqp.Delivery
	Channel   *amqp.Channel
	Services  Services
	ResultMap map[string]string
}

func (c *Context) Ack(msg amqp.Publishing) string {
	if c.QueueObj.ReplyTo != "" {
		if string(msg.Body) == "" {
			return fmt.Sprintf("传入内容为空，无法ack并回传内容至replyTo：%v", c.QueueObj.ReplyTo)
		}
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
func (c *Context) GetResultMap() map[string]string {
	c.ResultMap["队列名称"] = c.QueueObj.RoutingKey
	c.ResultMap["队列传入数据"] = string(c.QueueObj.Body)
	return c.ResultMap
}

// defer capture exceptions. Then Ack queue with key("IsReplyTo") in context.ResultMap
// if not ack without replayMsg to specify routingKey
// attention, the value of key("返回数据") and "IsReplyTo"  must be set before ending
func (c *Context) Defer(l *golog.Log) {
	var msg string
	// Capture panic
	if err := recover(); err != nil {
		msg = "程序异常" + err.(string)
		c.ResultMap["返回数据"] = msg
	}
	//
	resultStr, err := json.Marshal(c.ResultMap)
	if err != nil {
		fmt.Println(err)
	}
	l.PrintLocal(string(resultStr), "info")
	// if the value of IsReplyTo in context.ResultMap contains "y" or "Y", result will be push to replyTo
	if strings.Contains(strings.ToLower(c.ResultMap["IsReplyTo"]), "y") {
		pub := amqp.Publishing{
			//Headers:         nil,
			//ContentType:     "",
			//ContentEncoding: "",
			//DeliveryMode:    0,
			//Priority:        0,
			//CorrelationId:   "",
			//ReplyTo:         "",
			//Expiration:      "",
			//MessageId:       "",
			Timestamp: time.Time{},
			Body:      []byte(c.ResultMap["返回数据"]),
		}
		c.Ack(pub)
		return
	}

	fmt.Println(c.Ack(amqp.Publishing{}))
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

// 封装context 构造函数
func NewContext() *Context {
	c := new(Context)
	c.ResultMap = make(map[string]string)
	return c
}
