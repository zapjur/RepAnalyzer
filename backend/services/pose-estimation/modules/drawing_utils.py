import cv2
import numpy as np

def draw_keypoints(frame, keypoints, confidence_threshold, input_size):
    y, x, c = frame.shape
    input_h, input_w = input_size
    shaped = np.squeeze(np.multiply(keypoints, [input_h, input_w, 1]))

    scale_y = y / input_h
    scale_x = x / input_w

    for ky, kx, kp_conf in shaped:
        if kp_conf > confidence_threshold:
            draw_x = int(kx * scale_x)
            draw_y = int(ky * scale_y)
            cv2.circle(frame, (draw_x, draw_y), 2, (0,0,255), -1)


def draw_connections(frame, keypoints, edges, confidence_threshold, input_size):
    y, x, c = frame.shape
    input_h, input_w = input_size
    shaped = np.squeeze(np.multiply(keypoints, [input_h, input_w, 1]))

    scale_y = y / input_h
    scale_x = x / input_w

    for edge, color in edges.items():
        p1, p2 = edge
        y1, x1, c1 = shaped[p1]
        y2, x2, c2 = shaped[p2]

        if (c1 > confidence_threshold) and (c2 > confidence_threshold):
            draw_x1 = int(x1 * scale_x)
            draw_y1 = int(y1 * scale_y)
            draw_x2 = int(x2 * scale_x)
            draw_y2 = int(y2 * scale_y)
            cv2.line(frame, (draw_x1, draw_y1), (draw_x2, draw_y2), (255,255,255), 1)

def loop_through_people(frame, keypoints_with_scores, edges, confidence_threshold, input_size):
    for person in keypoints_with_scores:
        draw_connections(frame, person, edges, confidence_threshold, input_size)
        draw_keypoints(frame, person, confidence_threshold, input_size)

def draw_landmarks(smoothed, scores, conf_thresh, frame):
    skeleton = [
        (5,7), (7,9), (6,8), (8,10),
        (11,13), (13,15), (12,14), (14,16),
        (5,6), (11,12), (5,11), (6,12),
        (0,1), (1,2), (2,3), (3,4),
        (0,5), (0,6)
    ]

    face_indices = {0, 1, 2, 3, 4}
    for i, (x, y) in enumerate(smoothed):
        if i in face_indices:
            continue
        if scores[i] >= conf_thresh:
            cv2.circle(frame, (int(x), int(y)), 3, (0,0,255), -1)

    face_bones = {(0,1), (1,2), (2,3), (3,4)}
    for i, j in skeleton:
        if (i in face_indices and j in face_indices) or (i, j) in face_bones:
            continue
        if scores[i] >= conf_thresh and scores[j] >= conf_thresh:
            pt1 = (int(smoothed[i,0]), int(smoothed[i,1]))
            pt2 = (int(smoothed[j,0]), int(smoothed[j,1]))
            cv2.line(frame, pt1, pt2, (0,255,0), 2)
    return frame