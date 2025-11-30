package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"Lab1/utils"
	"golang.org/x/sync/semaphore"
)

func main() {
	// Read the command line to get port
	listenAddress, err := utils.GetAddress()
	if err != nil {
		fmt.Printf("Error get listening address %s", err)
		return
	}

	// Create the listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Printf("Error listening on %s, ErrorMessage: %s\n", listenAddress, err)
		return
	}
	// Ensure the listener to be closed and cleaning up resources
	defer listener.Close()
	fmt.Printf("Server is listening on %s\n", listenAddress)

	// Create a semaphore with a capacity of MAX_GOROUTINES.
	sem := semaphore.NewWeighted(int64(utils.MaxGoroutines))
	// Create a context to be used to cancel the requestHandler.
	ctx := context.Background()

	// Waiting for a connection
	for {
		// Accept() blocks until a new client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection, ErrorMessage: %s\n", err)
			continue
		}

		// Acquire semaphore. This will block if 10 goroutines are already running.
		fmt.Printf("⌛️[SEMAPHORE] Attempting to acquire semaphore for %s ...\n", conn.RemoteAddr())
		sem.Acquire(ctx, 1)
		fmt.Printf("✅[SEMAPHORE] Acquired semaphore for %s\n", conn.RemoteAddr())

		// spawn a new Goroutine to handle the request concurrently.
		go handleConnection(conn, sem, ctx)
	}
}

func handleConnection(conn net.Conn, sem *semaphore.Weighted, ctx context.Context) {
	defer conn.Close()
	// Release the semaphore when the request is handled.
	defer func() {
		sem.Release(1)
		fmt.Printf("👋[SEMAPHORE]Released semaphore for %s\n", conn.RemoteAddr())
	}()

	fmt.Printf("Handling new connection from %s\n", conn.RemoteAddr())

	// Create a buffer and read info from the connection
	bufferReader := bufio.NewReader(conn)

	// Read data from the connection.
	request, err := http.ReadRequest(bufferReader)
	if err != nil {
		// EOF error is common when client closes connection before sending complete request
		// This is normal behavior and not a critical error
		if errors.Is(err, io.EOF) || err.Error() == "EOF" {
			fmt.Printf("Client closed connection before sending complete request: %s\n", conn.RemoteAddr())
		} else {
			fmt.Printf("Error reading from connection, ErrorMessage: %s\n", err)
			response := utils.CreateResponse(http.StatusBadRequest, "Bad Request(read request error)", "text/plain")
			response.Write(conn)
		}
		return
	}

	// for debugging, print the request method
	fmt.Println("Request received:\n", request.Method, request.URL.String())
	bodyBytes, _ := io.ReadAll(request.Body)
	fmt.Printf("Request Body: %s\n", string(bodyBytes))

	switch request.Method {
	case utils.GetRequest:
		GetRequestHandler(conn, request)
	case utils.PostRequest:
		PostRequestHandler(conn, request)
	default:
		UnsupportedRequestHandler(conn, request)
	}
}
