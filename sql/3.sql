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
FROM generate_series(1,800000);

INSERT INTO conn_log (user_id, ip_addr, ts)
VALUES
(1, '1.1.1.1', now()),
(1, '2.2.2.2', now()),
(2, '1.1.1.1', now()),
(2, '2.2.2.2', now()),
(3, '1.1.1.1', now()),
(3, '2.2.2.2', now()),

(4, '1.1.1.17', now()),
(4, '2.2.2.17', now()),
(5, '1.1.1.17', now()),
(5, '2.2.2.17', now())

ON CONFLICT DO NOTHING;


COMMIT;
