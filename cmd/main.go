package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/raviand/test-project/cmd/server"
	"github.com/raviand/test-project/internal/audit"
	"github.com/raviand/test-project/internal/controller"
	"github.com/raviand/test-project/internal/repository"
	"github.com/raviand/test-project/internal/service"
)

func main() {

	fmt.Println("Starting server...")
	// db := repository.NewDatabase(os.Getenv("DB_FILE_PATH"))
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/movies_db", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
	mysql, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db := repository.NewDatabase(mysql)
	svc := service.NewProductService(db)
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
