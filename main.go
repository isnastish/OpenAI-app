package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

// TODO: Create openai package and move all the structs there
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

type FrontendRequestBody struct {
	OpenaiQuestion string `json:"openai-question"`
}

type OpenAIClient struct {
	OpenAIAPIKey string
	*http.Client
}

func NewOpenAIClient() (*OpenAIClient, error) {
	OPENAI_API_KEY, set := os.LookupEnv("OPENAI_API_KEY")
	if set == false || OPENAI_API_KEY == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	return &OpenAIClient{
		OpenAIAPIKey: OPENAI_API_KEY,
		Client:       &http.Client{},
	}, nil
}

func (c *OpenAIClient) AskOpenAI(message string) (*OpenAIResp, error) {
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

func main() {
	openaiClient, _ := NewOpenAIClient() // omit error for now

	// the server which will accept requests from the frontend
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "Fiber",
	})

	// CORS middleware
	app.Use("/", func(ctx *fiber.Ctx) error {
		fmt.Println("Middleware function was triggered")

		ctx.Set("Access-Control-Allow-Origin", "*")

		if ctx.Method() == "OPTIONS" {
			// NOTE: Without this header none of PUT/POST requests would work
			ctx.Set("Access-Control-Allow-Credentials", "true")
			ctx.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
			ctx.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			return nil
		}

		return ctx.Next()
	})

	app.Put("/api/openai/:message", func(ctx *fiber.Ctx) error {
		messsage := ctx.Params("message")
		fmt.Printf("Got a message: %s\n", messsage)

		reqBody := ctx.Request().Body()

		var reqData FrontendRequestBody
		if err := json.Unmarshal(reqBody, &reqData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		fmt.Printf("Openai question: %s\n", reqData.OpenaiQuestion)

		resp, err := openaiClient.AskOpenAI(reqData.OpenaiQuestion)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("openai: %s", err.Error()))
		}

		// The json will be returned to our frontend based on React
		// probably include status, message, API version
		return ctx.JSON(map[string]string{
			"model":  resp.Model,
			"openai": resp.Choices[0].Message.Content,
		}, "application/json")
	})

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3031"); err != nil {
			log.Fatalf("Server failed %v", err)
		}
	}()

	<-osSigChan
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown the server: %v", err)
	}

	os.Exit(0)
}
