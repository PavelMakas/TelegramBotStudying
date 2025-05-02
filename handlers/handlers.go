package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"telegram-bot/utils"
)

type StoryRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StoryResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func HandleMessage(message string) string {
	return utils.ReverseString(message)
}

func GenerateStory(style string, apiKey string) (string, error) {
	prompt := fmt.Sprintf("Generate a short %s story (maximum 400 characters).", style)

	reqBody := StoryRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a creative storyteller. Keep stories under 400 characters.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var storyResp StoryResponse
	if err := json.Unmarshal(body, &storyResp); err != nil {
		return "", err
	}

	if len(storyResp.Choices) == 0 {
		return "", fmt.Errorf("no story generated")
	}

	story := storyResp.Choices[0].Message.Content
	if len(story) > 400 {
		story = story[:397] + "..."
	}

	return story, nil
}
