class Search:
    def __init__(self, content):
        pass

    def search_similar(self, query, threshold=0.5):
        results = self.vector_store.search(query)  
        relevant_docs = [doc for doc, score in results if score >= threshold]
        if not relevant_docs:
            raise ValueError("No relevant documents found")
        return relevant_docs
