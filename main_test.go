package gohost

import (
	"testing"

	"github.com/wikensmith/gohost/queue"
)

func init() {
	Workers["YS.机票.国内.退票查询.wiken.DEBUG"] = func(context queue.Context) string {
		log := context.Services.NewLogger()
		defer context.Defer(log)
		context.ResultMap["返回数据"] = "test_value"
		context.ResultMap["IsReplyTo"] = "y"

		return "end"
	}
}

func TestUse(t *testing.T) {
	Start()

}
