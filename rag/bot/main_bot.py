import logging
from bot import BotHandler

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

if __name__ == "__main__":
    logger.info("Запуск бота...")
    bot_instance = BotHandler()
    bot_instance.run()
