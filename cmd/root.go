package cmd

import (
	"fmt"
	"os"
	"ter-tool/cmd/account"
	"ter-tool/cmd/gpt"
	"ter-tool/cmd/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wuliu",
	Short: "五六的工作管理工具",
}

func Execute() {

	rootCmd.AddCommand(
		server.NewServerCmd(),
		gpt.NewGptCmd(),
		NewDocCmd(),
		account.NewAccountCmd(),
		NewYmlCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
