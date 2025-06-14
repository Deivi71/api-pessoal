package main

import (
	"log"
	"net/http"

	"API-gestar-bem/src/router"
)

func main() {
	router.Configurar()
	log.Fatal(http.ListenAndServe(":5000", nil))
}
