package main

import (
	"flotify/internal/config"
	"flotify/internal/database"
	"flotify/internal/handler"
	"fmt"
	"net/http"
)

func main() {
	dbpool := database.GetDatabasePool()
	defer dbpool.Close()

	router := handler.InitRouter(dbpool)

	server_config := config.LoadServerConfig()
	addr := fmt.Sprintf("%s:%s", server_config.Host, server_config.Port)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}
	panic(server.ListenAndServe())

}
