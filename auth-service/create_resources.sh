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

# DB URL
DATABASE_URL=${DATABASE_URL:-postgres://admin:secret@db:5432/auth_db?sslmode=disable}

echo "==> LocalStack endpoint: $LS_ENDPOINT"
echo "==> Region: $AWS_DEFAULT_REGION"
echo "==> Function: $FUNC_NAME"
echo "==> API name: $API_NAME"
echo "==> Stage: $STAGE"

# ===== Build Lambda (Go custom runtime) =====
echo "==> Building lambda binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap .
chmod +x bootstrap
zip -j function.zip bootstrap >/dev/null

# ===== (Re)Create Lambda =====
echo "==> Recreating lambda: $FUNC_NAME"
aws --region "$AWS_DEFAULT_REGION" lambda delete-function \
  --function-name "$FUNC_NAME" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null 2>&1 || true

aws --region "$AWS_DEFAULT_REGION" lambda create-function \
  --function-name "$FUNC_NAME" \
  --runtime provided.al2 \
  --handler bootstrap \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --zip-file fileb://function.zip \
  --timeout 10 \
  --memory-size 256 \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

wait_lambda_active() {
  local name="$1"
  local tries="${2:-30}"  # 30 tentativas
  local sleep_s="${3:-1}"

  echo "==> Waiting lambda to become Active..."
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

wait_lambda_active "$FUNC_NAME"

echo "==> Setting lambda env vars (DATABASE_URL)..."
aws --region "$AWS_DEFAULT_REGION" lambda update-function-configuration \
  --function-name "$FUNC_NAME" \
  --environment "Variables={DATABASE_URL=$DATABASE_URL}" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

# ===== (Re)Create API Gateway REST API =====
echo "==> Removing existing REST APIs with name '$API_NAME' (if any)..."
EXISTING_IDS=$(aws --region "$AWS_DEFAULT_REGION" apigateway get-rest-apis \
  --endpoint-url "$LS_ENDPOINT" \
  --query "items[?name=='$API_NAME'].id" \
  --output text || true)

if [[ -n "${EXISTING_IDS// }" ]]; then
  for id in $EXISTING_IDS; do
    echo "   - deleting api: $id"
    aws --region "$AWS_DEFAULT_REGION" apigateway delete-rest-api \
      --rest-api-id "$id" \
      --endpoint-url "$LS_ENDPOINT" >/dev/null || true
  done
fi

echo "==> Creating REST API: $API_NAME"
API_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-rest-api \
  --name "$API_NAME" \
  --query 'id' \
  --output text \
  --endpoint-url "$LS_ENDPOINT")

echo "==> API_ID: $API_ID"

ROOT_ID=$(aws --region "$AWS_DEFAULT_REGION" apigateway get-resources \
  --rest-api-id "$API_ID" \
  --query 'items[?path==`/`].id | [0]' \
  --output text \
  --endpoint-url "$LS_ENDPOINT")

echo "==> ROOT_ID: $ROOT_ID"

# Helper to create resource+method+integration+permission
create_route () {
  local path_part="$1"       # register | login
  local http_method="$2"     # POST
  local statement_id="$3"    # apigw-register | apigw-login

  echo "==> Creating resource '/$path_part'..."
  local resource_id
  resource_id=$(aws --region "$AWS_DEFAULT_REGION" apigateway create-resource \
    --rest-api-id "$API_ID" \
    --parent-id "$ROOT_ID" \
    --path-part "$path_part" \
    --query 'id' \
    --output text \
    --endpoint-url "$LS_ENDPOINT")

  echo "   - resource_id: $resource_id"

  echo "==> Putting method $http_method on '/$path_part'..."
  aws --region "$AWS_DEFAULT_REGION" apigateway put-method \
    --rest-api-id "$API_ID" \
    --resource-id "$resource_id" \
    --http-method "$http_method" \
    --authorization-type NONE \
    --no-api-key-required \
    --endpoint-url "$LS_ENDPOINT" >/dev/null

  echo "==> Integrating $http_method '/$path_part' -> Lambda (AWS_PROXY)..."
  aws --region "$AWS_DEFAULT_REGION" apigateway put-integration \
    --rest-api-id "$API_ID" \
    --resource-id "$resource_id" \
    --http-method "$http_method" \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri "arn:aws:apigateway:$AWS_DEFAULT_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_DEFAULT_REGION:000000000000:function:$FUNC_NAME/invocations" \
    --endpoint-url "$LS_ENDPOINT" >/dev/null

  # Permission: ignore if already exists (script recria API, mas a lambda pode manter policy entre runs dependendo do LocalStack)
  echo "==> Adding lambda permission for API Gateway ($statement_id)..."
  aws --region "$AWS_DEFAULT_REGION" lambda add-permission \
    --function-name "$FUNC_NAME" \
    --statement-id "$statement_id" \
    --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn "arn:aws:execute-api:$AWS_DEFAULT_REGION:000000000000:$API_ID/*/$http_method/$path_part" \
    --endpoint-url "$LS_ENDPOINT" >/dev/null 2>&1 || true
}

# Routes
create_route "register" "POST" "apigw-register"
create_route "login" "POST" "apigw-login"

# Deploy API 
echo "==> Deploying stage: $STAGE"
aws --region "$AWS_DEFAULT_REGION" apigateway create-deployment \
  --rest-api-id "$API_ID" \
  --stage-name "$STAGE" \
  --endpoint-url "$LS_ENDPOINT" >/dev/null

echo ""
echo "✅ Done!"
echo ""
echo "Register URL:"
echo "  $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/register"
echo "Login URL:"
echo "  $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/login"
echo ""
echo "Test (register):"
echo "  curl -i -X POST \\"
echo "    $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/register \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"name\":\"Joao\",\"email\":\"joao@email.com\",\"cpf\":\"000\",\"password\":\"123456\"}'"
echo ""
echo "Test (login):"
echo "  curl -i -X POST \\"
echo "    $LS_ENDPOINT/restapis/$API_ID/$STAGE/_user_request_/login \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"email\":\"joao@email.com\",\"password\":\"123456\"}'"
