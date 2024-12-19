import telebot
from telebot.types import Message
from time import time
import requests
from utils import SpamDetector
from config import TELEGRAM_BOT_TOKEN, FASTAPI_URL


class BotHandler:
    def __init__(self):
        self.bot = telebot.TeleBot(TELEGRAM_BOT_TOKEN)
        self.spam_detector = SpamDetector()

        
        self.bot.message_handler(commands=["start"])(self.start_command)
        self.bot.message_handler(content_types=["text", "sticker"])(self.handle_message)

    def start_command(self, message: Message):
        user_id = message.from_user.id
        if self.spam_detector.is_suspicious(user_id):
            self.bot.reply_to(message, "Друг, реши капчу")
        else:
            self.bot.reply_to(message, "Как я могу вам помочь?")

    def handle_message(self, message: Message):
        user_id = message.from_user.id
        current_time = time()

        if self.spam_detector.is_suspicious(user_id):
            if message.text and self.spam_detector.verify_captcha(user_id, message.text):
                self.spam_detector.remove_captcha(user_id)
                self.bot.reply_to(message, "Ладно")
            else:
                self.bot.reply_to(message, "пересасывай")
            return

        if self.spam_detector.is_blocked(user_id, current_time):
            self.bot.reply_to(message, "Вы временно заблокированы за спам. Попробуйте позже.")
            return

       
        self.spam_detector.update_message_count(user_id, current_time)

        
        if self.spam_detector.detect_spam(user_id, current_time):
            self.bot.reply_to(message, "Вы отправили слишком много сообщений. Вы временно заблокированы.")
            return

        
        if self.spam_detector.detect_suspicious_activity(message):
            captcha_answer = 8  
            self.spam_detector.add_captcha(user_id, captcha_answer)
            self.bot.reply_to(message, "Пожалуйста, решите капчу: сколько будет 5 + 3?")
            return

        
        if message.text:
            try:
                response = requests.post(FASTAPI_URL, json={"question": message.text})
                if response.status_code == 200:
                    response_data = response.json()
                    bot_response = response_data.get("response", "Не удалось получить ответ от сервера.")
                else:
                    bot_response = f"Ошибка сервера: {response.status_code}"
            except Exception as e:
                bot_response = f"Произошла ошибка: {e}"

            self.bot.reply_to(message, bot_response)
        else:
            self.bot.reply_to(message, "Сообщение не распознано. Отправьте текстовое сообщение.")

    def run(self):
        self.bot.polling(none_stop=True)
