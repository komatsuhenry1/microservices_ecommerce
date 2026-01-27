# Guia de Deploy Manual: AWS Lambda + API Gateway (LocalStack)

Este guia detalha cada passo dos comandos utilizados para configurar o ambiente serverless localmente.

## 1. Configuração de Variáveis de Ambiente
Antes de executar os comandos `aws`, definimos variáveis para facilitar a reutilização e garantir que o CLI aponte para o LocalStack.

```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export LAMBDA_ENDPOINT=http://localhost:4566
```
*   **AWS_ACCESS_KEY_ID / SECRET**: Credenciais dummy aceitas pelo LocalStack.
*   **LAMBDA_ENDPOINT**: Aponta os comandos para o LocalStack (`localhost:4566`) em vez da AWS real.

Definição de nomes e estágios:
```bash
export FUNCTION_NAME=auth-service
export API_NAME=ecommerce-api
export STAGE=dev
```

---

## 2. Build e Empacotamento (Go)
O Lambda precisa de um executável compatível com Linux.

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap .
zip -j function.zip bootstrap
```
*   **GOOS=linux GOARCH=amd64**: Compila para Linux/AMD64 (ambiente padrão do Lambda).
*   **-o bootstrap**: O runtime `provided.al2` da AWS espera um executável chamado exatamente `bootstrap`.
*   **zip**: Compacta o binário em um arquivo `.zip` para upload.

---

## 3. Criar a Função Lambda
Cria a função no LocalStack com o código zipado.

```bash
aws lambda create-function \
  --function-name "$FUNCTION_NAME" \
  --runtime provided.al2 \
  --handler bootstrap \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --zip-file fileb://function.zip \
  --timeout 10 \
  --memory-size 256 \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
*   **--runtime provided.al2**: Especifica que é um runtime customizado (necessário para Go moderno).
*   **--role**: Um ARN fictício de IAM Role (o LocalStack aceita qualquer um, mas é obrigatório).
*   **fileb://**: Indica que o arquivo é binário.

---

## 4. Atualizar Configuração (Variáveis de Ambiente)
Injeta a string de conexão do banco de dados na Lambda.

```bash
aws lambda update-function-configuration \
  --function-name "$FUNCTION_NAME" \
  --environment "Variables={DATABASE_URL=postgres://admin:secret@db:5432/auth_db?sslmode=disable}" \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
*   Isso permite que sua aplicação Go leia `os.Getenv("DATABASE_URL")` para conectar ao Postgres.

---

## 5. Criar API Gateway (Rest API)
Cria a entrada da API.

```bash
API_ID=$(aws apigateway create-rest-api \
  --name "$API_NAME" \
  --query 'id' \
  --output text \
  --endpoint-url "$LAMBDA_ENDPOINT")
echo "API_ID=$API_ID"
```
*   Cria a API e salva o `id` gerado na variável `API_ID` para uso nos próximos passos.

---

## 6. Obter o Recurso Raiz (Root)
Toda API tem um diretório raiz (`/`). Precisamos do ID dele para criar sub-rotas.

```bash
ROOT_ID=$(aws apigateway get-resources \
  --rest-api-id "$API_ID" \
  --query 'items[?path==`/`].id | [0]' \
  --output text \
  --endpoint-url "$LAMBDA_ENDPOINT")
```

---

## 7. Criar Recurso (Rota `/register`)
Cria o caminho `/register` dentro da API.

```bash
REGISTER_ID=$(aws apigateway create-resource \
  --rest-api-id "$API_ID" \
  --parent-id "$ROOT_ID" \
  --path-part register \
  --query 'id' \
  --output text \
  --endpoint-url "$LAMBDA_ENDPOINT")
```
*   **--parent-id**: Vincula ao `ROOT_ID`, criando a estrutura `/ -> /register`.

---

## 8. Definir Método HTTP (POST)
Define que a rota `/register` aceita o verbo `POST`.

```bash
aws apigateway put-method \
  --rest-api-id "$API_ID" \
  --resource-id "$REGISTER_ID" \
  --http-method POST \
  --authorization-type NONE \
  --no-api-key-required \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
*   **--authorization-type NONE**: A API é pública (sem auth por enquanto).

---

## 9. Integração (Vincular API -> Lambda)
Este é o passo mais complexo. Configura o API Gateway para repassar a requisição para a Lambda.

```bash
aws apigateway put-integration \
  --rest-api-id "$API_ID" \
  --resource-id "$REGISTER_ID" \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:auth-service/invocations" \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
*   **--type AWS_PROXY**: A Lambda recebe o evento bruto (headers, body completo) e deve retornar o formato esperado pelo API Gateway.
*   **--integration-http-method POST**: O API Gateway invoca a Lambda sempre via POST (padrão interno da AWS).
*   **--uri**: O ARN de invocação "mágico" construido com o ARN da função. Observe a estrutura: `arn:aws:apigateway:<region>:lambda:path/.../functions/<lambda-arn>/invocations`.

---

## 10. Permissões de Invocação
Autoriza o serviço API Gateway a executar a Lambda.

```bash
aws lambda add-permission \
  --function-name auth-service \
  --statement-id apigw-register \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:us-east-1:000000000000:$API_ID/*/POST/register" \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
*   Sem isso, a chamada daria erro de permissão (403 ou 500 interno) ao tentar invocar a função.

---

### Passo Adicional Recomendado: Deploy da API
Seus comandos terminavam na permissão, mas para a API ficar acessível via URL de stage (ex: `/dev/register`), geralmente é necessário criar um deployment:

```bash
aws apigateway create-deployment \
  --rest-api-id "$API_ID" \
  --stage-name "$STAGE" \
  --endpoint-url "$LAMBDA_ENDPOINT"
```
