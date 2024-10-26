package db

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteConn *sql.DB

func MustSetup(dsn string, logger *slog.Logger) {
	var err error
	sqliteConn, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	sqliteConn.SetMaxOpenConns(1)

	if err = sqliteConn.Ping(); err != nil {
		log.Fatal(err)
	}

	logger.Info("Database is ready")
}

func GetDBConnection() *sql.DB {
	return sqliteConn
}
