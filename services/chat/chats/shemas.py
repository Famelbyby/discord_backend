from pydantic import BaseModel, Field, field_validator, model_validator
from typing import Optional, Literal, List

from core.db import db

from models.chat import ChatORM, Association
from typing import Optional
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from sqlalchemy.orm import selectinload

RuleType = Literal['camera', 'chat', 'call']
MType = Literal['SEND_MESSAGE','DELETE_MESSAGE','EDIT_MESSAGE']
DEFAULT_RULES: List[RuleType] = ['camera', 'chat', 'call']
MAX_USERS = 10


class UserInChat(BaseModel):
    id: str
    rules: Optional[List[RuleType]] = Field(
        description="Права пользователя в чате", default=DEFAULT_RULES
    )

    # @field_validator('rules')
    # @classmethod
    # def validate_rules(cls, rules: List[RuleType]) -> List[RuleType]:
    #     if not rules:
    #         return DEFAULT_RULES
    #     return rules


class ChatCreate(BaseModel):

    name: Optional[str] = Field(
        description="Название чата",
        max_length=100,
        default="Название чата",  # потом переделать в default_factory
    )
    lead_id: str
    users: List[UserInChat]

    @field_validator("users")
    @classmethod
    def validate_users(cls, users: List[UserInChat]):

        if len(users) > MAX_USERS:
            raise ValueError(
                f"Maximum {MAX_USERS} users allowed, got {len(users)}"
            )
        return users

    @model_validator(mode='after')
    def add_lead_user(self) -> 'ChatCreate':

        lead_exists = any(user.id == self.lead_id for user in self.users)

        if not lead_exists:

            lead_user = UserInChat(id=self.lead_id, rules=list(DEFAULT_RULES))
            self.users.append(lead_user)
        else:

            for user in self.users:
                if user.id == self.lead_id:
                    missing_rules = set(DEFAULT_RULES) - set(user.rules)
                    if missing_rules:
                        raise ValueError(
                            f"Lead user must have all rules. Missing: {list(missing_rules)}"
                        )
                    break

        return self


class ChatDisplay(ChatCreate):
    id: str


class ChatUserResponse(BaseModel):
    id: str
    rules: List[str]


class ChatResponse(BaseModel):
    id: str
    name: str
    lead_id: str
    users: List[ChatUserResponse]


class ChatUpdateRequest(BaseModel):
    name: Optional[str] = Field(
        description="Название чата",
        max_length=100,
        default="Название чата",
    )
    users: List[UserInChat] = Field(default_factory=list)

    @field_validator('users')
    @classmethod
    def validate_users_count(cls, v):
        if len(v) > 10:
            raise ValueError('Maximum 10 users allowed in chat')
        return v


class ChatUpdateResponse(BaseModel):
    id: str
    name: str
    lead_id: str
    users: List[UserInChat]


class MessaggePayload(BaseModel):
    message: str
    messageId: Optional[str] = Field(default=None)


class KafkaProduceMessage(BaseModel):
    userId: str
    type: MType
    chatId: str
    receiverIds: List[str]
    payload: MessaggePayload
    
async def format_chat_response(
    chat: ChatORM, db: AsyncSession
) -> ChatUpdateResponse:
    """Форматирует ответ с информацией о чате"""
    # Загружаем пользователей чата
    stmt = (
        select(Association)
        .where(Association.chat_id == chat.id)
        .options(selectinload(Association.user))
    )
    result = await db.execute(stmt)
    associations = result.scalars().all()

    users_response = []
    for association in associations:
        users_response.append(
            UserInChat(
                id=association.user.main_id,
                rules=(
                    association.rules.split(',') if association.rules else []
                ),
            )
        )

    return ChatUpdateResponse(
        id=chat.id, name=chat.name, lead_id=chat.lead_id, users=users_response
    )
