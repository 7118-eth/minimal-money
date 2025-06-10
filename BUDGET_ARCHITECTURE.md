# Budget Tracking Architecture

## Overview

Transform the current Asset Tracker into a comprehensive Personal Finance Manager by adding budget tracking capabilities while maintaining the existing asset tracking functionality.

## Core Concepts

### 1. Transactions
The foundation of budget tracking - recording money flow in and out.

```go
type Transaction struct {
    ID            uint
    Type          TransactionType  // income, expense, transfer
    Amount        float64
    Currency      string          // AED, USD, EUR, etc.
    CategoryID    uint
    AccountID     uint           // Links to existing Account model
    Description   string
    Date          time.Time
    IsRecurring   bool
    RecurrenceID  uint           // Links to RecurrenceRule
    Tags          []string       // Flexible labeling
}

type TransactionType string
const (
    TransactionIncome   TransactionType = "income"
    TransactionExpense  TransactionType = "expense"
    TransactionTransfer TransactionType = "transfer"
)
```

### 2. Categories
Organize transactions for better insights.

```go
type Category struct {
    ID          uint
    Name        string        // "Rent", "Groceries", "Salary"
    Type        TransactionType
    Icon        string        // Emoji for UI
    Color       string
    ParentID    *uint        // For subcategories
    IsSystem    bool         // System categories can't be deleted
}
```

### 3. Budgets
Set spending limits and track progress.

```go
type Budget struct {
    ID          uint
    Name        string
    Period      BudgetPeriod  // monthly, quarterly, yearly
    StartDate   time.Time
    EndDate     time.Time
    CategoryID  uint
    Amount      float64       // Budget limit in USD
    Currency    string        // Original currency
}

type BudgetPeriod string
const (
    BudgetMonthly   BudgetPeriod = "monthly"
    BudgetQuarterly BudgetPeriod = "quarterly"
    BudgetYearly    BudgetPeriod = "yearly"
    BudgetCustom    BudgetPeriod = "custom"
)
```

### 4. Recurrence Rules
Handle recurring transactions efficiently.

```go
type RecurrenceRule struct {
    ID          uint
    Frequency   RecurrenceFrequency  // daily, weekly, monthly, yearly
    Interval    int                  // Every N frequency units
    EndDate     *time.Time          // Optional end date
    NextDate    time.Time           // Next occurrence
}
```

## UI/UX Design

### Navigation Structure
```
Asset Tracker (current view)
├── Portfolio Overview
├── Add/Edit/Delete Assets
└── Price Refresh

Budget Tracker (new view - press 'b')
├── Dashboard
│   ├── Month Summary (Income vs Expenses)
│   ├── Budget Progress Bars
│   └── Recent Transactions
├── Transactions (press 't')
│   ├── List View with Filters
│   ├── Add Transaction (press 'n')
│   └── Edit/Delete (press 'e'/'d')
├── Categories (press 'c')
│   ├── Income Categories
│   ├── Expense Categories
│   └── Manage Categories
├── Budgets (press 'u')
│   ├── Active Budgets
│   ├── Create Budget
│   └── Budget vs Actual Report
└── Reports (press 'r')
    ├── Monthly Summary
    ├── Category Breakdown
    ├── Trends Over Time
    └── Net Worth (Assets - Liabilities)
```

### Key Bindings
- `a` - Switch to Asset Tracker
- `b` - Switch to Budget Tracker
- `n` - New (context-aware: asset/transaction/budget)
- `e` - Edit selected item
- `d` - Delete selected item
- `f` - Filter/Search
- `Tab` - Navigate between sections

## Implementation Phases

### Phase 1: Core Transaction System
1. Create transaction models and database tables
2. Implement transaction CRUD operations
3. Add basic transaction list view
4. Create "Add Transaction" modal
5. Link transactions to existing accounts

### Phase 2: Categories & Organization
1. Implement category system with defaults
2. Add category management UI
3. Enable transaction categorization
4. Add transaction search and filters

### Phase 3: Budget Management
1. Create budget models
2. Implement budget creation/editing
3. Add budget vs actual calculations
4. Create budget progress visualization

### Phase 4: Reporting & Analytics
1. Monthly/yearly summary reports
2. Category breakdown charts (ASCII)
3. Spending trends analysis
4. Net worth tracking (assets - expenses)

### Phase 5: Advanced Features
1. Recurring transactions
2. Multi-currency budget tracking
3. Savings goals
4. Bill reminders
5. Export functionality (CSV)

## Data Integration

### Linking Assets and Budget
- Asset purchases/sales create transactions automatically
- Investment income tracked as income transactions
- Net worth = Total Assets - Total Liabilities

### Currency Handling
- All amounts stored in original currency
- Display currency preference (default USD)
- Use existing exchange rate system
- Historical rates for past transactions

## Database Schema

### New Tables
```sql
-- Categories
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    icon TEXT,
    color TEXT,
    parent_id INTEGER,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);

-- Transactions
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    category_id INTEGER,
    account_id INTEGER,
    description TEXT,
    date TIMESTAMP NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    recurrence_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (recurrence_id) REFERENCES recurrence_rules(id)
);

-- Transaction Tags
CREATE TABLE transaction_tags (
    transaction_id INTEGER,
    tag TEXT,
    PRIMARY KEY (transaction_id, tag),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- Budgets
CREATE TABLE budgets (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    period TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    category_id INTEGER,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Recurrence Rules
CREATE TABLE recurrence_rules (
    id INTEGER PRIMARY KEY,
    frequency TEXT NOT NULL,
    interval INTEGER NOT NULL,
    end_date TIMESTAMP,
    next_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## Default Categories

### Income Categories
- 💼 Salary
- 💰 Investment Income
- 🎯 Freelance/Contract
- 🎁 Gifts
- 💸 Other Income

### Expense Categories
- 🏠 Housing (Rent/Mortgage)
- 🍔 Food & Dining
- 🚗 Transportation
- 🛒 Shopping
- 💊 Healthcare
- 🎮 Entertainment
- 📚 Education
- 💡 Utilities
- 📱 Subscriptions
- ✈️ Travel
- 💳 Debt Payments
- 🏦 Savings/Investments
- 🎯 Other Expenses

## Technical Considerations

1. **Performance**
   - Index on transaction date and category
   - Cache budget calculations
   - Pagination for transaction lists

2. **Data Integrity**
   - Transaction amounts must be positive
   - Enforce currency consistency within budgets
   - Soft deletes for audit trail

3. **User Experience**
   - Quick entry shortcuts (e.g., "50 AED groceries")
   - Smart date parsing ("yesterday", "last monday")
   - Autocomplete for descriptions
   - Remember last used category

4. **Future Extensibility**
   - API-ready structure
   - Plugin system for bank imports
   - Mobile app compatibility
   - Multi-user/family budgets

## Success Metrics

1. Can track daily expenses in under 10 seconds
2. Monthly budget overview at a glance
3. Clear visual feedback on overspending
4. Seamless integration with asset tracking
5. Accurate multi-currency calculations

## Next Steps

1. Review and refine architecture
2. Create transaction models and repositories
3. Build transaction entry UI
4. Implement basic reporting
5. Add budget management features

This architecture provides a solid foundation for comprehensive personal finance management while maintaining the simplicity and speed of the current asset tracker.