package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

// ContributorCategory is a category of contributors, optionally with a role.
type ContributorCategory struct {
	ID     int64
	Name   string
	RoleID *discord.RoleID
}

// Contributor is a single contributor.
type Contributor struct {
	UserID   discord.UserID
	Category int64
	Name     string
	Override *string
}

// AddContributorCategory adds a contributor category.
func (db *DB) AddContributorCategory(name string, roleID *discord.RoleID) (cat ContributorCategory, err error) {
	err = pgxscan.Get(context.Background(), db, &cat, "insert into contributor_categories (name, role_id) values ($1, $2) returning *", name, roleID)
	return
}

// CategoryFromRole gets a contributor category from a role ID.
func (db *DB) CategoryFromRole(id discord.RoleID) *ContributorCategory {
	var c ContributorCategory
	err := pgxscan.Get(context.Background(), db, &c, "select * from contributor_categories where role_id = $1", id)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			db.Sugar.Errorf("Error getting contributor category: %v", err)
		}
		return nil
	}
	return &c
}

// ContributorCategories ...
func (db *DB) ContributorCategories() ([]ContributorCategory, error) {
	var cats []ContributorCategory

	err := pgxscan.Select(context.Background(), db, &cats, "select * from contributor_categories order by id")
	return cats, err
}

// ContributorCategory ...
func (db *DB) ContributorCategory(name string) *ContributorCategory {
	var c ContributorCategory
	err := pgxscan.Get(context.Background(), db, &c, "select * from contributor_categories where $1 ilike name", name)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			db.Sugar.Errorf("Error getting contributor category: %v", err)
		}
		return nil
	}
	return &c
}

// AddContributor adds a contributor.
func (db *DB) AddContributor(cat int64, userID discord.UserID, name string) (err error) {
	_, err = db.Exec(context.Background(), "insert into contributors (user_id, category, name) values ($1, $2, $3) on conflict do nothing", userID, cat, name)
	return err
}

// UpdateContributorName updates the contributor's Discord name.
func (db *DB) UpdateContributorName(userID discord.UserID, newName string) (err error) {
	_, err = db.Exec(context.Background(), "update contributors set name = $1 where user_id = $2", newName, userID)
	return err
}

// OverrideContributorName overrides the contributor's name.
func (db *DB) OverrideContributorName(userID discord.UserID, override *string) (err error) {
	_, err = db.Exec(context.Background(), "update contributors set override = $1 where user_id = $2", override, userID)
	return
}

// Contributors ...
func (db *DB) Contributors(id int64) ([]Contributor, error) {
	var c []Contributor
	err := pgxscan.Select(context.Background(), db, &c, "select * from contributors where category = $1 order by override nulls last, name, user_id", id)
	return c, err
}
