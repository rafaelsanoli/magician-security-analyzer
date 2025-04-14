// Estrutura inicial do projeto Magician Software Analyzer - Fase 1

// main.go
package main

import (
	"magician-analyzer/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	rootCmd.Execute()
}
