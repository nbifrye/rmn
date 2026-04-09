---
title: 使い方
description: rmnでRedmineチケットを管理 — 一覧表示、閲覧、作成、更新、クローズ、削除。
---

# 使い方

## チケットの一覧表示

```bash
rmn issue list                                    # オープンなチケットを一覧表示（デフォルト上限: 25）
rmn issue list -p my-project                      # プロジェクトでフィルタ
rmn issue list -s closed                          # ステータスでフィルタ（open, closed, *, またはステータスID）
rmn issue list -a me                              # 自分に割り当てられたチケット
rmn issue list -t 2                               # トラッカーIDでフィルタ
rmn issue list --sort updated_on:desc             # カラムでソート
rmn issue list -l 50 --offset 100                 # ページネーション
rmn issue list -p my-project -s closed -a me      # フィルタを組み合わせ
```

## チケットの閲覧

```bash
rmn issue view 42                                 # チケットの詳細を表示
```

## チケットの作成

```bash
rmn issue create -p my-project -s "Bug report"
rmn issue create -p my-project -s "Feature request" -d "Detailed description" \
  -t 2 --priority 3 -a 5 --start-date 2025-01-01 --due-date 2025-03-31
```

作成時の全フラグ: `--project/-p`（必須）、`--subject/-s`（必須）、`--description/-d`、`--tracker/-t`、`--priority`、`--assignee/-a`、`--category`、`--version`、`--parent`、`--start-date`、`--due-date`、`--estimated-hours`、`--done-ratio`。

## チケットの更新

```bash
rmn issue update 42 --status 3                    # ステータスを変更
rmn issue update 42 -n "Work in progress"         # ノートを追加
rmn issue update 42 --done-ratio 50 --priority 2  # 複数フィールドを更新
```

指定されたフィールドのみ変更され、省略されたフィールドは変更されません。作成時の全フラグに加えて、`--status` と `--notes/-n` が使用可能です。

## チケットのクローズ

```bash
rmn issue close 42                                # クローズ（デフォルトのステータスID 5）
rmn issue close 42 --status 6                     # カスタムステータスIDでクローズ
rmn issue close 42 -n "Fixed in v1.2"             # ノート付きでクローズ
```

## チケットの削除

```bash
rmn issue delete 42                               # 確認プロンプト付きで削除
rmn issue delete 42 -y                            # 確認をスキップ
```

## コマンドエイリアス

| コマンド             | エイリアス     |
|----------------------|----------------|
| `rmn issue list`     | `ls`           |
| `rmn issue view`     | `show`, `get`  |
| `rmn issue create`   | `new`          |
| `rmn issue delete`   | `rm`           |

```bash
rmn issue ls                    # rmn issue list と同じ
rmn issue show 42               # rmn issue view 42 と同じ
rmn issue new -p proj -s "Bug"  # rmn issue create ... と同じ
rmn issue rm 42                 # rmn issue delete 42 と同じ
```

## グローバルフラグ

| フラグ           | 説明                                     |
|------------------|------------------------------------------|
| `--output`       | 出力形式: `table`（デフォルト）または `json` |
| `--redmine-url`  | RedmineインスタンスURLをオーバーライド   |
| `--api-key`      | Redmine APIキーをオーバーライド          |

## JSON出力

`--output json` を任意のコマンドに指定すると、スクリプトやパイプに便利な機械可読出力が得られます：

```bash
rmn issue list --output json                      # チケットのJSON配列
rmn issue view 42 --output json                   # チケット全体をJSONで表示
rmn issue list -p my-project --output json | jq '.issues[].subject'
```
