package gohost

import (
	"fmt"
	"github.com/wikensmith/gohost/queue"
	"github.com/wikensmith/gohost/structs"

	//"github.com/wikensmith/gohost/structs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

var Workers = make(map[string](func(context queue.Context)), 0)

var conn *amqp.Connection

//var Prefetch = 3
//var URI = "amqp://ys:ysmq@192.168.0.100:5672/"
//var Consumer = "goHost" // 队列消费者名称
var Params = &structs.Param{
	Prefetch:    3,
	Consumer:    "gohost",
	MqURI:       "amqp://ys:ysmq@192.168.0.100:5672/", // mq 地址
	Project:     "TestCenter",                         // 日志模块名称
	Module:      "test",
	User:        "7921",
	LogURI:      "http://log.ys.com/log/save", // 日志中心地址
	HealthyPort: "9000",
}

func ReConnection(c *queue.Context) {
	for {
		//isClosed := c.Conn.IsClosed()
		//fmt.Println("是否关闭: ", isClosed)
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
		conn, err := GetConnection()
		if err != nil {
			fmt.Printf("连接失败, 原因: %s\n", err.Error())
		} else {
			c.Delivery, err = GetMsgChan(c.Conn, c.QueueName)
			if err != nil {
				fmt.Printf("连接失败, 原因: %s", err.Error())
				//log.Fatalf("reconnect error , error: [%s]", err.Error())
			} else {
				_ = conn.Close()
				Start()
				return
			}
		}
		fmt.Println("连接中... ...")
		time.Sleep(time.Second)
	}
}

func GetConnection() (conn *amqp.Connection, err error) {
	conn, err = amqp.Dial(Params.MqURI)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func GetMsgChan(conn *amqp.Connection, queueName string) (<-chan amqp.Delivery, error) {
	ch, _ := conn.Channel()
	notifyClose = ch.NotifyClose(errChan)
	if err := ch.Qos(Params.Prefetch, 0, true); err != nil {
		fmt.Println("error in ch.Qos():, error is :", err.Error())
	}
	fmt.Println("队列名称: ", queueName)
	return ch.Consume(
		queueName,
		Params.Consumer,
		false,
		false,
		false,
		false,
		nil)
}

func connect(queueName string, f func(queue.Context)) {
	var err error
	conn, err = GetConnection()
	//errChan = conn.NotifyClose(errChan)
	if err != nil {
		log.Fatalf("error in gohost.connect.GetConnect, error:[%s]", err.Error())
		return
	}
	isConnectClosed = false
	defer conn.Close()
	msgChan, err := GetMsgChan(conn, queueName)
	if err != nil {
		log.Fatalf("error in gohost.connect.GetMsgChan, error:[%s]\n", err.Error())
	}
	context := queue.Context{
		Conn:     conn,
		Delivery: msgChan,
	}
	for msg := range context.Delivery {
		context.Cxt = queue.NewContext()
		context.QueueObj = msg
		context.QueueName = queueName
		context.Connection = conn
		context.Log.PrintInput(string(msg.Body))
		context.Log.Project = Params.Project
		context.Log.Module = Params.Module
		context.Log.User = Params.User
		context.Log.LogURL = Params.LogURI
		go f(context)
	}
}

//func Share(worker1 string, worker2 string) {
//	Workers[worker1] = Workers[worker2]
//}
func forever() chan struct{} {
	ch := make(chan struct{})

	go func() {
		c1 := make(chan os.Signal, 1)
		signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c1
		ch <- struct{}{}
	}()
	return ch
}

func Start() {
	go func() {
		time.Sleep(time.Second * 3)
		if !conn.IsClosed() {
			HealthyCheck()
		}
	}()
	start()
}
func start() {
	for queueName, f := range Workers {
		go connect(queueName, f)
	}
	<-forever()
	fmt.Println("程序结束")

}
