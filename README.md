# 🚀 Ecommerce Microservices

> Um projeto de E-commerce baseado em **Microsserviços**, construído com **Go**, **AWS Lambda** e **LocalStack**.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![LocalStack](https://img.shields.io/badge/LocalStack-3.0+-blue?style=for-the-badge&logo=amazonaws&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15.0+-336791?style=for-the-badge&logo=postgresql&logoColor=white)

---

## 🏗️ Arquitetura

O projeto segue uma arquitetura baseada em eventos e funções serverless, simulando um ambiente AWS localmente.

```mermaid
graph LR
    Client[Client / Postman] -->|HTTP| APIG[API Gateway (LocalStack)]
    APIG -->|Invoke| Lambda[AWS Lambda (Go)]
    Lambda -->|Proxy Request| Gin[Gin Router]
    Gin -->|Handle| Controller[Handler Layer]
    Controller -->|Business Logic| Service[Service Layer]
    Service -->|Data Access| Repo[Repository Layer]
    Repo -->|SQL| DB[(PostgreSQL)]
```

### 🛠️ Tech Stack

- **Linguagem Principal**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin) (com adapter para Lambda)
- **Database**: PostgreSQL (Driver `pgx`)
- **Infraestrutura**:
  - **Docker Compose**: Orquestração dos containers.
  - **LocalStack**: Emulação de serviços AWS (Lambda, API Gateway, S3, SQS).
- **Ferramentas**:
  - AWS CLI (para interagir com o LocalStack).

---

## 📂 Estrutura do Projeto

```bash
ecommerce-microservices/
├── auth-service/           # Microsserviço de Autenticação
│   ├── cmd/               # Entrypoints
│   ├── handler/           # Controladores HTTP
│   ├── service/           # Regras de Negócio
│   ├── repository/        # Acesso ao Banco de Dados
│   └── scripts/           # Scripts de automação (deploy local)
├── docker-compose.yml     # Definição da Infraestrutura (LocalStack + Postgres)
└── README.md              # Documentação
```

---

## 🚀 Como Rodar o Projeto

Siga os passos abaixo para levantar o ambiente e fazer o deploy da função Lambda localmente.

### 1️⃣ Pré-requisitos
- **Docker** e **Docker Compose** instalados.
- **Go 1.25+** instalado.
- **AWS CLI** (opcional, mas recomendado para debug).
- **Zip** (para empacotar a lambda).

### 2️⃣ Iniciar a Infraestrutura
Na raiz do projeto, suba os containers do LocalStack e Postgres:

```bash
docker compose up -d
```
> **Nota:** Certifique-se de que as portas `4566` (LocalStack) e `5433` (Postgres Host) estão livres.

### 3️⃣ Build & Deploy da Lambda
Utilize o script facilitador para compilar o código Go, criar o zip e fazer o deploy no LocalStack:

```bash
cd auth-service
sh scripts/bootstrap-localstack.sh
```

O script irá:
1. Compilar o binário `bootstrap` (Linux/AMD64).
2. Criar o arquivo `function.zip`.
3. Criar a **Lambda Function** no LocalStack.
4. Criar o **API Gateway** e configurar a rota `/register`.
5. Retornar a **URL completa** para teste.

---

## 🧪 Testando a API

Após rodar o script de deploy, você receberá uma URL. Exemplo de chamada `POST /register`:

```bash
curl -X POST \
  http://localhost:4566/restapis/<API_ID>/dev/_user_request_/register \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Henry Komatsu",
    "email": "henry@example.com",
    "password": "secure_password"
  }'
```

---

## ⚙️ Variáveis de Ambiente

As configurações principais estão no `docker-compose.yml` e no script de deploy.

| Serviço | Variável | Valor Padrão | Descrição |
|---------|----------|--------------|-----------|
| **Postgres** | `POSTGRES_USER` | `admin` | Usuário do DB |
| **Postgres** | `POSTGRES_PASSWORD` | `secret` | Senha do DB |
| **Postgres** | `POSTGRES_DB` | `micro_db` | Nome do Database |
| **App** | `DATABASE_URL` | `postgres://...` | Connection String (Injetada na Lambda) |

---
*Documentação gerada automaticamente por Antigravity.*
