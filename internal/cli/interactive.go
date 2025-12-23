package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"ludwig/internal/storage"
	"ludwig/internal/types"
	"ludwig/internal/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// StartInteractive runs the interactive bubbletea UI.
func StartInteractive() {
	taskStore, err := storage.NewFileTaskStorage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing task storage: %v\n", err)
		os.Exit(1)
	}

	m := NewModel(taskStore)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// Model represents the state of the application.
type Model struct {
	taskStore *storage.FileTaskStorage
	tasks     []types.Task
	textInput textinput.Model
	commands  []utils.Command
	err       error
	message   string
}

// tickMsg is a message sent on a timer to trigger a refresh.
type tickMsg time.Time

// NewModel creates a new model with initial values.
func NewModel(taskStore *storage.FileTaskStorage) *Model {
	ti := textinput.New()
	ti.Placeholder = "...Enter command (e.g., 'add <task>', 'exit', 'help')"
	ti.Width = 50
	ti.Focus()

	tasks, err := taskStore.ListTasks()
	if err != nil {
		// This error will be displayed in the view.
		return &Model{err: fmt.Errorf("could not load tasks: %w", err)}
	}

	m := &Model{
		taskStore: taskStore,
		tasks:     utils.PointerSliceToValueSlice(tasks),
		textInput: ti,
	}
	m.commands = PalleteCommands(taskStore)
	return m
}


// Init initializes the application with a command to start the timer.
func (m *Model) Init() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles incoming messages and updates the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textInput.Value())
			parts := strings.Fields(input)
			m.textInput.SetValue("")
			//m.message = "" // Clear previous message
			m.err = nil     // Clear previous error

			if len(parts) == 0 {
				return m, nil
			}

			commandText := parts[0]
			if commandText == "exit" {
				return m, tea.Quit
			}

			for _, cmd := range m.commands {
				if cmd.Text == commandText {
					// Execute the command's action.
					if cmd.Action != nil {
						output := cmd.Action(strings.Join(parts, " "))
						m.message = output
					}
					// After action, refresh tasks immediately.
				tasks, err := m.taskStore.ListTasks()
					if err != nil {
						m.err = err
					} else {
						m.tasks = utils.PointerSliceToValueSlice(tasks)
					}
					return m, nil
				}
			}
			m.err = fmt.Errorf("command not found: %q", commandText)
			return m, nil
		}
	case tickMsg:
		// On each tick, reload tasks from storage.
		tasks, err := m.taskStore.ListTasks()
		if err != nil {
			m.err = err
		} else {
			if len(tasks) != len(m.tasks) {
				m.tasks = utils.PointerSliceToValueSlice(tasks)
			}
		}
		// Return a new tick command to continue polling.
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m *Model) View() string {
	var s strings.Builder
	// Render the Kanban board.
	s.WriteString(RenderKanban(m.tasks))
	//s.WriteString("\n")
	
	// Render status messages and errors before the input area
	s.WriteString("  " + m.message + "\n")

	if m.err != nil {
		s.WriteString("\nError: " + m.err.Error() + "\n")
	}
	
	// Render the text input for commands with bubble border.
	termWidth := utils.TermWidth()
	s.WriteString(utils.GenerateTopBubbleBorder(termWidth))
	
	// Calculate padding for the input to fit within the bubble
	inputWidth := termWidth - 6 // Account for borders and padding

	inputWidth = max(inputWidth, 20)

	
	// Set the textinput width to fit within the bubble
	m.textInput.Width = inputWidth
	
	// Render the middle of the bubble with the input
	inputLine := " │ " + m.textInput.View()
	// Pad the line to fill the bubble width
	neededPadding := termWidth - len(inputLine) - 1
	if neededPadding > 0 {
		inputLine += strings.Repeat(" ", neededPadding)
	}
	inputLine += "│ \n"
	s.WriteString(inputLine)
	
	s.WriteString(utils.GenerateBottomBubbleBorder(termWidth))
	
	return s.String()
}
