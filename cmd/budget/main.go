package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Build-time variables (set via -ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Handle version flag
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Minimal Money %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		os.Exit(0)
	}
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
