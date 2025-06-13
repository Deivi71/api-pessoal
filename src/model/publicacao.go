package model

import (
	"errors"
	"strings"
	"time"
)

// Publicacao - representa uma publicação no sistema
type Publicacao struct {
	ID          uint64    `json:"id,omitempty"`
	Titulo      string    `json:"titulo,omitempty"`
	Conteudo    string    `json:"conteudo,omitempty"`
	AutorID     uint64    `json:"autorId,omitempty"`
	AutorNick   string    `json:"autorNick,omitempty"`
	AutorNome   string    `json:"autorNome,omitempty"`
	Curtidas    uint64    `json:"curtidas"`
	Comentarios uint64    `json:"comentarios"`
	CriadoEm    time.Time `json:"criadoem,omitempty"`
}

// Preparar - valida e formata os dados da publicação
func (publicacao *Publicacao) Preparar() error {
	if erro := publicacao.validar(); erro != nil {
		return erro
	}

	publicacao.formatar()
	return nil
}

// validar - verifica se os dados da publicação são válidos
func (publicacao *Publicacao) validar() error {
	if publicacao.Titulo == "" {
		return errors.New("o título é obrigatório e não pode estar em branco")
	}

	if publicacao.Conteudo == "" {
		return errors.New("o conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

// formatar - formata os dados da publicação
func (publicacao *Publicacao) formatar() {
	publicacao.Titulo = strings.TrimSpace(publicacao.Titulo)
	publicacao.Conteudo = strings.TrimSpace(publicacao.Conteudo)
}
