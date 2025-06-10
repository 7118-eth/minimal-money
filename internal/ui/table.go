package ui

import (
	"fmt"
	"strings"

	"github.com/bioharz/budget/internal/models"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	totalStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))
)

func (m *Model) setupTable() {
	columns := []table.Column{
		{Title: "Account", Width: 12},
		{Title: "Asset", Width: 8},
		{Title: "Amount", Width: 12},
		{Title: "Price", Width: 10},
		{Title: "Value", Width: 12},
		{Title: "24h%", Width: 8},
		{Title: "P&L", Width: 10},
	}

	rows := m.buildTableRows()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle
	t.SetStyles(s)

	m.table = t
}

func (m *Model) buildTableRows() []table.Row {
	var rows []table.Row
	
	for _, holding := range m.holdings {
		account := m.getAccountByID(holding.AccountID)
		asset := m.getAssetByID(holding.AssetID)
		price := m.prices[holding.AssetID]
		value := holding.Amount * price
		
		var pl float64
		if holding.PurchasePrice > 0 {
			pl = (price - holding.PurchasePrice) * holding.Amount
		}
		
		row := table.Row{
			account.Name,
			asset.Symbol,
			fmt.Sprintf("%.4f", holding.Amount),
			fmt.Sprintf("$%.2f", price),
			fmt.Sprintf("$%.2f", value),
			"+0.0%", // TODO: implement 24h change
			fmt.Sprintf("$%.2f", pl),
		}
		rows = append(rows, row)
	}
	
	// Add sample data if no holdings
	if len(rows) == 0 {
		rows = append(rows, table.Row{
			"Example",
			"BTC",
			"0.5",
			"$42,000",
			"$21,000",
			"+2.3%",
			"+$500",
		})
	}
	
	return rows
}

func (m *Model) updateTableData() {
	rows := m.buildTableRows()
	m.table.SetRows(rows)
}

func (m *Model) calculateTotal() float64 {
	var total float64
	for _, holding := range m.holdings {
		price := m.prices[holding.AssetID]
		total += holding.Amount * price
	}
	return total
}

func (m *Model) getAccountByID(id uint) models.Account {
	for _, acc := range m.accounts {
		if acc.ID == id {
			return acc
		}
	}
	return models.Account{Name: "Unknown"}
}

func (m *Model) getAssetByID(id uint) models.Asset {
	for _, asset := range m.assets {
		if asset.ID == id {
			return asset
		}
	}
	return models.Asset{Symbol: "???"}
}

func (m *Model) tableView() string {
	total := m.calculateTotal()
	
	var b strings.Builder
	
	// Header
	header := fmt.Sprintf("💰 Budget Tracker%sTotal: $%.2f", 
		strings.Repeat(" ", 40), total)
	b.WriteString(totalStyle.Render(header) + "\n\n")
	
	// Table
	b.WriteString(baseStyle.Render(m.table.View()) + "\n\n")
	
	// Footer
	footer := "[n]ew  [e]dit  [d]elete  [r]efresh  [h]istory  [q]uit"
	b.WriteString(footer)
	
	return b.String()
}