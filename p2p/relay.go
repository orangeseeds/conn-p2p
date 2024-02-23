package p2p

type Relay struct {
	Node
	Connections map[string]*Peer
}

func NewRelay(laddr string) *Relay {
	return &Relay{
		Node:        *NewNode(laddr),
		Connections: make(map[string]*Peer),
	}
}

func Forward(msg Message) error {
    return nil
}
