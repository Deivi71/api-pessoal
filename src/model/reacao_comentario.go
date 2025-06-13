package model

import "time"

// ReacaoComentario representa uma reação a um comentário
type ReacaoComentario struct {
	ID           uint64    `json:"id,omitempty"`
	ComentarioID uint64    `json:"comentario_id,omitempty"`
	UsuarioID    uint64    `json:"usuario_id,omitempty"`
	TipoReacao   string    `json:"tipo_reacao,omitempty"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
}

// ReacaoContador representa a contagem de reações por tipo
type ReacaoContador struct {
	TipoReacao string `json:"tipo"`
	Quantidade uint64 `json:"quantidade"`
	EuReagi    bool   `json:"eu_reagi"`
}
