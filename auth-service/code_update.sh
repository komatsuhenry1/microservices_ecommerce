set -euo pipefail

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap .
zip -j function.zip bootstrap >/dev/null

aws lambda update-function-code \
  --function-name auth-service \
  --zip-file fileb://function.zip \
  --endpoint-url http://localhost:4566

echo "✅ Code updated"