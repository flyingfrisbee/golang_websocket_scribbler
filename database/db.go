package database

import (
	env "GithubRepository/golang_websocket_scribbler/environment"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DBConn struct {
	Conn *sql.DB
	sync.Mutex
}

func (dbc *DBConn) CreateUser(username string) (int, error) {
	row := dbc.Conn.QueryRow(
		`SELECT COUNT(*) FROM user
		WHERE username = ?`,
		username,
	)

	var a int
	err := row.Scan(&a)
	if err != nil {
		return -1, err
	}

	userExist := a != 0
	if userExist {
		return -1, fmt.Errorf("username %s already exist", username)
	}

	dbc.Lock()
	defer dbc.Unlock()

	result, err := dbc.Conn.Exec(
		`INSERT INTO user (username)
		VALUES (?)`,
		username,
	)
	if err != nil {
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(userID), nil
}

func (dbc *DBConn) CloseConnection() {
	dbc.Conn.Close()
}

func SetupDB() {
	DB = createDBConnection()
}

func createDBConnection() DBConn {
	db, err := sql.Open("sqlite3", env.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fileBytes, err := os.ReadFile(env.SQLFilePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(fileBytes))
	if err != nil {
		log.Fatal(err)
	}

	return DBConn{
		Conn: db,
	}
}

var (
	DB DBConn
)
