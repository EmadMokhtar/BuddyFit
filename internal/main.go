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

	if openAIKey == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY environment variable is not set\n")
		os.Exit(1)
	}

	// Set PGOPTIONS environment variable
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse DSN: %v\n", err)
		os.Exit(1)
	}
	connConfig.RuntimeParams["options"] = fmt.Sprintf("-c ai.openai_api_key=%s", openAIKey)

	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
