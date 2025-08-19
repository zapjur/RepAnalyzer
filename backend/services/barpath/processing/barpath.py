import os, io, csv, subprocess
import numpy as np
import cv2
from ultralytics import YOLO
from config import MODEL_PATH, MINIO_CLIENT, DEFAULT_METERS_PER_PIXEL
from rabbit.publisher import publish_result

MODEL = YOLO(MODEL_PATH)
ASSUMED_PLATE_DIAM_M = 0.45

def _barpath_keys_from_original(original_key: str):
    mp4_key = original_key.replace("/original/", "/barpath/", 1)
    csv_key = os.path.splitext(mp4_key)[0] + ".csv"
    return mp4_key, csv_key

def _moving_average(a, k=9):
    k = max(3, k | 1)
    pad = k // 2
    a_pad = np.pad(a, (pad, pad), mode="edge")
    ker = np.ones(k) / k
    return np.convolve(a_pad, ker, mode="valid")

def _mpp_from_plate_diams(diams_px):
    if not diams_px:
        return None
    med = float(np.median(np.asarray(diams_px, dtype=float)))
    if med <= 0:
        return None
    return ASSUMED_PLATE_DIAM_M / med

def _h264_faststart(src, dst):
    subprocess.run([
        "ffmpeg","-y","-i",src,
        "-c:v","libx264","-preset","veryfast","-crf","23",
        "-pix_fmt","yuv420p",
        "-profile:v","baseline","-level","3.0",
        "-movflags","+faststart",
        "-c:a","aac","-b:a","128k",
        dst
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

async def process_task(data):
    video_id = data["video_id"]
    bucket = data["bucket"]
    object_key = data["object_key"]
    reply_queue = data["reply_queue"]

    tmp_in = f"/tmp/{video_id}.mp4"
    tmp_out = f"/tmp/{video_id}_out.mp4"

    try:
        obj = MINIO_CLIENT.get_object(bucket, object_key)
        try:
            with open(tmp_in, "wb") as f:
                for chunk in obj.stream(32 * 1024):
                    f.write(chunk)
        finally:
            obj.close(); obj.release_conn()

        cap = cv2.VideoCapture(tmp_in)
        fps = cap.get(cv2.CAP_PROP_FPS) or 25.0
        width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH) or 1280)
        height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT) or 720)
        fourcc = cv2.VideoWriter_fourcc(*"mp4v")
        out = cv2.VideoWriter(tmp_out, fourcc, fps, (width, height))

        traj = []
        frames, xs, ys = [], [], []
        plate_diams = []
        i = 0

        while cap.isOpened():
            ret, frame = cap.read()
            if not ret:
                break

            results = MODEL.predict(source=frame, conf=0.5, iou=0.4, verbose=False)
            boxes = results[0].boxes

            if boxes is not None and len(boxes) > 0:
                xyxy = boxes.xyxy.cpu().numpy()
                areas = (xyxy[:,2]-xyxy[:,0]) * (xyxy[:,3]-xyxy[:,1])
                j = int(np.argmax(areas))
                x1, y1, x2, y2 = map(int, xyxy[j])
                cx = (x1 + x2) // 2
                cy = (y1 + y2) // 2
                traj.append((cx, cy))
                xs.append(cx); ys.append(cy); frames.append(i)

                diam = float(max(x2 - x1, y2 - y1))
                if 60 <= diam <= 2000:
                    plate_diams.append(diam)

                cv2.circle(frame, (cx, cy), 4, (0, 255, 0), -1)

            for k in range(1, len(traj)):
                cv2.line(frame, traj[k - 1], traj[k], (255, 0, 0), 2)

            out.write(frame)
            i += 1

        cap.release()
        out.release()
        cv2.destroyAllWindows()

        tmp_out_h264 = f"/tmp/{video_id}_out_h264.mp4"
        try:
            _h264_faststart(tmp_out, tmp_out_h264)
            upload_path = tmp_out_h264
        except Exception:
            upload_path = tmp_out

        out_mp4_key, out_csv_key = _barpath_keys_from_original(object_key)
        MINIO_CLIENT.fput_object(bucket, out_mp4_key, upload_path, content_type="video/mp4")

        mpp = _mpp_from_plate_diams(plate_diams)
        if mpp is None or not (0.0005 <= mpp <= 0.01):
            mpp = float(DEFAULT_METERS_PER_PIXEL)

        if len(frames) >= 2:
            t = np.array(frames, dtype=float) / float(fps)
            x = np.array(xs, dtype=float)
            y = np.array(ys, dtype=float)
            vy_px = -np.gradient(y, t)
            vy_px_s = _moving_average(vy_px, k=9)
            vy_m_s = vy_px * mpp
            vy_smooth_m_s = vy_px_s * mpp

            buf = io.StringIO()
            w = csv.writer(buf)
            w.writerow(["frame","t","x_px","y_px","vy_m_s","vy_smooth_m_s","meters_per_pixel"])
            for j in range(len(t)):
                w.writerow([int(frames[j]), float(t[j]), float(x[j]), float(y[j]), float(vy_m_s[j]), float(vy_smooth_m_s[j]), float(mpp)])
            data_bytes = buf.getvalue().encode("utf-8")

            MINIO_CLIENT.put_object(
                bucket,
                out_csv_key,
                io.BytesIO(data_bytes),
                len(data_bytes),
                content_type="text/csv",
                metadata={"x-amz-meta-velocity":"m/s","x-amz-meta-meters-per-pixel":str(mpp)},
            )

        result = {"video_id": video_id, "status": "success", "bucket": bucket, "object_key": out_mp4_key}
    except Exception as e:
        result = {"video_id": video_id, "status": "error", "message": str(e)}
    finally:
        try:
            if os.path.exists(tmp_in): os.remove(tmp_in)
            if os.path.exists(tmp_out): os.remove(tmp_out)
            tmp_out_h264 = f"/tmp/{video_id}_out_h264.mp4"
            if os.path.exists(tmp_out_h264): os.remove(tmp_out_h264)
        except Exception:
            pass

    await publish_result(reply_queue, result)
