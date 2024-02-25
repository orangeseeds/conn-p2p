# Hole Punching & NAT Traversal

An implementation of a UDP Hole Punching which uses synchronization techniques like the ones used in [DCUtR](https://github.com/libp2p/specs/blob/master/relay/DCUtR.md).

## Usage

This is just and experimental project to learn more about p2p, hole-punching & NAT Traversal.

#### Running the project (WIP)

```
Usage of ./main:
  -c	         node is client
  -p string      port for local addr (default "1111")
  -rAddr string  relay address (default "192.168.1.71:5173")
  -s	         is relay server
```

**Running Node Client**
```sh
go run ./cmd/ -c -p :[CLIENT_PORT] -rAddr [RELAY_IP]:[RELAY_PORT]
```

**Running Relay Server**
```
go run ./cmd/ -s -p :[PORT]
```
