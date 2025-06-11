# Homebrew Tap Setup Guide

This guide explains how to set up the Homebrew tap for minimal-money.

## Prerequisites
- GitHub account with access to create repositories
- A published release with binaries and checksums

## Step 1: Create the Tap Repository

1. Create a new repository named `homebrew-minimal-money` on GitHub:
   ```bash
   gh repo create 7118-eth/homebrew-minimal-money --public --description "Homebrew tap for minimal-money"
   ```

2. Clone the repository:
   ```bash
   git clone git@github.com:7118-eth/homebrew-minimal-money.git
   cd homebrew-minimal-money
   ```

3. Create the Formula directory:
   ```bash
   mkdir -p Formula
   ```

## Step 2: Get Release Information

1. Get the latest release info:
   ```bash
   # Get latest release tag
   LATEST_TAG=$(gh release list --repo 7118-eth/minimal-money --limit 1 | cut -f3)
   echo "Latest release: $LATEST_TAG"
   
   # Download checksums
   gh release download $LATEST_TAG --repo 7118-eth/minimal-money --pattern "*.sha256"
   ```

2. Extract checksums for each platform:
   ```bash
   # Example commands to get checksums
   DARWIN_ARM64_SHA=$(cat minimal-money_*_darwin_arm64.tar.gz.sha256 | cut -d' ' -f1)
   DARWIN_AMD64_SHA=$(cat minimal-money_*_darwin_amd64.tar.gz.sha256 | cut -d' ' -f1)
   LINUX_ARM64_SHA=$(cat minimal-money_*_linux_arm64.tar.gz.sha256 | cut -d' ' -f1)
   LINUX_AMD64_SHA=$(cat minimal-money_*_linux_amd64.tar.gz.sha256 | cut -d' ' -f1)
   ```

## Step 3: Create the Formula

1. Copy the template and update it:
   ```bash
   # Copy from minimal-money repo
   cp ../minimal-money/packaging/homebrew/minimal-money.rb.template Formula/minimal-money.rb
   
   # Update version and checksums
   VERSION=${LATEST_TAG#v}  # Remove 'v' prefix
   
   # Use sed or manually edit the file to replace placeholders
   sed -i '' "s/VERSION_PLACEHOLDER/$VERSION/g" Formula/minimal-money.rb
   sed -i '' "s/SHA256_DARWIN_ARM64_PLACEHOLDER/$DARWIN_ARM64_SHA/g" Formula/minimal-money.rb
   sed -i '' "s/SHA256_DARWIN_AMD64_PLACEHOLDER/$DARWIN_AMD64_SHA/g" Formula/minimal-money.rb
   sed -i '' "s/SHA256_LINUX_ARM64_PLACEHOLDER/$LINUX_ARM64_SHA/g" Formula/minimal-money.rb
   sed -i '' "s/SHA256_LINUX_AMD64_PLACEHOLDER/$LINUX_AMD64_SHA/g" Formula/minimal-money.rb
   ```

## Step 4: Test the Formula Locally

1. Test the formula:
   ```bash
   # Test installation
   brew install --build-from-source Formula/minimal-money.rb
   
   # Test the installed binary
   minimal-money --version
   
   # Run brew tests
   brew test Formula/minimal-money.rb
   
   # Audit the formula
   brew audit --strict Formula/minimal-money.rb
   ```

## Step 5: Commit and Push

1. Add and commit the formula:
   ```bash
   git add Formula/minimal-money.rb
   git commit -m "Add minimal-money formula v$VERSION"
   git push origin main
   ```

## Step 6: Test the Tap

Users can now install via:
```bash
# Add the tap
brew tap 7118-eth/minimal-money

# Install
brew install minimal-money
```

## Automation with GoReleaser

Once the tap repository exists, GoReleaser can automatically update it on each release:

1. Create a GitHub token with repo access to the tap
2. Add it as a secret: `HOMEBREW_TAP_GITHUB_TOKEN`
3. GoReleaser will handle formula updates automatically

## Manual Update Process

For each new release:
1. Download the new checksums
2. Update Formula/minimal-money.rb with new version and checksums
3. Commit and push to the tap repository

## Example Script for Updates

Create `update-formula.sh` in the tap repository:

```bash
#!/bin/bash
set -e

REPO="7118-eth/minimal-money"
LATEST_TAG=$(gh release list --repo $REPO --limit 1 | cut -f3)
VERSION=${LATEST_TAG#v}

echo "Updating to version $VERSION..."

# Download checksums
gh release download $LATEST_TAG --repo $REPO --pattern "*.sha256" --clobber

# Extract checksums
DARWIN_ARM64_SHA=$(grep darwin_arm64 *.sha256 | cut -d' ' -f1)
DARWIN_AMD64_SHA=$(grep darwin_amd64 *.sha256 | cut -d' ' -f1)
LINUX_ARM64_SHA=$(grep linux_arm64 *.sha256 | cut -d' ' -f1)
LINUX_AMD64_SHA=$(grep linux_amd64 *.sha256 | cut -d' ' -f1)

# Update formula
sed -i '' "s/version \".*\"/version \"$VERSION\"/" Formula/minimal-money.rb
sed -i '' "s/sha256 \".*\" # darwin_arm64/sha256 \"$DARWIN_ARM64_SHA\" # darwin_arm64/" Formula/minimal-money.rb
sed -i '' "s/sha256 \".*\" # darwin_amd64/sha256 \"$DARWIN_AMD64_SHA\" # darwin_amd64/" Formula/minimal-money.rb
sed -i '' "s/sha256 \".*\" # linux_arm64/sha256 \"$LINUX_ARM64_SHA\" # linux_arm64/" Formula/minimal-money.rb
sed -i '' "s/sha256 \".*\" # linux_amd64/sha256 \"$LINUX_AMD64_SHA\" # linux_amd64/" Formula/minimal-money.rb

# Clean up
rm *.sha256

echo "Formula updated to version $VERSION"
```

## Troubleshooting

### "SHA256 mismatch"
- Ensure you're using the correct checksums from the release
- Check that the URL pattern matches the actual release assets

### "Formula not found"
- Make sure the repository is public
- Check that Formula directory exists with correct capitalization

### "Download failed"
- Verify the release assets exist and are public
- Check the URL construction in the formula