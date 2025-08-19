package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/pressly/goose/v3"
)

type Database struct {
	conn *pgxpool.Pool
	cfg  *config.Config
	ctx  context.Context
	User *User
	Home *Home
}

type User struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

type Home struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *Database) doMigrations() error {
	// Get a standard sql.DB from the pool for goose
	sqlDB := stdlib.OpenDBFromPool(db.conn)
	defer sqlDB.Close()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
}

func New(cfg *config.Config, ctx context.Context) (*Database, error) {
	conn, err := pgxpool.New(ctx, cfg.DatabaseAddr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	config := conn.Config()
	config.MaxConns = 20
	config.MinConns = 2
	user := &User{
		conn: conn,
		ctx:  ctx,
	}

	home := &Home{
		conn: conn,
		ctx:  ctx,
	}

	db := &Database{
		conn: conn,
		cfg:  cfg,
		ctx:  ctx,
		User: user,
		Home: home,
	}

	// Run migrations
	if err := db.doMigrations(); err != nil {
		return nil, err
	}
	fmt.Println("migrations done")

	return db, nil
}
