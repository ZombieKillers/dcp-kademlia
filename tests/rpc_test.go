package rpc_test

import (
	"testing"
	"net"
	"../nodes"
	"fmt"
	"math/rand"
	"time"
)



func TestPingServer(t *testing.T){

	// Init
	rand.Seed(time.Now().UTC().UnixNano())
	ServerAddr, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		t.Fatal(err)
	}

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	if err != nil {
		t.Fatal(err)
	}


	// Writing
	msg := "PING " + nodes.NewRandomNodeId().String()
	buf := []byte(msg)

	_, err = Conn.Write(buf)
	if err != nil {
		t.Fatal(err)
	}
	Conn.Close()

	LocalAddr, err := net.ResolveUDPAddr("udp", Conn.LocalAddr().String())
	Listener, e := net.ListenUDP("udp", LocalAddr)
	if e != nil {
		t.Fatal(e)
	}

	// Now it's time to read back
	fmt.Println("Sent successfully, now reading back the reply... ")
	buf = make([]byte, 1024)
	n, addr, err := Listener.ReadFromUDP(buf)
	if err != nil {
		t.Fatal(err)
	}
	reply := string(buf[0:n])
	fmt.Println("Message:", reply, "from ", addr)
	return;

}
