from pydantic_settings import BaseSettings


class Config(BaseSettings):

    db_name: str = "chat_db"

    db_link: str = f"postgres:postgres@localhost:5430/{db_name}"

    async_db_url: str = f"postgresql+asyncpg://{db_link}"

    kafka_url: str = "localhost:9092"

config = Config()
