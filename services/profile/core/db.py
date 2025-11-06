from core.config import config

from fastapi import Depends

from sqlalchemy import create_engine
from sqlalchemy.ext.asyncio import (
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)
from sqlalchemy.orm import sessionmaker


class DataBase:
    """Класс для получения сессий и подключения к бд

    Параметры для иницилизации:
    db_url - ссылка подключения к бд
    async_db_url - ссылка для асинхронного подключения к бд
    """

    def __init__(self, async_db_url: str):

        self.async_engine = create_async_engine(async_db_url)

        self.async_session_factory = async_sessionmaker(
            bind=self.async_engine,
            autoflush=False,
            autocommit=False,
            expire_on_commit=False,
        )

    async def get_async_session(self):
        async with self.async_session_factory() as session:
            yield session


db = DataBase(config.async_db_url)
