package git

import (
	"fmt"
	"os/exec"
	"time"
)

func CreatePullRequest(branchName, title, body string) error {
	timestamp := time.Now().Format("20060102-150405")
	branch := fmt.Sprintf("%s-%s", branchName, timestamp)

	// 1. git checkout -b nova-branch
	if err := run("git", "checkout", "-b", branch); err != nil {
		return fmt.Errorf("falha ao criar branch: %v", err)
	}

	// 2. git add .
	if err := run("git", "add", "."); err != nil {
		return fmt.Errorf("falha ao adicionar arquivos: %v", err)
	}

	// 3. git commit
	commitMsg := "[magician] correções automáticas de segurança"
	if err := run("git", "commit", "-m", commitMsg); err != nil {
		return fmt.Errorf("falha ao fazer commit: %v", err)
	}

	// 4. git push
	if err := run("git", "push", "-u", "origin", branch); err != nil {
		return fmt.Errorf("falha ao fazer push: %v", err)
	}

	// 5. gh pr create
	args := []string{
		"pr", "create",
		"--title", title,
		"--body", body,
		"--head", branch,
		"--base", "main",
	}
	if err := run("gh", args...); err != nil {
		return fmt.Errorf("falha ao criar PR via GitHub CLI: %v", err)
	}

	fmt.Printf("✅ Pull Request criado com sucesso para a branch %s\n", branch)
	return nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run()
}
