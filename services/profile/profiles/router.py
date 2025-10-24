from fastapi import (
    APIRouter,
    Depends,
    UploadFile,
    Body,
    status as st,
    HTTPException,
)
from sqlalchemy.ext.asyncio import AsyncSession
from typing import Annotated, Optional
from core.db import db
from models.profile import ProfileORM
from sqlalchemy import select

from sqlalchemy.exc import IntegrityError
from pydantic import Field, EmailStr

profile_router = APIRouter(prefix="/profile")


async def check_if_profile_exists(
    id: str,
    session: Annotated[AsyncSession, Depends(db.get_async_session)],
):

    check_stmt = select(ProfileORM).where(ProfileORM.id == id)
    result = await session.execute(check_stmt)
    profile = result.scalar_one_or_none()

    if not profile:
        raise HTTPException(
            status_code=st.HTTP_404_NOT_FOUND, detail="Профиль не найден"
        )

    return profile


@profile_router.post("/", status_code=st.HTTP_201_CREATED)
async def save_profile(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    mail: Annotated[EmailStr, Field(max_length=50), Body()],
    username: Annotated[str, Field(max_length=25, min_length=4), Body()],
    avatar: Annotated[
        UploadFile | str | None,
        Field(default="https://elitclub.kz/upload/images/11333_134586_14.jpeg"),
    ] = None,
    status: Annotated[
        Optional[str],
        Field(max_length=150, default=None),
        Body(),
    ] = None,
):

    try:
        profile_orm = ProfileORM(
            mail=mail,
            username=username,
            avatar="https://elitclub.kz/upload/images/11333_134586_14.jpeg",
            status=status,
        )

        print(profile_orm.avatar)
        db.add(profile_orm)
        await db.commit()
        return {"profile": profile_orm}
    except IntegrityError as e:
        raise HTTPException(
            status_code=st.HTTP_400_BAD_REQUEST,
            detail="Такой mail уже существует",
        )

    except Exception as e:
        raise HTTPException(
            status_code=st.HTTP_400_BAD_REQUEST,
            detail=str(e),
        )


@profile_router.get("/{id}", status_code=st.HTTP_200_OK)
async def get_profile(
    profile: ProfileORM = Depends(check_if_profile_exists),
    db: AsyncSession = Depends(db.get_async_session),
):

    return {"profie": profile}


@profile_router.get("/", status_code=st.HTTP_200_OK)
async def get_profile_pagination(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    name: str = "",
    limit: int = 30,
    page: int = 1,
):

    offset = (page - 1) * limit

    statement = (
        select(ProfileORM)
        .where(ProfileORM.username.startswith(name))
        .offset(offset)
        .limit(limit)
    )
    result = await db.execute(statement)
    profile = result.scalars().all()
    if not profile:
        return {"profie": []}
    return {"profie": profile}


@profile_router.delete("/{id}", status_code=st.HTTP_204_NO_CONTENT)
async def delete_profile(
    profile: Annotated[ProfileORM, Depends(check_if_profile_exists)],
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
):

    await db.delete(profile)
    await db.commit()


@profile_router.put("/{id}", status_code=st.HTTP_201_CREATED)
async def edit_profile(
    db: Annotated[AsyncSession, Depends(db.get_async_session)],
    profile: Annotated[ProfileORM, Depends(check_if_profile_exists)],
    avatar: Annotated[
        Optional[UploadFile | str],
        Field(default="https://elitclub.kz/upload/images/11333_134586_14.jpeg"),
    ] = "https://elitclub.kz/upload/images/11333_134586_14.jpeg",
    username: Annotated[
        Optional[str], Field(max_length=25, min_length=4, default=None), Body()
    ] = None,
    status: Annotated[
        Optional[str], Field(max_length=150, default=None), Body()
    ] = None,
):

    if status:
        profile.status = status
    if username:
        profile.username = username
    profile.avatar = avatar
    await db.commit()
    await db.refresh(profile)

    return profile
