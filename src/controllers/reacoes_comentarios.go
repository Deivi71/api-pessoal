package controllers

import (
	"API-gestar-bem/src/autentication"
	"API-gestar-bem/src/banco"
	"API-gestar-bem/src/repositorys"
	"API-gestar-bem/src/respostas"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AdicionarReacaoComentario adiciona uma reação a um comentário
func AdicionarReacaoComentario(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	comentarioID, erro := strconv.ParseUint(parametros["comentarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var dadosReacao struct {
		TipoReacao string `json:"tipo_reacao"`
	}

	if erro = json.Unmarshal(corpoRequisicao, &dadosReacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar tipo de reação
	tiposValidos := map[string]bool{
		"like":  true,
		"love":  true,
		"laugh": true,
		"wow":   true,
		"sad":   true,
		"angry": true,
	}

	if !tiposValidos[dadosReacao.TipoReacao] {
		respostas.Erro(w, http.StatusBadRequest, nil)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorys.NewRepositoryReacoesComentarios(db)
	if erro = repositorio.AdicionarReacao(comentarioID, usuarioID, dadosReacao.TipoReacao); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, nil)
}

// RemoverReacaoComentario remove uma reação de um comentário
func RemoverReacaoComentario(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	comentarioID, erro := strconv.ParseUint(parametros["comentarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorys.NewRepositoryReacoesComentarios(db)
	if erro = repositorio.RemoverReacao(comentarioID, usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, nil)
}

// BuscarReacoesComentario busca as reações de um comentário
func BuscarReacoesComentario(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	comentarioID, erro := strconv.ParseUint(parametros["comentarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorys.NewRepositoryReacoesComentarios(db)
	reacoes, erro := repositorio.BuscarReacoesPorComentario(comentarioID, usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, reacoes)
}
