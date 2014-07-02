# Crystal

Crystal is a decentralized, k-ordered id generation service. Run crystal on
each node in your infrastructure to generate conflict-free ids on-demand
without coordination.

Crystal produces IDs in the same format as
(flake)[https://github.com/boundary/flake], but crystal is implemented in Go to
make deployment and interfacing with it a tad easier.

For more information on the advantages of decentralized k-ordered IDs, see the
[Bounary blog post on
flake](http://www.boundary.com/blog/2012/01/12/flake-a-decentralized-k-ordered-unique-id-generator-in-erlang/).
Crystal is very similar to flake, but has no runtime dependencies on the Erlang
VM and can be interfaced with more easily in any language.

## What's a k-ordered ID?

Crystal produces 128 bit k-ordered IDs. The first 64 bits is a timestamp, so
IDs are time-ordered lexically.

The next 48 bits is a unique worker ID, usually set to the MAC address of a
machine's network interface, so IDs are conflict-free but without coordination.
You should run crystal locally on every node that needs to generate IDs.

The final 16 bits is a sequence ID, incremented each time an ID is generated in
the same millisecond as a previous ID. The sequence is reset to 0 when time
advances.

```
|--------------------------------|-------------------------|--------|
|              64                |            48           |   16   |
|          Timestamp (ms)        |   Worker ID (MAC addr)  |  Seqn  |
|--------------------------------|-------------------------|--------|
```
