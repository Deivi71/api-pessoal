package middleware

import (
	"net/http"

	"API-gestar-bem/src/autentication"
	"API-gestar-bem/src/responses"
)

// Autenticar verifica se o usuário está autenticado
func Autenticar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := autentication.ValidarToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r)
	}
}
