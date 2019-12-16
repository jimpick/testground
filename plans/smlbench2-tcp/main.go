package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/ipfs/testground/sdk/runtime"
	sdksync "github.com/ipfs/testground/sdk/sync"
)

var (
        MetricBytesSent    = &runtime.MetricDefinition{Name: "sent_bytes", Unit: "bytes", ImprovementDir: -1}
        MetricBytesReceived = &runtime.MetricDefinition{Name: "received_bytes", Unit: "bytes", ImprovementDir: -1}
        MetricTimeToSend     = &runtime.MetricDefinition{Name: "time_to_send", Unit: "ms", ImprovementDir: -1}
        MetricTimeToReceive = &runtime.MetricDefinition{Name: "time_to_receive", Unit: "ms", ImprovementDir: -1}
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
		// fmt.Fprintln(os.Stderr, "Receiver:", peerIP)

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
		// fmt.Fprintln(os.Stderr, "State: ready")

		var wg sync.WaitGroup
		wg.Add(runenv.TestInstanceCount - 1)
		// fmt.Fprintln(os.Stderr, "Waiting for connections:", runenv.TestInstanceCount-1)

		go func() {
			for {
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
				go func(c net.Conn) {
					defer c.Close()
					bytesRead := 0
					buf := make([]byte, 128*1024)
					tstarted := time.Now()
					for {
						n, err := c.Read(buf)
						bytesRead += n
						// fmt.Fprintln(os.Stderr, "Received", n)
						if err == io.EOF {
							break
						}
					}
					runenv.EmitMetric(MetricBytesReceived, float64(bytesRead))
					runenv.EmitMetric(MetricTimeToReceive, float64(time.Now().Sub(tstarted)/time.Millisecond))
					wg.Done()
				}(conn)
			}
		}()

		wg.Wait()

	case seq == 2: // sender
		// fmt.Fprintln(os.Stderr, "Sender:", peerIP)

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

		// Wait until ready state is signalled.
		// fmt.Fprintln(os.Stderr, "Waiting for ready")
		readyCh := watcher.Barrier(ctx, ready, 1)
		if err := <-readyCh; err != nil {
			panic(err)
		}
		// fmt.Fprintln(os.Stderr, "State: ready")

		for _, peerIPToDial := range peerIPsToDial {
			// fmt.Fprintln(os.Stderr, "Dialing", peerIPToDial)
			conn, err := net.Dial("tcp", peerIPToDial.String()+":2000")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error", err)
				return
			}
			buf := make([]byte, 100*1024)
			for i := 0; i < len(buf); i++ {
				buf[i] = byte(i)
			}
			bytesWritten := 0
			tstarted := time.Now()
			for i := 0; i < 100; i++ {
				n, err := conn.Write(buf)
				// fmt.Fprintln(os.Stderr, "Sent", n)
				bytesWritten += n
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error", err)
					break
				}
			}
			runenv.EmitMetric(MetricBytesSent, float64(bytesWritten))
			runenv.EmitMetric(MetricTimeToSend, float64(time.Now().Sub(tstarted)/time.Millisecond))
			conn.Close()
		}

	default:
		runenv.Abort(fmt.Errorf("Unexpected seq: %v", seq))
	}

	runenv.OK()
}
