package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/commands"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	osExit = os.Exit
)

func run(ctx context.Context, version string) error {
	f := cmdutil.NewFactory()
	cmd := commands.NewCmdRoot(f, version)
	return cmd.ExecuteContext(ctx)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		osExit(1)
	}
}
