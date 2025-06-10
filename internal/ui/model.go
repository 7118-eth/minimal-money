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
		case "esc":
			if m.inputMode {
				m.inputMode = false
				m.inputBuffer = ""
			}
			m.view = ViewMain
			return m, nil
		case "enter":
			if m.inputMode && m.view == ViewAddAsset {
				// TODO: Process the input
				m.inputMode = false
				m.inputBuffer = ""
				m.view = ViewMain
			}
			return m, nil
		case "backspace":
			if m.inputMode && len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
			return m, nil
		default:
			if m.inputMode {
				m.inputBuffer += msg.String()
				return m, nil
			}
			
			switch msg.String() {
			case "1":
				if m.view == ViewMain {
					m.view = ViewAssets
				}
			case "2":
				if m.view == ViewMain {
					m.view = ViewAddAsset
					m.inputMode = true
				}
			case "3":
				if m.view == ViewMain {
					m.view = ViewHistory
				}
			}
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
💰 Budget Tracker

1. View Assets
2. Add Asset
3. View History
q. Quit

Choose an option: `
}

func (m Model) assetsView() string {
	content := "📊 Your Assets\n\n"
	content += "No assets yet. Press 2 from main menu to add assets.\n\n"
	content += "Press ESC to go back"
	return content
}

func (m Model) addAssetView() string {
	content := "➕ Add New Asset\n\n"
	content += "Enter asset symbol (e.g., BTC, ETH, USD, EUR): " + m.inputBuffer
	if m.inputMode {
		content += "█"
	}
	content += "\n\n"
	content += "Press ENTER to add, ESC to cancel"
	return content
}

func (m Model) historyView() string {
	content := "📈 Price History\n\n"
	content += "No history data available yet.\n\n"
	content += "Press ESC to go back"
	return content
}