from ultralytics import YOLO
import cv2
import os

VIDEO_PATH = "videos_to_test_mp4/1.mp4" 
MODEL_PATH = "runs/detect/train6/weights/best.pt"
OUTPUT_PATH = "output_barpath.mp4"

model = YOLO(MODEL_PATH)
cap = cv2.VideoCapture(VIDEO_PATH)
traj = []

fps = cap.get(cv2.CAP_PROP_FPS)
width  = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
fourcc = cv2.VideoWriter_fourcc(*'mp4v')
out = cv2.VideoWriter(OUTPUT_PATH, fourcc, fps, (width, height))

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
print(f"Saved to {OUTPUT_PATH}")
