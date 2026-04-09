---
title: 開発
description: rmnのビルド、テスト、コントリビューション。
---

# 開発

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
4. テストが通ることを確認：`make test && make vet`
5. コミットしてプッシュ
6. プルリクエストを作成

既存のコードスタイルに従ってください：サブコマンドごとに1ファイル、`httptest.NewServer` を使ったテーブル駆動テスト、エラーラッピングには `fmt.Errorf("context: %w", err)` を使用。

## ドキュメント

ユーザー向けの動作に影響する変更（新しいコマンド、フラグの変更、新しいMCPツール）を行った場合は、以下を更新してください：
1. `README.md` — 該当するセクション
2. `docs/` — 対応するVitePress ページ
3. `docs/public/llms.txt` と `llms-full.txt` — CLIコマンドやMCPツールが変更された場合

詳細なコードとドキュメントの対応表は `CLAUDE.md` を参照してください。

CIは、ユーザー向けコードの変更にドキュメントの変更が伴わないプルリクエストにリマインダーコメントを投稿します。
