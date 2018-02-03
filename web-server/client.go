/*
 * Client
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s host port filename", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1] + ":" + os.Args[2]

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	reqMsg := []byte("GET /" + os.Args[3] + " HTTP/1.1\r\n\r\n")
	_, err = conn.Write(reqMsg)
	checkError(err)

	res, err := ioutil.ReadAll(conn)
	checkError(err)

	fmt.Println(string(res))

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
