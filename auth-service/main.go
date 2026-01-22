package main

import (
	"auth-service/di"
	"auth-service/router"
	"auth-service/utils"
	"context"
	"log"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

//no go init() é uma função especial que roda antes da main()

func init() {
	// 1) Conecta no DB (do jeito que você já faz hoje, ou via utils/connect)
	db := utils.ConnectRDS()
	container := di.NewContainer(db)

	// 2) Sobe o router Gin (mas SEM Run/Listen)
	r := router.SetupRoutes(container)

	// 3) Cria o adapter Gin <-> API Gateway
	ginLambda = ginadapter.New(r)

	log.Println("Lambda init ok")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Esse adapter entende o evento do API Gateway e roda seu Gin
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
