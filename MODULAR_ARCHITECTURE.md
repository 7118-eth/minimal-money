# Modular Architecture Decision

## Why Separate Applications?

After implementing the asset tracker and planning the budget tracker, we've decided to build them as separate applications rather than a monolithic "Personal Finance Manager". This decision is based on both technical and practical considerations.

## Benefits of Separation

### 1. Clarity for AI Development
- Each AI agent works on a focused, single-purpose application
- No confusion about which features belong where
- Clear boundaries prevent accidental coupling
- Easier to verify "all features work" in smaller scope

### 2. Maintainability
- Changes to asset tracking don't risk breaking budget features
- Each app can evolve independently
- Simpler codebases are easier to understand
- Less cognitive load when working on either app

### 3. User Experience
- Users who only want asset tracking aren't burdened with budget features
- Faster startup times
- Smaller binaries
- Can run both simultaneously if needed

### 4. Testing
- Focused test suites
- Faster test execution
- Clearer test coverage metrics
- No cross-contamination of test data

## Current Structure

```
asset-tracker/          # Current repository
├── cmd/
│   └── budget/        # Will be renamed to 'assets'
├── internal/
│   ├── models/        # Asset-specific models
│   ├── repository/    # Asset data access
│   ├── service/       # Price fetching
│   └── ui/            # Asset UI
└── README.md

budget-tracker/         # New, separate repository
├── cmd/
│   └── budget/
├── internal/
│   ├── models/        # Transaction, Category, Budget
│   ├── repository/    # Budget data access
│   ├── service/       # Budget calculations
│   └── ui/            # Budget UI
└── README.md
```

## Shared Concepts

While separate, both applications share some concepts that could be extracted later:

### Common Models
- Currency conversion rates
- Date/time utilities
- Money formatting

### Potential Shared Package
```go
// github.com/user/finance-common
package common

type Currency string
type Money struct {
    Amount   float64
    Currency Currency
}

func ConvertToUSD(m Money, rates map[Currency]float64) Money
func FormatMoney(m Money) string
```

## Future Integration Options

### 1. Launcher Application
A simple menu to choose which app to run:
```
💰 Personal Finance Manager

[1] Asset Tracker
[2] Budget Tracker
[3] Combined Dashboard (future)

Choose: _
```

### 2. Data Sharing
- Export/import between applications
- Shared currency rate cache
- Common configuration file

### 3. Combined Reporting
A third application that reads from both databases:
- Net worth = Assets - Liabilities
- Investment income tracked in budget
- Asset purchases reflected as expenses

## Migration Path

### Current State
- Asset tracker is complete and functional
- Budget tracker is planned but not started

### Next Steps
1. Rename current binary from 'budget' to 'assets'
2. Update documentation to reflect asset-only focus
3. Create new repository for budget tracker
4. Implement budget tracker independently

### Future Consideration
After both apps are stable:
- Evaluate user feedback
- Consider creating finance-common package
- Potentially build unified dashboard

## Design Principles

### Do One Thing Well
- Asset tracker: Track portfolio value
- Budget tracker: Track income/expenses
- Each excels at its purpose

### Loose Coupling
- No direct dependencies between apps
- Communication through files/exports if needed
- Independent release cycles

### User Choice
- Users can adopt one or both
- No forced complexity
- Gradual adoption path

## Conclusion

Separating the applications provides immediate benefits for development, testing, and maintenance. It aligns with Unix philosophy of "do one thing well" and makes the codebase more manageable for both humans and AI agents.

The modular approach doesn't prevent future integration - it actually makes it easier by ensuring clean boundaries from the start.