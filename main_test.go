package gohost

import (
	"fmt"
	"testing"

	"github.com/wikensmith/gohost/queue"
)

func init() {
	Prefetch = 1

	Workers["YS.机票.国内.询价.wiken.DEBUG"] = func(context queue.Context) string {
		//context.Ack([]byte("test_result"))
		//context.NextTo("system.request", "YS.机票.国内.退票.wiken.DEBUG", []byte("test_result"))
		context.Ack(nil)
		fmt.Println("here")
		context.LocalLog("test_reslaaaaa", "error")
		return "end"
	}
}

func TestUse(t *testing.T) {
	Start()

}
