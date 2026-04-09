---
layout: home
title: rmn — Redmine CLIツール
titleTemplate: Redmine用コマンドラインクライアント

hero:
  name: rmn
  text: Redmine CLIツール
  tagline: Go言語で書かれた高速なオープンソースのRedmineコマンドラインクライアント。ターミナルからチケットを管理。内蔵MCPサーバーでAIエージェント連携も可能。
  actions:
    - theme: brand
      text: はじめる
      link: /ja/guide/installation
    - theme: alt
      text: GitHubで見る
      link: https://github.com/nbifrye/rmn

features:
  - title: チケットのライフサイクル管理
    details: Redmineのチケットをコマンドラインから一覧表示、閲覧、作成、更新、クローズ、削除。
  - title: AIエージェント向けMCPサーバー
    details: Model Context Protocolを通じて、Claude CodeなどのAIアシスタントにRedmine操作を公開。
  - title: 複数の出力形式
    details: 人間が読みやすいテーブル形式（デフォルト）と、スクリプトや自動化に便利なJSON形式に対応。
  - title: 柔軟なフィルタリング
    details: プロジェクト、ステータス、担当者、トラッカーでフィルタリング。ソートやページネーションにも対応。
  - title: GitLab CLI風エイリアス
    details: おなじみの短縮コマンド（ls, show, get, new, rm）でワークフローを高速化。
  - title: 6つのインストール方法
    details: Homebrew、mise、Nix、Go install、ビルド済みバイナリ、ソースからビルド。
  - title: シェル補完
    details: Bash、Zsh、Fish、PowerShellの自動補完に対応。
  - title: XDG準拠の設定
    details: $XDG_CONFIG_HOMEに準拠した設定ファイルの配置。
  - title: セキュリティ強化
    details: TLS 1.2+の強制、安全なファイルパーミッション（0600）、リダイレクト時のAPIキー保護。
  - title: テストカバレッジ100%
    details: すべてのプルリクエストでCIにより強制。
  - title: クロスプラットフォーム
    details: Linux、macOS、Windowsのamd64およびarm64向けビルド済みバイナリを提供。
---

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
