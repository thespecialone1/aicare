package api

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type QAService struct {
	gemini  *GeminiClient
	repo    *QuestionRepo
	msgRepo *MessageRepo
}

func NewQAService(g *GeminiClient, qRepo *QuestionRepo, mRepo *MessageRepo) *QAService {
	return &QAService{
		gemini:  g,
		repo:    qRepo,
		msgRepo: mRepo,
	}
}

func (s *QAService) AskAndSave(ctx context.Context, userID int, question string, mode string) (string, error) {
	// 1. Presist the user's question
	if err := s.msgRepo.Save(userID, "user", question); err != nil {
		return "", fmt.Errorf("saving user messages: %w", err)
	}

	// 2. Load full history
	historyRows, err := s.msgRepo.ListMessages(userID)
	if err != nil {
		return "", fmt.Errorf("listing history: %w", err)
	}

	// 3. Select the system prompt based on mode
	var systemPrompt string
	switch mode {
	case "diag":
		systemPrompt = `You are a medical diagnostic assistant. Conduct a structured interview: first ask for the patient's age, gender, and main complaint. Then ask about duration, severity, and relevant history. Only after gathering all details, offer possible diagnoses.`
	default:
		systemPrompt = `You are a qualified medical assistant. Answer medical questions concisely in simple language without citing sources.`
	}

	// 4. Build genai contents with a medical system prompt
	contents := []*genai.Content{
		genai.NewContentFromText(systemPrompt, genai.RoleUser),
	}

	for _, m := range historyRows {
		var role genai.Role
		switch m.Role {
		case "assistant":
			role = genai.RoleModel
		default:
			role = genai.RoleUser
		}
		contents = append(contents, genai.NewContentFromText(m.Content, role))
	}

	// 5. Create the chat and sends the message
	chat, err := s.gemini.client.Chats.Create(ctx, s.gemini.model, nil, contents)
	if err != nil {
		return "", fmt.Errorf("creating chat: %w", err)
	}
	res, err := chat.SendMessage(ctx, genai.Part{Text: question})
	if err != nil {
		return "", fmt.Errorf("seding message: %w", err)
	}

	// 6. Extract the assistant's answer
	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no answer returned")
	}
	answer := res.Candidates[0].Content.Parts[0].Text

	// 7. Save the assistant's answer
	if err := s.msgRepo.Save(userID, "assistant", answer); err != nil {
		return "", fmt.Errorf("saving assistant message: %w", err)
	}

	// 8. (Optional) also record in questions table
	if err := s.repo.Save(userID, question, answer); err != nil {
		return "", fmt.Errorf("saving Q&A: %w", err)
	}
	return answer, nil
}
