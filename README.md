# GoCurl - A curl-like tool implemented in Go

## Overview
GoCurl is a command-line HTTP client tool inspired by curl, implemented in Go. It provides functionality for making HTTP requests with various methods, handling custom headers, supporting SSL/TLS, and managing file uploads and downloads.

## Features
- Support for different HTTP methods (GET, POST, PUT, DELETE, etc.)
- Custom header support
- JSON output formatting
- SSL/TLS support
- File upload functionality
- File download functionality
- Verbose output mode

## Installation
To build GoCurl from source:

```bash
git clone https://github.com/user/gocurl.git
cd gocurl
go build -o gocurl ./cmd/gocurl
```

## Usage

### Basic GET request
```bash
./gocurl https://example.com
```

### Using different HTTP methods
```bash
./gocurl -X POST https://example.com
./gocurl -X PUT https://example.com
./gocurl -X DELETE https://example.com
```

### Adding custom headers
```bash
./gocurl -H "Content-Type: application/json" -H "Authorization: Bearer token" https://example.com
```

### Sending data in the request body
```bash
./gocurl -X POST -d '{"key":"value"}' https://example.com
```

### JSON output formatting
```bash
./gocurl --json https://api.example.com
```

### Allowing insecure SSL connections
```bash
./gocurl -k https://self-signed.example.com
```

### Verbose output
```bash
./gocurl -v https://example.com
```

### Saving response to a file
```bash
./gocurl -o response.txt https://example.com
```

### Uploading a file
```bash
./gocurl -F "file=@/path/to/file.txt" https://upload.example.com
```

### Adding form fields with file upload
```bash
./gocurl -F "file=@/path/to/file.txt" -form "name=value" https://upload.example.com
```

## Command Line Options
```
-F string
    Upload file (format: fieldname=@filename)
-H string
    Headers to include (can be used multiple times)
-X string
    HTTP method to use (default "GET")
-d string
    Data to send in the request body
-form string
    Add form field (format: name=value)
-json
    Format output as JSON
-k
    Allow insecure SSL connections
-o string
    Write output to file instead of stdout
-v
    Enable verbose output
```

## Project Structure
```
gocurl/
├── cmd/
│   └── gocurl/           # Main application entry point
│       └── main.go       # Main function and CLI handling
├── pkg/
│   ├── http/             # HTTP request handling
│   │   ├── client.go     # HTTP client implementation
│   │   ├── request.go    # Request building
│   │   └── response.go   # Response handling
│   ├── output/           # Output formatting
│   │   ├── formatter.go  # Output formatter interface
│   │   └── json.go       # JSON output implementation
│   └── file/             # File operations
│       ├── upload.go     # File upload functionality
│       └── download.go   # File download functionality
└── go.mod                # Go module definition
```
