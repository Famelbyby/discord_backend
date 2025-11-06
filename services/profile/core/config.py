from pydantic_settings import BaseSettings


class Config(BaseSettings):

    db_name: str = "discord"

    db_link: str = f"postgres:postgres@postgres:5432/{db_name}"

    async_db_url: str = f"postgresql+asyncpg://{db_link}"


config = Config()
