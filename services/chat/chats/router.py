from fastapi import APIRouter, Depends, HTTPException, status, Path, Query
from core.db import db
from chats.shemas import (
    ChatCreate,
    ChatDisplay,
    ChatResponse,
    ChatUserResponse,
    ChatUpdateRequest,
    ChatUpdateResponse,
    format_chat_response,
)
from models.chat import ChatORM, UserORM, Association
from typing import Annotated, Optional
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, delete
from sqlalchemy.orm import selectinload

chat_router = APIRouter(prefix="/chat")


async def get_chat_by_id(
    chat_id: Annotated[str, Path()],
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
) -> ChatORM:
    statement = (
        select(ChatORM)
        .options(selectinload(ChatORM.users).selectinload(Association.user))
        .where(ChatORM.id == chat_id)
    )
    result = await db.execute(statement)
    chat = result.scalar_one_or_none()
    return chat


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


@chat_router.get("/{chat_id}")
async def get_chat(
    chat: Annotated[ChatORM, Depends(get_chat_by_id)], 
    user_id: str,
):

    if not chat:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Chat not found"
        )

    if chat.lead_id != user_id:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Chat not found"
        )

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


@chat_router.delete(
    "/{chat_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_chat(
    chat: Annotated[ChatORM, Depends(get_chat_by_id)],
    user_id: str,
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
):

    if not chat:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Chat not found"
        )

    if chat.lead_id != user_id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You're not permitted to delete this chat",
        )

    await db.delete(chat)
    await db.commit()


@chat_router.get("/")
async def get_user_chats(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    user_id: str,
    chat_name: Optional[str] = None,
):

    user_stmt = select(UserORM.id).where(UserORM.main_id == user_id)
    user_result = await db.execute(user_stmt)
    user_id_internal = user_result.scalar_one_or_none()

    if not user_id_internal:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="User not found"
        )

    chat_stmt = (
        select(ChatORM)
        .join(Association, ChatORM.id == Association.chat_id)
        .where(Association.user_id == user_id_internal)
    )

    if chat_name:
        chat_stmt = chat_stmt.where(ChatORM.name.startswith(chat_name))

    chat_stmt = chat_stmt.options(
        selectinload(ChatORM.users).selectinload(Association.user)
    )

    chat_result = await db.execute(chat_stmt)
    chats_orm = chat_result.scalars().all()

    chats = []
    for chat in chats_orm:
        chat_users = []
        for association in chat.users:
            chat_users.append(
                ChatUserResponse(
                    id=association.user.main_id,
                    rules=(
                        association.rules.split(",")
                        if association.rules
                        else []
                    ),
                )
            )

        chats.append(
            ChatResponse(
                id=chat.id,
                name=chat.name,
                lead_id=chat.lead_id,
                users=chat_users,
            )
        )

    return chats


@chat_router.delete(
    "/{chat_id}/delete_user",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_user_from_chat(
    user_id: str,
    lead_id: str,
    chat: Annotated[ChatORM, Depends(get_chat_by_id)],
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
):
    if not chat:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Chat doesn't exist")

    if chat.lead_id != lead_id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You are not leader of chat",
        )
        
    if chat.lead_id == user_id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You can not delete yourself"
        )

    user_stmt = select(UserORM).where(UserORM.main_id == user_id)
    user_result = await db.execute(user_stmt)
    user = user_result.scalar_one_or_none()

    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found",
        )

    association_stmt = (
        select(Association)
        .where(Association.chat_id == chat.id)
        .where(Association.user_id == user.id)
    )
    association_result = await db.execute(association_stmt)
    association = association_result.scalar_one_or_none()

    if not association:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User is not a member of this chat",
        )

    delete_stmt = (
        delete(Association)
        .where(Association.chat_id == chat.id)
        .where(Association.user_id == user.id)
    )

    await db.execute(delete_stmt)
    await db.commit()


@chat_router.put(
    "/{chat_id}",
    response_model=ChatUpdateResponse,
    status_code=status.HTTP_200_OK,
)
async def update_chat(
    chat_update: ChatUpdateRequest,
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    chat: Annotated[ChatORM, Depends(get_chat_by_id)],
    user_id: str = Query(),
):

    if not chat:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Chat not found"
        )

    if chat.lead_id != user_id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You are not the leader of this chat",
        )

    try:

        if chat_update.name != "Название чата":
            chat.name = chat_update.name

        lead_user_stmt = select(UserORM.id).where(
            UserORM.main_id == chat.lead_id
        )
        lead_user_result = await db.execute(lead_user_stmt)
        lead_user_id = lead_user_result.scalar_one()

        await db.execute(
            delete(Association).where(
                Association.chat_id == chat.id,
                Association.user_id != lead_user_id,
            )
        )

        for user_data in chat_update.users:

            if user_data.id == chat.lead_id:
                continue

            user_stmt = select(UserORM).where(UserORM.main_id == user_data.id)
            user_result = await db.execute(user_stmt)
            user = user_result.scalar_one_or_none()

            if not user:
                user = UserORM(main_id=user_data.id)
                db.add(user)
                await db.flush()

            association = Association(
                chat_id=chat.id,
                user_id=user.id,
                rules=','.join(user_data.rules),
            )
            db.add(association)

        await db.commit()
        await db.refresh(chat)

        return await format_chat_response(chat, db)

    except Exception as e:
        await db.rollback()
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Error updating chat: {str(e)}",
        )
