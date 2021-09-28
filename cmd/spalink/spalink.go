package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	// Subcommands
	_ "github.com/freman/spanet/subcmd/connect"
	_ "github.com/freman/spanet/subcmd/status"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
