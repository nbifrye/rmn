---
title: インストール
description: Homebrew、mise、Nix、Go、バイナリダウンロード、ソースからビルドでrmnをインストール。
---

# インストール

Linux、macOS、Windows向けのamd64およびarm64アーキテクチャ用ビルド済みバイナリが利用可能です。

## Homebrew（macOS/Linux）

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

## バイナリダウンロード

[GitHub Releases](https://github.com/nbifrye/rmn/releases)からビルド済みバイナリをダウンロードできます。

## ソースからビルド

Go 1.24以降が必要です。

```bash
make build
```
