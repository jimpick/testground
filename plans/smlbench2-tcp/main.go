package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"

	"github.com/ipfs/testground/sdk/runtime"
	sdksync "github.com/ipfs/testground/sdk/sync"
)

const (
	// sizeBytes = 10 * 1024 * 1024 // 10mb
	sizeBytes = 100 * 1024 * 1024 // 100mb
)

var peerIPSubtree = &sdksync.Subtree{
	GroupKey:    "peerIPs",
	PayloadType: reflect.TypeOf(&net.IP{}),
	KeyFunc: func(val interface{}) string {
		return val.(*net.IP).String()
	},
}

func main() {
	runenv := runtime.CurrentRunEnv()

	ifaddrs, err := net.InterfaceAddrs()
	if err != nil {
		runenv.Abort(err)
		return
	}
	// fmt.Fprintln(os.Stderr, "Addrs:", ifaddrs)

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
		// fmt.Fprintln(os.Stderr, "IP:", ip)
		if localnet.Contains(ip) {
			peerIP = ip
			break
		}
	}
	//fmt.Fprintln(os.Stderr, "Matched IP:", peerIP)
	if peerIP == nil {
		runenv.Abort("No IP match")
		return
	}

	watcher, writer := sdksync.MustWatcherWriter(runenv)
	defer watcher.Close()

	ctx := context.Background()

	seq, err := writer.Write(peerIPSubtree, &peerIP)
	if err != nil {
		runenv.Abort(err)
		return
	}
	defer writer.Close()

	// States
	ready := sdksync.State("ready")

	switch {
	case seq == 1: // receiver
		fmt.Fprintln(os.Stderr, "Receiver:", peerIP)

		quit := make(chan int)

		l, err := net.Listen("tcp", peerIP.String()+":2000")
		if err != nil {
			runenv.Abort(err)
			return
		}
		defer func() {
			close(quit)
			l.Close()
		}()

		// Signal we're ready
		_, err = writer.SignalEntry(ready)
		if err != nil {
			runenv.Abort(err)
			return
		}
		fmt.Fprintln(os.Stderr, "State: ready")

		var wg sync.WaitGroup
		wg.Add(runenv.TestInstanceCount - 1)
		fmt.Fprintln(os.Stderr, "Waiting for connections:", runenv.TestInstanceCount-1)

		go func() {
			for {
				// Wait for a connection.
				conn, err := l.Accept()
				if err != nil {
					select {
					case <-quit:
						return
					default:
						runenv.Abort(err)
						return
					}
				}
				// Handle the connection in a new goroutine.
				// The loop then returns to accepting, so that
				// multiple connections may be served concurrently.
				go func(c net.Conn) {
					defer c.Close()
					bytesRead := 0
					buf := make([]byte, 128*1024)
					for {
						n, err := c.Read(buf)
						bytesRead += n
						fmt.Fprintln(os.Stderr, "Received", n)
						if err == io.EOF {
							break
						}
					}
					fmt.Fprintln(os.Stderr, "Bytes read:", bytesRead)
					wg.Done()
				}(conn)
			}
		}()

		wg.Wait()

	case seq == 2: // sender
		fmt.Fprintln(os.Stderr, "Sender:", peerIP)

		// Connect to other peers
		peerIPCh := make(chan *net.IP, 16)
		cancel, err := watcher.Subscribe(peerIPSubtree, peerIPCh)
		if err != nil {
			runenv.Abort(err)
		}
		defer cancel()

		var peerIPsToDial = make([]net.IP, 0)
		for i := 0; i < runenv.TestInstanceCount; i++ {
			receivedPeerIP := <-peerIPCh
			if receivedPeerIP.String() == peerIP.String() {
				continue
			}
			peerIPsToDial = append(peerIPsToDial, *receivedPeerIP)
		}
		fmt.Fprintln(os.Stderr, "Waiting for ready")

		// Set a state barrier.
		readyCh := watcher.Barrier(ctx, ready, 1)

		// Wait until ready state is signalled.
		if err := <-readyCh; err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, "State: ready")

		for _, peerIPToDial := range peerIPsToDial {
			fmt.Fprintln(os.Stderr, "Dialing", peerIPToDial)
			conn, err := net.Dial("tcp", peerIPToDial.String()+":2000")
			if err != nil {
				// handle error
				fmt.Fprintln(os.Stderr, "Error", err)
				return
			}
			buf := make([]byte, 100*1024)
			for i := 0; i < len(buf); i++ {
				buf[i] = byte(i)
			}
			bytesWritten := 0
			for i := 0; i < 10; i++ {
				n, err := conn.Write(buf)
				fmt.Fprintln(os.Stderr, "Sent", n)
				bytesWritten += n
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error", err)
					break
				}
			}
			fmt.Fprintln(os.Stderr, "Bytes written:", bytesWritten)
			conn.Close()
		}

		/*
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
		*/
	default:
		runenv.Abort(fmt.Errorf("Unexpected seq: %v", seq))
	}

	runenv.OK()
}
