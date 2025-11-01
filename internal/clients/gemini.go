package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient() (*GeminiClient, error) {
	ctx := context.Background()

	// Check if API key is set
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{client: client}, nil
}

func (g *GeminiClient) GenerateBirthdayWish(name string, age int) (string, error) {
	// Use timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prompt := fmt.Sprintf("Generate a warm and personalized birthday wish for %s who is turning %d years old. Make it heartfelt, positive, and celebratory. Keep it under 100 words.", name, age)

	result, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate birthday wish: %w", err)
	}

	if result.Text() == "" {
		return fmt.Sprintf("🎉 Happy %d Birthday, %s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! May this new year bring you endless blessings and amazing adventures! 🌟", age, name), nil
	}

	return result.Text(), nil
}

func (g *GeminiClient) GenerateGenericBirthdayWish(name string) (string, error) {
	// Use timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prompt := fmt.Sprintf("Generate a warm and personalized birthday wish for %s. Make it heartfelt, positive, and celebratory. Keep it under 100 words.", name)

	result, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate birthday wish: %w", err)
	}

	if result.Text() == "" {
		return fmt.Sprintf("🎉 Happy Birthday, %s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! May this year bring you endless blessings and amazing adventures! 🌟", name), nil
	}

	return result.Text(), nil
}
