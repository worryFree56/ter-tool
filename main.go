package main

import (
	"ter-tool/cmd"
	"ter-tool/config"
)

func main() {
	//加载本地配置
	config.NewConfig()

	cmd.Execute()
}
