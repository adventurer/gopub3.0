package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Setting struct {
	Title    string
	Database Database
}

type Database struct {
	User     string
	Host     string
	Port     string
	Password string
}

var Variable Setting

func init() {
	if _, err := toml.DecodeFile("./setting.toml", &Variable); err != nil {
		fmt.Println("初始化配置文件出错：检查setting.toml文件", err)
		return
	}
}
