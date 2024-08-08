package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/raviand/test-project/cmd/server"
	"github.com/raviand/test-project/internal/audit"
	"github.com/raviand/test-project/internal/controller"
	"github.com/raviand/test-project/internal/repository"
	"github.com/raviand/test-project/internal/service"
)

func main() {
	fmt.Println("Starting server...")
	mysql, err := InitDB()
	if err != nil {
		panic(err)
	}
	dynamo, err := InitDynamo()
	if err != nil {
		panic(err)
	}
	ddb := repository.NewDynamoRepository(dynamo, "User")
	db := repository.NewDatabase(mysql)
	svc := service.NewProductService(db, ddb)
	var wg sync.WaitGroup
	notiChannel, auditor := audit.NewAuditRoutine(&wg)
	ctrl := controller.NewHandler(svc, notiChannel)
	wg.Add(1)
	go auditor.Run()
	if err := server.NewServer(ctrl).Start(); err != nil {
		panic(err)
	}
	close(notiChannel)
	wg.Wait()
}

func InitDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/movies_db", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
	mysql, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return mysql, nil
}

func InitDynamo() (*dynamodb.DynamoDB, error) {
	region := "us-west-2"
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	cred := credentials.NewStaticCredentials("local", "local", "")
	sess, err := session.NewSession(aws.NewConfig().WithEndpoint(endpoint).WithRegion(region).WithCredentials(cred))
	if err != nil {
		return nil, err
	}
	dynamo := dynamodb.New(sess)
	return dynamo, nil
}
