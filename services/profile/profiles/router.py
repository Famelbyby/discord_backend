from fastapi import APIRouter


profile_router = APIRouter(prefix="/profile")


@profile_router.get("/")
async def get_profile():
    return {"profile": "name"}
