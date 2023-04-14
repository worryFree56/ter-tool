package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"ter-tool/config"
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
	msg = ChatCompletionRequest{
		Model:  "gpt-3.5-turbo",
		Stream: true,
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
		umsg := ChatCompletionMessage{
			Role:    ChatMessageRoleUser,
			Content: userInput,
		}
		msg.Messages = append(msg.Messages, umsg)
		stream, err := sendAPIRequest(proxy, apiEndpoint+"/chat/completions", apiKey, msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}
		defer stream.response.Body.Close()
		fmt.Print("AI: ")
		var echoMsg string
		for {
			resposne, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// fmt.Println("\nStream finished")
				break
			}
			if err != nil {
				fmt.Printf("\nStream error:%v\n", err)
				break
			}
			con := resposne.Choices[0].Delta.Content
			echoMsg += con
			print(con)
		}
		msg.Messages = append(msg.Messages, ChatCompletionMessage{
			Role:    ChatMessageRoleAssistant,
			Content: echoMsg,
		})
	}

	return nil
}

// 发送请求
func sendAPIRequest(proxy, endpoint string, apiKey string, body interface{}) (stream *StreamReader, err error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	//stream = true
	req.Header.Set("Accept", "text/event-stream")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	stream = &StreamReader{
		reader:   bufio.NewReader(resp.Body),
		response: resp,
		buffer:   &bytes.Buffer{},
	}
	return
}
