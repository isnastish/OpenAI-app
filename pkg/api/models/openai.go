package models

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
