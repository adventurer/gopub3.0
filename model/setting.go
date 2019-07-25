package model

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Setting struct {
	Title string
	Git   GitConfig
}

type GitConfig struct {
	Server string
	Port   string
}

var Config Setting

func init() {
	if _, err := toml.DecodeFile("./setting.toml", &Config); err != nil {
		fmt.Println("初始化配置文件出错：检查setting.toml文件", err)
		return
	}
}
