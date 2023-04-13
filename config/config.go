package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	ServerEnable = "server"
)

type Config struct {
	Servers map[string]Server
}

type Server struct {
	IP         string
	Name       string
	User       string
	Pwd        string
	NodeDir    string //节点目录
	Desc       string //描述
	Bt         string //宝塔面板
	ValAddress string //验证者地址
}

var (
	//配置
	Cfg Config
)

func init() {
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic("载入本地配置失败：" + err.Error())
	}
}

func NewConfig() {
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c Config) GetServer() map[string]Server {
	return c.Servers
}

func (c Config) GetServerNames() string {
	var values = "(|"
	for k := range c.GetServer() {
		values += k + "|"
	}
	values += ")"
	return values
}
