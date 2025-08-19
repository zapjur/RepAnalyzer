import os
from minio import Minio

MINIO_ENDPOINT   = os.getenv("MINIO_ENDPOINT")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY")
MINIO_USE_SSL    = os.getenv("MINIO_USE_SSL", "false").lower() == "true"
DEFAULT_METERS_PER_PIXEL = 0.0025

MINIO_CLIENT = Minio(
    MINIO_ENDPOINT,
    access_key=MINIO_ACCESS_KEY,
    secret_key=MINIO_SECRET_KEY,
    secure=MINIO_USE_SSL,
)

MODEL_PATH = "runs/detect/train6/weights/best.pt"
TASK_QUEUE = "bar_path_queue"
REPLY_QUEUE = "bar_path_results_queue"
MAX_CONCURRENT_TASKS = 2
