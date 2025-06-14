package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"API-gestar-bem/src/config"

	_ "github.com/go-sql-driver/mysql"
)

// ExecutarMigracoes executa as migrações do banco de dados
func ExecutarMigracoes() error {
	db, err := sql.Open("mysql", config.ConnectBD)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Ler o arquivo de migrações
	content, err := ioutil.ReadFile("src/database/migrations.sql")
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de migrações: %v", err)
	}

	// Dividir o arquivo em comandos SQL individuais
	commands := strings.Split(string(content), ";")

	// Executar cada comando SQL
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		_, err = db.Exec(cmd)
		if err != nil {
			return fmt.Errorf("erro ao executar migração: %v\nComando: %s", err, cmd)
		}
	}

	log.Println("Migrações executadas com sucesso!")
	return nil
}
