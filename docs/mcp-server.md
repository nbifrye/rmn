---
title: MCP Server
description: Use rmn's built-in MCP server to expose Redmine operations to AI agents like Claude Code.
---

# MCP Server

rmn includes a built-in [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that exposes Redmine operations to AI agents. This allows AI assistants like Claude Code to manage Redmine issues, projects, users, versions, time entries, memberships, and wiki pages through natural language.

## Starting the MCP server

```bash
rmn mcp serve
```

This starts a stdio-based MCP server.

## Available MCP tools

| Tool                         | Description                          | Read-only | Destructive |
|------------------------------|--------------------------------------|-----------|-------------|
| `list_issues`                | List and filter Redmine issues       | Yes       | No          |
| `get_issue`                  | Get full details of an issue         | Yes       | No          |
| `create_issue`               | Create a new issue                   | No        | No          |
| `update_issue`               | Update an existing issue             | No        | No          |
| `delete_issue`               | Permanently delete an issue          | No        | Yes         |
| `list_projects`              | List and filter projects             | Yes       | No          |
| `get_project`                | Get full details of a project        | Yes       | No          |
| `create_project`             | Create a new project                 | No        | No          |
| `update_project`             | Update an existing project           | No        | No          |
| `archive_project`            | Archive a project (reversible)       | No        | No          |
| `unarchive_project`          | Unarchive a project                  | No        | No          |
| `delete_project`             | Permanently delete a project         | No        | Yes         |
| `list_users`                 | List and filter users                | Yes       | No          |
| `get_user`                   | Get full details of a user           | Yes       | No          |
| `get_current_user`           | Get the user for the current API key | Yes      | No          |
| `list_versions`              | List versions of a project           | Yes       | No          |
| `get_version`                | Get full details of a version        | Yes       | No          |
| `create_version`             | Create a new version                 | No        | No          |
| `update_version`             | Update an existing version           | No        | No          |
| `delete_version`             | Permanently delete a version         | No        | Yes         |
| `list_time_entries`          | List and filter time entries         | Yes       | No          |
| `get_time_entry`             | Get full details of a time entry     | Yes       | No          |
| `create_time_entry`          | Log time on an issue or project      | No        | No          |
| `update_time_entry`          | Update an existing time entry        | No        | No          |
| `delete_time_entry`          | Permanently delete a time entry      | No        | Yes         |
| `list_memberships`           | List project memberships             | Yes       | No          |
| `get_membership`             | Get full details of a membership     | Yes       | No          |
| `create_membership`          | Add a user to a project              | No        | No          |
| `update_membership`          | Update membership roles              | No        | No          |
| `delete_membership`          | Remove a membership                  | No        | Yes         |
| `list_wiki_pages`            | List wiki pages in a project         | Yes       | No          |
| `get_wiki_page`              | Get the content of a wiki page       | Yes       | No          |
| `create_or_update_wiki_page` | Create or update a wiki page         | No        | No          |
| `delete_wiki_page`           | Permanently delete a wiki page       | No        | Yes         |
| `list_trackers`              | List trackers (issue types)          | Yes       | No          |
| `list_issue_statuses`        | List issue statuses                  | Yes       | No          |

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
