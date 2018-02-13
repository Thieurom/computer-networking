package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	fmt.Println("Start sending 10 ping messages to", service)

	for i := 1; i <= 10; i++ {
		conn, err := net.DialUDP("udp", nil, udpAddr)
		checkError(err)

		sendTime := time.Now()
		message := "Ping " + strconv.Itoa(i) + " time=" + sendTime.Format("15:04:05.0000")
		_, err = conn.Write([]byte(message))
		checkError(err)

		deadline := sendTime.Add(time.Second)
		err = conn.SetReadDeadline(deadline)
		checkError(err)

		var res [512]byte
		n, err := conn.Read(res[0:])
		recvTime := time.Now()
		if err2, ok := err.(net.Error); ok && err2.Timeout() {
			fmt.Println("Request timed out")
			continue
		} else {
			fmt.Println(string(res[0:n]) + "; RTT=" + recvTime.Sub(sendTime).String())
			checkError(err)
		}
	}

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
