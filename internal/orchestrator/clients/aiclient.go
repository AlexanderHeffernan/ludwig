package clients

type AIClient interface {
	SendPrompt(prompt string) (string, error)
}
