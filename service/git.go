package service

import (
	"strings"

	"gopub3.0/cmd"
)

func GitApi() {
	// 	#git初始化
	// git init
	// #设置remote地址
	// git remote add  origin 地址
	// 得到干净的文件列表
	// git log -1 --name-only --pretty=format:'' 版本号
}

func GetBranchs(directory string) (result []string) {
	return
}

func GetVersions(directory string) (result []string, err error) {
	output, err := cmd.RunLocal("cd " + directory + " && " + "git log -20 --pretty=\"%h - %an - %s - %cD\"")
	if err != nil {
		return
	}
	ouputArr := strings.Split(output, "\n")
	return ouputArr, nil
}

func GetVersionInfo(directory string, hash string) (result string, err error) {
	output, err := cmd.RunLocal("cd " + directory + " && " + "git log --stat -1 " + hash)
	if err != nil {
		return
	}
	return output, nil
}
