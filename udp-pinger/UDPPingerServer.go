package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

func main() {
	service := ":12000"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	var req [512]byte

	n, addr, err := conn.ReadFromUDP(req[0:])
	if err != nil {
		return
	}

	res := strings.ToUpper(string(req[0:n]))
	r := rand.Intn(10)
	if r < 4 {
		return
	}
	conn.WriteToUDP([]byte(res), addr)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
