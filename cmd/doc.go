package cmd

import (
	"log"
	"os"
	"ter-tool/config"
	"ter-tool/tool"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewDocCmd() *cobra.Command {
	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "生成当前工具的使用文档",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				log.Println(err)
				return
			}
			ft, _ := cmd.Flags().GetString("type")
			switch ft {
			case "md":
				err = doc.GenMarkdownTree(rootCmd, dir)
			default:
				err = doc.GenYamlTree(rootCmd, dir)
			}
			if err != nil {
				log.Println(err)
				return
			}
		},
	}
	dir, _ := os.UserHomeDir()
	docCmd.PersistentFlags().StringP("dir", "d", dir+"/wuliu-doc", "生成的文档所在目录,如果打开文件失败，请先创建wuliu-doc目录")
	docCmd.PersistentFlags().StringP("type", "t", "md", "生成文档文件类型")
	return docCmd
}

func NewYmlCmd() *cobra.Command {
	yCmd := &cobra.Command{
		Use:   "config",
		Short: "yml配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			key, err := cmd.Flags().GetString("secret")
			if err != nil {
				log.Println(err)
				return
			}
			if key == "" || len(key) != 16 {
				log.Println("secret必填项且长度为16")
				return
			}
			log.Println(tool.PrettyPrint(config.Cfg))
		},
	}
	yCmd.PersistentFlags().String("secret", "", "操作需要的密钥（必须项）,长度16位")

	return yCmd
}
