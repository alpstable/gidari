CREATE TABLE tests (
	id VARCHAR(255) NOT NULL,
	test_string VARCHAR(255) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE accounts (
	available DECIMAL(20, 8) NOT NULL,
	balance DECIMAL(20, 8) NOT NULL,
	hold DECIMAL(20, 8) NOT NULL,
	id VARCHAR(255) NOT NULL,
	currency VARCHAR(255) NOT NULL,
	profile_id VARCHAR(255) NOT NULL,
	trading_enabled BOOLEAN NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE btc_stats (
	id SERIAL PRIMARY KEY,
	high DECIMAL(20, 8) NOT NULL,
	last DECIMAL(20, 8) NOT NULL,
	low DECIMAL(20, 8) NOT NULL,
	open DECIMAL(20, 8) NOT NULL,
	volume DECIMAL(20, 8) NOT NULL,
	volume_30day DECIMAL(20, 8) NOT NULL
);

CREATE TABLE candle_minutes (
	price_close DECIMAL(20, 8) NOT NULL,
	price_high DECIMAL(20, 8) NOT NULL,
	price_low DECIMAL(20, 8) NOT NULL,
	price_open DECIMAL(20, 8) NOT NULL,
	product_id VARCHAR(255) NOT NULL,
	unix BIGINT NOT NULL,
	volume DECIMAL(20, 8) NOT NULL,
	PRIMARY KEY (unix, product_id)
);
