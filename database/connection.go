package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig set configuration for Connect method purpose
type DBConfig struct {
	Type   string
	Config string
}

// GetType of database type, like MySQL, PostgreSQL or others
// func GetDBType(dbType string) string {
// 	return dbType
// }

func GetDBConfig() string {
	return "postgres://postgres:rootatdawn@localhost/company_profile?sslmode=disable"
}

// OpenDB to the postgresql database
func OpenDB() (*sql.DB, error) {
	var dbc DBConfig
	dbc.Type = "postgres"
	dbc.Config = GetDBConfig()

	db, err := sql.Open(dbc.Type, dbc.Config)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Second * 10)
	// db.SetConnMaxLifetime(time.Second * 3)
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
