# Budget Tracker - Development Progress

## Current State (2025-01-06)

### ✅ Implemented
- Basic project structure with Go modules
- Database models (Asset, Holding, PriceHistory, PortfolioSnapshot)
- SQLite connection with GORM
- Basic Bubble Tea UI with menu navigation
- Text input handling for asset symbols
- ESC key navigation between views
- API client structure with caching

### ⚠️ Partially Implemented
- Asset creation (can enter symbol, but can't save)
- Navigation system (works but needs table view)
- Price fetching (structure exists, not connected)

### ❌ Not Implemented
- **Accounts concept** - No way to organize by platform
- **Asset amounts** - Can only enter symbol, not quantity
- **Purchase price tracking** - No P&L calculations
- **Table view** - Still using menu system
- **Price updates** - API not connected to UI
- **Data persistence** - Input not saved to database
- **History view** - Empty placeholder
- **Edit/Delete** - No way to modify holdings

## Known Issues
1. **Can't add asset values** - Only symbol input exists
2. **No account selection** - Missing account management
3. **No data persistence** - Entries aren't saved
4. **No price display** - API not integrated with UI
5. **Menu-based navigation** - Should be table-based

## Next Steps (Priority Order)

### 1. Update Data Models
- [ ] Add Account model
- [ ] Update Holding to include AccountID
- [ ] Add PurchasePrice and PurchaseDate fields

### 2. Implement Table View
- [ ] Replace menu system with table component
- [ ] Show all holdings in main view
- [ ] Add keyboard navigation for table

### 3. Create Input Modal
- [ ] Account selection dropdown
- [ ] Asset selection/creation
- [ ] Amount input field
- [ ] Purchase price field

### 4. Connect Database
- [ ] Save new holdings
- [ ] Load holdings on startup
- [ ] Update operations

### 5. Integrate Price API
- [ ] Fetch prices on startup
- [ ] Background price updates
- [ ] Calculate portfolio value

## Decision Log

### 2025-01-06: Initial UI Approach
- **Decision**: Start with menu-based navigation
- **Reason**: Simpler to implement initially
- **Result**: Too limiting, need table view

### 2025-01-06: Added Accounts Concept
- **Decision**: Add Account model to group holdings
- **Reason**: Users organize assets by platform (NeoBank, hardware wallet, etc.)
- **Result**: Better matches real-world usage

### 2025-01-06: Switch to Table UI
- **Decision**: Replace menu with always-visible table
- **Reason**: Better overview, faster workflow
- **Result**: (Pending implementation)

## Testing Notes
- Manual testing only (no MCP server yet)
- User provides feedback on UI behavior
- Iterative development based on testing

## Questions for User
1. Should we add categories (e.g., "Crypto", "Cash", "Stocks")?
2. Do you want transaction history or just current holdings?
3. Should we support multiple currencies or convert everything to USD?
4. Any specific exchanges/wallets we should prioritize?