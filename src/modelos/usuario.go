package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa a entidade usuario do db
type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoem,omitempty"`
}

// Preparar irá validar e formatar os dados do usuário recebido
func (usuario *Usuario) Preparar(momento string) error {
	if erro := usuario.validar(momento); erro != nil {
		return erro
	}
	if erro := usuario.formatar(momento); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(momento string) error {
	if usuario.Nome == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}
	if usuario.Nick == "" {
		return errors.New("o nick é obrigatório e não pode estar em branco")
	}
	if usuario.Email == "" {
		return errors.New("o email é obrigatório e não pode estar em branco")
	}
	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("o email é inserido é inváido")
	}
	if momento == "cadastro" && usuario.Senha == "" {
		return errors.New("a senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (usuario *Usuario) formatar(momento string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)
	if momento == "cadastro" {
		senhaHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}
		usuario.Senha = string(senhaHash)
	}
	return nil
}
