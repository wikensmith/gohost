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

// Ack 实现amqp的ack 方法并封装了日志消停到context
func (c *Context) Ack() {
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
