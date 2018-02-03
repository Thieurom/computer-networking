# Socket Programming Assignment 1: Web Server

You will develop a web server that handles one HTTP request at a time. Your web server should accept and parse the HTTP request, get the requested file from the server’s file system, create an HTTP response message consisting of the requested file preceded by header lines, and then send the response directly to the client. If the requested file is not present in the server, the server should send an HTTP “404 Not Found” message back to the client.


### Run the Server

```
go run web_server.go
```

### Test the Server

#### From command line

```
go run client.go $SERVER_IP 6789 index.html
```
$SERVER_IP is the IP address of the host that is running the server

#### From web browser

Open web browser with url: `http://$SERVER_IP:6789/index.html`
