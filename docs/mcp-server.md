---
title: MCP Server
description: Use rmn's built-in MCP server to expose Redmine operations to AI agents like Claude Code.
---

# MCP Server

rmn includes a built-in [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that exposes Redmine operations to AI agents. This allows AI assistants like Claude Code to manage Redmine issues through natural language.

## Starting the MCP server

```bash
rmn mcp serve
```

This starts a stdio-based MCP server.

## Available MCP tools

| Tool             | Description                          | Read-only | Destructive |
|------------------|--------------------------------------|-----------|-------------|
| `list_issues`    | List and filter Redmine issues       | Yes       | No          |
| `get_issue`      | Get full details of an issue         | Yes       | No          |
| `create_issue`   | Create a new issue                   | No        | No          |
| `update_issue`   | Update an existing issue             | No        | No          |
| `delete_issue`   | Permanently delete an issue          | No        | Yes         |

Each tool includes MCP annotations (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) to help AI agents understand the impact of each operation.

## Claude Code integration

Add the following to your MCP configuration (e.g. `~/.claude/claude_desktop_config.json` or project `.mcp.json`):

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

Once configured, your AI agent can list, create, update, and close Redmine issues through conversational commands.
