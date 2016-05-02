package server

import (
	"net"
	"fmt"
	"strings"
	"errors"
	"../nodes"
)
func GetStates() []string{
	states := []string {"PING", "STORE", "FIND_NODE", "FIND_VALUE"}
	return states
}

func getLocalAddress() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", "127.0.0.1:0")

}

func handlePing(message []string, address *net.UDPAddr) error {
	messageId, err := nodes.NewNodeId(message[0])
	if err != nil {
		return err
	}

	LocalAddr, err := getLocalAddress()
	if err != nil {
		return err
	}

	Conn, err := net.DialUDP("udp", LocalAddr, address)
	if err != nil {
		return err
	}

	reply := "PONG " + messageId.String()
	_, err = Conn.Write([]byte(reply))
	if err != nil {
		return  err
	}
	Conn.Close()
	return nil
}


func HandleMessage(splitMessage []string, address *net.UDPAddr) error {
	states := GetStates()
	if len(splitMessage) < 2 {
		return errors.New("Empty message or unrecognized RPC");
	}
	err := error(nil)
	procedure := splitMessage[0]
	switch procedure {
	case states[0]:
		err = handlePing(splitMessage[1:], address)
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
		err = errors.New("RPC not found!")
	}

	if err != nil {
		fmt.Println("Error: ", err)
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

		HandleMessage(strings.Split(msg, " "), addr)
	}
	return nil
}


func StartServer() error {
	ServerAddr, e := net.ResolveUDPAddr("udp",":12345")
	if e != nil {
		return e
	}
	l, e := net.ListenUDP("udp", ServerAddr)
	if e != nil {
		return e
	}

	ListenForMessages(l)
	return nil
}