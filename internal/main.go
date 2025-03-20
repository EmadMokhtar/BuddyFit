package internal

import (
	"fmt"
	"os"

	"github.com/EmadMokhtar/BuddyFit/internal/agent"
	"github.com/EmadMokhtar/BuddyFit/internal/config"
)

func AskAI(usrPrompt string, model string) chan string {
	openAIKey := config.GetEnvWithDefault("OPENAI_API_KEY", "")
	ollamaHost := config.GetEnvWithDefault("OLLAMA_HOST", "")
	ollamaPort := config.GetEnvWithDefault("OLLAMA_PORT", "11434")

	if openAIKey == "" && ollamaHost == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY or OLLAMA_HOST environment variable is not set\n")
		os.Exit(1)
	}

	var aiProviderConfig *config.AIProviderConfig
	if ollamaHost != "" {
		aiProviderConfig = config.NewOllamaConfig(ollamaHost, ollamaPort)
	} else if openAIKey != "" {
		aiProviderConfig = config.NewOpenAIConfig(openAIKey)
	}

	buddyfit := agent.NewBuddyFitAgent(model, *aiProviderConfig)
	buddyfit.AddUserMessage(usrPrompt)
	return buddyfit.CompleteChat()
}
