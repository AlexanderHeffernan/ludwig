package orchestrator

import (
	"fmt"
	"ludwig/internal/orchestrator/clients"
)

func Start() {
	gemini := &clients.GeminiClient{}
	prompt := "Say Hello, World!"
	response, err := gemini.SendPrompt(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Gemini response:", response)
}
