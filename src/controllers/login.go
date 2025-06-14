package controllers

import (
	"API-gestar-bem/src/autentication"
	"API-gestar-bem/src/banco"
	"API-gestar-bem/src/model"
	"API-gestar-bem/src/repositorys"
	"API-gestar-bem/src/respostas"
	"API-gestar-bem/src/security"
	"encoding/json"
	"io"
	"net/http"
)

// Login - vai autenticar um usuário na API
func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario model.Usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositorys.NewRepositoryUsuarios(db)
	usuarioSalvoNoBanco, erro := repository.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//Comparando a senha que o usuário enviou com a senha que está no banco
	if erro = security.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//Se a senha estiver correta, vamos criar o token
	token, erro := autentication.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	w.Write([]byte(token))
}
