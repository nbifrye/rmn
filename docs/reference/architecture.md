---
title: Architecture
description: rmn project architecture — module structure and key dependencies.
---

# Architecture

```
cmd/rmn/main.go          Entry point (signal handling, factory, root command)
internal/api/             Redmine HTTP client + domain types
internal/commands/        Cobra command tree (root, auth, issue, project, user,
                          version, timeentry, membership, wiki, tracker, status, mcp)
internal/cmdutil/         Factory (dependency injection), IOStreams
internal/config/          XDG-compliant JSON config (~/.config/rmn/config.json)
```

rmn uses [Cobra](https://github.com/spf13/cobra) for the CLI framework and [go-sdk](https://github.com/modelcontextprotocol/go-sdk) for the MCP server implementation. The codebase follows a factory pattern for dependency injection, making all commands testable with mock HTTP servers.
