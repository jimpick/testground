# Networking

All testground runners _except_ for the local runner have two networks: a
control and a data network.

* Test instances communicate with each other over the data network.
* Test instances communicate with the sync service, and _only_ the sync service,
  over the control network.

The local runner will use your machines local network interfaces.

## Control Network

The "control network" runs on 192.168.0.1/16 and should only be used to
communicate with the sync service.

After the sidecar is finished [initializing the
network](https://github.com/ipfs/testground/blob/master/docs/SIDECAR.md#initialization),
it should be impossible to use this network to communicate with other nodes.
However, a good test plan should avoid listening on and/or announcing this
network _anyways_ to ensure that it doesn't interfere with the test.

## Data Network

The "data network", used for all inter-instance communication, runs on:

- 8.0.0.0/8  -- public
- 9.0.0.0/8  -- public
- 10.0.0.0/8 -- private

With a 8.0.0.1 as the gateway.

On start, your test instance will be assigned an IP address from one of these
subnets. You can change your IP address (within this range) at any time [using
the
sidecar](https://github.com/ipfs/testground/blob/master/docs/SIDECAR.md#ip-addresses).
