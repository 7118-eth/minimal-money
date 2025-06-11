package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bioharz/budget/internal/models"
)

func (m Model) formatHoldingChange(log models.AuditLog) string {
	var result strings.Builder

	switch log.Action {
	case models.AuditActionCreate:
		// Parse new value
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(log.NewValue), &data); err == nil {
			accountID := uint(data["account_id"].(float64))
			assetID := uint(data["asset_id"].(float64))
			amount := data["amount"].(float64)

			account := m.getAccountByID(accountID)
			asset := m.getAssetByID(assetID)

			result.WriteString(fmt.Sprintf("  Added %.4f %s to %s\n", amount, asset.Symbol, account.Name))
			if purchasePrice, ok := data["purchase_price"].(float64); ok && purchasePrice > 0 {
				result.WriteString(fmt.Sprintf("  Purchase price: $%.2f\n", purchasePrice))
			}
		}

	case models.AuditActionUpdate:
		// Parse old and new values
		var oldData, newData map[string]interface{}
		if err := json.Unmarshal([]byte(log.OldValue), &oldData); err == nil {
			if err := json.Unmarshal([]byte(log.NewValue), &newData); err == nil {
				// Get the account and asset info
				accountID := uint(newData["account_id"].(float64))
				assetID := uint(newData["asset_id"].(float64))
				account := m.getAccountByID(accountID)
				asset := m.getAssetByID(assetID)

				result.WriteString(fmt.Sprintf("  Updated %s in %s:\n", asset.Symbol, account.Name))

				// Check what changed
				oldAmount := oldData["amount"].(float64)
				newAmount := newData["amount"].(float64)
				if oldAmount != newAmount {
					result.WriteString(fmt.Sprintf("  Amount: %.4f → %.4f\n", oldAmount, newAmount))
				}

				// Check if account changed
				oldAccountID := uint(oldData["account_id"].(float64))
				if oldAccountID != accountID {
					oldAccount := m.getAccountByID(oldAccountID)
					result.WriteString(fmt.Sprintf("  Moved from: %s → %s\n", oldAccount.Name, account.Name))
				}
			}
		}

	case models.AuditActionDelete:
		// Parse old value
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(log.OldValue), &data); err == nil {
			accountID := uint(data["account_id"].(float64))
			assetID := uint(data["asset_id"].(float64))
			amount := data["amount"].(float64)

			account := m.getAccountByID(accountID)
			asset := m.getAssetByID(assetID)

			result.WriteString(fmt.Sprintf("  Removed %.4f %s from %s\n", amount, asset.Symbol, account.Name))
		}
	}

	return result.String()
}
