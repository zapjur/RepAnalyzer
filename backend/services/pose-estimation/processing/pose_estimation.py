import os, io, csv, subprocess
import numpy as np
import cv2
from config import MINIO_CLIENT

from modules.model_utils import load_person_model, load_pose_model
from modules.preprocess import detect_persons, get_poses, get_person
from modules.drawing_utils import draw_landmarks
from modules.postprocess import kf_bank, kf_initialized
from modules.csv_utils import get_csv_writer

DEVICE = "cpu"
PERSON_MODEL, PERSON_PROC = load_person_model()
POSE_MODEL, POSE_PROC = load_pose_model()

def _pose_keys_from_original(original_key: str):
    mp4_key = original_key.replace("/original/", "/pose/", 1)
    csv_key = os.path.splitext(mp4_key)[0] + ".csv"
    return mp4_key, csv_key

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

def run_pipeline_sync(data: dict) -> dict:
    video_id    = data["video_id"]
    bucket      = data["bucket"]
    object_key  = data["object_key"]
    auth0id    = data["auth0_id"]
    exercise_name = data["exercise_name"]
    reply_queue = data["reply_queue"]


    tmp_in  = f"/tmp/pose_{video_id}.mp4"
    tmp_out = f"/tmp/pose_{video_id}_raw.mp4"
    tmp_h264= f"/tmp/pose_{video_id}.mp4"
    try:
        obj = MINIO_CLIENT.get_object(bucket, object_key)
        try:
            with open(tmp_in, "wb") as f:
                for chunk in obj.stream(32*1024):
                    f.write(chunk)
        finally:
            obj.close(); obj.release_conn()

        cap = cv2.VideoCapture(tmp_in)
        fps = cap.get(cv2.CAP_PROP_FPS) or 25.0
        W   = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH) or 1280)
        H   = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT) or 720)
        fourcc = cv2.VideoWriter_fourcc(*"mp4v")
        out = cv2.VideoWriter(tmp_out, fourcc, fps, (W,H))

        keypoints = []
        conf_thresh = 0.30
        frame_idx = 0

        import modules.postprocess as _pp
        _pp.kf_initialized = False

        while True:
            ok, frame = cap.read()
            if not ok: break
            frame_idx += 1

            image = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            person_boxes, person_boxes_xyxy = detect_persons(image, PERSON_MODEL, PERSON_PROC, DEVICE, H, W)
            if len(person_boxes_xyxy) == 0:
                out.write(frame); continue

            persons = get_poses(image, person_boxes, POSE_MODEL, POSE_PROC, DEVICE)
            if len(persons) == 0:
                out.write(frame); continue

            person = get_person(persons)
            kps    = person["keypoints"].cpu().numpy()
            scores = person["scores"].cpu().numpy()

            dt = 1.0 / max(fps, 1.0)
            for kf in kf_bank:
                A = np.array([[1,0,dt,0],
                              [0,1,0,dt],
                              [0,0,1, 0],
                              [0,0,0, 1]], dtype=np.float32)
                kf.transitionMatrix = A

            if not _pp.kf_initialized:
                for i in range(17):
                    x, y = kps[i]
                    kf_bank[i].statePost = np.array([[x],[y],[0.0],[0.0]], dtype=np.float32)
                _pp.kf_initialized = True

            smoothed = np.zeros_like(kps, dtype=np.float32)
            for i in range(17):
                pred = kf_bank[i].predict()
                if scores[i] >= conf_thresh and np.isfinite(kps[i]).all():
                    meas = np.array([[np.float32(kps[i,0])],[np.float32(kps[i,1])]], dtype=np.float32)
                    est = kf_bank[i].correct(meas)
                    smoothed[i,0], smoothed[i,1] = est[0,0], est[1,0]
                else:
                    smoothed[i,0], smoothed[i,1] = pred[0,0], pred[1,0]

            frame = draw_landmarks(smoothed, scores, conf_thresh, frame)
            keypoints.append(smoothed)
            out.write(frame)

        cap.release(); out.release()

        try:
            _h264_faststart(tmp_out, tmp_h264)
            upload_video = tmp_h264
        except Exception:
            upload_video = tmp_out

        out_mp4_key, out_csv_key = _pose_keys_from_original(object_key)

        MINIO_CLIENT.fput_object(bucket, out_mp4_key, upload_video, content_type="video/mp4")

        tmp_csv = f"/tmp/pose_{video_id}.csv"
        get_csv_writer(np.array(keypoints), os.path.dirname(tmp_csv))
        if not os.path.exists(tmp_csv):
            with open(tmp_csv, "w", newline="") as f:
                w = csv.writer(f)
                w.writerow(["frame"] + [f"k{i}_x" for i in range(17)] + [f"k{i}_y" for i in range(17)])
                for fi, sm in enumerate(keypoints):
                    row = [fi] + [float(sm[i,0]) for i in range(17)] + [float(sm[i,1]) for i in range(17)]
                    w.writerow(row)

        MINIO_CLIENT.fput_object(bucket, out_csv_key, tmp_csv, content_type="text/csv")

        return {
            "video_id": video_id,
            "status": "success",
            "bucket": bucket,
            "object_key": out_mp4_key,
            "auth0_id": auth0id,
            "exercise_name": exercise_name,
            "reply_queue": reply_queue
        }
    except Exception as e:
        return {"video_id": data.get("video_id"), "status": "error", "message": str(e)}
    finally:
        for p in (tmp_in, tmp_out, tmp_h264, locals().get("tmp_csv", None)):
            try:
                if p and os.path.exists(p): os.remove(p)
            except Exception:
                pass
