package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/vp2305/common"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	handler := NewHandler()
	mux := handler.registerRoutes()

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
