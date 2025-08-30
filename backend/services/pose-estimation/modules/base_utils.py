import os
import cv2

def get_writer_str(uploaded_file, TEMP_FOLDER, frame_width, frame_height, fps):
    output_path = os.path.join(TEMP_FOLDER, f"processed_{uploaded_file.name}")
    out = cv2.VideoWriter(
        output_path,
        cv2.VideoWriter_fourcc(*'mp4v'),
        fps,
        (frame_width, frame_height)
    )
    return output_path, out


def get_writer(uploaded_file, TEMP_FOLDER, frame_width, frame_height, fps):
    os.makedirs(TEMP_FOLDER, exist_ok=True)

    output_path = os.path.join(TEMP_FOLDER, f"processed_{uploaded_file}")

    out = cv2.VideoWriter(
        output_path,
        cv2.VideoWriter_fourcc(*'mp4v'),
        fps,
        (frame_width, frame_height)
    )

    return output_path, out


def initialize_video(video_path):
    cap = cv2.VideoCapture(video_path)
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    fps = int(cap.get(cv2.CAP_PROP_FPS))
    frame_count = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    frame_count = max(frame_count, 1)
    return cap, width, height, fps, frame_count

