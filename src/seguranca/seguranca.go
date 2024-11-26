package seguranca

import "golang.org/x/crypto/bcrypt"

//Hash recebe uma senha string e coloca hash nela
func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

//VerificarSenha compara uma senha string com uma c hash e retorna se são iguais
func VerificarSenha(senhaHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senhaString))
}
