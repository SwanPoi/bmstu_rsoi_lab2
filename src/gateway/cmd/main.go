package main

import (
	"log"

	server "github.com/SwanPoi/bmstu_rsoi_lab2/src/gateway/server"
	handler "github.com/SwanPoi/bmstu_rsoi_lab2/src/gateway/handler"
	services "github.com/SwanPoi/bmstu_rsoi_lab2/src/gateway/services"
)

func main() {
	services := services.NewServices()
	handler := handler.NewHandler(services)

	srv := new(server.CommonServer)


	if err := srv.Run("8080", handler.SetupRoutes()); err != nil {
		log.Fatal("Fail during gateway server start: ", err)
		return
	}
}