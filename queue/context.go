package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
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
	LogMsg map[string]interface{}
}

type Context struct {
	logCenter
	KeysMutex  *sync.RWMutex
	QueueObj   amqp.Delivery
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Services   Services
	log        *golog.Log
	Result     []byte
	StarTime   time.Time              // 程序开始时间
	Keys       map[string]interface{} // 属性
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/
func (c *Context) Set(key string, value interface{}) {
	if c.KeysMutex == nil {
		c.KeysMutex = &sync.RWMutex{}
	}

	c.KeysMutex.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.KeysMutex.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.KeysMutex == nil {
		c.KeysMutex = &sync.RWMutex{}
	}
	c.KeysMutex.RLock()
	value, exists = c.Keys[key]
	c.KeysMutex.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
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
	c.Keys = make(map[string]interface{}, 0)
	c.StarTime = time.Now()
	return &c
}
