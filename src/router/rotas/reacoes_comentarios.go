package rotas

import (
	"API-gestar-bem/src/controllers"
	"net/http"
)

var rotasReacoesComentarios = []Rota{
	{
		URI:                "/comentarios/{comentarioId}/reacoes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarReacaoComentario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/comentarios/{comentarioId}/reacoes",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverReacaoComentario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/comentarios/{comentarioId}/reacoes",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarReacoesComentario,
		RequerAutenticacao: true,
	},
}
