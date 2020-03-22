-- +migrate Up
create table events (
    id uuid primary key,
    heading text,
    start_date date not null,
    start_time time,
    end_date date not null,
    end_time time,
    descr text,
    owner text
);
create index start_idx on events using btree (start_date, start_time);

-- +migrate Down
DROP TABLE events;