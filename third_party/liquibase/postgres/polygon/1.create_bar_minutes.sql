--liquibase formatted sql
--changeset <system>:create_bar_minutes
CREATE TABLE bar_minutes (
		ticker VARCHAR(50),
		adjusted BOOLEAN,
		c NUMERIC,
		h NUMERIC,
		l NUMERIC,
		n NUMERIC,
		o NUMERIC,
		t NUMERIC,
		v NUMERIC,
		vw NUMERIC,
    PRIMARY KEY(ticker, adjusted, t)
);

--rollback DROP TABLE
--rollback bar_minutes
