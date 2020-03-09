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

func (c *Context) Ack(msg []byte) {
	replyTo := ""
	if c.QueueObj.Headers["ReplyTo"] != nil {
		replyTo = c.QueueObj.Headers["ReplyTo"].(string)
	}

	if replyTo != "" {
		info := c.NextTo(c.QueueObj.Exchange, replyTo, msg)
		fmt.Println(info)
	}
	err := c.QueueObj.Ack(false)
	if err != nil {
		fmt.Println("Ack false", err)
	}
}

func (c *Context) Nack() {
	err := c.QueueObj.Nack(false, true)

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
// c.ResultMap["IsReplyTo"] = "y",  replyTo is default
// c.ResultMap["LogToNet"] = "y",  日志默认打印网络和本地, LogToNet 不包含“y”  可不打印网络日志
func (c *Context) Defer(l *golog.Log) {
	var msg string
	logLevel := "info"
	logCode := 200
	// Capture panic
	if err := recover(); err != nil {
		msg = "程序异常" + err.(string)
		c.ResultMap["返回数据"] = msg
		c.ResultMap["IsReplyTo"] = "y"
		c.ResultMap["LogToNet"] = "y"
		logLevel = "error"
		logCode = 400
	}
	//
	resultStr, err := json.Marshal(c.ResultMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, ok := c.ResultMap["LogToNet"]; !ok {
		c.ResultMap["LogToNet"] = "y"

	}
	// 日志默认打印网络和本地, LogToNet 不包含“y”  可不打印网络日志
	if strings.Contains(c.ResultMap["LogToNet"], "y") {
		queueInfo := strings.Split(c.QueueObj.RoutingKey, ".")
		processSage := queueInfo[len(queueInfo)-2]
		application := queueInfo[len(queueInfo)-1]
		l.PrintAll(logLevel, c.ResultMap, logCode, processSage, "wiken", application)
	} else {
		l.PrintLocal(strings.Replace(string(resultStr), "\r\n", "", -1), logLevel)
	}
}

func (c *Context) Text() []byte {
	return c.QueueObj.Body
}

func (c *Context) NextTo(exchangeName string, routingKey string, msg []byte) string {
	returnMsg := amqp.Publishing{
		Headers:         nil,
		ContentType:     "application/json",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            msg,
	}

	err := c.Channel.Publish(exchangeName, routingKey, false, false, returnMsg)
	if err != nil {
		fmt.Println("MQ 消息发送失败")
		return "MQ 消息发送失败"
	} else {
		fmt.Println("MQ 消息发送成功")
		return "MQ 消息发送成功"
	}
}

// 封装context 构造函数
func NewContext() *Context {
	c := new(Context)
	c.ResultMap = make(map[string]string)
	return c
}
