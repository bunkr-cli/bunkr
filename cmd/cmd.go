package cmd

import (
	"github.com/bunkr-cli/bunkr/cmd/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bunkr",
	}
	cmd.AddCommand(
		tui.NewCommand(),
	)
	return cmd
}
