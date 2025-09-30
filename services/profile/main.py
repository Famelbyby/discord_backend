from fastapi import FastAPI
import uvicorn
from profile.router import main_router

app = FastAPI()


app.include_router(main_router)


if __name__ == "__main__":

    uvicorn.run(
        app="profile.main:app",
        host="127.0.0.1",
        port=9000,
        reload=True,
    )
