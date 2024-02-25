# Hole Punching & NAT Traversal

An implementation of a UDP Hole Punching which uses synchronization techniques like the ones used in [DCUtR](https://github.com/libp2p/specs/blob/master/relay/DCUtR.md).

## Usage

This is just and experimental project to learn more about p2p, hole-punching & NAT Traversal.

#### Running the project (WIP)

```
Usage of ./main:
  -c	         node as client(c)
  -p string      port for local addr (default ":1111")
  -rAddr string  relay address (default ":5173")
  -s	         node as server(s)
```

**Running Node Client**
```sh
go run ./cmd/ -c -p :[CLIENT_PORT] -rAddr [RELAY_IP]:[RELAY_PORT]
```

**Running Relay Server**
```
go run ./cmd/ -s -p :[PORT]
```
