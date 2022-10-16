CREATE TABLE tests1 (
	id VARCHAR(255) NOT NULL,
	test_string VARCHAR(255) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE tests2 (
	id VARCHAR(255) NOT NULL,
	test_string VARCHAR(255) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE lttests1 (
	id VARCHAR(255) NOT NULL,
	test_string VARCHAR(255) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE pktests1 (
	test_string VARCHAR(255) NOT NULL,
	test_int INT NOT NULL,
	PRIMARY KEY (test_string)
);
