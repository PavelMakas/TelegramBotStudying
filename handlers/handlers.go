package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"

	"telegram-bot/utils"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func HandleMessage(message string) string {
	return utils.ReverseString(message)
}

func GenerateStory(style string, apiKey string) (string, error) {
	log.Printf("Starting story generation with style: %s", style)

	if apiKey == "" {
		log.Printf("Error: API key is empty")
		return "", fmt.Errorf("API key is required")
	}

	log.Printf("API Key length: %d", len(apiKey))
	prompt := fmt.Sprintf("Write a short %s story. The story should be engaging and creative, but not exceed 1000 characters (excluding spaces and punctuation). Focus on making it memorable and complete:", style)
	log.Printf("Generated prompt: %s", prompt)

	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Add status code logging
	log.Printf("API Response Status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	// Log the response for debugging
	log.Printf("API Response Body: %s", string(body))

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &openaiResp); err != nil {
		log.Printf("Error unmarshaling response: %v\nResponse body: %s", err, string(body))
		return "", err
	}

	if openaiResp.Error != nil {
		log.Printf("OpenAI API error: %s\nFull response: %s", openaiResp.Error.Message, string(body))
		return "", fmt.Errorf("OpenAI API error: %s", openaiResp.Error.Message)
	}

	if len(openaiResp.Choices) == 0 {
		log.Printf("No story generated")
		return "", fmt.Errorf("no story generated")
	}

	story := strings.TrimSpace(openaiResp.Choices[0].Message.Content)

	// Remove spaces and punctuation for character count
	storyWithoutSpaces := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, story)

	if len(storyWithoutSpaces) > 1000 {
		// Find the last complete word within the limit
		lastSpace := strings.LastIndex(story[:1000], " ")
		if lastSpace > 0 {
			story = story[:lastSpace] + "..."
		} else {
			story = story[:997] + "..."
		}
	}

	return story, nil
}
