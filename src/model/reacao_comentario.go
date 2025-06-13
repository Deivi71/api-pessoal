package model

import (
	"errors"
	"time"
)

// ReacaoComentario representa uma reação a um comentário no sistema
type ReacaoComentario struct {
	ID           string    `json:"id"`
	Tipo         string    `json:"tipo"`
	UsuarioID    string    `json:"usuario_id"`
	ComentarioID string    `json:"comentario_id"`
	CriadoEm     time.Time `json:"criado_em"`
}

// ReacaoContador representa o contador de reações de um comentário
type ReacaoContador struct {
	ComentarioID uint64 `json:"comentario_id"`
	TipoReacao   string `json:"tipo_reacao"`
	Quantidade   int    `json:"quantidade"`
	EuReagi      bool   `json:"eu_reagi"`
}

// PrepararReacaoComentario prepara uma reação para ser salva
func (reacao *ReacaoComentario) PrepararReacaoComentario() error {
	if erro := reacao.validarReacaoComentario(); erro != nil {
		return erro
	}
	return nil
}

// validarReacaoComentario valida os campos da reação
func (reacao *ReacaoComentario) validarReacaoComentario() error {
	if reacao.Tipo == "" {
		return errors.New("o tipo da reação é obrigatório")
	}

	// Validar se o tipo é um dos tipos permitidos
	tiposValidos := map[string]bool{
		"like":    true,
		"dislike": true,
		"love":    true,
		"hate":    true,
		"laugh":   true,
		"angry":   true,
	}

	if !tiposValidos[reacao.Tipo] {
		return errors.New("tipo de reação inválido")
	}

	return nil
}
