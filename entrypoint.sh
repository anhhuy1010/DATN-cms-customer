#!/bin/sh

if [ -z "$PORT" ]; then
  echo "PORT environment variable not set. Default to 8080"
  export PORT=8080
fi

# Thay thế port trong file envoy.yaml bằng port từ Heroku
sed -i "s/port_value: 8080/port_value: $PORT/" /app/envoy.yaml

# Chạy Go app (gRPC server) trên port 9003 (cấu hình trong app của bạn phải đúng)
./main &

# Chạy Envoy proxy (lắng nghe port $PORT)
envoy -c /app/envoy.yaml
