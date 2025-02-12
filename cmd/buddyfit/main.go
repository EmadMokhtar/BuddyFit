package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EmadMokhtar/BuddyFit/internal"
)

func main() {
	// Define a flag for the required argument
	prompt := flag.String("prompt", "", "Prompt for the AI")
	p := flag.String("p", "", "Alias for prompt")
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

	responseChan := internal.AskAI(*prompt)

	for response := range responseChan {
		fmt.Print(response)
		// FIXME: This is the glamour code doesn't support the rendering of the stream response.
		//out, err := r.Render(response)
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
		//fmt.Print(out)
	}
}
