package api

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

// NewGeminiClient creates a new client with your API key.
func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	ctx := context.Background()
	cli, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating genai client: %w", err)
	}

	return &GeminiClient{
		client: cli,
		model:  "gemini-2.0-flash",
	}, nil
}

// ask sends the user's question to gemini and returns its answer
func (g *GeminiClient) Ask(ctx context.Context, question string) (string, error) {
	// 1. Marshal the body request
	chat, err := g.client.Chats.Create(ctx, g.model, nil, []*genai.Content{
		genai.NewContentFromText(question, genai.RoleUser),
	})
	if err != nil {
		return "", fmt.Errorf("creating gemini chat: %w", err)
	}

	// 3. Execute
	resp, err := chat.SendMessage(ctx, genai.Part{Text: question})
	if err != nil {
		return "", fmt.Errorf("sending message: %w", err)
	}
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}
	return resp.Candidates[0].Content.Parts[0].Text, nil
}
