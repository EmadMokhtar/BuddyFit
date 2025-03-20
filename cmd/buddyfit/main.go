package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"

	"github.com/EmadMokhtar/BuddyFit/internal"
)

func main() {
	// Define a flag for the required argument
	prompt := flag.String("prompt", "", "Prompt for the AI")
	p := flag.String("p", "", "Alias for prompt")
	model := flag.String("model", "llama3.1:latest", "Model for the AI")
	flag.Parse()
	noPrompt := *prompt == "" && *p == ""

	// Check if the argument is provided
	if noPrompt {
		fmt.Println("Error: -prompt is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if the prompt is empty and the alias is not
	// If so, set the alias to the prompt
	if *prompt == "" && *p != "" {
		prompt = p
	}

	responseChan := internal.AskAI(*prompt, *model)
	var fullResp strings.Builder

	for response := range responseChan {
		fmt.Print(response)
		fullResp.WriteString(response)
	}
	fmt.Print("\033[2J")
	out, err := glamour.Render(fullResp.String(), "dark")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(out)
}
