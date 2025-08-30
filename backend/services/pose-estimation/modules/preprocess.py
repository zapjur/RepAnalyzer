import torch
import numpy as np

def detect_persons(image, person_model, person_processor, device, height, width):
    inputs = person_processor(images=image, return_tensors="pt").to(device)
    with torch.no_grad():
        outputs = person_model(**inputs)
    results = person_processor.post_process_object_detection(
        outputs, target_sizes=torch.tensor([(height, width)]), threshold=0.5
    )
    result = results[0]
    person_boxes_xyxy = result["boxes"][result["labels"] == 0]

    person_boxes_xyxy = person_boxes_xyxy.cpu().numpy()
    person_boxes = person_boxes_xyxy.copy()
    person_boxes[:, 2] = person_boxes[:, 2] - person_boxes[:, 0]
    person_boxes[:, 3] = person_boxes[:, 3] - person_boxes[:, 1]
    return person_boxes, person_boxes_xyxy

def get_poses(image, persons, pose_model, pose_processor, device):
    pose_inputs = pose_processor(image, boxes=[persons], return_tensors="pt").to(device)
    if pose_model.config.backbone_config.num_experts > 1:
        dataset_index = torch.tensor([0] * len(pose_inputs["pixel_values"]))
        dataset_index = dataset_index.to(pose_inputs["pixel_values"].device)
        pose_inputs["dataset_index"] = dataset_index
    with torch.no_grad():
        pose_outputs = pose_model(**pose_inputs)
    pose_results = pose_processor.post_process_pose_estimation(pose_outputs, boxes=[persons])
    persons = pose_results[0]
    return persons

def get_person(persons):
    areas = []
    for p in persons:
        x1, y1, x2, y2 = p["bbox"]
        area = (x2 - x1) * (y2 - y1)
        areas.append(area)

    best_idx = int(np.argmax(areas))
    person = persons[best_idx]
    return person