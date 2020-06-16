package gohost

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
	"github.com/wikensmith/gohost/queue"
)

var Workers = make(map[string](func(context queue.Context)), 0)
var Prefetch = 3
var URI = "amqp://ys:ysmq@192.168.0.100:5672/"

func GetConnection() *amqp.Connection {
	conn, err := amqp.Dial(URI)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func connect(queueName string, f func(queue.Context)) {
	var conn *amqp.Connection
	conn = GetConnection()
	defer conn.Close()

	ch, _ := conn.Channel()
	if err := ch.Qos(Prefetch, 0, true); err != nil {
		fmt.Println("error in ch.Qos():, error is :", err.Error())
	}
	fmt.Println("队列名称: ", queueName)
	msgChan, _ := ch.Consume(
		queueName,
		"goHost",
		false,
		false,
		false,
		false,
		nil)
	for msg := range msgChan {
		if conn.IsClosed() {
			conn = GetConnection()
		}
		context := queue.NewContext()
		context.QueueObj = msg
		context.Channel = ch
		context.Connection = conn
		go f(*context)
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
	for queueName, f := range Workers {
		go connect(queueName, f)
	}
	<-forever()
	fmt.Println("程序结束")
}
