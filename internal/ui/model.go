package ui

import (
	"github.com/bioharz/budget/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type View string

const (
	ViewMain      View = "main"
	ViewAssets    View = "assets"
	ViewAddAsset  View = "add_asset"
	ViewHistory   View = "history"
)

type Model struct {
	view         View
	assets       []models.Asset
	holdings     []models.Holding
	prices       map[uint]float64
	table        table.Model
	cursor       int
	width        int
	height       int
	err          error
	inputBuffer  string
	inputMode    bool
}

func InitialModel() Model {
	return Model{
		view:   ViewMain,
		prices: make(map[uint]float64),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.inputMode {
				return m, tea.Quit
			}
		case "1":
			if m.view == ViewMain && !m.inputMode {
				m.view = ViewAssets
			}
		case "2":
			if m.view == ViewMain && !m.inputMode {
				m.view = ViewAddAsset
				m.inputMode = true
			}
		case "3":
			if m.view == ViewMain && !m.inputMode {
				m.view = ViewHistory
			}
		case "escape":
			if m.inputMode {
				m.inputMode = false
				m.inputBuffer = ""
			}
			m.view = ViewMain
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	switch m.view {
	case ViewMain:
		return m.mainView()
	case ViewAssets:
		return m.assetsView()
	case ViewAddAsset:
		return m.addAssetView()
	case ViewHistory:
		return m.historyView()
	default:
		return "Unknown view"
	}
}

func (m Model) mainView() string {
	return `
Budget Tracker

1. View Assets
2. Add Asset
3. View History
q. Quit

Choose an option: `
}

func (m Model) assetsView() string {
	return "Assets view (ESC to go back)"
}

func (m Model) addAssetView() string {
	return "Add asset view (ESC to go back)\n\nEnter asset symbol: " + m.inputBuffer
}

func (m Model) historyView() string {
	return "History view (ESC to go back)"
}