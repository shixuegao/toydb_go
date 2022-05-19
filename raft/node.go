package raft

type Node struct {
	id          string
	peers       []string
	term        uint64
	log         *Log
	msgSender   chan<- *Message
	stateSender chan<- Instruction
	queuedReqs  []struct {
		event Event
		addr  Addr
	}
	proxiedReqs map[string]Addr
	role        *Role
}

func (n *Node) abortProxied() error {
	return nil
}

func (n *Node) forwardQueued(leader Addr) error {
	return nil
}

func (n *Node) send(to Addr, event Event) error {
	return nil
}

func (n *Node) quorum() int {
	return (len(n.peers)+1)/2 + 1
}
