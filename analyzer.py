import os
from llama_cpp import Llama
import openai
import json

USE_OPENAI = os.getenv("USE_OPENAI", "false").lower() == "true"
OPENAI_KEY = os.getenv("OPENAI_API_KEY", "")
MODEL_PATH = os.getenv("LLAMA_MODEL_PATH", "models/llama-2-7b-32k-instruct.Q4_K_M.gguf")
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
                temperature=0.8,
                max_tokens=2048,
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

def generate_json_prompt(code: str, language: str) -> str:
    return f"""
Você é um especialista em segurança de software. Analise o código a seguir e retorne as vulnerabilidades detectadas no formato JSON.

Responda SOMENTE com o JSON.

Formato esperado:
[
  {{
    "tipo_de_falha": "descrição da falha",
    "trecho": "trecho vulnerável",
    "recomendacao": "como corrigir ou melhorar",
    "severidade": "baixa | média | alta | crítica"
  }}
]

Código:
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
        


def analyze_code_with_llm_json(code_snippet: str, language: str) -> dict:
    prompt = generate_json_prompt(code_snippet, language)

    try:
        if USE_OPENAI:
            response = openai.ChatCompletion.create(
                model="gpt-4",
                messages=[{"role": "user", "content": prompt}],
                temperature=0.3
            )
            content = response['choices'][0]['message']['content']
        else:
            if llm_local is None:
                return {"error": "Modelo local não foi carregado."}
            output = llm_local(prompt, max_tokens=2048, temperature=0.3)
            content = output["choices"][0]["text"]

        # Tentar converter para JSON
        json_data = json.loads(content)
        return {"analysis_json": json_data}

    except json.JSONDecodeError:
        return {"error": "Falha ao decodificar JSON da IA.", "raw_output": content}
    except Exception as e:
        return {"error": str(e)}