echo "Testing rate limiter..."
echo "Должно занять ~2 секунды для 200 запросов при лимите 100 RPS"
echo ""

time (
  for i in {1..200}; do
    grpcurl -plaintext \
      -d '{"username": "denis", "password": "12345678"}' \
      localhost:50051 \
      api.Auth/Login > /dev/null 2>&1 &
  done
  wait
)