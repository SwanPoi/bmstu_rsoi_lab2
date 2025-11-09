package main

import (
	"log"

	handler "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/handler"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/repositories"
	server "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/server"
	services "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/services"
)

func main() {
	connString := repo.GetConnectionString(&repo.DatabaseConfig{
		Host: "postgres",
		Port: 5432,
		User: "program",
		Password: "test",
		Database: "cars",
	})

	db, err := repo.InitDb(connString)

	if err != nil {
		log.Fatal("Fail during db connection", err)
		return
	}

	repos := repo.NewRepository(db)
	service := services.NewServices(repos)
	handler := handler.NewHandler(service)

	srv := new(server.CommonServer)

	if err := srv.Run("8070", handler.SetupRoutes()); err != nil {
		log.Fatal("Fail during car server start", err)
		return
	}
}