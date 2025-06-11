package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
