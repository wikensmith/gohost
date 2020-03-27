package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/streadway/amqp"
	"github.com/wikensmith/gohost/golog"
	"github.com/wikensmith/gohost/structs"
)

// struct for context in callable function that defined by yourself.

var logCenterUrl = "http://192.168.0.212:8081/log/save"

type Services struct {
}

func (s *Services) NewLogger() *golog.Log {
	return new(golog.Log).New()
}

type logCenter struct {
	Level  string
	LogMsg map[string]interface{}
	Field1 string
	Field2 string
	Field3 string
	Field4 string
	Field5 string
}

type Context struct {
	logCenter
	QueueObj amqp.Delivery
	Channel  *amqp.Channel
	Services Services
	log      *golog.Log
	Result   []byte
	StarTime time.Time // 程序开始时间
}

// Ack 实现amqp的ack 方法并封装了日志消停到context,
// isReply: bool, true 返回msg至ReplyTo的消息队列中,false 不返回任何消息(其他程序会返回消息到该队列)
func (c *Context) Ack(level string, msg []byte, IsReply bool, headers map[string]interface{}) {
	// 设置日志等级及消息
	m := make(map[string]interface{})
	_ = json.Unmarshal(c.QueueObj.Body, &m)

	c.Level = level
	c.LogMsg["b队列名称"] = c.QueueObj.RoutingKey
	c.LogMsg["c传入数据"] = m

	if msg != nil {
		c.LogMsg["a返回数据"] = string(msg)
		c.Level = "error"
	} else {
		msg = c.Result
		c.Level = "info"
	}

	// 如果有replyTo 和 并且需要返回消息,调用NextTo
	replyTo := c.QueueObj.Headers["replyTo"]
	if replyTo != nil && IsReply == true {
		ExchangeName := c.QueueObj.Headers["exchangeName"].(string)
		info := c.NextTo(ExchangeName, replyTo.(string), msg, headers)
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

// 保存日志至本地文件
func (c *Context) LocalLog(msg, level string) {
	c.log.PrintLocal(msg, level)
}

// 保存日志至日志中心
func (c *Context) LogCenter(msg *structs.LogCenterStruct) {
	msgByte, _ := json.Marshal(msg)
	resp, err := http.Post(logCenterUrl, "application/json", bytes.NewReader(msgByte))
	if err != nil {
		c.log.PrintLocal("传入日志中心异常:  日志信息: "+string(msgByte)+"异常信息: "+err.Error(), "error")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	respStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.log.PrintLocal("传入日志中心异常:  日志信息: "+string(respStr)+"异常信息: "+err.Error(), "error")
	}
}

func (c *Context) Text() []byte {
	return c.QueueObj.Body
}

func (c *Context) NextTo(exchangeName string, routingKey string, msg []byte, headers map[string]interface{}) string {
	returnMsg := amqp.Publishing{
		Headers:         headers,
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
		//fmt.Println("MQ 消息发送失败")
		return "MQ 消息发送失败"
	} else {
		//fmt.Println("MQ 消息发送成功")
		return "MQ 消息发送成功"
	}
}

// GetElapsedTime 获取耗时
func (c *Context) GetElapsedTime() int64 {
	return time.Now().Sub(c.StarTime).Microseconds()
}

// 封装context 构造函数
func NewContext() *Context {
	c := Context{
		log: new(golog.Log).New(),
	}
	c.LogMsg = make(map[string]interface{})
	c.StarTime = time.Now()
	return &c
}
