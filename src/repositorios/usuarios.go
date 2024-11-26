package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Usuarios representa o repositório de usuários
type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositorio de usuarios
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// funcs abaixo são metodos do struct usuarios q representa um repositorio. O repositorio contém o db devbook(pq foi o db aberto no controller) e metodos.
// o parametro usuario do tipo modelos.Usuario é o struct obtido no controller do corpo da requisição

// Criar insere um usuário no banco de dados
func (repositorio Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	//criando declaração de inserção e a executando
	statement, erro := repositorio.db.Prepare(
		"insert into usuarios (nome,nick,email,senha) values (?,?,?,?)")
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
	//retorna o id do usuario inserido
	return uint64(ultimoIDInserido), nil

}

// Buscar traz todos os usuários que atendem o filtro de nome ou nicks
func (repositorio Usuarios) Buscar(nomeOUnick string) ([]modelos.Usuario, error) {
	nomeOUnick = fmt.Sprintf("%%%s%%", nomeOUnick) // %nomeOUnick% pra usar o comando alike do sql
	//pegando usuarios que tenham nome ou nick igual ou contendo nomeOUnick
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoem from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOUnick, nomeOUnick,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	var usuarios []modelos.Usuario
	//percorrendo linhas obtidas e adcionando num slice um struct para cada usuario
	for linhas.Next() {
		var usuario modelos.Usuario
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

// BuscarPorID traz os dados de um usuário por seu id
func (repositorio Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	//selecionando usuario que tenha o id recebido
	linha, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoem from usuarios where id = ?", ID)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()
	//passando os dados do usuario para uma struct e a retornando
	var usuario modelos.Usuario
	if linha.Next() {
		if erro = linha.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Atualizar atualiza os dados de usuario exceto a senha
func (repositorio Usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID)
	if erro != nil {
		return erro
	}
	return nil
}

// Deletar deleta os dados de um usuário
func (repositorio Usuarios) Deletar(ID uint64) error {
	//criando declaração de deletar e a executando
	statement, erro := repositorio.db.Prepare(
		"delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(ID)
	if erro != nil {
		return erro
	}
	return nil
}

// BuscarPorEmail busca o id e senha de um usuario do banco usando email
func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	//selecionando usuario que tenha o email recebido
	linha, erro := repositorio.db.Query(
		"select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()
	//passando os dados do usuario para uma struct e a retornando
	var usuario modelos.Usuario
	if linha.Next() {
		if erro = linha.Scan(
			&usuario.ID,
			&usuario.Senha,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Seguir faz o usuário de id seguidorID seguir o usuário de id usuarioID
func (repositorio Usuarios) Seguir(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare("insert ignore into seguidores (usuario_id, seguidor_id) values (?,?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(usuarioID, seguidorID)
	if erro != nil {
		return erro
	}
	return nil
}

// PararDeSeguir faz o usuário de id seguidorID parar de seguir o usuário de id usuarioID
func (repositorio Usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from seguidores where usuario_id=? and seguidor_id=?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(usuarioID, seguidorID)
	if erro != nil {
		return erro
	}
	return nil
}

// BuscarSeguidores busca todos seguidores de um usuario de id usuarioID
func (repositorio Usuarios) BuscarSeguidores(usuarioID uint64) ([]modelos.Usuario, error) {
	//selecionando linhas que tenha o usuarioID como seguido (campo usuario_id)
	linhas, erro := repositorio.db.Query(
		"select u.id, u.nome, u.nick, u.email, u.criadoem from usuarios u inner join seguidores s on u.id = s.seguidor_id where s.usuario_id=?", usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	var seguidores []modelos.Usuario
	for linhas.Next() {
		var seguidor modelos.Usuario
		if erro = linhas.Scan(
			&seguidor.ID,
			&seguidor.Nome,
			&seguidor.Nick,
			&seguidor.Email,
			&seguidor.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		seguidores = append(seguidores, seguidor)
	}
	return seguidores, nil
}

// BuscarSeguindo traz todos usuários que um usuário de id usuarioID está seguindo
func (repositorio Usuarios) BuscarSeguindo(usuarioID uint64) ([]modelos.Usuario, error) {
	//selecionando linhas que tenha o usuarioID como seguidor (campo seguidor_id)
	linhas, erro := repositorio.db.Query(
		"select u.id, u.nome, u.nick, u.email, u.criadoem from usuarios u inner join seguidores s on u.id = s.usuario_id where s.seguidor_id=?", usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	var seguindo []modelos.Usuario
	for linhas.Next() {
		var seguido modelos.Usuario
		if erro = linhas.Scan(
			&seguido.ID,
			&seguido.Nome,
			&seguido.Nick,
			&seguido.Email,
			&seguido.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		seguindo = append(seguindo, seguido)
	}
	return seguindo, nil
}

// BuscarSenha busca a senha de um usuario do banco usando id
func (repositorio Usuarios) BuscarSenha(ID uint64) (string, error) {
	//selecionando usuario que tenha o id recebido
	linha, erro := repositorio.db.Query(
		"select senha from usuarios where id = ?", ID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()
	//convertendo pra string
	var senha modelos.Senha
	if linha.Next() {
		if erro = linha.Scan(
			&senha.Atual,
		); erro != nil {
			return "", erro
		}
	}

	return senha.Atual, nil
}

// AtualizarSenha atualiza a senha de um usuario
func (repositorio Usuarios) AtualizarSenha(ID uint64, senha string) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare(
		"update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(senha, ID)
	if erro != nil {
		return erro
	}
	return nil
}
