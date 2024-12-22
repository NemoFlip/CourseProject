from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse
import os
import cv2
import numpy as np
import torch
from insightface.app import FaceAnalysis
from sklearn.metrics.pairwise import cosine_similarity
from database import save_embedding, load_embedding, load_all_embeddings
import uuid

app = FastAPI()

face_app = FaceAnalysis(providers=['CPUExecutionProvider'])
face_app.prepare(ctx_id=0)

THRESHOLD = 0.6

DATA_DIR = ""
os.makedirs(DATA_DIR, exist_ok=True)


def save_image(user_id: str, image_bytes: bytes, action: str) -> str:
    image_path = os.path.join(DATA_DIR, f"{user_id}_{action}.jpg")
    nparr = np.frombuffer(image_bytes, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    cv2.imwrite(image_path, img)
    return image_path


def extract_face_embedding(image_bytes: bytes):
    nparr = np.frombuffer(image_bytes, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    faces = face_app.get(img)
    return faces[0].embedding if faces else None


@app.post("/register/")
async def register_user(file: UploadFile = File(...)):
    image_bytes = await file.read()
    embedding = extract_face_embedding(image_bytes)
    if embedding is None:
        return JSONResponse(status_code=400, content={"error": "Лицо не обнаружено на изображении"})
    
    user_id = str(uuid.uuid4())  
    save_embedding(user_id, embedding)
    save_image(user_id, image_bytes, action="registration")
    return {"message": "✅ Пользователь зарегистрирован"}


@app.post("/verify/")
async def verify_user(file: UploadFile = File(...)):
    image_bytes = await file.read()
    embedding = extract_face_embedding(image_bytes)
    if embedding is None:
        return JSONResponse(status_code=400, content={"error": "Лицо не обнаружено на изображении"})
    
    save_image("verification_attempt", image_bytes, action="verification")
    
    all_embeddings = load_all_embeddings()
    for user_id, user_embedding in all_embeddings.items():
        similarity = cosine_similarity([embedding], [user_embedding])[0][0]
        if similarity > THRESHOLD:
            return {"message": f"✅ Пользователь верифицирован"}
    
    return {"message": "Пользователь не опознан"}


@app.get("/")
async def root():
    return {"message": " /register/  /verify/"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8003)
