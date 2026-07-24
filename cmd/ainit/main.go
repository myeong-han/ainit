package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/myeong-han/ainit/pkg/config"
	"github.com/myeong-han/ainit/pkg/tui"
)

func main() {
	cfg := config.NewDefaultConfig()

	p := tea.NewProgram(tui.NewModel(cfg), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running Agentic-Init (ainit) TUI: %v\n", err)
		os.Exit(1)
	}
}
