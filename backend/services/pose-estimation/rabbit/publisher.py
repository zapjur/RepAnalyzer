import os
import json
import asyncio
import aio_pika
from aio_pika import Message, DeliveryMode

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq/")
_conn = None
_chan = None

async def _get_channel():
    global _conn, _chan
    if _conn is None or _conn.is_closed:
        _conn = await aio_pika.connect_robust(RABBITMQ_URL, heartbeat=600, timeout=30)
    if _chan is None or _chan.is_closed:
        _chan = await _conn.channel(publisher_confirms=True)
    return _chan

async def publish_result(reply_queue: str, result: dict, retries: int = 5):
    body = json.dumps(result).encode()
    delay = 0.5
    for _ in range(retries):
        try:
            ch = await _get_channel()
            await ch.declare_queue(reply_queue, durable=True)
            await ch.default_exchange.publish(
                Message(
                    body=body,
                    content_type="application/json",
                    delivery_mode=DeliveryMode.PERSISTENT,
                ),
                routing_key=reply_queue,
            )
            return
        except asyncio.CancelledError:
            raise
        except Exception:
            global _conn, _chan
            _chan = None
            _conn = None
            await asyncio.sleep(delay)
            delay = min(delay * 2, 5.0)

    print("[publish_result] failed after retries")
