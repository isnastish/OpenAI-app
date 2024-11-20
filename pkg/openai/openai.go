package openai

import (
	"bytes"
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
	OpenAIAPIKey string
	*http.Client
}

func NewOpenAIClient() (*Client, error) {
	OPENAI_API_KEY, set := os.LookupEnv("OPENAI_API_KEY")
	if set == false || OPENAI_API_KEY == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	return &Client{
		OpenAIAPIKey: OPENAI_API_KEY,
		Client:       &http.Client{},
	}, nil
}

func (c *Client) AskOpenAI(message string) (*OpenAIResp, error) {
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

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to create a request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.OpenAIAPIKey))

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
