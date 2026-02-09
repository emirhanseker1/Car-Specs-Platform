package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ChatService interface {
	ProcessMessage(message string) (string, error)
}

type chatService struct {
	apiKey  string
	baseURL string
	model   string
}

func NewChatService() ChatService {
	// Try Groq first
	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")
	model := os.Getenv("GROQ_MODEL")

	// Fallback to OpenAI if Groq is not set
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if baseURL == "" {
		baseURL = os.Getenv("OPENAI_BASE_URL")
	}
	if model == "" {
		model = os.Getenv("OPENAI_MODEL")
	}

	// Defaults
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	}
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	return &chatService{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
	}
}



func (s *chatService) ProcessMessage(message string) (string, error) {
	// Define the Python Microservice URL
	pythonServiceURL := "http://localhost:8000/ask"

	// Create request payload
	reqBody, err := json.Marshal(map[string]string{
		"message": message,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", pythonServiceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Fallback error message if service is down
		return "⚠️ Bağlantı Hatası: Python AI servisi yanıt vermiyor. Lütfen servisin çalıştığından emin olun (port 8000).", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Python service error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var pythonResp struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pythonResp); err != nil {
		return "", fmt.Errorf("failed to decode python service response: %v", err)
	}

	return pythonResp.Response, nil
}
