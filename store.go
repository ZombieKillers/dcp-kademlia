package kademlia

import (
	"net"
	"strings"
)

type KeyValuePair struct{
	Key string
	Value string
}

//Clientside
func (k *Kademlia) Store(contact *Contact, request *KeyValuePair, done chan bool) {
	msg := "STORE " + k.self.Id.String() + " " + NewRandomNodeId().String() + " " + request.Key + ":" + request.Value
	go k.server.sendMessage(contact, msg, done)
}


//Serverside
func (ks *KademliaServer) handleStore(msg []string, addr *net.UDPAddr) (err error){
	otherNodeId, err := NewNodeId(msg[0])
	if err != nil{
		return
	}
	ks.PingContacts <- NewContact(otherNodeId, addr.IP.String(), addr.Port)
	keyValuePair := strings.Split(msg[len(msg)-1], ":")
	ks.KeyValuePairs <- &KeyValuePair{Key:keyValuePair[0], Value:keyValuePair[1]}
	return
}