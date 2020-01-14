package gohost

import (
	"fmt"
	"github.com/wikensmith/gohost/queue"
	"testing"
)

func init() {
	Workers["YS.机票.国内.询价.wiken.DEBUG"] = func(context queue.Context) string {
		fmt.Println(string(context.QueueObj.Body))
		return "end"
	}
}

func TestUse(t *testing.T) {
	Start()

}
