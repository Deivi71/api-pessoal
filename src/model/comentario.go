package model

import (
	"errors"
	"strings"
	"time"
)

// Comentario - representa um comentário no sistema
type Comentario struct {
	ID           uint64           `json:"id,omitempty"`
	Conteudo     string           `json:"conteudo,omitempty"`
	AutorID      uint64           `json:"autorId,omitempty"`
	AutorNick    string           `json:"autorNick,omitempty"`
	AutorNome    string           `json:"autorNome,omitempty"`
	PublicacaoID uint64           `json:"publicacaoId,omitempty"`
	CriadoEm     time.Time        `json:"criadoem,omitempty"`
	Reacoes      []ReacaoContador `json:"reacoes,omitempty"`
}

// Preparar - valida e formata os dados do comentário
func (comentario *Comentario) Preparar() error {
	if erro := comentario.validar(); erro != nil {
		return erro
	}

	comentario.formatar()
	return nil
}

// validar - verifica se os dados do comentário são válidos
func (comentario *Comentario) validar() error {
	if comentario.Conteudo == "" {
		return errors.New("o conteúdo do comentário é obrigatório e não pode estar em branco")
	}

	if len(comentario.Conteudo) > 500 {
		return errors.New("o comentário não pode ter mais de 500 caracteres")
	}

	return nil
}

// formatar - formata os dados do comentário
func (comentario *Comentario) formatar() {
	comentario.Conteudo = strings.TrimSpace(comentario.Conteudo)
}
