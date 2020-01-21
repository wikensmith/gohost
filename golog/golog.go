package golog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 网络日志发往的地址
const URL = "http://192.168.0.100:5000/api/LogCenter/NewLog"

func init() {
	err := os.MkdirAll("./logs/error", 0666)
	if err != nil {
		fmt.Println("err in golog init", err)
	}
	err = os.MkdirAll("./logs/info", 0666)
	if err != nil {
		fmt.Println("err in golog init", err)
	}
}

// 日志打印功能
//使用方法：
// l := context.services.NewLogger
// l.PrintLocal(msg, level)  // 只打印本地日志
// var m = make(map[string]string)
//		m["队列名称"] = "sss"  // 该字段不会自动加， 需要手动加入
//		m["传入数据"] = "input data"
//		m["返回数据"] = body  // 这个字段必须有，是打入网络日志的内容
//		m["请求参数"] = "requests data"
//		m["响应参数"] = "response data"
//		log.PrintAll(
//			"INFO",
//			m,
//			200,
//			"退票",
//			"wiken",
//			"test-wiken",)
type Log struct {
	Logger *log.Logger
	Mu     sync.Mutex
}

func (l *Log) New() *Log {
	l.Logger = log.New(new(MyWriter).New(), "", 7)
	return l
}

//level: INFO/ERROR
// msg: 消息内容
// Code: 状态码  200 成功  非200失败， 0 过程日志
// msg: 消息主体
// processStage： 日志筛选关键字
// applicationName: 程序名称
// 功能： 同时将日志定出本地和日志中心， 本地文件命名为日期.log
func (l *Log) PrintAll(level string, msgMap map[string]string, Code int, processStage string,
	auth string, applicationName string) {
	msg, err := json.Marshal(msgMap)
	if err != nil {
		fmt.Println("error in log.Print for json.Marshal")
		return
	}
	data := NetLogParam{
		Timestamp: time.Unix(time.Now().Unix(), 0).Format("2016-01-02T15:04:05"),
		Level:     level,
		Message:   string(msg),
		Properties: Properties{
			TraceProperty: TraceProperty{
				TraceId:      uuid.NewV1().String(),
				ProcessId:    uuid.NewV1().String(),
				ProcessStage: processStage,
			},
			ApplicationProperty: ApplicationProperty{
				ApplicationName:    applicationName,
				ApplicationVersion: "",
				ApplicationModule:  "",
				Author:             auth,
			},
			DataProperty: DataProperty{
				StatusCode: Code,
				StatusDesc: string(msg),
			},
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error in myLog print:", err)
	}
	// 写入日志
	l.Logger.Print(string(b))
}
func (l *Log) Panic(v interface{}) {
	l.Logger.Panic(v)
}

// 只打印本地日志
func (l *Log) PrintLocal(msg string, level string) {
	var file *os.File
	var err error

	msg = "\n" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02T15:04:05") +
		strings.Replace(msg, "\r\n", "", -1)
	msg = strings.Replace(msg, " ", "", -1)

	l.Mu.Lock()
	defer l.Mu.Unlock()

	if strings.ToUpper(level) == "INFO" {
		file, err = os.OpenFile("./logs/info/"+time.Unix(time.Now().Unix(), 0).Format("2006-01-02")+".log",
			os.O_CREATE|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile("./logs/error/"+time.Unix(time.Now().Unix(), 0).Format("2006-01-02")+".log",
			os.O_CREATE|os.O_APPEND, 0666)
	}
	if err != nil {
		fmt.Println("err when open file:", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("error in close file:", err)
		}
	}()

	n, err := file.Write([]byte(msg))
	if err != nil {
		fmt.Println("err when open file:", err)
	}

	fmt.Printf("本地写入成功，写入%d个字节", n)
}

type MyWriter struct {
	Mu sync.Mutex
}

// 日志写入本地
func (m *MyWriter) WriteToFile(p []byte) {
	var fileName string

	data := new(NetLogParam)
	err := json.Unmarshal(p[27:], &data)
	if err != nil {
		fmt.Println("error in Unmarshal data for write to file:", err)
		return
	}

	if strings.ToUpper(data.Level) == "INFO" {
		fileName = "./logs/info/" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02") + ".log"
	} else {
		fileName = "./logs/error/" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02") + ".log"
	}

	m.Mu.Lock()
	defer m.Mu.Unlock()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error in Write log to File:", err)
	}
	defer file.Close()

	d := "\n" + string(p[:27]) + data.Message
	d = strings.Replace(d, "\r\n", "", -1)
	d = strings.Replace(d, " ", "", -1)

	n, err := file.Write([]byte(d))
	if err != nil {
		fmt.Println("error in Write log to File:", err)
	}

	fmt.Printf("本地日志写入成功， 写入字节数%d", n)
}

// 日志写入网络日志中心， msg 为return msg
func (m *MyWriter) WriteToNet(p []byte) {
	data := new(NetLogParam)

	err := json.Unmarshal([]byte(string(p[27:])), &data)
	if err != nil {
		fmt.Println("error in log .writeToNet:", err)
		return
	}
	msg := data.Message
	var msgMap = make(map[string]string)
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		fmt.Println("error in WriteToFile:", err)
		return
	}
	data.Message = msgMap["返回数据"]
	data.Properties.DataProperty.StatusDesc = msgMap["返回数据"]
	p, err = json.Marshal(data)
	if err != nil {
		fmt.Println("error in log .writeToNet 2:", err)
		return
	}
	res, err := http.Post(URL, "application/json", bytes.NewReader(p))
	if err != nil {
		fmt.Println("error in send log to net:", err)
	}
	defer res.Body.Close()
	fmt.Println("日志写入网络成功")
}

func (m *MyWriter) Write(p []byte) (n int, err error) {
	go m.WriteToFile(p)
	go m.WriteToNet(p)
	time.Sleep(time.Second * 1)
	return 0, nil
}

func (m *MyWriter) New() *MyWriter {
	return &MyWriter{}
}
