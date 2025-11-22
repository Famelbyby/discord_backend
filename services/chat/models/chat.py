from models.base import Base
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy import ForeignKey, String
import uuid
from typing import List


class Association(Base):
    __tablename__ = "association_table"
    chat_id: Mapped[int] = mapped_column(
        ForeignKey("chats.id"), primary_key=True
    )
    user_id: Mapped[int] = mapped_column(
        ForeignKey("users.id"), primary_key=True
    )
    rules: Mapped[str] = mapped_column(nullable=False)
    chat: Mapped["ChatORM"] = relationship(back_populates="users")
    user: Mapped["UserORM"] = relationship(back_populates="chats")


class ChatORM(Base):

    __tablename__ = "chats"

    id: Mapped[str] = mapped_column(
        primary_key=True, default=lambda: str(uuid.uuid4())
    )
    name: Mapped[str] = mapped_column(String(100), nullable=False)
    lead_id: Mapped[str] = mapped_column(nullable=False)
    users: Mapped[List["Association"]] = relationship(back_populates="chat")


class UserORM(Base):

    __tablename__ = "users"

    id: Mapped[str] = mapped_column(
        primary_key=True, default=lambda: str(uuid.uuid4())
    )
    main_id: Mapped[str] = mapped_column(nullable=False)
    chats: Mapped[List["Association"]] = relationship(back_populates="user")
