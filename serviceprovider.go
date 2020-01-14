package gohost

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wikensmith/gohost/queue"
	"github.com/streadway/amqp"
)

var Workers = make(map[string](func(context queue.Context) string), 2)

func connect(queueName string, f func(queue.Context) string) {
	conn, err := amqp.Dial("amqp://ys:ysmq@192.168.0.100:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, _ := conn.Channel()
	ch.Qos(3, 0, true)
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
		context := queue.Context{QueueObj: msg, Channel: ch}
		go f(context)
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
		connect(queueName, f)
	}
	fmt.Println("程序结束")
	<-forever()
}
