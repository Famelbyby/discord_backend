from fastapi import APIRouter, Depends, HTTPException, status
from core.db import db
from chats.shemas import ChatCreate, ChatDisplay
from models.chat import ChatORM, UserORM, Association
from typing import Annotated
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from sqlalchemy.orm import selectinload

chat_router = APIRouter(prefix="/chat")


@chat_router.post("/", status_code=status.HTTP_201_CREATED)
async def create_chat(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    chat: ChatCreate,
):
    try:
        created_chat = ChatORM(name=chat.name, lead_id=chat.lead_id)

        db.add(created_chat)
        await db.flush()

        for user_data in chat.users:

            stmt = select(UserORM).where(UserORM.main_id == user_data.id)
            result = await db.execute(stmt)
            user = result.scalar_one_or_none()

            if user is None:
                user = UserORM(main_id=user_data.id)
                db.add(user)
                await db.flush()

            association = Association(
                chat_id=created_chat.id,
                user_id=user.id,
                rules=','.join(user_data.rules),
            )

            db.add(association)
        await db.commit()
        await db.refresh(created_chat)

        return {
            "id": created_chat.id,
            "name": created_chat.name,
            "lead_id": created_chat.lead_id,
            "users": chat.users,
        }

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@chat_router.get("/{id}")
async def get_chat_by_id(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    id: str,
):

    stmt = (
        select(ChatORM)
        .options(selectinload(ChatORM.users).selectinload(Association.user))
        .where(ChatORM.id == id)
    )
    result = await db.execute(stmt)
    chat = result.scalar_one()

    users = [
        {
            "id": user.user.main_id,
            "rules": (
                user.rules.split(",") if user.rules.split(",") != [""] else []
            ),
        }
        for user in chat.users
    ]
    chat_display = ChatDisplay(
        id=chat.id, name=chat.name, lead_id=chat.lead_id, users=users
    )
    return {"chat": chat_display}
