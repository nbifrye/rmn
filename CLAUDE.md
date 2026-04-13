# rmn - Redmine CLI

## Build & Test

```bash
go build -o rmn ./cmd/rmn/       # Build binary
go test ./...                     # Run all tests
go test -run TestFoo ./...        # Run single test
go vet ./...                      # Static analysis
make lint                         # Linter (golangci-lint required)
```

## Verify Changes

After any code change, always run:
1. `go vet ./...` — must pass with zero warnings
2. `go test ./...` — must pass with zero failures

## Architecture

```
cmd/rmn/main.go          Entry point (signal handling, factory, root command)
internal/api/             Redmine HTTP client + domain types
internal/commands/        Cobra command tree (root, auth, issue, project, user,
                          version, timeentry, membership, wiki, tracker, status, mcp)
internal/cmdutil/         Factory (DI), IOStreams
internal/config/          XDG-compliant JSON config (~/.config/rmn/config.json)
```

## Code Style

- Error wrapping: `fmt.Errorf("context: %w", err)` — always add context
- Cobra commands: one file per subcommand, constructor returns `*cobra.Command`
- Use `cmd.Context()` for context propagation
- Use `f.IO.Out` / `f.IO.ErrOut` for output, never bare `fmt.Println`
- Pointer types (`*int`) for optional update params to distinguish "not set" from zero
- Table-driven tests preferred, use `httptest.NewServer` for API mocking
- Test factory: use `newTestFactory(srv)` from `testutil_test.go` in command tests

## Conventions

- Command aliases follow GitLab CLI patterns (ls, show, new, rm)
- MCP tool names use snake_case matching CLI subcommands
- Config validation happens in Factory.APIClient(), not in individual commands
- Global `--output` flag: "table" (default) or "json"

## Adding a New Command

1. Create `internal/commands/<group>/<name>.go` with `NewCmd<Name>(f *cmdutil.Factory) *cobra.Command`
2. Register in the parent command group's `.go` file via `cmd.AddCommand()`
3. Support `--output json` for machine-readable output
4. Create `<name>_test.go` with httptest mock server
5. Add corresponding MCP tool in `internal/commands/mcp/serve.go` with annotations
6. Add usage examples to `README.md` under the appropriate section
7. Add usage examples to `docs/guide/usage.md` (or the relevant VitePress page)
8. Update `docs/public/llms.txt` and `docs/public/llms-full.txt` with the new command

## MCP Tool Guidelines

- Every tool MUST have annotations (readOnlyHint, destructiveHint, idempotentHint, openWorldHint)
- Read-only tools: `WithReadOnlyHintAnnotation(true)`, `WithDestructiveHintAnnotation(false)`
- Destructive tools: `WithDestructiveHintAnnotation(true)`, `WithReadOnlyHintAnnotation(false)`
- Tool descriptions must include: what it does, what it returns, and edge case behavior
- Do NOT hardcode Redmine-instance-specific IDs in descriptions; note they are configurable

## Documentation

### Code → Documentation Mapping

When modifying code in these directories, update the corresponding documentation:

| Source | Documentation Files |
|---|---|
| `internal/commands/issue/` | `docs/guide/usage.md` |
| `internal/commands/mcp/serve.go` | `docs/mcp-server.md`, `docs/public/llms.txt`, `docs/public/llms-full.txt` |
| `internal/commands/auth/` | `README.md` (Configuration section), `docs/guide/configuration.md` |
| `internal/config/` | `README.md` (Configuration section), `docs/guide/configuration.md` |
| `internal/api/` (security features) | `docs/reference/security.md` |
| `internal/commands/root.go` (global flags) | `docs/guide/usage.md` |
| `cmd/rmn/` | `docs/reference/architecture.md` |
| `internal/cmdutil/` | `docs/reference/architecture.md` |
| `Makefile` | `README.md` (Development section), `docs/development.md` |

### Japanese Translation (i18n)

The documentation supports English (default) and Japanese. When updating English docs, also update the corresponding Japanese files:

| English | Japanese |
|---|---|
| `README.md` | `README.ja.md` |
| `docs/<page>.md` | `docs/ja/<page>.md` |
| `docs/guide/<page>.md` | `docs/ja/guide/<page>.md` |
| `docs/reference/<page>.md` | `docs/ja/reference/<page>.md` |
| `docs/public/llms.txt` | `docs/public/llms-ja.txt` |
| `docs/public/llms-full.txt` | `docs/public/llms-full-ja.txt` |

VitePress i18n is configured in `docs/.vitepress/config.mts` via the `locales` object. The Japanese locale (`ja`) has its own nav, sidebar, and UI labels.

### When to Update Documentation

**Always update docs when:**
- Adding a new CLI command or subcommand
- Adding or removing command flags
- Adding or changing MCP tools
- Changing config file format or location
- Changing authentication behavior
- Changing security behavior (TLS, permissions, redirects)
- Changing installation methods or build requirements
- Changing command aliases

**No doc update needed for:**
- Internal refactors that do not change user-facing behavior
- Test-only changes (new tests, test fixes)
- CI/workflow changes (unless they affect the Development section)
- Dependency updates (go.mod/go.sum) without behavior changes
- Code style or linting fixes

### Documentation Checklist

When docs are needed, update all applicable locations:
1. `README.md` — only if the change affects sections still in the README (features, installation, configuration basics)
2. `docs/` — the corresponding VitePress page
3. `docs/public/llms.txt` and `docs/public/llms-full.txt` — if CLI commands or MCP tools changed
4. `CLAUDE.md` — if the change affects development workflow or conventions
5. `README.ja.md` and `docs/ja/` — the corresponding Japanese translations
6. `docs/public/llms-ja.txt` and `docs/public/llms-full-ja.txt` — if the English LLM files were updated
