from fastapi import FastAPI
import uvicorn
from router import main_router
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()


app.include_router(main_router)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://127.0.0.1:3000", "http://clients_go:8000"],
    allow_credentials=True,
    allow_methods=["DELETE", "OPTIONS", "PUT", "POST", "GET"],
    allow_headers=["Accept", "Accept-Language", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "X-CSRF-Token"],
)

if __name__ == "__main__":

    uvicorn.run(
        app="main:app",
        host="0.0.0.0",
        port=9999,
        reload=True,
    )