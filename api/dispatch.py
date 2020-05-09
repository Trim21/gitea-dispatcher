import os

from fastapi import FastAPI, Header
from fastapi.responses import JSONResponse
import httpx
import pydantic

secret = os.getenv("SECRET")
githubToken = "token " + os.getenv("GH_TOKEN")
client = httpx.AsyncClient(headers={"Authorization": githubToken})

app = FastAPI()


class Payload(pydantic.BaseModel):
    secret: str


@app.post("/dispatch", response_class=JSONResponse)
async def post(payload: Payload, event=Header(..., alias="X-Gitea-Event")):
    if secret == payload.secret:
        r = await client.post(
            "https://api.github.com/repos/Trim21/actions-cron/dispatches",
            json={"event_type": "%s push"},
        )
        return r.json()
    return {"error": "secret mismatch"}


@app.get("/dispatch", response_class=JSONResponse)
async def index():
    return {"hello": "world"}
