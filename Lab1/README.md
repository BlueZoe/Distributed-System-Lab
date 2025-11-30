# Running the server

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Navigate to the proxy directory and build:
   ```bash
   cd server
   go build -o http_server
   ```

3. Run the server:
   ```bash
   ./http_server [port]
   ```

# Running the proxy
1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Navigate to the proxy directory and build:
   ```bash
   cd proxy
   go build -o http_proxy
   ```

3. Run the proxy:
   ```bash
   ./http_proxy [port]
   ```

# Testing

## Test with curl

Use `curl` to test the server manually:

**GET Request**
```bash
# Get a text file
curl http://localhost:8080/resources/hello.txt

# Get an HTML file
curl http://localhost:8080/resources/wiki.html

# Get an image file and save
curl http://localhost:8080/resources/football.jpeg -o downloaded_image.jpeg
```

**POST Request**
```bash
# Upload a text file
curl -X POST http://localhost:8080/uploaded/hello.txt -d "Hello from curl"
```

**Error Cases:**
```bash
# File not found
curl -v http://localhost:8080/resources/nonexistent.txt

# Unsupported file extension
curl -X POST http://localhost:8080/uploaded/test.pdf -d "test"

# Unsupported HTTP method
curl -X PUT http://localhost:8080/resources/test.txt
```

**Request through proxy**
```bash
# Get a text file
curl http://localhost:8080/resources/hello.txt -x localhost:8081

# Get an HTML file
curl http://localhost:8080/resources/wiki.html -x localhost:8081

# Get an image file and save
curl http://localhost:8080/resources/football.jpeg -o downloaded_image.jpeg -x localhost:8081
```

**Error Cases:**
```bash
# Unsupported HTTP method for proxy
curl -X POST http://localhost:8080/uploaded/hello.txt -d "Hello from curl" -x localhost:8081

# Failed to connect to server
[server terminated] curl.exe http://localhost:8080/resources/hello.txt -x localhost:8081

# File not found
curl http://localhost:8080/resources/nonexistent.txt -x localhost:8081

# Unsupported file extension
curl http://localhost:8080/resources/test.pdf -x localhost:8081
```


## Test with script

A comprehensive test script is provided to verify all server functionality.

From the `server/` directory:

```bash
cd server
./test.sh [port] [host]
```

**Examples:**
```bash
# Use default port 8080
./test.sh

# Specify port
./test.sh 8080

# Specify port and host
./test.sh 8080 localhost
```
