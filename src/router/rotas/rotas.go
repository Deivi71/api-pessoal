package rotas

import (
	"API-gestar-bem/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Rotas - vai receber as rotas
type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(w http.ResponseWriter, r *http.Request)
	RequerAutenticacao bool
}

// Configurar - vai configurar as rotas do router
func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsuarios
	rotas = append(rotas, rotaLogin)
	rotas = append(rotas, rotasPublicacoes...)
	rotas = append(rotas, rotasComentarios...)

	for _, rota := range rotas {

		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI, middlewares.CORS(middlewares.Logger(middlewares.Autenticar(rota.Funcao)))).Methods(rota.Metodo, "OPTIONS")
		} else {
			r.HandleFunc(rota.URI, middlewares.CORS(middlewares.Logger(rota.Funcao))).Methods(rota.Metodo, "OPTIONS")
		}
	}
	return r
}
