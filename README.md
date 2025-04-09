# Magician Security Analyzer — Auditoria de Segurança Automatizada para Repositórios de Código

Uma plataforma inteligente que analisa repositórios de código, Dockerfiles, pipelines CI/CD e busca por segredos sensíveis — com suporte a correções automáticas, geração de relatórios visuais e criação de pull requests.  
Ideal para DevSecOps, segurança de software, revisão de código e integrações CI.

---

## ✨ Funcionalidades

| Recurso                          | Descrição                                                                 |
|----------------------------------|---------------------------------------------------------------------------|
| ✅ **Análise de Código**         | Executa Semgrep e GoSec para detectar vulnerabilidades em código-fonte   |
| 🐳 **Scan de Dockerfile**        | Identifica práticas inseguras como `USER root`, falta de `--no-cache`    |
| 🔑 **Verificação de Segredos**   | Usa Gitleaks para encontrar tokens, senhas, e chaves hardcoded           |
| ⚙️ **Análise de CI/CD**          | Detecta usos inseguros em `.github/workflows` e `.gitlab-ci.yml`         |
| 🔧 **Correções Automáticas**     | Corrige problemas simples com `--fix`, como falta de flags ou SHA        |
| ↺ **Pull Request Automático**   | Cria branch, commit, push e PR via GitHub CLI com `--pr`                 |
| 📊 **Relatórios Bonitos**        | Gera relatório HTML/PDF com Tailwind e Reveal.js                         |
| 🌐 **API RESTful**               | Upload de projetos e análise via FastAPI + retorno de relatório          |

---

## 🚀 Como usar (CLI)

### 1. Clone o projeto

```bash
git clone https://github.com/seuusuario/magician-analyzer.git
cd magician-analyzer
```

### 2. Instale as dependências

```bash
# Go
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/returntocorp/semgrep/cmd/semgrep@latest

# Python
pip install jinja2 weasyprint fastapi uvicorn
```

> Para o modo `--pr`, você também precisa do [GitHub CLI](https://cli.github.com/):

```bash
gh auth login
```

### 3. Execute a análise no terminal

```bash
go run main.go scan
```

#### ⚙️ Flags úteis:

```bash
--fix      # Aplica correções automáticas quando possível
--pr       # Cria um pull request com as correções (requer --fix)
```

---

## 🌐 Como usar (API)

### 1. Inicie o servidor

```bash
uvicorn api.main:app --reload
```

### 2. Faça upload de um projeto `.zip` via curl ou frontend:

```bash
curl -X POST http://localhost:8000/scan \
  -F "file=@meu_projeto.zip" \
  -F "fix=true" \
  -F "pr=false"
```

### 3. Acesse o relatório

A API responde com:

```json
{
  "status": "ok",
  "findings": [...],
  "report_url": "/static/reports/report_20250409154321.html"
}
```

Abra o link gerado para ver os resultados em slides interativos Reveal.js.

---

## 📦 Estrutura do Projeto

```
magician-analyzer/
├── cmd/                 # Comando CLI (cobra/viper)
├── internal/
│   ├── scan/            # Scanners: semgrep, gosec, docker, ci, secrets
│   ├── fix/             # Corretores automáticos (--fix)
│   └── git/             # Geração de pull request (--pr)
├── scripts/             # render.py (relatório HTML/PDF)
├── templates/           # Template Reveal.js + Tailwind
├── api/                 # Servidor FastAPI
│   ├── routes/scan.py
│   └── static/reports/
├── results.json         # Saída bruta da análise
├── README.md
└── go.mod
```

---

## 📝 Exemplo de Correções com `--fix`

- Dockerfile:
    - `apk add` → `apk add --no-cache`
    - `apt-get install` → `apt-get install -y`
- GitHub Actions:
    - `uses: actions/checkout` → `uses: actions/checkout@latest`

---

## 🔮 Roadmap

| Fase           | Descrição                                                                 |
|----------------|---------------------------------------------------------------------------|
| ✅ Fase 1       | CLI + Semgrep/GoSec + relatório HTML                                      |
| ✅ Fase 2       | Dockerfile + CI/CD + Gitleaks                                             |
| ✅ Fase 3       | Correções automáticas e criação de PRs                                    |
| ✅ Fase 4       | API REST com upload, execução e retorno de relatório                     |
| 🔢 Fase 5       | Dashboard web, login, histórico de análises, multiusuário                 |

---

## ✅ Requisitos

- [Go](https://golang.org/dl/) 1.20+
- [Python 3.10+](https://www.python.org/)
- [semgrep](https://semgrep.dev/)
- [gosec](https://github.com/securego/gosec)
- [gitleaks](https://github.com/gitleaks/gitleaks)
- [gh CLI](https://cli.github.com/) (opcional)

---

## 📜 Licença

MIT License © 2025  
Feito com 💻 e 🛡️ por quem acredita em software seguro por padrão.

