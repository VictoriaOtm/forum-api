package database

import (
	"log"

	"io/ioutil"

	"github.com/jackc/pgx"
)

func newConnPool() *pgx.ConnPool {
	DBConnPool, err := pgx.NewConnPool(dbConnPoolConfig)
	if err != nil {
		log.Fatalln(err)
	}

	return DBConnPool
}

var DB = newConnPool()

func TxMustBegin() *pgx.Tx {
	tx, err := DB.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

func InitSchema(pathToSchemaFile string) {
	schemaStr, err := ioutil.ReadFile(pathToSchemaFile)
	if err != nil {
		panic(err)
	}

	tx := TxMustBegin()
	defer tx.Commit()

	_, err = tx.Exec(string(schemaStr))
	if err != nil {
		panic(err)
	}
}
