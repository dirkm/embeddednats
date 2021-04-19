package embeddednats

import (
	"fmt"
	"os"
	"sync"
	"testing"

	natsstreaming "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nats.go"
	"github.com/phayes/freeport"
)

const (
	MessageBusHostEnv      = "MESSAGE_BUS_HOST"
	MessageBusClusterIDEnv = "MESSAGE_BUS_CLUSTER_ID"
)

var streamingInitialized = sync.Once{}
var initNatsPort = -1
var initClusterListenPort = -1

func InitNATSStreaming(t testing.TB) (natsPort int, clusterListenPort int) {
	streamingInitialized.Do(func() { initNatsPort, initClusterListenPort = StartNATSStreaming(t) })
	return initNatsPort, initClusterListenPort
}

func StartNATSStreaming(t testing.TB) (natsPort int, clusterListenPort int) {
	freeports, err := freeport.GetFreePorts(3)
	if err != nil {
		t.Fatalf("could not create message bus factory %v", err)
	}
	natsPort, clusterListenPort, routePort := freeports[0], freeports[1], freeports[2]
	// Start a streaming server, and setup a route
	nOpts := natsstreaming.DefaultNatsServerOptions
	nOpts.Host = "0.0.0.0"
	nOpts.Port = natsPort
	nOpts.Cluster.ListenStr = fmt.Sprintf("nats://0.0.0.0:%d", clusterListenPort)
	nOpts.RoutesStr = fmt.Sprintf("nats://0.0.0.0:%d", routePort)

	sOpts := natsstreaming.GetDefaultOptions()
	sOpts.ID = "test-cluster"
	sOpts.Debug = true
	sOpts.NATSClientOpts = append(sOpts.NATSClientOpts, nats.UseOldRequestStyle())

	_, err = natsstreaming.RunServerWithOpts(sOpts, &nOpts)
	if err != nil {
		t.Fatalf("could not create server %v", err)
	}

	if err = os.Setenv(MessageBusHostEnv, fmt.Sprintf("0.0.0.0:%d", natsPort)); err != nil {
		t.Fatalf("could not set variable %v", err)
	}

	if err = os.Setenv(MessageBusClusterIDEnv, sOpts.ID); err != nil {
		t.Fatalf("could not set variable %v", err)
	}
	return
}
