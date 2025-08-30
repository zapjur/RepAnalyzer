# models.py
import torch
from transformers import RTDetrForObjectDetection, VitPoseForPoseEstimation, AutoProcessor

DEVICE = "cuda" if torch.cuda.is_available() else "cpu"

_person_model = None
_person_processor = None
_pose_model = None
_pose_processor = None

def load_person_model():
    global _person_model, _person_processor
    if _person_model is None or _person_processor is None:
        _person_processor = AutoProcessor.from_pretrained("PekingU/rtdetr_r50vd_coco_o365")
        _person_model = RTDetrForObjectDetection.from_pretrained(
            "PekingU/rtdetr_r50vd_coco_o365"
        ).to(DEVICE)
    return _person_model, _person_processor

def load_pose_model():
    global _pose_model, _pose_processor
    if _pose_model is None or _pose_processor is None:
        _pose_processor = AutoProcessor.from_pretrained("usyd-community/vitpose-base-simple")
        _pose_model = VitPoseForPoseEstimation.from_pretrained(
            "usyd-community/vitpose-base-simple"
        ).to(DEVICE)
    return _pose_model, _pose_processor
