package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPublicacao adciona uma nova publicacao no db
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//lendo requisição
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	//passando pra struct
	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequest, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	publicacao.AutorID = usuarioID
	//fazendo verificações
	if erro = publicacao.Preparar(); erro != nil {
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao.ID, erro = repositorio.Criar(publicacao)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusCreated, publicacao)

}

// BuscarPublicacoes traz as publicacoes que terao no feed do usuario(sua e de quem segue)
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.Buscar(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacoes)

}

// BuscarPublicacao traz uma publicacao pelo seu id
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro publicacaoId
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao, erro := repositorio.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacao)

}

// AtualizarPublicacao edita uma publicacao
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro publicacaoId
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacaoSalva, erro := repositorio.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//vendo se o id de quem fez a publi é o mesmo de quem ta logado
	if publicacaoSalva.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar uma puclicação que não seja sua"))
		return
	}
	//lendo requisição
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	//passando pra struct
	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequest, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	publicacao.AutorID = usuarioID
	//fazendo verificações
	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//usando repositorios denovo para agora atualizar de fato
	if erro = repositorio.Atualizar(publicacaoID, publicacao); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)

}

// DeletarPublicacao deleta uma puclicacao
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
	//Obtendo ID do token pra saber qual usuario está logado
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro publicacaoId
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacaoSalva, erro := repositorio.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//vendo se o id de quem fez a publi é o mesmo de quem ta logado
	if publicacaoSalva.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível deletar uma puclicação que não seja sua"))
		return
	}
	//usando repositorios para deletar de fato a publicacao
	if erro = repositorio.Deletar(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarPublicacoesPorUsuario traz todas publicações de um usuário
func BuscarPublicacoesPorUsuario(w http.ResponseWriter, r *http.Request) {
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro usuarioId
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacoes)
}

// CurtirPublicacao incrementa o número de curtidas de uma publicacao
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro publicacaoId
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	erro = repositorio.Curtir(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

// DescurtirPublicacao decrementa o número de curtidas de uma publicacao
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	//pegando parametro (url/{parametro})
	parametros := mux.Vars(r)
	//aqui alem de converter de str p uint to pegando só o parametro publicacaoId
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	erro = repositorio.Descurtir(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}
