package rotas

import (
	"API-gestar-bem/src/controllers"
	"net/http"
)

var rotasComentarios = []Rota{
	{
		URI:                "/publicacoes/{publicacaoId}/comentarios",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarComentario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/publicacoes/{publicacaoId}/comentarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarComentarios,
		RequerAutenticacao: true,
	},
	{
		URI:                "/comentarios/{comentarioId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarComentario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/comentarios/{comentarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarComentario,
		RequerAutenticacao: true,
	},
}
