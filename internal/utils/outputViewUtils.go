package utils

import (
	"strings"
	"encoding/json"
	"time"
	//tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func removeSurroundingQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

func OutputLine(line string) string {
	// remove surrounding {}
	if len(line) == 0 {
		return line
	}
	if line[0] != '{' || len(line) < 2 {
		return line
	}
	var object map[string]any
	err := json.Unmarshal([]byte(line), &object)
	builder := strings.Builder{}
	if err != nil {
		return ""
	}
	if object["type"] == "init" {
		return ""
	}
	if object["type"] == "message" {
		timestamp := object["timestamp"].(string)
		time, err := time.Parse(time.RFC3339, timestamp)
		if err == nil {
			layout := "2006-01-02 15:04:05"
			formattedTime := time.Format(layout)
			style := lipgloss.NewStyle().Faint(true)
			formattedTime = style.Render(formattedTime)
			builder.WriteString(ColoredString(formattedTime, "240"))
			builder.WriteString("\n")
		}
		builder.WriteString(object["content"].(string))
		return builder.String()
	}

	return line
}

func OutputLines(lines []string) string {
	output := strings.Builder{}
	if len(lines) == 0 {
		return "no output"
	}
	started := false
	for _, line := range lines {
		if line == "---" {
			started = true
			continue
		}
		if !started {
			continue
		}
		if line == "" {
			continue
		}
		output.WriteString(OutputLine(line))
		output.WriteString("\n")
	}
	return output.String()
}
