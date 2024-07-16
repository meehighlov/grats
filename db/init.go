package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Client *sql.DB

func init() {
	var err error
	Client, err = sql.Open("sqlite3", "grats.db")
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}

	create_tables()

	log.Println("Database is ready")
}

func create_table(create_table_sql string) error {
	_, err := Client.Exec(create_table_sql)
	if err != nil {
		log.Println("Error when trying to prepare statement during creating tables")
		log.Println(err)
		return err
	}

	return nil
}

func create_tables() error {
	create_user_table_sql := `CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY,
		name VARCHAR,
		tgusername VARCHAR,
		chatid VARCHAR,
		birthday VARCHAR,
		isadmin INTEGER
	);`

	create_friend_table_sql := `CREATE TABLE IF NOT EXISTS friend (
		id INTEGER PRIMARY KEY,
		name VARCHAR,
		birthday VARCHAR,
		chatid INTEGER,
		userid INTEGER
	);`

	for _, table := range []string{
		create_user_table_sql,
		create_friend_table_sql,
	} {
		err := create_table(table)
		if err != nil {
			return err
		}
	}

	return nil
}
