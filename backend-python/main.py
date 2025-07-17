from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import asyncpg
import redis.asyncio as redis
import os
import json

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

DB_DSN = f"postgresql://{os.getenv('DB_USER')}:{os.getenv('DB_PASSWORD')}@{os.getenv('DB_HOST')}:{os.getenv('DB_PORT')}/{os.getenv('DB_NAME')}"
REDIS_URL = f"redis://{os.getenv('REDIS_HOST')}:{os.getenv('REDIS_PORT')}"

class User(BaseModel):
    name: str

@app.on_event("startup")
async def startup():
    app.state.db = await asyncpg.create_pool(DB_DSN)
    app.state.redis = redis.from_url(REDIS_URL, decode_responses=True)
    async with app.state.db.acquire() as conn:
        await conn.execute("""
            CREATE TABLE IF NOT EXISTS users_python (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL
            )
        """)

@app.get("/users")
async def get_users():
    cache = await app.state.redis.get("users_python_cache")
    if cache:
        return json.loads(cache)

    rows = await app.state.db.fetch("SELECT id, name FROM users_python")
    users = [{"id": r["id"], "name": r["name"]} for r in rows]
    await app.state.redis.set("users_python_cache", json.dumps(users))
    return users

@app.post("/users")
async def create_user(user: User):
    async with app.state.db.acquire() as conn:
        await conn.execute("INSERT INTO users_python(name) VALUES($1)", user.name)
    await app.state.redis.delete("users_python_cache")
    return {"status": "created"}
