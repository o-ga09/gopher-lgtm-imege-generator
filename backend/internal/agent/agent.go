package agent

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func NewAgent(ctx context.Context, model model.LLM, tools ...tool.Tool) (agent.Agent, error) {
	// Create agent with tools
	a, err := llmagent.New(llmagent.Config{
		Name:        "gopher-lgtm-image-generator-agent",
		Model:       model,
		Description: "The agent that generates LGTM images based on user requests.",
		Instruction: `You are an agent that generates LGTM images based on user requests. 
		- Use the provided tools to create and manipulate images as needed
		- Save the final image to the Cloudflare R2 bucket
		- Respond with the URL of the saved image
		**【MUST】** Always provide the URL of the saved image in your final response. 
		The URL should be formatted as: https://pub-0a072bc79aa54f28b971e7bd751566a4.r2.dev/[uploaded_path]`,
		Tools: tools,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	return a, nil
}
