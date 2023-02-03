package tui

import (
	"github.com/bunkr-cli/bunkr/cmd/tui/models"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launches a terminal UI",
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	model, err := models.NewRoot()
	if err != nil {
		return err
	}

	zone.NewGlobal()
	_, err = tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run()
	return err
}
