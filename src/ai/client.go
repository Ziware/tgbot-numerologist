package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type Client struct {
	apiKey     string
	apiURL     string
	proxyURL   string
	httpClient *http.Client
}

var model string = "gpt-4.1"

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}
type Choice struct {
	Message Message `json:"message"`
}

type ErrorResponse struct {
	Error ErrorObject `json:"error"`
}
type ErrorObject struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

var GPTClient *Client

func Init(apiKey, proxyURL string) error {
	dialer, err := proxy.SOCKS5("tcp", proxyURL, nil, proxy.Direct)
	if err != nil {
		return err
	}
	transport := &http.Transport{
		Dial: dialer.Dial,
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	GPTClient = &Client{
		apiKey:     apiKey,
		apiURL:     "https://api.openai.com/v1/chat/completions",
		proxyURL:   proxyURL,
		httpClient: httpClient,
	}

	return nil
}

func (c *Client) SendMessage(messages []Message) (string, error) {
	reqBody := Request{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	var lastErrBody string

	for range 5 {
		req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var errResp ErrorResponse
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErrBody = "unable to parse error body"
				time.Sleep(1)
				continue
			}
			if err := json.Unmarshal([]byte(body), &errResp); err == nil && errResp.Error.Message != "" {
				lastErrBody = fmt.Sprintf("API error, type: %s: %s", errResp.Error.Type, errResp.Error.Message)
				continue
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(1)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("API error status %d, type: %s: %s", resp.StatusCode, errResp.Error.Type, errResp.Error.Message)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var r Response
		if err := json.Unmarshal([]byte(body), &r); err != nil {
			return "", err
		}

		if len(r.Choices) > 0 {
			content := r.Choices[0].Message.Content
			return content, nil
		}
		return "", errors.New("No choices found in json")
	}

	return "", errors.New("Retries exceeded: " + lastErrBody)
}
