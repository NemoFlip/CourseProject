from langchain.llms import OpenAI
from langchain.callbacks import get_openai_callback
from langchain.chains.question_answering import load_qa_chain

class QAService:
    """Класс для ответа на вопросы с использованием LangChain."""
    def __init__(self):
        self.llm = OpenAI(temperature=0)
        self.chain = load_qa_chain(self.llm, chain_type="refine")

    def get_answer(self, documents, question):
        """Возвращает ответ на вопрос на основе документов."""
        with get_openai_callback():
            response = self.chain.run(input_documents=documents, question=question)
        return response
