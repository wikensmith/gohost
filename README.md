1. 接收YATP消息
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
2. gohost 使用示例
```go
import (
	"github.com/wikensmith/gohost"
	"github.com/wikensmith/gohost/queue"
)
func myFunc(c *queue.Context){    
    defer func(){}{
        // 返回结果类型为interface{}, 
        // 写入数据格式有
        // c.LogMsg["传入数据"] = "" 
        // c.LogMsg["队列名称"] = "" 这两个值在gohost中已经内部赋值 
        c.LogMsg["返回数据"] = make(map[string]string{
            "message": "return_message
        })
        c.Level = "info"                
    }
    
    
}


func InitMQ(){
    gohost.Prefetch = c.GetConfig().RabbitMq.Prefetch
    gohost.Workers[c.GetConfig().RabbitMq.QueueName] = func(context queue.Context) string{  
        // 显示调用写入日志, 因为project, module等需要按项目需要传入
        defer func() {
			msgB, _ := json.Marshal(context.LogMsg)
			// 发送消息到日志中心
			context.LogCenter(&HostStructs.LogCenterStruct{
				Project: "wikenTest",  // 日志项目名
				Module:  "test",  // 日志模块名
				User:    "7921",  // 工号
				Level:   context.Level,  // 日志等级
				Message: string(msgB),  
				Time:    time.Now().Format("2006-01-02T15:04:05+08:00"),  // 传入本地时间
                Field1:  context.Field1, // 示例:平台订单号
                Field2:  context.Field2, // 示例:是否自愿}
                Field3:  context.Field3, // 示例:平台退票单号
                Field4:  context.Field4, // 示例:渠道名称
                Field5:  context.Field5, // 示例:队列名称
			})
		}()
        myfunc(&context)
        
        if xx {
            context.ack("xxx异常", false)
        }
        context.Ack(nil,true)
        return ""
    }
    
    

}
func (main){
    gohost.Start()  // 开始项目
}

```

