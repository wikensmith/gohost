package golog

// net log struct
type TraceProperty struct {
	TraceId      string `json:"traceId"`
	ProcessId    string `json:"processId"`
	ProcessStage string `json:"processStage"` // 筛选关键词
}
type ApplicationProperty struct {
	ApplicationName    string `json:"applicationName"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationModule  string `json:"applicationModule"`
	Author             string `json:"author"`
}
type DataProperty struct {
	StatusCode int // 200 成功 非200 失败  0 过程日志
	StatusDesc string
}

type Properties struct {
	TraceProperty       `json:"traceProperty"`
	ApplicationProperty `json:"applicationProperty"`
	DataProperty        `json:"dataProperty"`
}

// 网络日志结构体
type NetLogParam struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	Properties `json:"properties"`
}

// 日志中心配置结构体
type LogCenterSettingStruct struct {
	Project     string `json:"project binding:required"`
	Module      string `json:"module binding:required"`
	Unit        string `json:"unit binding:required"`
	Number      int    `json:"number binding:required"`
	MaxAlarmNum int    `json:"max_alarm_num binding:required"`
	ExtraField  string `json:"extra_field binding:required"`
}
