package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as informações do usuário
func CriarToken(usuarioId uint64) (string, error) {
	//aqui vão ser as informações que o token vai conter
	permissoes := jwt.MapClaims{}
	//o usuario tem autorização
	permissoes["authorized"] = true
	//o tempo logado tem expiração de 6 horas
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	//O id do usuário logado
	permissoes["usuarioId"] = usuarioId
	//gerando token com um secret key (chave para encriptografar o token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))
}

// ValidarToken verifica se o token passado na requisição é váido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return erro
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("token inválido")
}

// ExtrairUSuarioID retorna o ID do usuarioID que está no token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return 0, erro
	}
	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64)
		if erro != nil {
			return 0, erro
		}
		return usuarioID, nil
	}
	return 0, errors.New("token inválido")

}

// extrairToken obtem o token no formato correto
func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	//token normalmente vem com duas palavras, uma bearer e outra o token em si
	//por isso usar um split por espaço pra saber se ele veio como deveria e só verificar o token
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

// verificando se a chave recebida é da família da secret key
func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
