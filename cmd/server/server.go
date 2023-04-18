package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"ter-tool/config"
	"ter-tool/tool"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// var serverPrefix = "servers."

func NewServerCmd() *cobra.Command {

	serCmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"sv", "s"},
		Short:   "服务器管理",
	}
	serCmd.AddCommand(
		NewListCmd(),
		NewLoginCmd(),
	)

	return serCmd
}
func NewLoginCmd() *cobra.Command {

	var loginCmd = &cobra.Command{
		Use:     "login",
		Aliases: []string{"l"},
		Short:   "登录给定信息服务器ssh",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serS, err := cmd.Flags().GetString(config.ServerEnable)
			if err != nil {
				fmt.Println("服务器列表加载失败")
				return
			}
			if !strings.Contains(serS, "|"+args[0]+"|") {
				fmt.Println("参数服务器不存在列表里面")
				return
			}
			sshInfo := config.Cfg.GetServer()[args[0]]
			cli, err := ssh.Dial("tcp", sshInfo.IP, &ssh.ClientConfig{
				User: sshInfo.User,
				Auth: []ssh.AuthMethod{ssh.Password(sshInfo.Pwd)},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			defer cli.Close()
			session, err := cli.NewSession()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer session.Close()
			fd := int(os.Stdin.Fd())
			oldState, err := term.MakeRaw(fd)
			if err != nil {
				fmt.Println("创建文件描述出错", err)
				return
			}
			defer term.Restore(fd, oldState)
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			session.Stdin = os.Stdin
			modes := ssh.TerminalModes{
				ssh.ECHO:          1,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}
			termType := os.Getenv("TERM")
			if termType == "" {
				termType = "xterm-256color"
			}
			if err := session.RequestPty(termType, 100, 100, modes); err != nil {
				fmt.Println("创建终端失败", err.Error())
				return
			}
			listenWindowChange(fd, session)

			//启动远程shell
			if err = session.Shell(); err != nil {
				fmt.Println(err)
				return
			}

			//等待远程命令（终端）退出
			if err = session.Wait(); err != nil {
				fmt.Println(err)
				return
			}

		},
	}
	loginCmd.Flags().StringP(config.ServerEnable, "s", config.Cfg.GetServerNames(), "可选地服务器项目名称")
	return loginCmd
}

func NewListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "lists",
		Aliases: []string{"ls"},
		Short:   "展示保存的服务器列表",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(tool.PrettyPrint(config.Cfg.GetServer()))
		},
	}
	return listCmd
}

func listenWindowChange(fd int, session *ssh.Session) {
	// 获取窗口大小
	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Println("获取窗口大小出错：", err.Error())
		return
	}
	// 发送大小信息到 SSH 会话的伪终端中
	session.WindowChange(height, width)
}
