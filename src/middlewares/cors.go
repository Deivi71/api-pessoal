package middlewares

import (
	"net/http"
)

// CORS - middleware para permitir requisições de diferentes origens
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Permite requisições de qualquer origem (em produção, especifique domínios específicos)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Métodos HTTP permitidos
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Headers permitidos
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Permite credentials (cookies, headers de autorização)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Cache do preflight por 24 horas
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Se for uma requisição OPTIONS (preflight), responde com 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continua para o próximo handler
		next(w, r)
	}
}
