package account

import (
	"encoding/hex"
	"errors"
	"log"
	"ter-tool/config"
	"ter-tool/tool"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// yaml存储信息前缀
var accountPrefix = "accounts."

// 账号存储
func NewAccountCmd() *cobra.Command {

	accCmd := &cobra.Command{
		Use:   "account",
		Short: "各个平台的账号管理",
	}
	accCmd.PersistentFlags().String("secret", "", "操作账户需要的密钥（必须项）,长度16位")
	accCmd.PersistentFlags().String("list", config.Cfg.GetAccountLists(), "已经存在的账号标识")

	accCmd.AddCommand(
		WriteCmd(),
		UpdateCmd(),
		LookCmd(),
	)
	return accCmd
}

func WriteCmd() *cobra.Command {
	wCmd := &cobra.Command{
		Use:   "write",
		Short: "写入新的平台账号加密数据",
		Run: func(cmd *cobra.Command, args []string) {
			secret, _ := cmd.Flags().GetString("secret")
			if err := ValSecret(secret); err != nil {
				log.Println(err)
				return
			}
			logo, _ := cmd.Flags().GetString("logo")
			info, _ := cmd.Flags().GetString("info")

			if _, ok := config.Cfg.Accounts[logo]; ok {
				log.Println("下标:" + logo + "已经存在，如果修改配置，请使用update方法")
				return
			}
			//进行账户信息加密
			aesInfo := []byte(info)
			aesSecret := []byte(secret)
			s := tool.AesEncryptCFB(aesInfo, aesSecret)

			viper.Set(accountPrefix+logo, hex.EncodeToString(s))
			err := viper.WriteConfig()
			if err != nil {
				log.Println("添加新账号出错：", err)
			}
		},
	}
	wCmd.Flags().String("logo", "", "当前信息的标识，可通过此标识查看原文")
	wCmd.Flags().String("info", "", "账号信息格式：用户名：xxx|密码：xxx|资产密码：xxx|平台：xxxx")
	return wCmd
}

func UpdateCmd() *cobra.Command {
	Ucmd := &cobra.Command{
		Use:   "update",
		Short: "修改平台账号加密数据",
		Run: func(cmd *cobra.Command, args []string) {
			secret, err := cmd.Flags().GetString("secret")
			if err != nil {
				log.Println("获取 secret 失败：", err)
				return
			}
			if secret == "" || len(secret) != 16 {
				log.Println("secret 必填项,切长度为16")
				return
			}
			logo, _ := cmd.Flags().GetString("logo")
			info, _ := cmd.Flags().GetString("info")

			if _, ok := config.Cfg.Accounts[logo]; !ok {
				log.Println("下标:" + logo + "不存在，如果添加配置，请使用write方法")
				return
			}
			//进行账户信息加密
			aesInfo := []byte(info)
			aesSecret := []byte(secret)
			s := tool.AesEncryptCFB(aesInfo, aesSecret)

			viper.Set(accountPrefix+logo, hex.EncodeToString(s))
			err = viper.WriteConfig()
			if err != nil {
				log.Println("修改信息出错：", err)
			}
		},
	}
	Ucmd.Flags().String("logo", "", "当前修改信息的标识，可通过此标识查看原文")
	Ucmd.Flags().String("info", "", "账号信息格式：用户名：xxx|密码：xxx|资产密码：xxx|平台：xxxx")
	return Ucmd
}

func LookCmd() *cobra.Command {
	lCmd := &cobra.Command{
		Use:   "look",
		Short: "查看账号的原文信息",
		Run: func(cmd *cobra.Command, args []string) {
			secret, err := cmd.Flags().GetString("secret")
			if err != nil {
				log.Println("获取 secret 失败：", err)
				return
			}
			if secret == "" || len(secret) != 16 {
				log.Println("secret 必填项,切长度为16")
				return
			}
			logo, _ := cmd.Flags().GetString("logo")
			if _, ok := config.Cfg.Accounts[logo]; !ok {
				log.Println("下标:" + logo + "不存在，无法查看信息")
				return
			}
			info := config.Cfg.Accounts[logo]
			aesInfo, _ := hex.DecodeString(info)
			aesSecret := []byte(secret)
			s := tool.AesDecryptCFB(aesInfo, aesSecret)

			log.Println(string(s))
		},
	}
	lCmd.Flags().String("logo", "", "查看信息的标识，值来源于--lists")
	return lCmd
}

// 验证密钥
func ValSecret(key string) error {
	if key == "" || len(key) != 16 {
		return errors.New("secret必填项且长度为16")
	}
	return nil
}
