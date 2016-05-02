package server

import (
	"net"
	"fmt"
	"reflect"
	"strings"
	"errors"
)
func GetStates() []string{
	states := []string {"PING", "STORE", "FIND_NODE", "FIND_VALUE"}
	return states
}

func HandleMessage(splitMessage []string, address *net.UDPAddr) error {
	states := GetStates()
	if len(splitMessage) < 1 {
		return errors.New("Empty message or unrecognized RPC");
	}
	procedure := splitMessage[0]
	switch procedure {
	case states[0]:
		fmt.Println("I got a ping message!")
		break
	case states[1]:
		fmt.Println("I got a store message!")
		break
	case states[2]:
		fmt.Println("I got a find node message!")
		break
	case states[3]:
		fmt.Println("I got a find value message!")
		break
	default:
		return errors.New("RPC not found!")
	}
	return nil
}

func ListenForMessages(server *net.UDPConn) error {
	defer server.Close()

	fmt.Println("Listening for contacts here")
	buf := make([]byte, 1024)
	for {
		n,addr,err := server.ReadFromUDP(buf)
		msg := string(buf[0:n])
		fmt.Println("Received ", msg, " from ",addr)
		if err != nil {
			fmt.Println("Error: ",err)
			return err
		}

		go HandleMessage(strings.Split(msg, " "), addr)
	}
	return nil
}


func StartServer() error {
	ServerAddr, e := net.ResolveUDPAddr("udp",":12345")
	if e != nil {
		return e
	}
	l, e := net.ListenUDP("udp", ServerAddr)
	fmt.Println(reflect.TypeOf(l))
	if e != nil {
		return e
	}

	ListenForMessages(l)
	return nil
}