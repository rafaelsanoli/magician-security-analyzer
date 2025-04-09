package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type CiFinding struct {
	File     string
	Line     int
	Severity string
	Message  string
}

// AnalyzeCIPipelines procura por problemas comuns em GitHub Actions e GitLab CI
func AnalyzeCIPipelines(root string) []CiFinding {
	var findings []CiFinding

	// Procurar arquivos de CI comuns
	filepaths := []string{}

	// GitHub Actions
	_ = filepath.Walk(filepath.Join(root, ".github", "workflows"), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".yml") {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	// GitLab CI
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".yml") && strings.Contains(path, ".gitlab-ci") {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	for _, path := range filepaths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "sudo") {
				findings = append(findings, CiFinding{
					File:     path,
					Line:     lineNum,
					Severity: "MEDIUM",
					Message:  "Uso de 'sudo' em pipelines pode representar riscos.",
				})
			}

			if strings.Contains(line, "curl") || strings.Contains(line, "wget") {
				if strings.Contains(line, "|") || strings.Contains(line, "sh") {
					findings = append(findings, CiFinding{
						File:     path,
						Line:     lineNum,
						Severity: "HIGH",
						Message:  "Scripts sendo executados diretamente via curl/wget.",
					})
				}
			}

			if strings.Contains(line, "uses:") && !strings.Contains(line, "@") {
				findings = append(findings, CiFinding{
					File:     path,
					Line:     lineNum,
					Severity: "MEDIUM",
					Message:  "GitHub Action sem tag ou SHA explícito — pode causar builds não determinísticos.",
				})
			}
			lineNum++
		}
	}

	return findings
}
