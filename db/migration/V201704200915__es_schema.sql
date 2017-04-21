CREATE SCHEMA IF NOT EXISTS es;

CREATE TABLE IF NOT EXISTS es.t_aeev_events(
id bigserial,
event_time TIMESTAMP(6) WITHOUT TIME ZONE DEFAULT CLOCK_TIMESTAMP(),
aggregate_id CHARACTER VARYING(60) NOT NULL,
version NUMERIC(38,0) NOT NULL,
typecode CHARACTER VARYING(30) NOT NULL,
payload BYTEA,
primary key(aggregate_id, version)
)
WITH (
OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS es.t_aepb_publish(
aggregate_id CHARACTER VARYING(60) NOT NULL,
version NUMERIC(38,0) NOT NULL,
primary key(aggregate_id, version)
)
WITH (
OIDS=FALSE
);