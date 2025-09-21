# -*- coding: utf-8 -*-
import os, io, csv, json, subprocess
import numpy as np
import cv2
from config import MINIO_CLIENT

from modules.model_utils import load_person_model, load_pose_model
from modules.preprocess import detect_persons, get_poses, get_person
from modules.drawing_utils import draw_landmarks
from modules.postprocess import kf_bank
from modules.csv_utils import COCO_KEYPOINTS  # użyjemy nazw COCO

DEVICE = "cpu"
PERSON_MODEL, PERSON_PROC = load_person_model()
POSE_MODEL, POSE_PROC = load_pose_model()

BODY_IDX = list(range(5, 17))  # barki -> kostki

def _pose_keys_from_original(original_key: str):
    mp4_key = original_key.replace("/original/", "/pose/", 1)
    csv_key = os.path.splitext(mp4_key)[0] + ".csv"
    meta_key = os.path.splitext(mp4_key)[0] + "_meta.json"
    return mp4_key, csv_key, meta_key

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

def _write_pose_csv(csv_path: str, fps: float, keypoints_list: list):
    os.makedirs(os.path.dirname(csv_path), exist_ok=True)
    with open(csv_path, "w", newline="") as f:
        w = csv.writer(f)
        header = ["frame_idx", "t"]
        for idx in BODY_IDX:
            header += [f"{COCO_KEYPOINTS[idx]}_x", f"{COCO_KEYPOINTS[idx]}_y"]
        w.writerow(header)
        for fi, sm in enumerate(keypoints_list):
            t = fi / max(fps, 1.0)
            row = [fi, float(t)]
            for idx in BODY_IDX:
                x, y = float(sm[idx,0]), float(sm[idx,1])
                row += [x, y]
            w.writerow(row)

def run_pipeline_sync(data: dict) -> dict:
    video_id    = data["video_id"]
    bucket      = data["bucket"]
    object_key  = data["object_key"]
    auth0id     = data["auth0_id"]
    exercise_name = data["exercise_name"]
    reply_queue = data["reply_queue"]

    tmp_in   = f"/tmp/pose_{video_id}.mp4"
    tmp_raw  = f"/tmp/pose_{video_id}_raw.mp4"
    tmp_h264 = f"/tmp/pose_{video_id}.mp4"
    tmp_csv  = f"/tmp/pose_{video_id}.csv"
    tmp_meta = f"/tmp/pose_{video_id}_meta.json"

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
        out = cv2.VideoWriter(tmp_raw, fourcc, fps, (W,H))

        keypoints = []
        conf_thresh = 0.30

        # (Re)initialize KF
        import modules.postprocess as _pp
        _pp.kf_initialized = False

        frame_idx = 0
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
            kps    = person["keypoints"].cpu().numpy()  # piksele w oryginalnej ramce
            scores = person["scores"].cpu().numpy()

            dt = 1.0 / max(fps, 1.0)
            for kf in kf_bank:
                kf.transitionMatrix = np.array([[1,0,dt,0],[0,1,0,dt],[0,0,1,0],[0,0,0,1]], dtype=np.float32)

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
            _h264_faststart(tmp_raw, tmp_h264)
            upload_video = tmp_h264
        except Exception:
            upload_video = tmp_raw

        out_mp4_key, out_csv_key, out_meta_key = _pose_keys_from_original(object_key)

        # 1 CSV – zawsze ten sam format
        _write_pose_csv(tmp_csv, fps, keypoints)

        # META JSON dla pewności układu
        meta = {
            "fps": float(fps),
            "width": int(W),
            "height": int(H),
            "origin": {"x": "right", "y": "down", "units": "pixel"},
            "dataset": "COCO-17",
            "indices": {COCO_KEYPOINTS[i]: i for i in range(17)},
            "body_indices_written": {COCO_KEYPOINTS[i]: i for i in BODY_IDX},
            "notes": "keypoints are in original video pixel coordinates"
        }
        with open(tmp_meta, "w", encoding="utf-8") as f:
            json.dump(meta, f, ensure_ascii=False)

        # Uploady
        MINIO_CLIENT.fput_object(bucket, out_mp4_key, upload_video, content_type="video/mp4")
        MINIO_CLIENT.fput_object(bucket, out_csv_key, tmp_csv, content_type="text/csv")
        MINIO_CLIENT.fput_object(bucket, out_meta_key, tmp_meta, content_type="application/json")

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
        for p in (tmp_in, tmp_raw, tmp_h264, tmp_csv, tmp_meta):
            try:
                if p and os.path.exists(p): os.remove(p)
            except Exception:
                pass
