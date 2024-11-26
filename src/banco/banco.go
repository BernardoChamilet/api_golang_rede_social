package banco

import (
	"api/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Conectar abre uma conexao com db
func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.Conexao)
	if erro != nil {
		return nil, erro
	}
	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}
	return db, nil
}
