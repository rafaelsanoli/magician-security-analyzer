import json
import sys
from jinja2 import Environment, FileSystemLoader
from weasyprint import HTML
from datetime import datetime
import os

# Caminho base
BASE_DIR = os.path.abspath(os.path.dirname(__file__))
TEMPLATE_DIR = os.path.join(BASE_DIR, "../templates")
OUTPUT_DIR = os.path.join(BASE_DIR, "../report")

# Criar pasta de saída se necessário
os.makedirs(OUTPUT_DIR, exist_ok=True)

# Carregar resultados
def load_results(json_path):
    with open(json_path, 'r') as f:
        return json.load(f)

# Gerar HTML com Jinja2
def generate_html(data, output_html):
    env = Environment(loader=FileSystemLoader(TEMPLATE_DIR))
    template = env.get_template("report.html.j2")
    rendered = template.render(results=data, date=datetime.now())
    with open(output_html, "w", encoding="utf-8") as f:
        f.write(rendered)
    print(f"[+] HTML gerado em: {output_html}")

# Exportar PDF com WeasyPrint
def generate_pdf(html_path, output_pdf):
    HTML(html_path).write_pdf(output_pdf)
    print(f"[+] PDF gerado em: {output_pdf}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Uso: python render.py results.json")
        sys.exit(1)

    input_json = sys.argv[1]
    html_path = os.path.join(OUTPUT_DIR, "report.html")
    pdf_path = os.path.join(OUTPUT_DIR, "report.pdf")

    results = load_results(input_json)
    generate_html(results, html_path)
    generate_pdf(html_path, pdf_path)
