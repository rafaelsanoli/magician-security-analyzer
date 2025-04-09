package scan

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type SecretFinding struct {
	File     string
	Line     int
	Severity string
	Message  string
}

func AnalyzeSecrets(path string) []SecretFinding {
	cmd := exec.Command("gitleaks", "detect", "--source="+path, "--report-format=json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Erro ao executar gitleaks: %v\n", err)
		return nil
	}

	var parsed struct {
		Findings []struct {
			Description string `json:"description"`
			File        string `json:"file"`
			StartLine   int    `json:"startLine"`
			Severity    string `json:"severity"`
		} `json:"findings"`
	}
	json.Unmarshal(output, &parsed)

	findings := []SecretFinding{}
	for _, f := range parsed.Findings {
		findings = append(findings, SecretFinding{
			File:     f.File,
			Line:     f.StartLine,
			Severity: f.Severity,
			Message:  f.Description,
		})
	}
	return findings
}
