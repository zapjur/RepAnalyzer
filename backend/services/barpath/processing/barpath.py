import cv2
import requests
import os
from ultralytics import YOLO
from config import MODEL_PATH, MINIO_CLIENT, REPLY_QUEUE
from urllib.parse import urlparse
from .utils import extract_filename
from rabbit.publisher import publish_result

async def process_task(data):
    video_id = data["video_id"]
    bucket = data["bucket"]
    object_key = data["object_key"]
    auth0_id = data["auth0_id"]
    exercise = data["exercise_name"]
    reply_queue = data["reply_queue"]

    try:
        video_url = f"http://minio:9000/{bucket}/{object_key}"
        response = requests.get(video_url, stream=True)
        temp_filename = f"/tmp/{video_id}.mp4"
        with open(temp_filename, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)

        model = YOLO(MODEL_PATH)

        cap = cv2.VideoCapture(temp_filename)
        traj = []

        fps = cap.get(cv2.CAP_PROP_FPS)
        width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')
        output_file = f"/tmp/{video_id}_out.mp4"
        out = cv2.VideoWriter(output_file, fourcc, fps, (width, height))

        while cap.isOpened():
            ret, frame = cap.read()
            if not ret:
                break

            results = model.predict(source=frame, conf=0.5, iou=0.4, classes=[0], verbose=False)
            boxes = results[0].boxes

            if boxes is not None and len(boxes) > 0:
                largest_box = max(boxes.data, key=lambda b: (b[2] - b[0]) * (b[3] - b[1]))
                x1, y1, x2, y2 = map(int, largest_box[:4])
                cx = (x1 + x2) // 2
                cy = (y1 + y2) // 2
                traj.append((cx, cy))
                cv2.circle(frame, (cx, cy), 4, (0, 255, 0), -1)

            for i in range(1, len(traj)):
                cv2.line(frame, traj[i - 1], traj[i], (255, 0, 0), 2)

            out.write(frame)

        cap.release()
        out.release()
        cv2.destroyAllWindows()

        filename = extract_filename(video_url)
        output_key = f"{auth0_id}/{exercise}/barpath/{filename}.mp4"
        MINIO_CLIENT.fput_object(bucket, output_key, output_file)

        result = {
            "video_id": video_id,
            "status": "success",
            "bucket": bucket,
            "object_key": output_key
        }

    except Exception as e:
        result = {
            "video_id": video_id,
            "status": "error",
            "message": str(e)
        }

    await publish_result(reply_queue, result)
