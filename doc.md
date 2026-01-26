fluxo end to end:

curl / Postman
   ↓
API Gateway (LocalStack)
   ↓
Lambda (Go)
   ↓
Gin Router
   ↓
Handler
   ↓
Service
   ↓
Repository
   ↓
Postgres


---

docker compose up -d

AWS fake → localhost:4566

Postgres → localhost:5432

---

comandos de inicio:

❯ aws lambda create-function \
  --function-name auth-service \
  --runtime provided.al2 \
  --handler bootstrap \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --zip-file fileb://function.zip \
  --endpoint-url=http://localhost:4566

  ❯ API_ID=$(aws apigateway create-rest-api \
  --name ecommerce-api \
  --query 'id' \
  --output text \
  --endpoint-url=http://localhost:4566)

  ❯ ROOT_ID=$(aws apigateway get-resources \
  --rest-api-id $API_ID \
  --query 'items[0].id' \
  --output text \
  --endpoint-url=http://localhost:4566)

  ❯ REGISTER_ID=$(aws apigateway create-resource \
  --rest-api-id $API_ID \
  --parent-id $ROOT_ID \
  --path-part register \
  --query 'id' \
  --output text \
  --endpoint-url=http://localhost:4566)

  aws apigateway put-method \
  --rest-api-id $API_ID \
  --resource-id $REGISTER_ID \
  --http-method POST \
  --authorization-type NONE \
  --endpoint-url=http://localhost:4566

aws apigateway put-integration \
  --rest-api-id hc1tnv47gs \
  --resource-id olc0ymus75 \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:auth-service/invocations \
  --endpoint-url=http://localhost:4566


aws apigateway create-deployment \
  --rest-api-id hc1tnv47gs \
  --stage-name dev \
  --endpoint-url=http://localhost:4566

  


----


codigo para atualizar o código de um microsserviço:
(ex auth-service)

aws lambda update-function-code \
  --function-name auth-service \
  --zip-file fileb://function.zip \
  --endpoint-url=http://localhost:4566 \
  --region us-east-1