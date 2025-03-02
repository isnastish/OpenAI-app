package ports

import (
	"context"

	"github.com/isnastish/openai/internal/domain/ai"
)

type AiRepository interface {
	AskAi(ctx context.Context, question string) (*ai.AiQueryResult, error)
}
