--liquibase formatted sql
--changeset <system>:create_accounts
CREATE TABLE accounts (
		available DECIMAL(20,8) NOT NULL,
		balance DECIMAL(20,8) NOT NULL,
		hold DECIMAL(20,8) NOT NULL,
		id VARCHAR(255) NOT NULL,
		currency VARCHAR(255) NOT NULL,
		profile_id VARCHAR(255) NOT NULL,
		trading_enabled BOOLEAN NOT NULL,
		PRIMARY KEY (id)
);
--rollback DROP TABLE
--rollback accounts;
