---
title: Development
description: Build, test, and contribute to rmn.
---

# Development

```bash
make build    # Build binary
make test     # Run all tests
make vet      # Static analysis
make lint     # Run linters (requires golangci-lint)
make cover    # Coverage report (enforces 100% coverage)
make install  # Install to $GOPATH/bin
make clean    # Remove build artifacts
```

All pull requests must pass CI, which enforces:
- `go vet ./...` with zero warnings
- 100% test coverage

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. **Update documentation** if your change affects user-facing behavior (see `CLAUDE.md` for the code → docs mapping)
5. Ensure tests pass: `make test && make vet`
6. Commit and push
7. Open a pull request

Please follow the existing code style: one file per subcommand, table-driven tests with `httptest.NewServer`, and `fmt.Errorf("context: %w", err)` for error wrapping.

## Documentation

When your changes affect user-facing behavior (new commands, changed flags, new MCP tools), update:
1. `README.md` — the relevant section
2. `docs/` — the corresponding VitePress page
3. `docs/public/llms.txt` and `llms-full.txt` — if CLI commands or MCP tools changed

See `CLAUDE.md` for the full code → documentation mapping.

CI will post a reminder comment on PRs that modify user-facing code without documentation changes.
