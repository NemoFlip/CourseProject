# main.py
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from rag.text_processor import TextProcessor
from rag.translator import Translator
from rag.vector_store_handler import VectorStoreHandler
from rag.qa_service import QAService
import os
import uvicorn
from dotenv import load_dotenv


app = FastAPI()
load_dotenv("telegram.env")

os.environ['OPENAI_API_KEY'] = os.getenv("OPEN_AI_KEY")

path_to_dataset=os.getenv("PATH_TO_DATASET")
path_to_processed_dataset=os.getenv("PATH_TO_PROCESSED_DATASET")



processor = TextProcessor(path_to_dataset, path_to_processed_dataset)
content = processor.process_dataset()

vector_store_handler = VectorStoreHandler(content)
vector_store_handler.create_vector_store()

qa_service = QAService()
translator = Translator()

class QueryRequest(BaseModel):
    question: str
   # modification: str

@app.post("/query")
async def query(request: QueryRequest):
    try:
        user_q_translate = translator.translate(request.question, src_lang="ru", dest_lang="en")
        modification_prompt_translate = translator.translate('Ответь очень кратко и емко . колличество слов должно быть меньше 30', src_lang="ru", dest_lang="en")

        modified_question = f"{user_q_translate}\n\nContext: {modification_prompt_translate}"

        docs = vector_store_handler.search_similar(modified_question)

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







