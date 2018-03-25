package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
)

var DBConnPool *pgx.ConnPool

var dbConnConfig = pgx.ConnConfig{
	Host:              "localhost",
	Port:              5432,
	Database:          "docker",
	User:              "docker",
	Password:          "docker",
	TLSConfig:         nil,
	UseFallbackTLS:    false,
	FallbackTLSConfig: nil,
}

var dbConnPoolConfig = pgx.ConnPoolConfig{
	ConnConfig:     dbConnConfig,
	MaxConnections: 50,
	AcquireTimeout: 0,
}

func InitConnPool() {
	var err error

	DBConnPool, err = pgx.NewConnPool(dbConnPoolConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func InitSchema(path string) {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := DBConnPool.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	_, err = tx.Exec(string(schema))
	if err != nil {
		log.Fatalln(err)
		tx.Rollback()
	}
}

func mustPrepare(name, stmt string) {
	_, err := DBConnPool.Prepare(name, stmt)
	if err != nil {
		log.Fatalln(name, err)
	}
}

func PrepareStatements() {
	postPrepareStatements()
	userPrepareStatements()
	forumPrepareStaements()
	prepareThreadStatements()
}
