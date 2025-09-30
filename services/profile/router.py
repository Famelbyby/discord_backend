from fastapi import APIRouter
from profile.profiles.router import profile_router

main_router = APIRouter(prefix="/api")

main_router.include_router(profile_router)
