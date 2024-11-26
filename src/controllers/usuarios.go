package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario recebe a requisição http, e trata para inserir dados no banco com auxilio do pacote repositorios.
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	//lendo requisição
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	//passando para struct
	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//fazendo verificações(cryptografando senha tm)
	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrindo banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// interagindo com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios retorna dados de todos os usuários do db
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	//r.URL.Get("algumacoisa") pega o algumacoisa que está em url/usuarios?x=algumacoisa?y=slaoq
	nomeOunick := strings.ToLower(r.URL.Query().Get("usuario"))
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOunick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuario retorna dados de um usuário do db
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	//pegando parametro (url/{parametro})
	//aqui vem todos parametros
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro usuarioId
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrino db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario atualiza os dados de um usuario no db
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	//lendo parametros
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioIDtoken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//se o usuario logado estiver tentando atualizar os dados de outro usuario
	if usuarioID != usuarioIDtoken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o logado"))
		return
	}

	//lendo requisição
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	//passando pra struct
	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//fazendo verificações
	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Atualizar(usuarioID, usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)

}

// DeletarUsuário deleta um usuário do db
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	//lendo parametros
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioIDtoken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//se o usuario logado estiver tentando deletar os dados de outro usuario
	if usuarioID != usuarioIDtoken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível deletar um usuário que não seja o logado"))
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Deletar(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// PararDeSeguirUsuario é utilizada quando um usuario já está logado para ele seguir um outro usuário
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado, o seguidor
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//lendo parametros para obter id do usuario que ele quer deixar de seguir
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//se o usuario logado estiver tentando deixar de seguir ele mesmo
	if usuarioID == seguidorID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível parar de seguir você mesmo"))
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.PararDeSeguir(usuarioID, seguidorID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuario é utilizada quando um usuario já está logado para ele seguir um outro usuário
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado, o seguidor
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//lendo parametros para obter id do usuario que ele quer seguir
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//se o usuario logado estiver tentando seguir ele mesmo
	if usuarioID == seguidorID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível seguir você mesmo"))
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Seguir(usuarioID, seguidorID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores traz todos os seguidores de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	//lendo parametros para obter id do usuario que quero ver seguidores
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo traz todos os usuários que um usuário está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	//lendo parametros para obter id do usuario que quero ver quem segue
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguindo, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, seguindo)
}

// AtualizarSenha atualiza a senha de um usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioIDtoken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//lendo parametros para obter id do usuario que terá senha atualizada
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//se o usuario logado estiver tentando atualizar senha de outro usuário
	if usuarioIDtoken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("só é possível atualizar sua própria senha"))
		return
	}
	//lendo requisição
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	//passando para struct
	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequest, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//abrindo db
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco para obter senha salva
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalva, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//vendo se a senha obtida do banco é igual a que o usuário digitou
	if erro = seguranca.VerificarSenha(senhaSalva, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("senha atual não condiz com que está no banco"))
		return
	}
	//colocando hash na senha nova obtida da requisição
	senhaHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//usando metodos do repositorio para interagir com banco para inserir senha nova
	if erro := repositorio.AtualizarSenha(usuarioID, string(senhaHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}
