from fastapi import FastAPI
import os

app = FastAPI()


@app.get("/health")
def health():
    return {"status": "ok", "env": os.getenv("ENV", "dev")}


@app.get("/")
def root():
    return {"message": "Backend Python funcionando"}
