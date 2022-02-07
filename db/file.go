package db

import (
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/snowflake/v2"
)

// File is a single file
type File struct {
	url string

	ID snowflake.ID

	Filename    string
	ContentType string

	Source      string
	Description string

	Data []byte
}

// URL ...
func (f File) URL() string {
	return f.url
}

// AddFile adds a file
func (db *DB) AddFile(filename, contentType string, data []byte) (f *File, err error) {
	f = &File{}

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Adding file with name %v, content type %v, data length %v", filename, contentType, len(data))

	err = pgxscan.Get(ctx, db.Pool, f, "insert into files (id, filename, content_type, data) values ($1, $2, $3, $4) returning *", db.Snowflake.Get(), filename, contentType, data)
	if err != nil {
		return nil, err
	}

	if db.Config.Bot.Website != "" {
		f.url = fmt.Sprintf("%vfile/%v/%v", db.Config.Bot.Website, f.ID, f.Filename)
	}

	return f, err
}

// File gets a file from the database
func (db *DB) File(id snowflake.ID) (f File, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting file with ID %v", id)

	err = pgxscan.Get(ctx, db.Pool, &f, "select * from files where id = $1", id)
	if err != nil {
		return
	}

	if db.Config.Bot.Website != "" {
		f.url = fmt.Sprintf("%vfile/%v/%v", db.Config.Bot.Website, f.ID, f.Filename)
	}

	return
}

// Files gets all files
func (db *DB) Files() (f []File, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting all files")

	err = pgxscan.Select(ctx, db.Pool, &f, "select id, filename, content_type, source, description from files order by filename asc")
	return
}

// FileName returns files with the given string in their name
func (db *DB) FileName(s string) (f []File, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting files containing %v", s)

	err = pgxscan.Select(ctx, db.Pool, &f, "select id, filename, content_type, source, description from files where position(lower($1) in lower(filename)) > 0 order by filename asc", s)
	return
}
