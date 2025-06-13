package main

import (
	"API-gestar-bem/src/config"
	"API-gestar-bem/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	config.Carregar()
	r := router.Gerar()

	fmt.Printf("Escutando na porta %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))

}
