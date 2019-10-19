module github.com/ipfs/testground/plans/smlbench2

go 1.13

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/ipfs/go-ipfs-api v0.0.2
	github.com/ipfs/testground/sdk/iptb v0.0.0-00010101000000-000000000000
	github.com/ipfs/testground/sdk/runtime v0.0.0-00010101000000-000000000000
	github.com/ipfs/testground/sdk/sync v0.0.0-20191017072543-376444a0dd33
	github.com/libp2p/go-libp2p v0.4.0
	github.com/libp2p/go-libp2p-core v0.2.3
)

replace (
	github.com/ipfs/testground/sdk/iptb => ../../sdk/iptb
	github.com/ipfs/testground/sdk/runtime => ../../sdk/runtime
	github.com/ipfs/testground/sdk/sync => ../../sdk/sync
)
