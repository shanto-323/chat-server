sleep 3

# retry until the exchange is successfully declared
until rabbitmqadmin declare exchange \
  --vhost=shanto323 \
  --username=shanto \
  --password=123456 \
  name=message.service \
  type=topic \
  durable=true; do
    echo "Waiting for RabbitMQ..."
    sleep 2
done

echo "Exchange message.service created!"