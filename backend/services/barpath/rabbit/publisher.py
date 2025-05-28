import aio_pika
import json

async def publish_result(reply_queue: str, result: dict):
    connection = await aio_pika.connect_robust("amqp://guest:guest@rabbitmq/")
    channel = await connection.channel()
    await channel.default_exchange.publish(
        aio_pika.Message(body=json.dumps(result).encode()),
        routing_key=reply_queue
    )
    await connection.close()
