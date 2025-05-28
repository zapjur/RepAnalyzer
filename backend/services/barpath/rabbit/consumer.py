import asyncio
import aio_pika
import json
import logging
from config import TASK_QUEUE, MAX_CONCURRENT_TASKS
from processing.barpath import process_task

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

semaphore = asyncio.Semaphore(MAX_CONCURRENT_TASKS)

async def handle_message(message: aio_pika.IncomingMessage):
    async with message.process():
        async with semaphore:
            try:
                payload = json.loads(message.body)
                logger.info(f"Received task: {payload.get('video_id')}")
                await process_task(payload)
                logger.info(f"Finished task: {payload.get('video_id')}")
            except Exception as e:
                logger.error(f"Error processing message: {e}")

async def start_consumer(queue_name: str):
    connection = await aio_pika.connect_robust("amqp://guest:guest@rabbitmq/")
    channel = await connection.channel()
    queue = await channel.declare_queue(queue_name, durable=True)
    await queue.consume(handle_message)
    logger.info(f" [*] Waiting for messages in {queue_name}.")
    await asyncio.Future()
