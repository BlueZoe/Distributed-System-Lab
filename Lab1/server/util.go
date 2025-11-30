package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"Lab1/utils"
)

func GetRequestHandler(conn net.Conn, request *http.Request) {
	fmt.Printf("🟢 GET request received for %s %s\n", request.URL.Path, request.Proto)

	localPath, _ := os.Getwd()
	// Use filepath.Join to properly construct the path, and clean it to remove any ".." or "." components
	filePath := filepath.Join(localPath, request.URL.Path)
	filePath = filepath.Clean(filePath)
	fmt.Printf("File path: %s\n", filePath)

	// Check if the file extension is supported
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")
	if !utils.IsFileExtensionSupported(extension) {
		fmt.Printf("Error unsupported file extension for %s, Extension: %s\n", filePath, extension)
		response := utils.CreateResponse(http.StatusBadRequest, "Unsupported file extension", "text/plain")
		response.Write(conn)
		return
	}

	// Check if the file exists (use filePath, not request.URL.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error file not found for %s, ErrorMessage: %s\n", filePath, err)
		response := utils.CreateResponse(http.StatusNotFound, "File not found", "text/plain")
		response.Write(conn)
		return
	}

	// Check if file could be read
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s, ErrorMessage: %s\n", filePath, err)
		response := utils.CreateResponse(http.StatusInternalServerError, "Internal server error (read file Error)", "text/plain")
		response.Write(conn)
		return
	}

	fileContent := string(file)
	// fmt.Println("File content", fileContent)
	response := utils.CreateResponse(http.StatusOK, fileContent, utils.MimeTypeMapping[extension])
	response.Write(conn)
}

func PostRequestHandler(conn net.Conn, request *http.Request) {
	fmt.Printf("🟡 POST request received for %s\n", request.URL.Path)

	localPath, _ := os.Getwd()
	// Use filepath.Join to properly construct the path, and clean it to remove any ".." or "." components
	filePath := filepath.Join(localPath, request.URL.Path)
	filePath = filepath.Clean(filePath)

	// Check if the file extension is supported
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")
	if !utils.IsFileExtensionSupported(extension) {
		fmt.Printf("Error unsupported file extension for %s, File extension: %s\n", filePath, extension)
		response := utils.CreateResponse(http.StatusBadRequest, "Unsupported file extension", "text/plain")
		response.Write(conn)
		return
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %s, ErrorMessage: %s\n", dir, err)
		response := utils.CreateResponse(http.StatusInternalServerError, "Directory not created successfully", "text/plain")
		response.Write(conn)
		return
	}

	// Create the local file
	content, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating local file: %s, ErrorMessage: %s\n", filePath, err)
		response := utils.CreateResponse(http.StatusInternalServerError, "File not created successfully", "text/plain")
		response.Write(conn)
		return
	}
	defer content.Close()

	// Copy the request body to the local file
	_, err = io.Copy(content, request.Body)
	if err != nil {
		fmt.Printf("Error copying request body to local file: %s, ErrorMessage: %s\n", filePath, err)
		response := utils.CreateResponse(http.StatusInternalServerError, "File not written successfully", "text/plain")
		response.Write(conn)
		return
	}

	response := utils.CreateResponse(http.StatusOK, "File created and written successfully", "text/plain")
	response.Write(conn)
}

func UnsupportedRequestHandler(conn net.Conn, request *http.Request) {
	fmt.Printf("🔴 Unsupported request received for %s\n", request.URL.Path)
	response := utils.CreateResponse(http.StatusNotImplemented, "Unsupported request", "text/plain")
	response.Write(conn)
}
