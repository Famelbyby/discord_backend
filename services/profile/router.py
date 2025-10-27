from fastapi import APIRouter
from profiles.router import profile_router

main_router = APIRouter(prefix="/api")

main_router.include_router(profile_router)
