package main

import (
	"Lab1/utils"
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	// Read the command line to get port
	listenAddress, err := utils.GetAddress()
	if err != nil {
		fmt.Printf("Proxy Error get listening address %s", err)
		return
	}

	// Create the listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Printf("Proxy Error listening on %s, ErrorMessage: %s\n", listenAddress, err)
		return
	}
	// Ensure the listener to be closed and cleaning up resources
	defer listener.Close()
	fmt.Printf("Proxy is listening on %s\n", listenAddress)

	// Waiting for a connection
	for {
		// Accept() blocks until a new client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Proxy Error accepting connection, ErrorMessage: %s\n", err)
			continue
		}
		// spawn a new Goroutine to handle the request concurrently
		go handleProxy(conn)
	}
}

func handleProxy(clientConn net.Conn) {
	defer clientConn.Close()

	fmt.Printf("Proxy Handling new connection from %s\n", clientConn.RemoteAddr())

	// Create a buffer and read info from the connection
	bufferReader := bufio.NewReader(clientConn)

	// Parse request from client
	request, err := http.ReadRequest(bufferReader)
	if err != nil {
		// EOF error is common when client closes connection before sending complete request
		// This is normal behavior and not a critical error
		if errors.Is(err, io.EOF) || err.Error() == "EOF" {
			fmt.Printf("Proxy Client closed connection before sending complete request: %s\n", clientConn.RemoteAddr())
		} else {
			fmt.Printf("Proxy Error reading from connection, ErrorMessage: %s\n", err)
			response := utils.CreateResponse(http.StatusBadRequest, "Bad Request(read request error)", "text/plain")
			response.Write(clientConn)
		}
		return
	}

	fmt.Println("Proxy request received:", request.Method, request.URL.String())

	switch request.Method {
	case utils.GetRequest:
		forwardRequest(clientConn, request)
	default:
		unsupportedRequestHandler(clientConn, request)
	}
}
