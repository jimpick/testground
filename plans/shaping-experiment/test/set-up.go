package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ipfs/testground/sdk/runtime"
	"github.com/ipfs/testground/sdk/sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
)

// SetUp sets up the elements necessary for the test cases
func SetUp(ctx context.Context, runenv *runtime.RunEnv, timeout time.Duration, randomWalk bool, bucketSize int, autoRefresh bool, watcher *sync.Watcher, writer *sync.Writer) (host.Host, *kaddht.IpfsDHT, []peer.AddrInfo, error) {
	// TODO: just put the hostname inside the runenv?
	hostname, err := os.Hostname()
	if err != nil {
		return nil, nil, nil, err
	}

	if runenv.TestSidecar {
		// Wait for the network to be ready.
		//
		// Technically, we don't need to do this as configuring the network will
		// block on it being ready.
		err := <-watcher.Barrier(ctx, "network-initialized", int64(runenv.TestInstanceCount))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to initialize network: %w", err)
		}

		writer.Write(sync.NetworkSubtree(hostname), &sync.NetworkConfig{
			Network: "default",
			Enable:  true,
			Default: sync.LinkShape{
				Latency:   100 * time.Millisecond,
				Bandwidth: 1 << 20, // 1Mib
			},
			State: "network-configured",
		})

		err = <-watcher.Barrier(ctx, "network-configured", int64(runenv.TestInstanceCount))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to configure network: %w", err)
		}
	}

	/// --- Set up

	node, dht, err := CreateDhtNode(ctx, runenv, bucketSize, autoRefresh)
	if err != nil {
		return nil, nil, nil, err
	}

	/// --- Warm up

	myNodeID := node.ID()
	runenv.Message("I am %s with addrs: %v", myNodeID, node.Addrs())

	if _, err = writer.Write(sync.PeerSubtree, host.InfoFromHost(node)); err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to get Redis Sync PeerSubtree %w", err)
	}

	// TODO: Revisit this - This assumed that it is ok to put in memory every single peer.AddrInfo that participates in this test
	peerCh := make(chan *peer.AddrInfo, 16)
	cancelSub, err := watcher.Subscribe(sync.PeerSubtree, peerCh)
	defer cancelSub()

	var toDial []peer.AddrInfo
	// Grab list of other peers that are available for this Run
	for i := 0; i < runenv.TestInstanceCount; i++ {
		select {
		case ai := <-peerCh:
			id1, _ := ai.ID.MarshalBinary()
			id2, _ := myNodeID.MarshalBinary()
			if bytes.Compare(id1, id2) >= 0 {
				// skip over dialing ourselves, and prevent TCP simultaneous
				// connect (known to fail) by only dialing peers whose peer ID
				// is smaller than ours.
				continue
			}
			toDial = append(toDial, *ai)

		case <-time.After(timeout):
			return nil, nil, nil, fmt.Errorf("no new peers in %d seconds", timeout/time.Second)
		}
	}

	// Dial to all the other peers
	for _, ai := range toDial {
		err = node.Connect(ctx, ai)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Error while dialing peer %v: %w", ai.Addrs, err)
		}
	}

	runenv.Message("Dialed all my buds")

	// Check if `random-walk` is enabled, if yes, run it 5 times
	for i := 0; randomWalk && i < 5; i++ {
		err = dht.Bootstrap(ctx)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Could not run a random-walk: %w", err)
		}
	}

Loop:
	for {
		select {
		case <-time.After(200 * time.Millisecond):
			if dht.RoutingTable().Size() > 0 {
				break Loop
			}
		case <-ctx.Done():
			return nil, nil, nil, fmt.Errorf("got no peers in routing table")
		}
	}

	return node, dht, toDial, nil
}
