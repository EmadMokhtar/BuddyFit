package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"

	"github.com/jackc/pgx/v5"
)

func AskAI(prompt string) string {
	dsn := os.Getenv("BF_DB_URL")

	if dsn == "" {
		fmt.Fprintf(os.Stderr, "BF_DB_URL environment variable is not set\n")
		os.Exit(1)
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if openAIKey == "" && ollamaHost == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY or OLLAMA_HOST environment variable is not set\n")
		os.Exit(1)
	}

	// Set PGOPTIONS environment variable
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse DSN: %v\n", err)
		os.Exit(1)
	}
	if openAIKey != "" {
		connConfig.RuntimeParams["options"] = fmt.Sprintf("-c ai.openai_api_key=%s", openAIKey)
	}

	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Ask pgai to generate a response
	var response string
	err = conn.QueryRow(ctx, "SELECT get_related_docs($1);", prompt).Scan(&response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get response: %v\n", err)
		os.Exit(1)
	}

	// TODO: call Ollama stream API to get the response
	sysProm := "You are a helpful gym personal trainer and professional bodybuilder. Use only the context provided to answer the question. Also mention the titles of the youtube videos you use to answer the question."
	userProm := fmt.Sprintf("Context: %s \n\n User Question: %s", response, prompt)
	model := "llama3.1:latest"
	chat := NewAIChat(model, sysProm, userProm)

	chatJSON, err := json.Marshal(chat)
	if err != nil {
		fmt.Printf("Error marshalling AIAgent to JSON: %v\n", err)
		os.Exit(1)
	}

	// Create a new request using http
	// TODO: read from env
	ollamaPort := "11434"
	ollamaChatURL := fmt.Sprintf("http://%s:%s/api/chat", ollamaHost, ollamaPort)
	log.Printf("Ollama Chat URL: %s\n", ollamaChatURL)
	req, err := http.NewRequest("POST", ollamaChatURL, bytes.NewBuffer(chatJSON))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the streaming response
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var streamResp struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done bool `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			continue
		}

		fullResponse.WriteString(streamResp.Message.Content)
		
		// If this is the last message, break the loop
		if streamResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	return fullResponse.String()
}

type AIChat struct {
	Model     string      `json:"model"`
	Messages  []AIMessage `json:"messages"`
	Stream    bool        `json:"stream"`
	KeepAlive string      `json:"keep_alive"`
}

func NewAIChat(model string, systemProm string, userProm string) *AIChat {
	return &AIChat{
		Model: model,
		Messages: []AIMessage{
			{Role: "system", Content: systemProm},
			{Role: "user", Content: userProm},
		},
		Stream:    true,
		KeepAlive: "30s",
	}
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
