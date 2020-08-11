package toLogCenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)
// 发往日志中心的信息结构体
type LogCenterStruct struct {
	Project string      `json:"project" binding:"required"`
	Module  string      `json:"module" binding:"required"`
	Level   string      `json:"level"`
	User    string      `json:"user"`
	Message interface{} `json:"message"`
	Time    string      `json:"time"`
	Field1  interface{}      `json:"field1"`
	Field2  interface{}     `json:"field2"`
	Field3  interface{}      `json:"field3"`
	Field4  interface{}      `json:"field4"`
	Field5  interface{}      `json:"field5"`
}

type logInfo struct{
	Content interface{}
	Func string
	Time string
}

var URL = "http://192.168.0.212:8081/log/save"

type Logger struct {
	lock *sync.Mutex
	cacheLst [][]interface{} // 日志缓存列表 [[key, *logInfo], [key2, *logInfo]]
	resultMap map[string]interface{}
	Project string // 日志中心项目名称
	Module string  // 日志中心模块名称
	LogURL string // 日志中心地址
	User string // 工号
	level string // 日志等级
	Field []interface{}
}

func (l *Logger)New() *Logger {
	l.lock = new(sync.Mutex)
	l.Field = make([]interface{}, 5)
	l.cacheLst = make([][]interface{}, 0)
	l.resultMap = make(map[string]interface{}, 0)
	if l.LogURL == ""{
		l.LogURL = URL
	}
	return l
}
func (l *Logger)Level(level string)  {
	l.lock.Lock()
	l.level = level
	l.lock.Unlock()
}
func (l *Logger)AddField(index int, val interface{})  {
	if index > 4 {
		return
	}
	l.Field[index] = val
}

func (l *Logger)runFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(4,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func (l *Logger)getMap()  {
	l.resultMap["process"] = l.cacheLst
}
func (l *Logger)PrintInput(content interface{})  {
	l.lock.Lock()
	l.resultMap["input"] = content
	l.lock.Unlock()
}

func (l *Logger)PrintReturn(content interface{})  {
	l.lock.Lock()
	l.resultMap["result"] = content
	l.lock.Unlock()

}

func (l *Logger)print(content interface{})  {
	fmt.Printf("%v\n", content)
	tempStruct := logInfo{
		Content: content,
		Func:    l.runFuncName(),
		Time:    time.Now().Format(time.RFC3339Nano),
	}
	l.lock.Lock()
	l.cacheLst = append(l.cacheLst, []interface{}{tempStruct})
	l.lock.Unlock()
}

func (l *Logger)Printf(format string, a ...interface{})  {
	content := fmt.Sprintf(format, a...)
	l.print(content)
}

func (l * Logger)Print(content interface{})  {
	l.print(content)

}

// 发送日志到退票中心
func (l *Logger)Send()  {
	if !strings.Contains("infoInfoErrorerrorWarnwarnFatalfatal", l.level){
		log.Println("输入日志等级错误, 日志等级仅支持:info,error,fatal,warn")
		return
	}
	l.getMap()
	if len(l.Field) > 5 {
		l.Field = l.Field[:5]
	}else{
		tempLst := make([]interface{}, 5-len(l.Field))
		l.Field = append(l.Field, tempLst...)
	}
	if l.level == ""{
		l.level = "error"
	}
	result, err := json.Marshal(l.resultMap)
	fmt.Print(err)

	msg := LogCenterStruct{
		Project: l.Project,
		Module:  l.Module,
		Level:   l.level,
		User:    l.User,
		Message: string(result),
		Time:    time.Now().Format(time.RFC3339Nano),
		Field1:  l.Field[0],
		Field2:  l.Field[1],
		Field3:  l.Field[2],
		Field4:  l.Field[3],
		Field5:  l.Field[4],
	}
	msgByte, err := json.Marshal(msg)
	fmt.Println(err)
	send(l.LogURL, msgByte)
}

func send(URL string, msg []byte)  {
	resp, err := http.Post(URL, "application/json", bytes.NewReader(msg))
	if err != nil {
		fmt.Println("error in send log post, error: ", err.Error())
	}

	defer func() {
		if resp == nil {
			return
		}
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if resp != nil {
		respStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error in send log post, error: ", err.Error())
		}
		fmt.Println(string(respStr))
	}else {
		fmt.Println("error in toLogCenter, 发送日志响应为nil, 请联系开发者检查日志服务是否正常！ 并检查日志地址是否正确？")
	}
}

