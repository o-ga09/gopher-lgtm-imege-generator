package main

import (
	"context"
	"log"

	"github.com/o-ga09/gopher-lgtm-image-generator/internal/agent"
	"github.com/o-ga09/gopher-lgtm-image-generator/internal/model"
	"github.com/o-ga09/gopher-lgtm-image-generator/internal/server"
	"github.com/o-ga09/gopher-lgtm-image-generator/internal/tools"
	"github.com/o-ga09/gopher-lgtm-image-generator/pkg/config"
	"github.com/o-ga09/gopher-lgtm-image-generator/pkg/logger"
	"google.golang.org/adk/tool/functiontool"
)

func main() {
	ctx, err := config.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	logger.Logger(ctx)

	m, err := model.NewModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	genImageTool, err := functiontool.New(functiontool.Config{
		Name:        "GenerateImageTool",
		Description: "Generates an LGTM image based on the provided prompt and saves it as an artifact.",
	}, tools.GenerateImage)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}
	saveImageTool, err := functiontool.New(functiontool.Config{
		Name:        "SaveImageTool",
		Description: "Saves the generated LGTM image artifact to local storage.",
	}, tools.SaveImage)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	a, err := agent.NewAgent(ctx, m, genImageTool, saveImageTool)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	server, err := server.NewServer(ctx, a)
	if err != nil {
		log.Fatalf("Failed to create agent server: %v", err)
	}

	env := config.GetCtxEnv(ctx)
	switch env.Env {
	case "local":
		if err := server.DebugServer(ctx); err != nil {
			log.Fatalf("Failed to start debug server: %v", err)
		}
	case "prod":
		if err := server.Start(ctx); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
