package main

import (
	"fmt"

	"github.com/raviand/test-project/cmd/server"
	"github.com/raviand/test-project/internal/controller"
	"github.com/raviand/test-project/internal/repository"
	"github.com/raviand/test-project/internal/service"
)

func main() {
	fmt.Println("Starting server...")
	db := repository.NewDatabase()
	svc := service.NewProductService(db)
	ctrl := controller.NewHandler(svc)
	if err := server.NewServer(ctrl).Start(); err != nil {
		panic(err)
	}
}
