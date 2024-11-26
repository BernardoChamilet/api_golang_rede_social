package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// Login faz o loginde um usuário
func Login(w http.ResponseWriter, r *http.Request) {
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
	//abrindo banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//usando metodos do repositorio para interagir com banco (detalhes na func criarusuario)
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvo, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//verificando se senha coincide
	if erro = seguranca.VerificarSenha(usuarioSalvo.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//gerando token do usuario e mandando na resposta(para testes)
	token, erro := autenticacao.CriarToken(usuarioSalvo.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	usuarioId := strconv.FormatUint(usuarioSalvo.ID, 10)

	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioId, Token: token})
}
