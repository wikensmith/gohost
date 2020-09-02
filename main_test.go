package gohost

import (
	"fmt"
	"github.com/wikensmith/gohost/queue"
	"testing"
	"time"
)

func init() {
	Params.Prefetch = 1
	Workers["YS.机票.国内.支付.wiken.DEBUG"] = func(context queue.Context) {

		// 日志使用

		//context.Log.AddField(0, "Field1_value") // 设置Field 1 的值
		//context.Log.AddField(4, "Field4_value") // 设置Field5 的值
		//context.Log.Print("test log msg")
		//context.Log.Printf("test log msg %s", "just test")
		//context.Log.PrintReturn("return msg") // context body 已经传入日志， 不需要额外传入
		//context.Log.Level("info")
		//
		////  获取队列body
		body := context.QueueObj.Body // []byte
		//context.Set("key1", "aaa")
		//aaa, ok := context.Get("key1")
		//fmt.Println(aaa)
		//if !ok {
		//	context.Log.Level("error") // 传入日志等级 默认为error
		//	context.Log.Printf("error in Get, error: [%s]", "该值不存在")
		//}
		//
		//// 将body传入指定队列
		//info := context.NextTo("YS.机票.询价", "YS.机票.国内.询价.wiken.DEBUG", body, nil)
		fmt.Println(time.Now().Format(time.RFC3339), ":", string(body))
		time.Sleep(time.Second * 2)
		context.Nack() //
		//context.Ack()  // Ack 队列
	}
}

func TestUse(t *testing.T) {
	Start()

}
