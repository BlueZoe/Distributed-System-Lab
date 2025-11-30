package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetAddress() (string, error) {
	args := os.Args

	if len(args) < 2 {
		return "", fmt.Errorf("missing port argument")
	}
	return ":" + args[1], nil
}

func CreateResponse(statusCode int, body string, contentType string) *http.Response {
	header := make(http.Header)
	header.Set("Content-Type", contentType)
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
}

func IsFileExtensionSupported(extension string) bool {
	_, ok := MimeTypeMapping[extension]
	return ok
}
