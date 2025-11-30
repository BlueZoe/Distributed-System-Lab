package main

import (
	"Lab1/utils"
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func forwardRequest(clientConn net.Conn, request *http.Request) {
	// Connect to remote server
	serverConn, err := net.Dial("tcp", request.Host)
	if err != nil {
		fmt.Printf("Proxy Error connecting to server, ErrorMessage: %s\n", err)
		response := utils.CreateResponse(http.StatusBadGateway, "Proxy cannot connect to target server", "text/plain")
		response.Write(clientConn)
		return
	}
	defer serverConn.Close()

	// Forward request to server
	err = request.Write(serverConn)
	if err != nil {
		fmt.Printf("Proxy Error forwarding request, ErrorMessage: %s\n", err)
		response := utils.CreateResponse(http.StatusBadGateway, "Proxy forward request error", "text/plain")
		response.Write(clientConn)
		return
	}

	// Read response from server
	serverReader := bufio.NewReader(serverConn)
	resp, err := http.ReadResponse(serverReader, request)
	if err != nil {
		fmt.Printf("Proxy Error reading response, ErrorMessage: %s\n", err)
		response := utils.CreateResponse(http.StatusBadGateway, "Proxy read response error", "text/plain")
		response.Write(clientConn)
		return
	}
	defer resp.Body.Close()

	// Forward server response to client
	resp.Write(clientConn)

	fmt.Println("Proxy forwarded successfully.")
}

func unsupportedRequestHandler(conn net.Conn, request *http.Request) {
	fmt.Printf("🔴 Proxy Unsupported request received for %s\n", request.URL.Path)
	response := utils.CreateResponse(http.StatusNotImplemented, "Unsupported request for proxy", "text/plain")
	response.Write(conn)
}
