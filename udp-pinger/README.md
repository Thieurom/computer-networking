# Socket Programming Assignment 2: UDP Pinger

### Run the program

First run the server:
```
go run UDPPingerServer.go
```
The server will be running at the port 12000 and waiting to receive UDP packets.

Then run the client to ping the server:
```
go run UDPPingerClient.go $HOST:12000
```
$HOST is the IP address of the host that is running the server or `localhost` in case you are running the server on the same machine which the client is running on.
