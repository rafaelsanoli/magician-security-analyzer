// Estrutura inicial do projeto Magician Software Analyzer - Fase 1

// main.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"magician-analyzer/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	rootCmd.Execute()
}
