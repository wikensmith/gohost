package gohost

import (
	"fmt"
	"github.com/wikensmith/gohost/queue"
	"testing"
	"time"
)

//type persion struct {
//	Name string
//	Age int
//}

func Log() {
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
		//fmt.Println(time.Now().Format(time.RFC3339), ":", body)
		time.Sleep(time.Second * 2)
		context.Log.PrintInput(string(body))
		//context.Nack() //
		context.Log.Print("adsf", "aa")
		context.Ack(true) // Ack 队列
	}
}

// 测试断线重连和nextTo
func connection() {
	Params.Project = "wikenTest"                         // 日志项目
	Params.Module = "test1"                              // 日志模块
	Params.User = "7921"                                 // 工号
	Params.LogURI = "http://192.168.0.212:8081/log/save" // 日志地址
	Params.IsHealthyCheck = true                         // 是否做健康检测, 默认不做
	Params.Prefetch = 1                                  // 并发数
	Params.IsReConnection = true                         // 是否断线生连、debug的时候，调度时间长了，会被认为连接已经断开
	Workers["YS.机票.国内.支付.wiken.DEBUG"] = func(context queue.Context) {
		defer func() {
			context.Ack(false) // 参数 true: 发送日志至日志中心 ; false: not send  to log center
			//context.Nack(true) // not ack msg and send log to log center
			fmt.Println("acked")
		}()
		body := string(context.QueueObj.Body) // 获取mq消息体
		header := context.QueueObj.Headers    // 获取mq消息头
		fmt.Println(body, header)
		// 发送至新的队列
		info := context.NextTo("system.request", "YS.机票.国内.退票查询.wiken.DEBUG", []byte("testinfo"), header)
		fmt.Println(info)
	}
}

func TestUse(t *testing.T) {
	//Log()
	connection()
	Start()
}
