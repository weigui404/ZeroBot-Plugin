package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	// baseURL  = "https://api.openai.com/v1/"
	proxyURL           = "https://api.deepseek.com/v1/"
	modelTurbo = "deepseek-reasoner"
	yunURL             = "https://api.deepseek.com/user/balance"
)

type yun struct {
	Data []struct {
	    Currency  string `json:"currency"`
		Total     string `json:"total_balance"`
		Granted string `json:"granted_balance"`
		Topped  string `json:"topped_up_balance"`
	} `json:"balance_infos"`
}

// chatResponseBody 响应体
type chatResponseBody struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
	Usage   chatUsage    `json:"usage"`
}

// chatRequestBody 请求体
type chatRequestBody struct {
	Model       string        `json:"model,omitempty"`
	Messages    []chatMessage `json:"messages,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	N           int           `json:"n,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// chatMessage 消息
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatChoice struct {
	Index        int `json:"index"`
	Message      chatMessage
	FinishReason string `json:"finish_reason"`
}

type chatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	},
	Timeout: time.Minute * 5,
}

// completions gtp3.5文本模型回复
// curl https://api.openai.com/v1/chat/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer YOUR_API_KEY"
// -d '{ "model": "gpt-3.5-turbo",  "messages": [{"role": "user", "content": "Hello!"}]}'
func completions(messages []chatMessage, apiKey string) (*chatResponseBody, error) {
	com := chatRequestBody{
		Messages: messages,
	}
	// default model
	if com.Model == "" {
		com.Model = modelTurbo
	}

	body, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, proxyURL+"chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("response error: " + strconv.Itoa(res.StatusCode))
	}

	v := new(chatResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
