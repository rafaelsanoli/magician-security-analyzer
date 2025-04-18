package cmd

import (
	"encoding/json"
	"fmt"
	"magician-analyzer/internal/fix"
	"magician-analyzer/internal/git"
	"magician-analyzer/internal/scan"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
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
	var createPR bool

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Executa análise de segurança com Semgrep, GoSec, Docker, Secrets e CI/CD",
		Run: func(cmd *cobra.Command, args []string) {
			results := []Finding{}

			fmt.Println("🔍 Executando Semgrep...")
			semgrepOutput := runTool("semgrep", []string{"--json", "--config=auto", "."})
			results = append(results, parseSemgrep(semgrepOutput)...)

			fmt.Println("🔍 Executando GoSec...")
			gosecOutput := runTool("gosec", []string{"-fmt=json", "./..."})
			results = append(results, parseGosec(gosecOutput)...)

			if autoFix {
				fmt.Println("🛠️  Aplicando correções automáticas:")

				if _, err := os.Stat("Dockerfile"); err == nil {
					result := fix.FixDockerfile("Dockerfile")
					if result.Modified {
						fmt.Printf("✔ Dockerfile corrigido: %s\n", result.File)
						for _, change := range result.Changes {
							fmt.Println("  - " + change)
						}
					}
				}

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

			if _, err := os.Stat("Dockerfile"); err == nil {
				for _, f := range scan.AnalyzeDockerfile("Dockerfile") {
					results = append(results, Finding{
						Tool:     "DockerScan",
						File:     f.File,
						Line:     f.Line,
						Severity: f.Severity,
						Message:  f.Message,
					})
				}
			}

			for _, f := range scan.AnalyzeSecrets(".") {
				results = append(results, Finding{
					Tool:     "Gitleaks",
					File:     f.File,
					Line:     f.Line,
					Severity: f.Severity,
					Message:  f.Message,
				})
			}

			for _, f := range scan.AnalyzeCIPipelines(".") {
				results = append(results, Finding{
					Tool:     "CIAnalyzer",
					File:     f.File,
					Line:     f.Line,
					Severity: f.Severity,
					Message:  f.Message,
				})
			}

			f, err := os.Create("results.json")
			if err != nil {
				fmt.Printf("❌ Erro ao criar results.json: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			if err := json.NewEncoder(f).Encode(results); err != nil {
				fmt.Printf("❌ Erro ao escrever JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("📄 Análise concluída. Resultados em results.json")

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
		},
	}

	cmd.Flags().BoolVar(&autoFix, "fix", false, "Aplica correções automáticas quando possível")
	cmd.Flags().BoolVar(&createPR, "pr", false, "Cria um pull request com as correções (requer --fix)")
	return cmd
}

func runTool(name string, args []string) []byte {
	cmd := exec.Command(name, args...)
	cmd.Dir, _ = filepath.Abs(".")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Erro ao executar %s: %v\n", name, err)
		os.Exit(1)
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

	if err := json.Unmarshal(output, &parsed); err != nil {
		fmt.Printf("❌ Erro ao parsear JSON do Semgrep: %v\n", err)
		return nil
	}

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

	if err := json.Unmarshal(output, &parsed); err != nil {
		fmt.Printf("❌ Erro ao parsear JSON do GoSec: %v\n", err)
		return nil
	}

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
