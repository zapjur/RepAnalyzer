import asyncio
import aio_pika
import json
import logging
from config import TASK_QUEUE, MAX_CONCURRENT_TASKS
from processing.pose_estimation import run_pipeline_sync
from rabbit.publisher import publish_result

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

semaphore = asyncio.Semaphore(MAX_CONCURRENT_TASKS)

async def handle_message(message: aio_pika.IncomingMessage):
    async with message.process(requeue=False):
        async with semaphore:
            payload = json.loads(message.body)
            video_id = payload.get("video_id")
            logger.info(f"Received task: {video_id}")
            try:
                result = await asyncio.to_thread(run_pipeline_sync, payload)
                await publish_result(payload["reply_queue"], result)
                logger.info(f"Finished task: {video_id}")
            except Exception as e:
                logger.exception(f"Error processing message {video_id}: {e}")
                err = {"video_id": video_id, "status": "error", "message": str(e)}
                try:
                    await publish_result(payload.get("reply_queue", "pose_results_queue"), err)
                except Exception:
                    logger.exception("Failed to publish error result")

async def start_consumer(queue_name: str):
    connection = await aio_pika.connect_robust("amqp://guest:guest@rabbitmq/", heartbeat=600, timeout=30)
    channel = await connection.channel()
    await channel.set_qos(prefetch_count=1)
    queue = await channel.declare_queue(queue_name, durable=True)
    await queue.consume(handle_message)
    logger.info(f" [*] Waiting for messages in {queue_name}.")
    await asyncio.Future()
