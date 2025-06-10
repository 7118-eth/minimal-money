package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	account := strings.TrimSpace(m.modalState.Fields[0].Value)
	asset := strings.TrimSpace(m.modalState.Fields[1].Value)
	amountStr := strings.TrimSpace(m.modalState.Fields[2].Value)
	priceStr := strings.TrimSpace(m.modalState.Fields[3].Value)

	if account == "" || asset == "" || amountStr == "" {
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

	// TODO: Save to database
	// For now, just close the modal
	m.view = ViewMain
	m.inputMode = false
	m.modalState = ModalState{}
}

func (m *Model) renderAddAssetModal() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add New Asset") + "\n\n")

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
			if m.inputMode && m.inputBuffer != "" {
				value = m.inputBuffer
			}
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