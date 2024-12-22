import os
from dotenv import load_dotenv

load_dotenv()

TELEGRAM_BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN", "7930398169:AAF3afBpBRDk53E_9lcw3NI4UTlx9nzDtcs")
FASTAPI_URL = os.getenv("FASTAPI_URL", "http://127.0.0.1:8004/query")
