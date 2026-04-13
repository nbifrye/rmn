# rmn

[![CI](https://github.com/nbifrye/rmn/actions/workflows/ci.yml/badge.svg)](https://github.com/nbifrye/rmn/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nbifrye/rmn)](https://github.com/nbifrye/rmn/blob/main/go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/nbifrye/rmn/blob/main/LICENSE)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://nbifrye.github.io/rmn/ja/)

**[English README](README.md)**

**rmn** は [Redmine](https://www.redmine.org/) 用の非公式コマンドラインクライアントで、Go言語で書かれています。ターミナルからRedmineのチケット、プロジェクト、ユーザー、バージョン、作業時間、メンバーシップ、Wikiページなどを直接管理するための高速で直感的なインターフェースを提供します。[GitLab CLI（glab）](https://gitlab.com/gitlab-org/cli)にインスパイアされ、rmnはRedmineエコシステムにおなじみのコマンドパターンを導入します。

rmnには[Model Context Protocol（MCP）](https://modelcontextprotocol.io/)サーバーも内蔵されており、Claude CodeなどのAIエージェントが自然言語でRedmineインスタンスとやり取りできるようにします。キーボードでもAIアシスタントでも、rmnはRedmine REST APIを通じてRedmineチケット管理を完全に制御します。

> **注意:** このプロジェクトはRedmineプロジェクトとは関係がなく、承認もされていません。独立したコミュニティ主導のツールです。

> **警告:** このプロジェクトは実験的な段階です。APIおよびCLIインターフェースはまだ安定しておらず、予告なく破壊的変更が加わる可能性があります。

## 機能

- **幅広いRedmine APIカバレッジ** -- チケット、プロジェクト、ユーザー、バージョン、作業時間、メンバーシップ、Wikiページ、トラッカー、チケットのステータスに対応
- **チケットのライフサイクル管理** -- Redmineチケットをコマンドラインから一覧表示、閲覧、作成、更新、クローズ、削除
- **AIエージェント向けMCPサーバー** -- Model Context Protocolを通じてClaude CodeなどのAIアシスタントにRedmine操作を公開
- **複数の出力形式** -- 人間が読みやすいテーブル形式（デフォルト）と、スクリプトや自動化に便利な機械可読JSON形式
- **柔軟なチケットフィルタリング** -- プロジェクト、ステータス、担当者、トラッカーでフィルタリング。ソートとページネーションにも対応
- **GitLab CLI風エイリアス** -- おなじみの短縮コマンド（`ls`、`show`、`get`、`new`、`rm`）でワークフローを高速化
- **6つのインストール方法** -- Homebrew、mise、Nix、Go install、ビルド済みバイナリ、ソースからビルド
- **シェル補完** -- Bash、Zsh、Fish、PowerShellの自動補完に対応
- **XDG準拠の設定** -- `$XDG_CONFIG_HOME` に準拠した設定ファイルの配置
- **セキュリティ強化** -- TLS 1.2+の強制、安全なファイルパーミッション（0600）、リダイレクト時のAPIキー保護
- **テストカバレッジ100%** -- すべてのプルリクエストでCIにより強制
- **クロスプラットフォーム** -- Linux、macOS、Windowsのamd64およびarm64向けビルド済みバイナリ

## クイックスタート

```bash
# 1. インストール
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn

# 2. 認証
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# 3. チケットを一覧表示
rmn issue list -a me

# 4. チケットの詳細を表示
rmn issue view 42
```

## インストール

Linux、macOS、Windows向けのamd64およびarm64アーキテクチャ用ビルド済みバイナリが利用可能です。

### Homebrew（macOS/Linux）

```bash
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn
```

### mise

```bash
mise use -g ubi:nbifrye/rmn
```

### Nix

```bash
nix profile install github:nbifrye/rmn
```

### Go

```bash
go install github.com/nbifrye/rmn/cmd/rmn@latest
```

### バイナリダウンロード

[GitHub Releases](https://github.com/nbifrye/rmn/releases)からビルド済みバイナリをダウンロードできます。

### ソースからビルド

Go 1.24以降が必要です。

```bash
make build
```

## 設定

### 認証

```bash
# フラグを使って設定
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# または対話形式で
rmn auth login

# 設定を確認
rmn auth status
```

### 設定ファイル

設定は `~/.config/rmn/config.json`（`$XDG_CONFIG_HOME` が設定されている場合は `$XDG_CONFIG_HOME/rmn/config.json`）に保存されます：

```json
{
  "redmine_url": "https://your-redmine.example.com",
  "api_key": "your-api-key-here"
}
```

設定ファイルは `0600` パーミッションで作成されます。rmnは安全でないパーミッションの設定ファイルの読み取りを拒否します。

## ドキュメント

詳細なドキュメントは[rmnドキュメントサイト](https://nbifrye.github.io/rmn/ja/)をご覧ください。

| トピック | 説明 |
|---------|------|
| [使い方ガイド](https://nbifrye.github.io/rmn/ja/guide/usage) | チケット、プロジェクト、ユーザー、バージョン、作業時間、メンバーシップ、Wikiページなどのコマンド |
| [MCPサーバー](https://nbifrye.github.io/rmn/ja/mcp-server) | AIエージェント連携用の内蔵Model Context Protocolサーバー（35ツール） |
| [設定](https://nbifrye.github.io/rmn/ja/guide/configuration) | コマンド単位のオーバーライドを含む設定リファレンス |
| [シェル補完](https://nbifrye.github.io/rmn/ja/reference/shell-completion) | Bash、Zsh、Fish、PowerShellの自動補完設定 |
| [セキュリティ](https://nbifrye.github.io/rmn/ja/reference/security) | TLS強制、ファイルパーミッション、APIキー保護 |
| [アーキテクチャ](https://nbifrye.github.io/rmn/ja/reference/architecture) | プロジェクト構成と主要な依存関係 |

## 開発

```bash
make build    # バイナリをビルド
make test     # 全テストを実行
make vet      # 静的解析
make lint     # リンター実行（golangci-lintが必要）
make cover    # カバレッジレポート（100%カバレッジを強制）
make install  # $GOPATH/binにインストール
make clean    # ビルド成果物を削除
```

すべてのプルリクエストはCIを通過する必要があります。CIでは以下が強制されます：
- `go vet ./...` で警告ゼロ
- テストカバレッジ100%

## コントリビューション

1. リポジトリをフォーク
2. フィーチャーブランチを作成（`git checkout -b feature/my-feature`）
3. 変更を実施
4. **ドキュメントを更新** -- ユーザー向けの動作に影響する変更の場合（コードとドキュメントの対応表は `CLAUDE.md` を参照）
5. テストが通ることを確認：`make test && make vet`
6. コミットしてプッシュ
7. プルリクエストを作成

既存のコードスタイルに従ってください：サブコマンドごとに1ファイル、`httptest.NewServer` を使ったテーブル駆動テスト、エラーラッピングには `fmt.Errorf("context: %w", err)` を使用。

## ライセンス

このプロジェクトは[MITライセンス](LICENSE)の下で公開されています。
