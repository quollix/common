package utils

import (
	"os"

	"github.com/spf13/cobra"
)

func ShowHelpCommand(cmd *cobra.Command) {
	if err := cmd.Help(); err != nil {
		Logger.Error("failed to render command help", "error", err)
		os.Exit(1)
	}
}
