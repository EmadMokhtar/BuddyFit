package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetOllamaAPIURL() string {
	ollamaHost := GetEnvWithDefault("OLLAMA_HOST", "localhost")
	ollamaPort := GetEnvWithDefault("OLLAMA_PORT", "11434")
	return fmt.Sprintf("http://%s:%s/api/chat", ollamaHost, ollamaPort)
}

func AskAI(prompt string) chan string {
	dsn := os.Getenv("BF_DB_URL")

	if dsn == "" {
		fmt.Fprintf(os.Stderr, "BF_DB_URL environment variable is not set\n")
		os.Exit(1)
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	ollamaHost := GetEnvWithDefault("OLLAMA_HOST", "")

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

	if ollamaHost != "" {
		connConfig.RuntimeParams["options"] = fmt.Sprintf("-c ai.ollama_host=%s", ollamaHost)
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

	// TODO: Create an Agent to encapsulate the chat logic
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
	ollamaChatURL := GetOllamaAPIURL()
	log.Printf("Ollama Chat URL: %s\n", ollamaChatURL)
	req, err := http.NewRequest("POST", ollamaChatURL, bytes.NewBuffer(chatJSON))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")
	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	//defer resp.Body.Close()

	// Read the streaming response
	scanner := bufio.NewScanner(resp.Body)
	responseChan := make(chan string)

	go func(resp *http.Response) {
		defer close(responseChan)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var streamResp struct {
				Message struct {
					Content string `json:"content"`
					Role    string `json:"role"`
				} `json:"message"`
				Done bool `json:"done"`
			}

			if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				continue
			}

			responseChan <- streamResp.Message.Content

			// If this is the last message, break the loop
			if streamResp.Done {
				resp.Body.Close()
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			os.Exit(1)
		}
	}(resp)

	return responseChan
}

type Options struct {
	Temperature int `json:"temperature"`
}

type AIChat struct {
	Model     string      `json:"model"`
	Messages  []AIMessage `json:"messages"`
	Stream    bool        `json:"stream"`
	KeepAlive string      `json:"keep_alive"`
	Options   Options     `json:"options"`
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
		Options: Options{
			Temperature: 0,
		},
	}
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
