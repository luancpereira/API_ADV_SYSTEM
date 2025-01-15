package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/luancpereira/APICheckout/core/database/sqlc"
)

var (
	DB_QUERIER sqlc.Querier
	CONN       *sql.DB
)

type Config struct{}

func (c Config) Start() {
	CONN = c.setup()
	DB_QUERIER = sqlc.New(CONN)
}

func (Config) setup() *sql.DB {
	dbSource := "postgres://apicheckout:apicheckout@0.0.0.0:5438/apicheckout?sslmode=disable"

	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		panic(err.Error())
	}

	err = conn.Ping()
	if err != nil {
		panic(err.Error())
	}

	return conn
}
