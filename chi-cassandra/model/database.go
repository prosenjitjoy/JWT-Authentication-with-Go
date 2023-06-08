package model

import (
	"log"

	"github.com/gocql/gocql"
)

// CREATE KEYSPACE chijwt WITH replication = {'class':'NetworkTopologyStrategy', 'datacenter1': 3};

const createTable = `
CREATE TABLE IF NOT EXISTS users (
	first_name text,
	last_name text,
	email text,
	password text,
	phone text,
	new_token text,
	refresh_token text,
	user_type text,
	user_id text,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY((last_name), email, phone)
)`

const createTableByUserID = `
CREATE MATERIALIZED VIEW IF NOT EXISTS getUserByID 
AS SELECT * FROM users 
WHERE user_id IS NOT NULL 
AND last_name IS NOT NULL 
AND email IS NOT NULL 
AND phone IS NOT NULL 
PRIMARY KEY(user_id, last_name, email, phone)
`

func DBSession() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1", "127.0.0.2", "127.0.0.3")
	cluster.Keyspace = "chijwt"

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("unable to connect to cassandra cluster:", err)
	}

	err = session.Query(createTable).Exec()
	if err != nil {
		log.Fatal("unable to create users table:", err)
	}

	err = session.Query(createTableByUserID).Exec()
	if err != nil {
		log.Fatal("unable to create user_id table:", err)
	}

	return session
}
