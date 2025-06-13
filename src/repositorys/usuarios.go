package repositorys

import (
	"API-gestar-bem/src/model"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

// NewRepositoryUsuarios - vai criar um novo repositório de usuários
func NewRepositoryUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// Criar - vai criar um usuário no banco de dados
func (r usuarios) Criar(usuario model.Usuario) (uint64, error) {
	statement, erro := r.db.Prepare("insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}
	return uint64(ultimoIDInserido), nil

}

func (r usuarios) Buscar(nomeOuNick string) ([]model.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)
	linhas, erro := r.db.Query(
		"select id, nome, nick, email, criadoem from usuarios where nome like ? or nick like ?",
		nomeOuNick, nomeOuNick,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	var usuarios []model.Usuario
	for linhas.Next() {
		var usuario model.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (r usuarios) BuscarPorID(ID uint64) (model.Usuario, error) {
	linhas, erro := r.db.Query(
		"select id, nome, nick, email, foto_perfil, criadoem from usuarios where id = ?",
		ID,
	)
	if erro != nil {
		return model.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario model.Usuario

	if linhas.Next() {
		var fotoPerfil sql.NullString
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&fotoPerfil,
			&usuario.CriadoEm,
		); erro != nil {
			return model.Usuario{}, erro
		}

		if fotoPerfil.Valid {
			usuario.FotoPerfil = fotoPerfil.String
		}
	} else {
		return model.Usuario{}, fmt.Errorf("usuário com ID %d não encontrado", ID)
	}
	return usuario, nil
}

func (r usuarios) Atualizar(ID uint64, usuario model.Usuario) error {
	statement, erro := r.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}
	return nil
}

func (r usuarios) Deletar(ID uint64) error {
	statement, erro := r.db.Prepare(
		"delete from usuarios where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}
	return nil
}

func (r usuarios) BuscarPorEmail(email string) (model.Usuario, error) {
	linha, erro := r.db.Query(
		"select id, senha from usuarios where email = ?",
		email,
	)
	if erro != nil {
		return model.Usuario{}, erro
	}
	defer linha.Close()

	var usuario model.Usuario

	if linha.Next() {
		if erro = linha.Scan(
			&usuario.ID,
			&usuario.Senha,
		); erro != nil {
			return model.Usuario{}, erro
		}
	}
	return usuario, nil
}

func (r usuarios) Seguir(usuarioID, seguidorID uint64) error {
	statement, erro := r.db.Prepare(
		"insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}
	return nil
}

func (r usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	statement, erro := r.db.Prepare(
		"delete from seguidores where usuario_id = ? and seguidor_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}
	return nil
}

// BuscarSeguidores - busca os seguidores de um usuário
func (r usuarios) BuscarSeguidores(usuarioID uint64) ([]model.Usuario, error) {
	linhas, erro := r.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem 
		FROM usuarios u 
		INNER JOIN seguidores s ON u.id = s.seguidor_id 
		WHERE s.usuario_id = ?`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []model.Usuario

	for linhas.Next() {
		var usuario model.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarSeguindo - busca os usuários que um usuário está seguindo
func (r usuarios) BuscarSeguindo(usuarioID uint64) ([]model.Usuario, error) {
	linhas, erro := r.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem 
		FROM usuarios u 
		INNER JOIN seguidores s ON u.id = s.usuario_id 
		WHERE s.seguidor_id = ?`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []model.Usuario

	for linhas.Next() {
		var usuario model.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// ContarSeguidores - conta quantos seguidores um usuário tem
func (r usuarios) ContarSeguidores(usuarioID uint64) (int, error) {
	linha, erro := r.db.Query(
		"SELECT COUNT(*) FROM seguidores WHERE usuario_id = ?",
		usuarioID,
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

// ContarSeguindo - conta quantos usuários um usuário está seguindo
func (r usuarios) ContarSeguindo(usuarioID uint64) (int, error) {
	linha, erro := r.db.Query(
		"SELECT COUNT(*) FROM seguidores WHERE seguidor_id = ?",
		usuarioID,
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

// AtualizarFotoPerfil - atualiza a foto de perfil de um usuário
func (r usuarios) AtualizarFotoPerfil(usuarioID uint64, caminhoFoto string) error {
	statement, erro := r.db.Prepare(
		"UPDATE usuarios SET foto_perfil = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(caminhoFoto, usuarioID); erro != nil {
		return erro
	}

	return nil
}
