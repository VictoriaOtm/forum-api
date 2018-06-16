package database

import "github.com/jackc/pgx"

var dbConnConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "docker",
	User:     "docker",
	Password: "docker",
}

var dbConnPoolConfig = pgx.ConnPoolConfig{
	ConnConfig:     dbConnConfig,
	MaxConnections: 50,
}
