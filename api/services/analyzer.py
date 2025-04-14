import os
from llama_cpp import Llama
import openai

USE_OPENAI = os.getenv("USE_OPENAI", "false").lower() == "true"
OPENAI_KEY = os.getenv("OPENAI_API_KEY", "")
MODEL_PATH = os.getenv("LLAMA_MODEL_PATH", "models/llama-2-7b.Q4_K_M.gguf")
llm_local = None

def init_local_model():
    global llm_local
    print(f"[INFO] Carregando modelo local: {MODEL_PATH}")
    try:
        llm_local = Llama(
            model_path=MODEL_PATH,
                 n_ctx=2048,           # Reduz o tamanho do contexto (mais rápido)
                n_threads=4,          # Usa menos threads se necessário
                n_gpu_layers=8,
        )
    except Exception as e:
        print(f"[ERRO] Falha ao carregar modelo local: {e}")
        llm_local = None

# Inicializa o modelo local apenas se for o modo local
if not USE_OPENAI:
    init_local_model()

# OpenAI setup
if USE_OPENAI and OPENAI_KEY:
    openai.api_key = OPENAI_KEY

def generate_prompt(code: str, language: str) -> str:
    return f"""
Você é um especialista em segurança. Analise o código abaixo e aponte:

1. Vulnerabilidades de segurança (injeção de comandos, autenticação, autorização, dados sensíveis).
2. Más práticas de programação.
3. Sugestões de melhorias e correções.

Se não encontrar nada, responda com "Nenhum problema encontrado".
\nCódigo:
```{language}
{code}
```"""

def analyze_code_with_llm(code_snippet: str, language: str) -> dict:
    prompt = generate_prompt(code_snippet, language)
    
    if USE_OPENAI:
        try:
            response = openai.ChatCompletion.create(
                model="gpt-4",
                messages=[{"role": "user", "content": prompt}],
                temperature=0.3
            )
            return {"analysis": response['choices'][0]['message']['content']}
        except Exception as e:
            return {"error": f"Erro na API OpenAI: {str(e)}"}
    else:
        if llm_local is None:
            return {"error": "Modelo local não foi carregado."}
        try:
            output = llm_local(prompt, max_tokens=2048, temperature=0.3)
            return {"analysis": output["choices"][0]["text"]}
        except Exception as e:
            return {"error": f"Erro no modelo local: {str(e)}"}
