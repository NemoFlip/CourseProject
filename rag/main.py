# main.py
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from text_processor import TextProcessor
from translator import Translator
from vector_store_handler import VectorStoreHandler
from qa_service import QAService
from search_similar import Search
import os
import uvicorn

app = FastAPI()

os.environ['OPENAI_API_KEY'] = ''

processor = TextProcessor('/Users/a1111/test/rag_prj/data/dataset.txt', '/Users/a1111/test/rag_prj/data/processed_dataset.txt')
content = processor.process_dataset()

vector_store_handler = VectorStoreHandler(content)
vector_store_handler.create_vector_store()

qa_service = QAService()
translator = Translator()

def validate_question_against_context(question, context):
    
    validation_prompt = f"Вопрос: {question}\n\nКонтекст: {context}\n\nЭтот вопрос относится к данному контексту? Ответь 'да' или 'нет'."
    response = qa_service.get_answer([context], validation_prompt)
    return response.lower().strip() == "да"


class QueryRequest(BaseModel):
    question: str
   # modification: str

'''@app.post("/query")
async def query(request: QueryRequest):
    try:
        user_q_translate = translator.translate(request.question, src_lang="ru", dest_lang="en")
        modification_prompt_translate = translator.translate('Ответь очень кратко и емко.колличество слов должно быть меньше 30 .Отвечай исключительно на основе данного текста ', src_lang="ru", dest_lang="en")

        modified_question = f"{user_q_translate}\n\nContext: {modification_prompt_translate}"

        docs = vector_store_handler.search_similar(modified_question)

        response = qa_service.get_answer(docs, modified_question)

        response_translate = translator.translate(response, src_lang="en", dest_lang="ru")

        return {"response": response_translate}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))'''
@app.post("/query")
async def query(request: QueryRequest):
    try:
        user_q_translate = translator.translate(request.question, src_lang="ru", dest_lang="en")
        modification_prompt_translate = translator.translate(
            'Ответь очень кратко и емко.колличество слов должно быть меньше 30. Отвечай исключительно на основе данного текста.',
            src_lang="ru", dest_lang="en"
        )

        modified_question = f"{user_q_translate}\n\nContext: {modification_prompt_translate}"

        docs = vector_store_handler.search_similar(modified_question)

        if not validate_question_against_context(user_q_translate, docs[0]):
            return {"response": "Вопрос не по теме"}

        response = qa_service.get_answer(docs, modified_question)
        response_translate = translator.translate(response, src_lang="en", dest_lang="ru")

        return {"response": response_translate}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    """Проверка статуса сервера."""
    return {"status": "OK"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8004)






