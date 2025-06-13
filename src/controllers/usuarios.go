package controllers

import (
	"API-gestar-bem/src/autentication"
	"API-gestar-bem/src/banco"
	"API-gestar-bem/src/model"
	"API-gestar-bem/src/repositorys"
	"API-gestar-bem/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// CriarUsuario - vai criar um usuário
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var usuario model.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("cadastro"); erro != nil {
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
	usuario.ID, erro = repository.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios - vai buscar todos os usuários
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositorys.NewRepositoryUsuarios(db)
	usuarios, erro := repository.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuarios)

}

// BuscarUsuario - vai buscar um usuário
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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
	repository := repositorys.NewRepositoryUsuarios(db)

	usuario, erro := repository.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario - vai atualizar um usuário
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	usuarioIDNoToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario model.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("edição"); erro != nil {
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
	if erro = repository.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario - vai deletar um usuário do DB
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível deletar um usuário que não seja o seu"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	respository := repositorys.NewRepositoryUsuarios(db)
	if erro = respository.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuario - permite que um usuário siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {

	seguidorID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositorys.NewRepositoryUsuarios(db)

	// Verificar se o usuário a ser seguido existe
	_, erro = repository.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}

	if erro = repository.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

func ParardeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível parar de seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositorys.NewRepositoryUsuarios(db)

	// Verificar se o usuário existe
	_, erro = repository.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}

	if erro = repository.PararDeSeguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores - busca os seguidores de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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

	repository := repositorys.NewRepositoryUsuarios(db)
	seguidores, erro := repository.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo - busca os usuários que um usuário está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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

	repository := repositorys.NewRepositoryUsuarios(db)
	seguindo, erro := repository.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguindo)
}

// BuscarEstatisticasUsuario - busca estatísticas de um usuário (seguidores e seguindo)
func BuscarEstatisticasUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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

	repository := repositorys.NewRepositoryUsuarios(db)

	seguidores, erro := repository.ContarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	seguindo, erro := repository.ContarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	estatisticas := struct {
		Seguidores int `json:"seguidores"`
		Seguindo   int `json:"seguindo"`
	}{
		Seguidores: seguidores,
		Seguindo:   seguindo,
	}

	respostas.JSON(w, http.StatusOK, estatisticas)
}

// UploadFotoPerfil - faz upload da foto de perfil do usuário
func UploadFotoPerfil(w http.ResponseWriter, r *http.Request) {
	// Verificar se o usuário está autenticado
	usuarioIDNoToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Verificar se o usuário está tentando atualizar sua própria foto
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar a foto de perfil de outro usuário"))
		return
	}

	// Parse do multipart form (limite de 10MB)
	erro = r.ParseMultipartForm(10 << 20)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("arquivo muito grande (máximo 10MB)"))
		return
	}

	// Obter o arquivo do form
	arquivo, handler, erro := r.FormFile("foto")
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("erro ao obter o arquivo"))
		return
	}
	defer arquivo.Close()

	// Verificar tipo de arquivo (apenas imagens)
	contentType := handler.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		respostas.Erro(w, http.StatusBadRequest, errors.New("apenas arquivos de imagem são permitidos"))
		return
	}

	// Criar diretório de uploads se não existir
	uploadDir := "uploads/perfil"
	if erro := os.MkdirAll(uploadDir, 0755); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao criar diretório de upload"))
		return
	}

	// Gerar nome único para o arquivo
	extensao := filepath.Ext(handler.Filename)
	nomeArquivo := fmt.Sprintf("perfil_%d_%d%s", usuarioID, time.Now().Unix(), extensao)
	caminhoCompleto := filepath.Join(uploadDir, nomeArquivo)

	// Criar arquivo no servidor
	arquivoDestino, erro := os.Create(caminhoCompleto)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao criar arquivo no servidor"))
		return
	}
	defer arquivoDestino.Close()

	// Copiar conteúdo do arquivo
	_, erro = io.Copy(arquivoDestino, arquivo)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao salvar arquivo"))
		return
	}

	// Atualizar banco de dados
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositorys.NewRepositoryUsuarios(db)
	if erro = repository.AtualizarFotoPerfil(usuarioID, nomeArquivo); erro != nil {
		// Se falhar ao atualizar BD, remover arquivo
		os.Remove(caminhoCompleto)
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	resposta := struct {
		Mensagem   string `json:"mensagem"`
		FotoPerfil string `json:"foto_perfil"`
	}{
		Mensagem:   "Foto de perfil atualizada com sucesso",
		FotoPerfil: nomeArquivo,
	}

	respostas.JSON(w, http.StatusOK, resposta)
}

// ServirFotoPerfil - serve as imagens de perfil dos usuários
func ServirFotoPerfil(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	nomeArquivo := parametros["arquivo"]

	// Validar nome do arquivo para evitar path traversal
	if strings.Contains(nomeArquivo, "..") || strings.Contains(nomeArquivo, "/") {
		respostas.Erro(w, http.StatusBadRequest, errors.New("nome de arquivo inválido"))
		return
	}

	caminhoArquivo := filepath.Join("uploads/perfil", nomeArquivo)

	// Verificar se o arquivo existe
	if _, erro := os.Stat(caminhoArquivo); os.IsNotExist(erro) {
		respostas.Erro(w, http.StatusNotFound, errors.New("imagem não encontrada"))
		return
	}

	// Servir o arquivo
	http.ServeFile(w, r, caminhoArquivo)
}
