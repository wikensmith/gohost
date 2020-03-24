package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wikensmith/gohost"

	"github.com/wikensmith/gohost/structs"

	"github.com/streadway/amqp"
	"github.com/wikensmith/gohost/golog"
)

// struct for context in callable function that defined by yourself.

var logCenterUrl = "http://192.168.0.212:8081/log/save"

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
	log       *golog.Log
}

func (c *Context) Ack(msg []byte) {
	replyTo := ""
	if c.QueueObj.Headers["ReplyTo"] != nil {
		replyTo = c.QueueObj.Headers["ReplyTo"].(string)
	}

	if replyTo != "" {
		info := c.NextTo(gohost.ReplyToExchangeName, replyTo, msg)
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

// 保存日志至本地文件
func (c *Context) LocalLog(msg, level string) {
	c.log.PrintLocal(msg, level)
}
func getTogether(a []byte, b []byte) []byte {
	for _, v := range b {
		a = append(a, v)
	}
	return a
}

// 保存日志至日志中心
func (c *Context) LogCenter(msg *structs.LogCenterStruct) {

	msgByte, err := json.Marshal(msg)
	if err != nil {
		c.log.PrintLocal("传入日志中心日志序列化异常:  日志信息: "+"异常信息: "+err.Error(), "error")
	}

	inputMsg := c.QueueObj.Body
	logMsg := getTogether(inputMsg, msgByte)

	resp, err := http.Post(logCenterUrl, "application/json", bytes.NewReader(logMsg))
	if err != nil {
		c.log.PrintLocal("传入日志中心异常:  日志信息: "+string(logMsg)+"异常信息: "+err.Error(), "error")
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
	c := Context{
		log: new(golog.Log).New(),
	}
	c.ResultMap = make(map[string]string)
	return &c
}

// 日志中心配置
func logCenterSetting(msg *structs.LogCenterStruct) {
	settings := golog.LogCenterSettingStruct{
		Project:     msg.Project,
		Module:      msg.Module,
		Unit:        "m",
		Number:      1,
		MaxAlarmNum: 100,
		ExtraField:  "",
	}
	http.Post("")

}
