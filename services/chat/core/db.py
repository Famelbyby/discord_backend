from core.config import config

from sqlalchemy.ext.asyncio import (
    async_sessionmaker,
    create_async_engine,
)


class DataBase:
    """Класс для получения сессий и подключения к бд

    Параметры для иницилизации:
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
