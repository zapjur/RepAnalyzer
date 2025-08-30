import cv2
import numpy as np

def make_kalman():
    kf = cv2.KalmanFilter(4, 2)
    kf.measurementMatrix = np.array([[1,0,0,0],
                                     [0,1,0,0]], dtype=np.float32)
    kf.processNoiseCov   = np.eye(4, dtype=np.float32) * 1e-2
    kf.measurementNoiseCov = np.eye(2, dtype=np.float32) * 5e-2
    kf.errorCovPost      = np.eye(4, dtype=np.float32) * 1.0
    
    return kf

kf_bank = [make_kalman() for _ in range(17)]
kf_initialized = False


def interpolate_keypoints(all_keypoints):
    all_keypoints = np.array(all_keypoints, dtype=np.float32)

    for j in range(all_keypoints.shape[1]):
        x = all_keypoints[:, j, 0]
        y = all_keypoints[:, j, 1]

        valid = ~np.isnan(x)
        frames = np.arange(len(x))

        if valid.sum() < 2:
            continue

        x_interp = np.interp(frames, frames[valid], x[valid])
        y_interp = np.interp(frames, frames[valid], y[valid])

        all_keypoints[:, j, 0] = x_interp
        all_keypoints[:, j, 1] = y_interp

    return all_keypoints