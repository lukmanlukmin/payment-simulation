#!/bin/bash
set -e

cleanup() {
  echo "🧹 Stopping services..."
#   docker compose down -v
}
trap cleanup EXIT

echo "🚀 Starting services..."
docker compose up -d doitpay-postgresql doitpay-kafka doitpay-zookeeper doitpay-api doitpay-worker

# Tunggu API benar-benar ready (cek port)
echo "⏳ Waiting for doitpay-api to be ready..."
until docker compose exec -T doitpay-api sh -c "netstat -tuln | grep -q ':80 '"; do
  sleep 2
done
echo "✅ doitpay-api is ready"

# Tunggu worker benar-benar jalan (cek proses run-worker)
echo "⏳ Waiting for doitpay-worker to be running..."
until docker compose exec -T doitpay-worker sh -c "pgrep -f 'run-worker' > /dev/null"; do
  sleep 2
done
echo "✅ doitpay-worker is running"

echo "✅ All services are up. Running k6 load test..."
docker compose run --rm doitpay-k6 run /scripts/transaction.js
