package tools

import "testing"

func TestGateway(t *testing.T) {
	t.Skip("These methods are unused")
	mongoDBContainer := "docker-mongo-1"
	gateway, err := gateway(mongoDBContainer)
	if err != nil {
		t.Fatalf("failed to get gateway IP address: %v", err)
	}
	t.Logf("gateway IP address: %s", gateway)
}
