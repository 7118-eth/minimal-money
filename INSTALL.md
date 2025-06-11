# Installation Guide

Minimal Money can be installed through various methods depending on your platform and preferences.

## Quick Install

### macOS/Linux (Recommended)
```bash
# Download latest release
curl -L https://github.com/7118-eth/minimal-money/releases/latest/download/minimal-money_$(uname -s)_$(uname -m).tar.gz | tar xz
chmod +x minimal-money
sudo mv minimal-money /usr/local/bin/
```

### Windows
Download the latest `.zip` file from [releases](https://github.com/7118-eth/minimal-money/releases) and extract it to a location in your PATH.

## Package Managers

### Homebrew (macOS/Linux)
```bash
# Once tap repository is available
brew tap 7118-eth/minimal-money
brew install minimal-money
```

### Snap (Linux)
```bash
# Once published to Snap Store
sudo snap install minimal-money
```

### AUR (Arch Linux)
```bash
# Using yay
yay -S minimal-money

# Using paru
paru -S minimal-money

# Manual installation
git clone https://aur.archlinux.org/minimal-money.git
cd minimal-money
makepkg -si
```

## Container

### Docker
```bash
# Run latest version
docker run -it --rm ghcr.io/7118-eth/minimal-money:latest

# Run with persistent data
docker run -it --rm \
  -v minimal-money-data:/home/minimal/data \
  ghcr.io/7118-eth/minimal-money:latest
```

### Docker Compose
```bash
# Clone repository
git clone https://github.com/7118-eth/minimal-money.git
cd minimal-money

# Run with docker-compose
docker-compose up -d
docker-compose exec minimal-money minimal-money
```

## Build from Source

### Prerequisites
- Go 1.21 or higher
- Git
- GCC (for SQLite support)

### Build Steps
```bash
# Clone repository
git clone https://github.com/7118-eth/minimal-money.git
cd minimal-money

# Build
make build

# Install
sudo make install

# Or run directly
./minimal-money
```

## Verify Installation

After installation, verify everything is working:

```bash
# Check version
minimal-money --version

# Start the application
minimal-money
```

## Platform-Specific Notes

### macOS
- If you see "cannot be opened because the developer cannot be verified", run:
  ```bash
  xattr -d com.apple.quarantine /usr/local/bin/minimal-money
  ```

### Windows
- Requires Windows Terminal or similar for best experience
- SQLite support may be limited (CGO disabled)

### Linux
- Requires terminal with 256-color support
- May need to install ca-certificates for API access

## Troubleshooting

### "Command not found"
Make sure the binary is in your PATH:
```bash
echo $PATH
which minimal-money
```

### Database Issues
The application creates a SQLite database at:
- Linux/macOS: `~/.local/share/minimal-money/budget.db`
- Windows: `%APPDATA%\minimal-money\budget.db`

### API Connection Issues
Ensure you have internet connectivity and ca-certificates installed:
```bash
# Debian/Ubuntu
sudo apt-get install ca-certificates

# Fedora/RHEL
sudo dnf install ca-certificates

# Arch
sudo pacman -S ca-certificates
```

## Uninstallation

### Manual Installation
```bash
sudo rm /usr/local/bin/minimal-money
rm -rf ~/.local/share/minimal-money
```

### Package Managers
```bash
# Homebrew
brew uninstall minimal-money

# Snap
sudo snap remove minimal-money

# AUR
yay -R minimal-money
```

### Docker
```bash
# Remove container
docker rm -f minimal-money

# Remove image
docker rmi ghcr.io/7118-eth/minimal-money:latest

# Remove volume
docker volume rm minimal-money-data
```