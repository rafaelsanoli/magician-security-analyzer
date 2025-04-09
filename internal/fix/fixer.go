package fix

import (
	"bufio"
	"os"
	"strings"
)

type FixResult struct {
	File     string
	Modified bool
	Changes  []string
}

// FixDockerfile aplica correções automáticas em instruções inseguras
func FixDockerfile(path string) FixResult {
	file, err := os.Open(path)
	if err != nil {
		return FixResult{File: path}
	}
	defer file.Close()

	var output []string
	var modified bool
	var changes []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Corrigir instalação sem flags seguras
		if strings.Contains(line, "apk add") && !strings.Contains(line, "--no-cache") {
			line += " --no-cache"
			modified = true
			changes = append(changes, "Adicionado '--no-cache' em linha com apk add")
		}
		if strings.Contains(line, "apt-get install") && !strings.Contains(line, "-y") {
			line += " -y"
			modified = true
			changes = append(changes, "Adicionado '-y' em linha com apt-get install")
		}

		// Sinalizar uso de root (mas não alterar automaticamente)
		if strings.Contains(line, "USER root") {
			changes = append(changes, "Aviso: 'USER root' detectado (não removido automaticamente)")
		}

		output = append(output, line)
	}

	if modified {
		os.WriteFile(path, []byte(strings.Join(output, "\n")+"\n"), 0644)
	}

	return FixResult{
		File:     path,
		Modified: modified,
		Changes:  changes,
	}
}

// FixCIPipeline corrige problemas comuns em arquivos de workflow
func FixCIPipeline(path string) FixResult {
	file, err := os.Open(path)
	if err != nil {
		return FixResult{File: path}
	}
	defer file.Close()

	var output []string
	var modified bool
	var changes []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Corrigir uses sem tag
		if strings.Contains(line, "uses:") && !strings.Contains(line, "@") {
			line = strings.Replace(line, "uses:", "uses:", 1) + "@latest"
			modified = true
			changes = append(changes, "Adicionado '@latest' em action sem versão fixada")
		}

		output = append(output, line)
	}

	if modified {
		os.WriteFile(path, []byte(strings.Join(output, "\n")+"\n"), 0644)
	}

	return FixResult{
		File:     path,
		Modified: modified,
		Changes:  changes,
	}
}
