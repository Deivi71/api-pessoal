package repositorys

import (
	"database/sql"
	"fmt"

	"API-gestar-bem/src/model"
)

// ReacoesComentarios representa um repositório de reações aos comentários
type ReacoesComentarios struct {
	db *sql.DB
}

// NewRepositoryReacoesComentarios cria um repositório de reações aos comentários
func NewRepositoryReacoesComentarios(db *sql.DB) *ReacoesComentarios {
	return &ReacoesComentarios{db}
}

// AdicionarReacao adiciona ou atualiza uma reação a um comentário
func (repositorio ReacoesComentarios) AdicionarReacao(comentarioID, usuarioID uint64, tipoReacao string) error {
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO reacoes_comentarios (comentario_id, usuario_id, tipo) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE tipo = VALUES(tipo)
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(comentarioID, usuarioID, tipoReacao)
	return erro
}

// RemoverReacao remove uma reação de um comentário
func (repositorio ReacoesComentarios) RemoverReacao(comentarioID, usuarioID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM reacoes_comentarios WHERE comentario_id = ? AND usuario_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(comentarioID, usuarioID)
	return erro
}

// BuscarReacoesPorComentario busca todas as reações de um comentário com contagem
func (repositorio ReacoesComentarios) BuscarReacoesPorComentario(comentarioID, usuarioID uint64) ([]model.ReacaoContador, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			tipo,
			COUNT(*) as quantidade,
			MAX(CASE WHEN usuario_id = ? THEN 1 ELSE 0 END) as eu_reagi
		FROM reacoes_comentarios 
		WHERE comentario_id = ? 
		GROUP BY tipo
		ORDER BY quantidade DESC
	`, usuarioID, comentarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var reacoes []model.ReacaoContador

	for linhas.Next() {
		var reacao model.ReacaoContador
		var euReagi int

		if erro = linhas.Scan(
			&reacao.TipoReacao,
			&reacao.Quantidade,
			&euReagi,
		); erro != nil {
			return nil, erro
		}

		reacao.ComentarioID = comentarioID
		reacao.EuReagi = euReagi == 1
		reacoes = append(reacoes, reacao)
	}

	return reacoes, nil
}

// BuscarReacoesPorComentarios busca reações para múltiplos comentários
func (repositorio ReacoesComentarios) BuscarReacoesPorComentarios(comentarioIDs []uint64, usuarioID uint64) (map[uint64][]model.ReacaoContador, error) {
	if len(comentarioIDs) == 0 {
		return make(map[uint64][]model.ReacaoContador), nil
	}

	// Construir a query com placeholders
	placeholders := ""
	args := []interface{}{usuarioID}
	for i, id := range comentarioIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT 
			comentario_id,
			tipo,
			COUNT(*) as quantidade,
			MAX(CASE WHEN usuario_id = ? THEN 1 ELSE 0 END) as eu_reagi
		FROM reacoes_comentarios 
		WHERE comentario_id IN (%s)
		GROUP BY comentario_id, tipo
		ORDER BY comentario_id, quantidade DESC
	`, placeholders)

	linhas, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	reacoesPorComentario := make(map[uint64][]model.ReacaoContador)

	for linhas.Next() {
		var comentarioID uint64
		var reacao model.ReacaoContador
		var euReagi int

		if erro = linhas.Scan(
			&comentarioID,
			&reacao.TipoReacao,
			&reacao.Quantidade,
			&euReagi,
		); erro != nil {
			return nil, erro
		}

		reacao.ComentarioID = comentarioID
		reacao.EuReagi = euReagi == 1
		reacoesPorComentario[comentarioID] = append(reacoesPorComentario[comentarioID], reacao)
	}

	return reacoesPorComentario, nil
}
