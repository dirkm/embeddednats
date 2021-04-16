package embeddednats_test

import (
	"testing"
	"time"

	"github.com/dirkm/embeddednats"
)

func TestNATSStreaming(t *testing.T) {
	natsPort, clusterListenPort := embeddednats.InitNATSStreaming(t)
	if natsPort == 0 {
		t.Fatalf("port not available %d", natsPort)
	}
	if clusterListenPort == 0 {
		t.Fatalf("cluster listen port not available %d", clusterListenPort)
	}
	time.Sleep(1 * time.Second)
}
