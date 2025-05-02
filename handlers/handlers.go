package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"telegram-bot/utils"
)

type StoryRequest struct {
	Inputs     string `json:"inputs"`
	Parameters struct {
		MaxLength int `json:"max_length"`
	} `json:"parameters"`
}

func HandleMessage(message string) string {
	return utils.ReverseString(message)
}

func GenerateStory(style string, apiKey string) (string, error) {
	prompt := fmt.Sprintf("Write a short %s story in 400 characters or less:", style)

	reqBody := StoryRequest{
		Inputs: prompt,
		Parameters: struct {
			MaxLength int `json:"max_length"`
		}{
			MaxLength: 400,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/gpt2", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer hf_DDmXqXqXqXqXqXqXqXqXqXqXqXqXqXqXqXq")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	// Log the response for debugging
	log.Printf("API Response: %s", string(body))

	var storyResp []struct {
		GeneratedText string `json:"generated_text"`
	}

	if err := json.Unmarshal(body, &storyResp); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return "", err
	}

	if len(storyResp) == 0 {
		log.Printf("No story generated")
		return "", fmt.Errorf("no story generated")
	}

	story := storyResp[0].GeneratedText
	// Clean up the story by removing the prompt and extra whitespace
	story = strings.TrimPrefix(story, prompt)
	story = strings.TrimSpace(story)

	if len(story) > 400 {
		story = story[:397] + "..."
	}

	return story, nil
}
