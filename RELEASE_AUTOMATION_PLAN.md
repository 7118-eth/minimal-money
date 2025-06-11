# Release Automation Enhancement Plan

## Overview
This document outlines a comprehensive plan to enhance the release automation for minimal-money, progressing from quick wins to advanced release engineering practices.

## Current State
- Basic GitHub Actions release workflow
- Triggers on version tags (v*.*.*)
- Builds for Linux/macOS (amd64/arm64)
- Creates GitHub releases with static release notes
- Generates checksums

## Phase 1: Quick Wins (1-2 days)

### 1.1 Add Windows Support
**File**: `.github/workflows/release.yml`
```yaml
# Add to matrix:
- os: windows-latest
  goos: windows
  goarch: amd64
  binary_name: minimal-money.exe
- os: windows-latest
  goos: windows
  goarch: arm64
  binary_name: minimal-money.exe
```

### 1.2 Fix Deprecated Commands
**Issue**: Using deprecated `set-output` command
**Fix**: Replace with environment files
```yaml
# Old:
echo "::set-output name=VERSION::${GITHUB_REF#refs/tags/}"
# New:
echo "VERSION=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT
```

### 1.3 Add Version Info to Binary
**File**: `cmd/budget/main.go`
```go
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func init() {
    if version != "dev" {
        fmt.Printf("Minimal Money %s (%s) built on %s\n", version, commit, date)
    }
}
```

**Build command update**:
```bash
go build -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
```

### 1.4 Generate Dynamic Changelog
**Tool**: git-cliff or custom script
```yaml
- name: Generate Changelog
  run: |
    echo "## What's Changed" > changelog.md
    git log $(git describe --tags --abbrev=0 HEAD^)..HEAD \
      --pretty=format:"* %s by @%an" >> changelog.md
```

## Phase 2: GoReleaser Migration (2-3 days)

### 2.1 Create GoReleaser Configuration
**File**: `.goreleaser.yml`
```yaml
version: 2

before:
  hooks:
    - go mod tidy
    - go test ./...

builds:
  - id: minimal-money
    main: ./cmd/budget
    binary: minimal-money
    env:
      - CGO_ENABLED=1
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
      - -X main.builtBy=goreleaser

archives:
  - id: default
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- .Arch }}

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  use: github
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
  groups:
    - title: 'Features'
      regexp: '^feat'
    - title: 'Bug Fixes'
      regexp: '^fix'
    - title: 'Performance'
      regexp: '^perf'

release:
  github:
    owner: 7118-eth
    name: minimal-money
  name_template: "{{.ProjectName}} v{{.Version}}"
  header: |
    ## Minimal Money v{{.Version}}
    
    {{.Date}}
    
    Welcome to this new release!
  footer: |
    ## Thanks!
    
    Those were the changes on {{ .Tag }}!
```

### 2.2 Update Release Workflow
**File**: `.github/workflows/release.yml`
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 2.3 Cross-Compilation Setup
**File**: `.github/workflows/release.yml`
```yaml
# Add cross-compilation tools for CGO
- name: Set up cross-compilation
  run: |
    sudo apt-get update
    sudo apt-get install -y \
      gcc-aarch64-linux-gnu \
      gcc-x86-64-linux-gnu \
      gcc-mingw-w64
```

## Phase 3: Package Managers (3-5 days)

### 3.1 Homebrew Tap
**Repository**: `7118-eth/homebrew-minimal-money`
**File**: `Formula/minimal-money.rb`
```ruby
class MinimalMoney < Formula
  desc "Beautiful terminal-based portfolio tracker"
  homepage "https://github.com/7118-eth/minimal-money"
  version "1.0.0"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/7118-eth/minimal-money/releases/download/v1.0.0/minimal-money_1.0.0_darwin_arm64.tar.gz"
      sha256 "HASH_HERE"
    else
      url "https://github.com/7118-eth/minimal-money/releases/download/v1.0.0/minimal-money_1.0.0_darwin_amd64.tar.gz"
      sha256 "HASH_HERE"
    end
  end
  
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/7118-eth/minimal-money/releases/download/v1.0.0/minimal-money_1.0.0_linux_arm64.tar.gz"
      sha256 "HASH_HERE"
    else
      url "https://github.com/7118-eth/minimal-money/releases/download/v1.0.0/minimal-money_1.0.0_linux_amd64.tar.gz"
      sha256 "HASH_HERE"
    end
  end
  
  def install
    bin.install "minimal-money"
  end
  
  test do
    system "#{bin}/minimal-money", "--version"
  end
end
```

**GoReleaser Config Addition**:
```yaml
brews:
  - tap:
      owner: 7118-eth
      name: homebrew-minimal-money
    folder: Formula
    homepage: https://github.com/7118-eth/minimal-money
    description: Beautiful terminal-based portfolio tracker
    test: |
      system "#{bin}/minimal-money", "--version"
```

