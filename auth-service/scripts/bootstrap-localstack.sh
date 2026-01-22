#!/usr/bin/env bash
set -euo pipefail

# ===== Config =====
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-test}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-test}
export AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION:-us-east-1}

LS_ENDPOINT=${LS_ENDPOINT:-http://localhost:4566}
FUNC_NAME=${FUNC_NAME:-auth-service}
API_NAME=${API_NAME:-ecommerce-api}
STAGE=${STAGE:-dev}

echo "==> LocalStack endpoint: $LS_ENDPOINT"
echo "==> Region: $AWS_DEFAULT_REGION"

# ===== Build Lambda (Go custom runtime) =====
echo "==> Building lambda binary..."
GOOS=linux GOARCH=amd64 go build -o bootstrap .
chmod +x bootstrap
zip -j function.zip bootstrap >/dev/null

# ===== (Re)Create Lambda =====
echo "==> Recreating lambda: $FUNC_NAME"
aws --region "$AWS_DEFAULT_REGION" lambda delete-function \
  --function-name "$FUNC_NAME" \
  --endpoint-url="$LS_ENDPOINT" >/dev/null 2>&1 || true

aws --region "$AWS_DEFAULT_REGION" lambda create-function \
  --function-name "$FUNC_NAME" \
  --runtime provided.al2 \
  --handler bootstrap \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --zip-file fileb://function.zip \
  --timeout 10 \
  --memory-size 128 \
  --endpoint-url="$LS_ENDPOINT" >/dev/null

# (Opcional) Env vars pro DB (AJUSTE depois conforme seu ConnectRDS)
echo "==> Setting lambda env vars (DB_*)..."
aws --region "$AWS_DEFAULT_REGION" lambda update-function-configuration \
  --function-name "$FUNC_NAME" \
  --environment "Variables={DB_HOST=postgres,DB_PORT=5432,DB_USER=admin,DB_PASSWORD=secret,DB_NAME=auth_db}" \
  --endpoint-url="$LS_ENDPOINT" >/dev/null || true

# ===== Create API Gateway =====
echo "==> Creating REST API: $API_NAME"
API_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-rest-api \
  --name "$API_NAME" \
  --endpoint-url="$LS_ENDPOINT" \
  --query 'id' --output text)

echo "==> API_ID: $API_ID"

ROOT_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway get-resources \
  --rest-api-id "$API_ID" \
  --endpoint-url="$LS_ENDPOINT" \
  --query 'items[0].id' --output text)

echo "==> ROOT_ID: $ROOT_ID"

REGISTER_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-resource \
  --rest-api-id "$API_ID" \
  --parent-id "$ROOT_ID" \
  --path-part register \
  --endpoint-url="$LS_ENDPOINT" \
  --query 'id' --output text)

echo "==> REGISTER_ID: $REGISTER_ID"

echo "==> Creating POST method..."
aws --region "$AWS_DEFAULT_REGION" apigateway put-method \
  --rest-api-id "$API_ID" \
  --resource-id "$REGISTER_ID" \
  --http-method POST \
  --authorization-type NONE \
  --endpoint-url="$LS_ENDPOINT" >/dev/null

echo "==> Integrating POST -> Lambda (AWS_PROXY)..."
aws --region "$AWS_DEFAULT_REGION" apigateway put-integration \
  --rest-api-id "$API_ID" \
  --resource-id "$REGISTER_ID" \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri "arn:aws:apigateway:$AWS_DEFAULT_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_DEFAULT_REGION:000000000000:function:$FUNC_NAME/invocations" \
  --endpoint-url="$LS_ENDPOINT" >/dev/null

echo "==> Deploying stage: $STAGE"
aws --region "$AWS_DEFAULT_REGION" apigateway create-deployment \
  --rest-api-id "$API_ID" \
  --stage-name "$STAGE" \
  --endpoint-url="$LS_ENDPOINT" >/dev/null

echo ""
echo "✅ Done!"
echo "API URL:"
echo "  $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/register"
echo ""
echo "Test:"
echo "  curl -X POST \\"
echo "    $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/register \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"name\":\"Joao\",\"email\":\"joao@email.com\",\"password\":\"123456\"}'"
