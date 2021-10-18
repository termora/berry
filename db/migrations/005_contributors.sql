-- +migrate Up

-- 2021-09-14: store contributors for the t;credits command

create table contributor_categories (
    id      serial  primary key,
    name    text    not null,
    role_id bigint  unique
);

create table contributors (
    user_id     bigint  not null,
    category    int     not null    references contributor_categories (id) on delete cascade,
    name        text    not null, -- updated by bot on member update event
    override    text,             -- possibly null, `name` is used if it's null

    primary key (user_id, category)
);
