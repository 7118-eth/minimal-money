package ui

import (
	"fmt"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"github.com/bioharz/budget/internal/service"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

type View string

const (
	ViewMain          View = "main"
	ViewAssets        View = "assets"
	ViewAddAsset      View = "add_asset"
	ViewHistory       View = "history"
	ViewDeleteConfirm View = "delete_confirm"
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
	priceService *service.PriceService
	auditService *service.AuditService
	deletingHoldingID uint
}

func InitialModel() Model {
	m := Model{
		view:         ViewMain,
		prices:       make(map[uint]float64),
		accounts:     []models.Account{},
		assets:       []models.Asset{},
		holdings:     []models.Holding{},
		priceService: service.NewPriceService(),
		auditService: service.NewAuditService(),
	}
	m.setupTable()
	return m
}

func InitialModelWithDB(db *gorm.DB) Model {
	m := Model{
		view:         ViewMain,
		prices:       make(map[uint]float64),
		accounts:     []models.Account{},
		assets:       []models.Asset{},
		holdings:     []models.Holding{},
		priceService: service.NewPriceServiceWithDB(db),
		auditService: service.NewAuditServiceWithDB(db),
	}
	m.setupTable()
	return m
}

func (m Model) Init() tea.Cmd {
	// Return a command to load data
	return m.loadDataCmd()
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
				if m.view == ViewAddAsset {
					m.handleModalInput("backspace")
				} else if len(m.inputBuffer) > 0 {
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
		
		// Handle delete confirmation view
		if m.view == ViewDeleteConfirm {
			switch msg.String() {
			case "y", "Y":
				m.confirmDelete()
			case "n", "N", "esc":
				m.view = ViewMain
				m.deletingHoldingID = 0
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
			m.editSelectedHolding()
		case "d":
			m.deleteSelectedHolding()
		case "r":
			return m, m.refreshPrices()
		case "h":
			m.view = ViewHistory
		case "esc":
			if m.view == ViewDeleteConfirm {
				m.deletingHoldingID = 0
			}
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
		
	case priceUpdateMsg:
		if msg.err == nil && msg.prices != nil {
			m.prices = msg.prices
			m.updateTableData()
		}
		
	case dataLoadedMsg:
		m.accounts = msg.accounts
		m.assets = msg.assets
		m.holdings = msg.holdings
		m.updateTableData()
		// After loading data, fetch prices
		if len(m.assets) > 0 {
			return m, m.refreshPrices()
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
	case ViewDeleteConfirm:
		return m.deleteConfirmView()
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
	content := "📜 Audit Trail (Last 50 changes)\n\n"
	
	// Get audit logs
	logs, err := m.auditService.GetAllLogs(50)
	if err != nil {
		content += fmt.Sprintf("Error fetching audit logs: %v\n\n", err)
		content += "Press ESC to go back"
		return content
	}
	
	if len(logs) == 0 {
		content += "No audit history available yet.\n"
		content += "Changes to your portfolio will be tracked here.\n\n"
		content += "Press ESC to go back"
		return content
	}
	
	// Display audit logs
	for _, log := range logs {
		// Format action with emoji
		actionIcon := ""
		switch log.Action {
		case models.AuditActionCreate:
			actionIcon = "➕"
		case models.AuditActionUpdate:
			actionIcon = "✏️"
		case models.AuditActionDelete:
			actionIcon = "🗑️"
		}
		
		content += fmt.Sprintf("%s %s - %s\n", 
			log.CreatedAt.Format("2006-01-02 15:04"),
			actionIcon,
			log.Action)
		
		// Parse and display the changes
		if log.EntityType == models.AuditEntityHolding {
			content += m.formatHoldingChange(log)
		}
		
		content += "───────────────────────────────────\n"
	}
	
	content += "\nPress ESC to go back"
	return content
}

func (m Model) deleteConfirmView() string {
	// Find the holding details
	var holding models.Holding
	for _, h := range m.holdings {
		if h.ID == m.deletingHoldingID {
			holding = h
			break
		}
	}
	
	account := m.getAccountByID(holding.AccountID)
	asset := m.getAssetByID(holding.AssetID)
	value := holding.Amount * m.prices[holding.AssetID]
	
	content := m.tableView() + "\n\n"
	content += "┌─────────────────────────────────────────────┐\n"
	content += "│          Confirm Delete                     │\n"
	content += "├─────────────────────────────────────────────┤\n"
	content += fmt.Sprintf("│ Account: %-34s │\n", account.Name)
	content += fmt.Sprintf("│ Asset:   %-34s │\n", asset.Symbol)
	content += fmt.Sprintf("│ Amount:  %-34.4f │\n", holding.Amount)
	content += fmt.Sprintf("│ Value:   $%-33.2f │\n", value)
	content += "│                                             │\n"
	content += "│ Are you sure you want to delete this?      │\n"
	content += "│                                             │\n"
	content += "│        [Y]es     [N]o / [ESC]              │\n"
	content += "└─────────────────────────────────────────────┘"
	
	return content
}

func (m *Model) loadData() {
	// Load accounts
	accountRepo := repository.NewAccountRepository()
	accounts, err := accountRepo.GetAll()
	if err == nil {
		m.accounts = accounts
	}

	// Load assets
	assetRepo := repository.NewAssetRepository()
	assets, err := assetRepo.GetAll()
	if err == nil {
		m.assets = assets
	}

	// Load holdings
	holdingRepo := repository.NewHoldingRepository()
	holdings, err := holdingRepo.GetAll()
	if err == nil {
		m.holdings = holdings
	}

	// Update table with new data
	m.updateTableData()
	
	// Fetch initial prices
	if len(m.assets) > 0 {
		prices, err := m.priceService.FetchPrices(m.assets)
		if err == nil {
			m.prices = prices
			m.updateTableData()
		}
	}
}

func (m Model) refreshPrices() tea.Cmd {
	return func() tea.Msg {
		if m.priceService == nil || len(m.assets) == 0 {
			return nil
		}
		
		prices, err := m.priceService.FetchPrices(m.assets)
		if err != nil {
			return priceUpdateMsg{err: err}
		}
		
		return priceUpdateMsg{prices: prices}
	}
}

type priceUpdateMsg struct {
	prices map[uint]float64
	err    error
}

type dataLoadedMsg struct {
	accounts []models.Account
	assets   []models.Asset
	holdings []models.Holding
}

func (m Model) loadDataCmd() tea.Cmd {
	return func() tea.Msg {
		// Load accounts
		accountRepo := repository.NewAccountRepository()
		accounts, _ := accountRepo.GetAll()

		// Load assets
		assetRepo := repository.NewAssetRepository()
		assets, _ := assetRepo.GetAll()

		// Load holdings with relationships
		holdingRepo := repository.NewHoldingRepository()
		holdings, _ := holdingRepo.GetAll()

		return dataLoadedMsg{
			accounts: accounts,
			assets:   assets,
			holdings: holdings,
		}
	}
}

func (m *Model) deleteSelectedHolding() {
	// Get selected row index
	selectedRow := m.table.Cursor()
	if selectedRow >= len(m.holdings) {
		return
	}

	// Get the holding to delete
	holding := m.holdings[selectedRow]
	m.deletingHoldingID = holding.ID
	m.view = ViewDeleteConfirm
}

func (m *Model) confirmDelete() {
	// Find the holding before deletion for audit log
	var holdingToDelete models.Holding
	for _, h := range m.holdings {
		if h.ID == m.deletingHoldingID {
			holdingToDelete = h
			break
		}
	}

	// Delete from database
	holdingRepo := repository.NewHoldingRepository()
	if err := holdingRepo.Delete(m.deletingHoldingID); err != nil {
		m.err = err
		return
	}

	// Log the deletion to audit trail
	if m.auditService != nil && holdingToDelete.ID != 0 {
		m.auditService.LogHoldingDelete(&holdingToDelete)
	}

	// Reload data and return to main view
	m.loadData()
	m.view = ViewMain
	m.deletingHoldingID = 0
}

func (m *Model) editSelectedHolding() {
	// Get selected row index
	selectedRow := m.table.Cursor()
	if selectedRow >= len(m.holdings) {
		return
	}

	// Get the holding to edit
	holding := m.holdings[selectedRow]
	account := m.getAccountByID(holding.AccountID)
	asset := m.getAssetByID(holding.AssetID)

	// Initialize edit modal with existing values
	m.view = ViewAddAsset
	m.inputMode = true
	m.initEditAssetModal(holding, account, asset)
}