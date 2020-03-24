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

var Workers = make(map[string](func(context queue.Context) string), 2)
var Prefetch = 3
var URI = "amqp://ys:ysmq@192.168.0.100:5672/"

// context ack 时 replyto 时候使用的交换机名 header
var ReplyToExchangeName = "system.request"

func connect(queueName string, f func(queue.Context) string) {
	conn, err := amqp.Dial(URI)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, _ := conn.Channel()
	err = ch.Qos(Prefetch, 0, true)
	if err != nil {
		fmt.Println("err in ch.Qos", err)
	}
	defer ch.Close()
	fmt.Println(queueName)
	msgChan, _ := ch.Consume(
		queueName,
		"goHost",
		false,
		false,
		false,
		false,
		nil)
	for msg := range msgChan {
		context := queue.NewContext()
		context.QueueObj = msg
		context.Channel = ch
		go f(*context)
	}
}

func Share(worker1 string, worker2 string) {
	Workers[worker1] = Workers[worker2]
}
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
		fmt.Println("sadsafgdsffgds")
		go connect(queueName, f)
	}
	<-forever()
	fmt.Println("程序结束")
}
