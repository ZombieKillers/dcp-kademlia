package kademlia

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
	"strings"
	"container/list"
)

type FindNodeRequest struct {
	destination *Contact
	target    NodeID
}

func NewFindNodeRequest(destination *Contact, target NodeID) (FindNodeRequest){
	return FindNodeRequest{destination, target}
}

func (k *Kademlia) FindNode(request FindNodeRequest, done chan []*ContactRecord){
	k.server.FindNode(request)
	select {
	case replyContact := <-k.server.FindNodeReplies:
		done <- replyContact
	case <-time.After(time.Second * 2):
		k.server.Errors <- errors.New("Request for ping timed out...")
		done <- []*ContactRecord{}
	}
}


func (k *Kademlia) IterativeFindNode(target NodeID, delta int) (ret *list.List){
	ret = new(list.List)
	done := make(chan []*ContactRecord)
	frontier := []*Contact{}
	seen := make(map[string] bool);

	for _, record := range k.routes.FindClosest(target, delta) {
		ret.PushFront(record)
		seen[record.node.Id.String()] = true;
		frontier = append(frontier, record.node)
	}

	pending := 0
	for i := 0; i < delta && len(frontier) > 0; i ++ {
		pending++
		front := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		go k.FindNode(FindNodeRequest{front, target}, done)
	}

	for pending > 0 {
		contactRecords := <-done
		pending--
		for _, contactRec := range contactRecords {
			if _, ok := seen[contactRec.node.Id.String()]; ok == false {
				ret.PushFront(contactRec)
				frontier = append(frontier, contactRec.node)
				seen[contactRec.node.Id.String()] = true
			}
		}


		for pending < delta && len(frontier) > 0 {
			front := frontier[len(frontier)-1]
			frontier = frontier[:len(frontier)-1]
			go k.FindNode(FindNodeRequest{front, target}, done)
			pending++
		}
	}
	return
}


// Serverside

func (ks *KademliaServer) handleFindNode(message []string, address *net.UDPAddr) (err error) {
	messageId, otherNodeId, err := ks.extractMessageAndOtherNodeId(message[0:2])
	if err != nil {
		return
	}
	target, err := NewNodeId(message[2])
	if err != nil {
		return
	}
	fmt.Println("MessageId:", messageId)
	contact := NewContact(*otherNodeId, address.IP.String(), address.Port)
	ks.FindNodeRequests <- FindNodeRequest{&contact, target}
	select {
	case response := <-ks.FindNodeReplies:
		ks.sendFindNodeReply(response, &contact)
	case <-time.After(time.Second * 2):
		return errors.New("Timeout in find node")
	}
	return
}

func buildNodeStr(records []*ContactRecord) string {
	s := ""
	for _, record := range records {
		s += record.node.String() + ";"
	}
	return s[0:len(s)-1]
}

func (ks *KademliaServer) sendFindNodeReply(nodes []*ContactRecord, contact *Contact) {
	replyString := buildNodeStr(nodes)
	fmt.Println("Sending reply to find node:", replyString)
	go func() {
		LocalAddr, err := ks.getLocalAddress()
		if err != nil {
			ks.Errors <- err
			return
		}
		ServerAddr, err := net.ResolveUDPAddr("udp", contact.Ip+":"+strconv.Itoa(contact.Port))
		if err != nil {
			ks.Errors <- err
			return
		}

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		if err != nil {
			ks.Errors <- err
		}

		// Writing
		msg := "FIND_NODE_REPLY " +
			ks.contact.Id.String() + " " +
			NewRandomNodeId().String() + " " +
			replyString

		buf := []byte(msg)

		_, err = Conn.Write(buf)
		if err != nil {
			ks.Errors <- err
		}
		Conn.Close()
	}()
}


func (ks *KademliaServer) FindNode(request FindNodeRequest) {
	go func(){

	LocalAddr, err := ks.getLocalAddress()
	if err != nil {
		fmt.Println(err)
		return
	}
	ServerAddr, err := net.ResolveUDPAddr("udp", request.destination.Ip+":"+strconv.Itoa(request.destination.Port))
	if err != nil {
		fmt.Println(err)
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	//ks.setReuseAddress(Conn)
	if err != nil {
		ks.Errors <- err
		return
	}

	// Writing
	msg := "FIND_NODE " + ks.contact.Id.String() + " " + NewRandomNodeId().String() + " " + request.target.String()
	buf := []byte(msg)

	_, err = Conn.Write(buf)
	if err != nil {
		ks.Errors <- err
	}
	Conn.Close()

	//LocalAddr, err = net.ResolveUDPAddr("udp", Conn.LocalAddr().String())
	Listener, e := net.ListenUDP("udp", LocalAddr)
	if e != nil {
		ks.Errors <- err
	}

	// Now it's time to read back
	buf = make([]byte, 1024)
	n, addr, err := Listener.ReadFromUDP(buf)
	if err != nil {
		ks.Errors <- err
	}
	reply := string(buf[0:n])
	fmt.Println("Message:", reply, "from ", addr)
	ks.handleFindNodeReply(strings.Split(reply, " ")[1:])
	Listener.Close()
	}()
	//ks.FindNodeReplies <- NewContact(otherNodeId, addr.IP.String(), addr.Port)
}


func (ks *KademliaServer) handleFindNodeReply(reply []string){

	_, _, err := ks.extractMessageAndOtherNodeId(reply)
	if err != nil {
		ks.Errors <- err
	}
	splitContacts := strings.Split(reply[2:][0], ";")
	contacts := ks.extractContacts(splitContacts)
	ks.FindNodeReplies <- contacts

}

func  (ks *KademliaServer) extractContacts(contacts []string) []*ContactRecord {
	res := make([]*ContactRecord, len(contacts))
	for i, contact := range contacts {
		contactString := strings.Split(strings.Split(strings.Split(contact, "(")[1], ")")[0], ",")
		nodeId, _ := NewNodeId(contactString[0])
		port, _ := strconv.Atoi(contactString[2])
		contact := NewContact(nodeId, contactString[1], port)

		fmt.Println("Foreign contact", contact.Id, "Current contact", ks.contact.Id)
		contactRecord := NewContactRecord(&contact, contact.Id.Distance(ks.contact.Id))
		res[i] = &contactRecord
	}
	return res
}
