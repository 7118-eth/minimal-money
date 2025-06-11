#!/bin/bash
# Install git hooks for Minimal Money

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔧 Installing git hooks..."

# Configure git to use our hooks directory
git config core.hooksPath .githooks

echo "${GREEN}✅ Git hooks installed successfully!${NC}"
echo ""
echo "The following hooks are now active:"
echo "  • pre-commit: Runs formatting and linting checks"
echo ""
echo "${YELLOW}💡 To disable hooks temporarily, use: git commit --no-verify${NC}"
echo "${YELLOW}💡 To uninstall hooks, run: git config --unset core.hooksPath${NC}"