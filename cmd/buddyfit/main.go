package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EmadMokhtar/BuddyFit/internal"

	"github.com/charmbracelet/glamour"
)

func main() {
	// Define a flag for the required argument
	prompt := flag.String("prompt", "", "Prompt for the AI")
	p := flag.String("p", "", "Alias for prompt")
	flag.Parse()

	// Check if the argument is provided
	if *prompt == "" {
		if *p == "" {
			fmt.Println("Error: -prompt is required")
			flag.Usage()
			os.Exit(1)
		} else {
			*prompt = *p
		}
	}

	// Call the AskAI function
	response := internal.AskAI(*prompt)

	out, err := glamour.Render(response, "dark")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(out)
}
