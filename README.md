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
    // 测试断线重连和nextTo
    func connection() {
    	Params.Project = "wikenTest"
    	Params.Module = "test1"
    	Params.User = "7921"
    	Params.IsHealthyCheck = true // 是否做健康检测, 默认不做
    	Params.Prefetch = 1
    	Params.IsReConnection = true // 是否断线生连、debug的时候，调度时间长了，会被认为连接已经断开
    	Workers["YS.机票.国内.支付.wiken.DEBUG"] = func(context queue.Context) {
    		defer func() {
    			context.Ack(false)
    			fmt.Println("acked")
    		}()
    		body := string(context.QueueObj.Body)
    		fmt.Println(body)
    		info := context.NextTo("system.request", "YS.机票.国内.退票查询.wiken.DEBUG", []byte("testinfo"), nil)
    		fmt.Println(info)
    	}
    }
}

func (main){
    gohost.Start()  // 开始项目
}

```

