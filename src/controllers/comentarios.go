package controllers

import (
	"API-gestar-bem/src/autentication"
	"API-gestar-bem/src/banco"
	"API-gestar-bem/src/model"
	"API-gestar-bem/src/repositorys"
	"API-gestar-bem/src/respostas"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarComentario - cria um novo comentário
func CriarComentario(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var comentario model.Comentario
	if erro = json.Unmarshal(corpoRequisicao, &comentario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	comentario.AutorID = usuarioID
	comentario.PublicacaoID = publicacaoID

	if erro = comentario.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Verificar se a publicação existe
	repositorioPublicacoes := repositorys.NewRepositoryPublicacoes(db)
	_, erro = repositorioPublicacoes.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, errors.New("publicação não encontrada"))
		return
	}

	repositorio := repositorys.NewRepositoryComentarios(db)
	comentario.ID, erro = repositorio.Criar(comentario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, comentario)
}

// BuscarComentarios - busca todos os comentários de uma publicação
func BuscarComentarios(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repositorio := repositorys.NewRepositoryComentarios(db)
	comentarios, erro := repositorio.BuscarPorPublicacao(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, comentarios)
}

// AtualizarComentario - atualiza um comentário
func AtualizarComentario(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorys.NewRepositoryComentarios(db)
	comentarioSalvoBanco, erro := repositorio.BuscarPorID(comentarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if comentarioSalvoBanco.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um comentário que não seja seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var comentario model.Comentario
	if erro = json.Unmarshal(corpoRequisicao, &comentario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = comentario.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.Atualizar(comentarioID, comentario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarComentario - deleta um comentário
func DeletarComentario(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorys.NewRepositoryComentarios(db)
	comentarioSalvoBanco, erro := repositorio.BuscarPorID(comentarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if comentarioSalvoBanco.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível deletar um comentário que não seja seu"))
		return
	}

	if erro = repositorio.Deletar(comentarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
