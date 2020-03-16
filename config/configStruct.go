package config

type configStruct struct {
	LogCenterSetting logCenterSetting
}
type logCenterSetting struct {
	Url         string `toml:"url"`
	Unit        string `toml:"unit"`
	Number      string `toml:"number"`
	MaxAlarmNum string `toml:"max_alarm_num"`
	ExtraField  string `toml:"extra_field"`
}
