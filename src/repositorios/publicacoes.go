package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Publicacoes representa o repositório de publicações
type Publicacoes struct {
	db *sql.DB
}

// NovoRepositorioDePublicacoes cria um repositorio de publicações
func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

// Criar insere uma publicação no banco de dados
func (repositorio Publicacoes) Criar(publicacao modelos.Publicacao) (uint64, error) {
	//criando declaração de inserção e a executando
	statement, erro := repositorio.db.Prepare(
		"insert into publicacoes (titulo,conteudo,autor_id) values (?,?,?)")
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
	//retorna o id da publicação inserido
	return uint64(ultimoIDInserido), nil
}

// BuscarPorID insere uma publicação no banco de dados
func (repositorio Publicacoes) BuscarPorID(publicacaoID uint64) (modelos.Publicacao, error) {
	//selecionando publicacao que tenha o id recebido
	linha, erro := repositorio.db.Query(
		"select p.*, u.nick from publicacoes p inner join usuarios u on u.id = p.autor_id where p.id=?", publicacaoID)
	if erro != nil {
		return modelos.Publicacao{}, erro
	}
	defer linha.Close()
	//passando os dados da publicacao para uma struct e a retornando
	var publicacao modelos.Publicacao
	if linha.Next() {
		if erro = linha.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return modelos.Publicacao{}, erro
		}
	}
	return publicacao, nil
}

// Buscar traz todas as publicações do usuario com usuarioID e de todos os usuários que ele segue
func (repositorio Publicacoes) Buscar(usuarioID uint64) ([]modelos.Publicacao, error) {
	//selecionando dados da tabela
	linhas, erro := repositorio.db.Query("select distinct p.*, u.nick from publicacoes p inner join usuarios u on u.id = p.autor_id left join seguidores s on p.autor_id = s.usuario_id where u.id=? or s.seguidor_id=? order by 1 desc", usuarioID, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	//passando os dados das publicações para um slice de structs e o retornando
	var publicacoes []modelos.Publicacao
	for linhas.Next() {
		var publicacao modelos.Publicacao
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}
	return publicacoes, nil
}

// Atualizar altera os dados de uma publicação no banco de dados
func (repositorio Publicacoes) Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare(
		"update publicacoes set titulo = ?, conteudo = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoID); erro != nil {
		return erro
	}
	return nil
}

// Deletar deleta os dados de uma publicação no banco de dados
func (repositorio Publicacoes) Deletar(publicacaoID uint64) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare(
		"delete from publicacoes where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}

// BuscarPorUsuario traz todas publicacoes de um usuario do banco de dados
func (repositorio Publicacoes) BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error) {
	//selecioando publicações
	linhas, erro := repositorio.db.Query("select p.*, u.nick from publicacoes p join usuarios u on u.id = p.autor_id where p.autor_id=?", usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	//passando os dados das publicações para um slice de structs e o retornando
	var publicacoes []modelos.Publicacao
	for linhas.Next() {
		var publicacao modelos.Publicacao
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}
	return publicacoes, nil

}

// Curtir incrementa o número de curtidas de uma publicação
func (repositorio Publicacoes) Curtir(publicacaoID uint64) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare("update publicacoes set curtidas = curtidas + 1 where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}

// Descurtir incrementa o número de curtidas de uma publicação
func (repositorio Publicacoes) Descurtir(publicacaoID uint64) error {
	//criando declaração de atualização e a executando
	statement, erro := repositorio.db.Prepare("update publicacoes set curtidas = CASE WHEN curtidas > 0 THEN curtidas - 1 ELSE curtidas END where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}
