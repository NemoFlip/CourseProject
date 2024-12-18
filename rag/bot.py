from aiogram import Bot, Dispatcher
from aiogram.types import Message
from aiogram.filters import Command
import asyncio
import logging
import requests
import os
from dotenv import load_dotenv


load_dotenv("telegram.env")
TELEGRAM_BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")
FASTAPI_URL= os.getenv("FASTAPI_URL")

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


bot = Bot(token=TELEGRAM_BOT_TOKEN)
dp = Dispatcher()


@dp.message(Command("start"))
async def start_command(message: Message):
    await message.answer("Таня на связи")


@dp.message()
async def handle_message(message: Message):
    user_question = message.text
    try:
        response = requests.post(
            FASTAPI_URL,
            json={"question": user_question}
        )
        if response.status_code == 200:
            response_data = response.json()
            bot_response = response_data.get("response", "Не удалось получить ответ от сервера.")
        else:
            bot_response = f"Ошибка сервера: {response.status_code}"
    except Exception as e:
        logger.error(f"Error: {e}")
        bot_response = "Произошла ошибка. Попробуйте позже."

    await message.answer(bot_response)


async def main():
    
    logger.info("Запуск бота...")
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())
