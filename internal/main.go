package internal

import (
	"context"
	"fmt"
	"os"
	"time"

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
	} else {
		// TODO: Fix the OLLAMA_HOST that missing the 'http://'
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
	err = conn.QueryRow(ctx, "SELECT generate_rag_response($1);", prompt).Scan(&response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get response: %v\n", err)
		os.Exit(1)
	}

	return response
}
