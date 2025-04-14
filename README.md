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
## Configuração do Modelo Llama-2-7B Q4_K_M para Análise de Segurança de Código

Este projeto utiliza modelos de linguagem para realizar análises automatizadas de segurança de código, detectando vulnerabilidades, más práticas e sugerindo melhorias. O modelo **Llama-2-7B Q4_K_M** é usado localmente para análise, e a API da OpenAI pode ser utilizada como fallback para obter melhores resultados quando necessário.

### 1. **Pré-requisitos**

- Python 3.8 ou superior
- Biblioteca `llama_cpp` para interação com o modelo Llama
- Biblioteca `openai` para integração com a API da OpenAI
- Acesso à internet para baixar o modelo e interagir com a API da OpenAI

### 2. **Instalação**

1. Clone o repositório:

    ```bash
    git clone https://github.com/seu-usuario/magician-security-analyzer.git
    cd magician-security-analyzer
    ```

2. Crie um ambiente virtual (recomendado):

    ```bash
    python3 -m venv venv
    source venv/bin/activate  # No Linux/macOS
    .\venv\Scripts\activate   # No Windows
    ```

3. Instale as dependências:

    ```bash
    pip install -r requirements.txt
    ```

### 3. **Configuração do Modelo Local (Llama-2-7B Q4_K_M)**

1. Coloque o arquivo do modelo **Llama-2-7B Q4_K_M** na pasta `models/` do seu projeto.

2. Defina a variável de ambiente `LLAMA_MODEL_PATH` com o caminho do modelo, por exemplo:

    ```bash
    export LLAMA_MODEL_PATH="/caminho/para/o/modelo/llama-2-7b-q4_k_m.gguf"
    ```

3. Defina a variável de ambiente `USE_OPENAI` para `false` para usar o modelo local:

    ```bash
    export USE_OPENAI="false"
    ```

### 4. **Configuração da API OpenAI (Fallback)**

1. Crie uma conta no [OpenAI](https://platform.openai.com/) e gere uma chave de API.

2. Defina a variável de ambiente `OPENAI_API_KEY` com sua chave de API:

    ```bash
    export OPENAI_API_KEY="sua-chave-da-api"
    ```

3. Defina a variável `USE_OPENAI` para `true` caso deseje usar a OpenAI como fallback:

    ```bash
    export USE_OPENAI="true"
    ```

### 5. **Rodando o Analisador de Código**

Para rodar a análise de segurança de código, basta chamar o script que você deseja, passando o código a ser analisado e o tipo de linguagem (por exemplo, Python):

```bash
python analyze_code.py --code "código_fonte_aqui" --language "python"
```

### 6. Modificando o Comportamento do Analisador

O comportamento do modelo pode ser ajustado via prompts para garantir que ele identifique vulnerabilidades de segurança, más práticas de programação e forneça sugestões de melhorias. Você pode ajustar o prompt na função generate_prompt no código, para alterar a forma como o modelo analisa os códigos.
### Exemplo de Prompt Personalizado
```bash
def generate_prompt(code: str, language: str) -> str:
    return f"""
    Você é um especialista em segurança de software e deve analisar o código abaixo, procurando especificamente:
    
    1. Vulnerabilidades de segurança, como:
        - Injeção de comandos
        - Falhas de autenticação
        - Falhas de autorização
        - Exposição de dados sensíveis
    
    2. Más práticas de programação (ex: uso de funções obsoletas, falta de validação de entrada).
    
    3. Sugestões de melhorias e correções, considerando as melhores práticas de segurança.
    
    Explique cada ponto identificado de forma clara e técnica.
    
    Código:
    ```{language}
    {code}
    ```
    """
```
### 7. Desempenho e Otimização

Para garantir que a análise de segurança de código seja rápida e eficiente, considere as seguintes otimizações:

    Uso de Modelos Menores: Se o modelo local estiver demorando muito para gerar as análises, considere usar modelos menores ou ajustar o número de camadas da GPU (caso você tenha uma placa gráfica com recursos limitados).

    Uso de OpenAI para Melhor Performance: Em casos onde a análise precisa ser feita rapidamente, você pode usar a OpenAI como fallback, configurando a variável USE_OPENAI para true.

### Contribuindo

Se você deseja contribuir para este projeto, fique à vontade para enviar pull requests. Fique atento para melhorias, otimizações de desempenho e novas funcionalidades!

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

