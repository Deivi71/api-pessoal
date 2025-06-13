package repositorys

import (
	"API-gestar-bem/src/model"
	"database/sql"
)

// Publicacoes - representa um repositório de publicações
type Publicacoes struct {
	db *sql.DB
}

// NewRepositoryPublicacoes - cria um repositório de publicações
func NewRepositoryPublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

// Criar - insere uma publicação no banco de dados
func (repositorio Publicacoes) Criar(publicacao model.Publicacao) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO publicacoes (titulo, conteudo, autor_id) VALUES (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AutorID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// BuscarPorID - busca uma publicação por ID
func (repositorio Publicacoes) BuscarPorID(publicacaoID uint64) (model.Publicacao, error) {
	linha, erro := repositorio.db.Query(`
		SELECT p.id, p.titulo, p.conteudo, p.autor_id, u.nick, u.nome, p.curtidas, 
		       (SELECT COUNT(*) FROM comentarios WHERE publicacao_id = p.id) as comentarios,
		       p.criadoEm 
		FROM publicacoes p 
		INNER JOIN usuarios u ON u.id = p.autor_id 
		WHERE p.id = ?`,
		publicacaoID,
	)
	if erro != nil {
		return model.Publicacao{}, erro
	}
	defer linha.Close()

	var publicacao model.Publicacao

	if linha.Next() {
		if erro = linha.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.AutorNick,
			&publicacao.AutorNome,
			&publicacao.Curtidas,
			&publicacao.Comentarios,
			&publicacao.CriadoEm,
		); erro != nil {
			return model.Publicacao{}, erro
		}
	}

	return publicacao, nil
}

// Buscar - busca todas as publicações dos usuários seguidos e próprias
func (repositorio Publicacoes) Buscar(usuarioID uint64) ([]model.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT DISTINCT p.id, p.titulo, p.conteudo, p.autor_id, u.nick, u.nome, p.curtidas,
		       (SELECT COUNT(*) FROM comentarios WHERE publicacao_id = p.id) as comentarios,
		       p.criadoEm
		FROM publicacoes p 
		INNER JOIN usuarios u ON u.id = p.autor_id 
		LEFT JOIN seguidores s ON p.autor_id = s.usuario_id 
		WHERE p.autor_id = ? OR s.seguidor_id = ?
		ORDER BY p.criadoEm DESC`,
		usuarioID, usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []model.Publicacao

	for linhas.Next() {
		var publicacao model.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.AutorNick,
			&publicacao.AutorNome,
			&publicacao.Curtidas,
			&publicacao.Comentarios,
			&publicacao.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// BuscarTodasPublicas - busca todas as publicações (feed público)
func (repositorio Publicacoes) BuscarTodasPublicas() ([]model.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT p.id, p.titulo, p.conteudo, p.autor_id, u.nick, u.nome, p.curtidas,
		       (SELECT COUNT(*) FROM comentarios WHERE publicacao_id = p.id) as comentarios,
		       p.criadoEm
		FROM publicacoes p 
		INNER JOIN usuarios u ON u.id = p.autor_id 
		ORDER BY p.criadoEm DESC`,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []model.Publicacao

	for linhas.Next() {
		var publicacao model.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.AutorNick,
			&publicacao.AutorNome,
			&publicacao.Curtidas,
			&publicacao.Comentarios,
			&publicacao.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// Atualizar - atualiza uma publicação no banco de dados
func (repositorio Publicacoes) Atualizar(publicacaoID uint64, publicacao model.Publicacao) error {
	statement, erro := repositorio.db.Prepare("UPDATE publicacoes SET titulo = ?, conteudo = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoID); erro != nil {
		return erro
	}

	return nil
}

// Deletar - deleta uma publicação do banco de dados
func (repositorio Publicacoes) Deletar(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM publicacoes WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}

// BuscarPorUsuario - busca publicações de um usuário específico
func (repositorio Publicacoes) BuscarPorUsuario(usuarioID uint64) ([]model.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT p.id, p.titulo, p.conteudo, p.autor_id, u.nick, u.nome, p.curtidas,
		       (SELECT COUNT(*) FROM comentarios WHERE publicacao_id = p.id) as comentarios,
		       p.criadoEm
		FROM publicacoes p 
		INNER JOIN usuarios u ON u.id = p.autor_id 
		WHERE p.autor_id = ?
		ORDER BY p.criadoEm DESC`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []model.Publicacao

	for linhas.Next() {
		var publicacao model.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.AutorNick,
			&publicacao.AutorNome,
			&publicacao.Curtidas,
			&publicacao.Comentarios,
			&publicacao.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// Curtir - adiciona uma curtida à publicação
func (repositorio Publicacoes) Curtir(publicacaoID, usuarioID uint64) error {
	// Primeiro, verifica se já curtiu
	statement, erro := repositorio.db.Prepare(
		"INSERT IGNORE INTO curtidas (usuario_id, publicacao_id) VALUES (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, publicacaoID); erro != nil {
		return erro
	}

	// Atualiza o contador de curtidas
	statement, erro = repositorio.db.Prepare(
		"UPDATE publicacoes SET curtidas = (SELECT COUNT(*) FROM curtidas WHERE publicacao_id = ?) WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID, publicacaoID); erro != nil {
		return erro
	}

	return nil
}

// Descurtir - remove uma curtida da publicação
func (repositorio Publicacoes) Descurtir(publicacaoID, usuarioID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM curtidas WHERE usuario_id = ? AND publicacao_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, publicacaoID); erro != nil {
		return erro
	}

	// Atualiza o contador de curtidas
	statement, erro = repositorio.db.Prepare(
		"UPDATE publicacoes SET curtidas = (SELECT COUNT(*) FROM curtidas WHERE publicacao_id = ?) WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID, publicacaoID); erro != nil {
		return erro
	}

	return nil
}
