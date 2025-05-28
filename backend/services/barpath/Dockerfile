FROM python:3.10-slim

RUN apt-get update && \
    apt-get install -y ffmpeg libsm6 libxext6 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY processing ./processing
COPY rabbit ./rabbit
COPY main.py .
COPY config.py .
COPY runs/detect/train6/weights/best.pt runs/detect/train6/weights/best.pt

CMD ["python", "main.py"]
