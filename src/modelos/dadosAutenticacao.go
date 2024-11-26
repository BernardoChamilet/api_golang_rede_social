package modelos

//DadosAutenticao contém token e id do usuáio autenticado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
