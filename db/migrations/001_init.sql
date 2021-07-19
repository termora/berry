-- +migrate Up

-- 2021-07-19: Initial database schema
-- Everything here should be idempotent assuming a clean database or an existing database using the old migrations in schemas.go

-- drop the old info table if it exists
drop table if exists info;

create table if not exists admins (
    user_id     text    primary key
);

create table if not exists categories (
    id      serial  primary key,
    name    text    not null
);

create table if not exists tags (
    normalized  text    primary key,
    display     text    not null
);

create table if not exists terms (
    id          serial      primary key,
    category    int         not null references categories (id) on delete cascade,
    name        text        not null,
    aliases     text[]      not null default array[]::text[],
    tags        text[]      not null default array[]::text[],
    description text        not null,
    source      text        not null default 'Unknown',
    image_url   text        not null default '',
    
    content_warnings    text        not null default '',
    note                text        not null default '',
    files               bigint[]    not null default array[]::bigint[],

    created         timestamp   not null default (current_timestamp at time zone 'utc'),
    last_modified   timestamp   not null default (current_timestamp at time zone 'utc'),

    flags           integer not null default 0,
    -- set when updating, because arrays aren't converted to tsvector correctly
    aliases_string  text    not null default '',

    searchtext  tsvector    generated always as (
        setweight(to_tsvector('english', "name"), 'A') ||
        setweight(to_tsvector('english', "description"), 'B') ||
        setweight(to_tsvector('english', "source"), 'C') ||
        setweight(to_tsvector('english', "aliases_string"), 'A')
    ) stored
);

create index if not exists term_names_alphabetical on terms (name, id);

create table if not exists pronouns (
    id          serial  primary key,
    subjective  text    not null default '',
    objective   text    not null default '',
    poss_det    text    not null default '',
    poss_pro    text    not null default '',
    reflexive   text    not null default '',

    sorting int not null default 5,

    unique (subjective, objective, poss_det, poss_pro, reflexive)
);

create index if not exists subjective_idx on pronouns (lower(subjective));
create index if not exists objective_idx on pronouns (lower(subjective), lower(objective));
create index if not exists poss_det_idx on pronouns (lower(subjective), lower(objective), lower(poss_det));
create index if not exists poss_pro_idx on pronouns (lower(subjective), lower(objective), lower(poss_det), lower(poss_pro));
create index if not exists reflexive_idx on pronouns (lower(subjective), lower(objective), lower(poss_det), lower(poss_pro), lower(reflexive));

create table if not exists explanations (
    id          serial      primary key,
    name        text        not null,
    aliases     text[]      not null default array[]::text[],
    description text        not null,
    created     timestamp   not null default (current_timestamp at time zone 'utc'),
    as_command  boolean     not null default false
);

create table if not exists servers
(
    id            text      primary key,
    blacklist     text[]    not null default array[]::text[],
    prefixes      text[]    not null default array[]::text[]
);

create table if not exists errors (
    id      uuid        primary key,
    command text        not null,
    user_id bigint      not null,
    channel bigint      not null,
    error   text        not null,
    time    timestamp   not null default (current_timestamp at time zone 'utc')
);

create table if not exists pronoun_msgs (
    message_id  bigint  primary key,
    subjective  text    not null default '',
    objective   text    not null default '',
    poss_det    text    not null default '',
    poss_pro    text    not null default '',
    reflexive   text    not null default ''
);

create table if not exists files (
    id  bigint  primary key,

    filename        text    not null,
    content_type    text    not null,

    source      text    not null default '',
    description text    not null default '',

    data    bytea   not null
);

-- other cleanup from the old migrations
drop table if exists admin_tokens;
