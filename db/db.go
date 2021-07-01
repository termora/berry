package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/starshine-sys/snowflake/v2"
	"github.com/termora/berry/structs"
	"go.uber.org/zap"
)

// Debug is a debug logging function
var Debug = func(template string, args ...interface{}) {}

var termCounter struct {
	count uint64
	mu    sync.RWMutex
}

// AddCount adds one to the term fetch count
func AddCount() uint64 {
	termCounter.mu.Lock()
	defer termCounter.mu.Unlock()
	termCounter.count++
	return termCounter.count
}

// GetCount ...
func GetCount() uint64 {
	termCounter.mu.RLock()
	defer termCounter.mu.RUnlock()
	return termCounter.count
}

// Db ...
type Db struct {
	Pool       *pgxpool.Pool
	Sugar      *zap.SugaredLogger
	GuildCache *ttlcache.Cache

	Config *structs.BotConfig

	Snowflake *snowflake.Generator

	Timeout time.Duration

	sentry    *sentry.Hub
	useSentry bool
}

// Init ...
func Init(url string, sugar *zap.SugaredLogger) (db *Db, err error) {
	guildCache := ttlcache.NewCache()
	guildCache.SetCacheSizeLimit(100)
	guildCache.SetTTL(10 * time.Minute)

	pool, err := initDB(sugar, url)
	if err != nil {
		return nil, err
	}

	db = &Db{
		Snowflake:  snowflake.NewGen(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		Pool:       pool,
		Sugar:      sugar,
		GuildCache: guildCache,
		Timeout:    10 * time.Second,
	}

	return
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

func initDB(s *zap.SugaredLogger, url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	if err := initDBIfNotInitialised(s, db); err != nil {
		return nil, err
	}
	err = updateDB(s, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDBIfNotInitialised(s *zap.SugaredLogger, db *pgxpool.Pool) error {
	var exists bool

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.QueryRow(ctx, "select exists (select from information_schema.tables where table_schema = 'public' and table_name = 'info')").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // the database has been initialised so we're done
	}

	// ...it's not initialised and we have to do that
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.Exec(ctx, initDBSql)
	if err != nil {
		return err
	}
	if s != nil {
		s.Infof("Successfully initialised the database.")
	}
	return nil
}

func updateDB(s *zap.SugaredLogger, db *pgxpool.Pool) (err error) {
	var dbVersion int

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.QueryRow(ctx, "select schema_version from public.info").Scan(&dbVersion)
	if err != nil {
		return err
	}
	initialDBVersion := dbVersion
	for dbVersion < DBVersion {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = db.Exec(ctx, DBVersions[dbVersion-1])
		if err != nil {
			return err
		}
		dbVersion++
		if s != nil {
			s.Infof("Updated database to version %v", dbVersion)
		}
	}
	if initialDBVersion < DBVersion && s != nil {
		s.Infof("Successfully updated database to target version %v", DBVersion)
	}
	return nil
}
