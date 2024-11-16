package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func main() {
	OPENAI_API_KEY, set := os.LookupEnv("OPENAI_API_KEY")
	if set == false || OPENAI_API_KEY == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	client := &http.Client{}

	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
		{
			"role":    "user",
			"content": "Please, generate simple Golang programm.",
		},
	}

	data := map[string]interface{}{
		"model":    "gpt-4o-mini-2024-07-18",
		"messages": messages,
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal the body: %s", err)
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Failed to create a request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", OPENAI_API_KEY))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make a request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Response status code: %d, message: %s", resp.StatusCode, resp.Status)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Faild to read response body: %s", err)
	}

	var openAIResp OpenAIResp
	err = json.Unmarshal(respBytes, &openAIResp)
	if err != nil {
		log.Fatalf("Failed to parse response body: %s", err)
	}

	fmt.Printf("Model: %s\n", openAIResp.Model)
	fmt.Printf("Choices: %v\n", openAIResp.Choices[0].Message.Content)

	os.Exit(0)
}
