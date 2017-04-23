ALTER TABLE es.t_aepb_publish add column typecode CHARACTER VARYING(30) NOT NULL;

ALTER TABLE es.t_aepb_publish add column
payload BYTEA;

ALTER TABLE es.t_aepb_publish add column
event_time TIMESTAMP(6) WITHOUT TIME ZONE DEFAULT CLOCK_TIMESTAMP();