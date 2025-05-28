from minio import Minio

MINIO_CLIENT = Minio(
    "minio:9000",
    access_key="admin",
    secret_key="admin123",
    secure=False
)

MODEL_PATH = "runs/detect/train6/weights/best.pt"
TASK_QUEUE = "bar_path_queue"
REPLY_QUEUE = "bar_path_results_queue"
MAX_CONCURRENT_TASKS = 2
