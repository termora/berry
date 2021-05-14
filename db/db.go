package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/structs"
	"go.uber.org/zap"
)

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

	sentry    *sentry.Hub
	useSentry bool

	Searcher
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
		Pool:       pool,
		Sugar:      sugar,
		GuildCache: guildCache,

		Searcher: NewPsqlSearcher(pool),
	}

	return
}

func initDB(s *zap.SugaredLogger, url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), url)
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
	err := db.QueryRow(context.Background(), "select exists (select from information_schema.tables where table_schema = 'public' and table_name = 'info')").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // the database has been initialised so we're done
	}

	// ...it's not initialised and we have to do that
	_, err = db.Exec(context.Background(), initDBSql)
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
	err = db.QueryRow(context.Background(), "select schema_version from public.info").Scan(&dbVersion)
	if err != nil {
		return err
	}
	initialDBVersion := dbVersion
	for dbVersion < DBVersion {
		_, err = db.Exec(context.Background(), DBVersions[dbVersion-1])
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
