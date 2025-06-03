package ai

import (
	"bytes"
	"encoding/json"
	"errors"
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

var model string = "gpt-4o"

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

		if resp.StatusCode == http.StatusTooManyRequests {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				body = make([]byte, 0)
			}
			lastErrBody = string(body)
			time.Sleep(1)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				body = make([]byte, 0)
			}
			return "", errors.New("non-OK HTTP status: " + resp.Status + ", error: " + string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}

	return "", errors.New("Retries exceeded: " + lastErrBody)
}
