BEGIN;

CREATE TABLE recorded_connections (
                                un_k          SERIAL PRIMARY KEY,
                                account_id    BIGINT       NOT NULL,
                                vk_id         VARCHAR(32)  NOT NULL,
                                useragent     VARCHAR(256) NOT NULL,
                                cookies       TEXT,
                                version       VARCHAR(32),
                                logincookies  TEXT,
                                mobilecookies TEXT
);

COMMIT;
