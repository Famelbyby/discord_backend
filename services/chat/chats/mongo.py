from pymongo import MongoClient
from core.config import config
from chats.shemas import MessageSave


class MongoDocument:
    def __init__(self):

        self.client = MongoClient(
            host=config.mongo_host,
            port=config.mongo_port,
        )
        self.db = self.client["messages-db"]
        self.collection = self.db["messasges"]

    def insert_chat_message(self, data: MessageSave):
        result = self.collection.insert_one(data.model_dump())
        return result

    def find_chat_messages(self, chat_id: str):
        """
        Найти сообщения чата по chatId и тексту payload

        Returns:
            List[dict]: список документов
        """
        query = {"chatId": chat_id}
        return self.collection.find(query)


mongo_document = MongoDocument()
