import cv2

def capture_face_image():
    cap = cv2.VideoCapture(0)
    ret, frame = cap.read()
    cap.release()
    if not ret:
        raise Exception("Ошибка при захвате изображения с камеры")
    return frame
