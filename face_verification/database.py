import os
import pickle

DATABASE_DIR = '/Users/a1111/test/face_verification/embeddings'
os.makedirs(DATABASE_DIR, exist_ok=True)

def save_embedding(user_id, embedding):
    file_path = os.path.join(DATABASE_DIR, f"{user_id}.pkl")
    with open(file_path, 'wb') as f:
        pickle.dump(embedding, f)

def load_embedding(user_id):
    file_path = os.path.join(DATABASE_DIR, f"{user_id}.pkl")
    if not os.path.exists(file_path):
        return None
    with open(file_path, 'rb') as f:
        return pickle.load(f)

def load_all_embeddings():
    embeddings = {}
    for file in os.listdir(DATABASE_DIR):
        if file.endswith(".pkl"):
            user_id = os.path.splitext(file)[0]
            with open(os.path.join(DATABASE_DIR, file), 'rb') as f:
                embeddings[user_id] = pickle.load(f)
    return embeddings

