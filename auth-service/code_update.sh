#!/usr/bin/env bash
set -euo pipefail

# Config (pode sobrescrever via env)
LS_ENDPOINT=${LS_ENDPOINT:-http://localhost:4566}
FUNC_NAME=${FUNC_NAME:-auth-service}
REGION=${AWS_DEFAULT_REGION:-us-east-1}

echo "==> Updating Lambda code"
echo "==> Function: $FUNC_NAME"
echo "==> Endpoint: $LS_ENDPOINT"
echo "==> Region:   $REGION"

# Build Go binary for Lambda custom runtime
GOWORK=off GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap .
chmod +x bootstrap
zip -j function.zip bootstrap >/dev/null

# Update only the code
aws --region "$REGION" lambda update-function-code \
  --function-name "$FUNC_NAME" \
  --zip-file fileb://function.zip \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

echo "✅ Code updated for '$FUNC_NAME'"
