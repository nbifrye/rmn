---
title: アーキテクチャ
description: rmnプロジェクトのアーキテクチャ — モジュール構成と主要な依存関係。
---

# アーキテクチャ

```
cmd/rmn/main.go          エントリーポイント（シグナル処理、ファクトリ、ルートコマンド）
internal/api/             Redmine HTTPクライアント + ドメイン型
internal/commands/        Cobraコマンドツリー（root, auth, issue, mcp）
internal/cmdutil/         ファクトリ（依存性注入）、IOStreams
internal/config/          XDG準拠のJSON設定（~/.config/rmn/config.json）
```

rmnはCLIフレームワークに[Cobra](https://github.com/spf13/cobra)を、MCPサーバーの実装に[go-sdk](https://github.com/modelcontextprotocol/go-sdk)を使用しています。コードベースは依存性注入にファクトリパターンを採用しており、すべてのコマンドをモックHTTPサーバーでテスト可能にしています。
