import requests

URL = "http://localhost:8000/scan/analyze"

code_example = """
import os
import sqlite3
import json

# Função 1: Injeção de SQL
def get_user_data(user_id):
    conn = sqlite3.connect('example.db')
    cursor = conn.cursor()
    query = f"SELECT * FROM users WHERE user_id = {user_id}"  # Injeção de SQL
    cursor.execute(query)
    return cursor.fetchall()

# Função 2: Injeção de comandos
def execute_system_command(command):
    os.system(f"echo {command}")  # Injeção de comando de shell

# Função 3: Falta de validação de entrada
def process_input(user_input):
    result = eval(user_input)  # Execução de código arbitrário sem validação

# Função 4: Exposição de dados sensíveis
def save_user_data(user_data):
    with open('user_data.json', 'w') as f:
        json.dump(user_data, f)  # Dados sensíveis podem ser expostos

# Função 5: Uso de funções obsoletas
def run_obsolete_function():
    os.popen("echo Hello World")  # Função obsoleta

# Função 6: Falhas de autenticação e autorização
def access_sensitive_data(user):
    if user != "admin":
        return "Acesso negado"  # Falta de validação de permissões
    return "Dados sensíveis"

# Função principal que usa as funções acima
def main():
    user_id = "1 OR 1=1"  # Testando injeção de SQL
    command = "ls"  # Testando injeção de comando
    user_input = "__import__('os').system('echo Hacked')"  # Testando execução de código arbitrário
    user_data = {"name": "John Doe", "password": "secret"}  # Dados sensíveis
    user = "guest"  # Falha de autenticação

    print(get_user_data(user_id))
    execute_system_command(command)
    process_input(user_input)
    save_user_data(user_data)
    run_obsolete_function()
    print(access_sensitive_data(user))

if __name__ == "__main__":
    main()
"""

payload = {
    "code": code_example,
    "lang": "python"
}

response = requests.post(URL, json=payload)

print("Resposta completa da API:") # adiciona o print da resposta completa
try:
    print(response.json())
except Exception:
    print("❌ Resposta não é JSON:")
    print(response.text) # adiciona o print da resposta completa

if response.status_code == 200:
    print("✅ Análise recebida:")
    analysis = response.json()["analysis"]
    print(analysis)
    if "injeção de comando" in analysis.lower():
        print("✅ Vulnerabilidade de injeção de comando detectada!")
    else:
        print("❌ Vulnerabilidade de injeção de comando não detectada!")
else:
    print("❌ Erro ao analisar:")
    print(response.status_code, response.text)