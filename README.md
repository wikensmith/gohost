## 1. 接收YATP消息
YATP向队列推送消息的时候如使用swagger 的 /api/MQ/PushMessage 方法时,推送消息格式如下:
````
{
  "exchangeName": "",  // 推送队列交换机,可以不写.但一定不能写错
  "routingKey": "YS.机票.国内.询价.wiken.DEBUG",  // 需要推送的目的队列
  "replyTo": "YS.机票.国内.退票查询.wiken.DEBUG",  // 返回消息存储队列名称
  "correlationId": "string",  
  "messageId": "string",
  "messageHeaders": {
    "exchangeName": "system.request",  // 返回消息队列的交换机名, 一定不能写错,写错服务器会挂的
    "additionalProp2": "string",
    "additionalProp3": "string"
  },
  "messageContent": "string",  // 推送消息内容
  "executionTime": "2020-03-24T02:01:26.187Z"  // 时间戳
}
```

### 2. gohost 使用示例

```go
import (
	"github.com/wikensmith/gohost"
	"github.com/wikensmith/gohost/queue"
)

func InitMQ(){
    Params.Prefetch = 1
	Workers["YS.机票.国内.支付.wiken.DEBUG"] = func(context queue.Context) {
		// 日志使用

		context.Log.AddField(0, "Field1_value") // 设置Field 1 的值
		context.Log.AddField(4, "Field4_value") // 设置Field5 的值
		context.Log.Print("test log msg")
		context.Log.Printf("test log msg %s", "just test")
		context.Log.PrintReturn("return msg") // context body 已经传入日志， 不需要额外传入
		context.Log.Level("info")

		//  获取队列body
		body := context.QueueObj.Body // []byte
		context.Set("key1", "aaa")
		aaa, ok := context.Get("key1")
		fmt.Println(aaa)
		if !ok {
			context.Log.Level("error") // 传入日志等级 默认为error
			context.Log.Printf("error in Get, error: [%s]", "该值不存在")
		}

		// 将body传入指定队列
		info := context.NextTo("YS.机票.询价", "YS.机票.国内.询价.wiken.DEBUG", body, nil)
		fmt.Println("info:", info)

		//context.Nack() //
		context.Ack() // Ack 队列
	}
}

func (main){
    gohost.Start()  // 开始项目
}

```

