from fastapi import APIRouter
from chats.router import chat_router

main_router = APIRouter(prefix="/api")

main_router.include_router(chat_router)
