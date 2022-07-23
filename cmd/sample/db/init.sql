CREATE SCHEMA IF NOT EXISTS sample;

CREATE TABLE IF NOT EXISTS sample.user (
    id bigserial primary key,
    added_at timestamp NOT NULL default clock_timestamp(),
    name text NOT NULL
);