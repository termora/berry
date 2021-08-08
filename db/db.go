package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/starshine-sys/snowflake/v2"
	"github.com/termora/berry/db/search"
	"github.com/termora/berry/db/search/pg"
	"github.com/termora/berry/structs"
	"go.uber.org/zap"

	// pgx driver for migrations
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Debug is a debug logging function
var Debug = func(template string, args ...interface{}) {}

// Db ...
type Db struct {
	// Embedded search methods
	search.Searcher

	Pool       *pgxpool.Pool
	Sugar      *zap.SugaredLogger
	GuildCache *ttlcache.Cache

	Config *structs.BotConfig

	Snowflake *snowflake.Generator

	Timeout time.Duration

	sentry    *sentry.Hub
	useSentry bool

	TermBaseURL string
}

// Init ...
func Init(url string, sugar *zap.SugaredLogger) (db *Db, err error) {
	guildCache := ttlcache.NewCache()
	guildCache.SetCacheSizeLimit(100)
	guildCache.SetTTL(10 * time.Minute)

	err = runMigrations(url, sugar)
	if err != nil {
		return nil, err
	}

	logger := sugar.Desugar()
	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %w", err)
	}
	conf.ConnConfig.LogLevel = pgx.LogLevelWarn
	conf.ConnConfig.Logger = zapadapter.NewLogger(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	db = &Db{
		Snowflake:  snowflake.NewGen(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		Pool:       pool,
		Sugar:      sugar,
		GuildCache: guildCache,
		Timeout:    10 * time.Second,
		Searcher:   pg.New(pool, Debug),
	}

	return
}

//go:embed migrations
var fs embed.FS

func runMigrations(url string, sugar *zap.SugaredLogger) (err error) {
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
		sugar.Infof("Performed %v migrations!", n)
	}

	err = db.Close()
	return err
}

// Time gets the time from a snowflake
func (db *Db) Time(s snowflake.ID) time.Time {
	t, _ := db.Snowflake.Parse(s)
	return t
}

// Context is a convenience method to get a context.Context with the database's timeout
func (db *Db) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), db.Timeout)
}
