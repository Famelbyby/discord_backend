from pydantic_settings import BaseSettings


class Config(BaseSettings):

    db_name: str = "company_db"

    db_link: str = f"postgres:postgres@company_db/{db_name}"

    async_db_url: str = f"postgresql+asyncpg://{db_link}"


config = Config()
