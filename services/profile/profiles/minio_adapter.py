from minio import Minio
from minio.error import S3Error
from core.minio import minio_settings
from fastapi import UploadFile, HTTPException, status
import uuid
from io import BytesIO
import json


class MinioAdapter:

    ALLOWED_TYPES = {
        'image/jpeg',
        'image/png',
        'image/gif',
        'image/webp',
        'image/svg',
    }

    def __init__(self):
        self.client = Minio(
            endpoint=minio_settings.minio_endpoint,
            secret_key=minio_settings.minio_secret_key,
            access_key=minio_settings.minio_access_key,
            secure=minio_settings.minio_secure,
        )
        self.bucket_name = minio_settings.bucket_name
        self._ensure_bucket_exists()

    def _ensure_bucket_exists(self):
        """Создает bucket если он не существует"""
        try:
            if not self.client.bucket_exists(self.bucket_name):
                self.client.make_bucket(self.bucket_name)

                print(f"Bucket '{self.bucket_name}' создан успешно")
                self._set_bucket_public_policy()
        except S3Error as e:
            print(f"Ошибка при создании bucket: {e}")

    def _set_bucket_public_policy(self):
        """Устанавливает политику публичного доступа для bucket"""
        public_policy = {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": "*"},
                    "Action": "s3:GetObject",
                    "Resource": f"arn:aws:s3:::{self.bucket_name}/*",
                },
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": "*"},
                    "Action": "s3:ListBucket",
                    "Resource": f"arn:aws:s3:::{self.bucket_name}",
                },
            ],
        }

        try:
            self.client.set_bucket_policy(
                self.bucket_name, json.dumps(public_policy)
            )
            print(f"Публичный доступ для bucket '{self.bucket_name}' настроен")
        except S3Error as e:
            print(f"Ошибка при установке политики: {e}")

    async def upload_file(self, file: UploadFile):

        file_type = file.content_type
        if file_type not in self.ALLOWED_TYPES:
            raise HTTPException(
                status_code=status.HTTP_422_UNPROCESSABLE_CONTENT,
                detail=f"Неправильный формат файла: {file_type}",
            )

        file_name = str(uuid.uuid4()) + "." + file_type.split("/")[-1]

        file_data = await file.read()
        file_bytes = BytesIO(file_data)

        self.client.put_object(
            self.bucket_name,
            object_name=file_name,
            data=file_bytes,
            length=file.size,
            content_type=file_type,
        )
        return f"http://{minio_settings.minio_url}/{minio_settings.bucket_name}/{file_name}"


minio_adapter = MinioAdapter()