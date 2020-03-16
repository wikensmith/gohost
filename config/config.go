package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var config = new(configStruct)

func InitConfig() {

	p := "./config.toml"

	if _, err := toml.DecodeFile(p, config); err != nil {
		fmt.Println("读取配置文件失败, 异常信息:", err.Error())
	}
}

func GetConfig() *configStruct {
	return config

}
