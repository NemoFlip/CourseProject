# README.md  

# RAG (Retrieval-Augmented Generation) и Telegram бот  
Проект представляет собой **RAG-систему**, которая обрабатывает текстовые данные, находит релевантные ответы на основе векторного поиска, а затем предоставляет краткие ответы пользователю через **Telegram-бота**.  

---

## 📋 **Описание проекта**  

1. **FastAPI сервер**  
   - Принимает вопросы от пользователя.  
   - Обрабатывает вопросы: переводит с **русского на английский**, модифицирует запрос с дополнительными условиями.  
   - Выполняет поиск релевантного контекста в данных с помощью **векторного хранилища**.  
   - Генерирует ответ с помощью системы вопросов и ответов (**QAService**).  
   - Переводит ответ обратно на **русский язык** и возвращает пользователю.  

2. **Telegram бот**  
   - Интегрирован с FastAPI сервером.  
   - Принимает вопросы пользователей в Telegram-чате.  
   - Делегирует запросы серверу и возвращает ответы пользователям.  

---

##  **Основные компоненты**  

### **1. FastAPI сервер**  
Файл: `main.py`  

- **Конечные точки**:  
   - `POST /query`: принимает вопрос и возвращает краткий ответ.  
   - `GET /health`: проверка статуса сервера.  

- **Основные зависимости**:  
   - `FastAPI` — веб-фреймворк для создания API.  
   - `Pydantic` — валидация данных.  
   - **Внешние модули**:  
      - `TextProcessor` — обработка текстового датасета.  
      - `Translator` — перевод текста.  
      - `VectorStoreHandler` — работа с векторным хранилищем.  
      - `QAService` — генерация ответов.  

---

### **2. Telegram бот**  
Файл: `bot.py`  

- Использует библиотеку **Aiogram** для создания асинхронного бота.  
- Передает текстовые запросы пользователя на FastAPI сервер.  
- Возвращает ответы пользователю в Telegram.  

- **Основные команды**:  
   - `/start`: приветственное сообщение.  
   - Простые текстовые сообщения обрабатываются и передаются серверу.  

---

### **4. Запуск FastAPI сервера**  

Запустите сервер с помощью Uvicorn:  

```bash
python main.py  
```

Проверьте статус сервера:  

```bash
curl http://127.0.0.1:8004/health  
```

---

### **5. Запуск Telegram бота**  

Запустите бота в отдельном терминале:  

```bash
python telegram_bot.py  
```

---

## 📄 **Примеры использования**  

1. **Запуск бота в Telegram**  
   - Откройте Telegram и найдите вашего бота по токену.  
   - Отправьте сообщение или команду `/start`.  

2. **Пример вопроса**  

   - Пользователь: `Какова длина Великой Китайской стены?`  
   - Ответ от бота:  
     > "Мой господин, длина Великой Китайской стены составляет около 21 196 км."  

---

## 🧩 **Зависимости**  

- **FastAPI**  
- **Pydantic**  
- **Aiogram**  
- **Requests**  
- **Uvicorn**  

---

## 🔗 **Структура проекта**  

```plaintext
project-name/  
│  
├── data/                       # Исходный и обработанный датасеты  
│   ├── dataset.txt  
│   └── processed_dataset.txt  
│  
├── main.py                     # FastAPI сервер  
├── telegram_bot.py             # Telegram бот  
├── text_processor.py           # Обработка текста  
├── translator.py               # Модуль для перевода  
├── vector_store_handler.py     # Векторное хранилище  
├── qa_service.py               # Генерация ответов  
│  
├── requirements.txt            # Список зависимостей  
└── README.md                   # Описание проекта  
```

---

## 💡 **Как это работает?**  

1. Telegram бот принимает запрос от пользователя и отправляет его FastAPI серверу.  
2. FastAPI обрабатывает запрос:  
   - Переводит вопрос с русского на английский.  
   - Добавляет модификацию запроса (например, ограничение на размер ответа).  
   - Выполняет поиск в векторном хранилище.  
   - Генерирует ответ на основе контекста.  
   - Переводит ответ обратно на русский язык.  
3. Ответ возвращается боту и отображается пользователю в Telegram.  

---

