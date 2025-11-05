from pydantic_settings import BaseSettings


class MinioSettins(BaseSettings):

    minio_endpoint: str = "minio:9000"
    minio_url: str = "localhost:9000"
    minio_access_key: str = "minioadmin"
    minio_secret_key: str = "minioadmin"
    minio_secure: bool = False
    bucket_name: str = "pictures"

    class Config:
        env_file = ".env"


minio_settings = MinioSettins()
