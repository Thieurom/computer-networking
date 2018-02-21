package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s your-gmail-address receiver-email-address", os.Args[0])
		os.Exit(1)
	}

	// mail server
	gmailServer := "smtp.gmail.com:587"

	// establish a TCP connection with mailserver and print server response
	tcpAddr, err := net.ResolveTCPAddr("tcp", gmailServer)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	var resp []byte
	var n int
	resp, n = recvFromServer(conn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "220")

	// send EHLO command and print server response
	sendToServer(conn, "EHLO Alice\r\n")
	resp, n = recvFromServer(conn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "250")

	// send STARTTLS command and print server response
	sendToServer(conn, "STARTTLS\r\n")
	resp, n = recvFromServer(conn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "220")

	// convert to TLS connection
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: "smtp.gmail.com",
	})

	// send EHLO command again (through TLS) and print server response
	sendToServer(tlsConn, "EHLO Alice\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "250")

	// send AUTH LOGIN command
	sendToServer(tlsConn, "AUTH LOGIN\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "334")

	// send username
	username := base64.StdEncoding.EncodeToString([]byte(os.Args[1]))
	sendToServer(tlsConn, username+"\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "334")

	// send password
	fmt.Print("Enter Password: ")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	checkError(err)

	password = base64.StdEncoding.EncodeToString([]byte(password))
	sendToServer(tlsConn, password+"\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "235")

	// send MAIL FROM command and print server response
	fromEmail := os.Args[1]
	sendToServer(tlsConn, "MAIL FROM: <"+fromEmail+">\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "250")

	// send RCPT TO command and print server response
	toEmail := os.Args[2]
	sendToServer(tlsConn, "RCPT TO: <"+toEmail+">\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "250")

	// send DATA command and print server response
	sendToServer(tlsConn, "DATA\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "354")

	// send message data
	text := "I love computer networks!"
	message := "From: " + fromEmail + "\r\n"
	message += "To: " + toEmail + "\r\n"
	message += "Subject: " + text + "\r\n\r\n"
	message += text + "\r\n"
	sendToServer(tlsConn, message)

	// message ends with a single period
	sendToServer(tlsConn, ".\r\n")
	resp, n = recvFromServer(tlsConn)
	fmt.Println(string(resp[0:n]))
	checkResponse(resp[0:], "250")

	// send QUIT command and print server response
	sendToServer(tlsConn, "QUIT\r\n")

	// close connection
	tlsConn.Close()

	os.Exit(0)
}

func sendToServer(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	checkError(err)
}

func recvFromServer(conn net.Conn) ([]byte, int) {
	var resp [512]byte

	n, err := conn.Read(resp[0:])
	checkError(err)

	return resp[0:], n
}

func checkResponse(resp []byte, code string) {
	if string(resp[0:3]) != code {
		fmt.Fprintf(os.Stderr, "%s reply not received from server.", code)
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
