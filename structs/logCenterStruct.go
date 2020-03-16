package structs

// 发往日志中心的信息结构体
type LogCenterStruct struct {
	Project string `json:"project binding: required"`
	Module  string `json:"module binding: required"`
	User    string
	Message string
	Time    string
	field1  string
	field2  string
	field3  string
	field4  string
	field5  string
}
