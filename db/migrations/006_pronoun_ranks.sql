-- +migrate Up notransaction

alter table pronouns add column uses bigint not null default 0;

create index concurrently pronouns_uses_idx on pronouns (uses);
