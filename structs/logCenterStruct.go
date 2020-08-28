package structs

// 发往日志中心的信息结构体
type LogCenterStruct struct {
	Project string `json:"project" binding:"required"`
	Module  string `json:"module" binding:"required"`
	Level   string `json:"level"`
	User    string `json:"user"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Field1  string `json:"field1"`
	Field2  string `json:"field2"`
	Field3  string `json:"field3"`
	Field4  string `json:"field4"`
	Field5  string `json:"field5"`
}

type Param struct {
	Prefetch    int    //
	Consumer    string //
	MqURI       string //
	Project     string
	Module      string
	User        string
	LogURI      string
	HealthyPort string // 健康检查端口地址
}
