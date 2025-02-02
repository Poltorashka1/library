package storage

import (
	"book/internal/config"
	"book/internal/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// todo add query in prarams

type DB interface {
	ExecContext(ctx context.Context, query Query) (sql.Result, error)
	QueryContext(ctx context.Context, query Query) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query Query) *sql.Row
}

type Query struct {
	QueryName string
	Query     string
	Args      []any
}

// todo mb optimize or use other version
func (q *Query) StringV1() string {
	queryString := fmt.Sprintf("sql: %s: query: %s", q.QueryName, q.Query)
	if len(q.Args) != 0 {
		for k, v := range q.Args {
			queryString = strings.Replace(queryString, fmt.Sprintf("$%d", k+1), fmt.Sprintf("%v", v), 1)
		}
	}

	return queryString
}
func (q *Query) StringV2() string {
	return fmt.Sprintf("sql: %s, query: %s, args: %v", q.QueryName, q.Query, q.Args)
}

type db struct {
	// todo mb ne nado
	cfg       config.DBConfig
	dbConnect *sql.DB
	logger    logger.Logger
}

func (db *db) ExecContext(ctx context.Context, query Query) (sql.Result, error) {
	db.logQuery(query)
	return db.dbConnect.ExecContext(ctx, query.Query, query.Args...)
}
func (db *db) QueryContext(ctx context.Context, query Query) (*sql.Rows, error) {
	db.logQuery(query)
	return db.dbConnect.QueryContext(ctx, query.Query, query.Args...)
}
func (db *db) QueryRowContext(ctx context.Context, query Query) *sql.Row {
	db.logQuery(query)
	return db.dbConnect.QueryRowContext(ctx, query.Query)
}

func NewQuery(queryName string, query string, args []any) Query {
	return Query{
		QueryName: queryName,
		Query:     query,
		Args:      args,
	}
}

func NewDB(ctx context.Context, logger logger.Logger, config config.DBConfig) DB {
	connect, err := sql.Open("sqlite3", config.DSN())
	if err != nil {
		logger.Fatal("error connect to database")
	}
	storage := &db{
		cfg:       config,
		dbConnect: connect,
		logger:    logger,
	}

	if err := storage.dbConnect.Ping(); err != nil {
		logger.Fatal("error connect to database")
	}

	return storage
}

func (db *db) logQuery(query Query) {
	db.logger.Debug(query.StringV1())
}
