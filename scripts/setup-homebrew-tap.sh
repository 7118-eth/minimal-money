#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Homebrew Tap Setup for minimal-money${NC}"
echo "======================================"

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed${NC}"
    echo "Install it with: brew install gh"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${RED}Error: Not authenticated with GitHub${NC}"
    echo "Run: gh auth login"
    exit 1
fi

# Step 1: Create tap repository
echo -e "\n${YELLOW}Step 1: Creating tap repository...${NC}"
if gh repo view 7118-eth/homebrew-minimal-money &> /dev/null; then
    echo "Repository already exists!"
else
    gh repo create 7118-eth/homebrew-minimal-money \
        --public \
        --description "Homebrew tap for minimal-money" \
        --clone
    cd homebrew-minimal-money
fi

# Step 2: Create Formula directory
echo -e "\n${YELLOW}Step 2: Setting up Formula directory...${NC}"
mkdir -p Formula

# Step 3: Get latest release info
echo -e "\n${YELLOW}Step 3: Getting latest release information...${NC}"
LATEST_TAG=$(gh release list --repo 7118-eth/minimal-money --limit 1 | cut -f3)
if [ -z "$LATEST_TAG" ]; then
    echo -e "${RED}Error: No releases found${NC}"
    echo "Please create a release first"
    exit 1
fi

VERSION=${LATEST_TAG#v}
echo "Latest release: $LATEST_TAG (version $VERSION)"

# Step 4: Download and extract checksums
echo -e "\n${YELLOW}Step 4: Downloading checksums...${NC}"
gh release download $LATEST_TAG --repo 7118-eth/minimal-money --pattern "*.sha256" --clobber || {
    echo -e "${RED}Error: Could not download checksums${NC}"
    echo "Are the release assets public?"
    exit 1
}

# Extract checksums
echo "Extracting checksums..."
DARWIN_ARM64_SHA=$(cat minimal-money_${VERSION}_darwin_arm64.tar.gz.sha256 2>/dev/null | cut -d' ' -f1 || echo "MISSING")
DARWIN_AMD64_SHA=$(cat minimal-money_${VERSION}_darwin_amd64.tar.gz.sha256 2>/dev/null | cut -d' ' -f1 || echo "MISSING")
LINUX_ARM64_SHA=$(cat minimal-money_${VERSION}_linux_arm64.tar.gz.sha256 2>/dev/null | cut -d' ' -f1 || echo "MISSING")
LINUX_AMD64_SHA=$(cat minimal-money_${VERSION}_linux_amd64.tar.gz.sha256 2>/dev/null | cut -d' ' -f1 || echo "MISSING")

# Step 5: Create formula from template
echo -e "\n${YELLOW}Step 5: Creating formula...${NC}"
cat > Formula/minimal-money.rb << EOF
class MinimalMoney < Formula
  desc "Beautiful terminal-based portfolio tracker that respects your time and privacy"
  homepage "https://github.com/7118-eth/minimal-money"
  version "${VERSION}"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/7118-eth/minimal-money/releases/download/v${VERSION}/minimal-money_${VERSION}_darwin_arm64.tar.gz"
      sha256 "${DARWIN_ARM64_SHA}"
    else
      url "https://github.com/7118-eth/minimal-money/releases/download/v${VERSION}/minimal-money_${VERSION}_darwin_amd64.tar.gz"
      sha256 "${DARWIN_AMD64_SHA}"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/7118-eth/minimal-money/releases/download/v${VERSION}/minimal-money_${VERSION}_linux_arm64.tar.gz"
      sha256 "${LINUX_ARM64_SHA}"
    else
      url "https://github.com/7118-eth/minimal-money/releases/download/v${VERSION}/minimal-money_${VERSION}_linux_amd64.tar.gz"
      sha256 "${LINUX_AMD64_SHA}"
    end
  end

  def install
    bin.install "minimal-money"
  end

  test do
    assert_match "Minimal Money", shell_output("#{bin}/minimal-money --version")
  end
end
EOF

# Clean up checksum files
rm -f *.sha256

# Step 6: Create update script
echo -e "\n${YELLOW}Step 6: Creating update script...${NC}"
cat > update-formula.sh << 'SCRIPT'
#!/bin/bash
set -e

REPO="7118-eth/minimal-money"
LATEST_TAG=$(gh release list --repo $REPO --limit 1 | cut -f3)
VERSION=${LATEST_TAG#v}

echo "Updating to version $VERSION..."

# Download checksums
gh release download $LATEST_TAG --repo $REPO --pattern "*.sha256" --clobber

# Extract checksums and update formula
for platform in darwin_arm64 darwin_amd64 linux_arm64 linux_amd64; do
    SHA=$(cat minimal-money_${VERSION}_${platform}.tar.gz.sha256 2>/dev/null | cut -d' ' -f1 || echo "MISSING")
    # Update both version and SHA in URL and sha256 lines
    sed -i '' "s|/v[0-9.]\+/minimal-money_[0-9.]\+_${platform}|/v${VERSION}/minimal-money_${VERSION}_${platform}|g" Formula/minimal-money.rb
    sed -i '' "/_${platform}\.tar\.gz\"/,/sha256/ s/sha256 \"[^\"]*\"/sha256 \"${SHA}\"/" Formula/minimal-money.rb
done

# Update version line
sed -i '' "s/version \"[^\"]*\"/version \"${VERSION}\"/" Formula/minimal-money.rb

# Clean up
rm -f *.sha256

echo "Formula updated to version $VERSION"
echo "Don't forget to commit and push!"
SCRIPT

chmod +x update-formula.sh

# Step 7: Commit and push
echo -e "\n${YELLOW}Step 7: Committing formula...${NC}"
git add Formula/minimal-money.rb update-formula.sh
git commit -m "Add minimal-money formula v${VERSION}"

echo -e "\n${YELLOW}Step 8: Pushing to GitHub...${NC}"
git push origin main 2>/dev/null || git push --set-upstream origin main

# Done!
echo -e "\n${GREEN}✅ Homebrew tap setup complete!${NC}"
echo -e "\nUsers can now install with:"
echo -e "  ${YELLOW}brew tap 7118-eth/minimal-money${NC}"
echo -e "  ${YELLOW}brew install minimal-money${NC}"

echo -e "\nTo update the formula for new releases:"
echo -e "  ${YELLOW}cd homebrew-minimal-money${NC}"
echo -e "  ${YELLOW}./update-formula.sh${NC}"

echo -e "\n${GREEN}Testing the tap locally:${NC}"
echo -e "  ${YELLOW}brew tap 7118-eth/minimal-money${NC}"
echo -e "  ${YELLOW}brew install minimal-money${NC}"
echo -e "  ${YELLOW}minimal-money --version${NC}"