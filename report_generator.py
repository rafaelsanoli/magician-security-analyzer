from jinja2 import Environment, FileSystemLoader
import uuid
import os

def gerar_relatorio_html(analysis_ia_json: list, output_dir="api/static/reports") -> str:
    env = Environment(loader=FileSystemLoader("api/templates"))
    template = env.get_template("report_template.html")

    report_id = str(uuid.uuid4())
    rendered = template.render(analysis_ia=analysis_ia_json)

    os.makedirs(output_dir, exist_ok=True)
    output_path = os.path.join(output_dir, f"report_{report_id}.html")
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(rendered)

    return output_path
