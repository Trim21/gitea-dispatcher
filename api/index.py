import os

from fastapi import FastAPI, Header
from fastapi.responses import JSONResponse

app = FastAPI()


@app.get("/", response_class=JSONResponse)
async def index():
    return {"hello": "world"}
