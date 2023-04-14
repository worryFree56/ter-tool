package gpt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"ter-tool/config"
	"ter-tool/openapi"
	"ter-tool/tool"

	"github.com/spf13/cobra"
)

const (
	apiEndpoint = "https://api.openai.com/v1"
)

var (
	apiKey string
	proxy  string
	//会话存储信息
	msg = openapi.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		// MaxTokens: 4096,
	}
	exitChan chan os.Signal
)

func saveMsgHandler() {
	s := <-exitChan
	fmt.Println("退出信号:", s)
	//写入聊天内容
	tool.WriteLog(tool.PrettyPrint(msg.Messages))
	os.Exit(1)
}

func NewGptCmd() *cobra.Command {

	gptCmd := &cobra.Command{
		Use:     "chatgpt",
		Aliases: []string{"cg"},
		Short:   "chat-gpt 对话客户端",
		Example: fmt.Sprintln("$ wuliu cg -x http://127.0.0.1:7890"),
		RunE:    run,
	}
	gptCmd.Flags().StringVar(&apiKey, "api-key", config.Cfg.Openapi["api-key"], "通过 https://platform.openai.com/account/api-keys 获取的api-Key")
	gptCmd.Flags().StringVarP(&proxy, "proxy", "x", config.Cfg.Openapi["proxy"], "代理地址")
	gptCmd.Flags().Bool("stream", false, "设置true，数据流输出")
	return gptCmd
}

func run(cmd *cobra.Command, args []string) error {
	if apiKey == "" {
		return fmt.Errorf("API Key is required")
	}
	if proxy == "" {
		return fmt.Errorf("proxy is required")
	}
	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	go saveMsgHandler()
	defer func() {
		fmt.Println("写入聊天记录")
		tool.WriteLog(tool.PrettyPrint(msg.Messages))
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// 获取用户输入
		fmt.Print("\n自己:")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if userInput == "" {
			continue
		} else {
			if userInput == "退出" || userInput == "exit" {
				os.Stdin.Close()
				break
			}
		}
		// 调用 GPT-3 API 获取响应
		umsg := openapi.ChatCompletionMessage{
			Role:    openapi.ChatMessageRoleUser,
			Content: userInput,
		}
		msg.Messages = append(msg.Messages, umsg)
		if ok, _ := cmd.Flags().GetBool("stream"); ok {
			msg.Stream = true

			resp, err := msg.SendChatStreamRequest(proxy, apiEndpoint+"/chat/completions", apiKey, msg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				continue
			}
			defer resp.Close()
			fmt.Print("AI: ")
			var echoMsg string
			for {
				resposne, err := resp.Recv()
				if errors.Is(err, io.EOF) {
					// fmt.Println("\nStream finished")
					break
				}
				if err != nil {
					fmt.Printf("\nStream error:1111%v\n", err.Error())
					break
				}
				con := resposne.Choices[0].Delta.Content
				echoMsg += con
				print(con)
			}
			msg.Messages = append(msg.Messages, openapi.ChatCompletionMessage{
				Role:    openapi.ChatMessageRoleAssistant,
				Content: echoMsg,
			})
		} else {
			resp, err := msg.SendChatRequest(apiEndpoint+"/chat/completions", apiKey, msg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				continue
			}
			fmt.Print("AI: ")
			print(resp.Choices[0].Message.Content)
			msg.Messages = append(msg.Messages, openapi.ChatCompletionMessage{
				Role:    openapi.ChatMessageRoleAssistant,
				Content: resp.Choices[0].Message.Content,
			})
		}

	}

	return nil
}
