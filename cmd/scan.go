package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"magician-analyzer/internal/fix"
	"magician-analyzer/internal/git"
	"magician-analyzer/internal/scan"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Finding struct {
	Tool     string `json:"tool"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func newScanCommand() *cobra.Command {
	var autoFix bool
	var autoFix bool
	var createPR bool

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Executa análise de segurança com Semgrep, GoSec, Docker, Secrets e CI/CD",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().BoolVar(&autoFix, "fix", false, "Aplica correções automáticas quando possível")
			cmd.Flags().BoolVar(&createPR, "pr", false, "Cria um pull request com as correções (requer --fix)")

			results := []Finding{}

			// Análise com Semgrep
			semgrepOutput := runTool("semgrep", []string{"--json", "--config=auto", "."})
			results = append(results, parseSemgrep(semgrepOutput)...)

			// Análise com GoSec
			gosecOutput := runTool("gosec", []string{"-fmt=json", "./..."})
			results = append(results, parseGosec(gosecOutput)...)

			// Correções automáticas se ativado
			if autoFix {
				fmt.Println("🛠️  Aplicando correções automáticas:")

				// Dockerfile
				if _, err := os.Stat("Dockerfile"); err == nil {
					result := fix.FixDockerfile("Dockerfile")
					if result.Modified {
						fmt.Printf("✔ Dockerfile corrigido: %s\n", result.File)
						for _, change := range result.Changes {
							fmt.Println("  - " + change)
						}
					}
				}

				// Workflows CI/CD
				ciPaths := scan.GetCIPipelinePaths(".")
				for _, path := range ciPaths {
					result := fix.FixCIPipeline(path)
					if result.Modified {
						fmt.Printf("✔ Workflow corrigido: %s\n", result.File)
						for _, change := range result.Changes {
							fmt.Println("  - " + change)
						}
					}
				}
			}

			// Dockerfile Scan
			for _, f := range scan.AnalyzeDockerfile("Dockerfile") {
				results = append(results, Finding{
					Tool:     "DockerScan",
					File:     f.File,
					Line:     f.Line,
					Severity: f.Severity,
					Message:  f.Message,
				})
			}

			// Secret Scan
			for _, f := range scan.AnalyzeSecrets(".") {
				results = append(results, Finding{
					Tool:     "Gitleaks",
					File:     f.File,
					Line:     f.Line,
					Severity: f.Severity,
					Message:  f.Message,
				})
			}

			// CI/CD Scan
			for _, f := range scan.AnalyzeCIPipelines(".") {
				results = append(results, Finding{
					Tool:     "CIAnalyzer",
					File:     f.File,
					Line:     f.Line,
					Severity: f.Severity,
					Message:  f.Message,
				})
			}

			// Exportar resultados
			f, _ := os.Create("results.json")
			defer f.Close()
			json.NewEncoder(f).Encode(results)
			fmt.Println("📄 Análise concluída. Resultados em results.json")
		},
	}
	if autoFix && createPR {
		err := git.CreatePullRequest(
			"magician-fix",
			"[magician] Correções automáticas de segurança",
			"Este PR foi gerado automaticamente pelo Magician Software Analyzer contendo correções seguras detectadas em análise.",
		)
		if err != nil {
			fmt.Printf("❌ Falha ao criar Pull Request: %v\n", err)
		}
	}

	cmd.Flags().BoolVar(&autoFix, "fix", false, "Aplica correções automáticas quando possível")
	return cmd
}

func runTool(name string, args []string) []byte {
	cmd := exec.Command(name, args...)
	cmd.Dir, _ = filepath.Abs(".")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Erro ao executar %s: %v\n", name, err)
		return nil
	}
	return output
}

func parseSemgrep(output []byte) []Finding {
	var parsed struct {
		Results []struct {
			CheckID string `json:"check_id"`
			Path    string `json:"path"`
			Start   struct {
				Line int `json:"line"`
			} `json:"start"`
			Extra struct {
				Message  string `json:"message"`
				Severity string `json:"severity"`
			} `json:"extra"`
		} `json:"results"`
	}
	json.Unmarshal(output, &parsed)

	findings := []Finding{}
	for _, r := range parsed.Results {
		findings = append(findings, Finding{
			Tool:     "Semgrep",
			File:     r.Path,
			Line:     r.Start.Line,
			Severity: r.Extra.Severity,
			Message:  r.Extra.Message,
		})
	}
	return findings
}

func parseGosec(output []byte) []Finding {
	var parsed struct {
		Issues []struct {
			Severity   string `json:"severity"`
			Confidence string `json:"confidence"`
			File       string `json:"file"`
			Line       string `json:"line"`
			Details    string `json:"details"`
		} `json:"Issues"`
	}
	json.Unmarshal(output, &parsed)

	findings := []Finding{}
	for _, issue := range parsed.Issues {
		line, _ := strconv.Atoi(issue.Line)
		findings = append(findings, Finding{
			Tool:     "GoSec",
			File:     issue.File,
			Line:     line,
			Severity: issue.Severity,
			Message:  issue.Details,
		})
	}
	return findings
}
