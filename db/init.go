package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var sqliteConn *sql.DB

func init() {
	var err error
	sqliteConn, err = sql.Open("sqlite3", "grats.db?cache=shared")
	sqliteConn.SetMaxOpenConns(1)

	if err != nil {
		panic(err)
	}
	if err = sqliteConn.Ping(); err != nil {
		panic(err)
	}

	err = create_tables()
	if err != nil {
		log.Println("Error configuring database")
		panic("Database is not configured")
	}

	log.Println("Database is ready")
}

func create_table(create_table_sql string) error {
	_, err := sqliteConn.Exec(create_table_sql)
	if err != nil {
		log.Println("Error when trying to prepare statement during creating tables")
		log.Println(err)
		return err
	}

	return nil
}

func create_tables() error {
	create_user_table_sql := `CREATE TABLE IF NOT EXISTS user (
		id VARCHAR PRIMARY KEY,
		tgid INTEGER,
		name VARCHAR,
		tgusername VARCHAR,
		chatid VARCHAR,
		birthday VARCHAR,
		isadmin INTEGER,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(tgid)
	);`

	create_friend_table_sql := `CREATE TABLE IF NOT EXISTS friend (
		id VARCHAR PRIMARY KEY,
		name VARCHAR,
		birthday VARCHAR,
		chatid INTEGER,
		userid INTEGER,
		notifyat VARCHAR,
		createdat VARCHAR,
		updatedat VARCHAR
	);`

	create_access_table_sql := `CREATE TABLE IF NOT EXISTS access (
		id VARCHAR PRIMARY KEY,
		tgusername VARCHAR,
		createdat VARCHAR,
		updatedat VARCHAR,

		UNIQUE(tgusername)
	);`

	for _, table := range []string{
		create_user_table_sql,
		create_friend_table_sql,
		create_access_table_sql,
	} {
		err := create_table(table)
		if err != nil {
			return err
		}
	}

	return nil
}
