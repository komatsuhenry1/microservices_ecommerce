#!/usr/bin/env bash
set -euo pipefail

# ===== Config =====
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-test}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-test}
export AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION:-us-east-1}

LS_ENDPOINT=${LS_ENDPOINT:-http://localhost:4566}

API_NAME=${API_NAME:-ecommerce-api}
STAGE=${STAGE:-dev}

FUNC_NAME=${FUNC_NAME:-auth-service}
DATABASE_URL=${DATABASE_URL:-postgres://admin:secret@db:5432/auth_db?sslmode=disable}

echo "==> LocalStack endpoint: $LS_ENDPOINT"
echo "==> Region: $AWS_DEFAULT_REGION"
echo "==> API name: $API_NAME"
echo "==> Stage: $STAGE"
echo "==> Lambda: $FUNC_NAME"

# ===== Helpers =====
wait_lambda_active() {
  local name="$1"
  local tries="${2:-30}"
  local sleep_s="${3:-1}"

  echo "==> Waiting lambda '$name' to become Active..."
  for i in $(seq 1 "$tries"); do
    state=$(aws --region "$AWS_DEFAULT_REGION" lambda get-function-configuration \
      --function-name "$name" \
      --endpoint-url "$LS_ENDPOINT" \
      --query 'State' --output text 2>/dev/null || echo "Unknown")

    echo "   - state: $state ($i/$tries)"

    if [[ "$state" == "Active" ]]; then
      return 0
    fi
    sleep "$sleep_s"
  done

  echo "❌ Lambda did not become Active in time."
  return 1
}

get_api_id_by_name() {
  aws --region "$AWS_DEFAULT_REGION" apigateway get-rest-apis \
    --endpoint-url "$LS_ENDPOINT" \
    --query "items[?name=='$API_NAME'].id | [0]" \
    --output text 2>/dev/null || true
}

get_resource_id_by_path() {
  local api_id="$1"
  local path="$2"
  aws --region "$AWS_DEFAULT_REGION" apigateway get-resources \
    --rest-api-id "$api_id" \
    --endpoint-url "$LS_ENDPOINT" \
    --query "items[?path=='$path'].id | [0]" \
    --output text 2>/dev/null || true
}

# ===== Build Lambda (Go custom runtime) =====
echo "==> Building lambda binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap .
chmod +x bootstrap
zip -j function.zip bootstrap >/dev/null

# ===== Create or Update Lambda =====
echo "==> Ensuring lambda exists..."
if aws --region "$AWS_DEFAULT_REGION" lambda get-function \
  --function-name "$FUNC_NAME" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null 2>&1; then

  echo "==> Lambda exists. Updating code..."
  aws --region "$AWS_DEFAULT_REGION" lambda update-function-code \
    --function-name "$FUNC_NAME" \
    --zip-file fileb://function.zip \
    --endpoint-url "$LS_ENDPOINT" >/dev/null

else
  echo "==> Lambda does not exist. Creating..."
  aws --region "$AWS_DEFAULT_REGION" lambda create-function \
    --function-name "$FUNC_NAME" \
    --runtime provided.al2 \
    --handler bootstrap \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --zip-file fileb://function.zip \
    --timeout 10 \
    --memory-size 256 \
    --endpoint-url "$LS_ENDPOINT" >/dev/null
fi

wait_lambda_active "$FUNC_NAME"

echo "==> Updating lambda env vars (DATABASE_URL)..."
aws --region "$AWS_DEFAULT_REGION" lambda update-function-configuration \
  --function-name "$FUNC_NAME" \
  --environment "Variables={DATABASE_URL=$DATABASE_URL}" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null || true

wait_lambda_active "$FUNC_NAME"

# ===== Create or Reuse API Gateway =====
API_ID=$(get_api_id_by_name)
if [[ -z "${API_ID// }" || "$API_ID" == "None" ]]; then
  echo "==> Creating REST API: $API_NAME"
  API_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-rest-api \
    --name "$API_NAME" \
    --query 'id' \
    --output text \
    --endpoint-url "$LS_ENDPOINT")
else
  echo "==> Reusing existing REST API: $API_ID ($API_NAME)"
fi

ROOT_ID=$(get_resource_id_by_path "$API_ID" "/")
echo "==> ROOT_ID: $ROOT_ID"

# ===== Ensure /auth resource =====
AUTH_ID=$(get_resource_id_by_path "$API_ID" "/auth")
if [[ -z "${AUTH_ID// }" || "$AUTH_ID" == "None" ]]; then
  echo "==> Creating resource '/auth'..."
  AUTH_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-resource \
    --rest-api-id "$API_ID" \
    --parent-id "$ROOT_ID" \
    --path-part auth \
    --query 'id' \
    --output text \
    --endpoint-url "$LS_ENDPOINT")
