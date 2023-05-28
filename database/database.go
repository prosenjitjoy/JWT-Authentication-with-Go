package database

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id int GENERATED ALWAYS AS IDENTITY,
    first_name text NOT NULL,
    last_name text NOT NULL,
	password text NOT NULL,
	email text NOT NULL,
	phone text NOT NULL,
	token text,
	user_type text NOT NULL,
	refresh_token text,
	created_at timestamp,
	updated_at timestamp,
	user_id text NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(email, phone, token, user_id)
);
`

func DBInstance() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	db.MustExec(schema)

	return db
}
