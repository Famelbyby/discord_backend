from pydantic_settings import BaseSettings


class Config(BaseSettings):

    db_name: str = "chat_db"

    db_link: str = f"postgres:postgres@postgres:5432/{db_name}"

    async_db_url: str = f"postgresql+asyncpg://{db_link}"

    kafka_url: str = "kafka-container:9092"

    mongo_host: str = "mongo_db"

    mongo_port: int = 27017


config = Config()
