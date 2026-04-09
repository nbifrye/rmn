---
title: MCPサーバー
description: rmnの内蔵MCPサーバーを使って、Claude CodeなどのAIエージェントにRedmine操作を公開。
---

# MCPサーバー

rmnには[Model Context Protocol（MCP）](https://modelcontextprotocol.io/)サーバーが内蔵されており、AIエージェントにRedmine操作を公開します。これにより、Claude CodeなどのAIアシスタントが自然言語でRedmineチケットを管理できるようになります。

## MCPサーバーの起動

```bash
rmn mcp serve
```

stdio方式のMCPサーバーが起動します。

## 利用可能なMCPツール

| ツール           | 説明                                 | 読み取り専用 | 破壊的 |
|------------------|--------------------------------------|-------------|--------|
| `list_issues`    | Redmineチケットの一覧表示・フィルタ  | はい        | いいえ |
| `get_issue`      | チケットの詳細を取得                  | はい        | いいえ |
| `create_issue`   | 新しいチケットを作成                  | いいえ      | いいえ |
| `update_issue`   | 既存のチケットを更新                  | いいえ      | いいえ |
| `delete_issue`   | チケットを完全に削除                  | いいえ      | はい   |

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
