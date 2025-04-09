from fastapi import APIRouter, UploadFile, File, Form
from fastapi.responses import JSONResponse
import subprocess
import tempfile
import shutil
import os
import json
from datetime import datetime

router = APIRouter()

@router.post("/")
async def scan_project(
        fix: bool = Form(False),
        pr: bool = Form(False),
        file: UploadFile = File(...)
):
    try:
        # Diretórios
        temp_dir = tempfile.mkdtemp()
        zip_path = os.path.join(temp_dir, file.filename)
        report_id = datetime.now().strftime("%Y%m%d%H%M%S")
        report_dir = os.path.abspath("api/static/reports")
        os.makedirs(report_dir, exist_ok=True)

        # Salvar e extrair ZIP
        with open(zip_path, "wb") as f:
            f.write(await file.read())
        shutil.unpack_archive(zip_path, temp_dir)

        # Rodar análise via Go
        args = ["go", "run", "main.go", "scan"]
        if fix:
            args.append("--fix")
        if pr:
            args.append("--pr")
        result = subprocess.run(args, cwd=temp_dir, capture_output=True, text=True)

        # Verificar e carregar resultados
        result_path = os.path.join(temp_dir, "results.json")
        if not os.path.exists(result_path):
            return JSONResponse(content={"error": "Análise falhou."}, status_code=500)
        with open(result_path) as f:
            findings = json.load(f)

        # Gerar HTML com render.py
        html_path = os.path.join(report_dir, f"report_{report_id}.html")
        subprocess.run([
            "python3", "scripts/render.py", result_path, html_path
        ])

        # Limpeza do tmp dir
        shutil.rmtree(temp_dir)

        return {
            "status": "ok",
            "findings": findings,
            "report_url": f"/static/reports/report_{report_id}.html"
        }

    except Exception as e:
        return JSONResponse(content={"error": str(e)}, status_code=500)
