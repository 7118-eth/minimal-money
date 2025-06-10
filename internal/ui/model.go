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
	accounts     []models.Account
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
	modalState   ModalState
}

func InitialModel() Model {
	m := Model{
		view:     ViewMain,
		prices:   make(map[uint]float64),
		accounts: []models.Account{},
		assets:   []models.Asset{},
		holdings: []models.Holding{},
	}
	m.setupTable()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "esc":
				m.inputMode = false
				m.inputBuffer = ""
				m.view = ViewMain
			case "enter":
				if m.view == ViewAddAsset {
					m.handleModalInput("enter")
					return m, nil
				}
			case "backspace":
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			case "tab":
				if m.view == ViewAddAsset {
					m.handleModalInput("tab")
				}
			default:
				if m.view == ViewAddAsset {
					m.handleModalInput(msg.String())
				} else {
					m.inputBuffer += msg.String()
				}
			}
			return m, nil
		}
		
		// Main table view keyboard handling
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			m.view = ViewAddAsset
			m.inputMode = true
			m.initAddAssetModal()
		case "e":
			// TODO: Edit selected holding
		case "d":
			// TODO: Delete selected holding
		case "r":
			// TODO: Refresh prices
		case "h":
			m.view = ViewHistory
		case "esc":
			m.view = ViewMain
		default:
			// Pass through to table for navigation
			if m.view == ViewMain {
				m.table, cmd = m.table.Update(msg)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.view == ViewMain {
			m.table.SetHeight(msg.Height - 8)
		}
	}

	return m, cmd
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
	return m.tableView()
}

func (m Model) assetsView() string {
	content := "📊 Your Assets\n\n"
	content += "No assets yet. Press 2 from main menu to add assets.\n\n"
	content += "Press ESC to go back"
	return content
}

func (m Model) addAssetView() string {
	return m.renderAddAssetModal()
}

func (m Model) historyView() string {
	content := "📈 Price History\n\n"
	content += "No history data available yet.\n\n"
	content += "Press ESC to go back"
	return content
}