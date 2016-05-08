package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	states := map[int]string{0: "PING", 1: "STORE", 2: "FIND_NODE", 3: "FIND_VALUE"}

	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)

	defer Conn.Close()
	i := 0
	for {
		msg := states[i]
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		i = (i + 1) % 4
		time.Sleep(time.Second * 1)
	}
}
