package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChoiceEntry struct {
	Index   int           `json:"index"`
	Message OpenAIMessage `json:"message"`
}

type OpenAIResp struct {
	Model   string              `json:"model"`
	Choices []OpenAIChoiceEntry `json:"choices"`
}

// This is not a request to OpenAI api, it's a request made from our frontend
// to the backend server
type OpenAIRequest struct {
	OpenaiQuestion string `json:"openai-question"`
}

type Client struct {
	openAIApiKey string
	*http.Client
}

func NewOpenAIClient() (*Client, error) {
	openAIApiKey, set := os.LookupEnv("OPENAI_API_KEY")
	if !set || openAIApiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	return &Client{
		openAIApiKey: openAIApiKey,
		Client:       &http.Client{},
	}, nil
}

// TODO: This should be rewritten in a more understandable way
// And the function should be renamed
func (c *Client) AskOpenAI(ctx context.Context, message string) (*OpenAIResp, error) {
	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
		{
			"role":    "user",
			"content": message,
		},
	}

	reqData := map[string]interface{}{
		"model":    "gpt-4o-mini-2024-07-18",
		"messages": messages,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to create a request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.openAIApiKey))

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %v", err)
	}

	defer resp.Body.Close()

	// TODO: Read API documentation for possible error codes
	// if resp.StatusCode != http.StatusOK {
	// 	// log.Fatalf("Response status code: %d, message: %s", resp.StatusCode, resp.Status)
	// }

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response body: %v", err)
	}

	var openAIResp OpenAIResp
	err = json.Unmarshal(respBytes, &openAIResp)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal the response body: %v", err)
	}

	return &openAIResp, nil
}
