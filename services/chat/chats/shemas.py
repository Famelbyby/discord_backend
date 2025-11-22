from pydantic import BaseModel, Field, field_validator, model_validator
from typing import Optional, Literal, List

RuleType = Literal['camera', 'chat', 'call', '']

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
