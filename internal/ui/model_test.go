package ui

import (
	"testing"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/helpers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	assert.Equal(t, ViewMain, model.view)
	assert.NotNil(t, model.prices)
	assert.NotNil(t, model.accounts)
	assert.NotNil(t, model.assets)
	assert.NotNil(t, model.holdings)
	assert.NotNil(t, model.priceService)
	assert.NotNil(t, model.table)
}

func TestModel_ViewTransitions(t *testing.T) {
	tests := []struct {
		name         string
		startView    View
		key          string
		expectedView View
		expectModal  bool
	}{
		{
			name:         "open add asset modal",
			startView:    ViewMain,
			key:          "n",
			expectedView: ViewAddAsset,
			expectModal:  true,
		},
		{
			name:         "open history view",
			startView:    ViewMain,
			key:          "h",
			expectedView: ViewHistory,
			expectModal:  false,
		},
		{
			name:         "escape from add asset",
			startView:    ViewAddAsset,
			key:          "esc",
			expectedView: ViewMain,
			expectModal:  false,
		},
		{
			name:         "escape from history",
			startView:    ViewHistory,
			key:          "esc",
			expectedView: ViewMain,
			expectModal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel()
			model.view = tt.startView
			if tt.startView == ViewAddAsset {
				model.inputMode = true
			}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "esc" {
				msg = tea.KeyMsg{Type: tea.KeyEscape}
			}

			newModel, _ := model.Update(msg)
			m := newModel.(Model)

			assert.Equal(t, tt.expectedView, m.view)
			assert.Equal(t, tt.expectModal, m.inputMode)
		})
	}
}

func TestModel_QuitHandling(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		inputMode  bool
		shouldQuit bool
	}{
		{
			name:       "quit with q",
			key:        "q",
			inputMode:  false,
			shouldQuit: true,
		},
		{
			name:       "quit with ctrl+c",
			key:        "ctrl+c",
			inputMode:  false,
			shouldQuit: true,
		},
		{
			name:       "don't quit in input mode",
			key:        "q",
			inputMode:  true,
			shouldQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel()
			model.inputMode = tt.inputMode

			var msg tea.Msg
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			_, cmd := model.Update(msg)

			if tt.shouldQuit {
				assert.NotNil(t, cmd)
				// Check if it's a quit command by running it
				msg = cmd()
				_, ok := msg.(tea.QuitMsg)
				assert.True(t, ok)
			} else {
				if cmd != nil {
					msg = cmd()
					_, ok := msg.(tea.QuitMsg)
					assert.False(t, ok)
				}
			}
		})
	}
}

func TestModel_ModalInput(t *testing.T) {
	model := InitialModel()
	model.view = ViewAddAsset
	model.inputMode = true
	model.initAddAssetModal()

	// Type some text
	inputs := []string{"h", "a", "r", "d", "w", "a", "r", "e", " ", "w", "a", "l", "l", "e", "t"}
	for _, char := range inputs {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(char)}
		newModel, _ := model.Update(msg)
		model = newModel.(Model)
	}

	assert.Equal(t, "hardware wallet", model.modalState.Fields[0].Value)

	// Test tab navigation
	msg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := model.Update(msg)
	model = newModel.(Model)
	assert.Equal(t, 1, model.modalState.ActiveField)

	// Type in second field
	model.modalState.ActiveField = 1
	inputs = []string{"B", "T", "C"}
	for _, char := range inputs {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(char)}
		newModel, _ := model.Update(msg)
		model = newModel.(Model)
	}
	assert.Equal(t, "BTC", model.modalState.Fields[1].Value)

	// Test backspace
	msg = tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ = model.Update(msg)
	model = newModel.(Model)
	assert.Equal(t, "BT", model.modalState.Fields[1].Value)
}

func TestModel_WindowResize(t *testing.T) {
	model := InitialModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
	// Table height should be adjusted (height - 8, but table may adjust it further)
	assert.LessOrEqual(t, m.table.Height(), 32)
}

func TestModel_RefreshPrices(t *testing.T) {
	db := helpers.SetupTestDB(t)
	model := InitialModelWithDB(db)

	// Add some assets
	model.assets = []models.Asset{
		{ID: 1, Symbol: "BTC", Type: models.AssetTypeCrypto},
		{ID: 2, Symbol: "USD", Type: models.AssetTypeFiat},
	}

	// Press 'p' to update prices
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")}
	_, cmd := model.Update(msg)

	// Should return a command
	assert.NotNil(t, cmd)

	// The command should fetch prices (we can't test the actual fetch without mocking)
	// But we can verify it returns a priceUpdateMsg
	result := cmd()
	_, ok := result.(priceUpdateMsg)
	assert.True(t, ok)
}

func TestModel_PriceUpdate(t *testing.T) {
	model := InitialModel()

	// Add some holdings for the table
	model.holdings = []models.Holding{
		{ID: 1, AccountID: 1, AssetID: 1, Amount: 0.5},
	}
	model.accounts = []models.Account{
		{ID: 1, Name: "Test Account"},
	}
	model.assets = []models.Asset{
		{ID: 1, Symbol: "BTC"},
	}

	// Send price update message
	prices := map[uint]float64{
		1: 50000.0,
	}
	msg := priceUpdateMsg{prices: prices}

	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	assert.Equal(t, 50000.0, m.prices[1])
}

func TestModel_ViewRendering(t *testing.T) {
	tests := []struct {
		name     string
		view     View
		contains string
	}{
		{
			name:     "main view shows budget tracker",
			view:     ViewMain,
			contains: "Minimal Money",
		},
		{
			name:     "assets view",
			view:     ViewAssets,
			contains: "Your Assets",
		},
		{
			name:     "add asset view",
			view:     ViewAddAsset,
			contains: "Add New Asset",
		},
		{
			name:     "history view",
			view:     ViewHistory,
			contains: "Audit Trail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var model Model
			if tt.view == ViewHistory {
				db := helpers.SetupTestDB(t)
				model = InitialModelWithDB(db)
			} else {
				model = InitialModel()
			}
			model.view = tt.view
			if tt.view == ViewAddAsset {
				model.initAddAssetModal()
			}

			output := model.View()
			assert.Contains(t, output, tt.contains)
		})
	}
}
