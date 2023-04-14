package openapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content"`
}

type ChatCompletionStreamChoice struct {
	Index        int                             `json:"index"`
	Delta        ChatCompletionStreamChoiceDelta `json:"delta"`
	FinishReason string                          `json:"finish_reason"`
}

type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

// 发送请求
func (chatReq *ChatCompletionRequest) SendChatStreamRequest(proxy, endpoint string, apiKey string, body interface{}) (stream *StreamReader, err error) {
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
