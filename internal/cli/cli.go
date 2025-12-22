package cli

import (
	"fmt"
	types "ludwig/internal/types"
)


func Execute() {
	fmt.Println("Starting AI Orchestrator CLI...")
	// Call orchestrator and mcp as needed
	var tasks = types.ExampleTasks()
	DisplayKanban(tasks)
}
