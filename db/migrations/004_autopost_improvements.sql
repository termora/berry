-- +migrate Up

-- 2021-09-11: improve autoposting
-- add an optional mention role to autoposts

alter table autopost add column role_id bigint;
