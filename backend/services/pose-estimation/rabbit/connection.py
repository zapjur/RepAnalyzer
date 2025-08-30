import aio_pika
import asyncio
import logging

logger = logging.getLogger(__name__)

async def connect_with_retries(url, max_retries=10, delay=3):
    for attempt in range(1, max_retries + 1):
        try:
            logger.info(f"Attempt {attempt}: Connecting to RabbitMQ...")
            connection = await aio_pika.connect_robust(url)
            logger.info("Connected to RabbitMQ.")
            return connection
        except Exception as e:
            logger.warning(f"Connection attempt {attempt} failed: {e}")
            await asyncio.sleep(delay)

    raise RuntimeError("Failed to connect to RabbitMQ after multiple attempts.")