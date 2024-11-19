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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
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

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

///////////////////////////////////////////////////////////////////
// JWT Auth

// Custom claims
type Claims struct {
	Email    string `json:"email"`
	Password string `json:"pwd"`
	jwt.RegisteredClaims
}

type TokensPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	app.Post("/api/login", func(ctx *fiber.Ctx) error {
		var userData UserData
		if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		// TODO: Email and password validation
		isValidData := true
		if !isValidData {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to validate data")
		}

		// Create token
		// Set the claims
		// Set an expiration time for jwt token
		// Create signed token
		// Create a refresh token and set the claims
		// Set the expiration time for refresh token
		// We have two token: an access_token and refresh_token

		// TODO: Generate JWT token with TTL 15 minutes,
		// and send it back to the client, so the next requests
		// should only be made including that token.
		// TODO: Determine which signing method to choose

		// Creat a new token with claims
		// TODO: Expiration time has to be an argument that we pass to a function
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			&Claims{
				Email:    userData.Email,
				Password: userData.Password,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now()),
					Issuer:    "test",
					Subject:   "somebody",
					ID:        "1",
					Audience:  []string{"somebody_else"},
				},
			})

		// Sign token using a secret key, it should be private key
		signedAccessToken, err := token.SignedString([]byte("my-secret"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to sign access token: %v", err))
		}

		// TODO: We should pass an expiration time as well.
		// Create refresh token with claims
		refreshToken := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			&Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "somebody", // supposed to be user ID?
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 45)),
				},
			},
		)

		signedRefreshToken, err := refreshToken.SignedString([]byte("my-secret"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to sign refresh token: %v", err))
		}

		tokensPair := TokensPairs{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
		}

		return ctx.JSON(tokensPair, "application/json")
	})

	app.Post("/api/signup", func(ctx *fiber.Ctx) error {
		var userData UserData
		if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		// TODO: Make sure that the user doesn't exist
		return fiber.NewError(fiber.StatusNotImplemented, "")
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
