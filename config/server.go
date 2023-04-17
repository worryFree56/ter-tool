package config

// wuliu server 模块使用模型
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
