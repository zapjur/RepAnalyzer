import asyncio
from config import TASK_QUEUE
from rabbit.consumer import start_consumer
from rabbit.connection import connect_with_retries

RABBITMQ_URL = "amqp://guest:guest@rabbitmq/"

async def main():
    await connect_with_retries(RABBITMQ_URL)
    await start_consumer(queue_name=TASK_QUEUE)

if __name__ == "__main__":
    asyncio.run(main())