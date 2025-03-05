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
	result, err := h.aiRepo.AskAi(ctx, question)
	if err != nil {
		return nil, err
	}
	return result, nil
}
