package cmd

import (
	"fmt"
	"log"
	"os"
	"ter-tool/cmd/gpt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var rootCmd = &cobra.Command{
	Use:   "wuliu",
	Short: "五六的工作管理工具",
}

func Execute() {
	rootCmd.AddCommand(
		// server.NewServerCmd(),
		gpt.NewGptCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 生成文档文件
func Docs(t string) {
	var err error
	dir, _ := os.Getwd()
	switch t {
	case "yml", "yaml":
		err = doc.GenYamlTree(rootCmd, dir+"/tmp")
	default:
		err = doc.GenMarkdownTree(rootCmd, dir+"/tmp")
	}
	if err != nil {
		log.Fatal(err)
	}
}
