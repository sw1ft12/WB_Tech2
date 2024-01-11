package main

import (
	"fmt"
	"log"
	"net/http"
	"task11/api"
	"task11/config"
	"task11/middleware"
)

func main() {
	cfg, err := config.ReadCfg()
	if err != nil {
		log.Fatal(err)
	}
	handler := api.NewHandler()
	mux := http.NewServeMux()
	wrappedMux := middleware.Logging(mux)
	handler.InitRoutes(mux)
	serverAddress := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	fmt.Println(serverAddress)
	err = http.ListenAndServe(serverAddress, wrappedMux)
	if err != nil {
		log.Fatal(err)
	}
}
