package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ServerEnable = "server"
)

type Config struct {
	Servers  map[string]Server
	Openapi  map[string]string
	Accounts map[string]string
}

var (
	//配置
	Cfg Config
)

func NewConfig() {

	//TODO: 设置配置文件目录
	viper.SetConfigName(".wuliu")
	viper.SetConfigType("yml")
	// viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//配置不存在，在当前用户目录创建配置
			dir, _ := os.UserHomeDir()
			if err := createCfgFile(dir + "/.wuliu.yml"); err != nil {
				panic("创建默认文件失败: " + err.Error())
			}
		} else {
			panic("载入本地配置失败：" + err.Error())
		}
	}

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

func (c Config) GetAccountLists() string {
	if len(c.Accounts) < 1 {
		return "()"
	}
	var values = "(|"
	for k := range c.Accounts {
		values += k + "|"
	}
	values += ")"
	return values
}

func createCfgFile(dir string) error {
	fs, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer fs.Close()

	data, err := yaml.Marshal(DefaultConfig())
	if err != nil {
		return err
	}
	return os.WriteFile(dir, data, 0666)
}

func DefaultConfig() Config {
	return Config{
		Servers: map[string]Server{
			"full-node": Server{
				IP:      "0.0.0.0",
				Name:    "test",
				Desc:    "test desc",
				User:    "root",
				Pwd:     "123456",
				Bt:      "http://",
				NodeDir: "/data",
			},
		},
		Openapi: map[string]string{
			"api-key": "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",
			"proxy":   "http://127.0.0.1:7890",
		},
		Accounts: map[string]string{
			"qq": "sdfasdfasdflflfnfnfnnn",
		},
	}
}
