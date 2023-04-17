package main

import (
	"ter-tool/cmd"
	"ter-tool/config"
)

func main() {
	// var randBytes = []byte{116, 101, 115, 108, 97, 119, 117, 108, 105, 117, 46, 33, 126, 64, 51, 50, 54}
	// log.Println([]byte("teslawuliu.!~@326326"))
	// return
	//加载本地配置
	config.NewConfig()

	cmd.Execute()
}
