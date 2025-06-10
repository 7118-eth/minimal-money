package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type InputField struct {
	Label       string
	Value       string
	Placeholder string
	Active      bool
}

type ModalState struct {
	Fields        []InputField
	ActiveField   int
	ShowError     bool
	ErrorMessage  string
	IsEdit        bool
	EditingHoldingID uint
}

var (
	modalStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		Width(30)

	activeInputStyle = inputStyle.Copy().
		BorderForeground(lipgloss.Color("62"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	buttonStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color("238")).
		Foreground(lipgloss.Color("255"))

	activeButtonStyle = buttonStyle.Copy().
		Background(lipgloss.Color("62"))
)

func (m *Model) initAddAssetModal() {
	m.modalState = ModalState{
		Fields: []InputField{
			{Label: "Account", Value: "", Placeholder: "e.g., hardware wallet, NeoBank"},
			{Label: "Asset", Value: "", Placeholder: "e.g., BTC, ETH, USD"},
			{Label: "Amount", Value: "", Placeholder: "e.g., 0.5"},
			{Label: "Purchase Price", Value: "", Placeholder: "e.g., 40000 (optional)"},
		},
		ActiveField: 0,
		IsEdit: false,
	}
}

func (m *Model) initEditAssetModal(holding models.Holding, account models.Account, asset models.Asset) {
	purchasePrice := ""
	if holding.PurchasePrice > 0 {
		purchasePrice = fmt.Sprintf("%.2f", holding.PurchasePrice)
	}
	
	m.modalState = ModalState{
		Fields: []InputField{
			{Label: "Account", Value: account.Name, Placeholder: "e.g., hardware wallet, NeoBank"},
			{Label: "Asset", Value: asset.Symbol, Placeholder: "e.g., BTC, ETH, USD"},
			{Label: "Amount", Value: fmt.Sprintf("%.6f", holding.Amount), Placeholder: "e.g., 0.5"},
			{Label: "Purchase Price", Value: purchasePrice, Placeholder: "e.g., 40000 (optional)"},
		},
		ActiveField: 0,
		IsEdit: true,
		EditingHoldingID: holding.ID,
	}
}

func (m *Model) handleModalInput(key string) {
	switch key {
	case "tab":
		m.modalState.ActiveField = (m.modalState.ActiveField + 1) % (len(m.modalState.Fields) + 2) // +2 for Save/Cancel buttons
	case "shift+tab":
		m.modalState.ActiveField--
		if m.modalState.ActiveField < 0 {
			m.modalState.ActiveField = len(m.modalState.Fields) + 1
		}
	case "enter":
		if m.modalState.ActiveField == len(m.modalState.Fields) { // Save button
			m.saveAsset()
		} else if m.modalState.ActiveField == len(m.modalState.Fields) + 1 { // Cancel button
			m.view = ViewMain
			m.inputMode = false
			m.inputBuffer = ""
			m.modalState = ModalState{}
		}
	case "backspace":
		if m.modalState.ActiveField < len(m.modalState.Fields) {
			field := &m.modalState.Fields[m.modalState.ActiveField]
			if len(field.Value) > 0 {
				field.Value = field.Value[:len(field.Value)-1]
			}
		}
	default:
		if m.modalState.ActiveField < len(m.modalState.Fields) && len(key) == 1 {
			m.modalState.Fields[m.modalState.ActiveField].Value += key
		}
	}
}

func (m *Model) saveAsset() {
	// Validate inputs
	accountName := strings.TrimSpace(m.modalState.Fields[0].Value)
	assetSymbol := strings.TrimSpace(m.modalState.Fields[1].Value)
	amountStr := strings.TrimSpace(m.modalState.Fields[2].Value)
	priceStr := strings.TrimSpace(m.modalState.Fields[3].Value)

	if accountName == "" || assetSymbol == "" || amountStr == "" {
		m.modalState.ShowError = true
		m.modalState.ErrorMessage = "Account, Asset, and Amount are required"
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		m.modalState.ShowError = true
		m.modalState.ErrorMessage = "Invalid amount"
		return
	}

	var purchasePrice float64
	if priceStr != "" {
		purchasePrice, err = strconv.ParseFloat(priceStr, 64)
		if err != nil || purchasePrice < 0 {
			m.modalState.ShowError = true
			m.modalState.ErrorMessage = "Invalid purchase price"
			return
		}
	}

	// Get or create account
	accountRepo := repository.NewAccountRepository()
	account, err := accountRepo.GetByName(accountName)
	if err == gorm.ErrRecordNotFound {
		account = models.Account{
			Name: accountName,
			Type: "unknown", // TODO: Add account type selection
		}
		if err := accountRepo.Create(&account); err != nil {
			m.modalState.ShowError = true
			m.modalState.ErrorMessage = "Failed to create account"
			return
		}
	} else if err != nil {
		m.modalState.ShowError = true
		m.modalState.ErrorMessage = "Database error"
		return
	}

	// Get or create asset
	assetRepo := repository.NewAssetRepository()
	asset, err := assetRepo.GetBySymbol(strings.ToUpper(assetSymbol))
	if err == gorm.ErrRecordNotFound {
		// Determine asset type based on symbol
		assetType := m.guessAssetType(assetSymbol)
		asset = models.Asset{
			Symbol: strings.ToUpper(assetSymbol),
			Name:   assetSymbol, // TODO: Fetch proper name from API
			Type:   assetType,
		}
		if err := assetRepo.Create(&asset); err != nil {
			m.modalState.ShowError = true
			m.modalState.ErrorMessage = "Failed to create asset"
			return
		}
	} else if err != nil {
		m.modalState.ShowError = true
		m.modalState.ErrorMessage = "Database error"
		return
	}

	holdingRepo := repository.NewHoldingRepository()
	
	if m.modalState.IsEdit {
		// Update existing holding
		holding := models.Holding{
			ID:            m.modalState.EditingHoldingID,
			AccountID:     account.ID,
			AssetID:       asset.ID,
			Amount:        amount,
			PurchasePrice: purchasePrice,
			// Keep original purchase date for edits
		}
		if err := holdingRepo.Update(&holding); err != nil {
			m.modalState.ShowError = true
			m.modalState.ErrorMessage = "Failed to update holding"
			return
		}
	} else {
		// Create new holding
		holding := models.Holding{
			AccountID:     account.ID,
			AssetID:       asset.ID,
			Amount:        amount,
			PurchasePrice: purchasePrice,
			PurchaseDate:  time.Now(),
		}
		if err := holdingRepo.Create(&holding); err != nil {
			m.modalState.ShowError = true
			m.modalState.ErrorMessage = "Failed to create holding"
			return
		}
	}

	// Success - reload data and close modal
	m.loadData()
	m.view = ViewMain
	m.inputMode = false
	m.modalState = ModalState{}
}

func (m *Model) renderAddAssetModal() string {
	var b strings.Builder

	title := "Add New Asset"
	if m.modalState.IsEdit {
		title = "Edit Asset"
	}
	b.WriteString(titleStyle.Render(title) + "\n\n")

	// Render input fields
	for i, field := range m.modalState.Fields {
		label := labelStyle.Render(field.Label + ":")
		
		value := field.Value
		if value == "" && m.modalState.ActiveField != i {
			value = field.Placeholder
		}
		
		style := inputStyle
		if m.modalState.ActiveField == i {
			style = activeInputStyle
			value += "█"
		}
		
		input := style.Render(value)
		b.WriteString(fmt.Sprintf("%-15s %s\n\n", label, input))
	}

	// Error message
	if m.modalState.ShowError {
		b.WriteString(errorStyle.Render(m.modalState.ErrorMessage) + "\n\n")
	}

	// Buttons
	saveStyle := buttonStyle
	cancelStyle := buttonStyle
	
	if m.modalState.ActiveField == len(m.modalState.Fields) {
		saveStyle = activeButtonStyle
	} else if m.modalState.ActiveField == len(m.modalState.Fields) + 1 {
		cancelStyle = activeButtonStyle
	}
	
	buttons := fmt.Sprintf("%s  %s", 
		saveStyle.Render("Save"),
		cancelStyle.Render("Cancel"))
	
	b.WriteString(strings.Repeat(" ", 15) + buttons)

	return modalStyle.Render(b.String())
}

func (m *Model) guessAssetType(symbol string) models.AssetType {
	symbol = strings.ToUpper(symbol)
	
	// Common fiat currencies
	fiatSymbols := []string{"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD", "NZD"}
	for _, fiat := range fiatSymbols {
		if symbol == fiat {
			return models.AssetTypeFiat
		}
	}
	
	// Common crypto currencies
	cryptoSymbols := []string{"BTC", "ETH", "USDT", "USDC", "BNB", "XRP", "SOL", "ADA"}
	for _, crypto := range cryptoSymbols {
		if symbol == crypto {
			return models.AssetTypeCrypto
		}
	}
	
	// Default to crypto for unknown symbols
	return models.AssetTypeCrypto
}