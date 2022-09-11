--liquibase formatted sql
--changeset <system>:create_tests
CREATE TABLE tests (
	id VARCHAR(255) NOT NULL,
	test_string VARCHAR(255) NOT NULL,
	PRIMARY KEY (id)
);
--rollback DROP TABLE
--rollback tests;
