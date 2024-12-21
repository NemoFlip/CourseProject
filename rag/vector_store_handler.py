from langchain.text_splitter import CharacterTextSplitter
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.vectorstores import FAISS

class VectorStoreHandler:
    def __init__(self, content):
        self.content = content
        self.vector_store = None

    def create_vector_store(self):
        text_splitter = CharacterTextSplitter(separator='\n', chunk_size=300, chunk_overlap=128, length_function=len)
        chunks = text_splitter.split_text(self.content)
        embeddings = OpenAIEmbeddings()
        self.vector_store = FAISS.from_texts(chunks, embeddings)
        self.vector_store.save_local("faiss_doc_idx")

    def search_similar(self, query):
        if not self.vector_store:
            self.vector_store = FAISS.load_local("faiss_doc_idx", OpenAIEmbeddings())
        return self.vector_store.similarity_search(query)