else
  echo "==> Resource '/auth' exists: $AUTH_ID"
fi

# ===== Ensure /auth/{proxy+} resource =====
AUTH_PROXY_ID=$(get_resource_id_by_path "$API_ID" "/auth/{proxy+}")
if [[ -z "${AUTH_PROXY_ID// }" || "$AUTH_PROXY_ID" == "None" ]]; then
  echo "==> Creating resource '/auth/{proxy+}'..."
  AUTH_PROXY_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-resource \
    --rest-api-id "$API_ID" \
    --parent-id "$AUTH_ID" \
    --path-part "{proxy+}" \
    --query 'id' \
    --output text \
    --endpoint-url "$LS_ENDPOINT")
else
  echo "==> Resource '/auth/{proxy+}' exists: $AUTH_PROXY_ID"
fi

LAMBDA_URI="arn:aws:apigateway:$AWS_DEFAULT_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_DEFAULT_REGION:000000000000:function:$FUNC_NAME/invocations"

# ===== Put ANY method + integration on /auth =====
echo "==> Configuring ANY on '/auth'..."
aws --region "$AWS_DEFAULT_REGION" apigateway put-method \
  --rest-api-id "$API_ID" \
  --resource-id "$AUTH_ID" \
  --http-method ANY \
  --authorization-type NONE \
  --no-api-key-required \
  --endpoint-url "$LS_ENDPOINT" >/dev/null || true

aws --region "$AWS_DEFAULT_REGION" apigateway put-integration \
  --rest-api-id "$API_ID" \
  --resource-id "$AUTH_ID" \
  --http-method ANY \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri "$LAMBDA_URI" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

# ===== Put ANY method + integration on /auth/{proxy+} =====
echo "==> Configuring ANY on '/auth/{proxy+}'..."
aws --region "$AWS_DEFAULT_REGION" apigateway put-method \
  --rest-api-id "$API_ID" \
  --resource-id "$AUTH_PROXY_ID" \
  --http-method ANY \
  --authorization-type NONE \
  --no-api-key-required \
  --request-parameters "method.request.path.proxy=true" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null || true

aws --region "$AWS_DEFAULT_REGION" apigateway put-integration \
  --rest-api-id "$API_ID" \
  --resource-id "$AUTH_PROXY_ID" \
  --http-method ANY \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri "$LAMBDA_URI" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

# ===== Lambda permissions (ignore if already exists) =====
echo "==> Adding lambda permissions (ignore if exists)..."
aws --region "$AWS_DEFAULT_REGION" lambda add-permission \
  --function-name "$FUNC_NAME" \
  --statement-id "apigw-auth-root" \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:$AWS_DEFAULT_REGION:000000000000:$API_ID/*/*/auth" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null 2>&1 || true

aws --region "$AWS_DEFAULT_REGION" lambda add-permission \
  --function-name "$FUNC_NAME" \
  --statement-id "apigw-auth-proxy" \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:$AWS_DEFAULT_REGION:000000000000:$API_ID/*/*/auth/*" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null 2>&1 || true

# ===== Deploy =====
echo "==> Deploying stage: $STAGE"
aws --region "$AWS_DEFAULT_REGION" apigateway create-deployment \
  --rest-api-id "$API_ID" \
  --stage-name "$STAGE" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

echo ""
echo "✅ Done!"
echo "Base URL (LocalStack):"
echo "  $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_"
echo ""
echo "Auth routes (examples):"
echo "  POST $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/auth/register"
echo "  POST $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/auth/login"
echo ""
echo "Test (register):"
echo "  curl -i -X POST \\"
echo "    $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/auth/register \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"name\":\"Joao\",\"email\":\"joao@email.com\",\"cpf\":\"000\",\"password\":\"123456\"}'"
echo ""
echo "Test (login):"
echo "  curl -i -X POST \\"
echo "    $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/auth/login \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"email\":\"joao@email.com\",\"password\":\"123456\"}'"
