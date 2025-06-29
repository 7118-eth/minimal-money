package ui

import (
	"fmt"
	"sort"
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
	// Calculate column widths based on terminal width
	availableWidth := m.width - 10 // Account for borders and padding
	if availableWidth < 80 {
		availableWidth = 80 // Minimum width
	}

	// Distribute width proportionally for new layout
	assetAccountWidth := int(float64(availableWidth) * 0.45)
	amountWidth := int(float64(availableWidth) * 0.25)
	valueWidth := int(float64(availableWidth) * 0.30)

	columns := []table.Column{
		{Title: "Asset/Account", Width: assetAccountWidth},
		{Title: "Amount", Width: amountWidth},
		{Title: "Value", Width: valueWidth},
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

	// Group holdings by asset
	assetHoldings := make(map[uint][]models.Holding)
	for _, holding := range m.holdings {
		assetHoldings[holding.AssetID] = append(assetHoldings[holding.AssetID], holding)
	}

	// Calculate total value per asset
	assetTotalValues := make(map[uint]float64)
	for assetID, holdings := range assetHoldings {
		price := m.prices[assetID]
		totalValue := 0.0
		for _, holding := range holdings {
			totalValue += holding.Amount * price
		}
		assetTotalValues[assetID] = totalValue
	}

	// Sort assets by total value (highest first)
	var assetIDs []uint
	for assetID := range assetHoldings {
		assetIDs = append(assetIDs, assetID)
	}
	sort.Slice(assetIDs, func(i, j int) bool {
		// Sort by total value descending
		valueI := assetTotalValues[assetIDs[i]]
		valueJ := assetTotalValues[assetIDs[j]]

		// If values are equal (or both zero), sort by symbol
		if valueI == valueJ {
			assetI := m.getAssetByID(assetIDs[i])
			assetJ := m.getAssetByID(assetIDs[j])
			return assetI.Symbol < assetJ.Symbol
		}

		return valueI > valueJ
	})

	// Build rows with tree structure
	for _, assetID := range assetIDs {
		holdings := assetHoldings[assetID]
		asset := m.getAssetByID(assetID)
		price := m.prices[assetID]
		totalValue := assetTotalValues[assetID]

		// Calculate total amount for this asset
		totalAmount := 0.0
		for _, holding := range holdings {
			totalAmount += holding.Amount
		}

		// Format total amount based on asset type
		amountStr := ""
		if asset.Type == models.AssetTypeFiat {
			amountStr = fmt.Sprintf("%.2f", totalAmount)
		} else {
			amountStr = fmt.Sprintf("%.4f", totalAmount)
		}

		// Add asset header row
		assetRow := table.Row{
			asset.Symbol,
			amountStr,
			fmt.Sprintf("$%.2f", totalValue),
		}
		rows = append(rows, assetRow)

		// Sort holdings within each asset by value (highest first)
		sort.Slice(holdings, func(i, j int) bool {
			valueI := holdings[i].Amount * price
			valueJ := holdings[j].Amount * price
			return valueI > valueJ
		})

		// Add account rows
		for i, holding := range holdings {
			account := m.getAccountByID(holding.AccountID)
			value := holding.Amount * price

			// Determine tree character
			var treeChar string
			if i == len(holdings)-1 {
				treeChar = "└─ "
			} else {
				treeChar = "├─ "
			}

			// Format amount based on asset type
			amountStr := ""
			if asset.Type == models.AssetTypeFiat {
				amountStr = fmt.Sprintf("%.2f", holding.Amount)
			} else {
				amountStr = fmt.Sprintf("%.4f", holding.Amount)
			}

			row := table.Row{
				"  " + treeChar + account.Name,
				amountStr,
				fmt.Sprintf("$%.2f", value),
			}
			rows = append(rows, row)
		}
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

	// Header with last update time
	headerLeft := "💰 Minimal Money"
	headerRight := fmt.Sprintf("Total: $%.2f", total)
	headerPadding := m.width - len(headerLeft) - len(headerRight) - 2
	if headerPadding < 1 {
		headerPadding = 1
	}
	header := headerLeft + strings.Repeat(" ", headerPadding) + headerRight
	b.WriteString(totalStyle.Render(header) + "\n")

	// Last update time
	if m.lastPriceUpdate != nil {
		updateText := fmt.Sprintf("Last Update: %s", m.lastPriceUpdate.Format("2006-01-02 15:04:05"))
		updatePadding := m.width - len(updateText) - 2
		if updatePadding < 0 {
			updatePadding = 0
		}
		updateTime := strings.Repeat(" ", updatePadding) + updateText
		b.WriteString(updateTime + "\n\n")
	} else {
		b.WriteString("\n")
	}

	// Table
	b.WriteString(baseStyle.Render(m.table.View()) + "\n\n")

	// Footer
	footer := "[n]ew  [e]dit  [d]elete  [p]rice update  [h]istory  [q]uit"
	b.WriteString(footer)

	return b.String()
}
