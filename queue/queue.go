package queue

import "github.com/streadway/amqp"

type Queue struct {
	RoutingKey   string
	ExchangeName string
	QueueType    string
	Prefetch     int
	MsgChan      <-chan amqp.Delivery
}
