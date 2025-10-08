from pydantic import BaseModel, Field, EmailStr
from typing import Annotated, Union
from fastapi import UploadFile


class ProfileCreateRequest(BaseModel):

    mail: Annotated[EmailStr, Field(max_length=50)]
    username: Annotated[str, Field(max_length=25, min_length=4)]
    status: Annotated[str, Field(max_length=150)]
    avatar: Annotated[
        Union[str, UploadFile],
        Field(default="https://elitclub.kz/upload/images/11333_134586_14.jpeg"),
    ]
