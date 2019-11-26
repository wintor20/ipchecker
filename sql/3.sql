\c checker_db checker

BEGIN;

CREATE OR REPLACE FUNCTION rand(h INT)
    RETURNS INT AS
$$
SELECT (floor(random()* (h + 1)))::INT
$$ LANGUAGE sql;

-- mock data
INSERT INTO conn_log (user_id, ip_addr, ts)
SELECT
       rand(10000),
       rand(255)||'.'||rand(255)||'.'||rand(255)||'.'||rand(255),
       now() + (rand(100)::varchar||'h')::interval
FROM generate_series(1, 50);

COMMIT;
