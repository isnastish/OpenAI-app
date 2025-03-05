package query

import (
	"context"

	"github.com/isnastish/aiclient/internal/domain/ai"
	"github.com/isnastish/aiclient/internal/ports"
)

type AskAiHandler struct {
	aiRepo ports.AiRepository
}

func NewAskAiHandler(aiRepo ports.AiRepository) AskAiHandler {
	return AskAiHandler{
		aiRepo: aiRepo,
	}
}

func (h AskAiHandler) Handle(ctx context.Context, question string) (*ai.AiQueryResult, error) {
	return nil, nil
}
