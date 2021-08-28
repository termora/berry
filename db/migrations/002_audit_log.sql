-- +migrate Up

-- 2021-08-05: detailed audit log
-- logs all actions done by term directors/admins

create type audit_log_entry_subject as enum ('term', 'pronouns', 'explanation');
create type audit_log_action as enum ('create', 'update', 'delete');

create table audit_log (
    id  serial  primary key,

    -- possibly 0 if the message didn't send correctly, or no log is set.
    private_message_id  bigint  not null    default 0,
    public_message_id   bigint  not null    default 0,

    subject_id  int                     not null,
    subject     audit_log_entry_subject not null,
    action      audit_log_action        not null,

    before  jsonb,
    after   jsonb,

    user_id bigint  not null,
    reason  text,

    timestamp   timestamp   not null default (current_timestamp at time zone 'utc')
);
