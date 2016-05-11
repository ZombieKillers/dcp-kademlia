package kademlia

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"errors"
)

func (k *Kademlia) Ping(contact *Contact) (ret *Contact) {
	k.server.SendPing(contact)
	select {
	case replyContact := <-k.server.PingReplies:
		ret = &replyContact
	case <-time.After(time.Second * 1):
		k.server.Errors <- errors.New("Request for ping timed out...")
		ret = nil
	}

	return
}

/// Server stuff

func (ks *KademliaServer) SendPing(contact *Contact) {
	go func() {
		LocalAddr, err := ks.getLocalAddress()
		if err != nil {
			fmt.Println(err)
			return
		}
		ServerAddr, err := net.ResolveUDPAddr("udp", contact.Ip+":"+strconv.Itoa(contact.Port))

		if err != nil {
			fmt.Println(err)
			return
		}

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		ks.setReuseAddress(Conn)
		if err != nil {
			fmt.Println(err)
		}

		// Writing
		msg := "PING " + ks.contact.Id.String() + " " + NewRandomNodeId().String()
		buf := []byte(msg)

		_, err = Conn.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
		Conn.Close()
		Listener, e := net.ListenUDP("udp", LocalAddr)
		ks.setReuseAddress(Listener)
		if e != nil {
			fmt.Println(e)
		}

		// Now it's time to read back
		fmt.Println("Sent successfully, now reading back the reply... ")
		buf = make([]byte, 1024)
		n, addr, err := Listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		reply := string(buf[0:n])

		splitMessage := strings.Split(reply, " ")
		if len(splitMessage) < 3 {
			return
		}

		otherNodeId, err := NewNodeId(splitMessage[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		ks.PingReplies <- NewContact(otherNodeId, addr.IP.String(), addr.Port)
		Listener.Close()
	}()
}

func (ks *KademliaServer) handlePing(message []string, address *net.UDPAddr) error {
	messageId, otherNodeId, err := ks.extractMessageAndOtherNodeId(message[0:2])

	LocalAddr, err := ks.getLocalAddress()
	if err != nil {
		return err
	}

	Conn, err := net.DialUDP("udp", LocalAddr, address)
	if err != nil {
		return err
	}

	reply := "PONG " + ks.contact.Id.String() + " " + messageId.String()
	_, err = Conn.Write([]byte(reply))
	if err != nil {
		return err
	}
	Conn.Close()
	ks.PingContacts <- NewContact(*otherNodeId, address.IP.String(), address.Port)
	return nil
}
