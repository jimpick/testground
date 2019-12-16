package main

import (
	// "context"
	"fmt"
	"net"
	"os"
	"reflect"
	// "time"

	"github.com/ipfs/testground/sdk/runtime"
	"github.com/ipfs/testground/sdk/sync"
)

const (
	// sizeBytes = 10 * 1024 * 1024 // 10mb
	sizeBytes = 100 * 1024 * 1024 // 100mb
)

var peerIPSubtree = &sync.Subtree{
	GroupKey:    "peerIPs",
	PayloadType: reflect.TypeOf(&net.IP{}),
	KeyFunc: func(val interface{}) string {
		return val.(*net.IP).String()
	},
}

/*
var cidSubtree = &sync.Subtree{
	GroupKey:    "cid",
	PayloadType: reflect.TypeOf(&cid.Cid{}),
	KeyFunc: func(val interface{}) string {
		return "cid"
	},
}
*/

func main() {
	runenv := runtime.CurrentRunEnv()

	ifaddrs, err := net.InterfaceAddrs()
	if err != nil {
		runenv.Abort(err)
		return
	}
	fmt.Fprintln(os.Stderr, "Addrs:", ifaddrs)

	_, localnet, _ := net.ParseCIDR("8.0.0.0/8")

	var peerIP net.IP
        for _, ifaddr := range ifaddrs {
                var ip net.IP
                switch v := ifaddr.(type) {
                case *net.IPNet:
                        ip = v.IP
                case *net.IPAddr:
                        ip = v.IP
                }
                fmt.Fprintln(os.Stderr, "IP:", ip)
                if localnet.Contains(ip) {
                        peerIP = ip
                        break
                }
        }
        fmt.Fprintln(os.Stderr, "Matched IP:", peerIP)
        if peerIP == nil {
		runenv.Abort("No IP match")
		return
        }

	/*
	timeout := func() time.Duration {
		if t, ok := runenv.IntParam("timeout_secs"); !ok {
			return 30 * time.Second
		} else {
			return time.Duration(t) * time.Second
		}
	}()
	*/

	watcher, writer := sync.MustWatcherWriter(runenv)
	defer watcher.Close()

	// ctx := context.Background()

	seq, err := writer.Write(peerIPSubtree, &peerIP)
	if err != nil {
		runenv.Abort(err)
		return
	}
	defer writer.Close()

	fmt.Fprintln(os.Stderr, "Jim1 seq", seq)

	/*
	// States
	added := sync.State("added")
	received := sync.State("received")

	switch {
	case seq == 1: // adder
		adder := client
		fmt.Fprintln(os.Stderr, "adder")

		// generate a random file of the designated size.
		file := utils.TempRandFile(runenv, ensemble.TempDir(), sizeBytes)
		defer os.Remove(file.Name())

		tstarted := time.Now()
		hash, err := adder.Add(file)
		if err != nil {
			runenv.Abort(err)
			return
		}
		fmt.Fprintln(os.Stderr, "cid", hash)

		c, err := cid.Parse(hash)
		if err != nil {
			runenv.Abort(err)
			return
		}

		runenv.EmitMetric(utils.MetricTimeToAdd, float64(time.Now().Sub(tstarted)/time.Millisecond))

		_, err = writer.Write(cidSubtree, &c)
		if err != nil {
			runenv.Abort(err)
		}

		// Signal we're done on the added state.
		_, err = writer.SignalEntry(added)
		if err != nil {
			runenv.Abort(err)
		}
		fmt.Fprintln(os.Stderr, "State: added")

		// Set a state barrier.
		receivedCh := watcher.Barrier(ctx, received, 1)

		// Wait until recieved state is signalled.
		if err := <-receivedCh; err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, "State: received")
		runenv.OK()
	case seq == 2: // getter
		getter := client
		fmt.Fprintln(os.Stderr, "getter")

		// Connect to other peers
		peerCh := make(chan *peer.AddrInfo, 16)
		cancel, err := watcher.Subscribe(sync.PeerSubtree, peerCh)
		if err != nil {
			runenv.Abort(err)
		}
		defer cancel()

		var events int
		for i := 0; i < runenv.TestInstanceCount; i++ {
			select {
			case ai := <-peerCh:
				events++
				if ai.ID == ID {
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				tstarted := time.Now()
				err = getter.SwarmConnect(ctx, ai.Addrs[0].String())
				if err != nil {
					runenv.Abort(err)
					return
				}

				runenv.EmitMetric(utils.MetricTimeToConnect, float64(time.Now().Sub(tstarted)/time.Millisecond))
				cancel()

			case <-time.After(timeout):
				// TODO need a way to fail a distributed test immediately. No point
				// making it run elsewhere beyond this point.
				panic(fmt.Sprintf("no new peers in %d seconds", timeout))
			}
		}

		// Set a state barrier.
		addedCh := watcher.Barrier(ctx, added, 1)

		// Wait until added state is signalled.
		if err := <-addedCh; err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, "State: added")

		cidCh := make(chan *cid.Cid, 0)
		cancel, err = watcher.Subscribe(cidSubtree, cidCh)
		if err != nil {
			runenv.Abort(err)
		}
		defer cancel()
		select {
		case c := <-cidCh:
			cancel()
			// Get the content from the adder node
			tstarted := time.Now()
			err = getter.Get(c.String(), ensemble.TempDir())
			if err != nil {
				runenv.Abort(err)
				return
			}
			runenv.EmitMetric(utils.MetricTimeToGet, float64(time.Now().Sub(tstarted)/time.Millisecond))
		case <-time.After(timeout):
			// TODO need a way to fail a distributed test immediately. No point
			// making it run elsewhere beyond this point.
			panic(fmt.Sprintf("no cid in %d seconds", timeout))
		}

		// Signal we're reached the received state.
		_, err = writer.SignalEntry(received)
		if err != nil {
			runenv.Abort(err)
		}
		fmt.Fprintln(os.Stderr, "State: received")
		runenv.OK()
	default:
		runenv.Abort(fmt.Errorf("Unexpected seq: %v", seq))
	}

	ensemble.Destroy()
	*/
	runenv.OK()
}
