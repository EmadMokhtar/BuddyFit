package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/EmadMokhtar/BuddyFit/internal/config"
)

type Options struct {
	Temperature int `json:"temperature"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Agent struct {
	Model     string      `json:"model"`
	Messages  []AIMessage `json:"messages"`
	Stream    bool        `json:"stream"`
	KeepAlive string      `json:"keep_alive"`
	Options   Options     `json:"options"`
	config    config.AIProviderConfig
}

func (a *Agent) GetContext(usrPrompt string) string {
	dsn := os.Getenv("BF_DB_URL")

	if dsn == "" {
		fmt.Fprintf(os.Stderr, "BF_DB_URL environment variable is not set\n")
		os.Exit(1)
	}
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set PGOPTIONS environment variable
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse DSN: %v\n", err)
		os.Exit(1)
	}

	switch a.config.Name {
	case "openai":
		connConfig.RuntimeParams["options"] = fmt.Sprintf("-c ai.openai_api_key=%s", a.config.Key)
	case "ollama":
		connConfig.RuntimeParams["options"] = fmt.Sprintf("-c ai.ollama_host=%s", a.config.Host)

	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	// Ask pgai to get related docs using RAG
	var retdDocs string
	err = conn.QueryRow(ctx, "SELECT get_related_docs($1);", usrPrompt).Scan(&retdDocs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get response: %v\n", err)
		os.Exit(1)
	}
	return retdDocs
}

func (a *Agent) AddUserMessage(usrPrompt string) {
	retdDocs := a.GetContext(usrPrompt)
	tmpl, err := template.ParseFiles("templates/prompt_template.tmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}
	var prompt bytes.Buffer
	err = tmpl.Execute(&prompt, map[string]string{"Context": retdDocs, "UserPrompt": usrPrompt})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	a.Messages = append(a.Messages, AIMessage{Role: "user", Content: prompt.String()})
}

func (a *Agent) CompleteChat() chan string {
	chatJSON, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Error marshalling AIAgent to JSON: %v\n", err)
		os.Exit(1)
	}

	ollamaChatURL := a.config.GetOllamaAPIURL()
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

func NewBuddyFitAgent(model string, providerConfig config.AIProviderConfig) *Agent {
	sysProm := "Your name is BuddyFit. You are an experienced personal trainer and professional bodybuilder with expertise in weight lifting techniques, bodybuilding programs, nutrition, and fitness science.\n\nWhen responding to questions:\n1. Use only the context provided to answer the question\n2. Always include both the titles AND URLs of any YouTube videos you reference\n3. Provide evidence-based advice when possible\n4. Prioritize safety and proper form in all recommendations\n5. Tailor advice to the user's experience level when specified\n6. Offer both practical advice and scientific reasoning behind recommendations\n7. Be encouraging and motivational in your tone\n\nFormat your entire response using clean, structured markdown:\n\n- Start with a top-level header: # BuddyFit Trainer\n\n- Use second-level headers for main sections: ## Your Training Plan\n\n- Use bold text for exercise names: **Barbell Bench Press**\n\n- Use proper markdown lists with indentation for exercise details:\n  * **Exercise Name**\n    * Sets: 3\n    * Reps: 8-12\n    * Rest: 60-90 seconds\n\n- Format the YouTube video section with a dedicated header:\n  ## Recommended Videos\n  * [Full Body Workout Guide](https://youtube.com/example)\n  * [Proper Squat Form](https://youtube.com/example2)\n\n- End with a horizontal rule and motivational quote:\n  ---\n  > Train hard, recover well, and stay consistent!\n\nIf you don't have sufficient context to answer a question safely, acknowledge the limitations and suggest what additional information would be helpful."

	return &Agent{
		Model: model,
		Messages: []AIMessage{
			{Role: "system", Content: sysProm},
		},
		Stream:    true,
		KeepAlive: "30s",
		Options: Options{
			Temperature: 0,
		},
		config: providerConfig,
	}
}
