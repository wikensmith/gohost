package golog

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := new(Log).New()
	var m = make(map[string]string)
	m["队列名称"] = "sss"
	m["传入数据"] = "input data"
	m["返回数据"] = "return data"
	m["请求参数"] = "requests data"
	m["响应参数"] = "response data"

	l.PrintAll(
		"INFO",
		m,
		200,
		"退票",
		"wiken",
		"test-wiken",
	)
	//a :=  `{"timestamp":"","level":"INFO","message":"{\"传入数据\":\"input data\",\"响应参数\":\"response data\",\"请求参数\":\"requests data\",\"返回数据\":\"return data\",\"队列名称\":\"sss\"}","property":{"traceId":"f7cd6b55-3aa4-11ea-a0ed-02004c4f4f50","processedId":"f7cd6b55-3aa4-11ea-a0ee-02004c4f4f50","processStage":"退票"},"applicationProperty":{"applicationName":"test-wiken","applicationVersion":"","applicatioModule":"","author":"wiken"},"dataProperty":{"StatusCode":200,"StatusDesc":"{\"传入数据\":\"input data\",\"响应参数\":\"response data\",\"请求参数\":\"requests data\",\"返回数据\":\"return data\",\"队列名称\":\"sss\"}"}}`
	//res, _ := http.Post("http://192.168.0.100:5000/api/LogCenter/NewLog", "text/plain", strings.NewReader(a))
	//defer fmt.Println(res.Body.Close())
	//fmt.Println(res.Body)

}
