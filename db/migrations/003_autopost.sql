-- +migrate Up

-- 2021-08-28: autoposting
-- automatically posts terms to a designated channel every [x] hours

create table autopost (
    -- stored as text because the servers table stores it as text
    guild_id    text    not null    references servers (id) on delete cascade,
    channel_id  bigint  primary key,

    next_post   timestamp   not null,
    interval    interval    not null,
    category_id int         references categories (id) on delete set null
);

create index autopost_guild_id_idx on autopost (guild_id);
