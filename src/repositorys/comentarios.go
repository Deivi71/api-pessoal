package repositorys

import (
	"API-gestar-bem/src/model"
	"database/sql"
)

// Comentarios - representa um repositório de comentários
type Comentarios struct {
	db *sql.DB
}

// NewRepositoryComentarios - cria um repositório de comentários
func NewRepositoryComentarios(db *sql.DB) *Comentarios {
	return &Comentarios{db}
}

// Criar - insere um comentário no banco de dados
func (repositorio Comentarios) Criar(comentario model.Comentario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO comentarios (conteudo, autor_id, publicacao_id) VALUES (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(comentario.Conteudo, comentario.AutorID, comentario.PublicacaoID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// BuscarPorPublicacao - busca todos os comentários de uma publicação
func (repositorio Comentarios) BuscarPorPublicacao(publicacaoID uint64) ([]model.Comentario, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT c.id, c.conteudo, c.autor_id, u.nick, u.nome, c.publicacao_id, c.criadoEm
		FROM comentarios c 
		INNER JOIN usuarios u ON u.id = c.autor_id 
		WHERE c.publicacao_id = ?
		ORDER BY c.criadoEm ASC`,
		publicacaoID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var comentarios []model.Comentario

	for linhas.Next() {
		var comentario model.Comentario

		if erro = linhas.Scan(
			&comentario.ID,
			&comentario.Conteudo,
			&comentario.AutorID,
			&comentario.AutorNick,
			&comentario.AutorNome,
			&comentario.PublicacaoID,
			&comentario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		comentarios = append(comentarios, comentario)
	}

	return comentarios, nil
}

// BuscarPorPublicacaoComReacoes - busca todos os comentários de uma publicação com reações
func (repositorio Comentarios) BuscarPorPublicacaoComReacoes(publicacaoID, usuarioID uint64) ([]model.Comentario, error) {
	comentarios, erro := repositorio.BuscarPorPublicacao(publicacaoID)
	if erro != nil {
		return nil, erro
	}

	if len(comentarios) == 0 {
		return comentarios, nil
	}

	// Buscar reações para todos os comentários
	comentarioIDs := make([]uint64, len(comentarios))
	for i, comentario := range comentarios {
		comentarioIDs[i] = comentario.ID
	}

	repositorioReacoes := NewRepositoryReacoesComentarios(repositorio.db)
	reacoesPorComentario, erro := repositorioReacoes.BuscarReacoesPorComentarios(comentarioIDs, usuarioID)
	if erro != nil {
		return nil, erro
	}

	// Adicionar reações aos comentários
	for i := range comentarios {
		if reacoes, existe := reacoesPorComentario[comentarios[i].ID]; existe {
			comentarios[i].Reacoes = reacoes
		} else {
			comentarios[i].Reacoes = []model.ReacaoContador{}
		}
	}

	return comentarios, nil
}

// BuscarPorID - busca um comentário por ID
func (repositorio Comentarios) BuscarPorID(comentarioID uint64) (model.Comentario, error) {
	linha, erro := repositorio.db.Query(`
		SELECT c.id, c.conteudo, c.autor_id, u.nick, u.nome, c.publicacao_id, c.criadoEm
		FROM comentarios c 
		INNER JOIN usuarios u ON u.id = c.autor_id 
		WHERE c.id = ?`,
		comentarioID,
	)
	if erro != nil {
		return model.Comentario{}, erro
	}
	defer linha.Close()

	var comentario model.Comentario

	if linha.Next() {
		if erro = linha.Scan(
			&comentario.ID,
			&comentario.Conteudo,
			&comentario.AutorID,
			&comentario.AutorNick,
			&comentario.AutorNome,
			&comentario.PublicacaoID,
			&comentario.CriadoEm,
		); erro != nil {
			return model.Comentario{}, erro
		}
	}

	return comentario, nil
}

// Atualizar - atualiza um comentário no banco de dados
func (repositorio Comentarios) Atualizar(comentarioID uint64, comentario model.Comentario) error {
	statement, erro := repositorio.db.Prepare("UPDATE comentarios SET conteudo = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(comentario.Conteudo, comentarioID); erro != nil {
		return erro
	}

	return nil
}

// Deletar - deleta um comentário do banco de dados
func (repositorio Comentarios) Deletar(comentarioID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM comentarios WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(comentarioID); erro != nil {
		return erro
	}

	return nil
}

// ContarPorPublicacao - conta quantos comentários uma publicação tem
func (repositorio Comentarios) ContarPorPublicacao(publicacaoID uint64) (int, error) {
	linha, erro := repositorio.db.Query(
		"SELECT COUNT(*) FROM comentarios WHERE publicacao_id = ?",
		publicacaoID,
	)
	if erro != nil {
		return 0, erro
	}
	defer linha.Close()

	var count int
	if linha.Next() {
		if erro = linha.Scan(&count); erro != nil {
			return 0, erro
		}
	}

	return count, nil
}
