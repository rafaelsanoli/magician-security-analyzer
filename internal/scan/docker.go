package scan

import (
	"bufio"
	"os"
	"strings"
)

type DockerFinding struct {
	File     string
	Line     int
	Severity string
	Message  string
}

func AnalyzeDockerfile(path string) []DockerFinding {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var findings []DockerFinding
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "USER root") {
			findings = append(findings, DockerFinding{
				File:     path,
				Line:     lineNum,
				Severity: "HIGH",
				Message:  "Uso de 'USER root' no Dockerfile pode representar risco de segurança.",
			})
		}
		if strings.Contains(line, "apk add") || strings.Contains(line, "apt-get install") {
			if !strings.Contains(line, "--no-cache") && !strings.Contains(line, "-y") {
				findings = append(findings, DockerFinding{
					File:     path,
					Line:     lineNum,
					Severity: "MEDIUM",
					Message:  "Instalação de pacotes sem flags seguras (ex: --no-cache, -y).",
				})
			}
		}
		lineNum++
	}
	return findings
}
