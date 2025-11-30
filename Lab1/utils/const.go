package utils

const MaxGoroutines = 10

const (
	GetRequest   = "GET"
	PostRequest  = "POST"
	OtherRequest = "OTHER"
)

var MimeTypeMapping = map[string]string{
	"html": "text/html",
	"txt":  "text/plain",
	"gif":  "image/gif",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"css":  "text/css",
}
