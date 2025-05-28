import pika
import requests
import cv2
from ultralytics import YOLO
from minio import Minio
import os
import uuid
import json
from urllib.parse import urlparse
import time
import logging

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")

MODEL_PATH = "runs/detect/train6/weights/best.pt"
MINIO_CLIENT = Minio("minio:9000", access_key="admin", secret_key="admin123", secure=False)

def extract_filename(url):
    path = urlparse(url).path
    base = os.path.basename(path)
    filename, _ = os.path.splitext(base)
    return filename

def process_video(bucket, video_url, video_id, auth0_id, exercise):
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
    object_key = f"{auth0_id}/{exercise}/barpath/{filename}.mp4"
    MINIO_CLIENT.fput_object(bucket, object_key, output_file)

    return bucket, object_key

def callback(ch, method, properties, body):
    data = json.loads(body)
    video_id = data["video_id"]
    bucket = data["bucket"]
    object_key = data["object_key"]
    reply_queue = data["reply_queue"]
    auth0_id = data["auth0_id"]
    exercise = data["exercise_name"]

    video_url = f"http://minio:9000/{bucket}/{object_key}"

    print(f"Received task {video_id}")
    try:
        bucket, object_key = process_video(bucket, video_url, video_id, auth0_id, exercise)
        result = {
            "video_id": video_id,
            "status": "success",
            "bucket": bucket,
            "object_key": object_key
        }
    except Exception as e:
        result = {
            "video_id": video_id,
            "status": "error",
            "message": str(e)
        }

    ch.basic_publish(exchange='', routing_key=reply_queue, body=json.dumps(result))
    ch.basic_ack(delivery_tag=method.delivery_tag)


def connect_rabbitmq():
    credentials = pika.PlainCredentials("guest", "guest")
    parameters = pika.ConnectionParameters(
        host="rabbitmq",
        port=5672,
        credentials=credentials
    )

    for attempt in range(10):
        try:
            logging.info(f"Attempt {attempt + 1}/10: Connecting to RabbitMQ...")
            connection = pika.BlockingConnection(parameters)
            logging.info("Connected to RabbitMQ!")
            return connection
        except pika.exceptions.AMQPConnectionError:
            logging.warning("RabbitMQ not ready yet. Retrying in 3 seconds...")
            time.sleep(3)

    raise RuntimeError("Failed to connect to RabbitMQ after multiple attempts")

def wait_for_queue(connection, queue_name):
    for attempt in range(10):
        try:
            channel = connection.channel()
            channel.queue_declare(queue=queue_name, passive=True)
            logging.info(f"Queue '{queue_name}' exists.")
            return channel
        except pika.exceptions.ChannelClosedByBroker:
            logging.warning(f"Queue '{queue_name}' not found. Retrying in 3 seconds...")
            time.sleep(3)
        except Exception as e:
            logging.error(f"Unexpected error while checking queue: {e}")
            time.sleep(3)
    raise RuntimeError(f"Queue '{queue_name}' did not appear after multiple attempts")

connection = connect_rabbitmq()
channel = wait_for_queue(connection, "bar_path_queue")

channel.basic_consume(queue="bar_path_queue", on_message_callback=callback)

logging.info("Waiting for tasks...")
channel.start_consuming()
