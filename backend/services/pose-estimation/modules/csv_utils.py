import csv
import os

COCO_KEYPOINTS = [
    "Nose", "Left Eye", "Right Eye", "Left Ear", "Right Ear",
    "Left Shoulder", "Right Shoulder",
    "Left Elbow", "Right Elbow",
    "Left Wrist", "Right Wrist",
    "Left Hip", "Right Hip",
    "Left Knee", "Right Knee",
    "Left Ankle", "Right Ankle"
]
def get_csv_writer(smoothed, csv_path):
    csv_path = os.path.join(csv_path, "keypoints_output.csv")
    body_indices = list(range(5, 17))  

    with open(csv_path, mode="w", newline="") as f:
        writer = csv.writer(f)

        header = ["frame_idx"]
        for idx in body_indices:
            header.append(f"{COCO_KEYPOINTS[idx]}_x")
            header.append(f"{COCO_KEYPOINTS[idx]}_y")
        writer.writerow(header)

        for frame_idx, kps in enumerate(smoothed):
            row = [frame_idx]
            for idx in body_indices:
                x, y = kps[idx]
                row.extend([float(x), float(y)])
            writer.writerow(row)

    print(f"Keypoints saved to {csv_path}")