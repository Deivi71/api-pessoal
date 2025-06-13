package main

import (
	"API-gestar-bem/src/config"
	"API-gestar-bem/src/database"
	"API-gestar-bem/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Carregar configurações
	config.Carregar()

	// Executar migrações do banco de dados
	if err := database.ExecutarMigracoes(); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	// Gerar rotas
	r := router.Gerar()

	// Iniciar servidor
	fmt.Printf("Escutando na porta %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
