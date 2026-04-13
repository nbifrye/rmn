---
title: MCPサーバー
description: rmnの内蔵MCPサーバーを使って、Claude CodeなどのAIエージェントにRedmine操作を公開。
---

# MCPサーバー

rmnには[Model Context Protocol（MCP）](https://modelcontextprotocol.io/)サーバーが内蔵されており、AIエージェントにRedmine操作を公開します。これにより、Claude CodeなどのAIアシスタントが自然言語でRedmineのチケット、プロジェクト、ユーザー、バージョン、作業時間、メンバーシップ、Wikiページを管理できるようになります。

## MCPサーバーの起動

```bash
rmn mcp serve
```

stdio方式のMCPサーバーが起動します。

## 利用可能なMCPツール

| ツール                         | 説明                                     | 読み取り専用 | 破壊的 |
|--------------------------------|------------------------------------------|-------------|--------|
| `list_issues`                  | Redmineチケットの一覧表示・フィルタ      | はい        | いいえ |
| `get_issue`                    | チケットの詳細を取得                     | はい        | いいえ |
| `create_issue`                 | 新しいチケットを作成                     | いいえ      | いいえ |
| `update_issue`                 | 既存のチケットを更新                     | いいえ      | いいえ |
| `delete_issue`                 | チケットを完全に削除                     | いいえ      | はい   |
| `list_projects`                | プロジェクトの一覧表示・フィルタ         | はい        | いいえ |
| `get_project`                  | プロジェクトの詳細を取得                 | はい        | いいえ |
| `create_project`               | 新しいプロジェクトを作成                 | いいえ      | いいえ |
| `update_project`               | 既存のプロジェクトを更新                 | いいえ      | いいえ |
| `archive_project`              | プロジェクトをアーカイブ（復元可能）     | いいえ      | いいえ |
| `unarchive_project`            | プロジェクトを復元                       | いいえ      | いいえ |
| `delete_project`               | プロジェクトを完全に削除                 | いいえ      | はい   |
| `list_users`                   | ユーザーの一覧表示・フィルタ             | はい        | いいえ |
| `get_user`                     | ユーザーの詳細を取得                     | はい        | いいえ |
| `get_current_user`             | 現在のAPIキーに対応するユーザーを取得    | はい        | いいえ |
| `list_versions`                | プロジェクトのバージョン一覧             | はい        | いいえ |
| `get_version`                  | バージョンの詳細を取得                   | はい        | いいえ |
| `create_version`               | 新しいバージョンを作成                   | いいえ      | いいえ |
| `update_version`               | 既存のバージョンを更新                   | いいえ      | いいえ |
| `delete_version`               | バージョンを完全に削除                   | いいえ      | はい   |
| `list_time_entries`            | 作業時間の一覧表示・フィルタ             | はい        | いいえ |
| `get_time_entry`               | 作業時間の詳細を取得                     | はい        | いいえ |
| `create_time_entry`            | チケット/プロジェクトに作業時間を記録    | いいえ      | いいえ |
| `update_time_entry`            | 既存の作業時間を更新                     | いいえ      | いいえ |
| `delete_time_entry`            | 作業時間を完全に削除                     | いいえ      | はい   |
| `list_memberships`             | プロジェクトメンバーシップを一覧表示     | はい        | いいえ |
| `get_membership`               | メンバーシップの詳細を取得               | はい        | いいえ |
| `create_membership`            | プロジェクトにユーザーを追加             | いいえ      | いいえ |
| `update_membership`            | メンバーシップのロールを更新             | いいえ      | いいえ |
| `delete_membership`            | メンバーシップを削除                     | いいえ      | はい   |
| `list_wiki_pages`              | プロジェクトのWikiページ一覧             | はい        | いいえ |
| `get_wiki_page`                | Wikiページの内容を取得                   | はい        | いいえ |
| `create_or_update_wiki_page`   | Wikiページを作成または更新               | いいえ      | いいえ |
| `delete_wiki_page`             | Wikiページを完全に削除                   | いいえ      | はい   |
| `list_trackers`                | トラッカー（チケットタイプ）を一覧表示   | はい        | いいえ |
| `list_issue_statuses`          | チケットステータスを一覧表示             | はい        | いいえ |

各ツールにはMCPアノテーション（`readOnlyHint`、`destructiveHint`、`idempotentHint`、`openWorldHint`）が含まれており、AIエージェントが各操作の影響を理解するのに役立ちます。

## Claude Code連携

MCP設定ファイル（例: `~/.claude/claude_desktop_config.json` またはプロジェクトの `.mcp.json`）に以下を追加します：

```json
{
  "mcpServers": {
    "rmn-redmine": {
      "command": "rmn",
      "args": ["mcp", "serve"]
    }
  }
}
```

設定が完了すると、AIエージェントが会話形式のコマンドでRedmineチケットの一覧表示、作成、更新、クローズを行えるようになります。
