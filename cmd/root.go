package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "magician",
		Short: "Magician Software Analyzer CLI",
	}

	cmd.AddCommand(newScanCommand())
	return cmd
}
