BEGIN;

CREATE TABLE conn_log (
    user_id bigint,
    ip_addr varchar(15),
    ts timestamp
);

COMMIT;
