# rmn

A CLI tool for interacting with [Redmine](https://www.redmine.org/), inspired by [GitLab CLI (glab)](https://gitlab.com/gitlab-org/cli).

## Installation

```bash
go install github.com/nbifrye/rmn/cmd/rmn@latest
```

Or build from source:

```bash
make build
```

## Configuration

```bash
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY
```

Or interactively:

```bash
rmn auth login
```

Verify your configuration:

```bash
rmn auth status
```

Configuration is stored in `~/.config/rmn/config.json`. Global flags `--redmine-url` and `--api-key` can override the stored configuration per-command.

## Usage

### Issues

```bash
rmn issue list                          # List open issues
rmn issue list -p my-project -s closed  # List closed issues in a project
rmn issue list -a me                    # List issues assigned to you
rmn issue view 42                       # View issue details
rmn issue create -p my-project -s "Bug report" -d "Description here"
rmn issue update 42 --status 3 --notes "In progress"
rmn issue close 42                      # Close issue (status ID 5)
rmn issue delete 42                     # Delete issue (with confirmation)
```

Use `--output json` on any command for machine-readable output.

### MCP Server

Expose Redmine operations to AI agents via the [Model Context Protocol](https://modelcontextprotocol.io/):

```bash
rmn mcp serve
```

This starts a stdio-based MCP server with tools: `list_issues`, `get_issue`, `create_issue`, `update_issue`, `delete_issue`.

To use with Claude Code, add to your MCP configuration:

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

### Shell Completion

```bash
source <(rmn completion bash)           # Bash
rmn completion zsh > "${fpath[1]}/_rmn" # Zsh
rmn completion fish | source            # Fish
```

## Development

```bash
make build    # Build binary
make test     # Run tests
make vet      # Static analysis
make lint     # Run linters (requires golangci-lint)
make install  # Install to $GOPATH/bin
```

## License

See [LICENSE](LICENSE) for details.
