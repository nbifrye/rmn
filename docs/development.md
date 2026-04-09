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
4. Ensure tests pass: `make test && make vet`
5. Commit and push
6. Open a pull request

Please follow the existing code style: one file per subcommand, table-driven tests with `httptest.NewServer`, and `fmt.Errorf("context: %w", err)` for error wrapping.
