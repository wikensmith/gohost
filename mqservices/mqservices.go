package mqservices

import "fmt"

const(
	BaseURL = "http://192.168.0.100:5000"
)
// 发送消息到消息队列
func PushMsg()  {
	fmt.Println(" in push msg")
}
// 删除队头的一条消息并返回该条消息
func PopMsg()  {
	fmt.Println(" in push msg")
}

// 向队列发送一条消息
func SendResponse()  {
	fmt.Println(" in push msg")
}

// 向指定队列发送一条,并阻塞等待回复
func SendRequest()  {
	fmt.Println(" in push msg")
}

// 通过 correlationd_id 获取 send_response 的结果
func GetResponse()  {
	fmt.Println(" in push msg")
}
// 发送push_msg 并循环获取数据结果，超时结果, 超时时间不能小于3s
// 使用方法：
// 1、传入参数data为字典类型， correlation_id = None 时， 为先push msg 再循环等待接收内容
// 2、当不传data， 只传入correlation_id 的时候，只进行针对该 id 的循环接收内容操作
func WaitResponse()  {
	fmt.Println(" in push msg")
}

