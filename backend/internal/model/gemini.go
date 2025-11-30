package model

import (
	"context"
	"fmt"

	"github.com/o-ga09/gopher-lgtm-image-generator/pkg/config"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func NewModel(ctx context.Context) (model.LLM, error) {
	env := config.GetCtxEnv(ctx)
	// Initialize Gemini model
	model, err := gemini.NewModel(ctx, env.GeminiModel, &genai.ClientConfig{
		APIKey: env.GeminiAPIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini: %w", err)
	}
	return model, nil
}
