package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/isnastish/aiclient/internal/domain/ai"
)

type OpenAiRepository struct {
	httpClient *http.Client
	apiKey     string
}

func NewOpenAiRepository( /* pass config which holds env variables */ ) *OpenAiRepository {
	apiKey, set := os.LookupEnv("OPENAI_API_KEY")
	if !set || apiKey == "" {
		panic("OPENAI_API_KEY is not set")
	}

	return &OpenAiRepository{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type openAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiChoiceEntry struct {
	Index   int           `json:"index"`
	Message openAiMessage `json:"message"`
}

type openAiResponse struct {
	Model   string              `json:"model"`
	Choices []openAiChoiceEntry `json:"choices"`
}

const openAiBaseURL = `https://api.openai.com/v1/chat/completions`

func (o OpenAiRepository) AskAi(ctx context.Context, question string) (*ai.AiQueryResult, error) {
	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
		{
			"role":    "user",
			"content": question,
		},
	}
	reqData := map[string]interface{}{
		"model":    "gpt-4o-mini-2024-07-18",
		"messages": messages,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body, error %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", openAiBaseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create a request, error %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: Read API documentation for possible error codes
	// if resp.StatusCode != http.StatusOK {
	// 	// log.Fatalf("Response status code: %d, message: %s", resp.StatusCode, resp.Status)
	// }

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body, error %v", err)
	}

	var res openAiResponse
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body, error %v", err)
	}

	return &ai.AiQueryResult{
		Data: res.Choices[0].Message.Content,
	}, nil
}
