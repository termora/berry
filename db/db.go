package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/starshine-sys/snowflake/v2"
	"github.com/termora/berry/common"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db/search"
	"github.com/termora/berry/db/search/pg"

	// pgx driver for migrations
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Debug is a debug logging function
var Debug = func(template string, args ...interface{}) {}

// Db ...
type DB struct {
	// Embedded search methods
	search.Searcher
	*pgxpool.Pool

	GuildCache *ttlcache.Cache

	Config common.Config

	Snowflake *snowflake.Generator

	Timeout time.Duration

	sentry    *sentry.Hub
	useSentry bool

	TermBaseURL string

	IncFunc func()
}

// Init ...
func Init(url string) (db *DB, err error) {
	guildCache := ttlcache.NewCache()
	guildCache.SetCacheSizeLimit(100)
	guildCache.SetTTL(10 * time.Minute)

	err = runMigrations(url)
	if err != nil {
		return nil, err
	}

	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %w", err)
	}
	conf.ConnConfig.LogLevel = pgx.LogLevelWarn
	conf.ConnConfig.Logger = zapadapter.NewLogger(log.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	db = &DB{
		Snowflake:  snowflake.NewGen(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		Pool:       pool,
		GuildCache: guildCache,
		Timeout:    10 * time.Second,
		Searcher:   pg.New(pool, Debug),
		IncFunc:    func() {},
	}

	return
}

//go:embed migrations
var fs embed.FS

func runMigrations(url string) (err error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: fs,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n != 0 {
		log.Infof("Performed %v migrations!", n)
	}

	err = db.Close()
	return err
}

// Time gets the time from a snowflake
func (db *DB) Time(s snowflake.ID) time.Time {
	t, _ := db.Snowflake.Parse(s)
	return t
}

// Context is a convenience method to get a context.Context with the database's timeout
func (db *DB) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), db.Timeout)
}

// QueryRow ...
func (db *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	go db.IncFunc()
	return db.Pool.QueryRow(ctx, sql, args...)
}

// Query ...
func (db *DB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	go db.IncFunc()
	return db.Pool.Query(ctx, sql, args...)
}

// Exec ...
func (db *DB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	go db.IncFunc()
	return db.Pool.Exec(ctx, sql, args...)
}
