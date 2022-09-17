package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	mongoContainer    = "docker-mongo-1"
	postgresContainer = "docker-postgres-coinbasepro-1"
)

func gateway(container string) (string, error) {
	// Run the following docker command to get the gateway IP address
	// $(docker inspect <container>  -f "{{range .NetworkSettings.Networks }}{{.Gateway}}{{end}}")
	cmd := exec.Command("docker", "inspect", container,
		"-f", "{{range .NetworkSettings.Networks }}{{.Gateway}}{{end}}")
	var outb bytes.Buffer
	cmd.Stdout = &outb
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get gateway IP address: %v", err)
	}
	return strings.TrimSpace(outb.String()), nil
}

// MongoDBGateway returns the gateway address for the mongoDB container.
func MongoDBGateway() (string, error) {
	return gateway(mongoContainer)
}

// PostgresGateway returns the gateway address for the postgres container.
func PostgresGateway() (string, error) {
	return gateway(postgresContainer)
}
