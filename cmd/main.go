package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/cmd/server/handler"
	"github.com/raviand/test-project/internal/repository/memory"
)

func main() {
	// Inicializacion del servicio
	file := os.Getenv("FILE_NAME")
	db := memory.NewDatabase(file)
	// mongo := mongodb.NewMongoDb()
	svc := handler.NewHandler(db)
	rt := chi.NewRouter()

	// declaracion de endpoints
	rt.Route("/product", func(rt chi.Router) {
		rt.Post("/", svc.CreateProduct)
		rt.Get("/code/{code}", svc.GetProductByCode)
		rt.Get("/id/{id}", svc.GetProductById)
		rt.Patch("/{id}", svc.PatchProduct)
	})

	// levantar el servicio
	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
