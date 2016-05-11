package kademlia

import (
	"fmt"
	"container/list"
	"time"
	"errors"
	"strings"
	"net"
	"strconv"
)


type FindValueReply struct {
	contacts []*ContactRecord
	found bool
	value string
}


type FindValueRequest struct {
	destination *Contact
	target    NodeID
}


func (k *Kademlia) IterativeFindValue(key string, delta int) (value string){
	hash := k.getHashForValue(key)
	contacts := new(list.List)

	valueToSearch, err := NewNodeId(fmt.Sprintf("%x", hash))
	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan []*ContactRecord)
	frontier := []*Contact{}
	seen := make(map[string] bool);

	for _, record := range k.routes.FindClosest(valueToSearch, delta) {
		contacts.PushFront(record)
		seen[record.node.Id.String()] = true;
		frontier = append(frontier, record.node)
	}

	pending := 0
	for i := 0; i < delta && len(frontier) > 0; i ++ {
		pending++
		front := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		go k.FindValue(FindNodeRequest{front, valueToSearch}, done)
	}

	for pending > 0 {
		contactRecords := <-done
		pending--
		for _, contactRec := range contactRecords {
			if _, ok := seen[contactRec.node.Id.String()]; ok == false {
				contacts.PushFront(contactRec)
				frontier = append(frontier, contactRec.node)
				seen[contactRec.node.Id.String()] = true
			}
		}


		for pending < delta && len(frontier) > 0 {
			front := frontier[len(frontier)-1]
			frontier = frontier[:len(frontier)-1]
			go k.FindValue(FindValueRequest{front, valueToSearch}, done)
			pending++
		}
	}
	return
}




func (k *Kademlia) FindValue(request FindValueRequest, done chan []*FindValueReply){
	k.server.FindValue(request)
	select {
	case findValueReply := <-k.server.FindValueReplies:
		done <- findValueReply
	case <-time.After(time.Second * 2):
		k.server.Errors <- errors.New("Request for ping timed out...")
		done <- &FindValueReply{found:false,
			contacts: []*ContactRecord{}}
	}
}



// Serverside

func (ks *KademliaServer) handleFindValue(message []string, address *net.UDPAddr) (err error) {
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
	ks.FindValueRequests <- FindValueRequest{&contact, target}
	select {
	case response := <-ks.FindValueReplies:
		ks.sendFindNodeReply(response, &contact)
	case <-time.After(time.Second * 2):
		return errors.New("Timeout in find node")
	}
	return
}





func (ks *KademliaServer) FindValue(request FindValueRequest) {
	go func(){

		// Writing
		msg := "FIND_VALUE " + ks.contact.Id.String() + " " + NewRandomNodeId().String() + " " + request.target.String()
		done := make(chan bool)
		ks.sendMessage(request.destination, msg, done)

		LocalAddr, err := ks.getLocalAddress()
		if err != nil {
			ks.Errors <- err
		}

		<- done
		Listener, e := net.ListenUDP("udp", LocalAddr)
		if e != nil {
			ks.Errors <- err
		}

		// Now it's time to read back
		buf := make([]byte, 1024)
		n, addr, err := Listener.ReadFromUDP(buf)
		if err != nil {
			ks.Errors <- err
		}
		reply := string(buf[0:n])
		fmt.Println("Message:", reply, "from ", addr)
		ks.handleFindValueReply(strings.Split(reply, " ")[1:])
		Listener.Close()
	}()
}

func (ks *KademliaServer) handleFindValueReply(reply []string){

	_, _, err := ks.extractMessageAndOtherNodeId(reply)
	if err != nil {
		ks.Errors <- err
	}
	splitContacts := strings.Split(reply[2:][0], ";")
	contacts := ks.extractContacts(splitContacts)
	ks.FindNodeReplies <- contacts

}
