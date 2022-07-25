--liquibase formatted sql
--changeset <system>:create_candle_minutes
CREATE TABLE candle_minutes (
		price_close NUMERIC,
		price_high NUMERIC,
		price_low NUMERIC,
		price_open NUMERIC,
    		product_id TEXT,
    		unix BIGINT,
		volume NUMERIC,
    PRIMARY KEY(product_id, unix)
);

--rollback DROP TABLE
--rollback candle_minutes
