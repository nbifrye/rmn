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
internal/commands/        Cobra command tree (root, auth, issue, mcp)
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

## MCP Tool Guidelines

- Every tool MUST have annotations (readOnlyHint, destructiveHint, idempotentHint, openWorldHint)
- Read-only tools: `WithReadOnlyHintAnnotation(true)`, `WithDestructiveHintAnnotation(false)`
- Destructive tools: `WithDestructiveHintAnnotation(true)`, `WithReadOnlyHintAnnotation(false)`
- Tool descriptions must include: what it does, what it returns, and edge case behavior
- Do NOT hardcode Redmine-instance-specific IDs in descriptions; note they are configurable
