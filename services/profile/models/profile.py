from profile.models.base import Base
from sqlalchemy.orm import Mapped, mapped_column
from datetime import datetime
from sqlalchemy import func, BigInteger
import uuid


class ProfileORM(Base):

    __tablename__ = "profiles"

    id: Mapped[str] = mapped_column(
        primary_key=True, default=lambda: str(uuid.uuid4())
    )
    mail: Mapped[str] = mapped_column(unique=True, nullable=False)
    username: Mapped[str] = mapped_column(nullable=False)
    created_at: Mapped[int] = mapped_column(
        BigInteger, default=lambda: int(datetime.now().timestamp() * 1000)
    )
    avatar: Mapped[str] = mapped_column()
    short_link: Mapped[str] = mapped_column(
        default=lambda: str(uuid.uuid4())[:9], unique=True
    )
    status: Mapped[str] = mapped_column()
