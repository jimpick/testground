package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	utils "github.com/ipfs/testground/plans/smlbench2/utils"
	"github.com/ipfs/testground/sdk/iptb"
	"github.com/ipfs/testground/sdk/runtime"
)

// simpleAddGetTC is a simple test that adds a file of the specified size to an IPFS node, and tries to fetch it from
// another local node, by swarm connecting to the former first. It measures: time to add, time to connect, time to get.
type SimpleAddGetTC struct {
	SizeBytes int64
}

var _ utils.SmallBenchmarksTestCase = (*SimpleAddGetTC)(nil)

func (tc *SimpleAddGetTC) Name() string {
	h := strings.ReplaceAll(strings.ToLower(humanize.IBytes(uint64(tc.SizeBytes))), " ", "")
	return fmt.Sprintf("simple-add-get-%s", h)
}

func (tc *SimpleAddGetTC) Configure(runenv *runtime.RunEnv, spec *iptb.TestEnsembleSpec) {
	spec.AddNodesDefaultConfig(iptb.NodeOpts{Initialize: true, Start: true}, "adder", "getter")
}

func (tc *SimpleAddGetTC) Execute(runenv *runtime.RunEnv, ensemble *iptb.TestEnsemble) {
	adder := ensemble.GetNode("adder").Client()
	getter := ensemble.GetNode("getter").Client()

	// generate a random file of the designated size.
	file := utils.TempRandFile(runenv, ensemble.TempDir(), tc.SizeBytes)
	defer os.Remove(file.Name())

	tstarted := time.Now()
	cid, err := adder.Add(file)
	if err != nil {
		runenv.Abort(err)
		return
	}

	runenv.EmitMetric(utils.MetricTimeToAdd, float64(time.Now().Sub(tstarted)/time.Millisecond))

	addrs, err := ensemble.GetNode("adder").SwarmAddrs()
	if err != nil {
		runenv.Abort(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tstarted = time.Now()
	err = getter.SwarmConnect(ctx, addrs[0])
	if err != nil {
		runenv.Abort(err)
		return
	}

	runenv.EmitMetric(utils.MetricTimeToConnect, float64(time.Now().Sub(tstarted)/time.Millisecond))

	tstarted = time.Now()
	err = getter.Get(cid, ensemble.TempDir())
	if err != nil {
		runenv.Abort(err)
		return
	}
	runenv.EmitMetric(utils.MetricTimeToGet, float64(time.Now().Sub(tstarted)/time.Millisecond))
	runenv.OK()
}
