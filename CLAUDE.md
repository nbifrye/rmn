# rmn - Redmine CLI

## Build & Test

```bash
go build -o rmn ./cmd/rmn/       # Build binary
go test ./...                     # Run all tests
go test ./internal/api/...        # Run API tests only
go test -run TestListIssues ./... # Run single test
go vet ./...                      # Static analysis
```

## Architecture

```
cmd/rmn/main.go          Entry point (signal handling, factory, root command)
internal/api/             Redmine HTTP client + domain types
internal/commands/        Cobra command tree (root, auth, issue, mcp)
internal/cmdutil/         Factory (DI), IOStreams
internal/config/          XDG-compliant JSON config (~/.config/rmn/config.json)
```

## Code Style

- Error wrapping: `fmt.Errorf("context: %w", err)` - always add context
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
