---
title: Installation
description: Install rmn via Homebrew, mise, Nix, Go, binary download, or build from source.
---

# Installation

Pre-built binaries are available for Linux, macOS, and Windows on both amd64 and arm64 architectures.

## Homebrew (macOS/Linux)

```bash
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn
```

## mise

```bash
mise use -g ubi:nbifrye/rmn
```

## Nix

```bash
nix profile install github:nbifrye/rmn
```

## Go

```bash
go install github.com/nbifrye/rmn/cmd/rmn@latest
```

## Binary download

Download pre-built binaries from [GitHub Releases](https://github.com/nbifrye/rmn/releases).

## Build from source

Requires Go 1.24 or later.

```bash
make build
```
