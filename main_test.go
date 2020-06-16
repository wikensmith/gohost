package gohost

import (
	"fmt"
	"testing"

	"github.com/wikensmith/gohost/queue"
)

func init() {
	Prefetch = 1

	Workers["YS.机票.国内.询价.wiken.DEBUG"] = func(context queue.Context) {
		//context.Ack([]byte("test_result"))
		//context.NextTo("system.request", "YS.机票.国内.退票.wiken.DEBUG", []byte("test_result"))
		msg := context.NextTo("system.reques", "YS.机票.国内.退票.wiken.DEBUG", []byte("aaa"), nil)
		fmt.Println(msg)

		context.Ack()

		if context.Connection.IsClosed() {
			fmt.Println("sdfsf")
		}
		fmt.Println("here:")
	}
}

func TestUse(t *testing.T) {
	Start()

}
