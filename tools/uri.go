package tools

import "fmt"

// MongoURI will return the full URI for accessing a mongo database.
func MongoURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("mongodb://%s%s:%s/%s", auth, host, port, database), nil
}

// PostgressURI will return the full URI for accessing a postgres database.
func PostgresURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" || username != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("postgresql://%s%s:%s/%s?sslmode=disable", auth, host, port, database), nil
}

// RedisURI will return the full URI for accessing a redis cache.
func RedisURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("redis://%s%s:%s/%s", auth, host, port, database), nil
}
