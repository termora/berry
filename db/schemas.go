package db

// DBVersion is the current database version
const DBVersion = 2

// DBVersions is a slice of schemas for every database version
var DBVersions []string = []string{
	`alter table public.terms add column flags integer not null default 0;
    update public.info set schema_version = 2;`,
}

// initDBSql is the initial SQL database schema
var initDBSql = `create table if not exists admins (
    user_id     text    primary key
);

create table if not exists categories (
    id      serial  primary key,
    name    text    not null
);

create table if not exists terms (
    id          serial      primary key,
    category    int         not null references categories (id) on delete cascade,
    name        text        not null,
    aliases     text[]      not null default array[]::text[],
    description text        not null,
	created     timestamp   not null default (current_timestamp at time zone 'utc'),
	source      text        not null default 'Unknown',
	searchtext  tsvector    generated always as (
        to_tsvector('english', "name") ||
        to_tsvector('english', "description") ||
        to_tsvector('english', "source") ||
        array_to_tsvector("aliases")
    ) stored
);

create table if not exists explanations (
    id          serial      primary key,
    name        text        not null,
    aliases     text[]      not null default array[]::text[],
    description text        not null,
    created     timestamp   not null default (current_timestamp at time zone 'utc')
);

create table if not exists servers
(
    id            text      primary key,
    blacklist     text[]    not null default array[]::text[]
);

create table if not exists info
(
    id                      int primary key not null default 1, -- enforced only equal to 1
    schema_version          int,
    constraint singleton    check (id = 1) -- enforce singleton table/row
);

insert into info (schema_version) values (1);`
