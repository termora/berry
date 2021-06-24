package db

// DBVersion is the current database version
const DBVersion = 20

// DBVersions is a slice of schemas for every database version
var DBVersions []string = []string{
	`alter table public.terms add column flags integer not null default 0;
    update public.info set schema_version = 2;`,

	`alter table public.terms drop column searchtext;
    alter table public.terms add column searchtext tsvector generated always as (
        setweight(to_tsvector('english', "name"), 'A') ||
        setweight(to_tsvector('english', "description"), 'B') ||
        setweight(to_tsvector('english', "source"), 'C') ||
        setweight(array_to_tsvector("aliases"), 'A')
    ) stored;
    update public.info set schema_version = 3;`,

	`alter table public.terms add column content_warnings text not null default '';
    update public.info set schema_version = 4;`,

	`create index term_names_alphabetical on public.terms (name, id);
    update public.info set schema_version = 5;`,

	`alter table public.terms add column last_modified timestamp;
    update public.terms set last_modified = created where last_modified is null;
    alter table public.terms alter column last_modified set default (current_timestamp at time zone 'utc');
    alter table public.terms alter column last_modified set not null;
    update public.info set schema_version = 6;`,

	`create table if not exists admin_tokens (
        user_id     text        primary key,
        token       text        not null,
        expires     timestamp   not null default (now() + interval '30 days')::timestamp
    );
    update public.info set schema_version = 7;`,

	`alter table public.terms add column note text not null default '';
    update public.info set schema_version = 8;`,

	`alter table public.explanations add column as_command boolean not null default false;
    update public.info set schema_version = 9;`,

	`alter table public.terms add column aliases_string text not null default '';
    
    alter table public.terms drop column searchtext;
    alter table public.terms add column searchtext tsvector generated always as (
        setweight(to_tsvector('english', "name"), 'A') ||
        setweight(to_tsvector('english', "description"), 'B') ||
        setweight(to_tsvector('english', "source"), 'C') ||
        setweight(to_tsvector('english', "aliases_string"), 'A')
    ) stored;

    update public.terms set aliases_string = array_to_string(aliases, ', ');

    update public.info set schema_version = 10;`,

	`create table if not exists errors (
        id      uuid        primary key,
        command text        not null,
        user_id bigint      not null,
        channel bigint      not null,
        error   text        not null,
        time    timestamp   not null default (current_timestamp at time zone 'utc')
    );
    
    update public.info set schema_version = 11;`,

	`alter table public.terms add column image_url text not null default '';
    
    update public.info set schema_version = 12;`,

	`create table if not exists pronouns (
        id          serial  primary key,
        subjective  text    not null default '',
        objective   text    not null default '',
        poss_det    text    not null default '',
        poss_pro    text    not null default '',
        reflexive   text    not null default '',

        unique (subjective, objective, poss_det, poss_pro, reflexive)
    );
    
    update public.info set schema_version = 13;`,

	`alter table public.terms add column tags text[] not null default array[]::text[];
    
    update public.info set schema_version = 14;`,

	`create index subjective_idx on public.pronouns (lower(subjective));
    create index objective_idx on public.pronouns (lower(subjective), lower(objective));
    create index poss_det_idx on public.pronouns (lower(subjective), lower(objective), lower(poss_det));
    create index poss_pro_idx on public.pronouns (lower(subjective), lower(objective), lower(poss_det), lower(poss_pro));
    create index reflexive_idx on public.pronouns (lower(subjective), lower(objective), lower(poss_det), lower(poss_pro), lower(reflexive));
    
    update public.info set schema_version = 15;`,

	`alter table public.servers add column prefixes text[] not null default array[]::text[];
    -- set default prefixes
    update public.servers set prefixes = array['t;', 't:'] where prefixes = array[]::text[];
    
    update public.info set schema_version = 16;`,

	`create table if not exists pronoun_msgs (
        message_id  bigint  primary key,
        subjective  text    not null default '',
        objective   text    not null default '',
        poss_det    text    not null default '',
        poss_pro    text    not null default '',
        reflexive   text    not null default ''
    );
    
    update public.info set schema_version = 17;`,

	`alter table public.pronouns add column sorting int not null default 5;
    update public.info set schema_version = 18;`,

	`create table if not exists tags (
        normalized  text    primary key,
        display     text    not null
    );

    update terms set tags = lower(tags::text)::text[];

    update public.info set schema_version = 19;`,

	`create table if not exists files (
        id  bigint  primary key,
    
        filename        text    not null,
        content_type    text    not null,
    
        source      text    not null default '',
        description text    not null default '',
    
        data    bytea   not null
    );

    alter table public.terms add column files bigint[] not null default array[]::bigint[];
    
    update public.info set schema_version = 20;`,
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
        setweight(to_tsvector('english', "name"), 'A') ||
        setweight(to_tsvector('english', "description"), 'B') ||
        setweight(to_tsvector('english', "source"), 'C') ||
        setweight(array_to_tsvector("aliases"), 'A')
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
