package gohost

import (
	"fmt"
	"github.com/wikensmith/gohost/queue"
	"testing"
)

func init() {
	Workers["YS.机票.国内.退票查询.wiken.DEBUG"] = func(context queue.Context) string {
		fmt.Println(string(context.QueueObj.Body))
		body := string(context.QueueObj.Body)
		log := context.Services.NewLogger()
		//log.PrintLocal(body, "INFo")

		var m = make(map[string]string)
		m["队列名称"] = "sss"
		m["传入数据"] = "input data"
		m["返回数据"] = body
		m["请求参数"] = "requests data"
		m["响应参数"] = "response data"
		log.PrintAll(
			"INFO",
			m,
			200,
			"退票",
			"wiken",
			"test-wiken",)
		return "end"
	}
}

func TestUse(t *testing.T) {
	Start()

}
