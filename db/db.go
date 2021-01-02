package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/jackc/pgx/v4/pgxpool"
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
}

// Init ...
func Init(url string, sugar *zap.SugaredLogger) (db *Db, err error) {
	guildCache := ttlcache.NewCache()
	guildCache.SetCacheSizeLimit(100)
	guildCache.SetTTL(10 * time.Minute)

	pool, err := initDB(url)
	if err != nil {
		return nil, err
	}

	db = &Db{
		Pool:       pool,
		Sugar:      sugar,
		GuildCache: guildCache,
	}

	return
}

func initDB(url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	if err := initDBIfNotInitialised(db); err != nil {
		return nil, err
	}
	err = updateDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDBIfNotInitialised(db *pgxpool.Pool) error {
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
	fmt.Printf("Successfully initialised the database.\n")
	return nil
}

func updateDB(db *pgxpool.Pool) (err error) {
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
		fmt.Printf("Updated database to version %v\n", dbVersion)
	}
	if initialDBVersion < DBVersion {
		fmt.Printf("Successfully updated database to target version %v\n", DBVersion)
	}
	return nil
}
