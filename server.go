package kademlia

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Server implementation

type KademliaServer struct {
	contact          *Contact
	PingContacts     chan Contact
	FindNodeRequests chan FindNodeRequest
	FindNodeReplies  chan []*ContactRecord
	PingReplies 	 chan Contact
	Done             chan bool
	Errors           chan error
}


func (ks *KademliaServer) GetStates() []string {
	states := []string{"PING", "STORE", "FIND_NODE", "FIND_VALUE"}
	return states
}

func (ks *KademliaServer) extractMessageAndOtherNodeId(message []string) (*NodeID, *NodeID, error) {
	messageId, err := NewNodeId(message[1])
	if err != nil {
		return nil, nil, err
	}

	otherNodeId, err := NewNodeId(message[0])
	if err != nil {
		return nil, nil, err
	}

	return &messageId, &otherNodeId, err
}

func (ks *KademliaServer) getLocalAddress() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", ks.contact.Ip+":"+strconv.Itoa(ks.contact.Port))

}

func (ks *KademliaServer) HandleMessage(splitMessage []string, address *net.UDPAddr) error {
	states := ks.GetStates()
	if len(splitMessage) < 3 {
		return errors.New("Empty message or unrecognized RPC")
	}
	err := error(nil)
	procedure := splitMessage[0]
	switch procedure {
	case states[0]:
		err = ks.handlePing(splitMessage[1:], address)
		break
	case states[1]:
		err = ks.handleFindNode(splitMessage[1:], address)
		break
	case states[2]:
		if len(splitMessage) < 4 {
			err = errors.New("Invalid amount of params for " + states[1])
		}
		err = ks.handleFindNode(splitMessage[1:], address)
		break
	case states[3]:
		fmt.Println("I got a find value message!")
		break
	default:
		err = errors.New("RPC not found!")
	}

	if err != nil {
		fmt.Println("Error: ", err)
	}
	return nil
}

func (ks *KademliaServer) ListenForMessages(server *net.UDPConn) {
	defer server.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := server.ReadFromUDP(buf)
		msg := string(buf[0:n])
		fmt.Println("Received ", msg, " from ", addr)
		if err != nil {
			ks.Errors <- err
		}
		if err = ks.HandleMessage(strings.Split(msg, " "), addr); err != nil {
			ks.Errors <- err
		}
	}
}

func (ks *KademliaServer) StartServer(self *Contact) error {
	ks.contact = self
	ks.Done = make(chan bool, 1)
	ks.PingContacts = make(chan Contact, 5)
	ks.Errors = make(chan error, 1)
	ks.FindNodeRequests = make(chan FindNodeRequest, 1)
	ks.FindNodeReplies = make(chan []*ContactRecord, 1)
	ks.PingReplies = make(chan Contact, 1)

	ServerAddr, e := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(self.Port))
	if e != nil {
		return e
	}
	l, e := net.ListenUDP("udp", ServerAddr)
	if e != nil {
		return e
	}

	go ks.ListenForMessages(l)
	return nil
}
