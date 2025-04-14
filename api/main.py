from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from api.routes import scan
from fastapi.staticfiles import StaticFiles
app = FastAPI()
app.include_router(scan.router)


app.mount("/static", StaticFiles(directory="api/static"), name="static")


app = FastAPI(
    title="Magician Software Analyzer API",
    description="API para análise automatizada de segurança em projetos de software.",
    version="0.1.0"
)

# CORS (liberar acesso ao frontend se necessário)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Roteadores
app.include_router(scan.router, prefix="/scan")

# Rota básica
@app.get("/ping")
def ping():
    return {"status": "ok"}
