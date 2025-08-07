import os
import cv2
from ultralytics import YOLO
from config import MODEL_PATH, MINIO_CLIENT
from rabbit.publisher import publish_result

MODEL = YOLO(MODEL_PATH)

def _stem_from_object_key(object_key: str) -> str:
    base = os.path.basename(object_key)
    return os.path.splitext(base)[0]

async def process_task(data):
    video_id   = data["video_id"]
    bucket     = data["bucket"]
    object_key = data["object_key"]
    auth0_id   = data["auth0_id"]
    exercise   = data["exercise_name"]
    reply_queue= data["reply_queue"]

    tmp_in  = f"/tmp/{video_id}.mp4"
    tmp_out = f"/tmp/{video_id}_out.mp4"

    try:
        obj = MINIO_CLIENT.get_object(bucket, object_key)
        try:
            with open(tmp_in, "wb") as f:
                for chunk in obj.stream(32 * 1024):
                    f.write(chunk)
        finally:
            obj.close()
            obj.release_conn()

        cap = cv2.VideoCapture(tmp_in)
        fps = cap.get(cv2.CAP_PROP_FPS) or 25.0
        width  = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH) or 1280)
        height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT) or 720)
        fourcc = cv2.VideoWriter_fourcc(*"mp4v")
        out = cv2.VideoWriter(tmp_out, fourcc, fps, (width, height))

        traj = []
        while cap.isOpened():
            ret, frame = cap.read()
            if not ret:
                break

            results = MODEL.predict(source=frame, conf=0.5, iou=0.4, classes=[0], verbose=False)
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

        stem = _stem_from_object_key(object_key)
        output_key = f"{auth0_id}/{exercise}/barpath/{stem}.mp4"

        MINIO_CLIENT.fput_object(
            bucket,
            output_key,
            tmp_out,
            content_type="video/mp4",
        )

        result = {
            "video_id": video_id,
            "status": "success",
            "bucket": bucket,
            "object_key": output_key,
        }

    except Exception as e:
        result = {
            "video_id": video_id,
            "status": "error",
            "message": str(e),
        }
    finally:
        try:
            if os.path.exists(tmp_in):
                os.remove(tmp_in)
            if os.path.exists(tmp_out):
                os.remove(tmp_out)
        except Exception:
            pass

    await publish_result(reply_queue, result)
