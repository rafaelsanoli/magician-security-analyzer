from fastapi import APIRouter, UploadFile, File, Form, HTTPException, Request
from fastapi.responses import JSONResponse
import subprocess
import tempfile
import shutil
import os
import json
from pydantic import BaseModel
from api.services.analyzer import analyze_code_with_llm
from datetime import datetime
from api.services.analyzer import analyze_code_with_llm
from api.services import analyzer

router = APIRouter() #inicializa o router aqui

class CodeInput(BaseModel):
    code: str
    lang: str

@router.post("/analyze")
async def analyze_code(payload: CodeInput):
    if not payload.code.strip():
        raise HTTPException(status_code=400, detail="Código não pode estar vazio.")

    result = analyze_code_with_llm(payload.code, payload.lang)

    if "error" in result:
        raise HTTPException(status_code=500, detail=result["error"])

    return {
        "language": payload.lang,
        "analysis": result["analysis"]
    }

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
