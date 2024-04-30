package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"subscription-service/utils"
)

func main() {
	cfg := utils.LoadConfig(".env")

	r := mux.NewRouter()
	address := cfg.IP + ":" + cfg.Port
	log.Printf("Server starting on %s", address)

	// Starting server
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