### 3.2 Snap Package
**File**: `snap/snapcraft.yaml`
```yaml
name: minimal-money
version: git
summary: Beautiful terminal-based portfolio tracker
description: |
  A beautiful terminal-based portfolio tracker that respects your time and privacy

base: core22
confinement: strict

apps:
  minimal-money:
    command: bin/minimal-money
    plugs:
      - home
      - network

parts:
  minimal-money:
    plugin: go
    source: .
    build-packages:
      - gcc
    stage-packages:
      - ca-certificates
```

### 3.3 AUR Package
**Repository**: `minimal-money-aur`
**File**: `PKGBUILD`
```bash
pkgname=minimal-money
pkgver=1.0.0
pkgrel=1
pkgdesc="Beautiful terminal-based portfolio tracker"
arch=('x86_64' 'aarch64')
url="https://github.com/7118-eth/minimal-money"
license=('MIT')
depends=('glibc')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  go build -o $pkgname -ldflags="-s -w -X main.version=$pkgver" cmd/budget/main.go
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 $pkgname "$pkgdir/usr/bin/$pkgname"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
```

## Phase 4: Container Support (1-2 days)

### 4.1 Dockerfile
**File**: `Dockerfile`
```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o minimal-money cmd/budget/main.go

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

COPY --from=builder /app/minimal-money /usr/local/bin/

# Create non-root user
RUN adduser -D -h /home/minimal minimal
USER minimal
WORKDIR /home/minimal

# Volume for persistent data
VOLUME ["/home/minimal/data"]

ENTRYPOINT ["minimal-money"]
```

### 4.2 GitHub Container Registry
**GoReleaser Addition**:
```yaml
dockers:
  - image_templates:
      - "ghcr.io/7118-eth/minimal-money:{{ .Tag }}"
      - "ghcr.io/7118-eth/minimal-money:latest"
    dockerfile: Dockerfile
    build_flag_templates:
      - "--pull"
      - "--label=org.opencontainers.image.created={{.Date}}"
      - "--label=org.opencontainers.image.revision={{.FullCommit}}"
      - "--label=org.opencontainers.image.version={{.Version}}"
```

## Phase 5: Advanced Features (2-3 days)

### 5.1 Release Signing
**GoReleaser Config**:
```yaml
signs:
  - artifacts: checksum
    cmd: cosign
    args:
      - sign-blob
      - "--output-certificate=${certificate}"
      - "--output-signature=${signature}"
      - "${artifact}"
      - "--yes"
```

### 5.2 Release Drafts
**File**: `.github/release-drafter.yml`
```yaml
name-template: 'v$RESOLVED_VERSION 🌟'
tag-template: 'v$RESOLVED_VERSION'
categories:
  - title: '🚀 Features'
    labels:
      - 'feature'
      - 'enhancement'
  - title: '🐛 Bug Fixes'
    labels:
      - 'fix'
      - 'bugfix'
      - 'bug'
  - title: '🧰 Maintenance'
    labels:
      - 'chore'
      - 'maintenance'
change-template: '- $TITLE @$AUTHOR (#$NUMBER)'
version-resolver:
  major:
    labels:
      - 'major'
  minor:
    labels:
      - 'minor'
  patch:
    labels:
      - 'patch'
  default: patch
template: |
  ## Changes

  $CHANGES
```

### 5.3 Semantic Release
**File**: `.github/workflows/semantic-release.yml`
```yaml
name: Semantic Release

on:
  push:
    branches:
      - main

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Semantic Release
        uses: cycjimmy/semantic-release-action@v4
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          semantic_version: 19
          extra_plugins: |
            @semantic-release/commit-analyzer
            @semantic-release/release-notes-generator
            @semantic-release/changelog
            @semantic-release/github
            @semantic-release/git
```

**File**: `.releaserc.json`
```json
{
  "branches": ["main"],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    ["@semantic-release/changelog", {
      "changelogFile": "CHANGELOG.md"
    }],
    ["@semantic-release/git", {
      "assets": ["CHANGELOG.md"],
      "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
    }],
    "@semantic-release/github"
  ]
}
```

## Implementation Order

1. **Week 1**: Phase 1 (Quick Wins)
   - Add Windows support
   - Fix deprecated commands
   - Add version info
   - Dynamic changelog

2. **Week 2**: Phase 2 (GoReleaser)
   - Create configuration
   - Update workflow
   - Test cross-compilation

3. **Week 3**: Phase 3 (Package Managers)
   - Set up Homebrew tap
   - Create Snap package
   - Submit to AUR

4. **Week 4**: Phase 4 & 5 (Containers & Advanced)
   - Docker support
   - Release signing
   - Semantic versioning

## Success Metrics

- [ ] All platforms have working binaries
- [ ] Release process takes < 10 minutes
- [ ] Users can install via package managers
- [ ] Changelog is automatically generated
- [ ] Releases are signed and verifiable
- [ ] Container images are available

## Rollback Plan

Each phase is independent. If issues arise:
1. Keep current working release.yml as backup
2. Test changes in a feature branch first
3. Use act for local testing before merge
4. Can revert to any previous phase

## Notes

- GoReleaser Pro features (like CGO cross-compilation) may require workarounds
- Homebrew tap requires separate repository
- AUR submission needs maintainer account
- Container registry requires authentication setup
- Signing requires GPG or cosign setup