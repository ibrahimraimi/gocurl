package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	customhttp "github.com/user/gocurl/pkg/http"
	"github.com/user/gocurl/pkg/file"
)

func main() {
	// Define command line flags
	method := flag.String("X", "GET", "HTTP method to use")
	headers := flag.String("H", "", "Headers to include (can be used multiple times)")
	data := flag.String("d", "", "Data to send in the request body")
	outputFile := flag.String("o", "", "Write output to file instead of stdout")
	jsonOutput := flag.Bool("json", false, "Format output as JSON")
	insecure := flag.Bool("k", false, "Allow insecure SSL connections")
	verbose := flag.Bool("v", false, "Enable verbose output")
	uploadFile := flag.String("F", "", "Upload file (format: fieldname=@filename)")
	formField := flag.String("form", "", "Add form field (format: name=value)")

	flag.Parse()

	// Check if URL is provided
	if flag.NArg() < 1 {
		fmt.Println("Error: URL is required")
		fmt.Println("Usage: gocurl [options] URL")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)

	// Create HTTP client
	clientOptions := customhttp.DefaultClientOptions()
	clientOptions.SkipTLSVerify = *insecure
	client := customhttp.NewClient(clientOptions)

	var stdReq *http.Request
	var err error

	// Handle file upload if specified
	if *uploadFile != "" {
		if !strings.Contains(*uploadFile, "=@") {
			fmt.Println("Error: File upload format should be fieldname=@filename")
			os.Exit(1)
		}

		parts := strings.SplitN(*uploadFile, "=@", 2)
		fieldName := parts[0]
		filePath := parts[1]

		// Parse form fields
		extraFields := make(map[string]string)
		if *formField != "" {
			formParts := strings.SplitN(*formField, "=", 2)
			if len(formParts) == 2 {
				extraFields[formParts[0]] = formParts[1]
			}
		}

		// Create multipart request
		files := []file.UploadFile{
			{
				FieldName: fieldName,
				FilePath:  filePath,
			},
		}

		// Force POST method for file uploads if not specified
		if *method == "GET" {
			*method = "POST"
		}

		// Create the multipart request
		stdReq, err = file.CreateMultipartRequest(url, files, extraFields)
		if err != nil {
			fmt.Printf("Error creating multipart request: %v\n", err)
			os.Exit(1)
		}

		// Add headers
		if *headers != "" {
			for _, header := range strings.Split(*headers, "\n") {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					stdReq.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}
	} else {
		// Regular request
		customReq := customhttp.NewRequest(*method, url)

		// Add headers
		if *headers != "" {
			for _, header := range strings.Split(*headers, "\n") {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					customReq.SetHeader(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}

		// Add body if provided
		if *data != "" {
			customReq.SetBody([]byte(*data))
			// Set content-type if not already set
			if _, exists := customReq.Headers["Content-Type"]; !exists {
				customReq.SetHeader("Content-Type", "application/x-www-form-urlencoded")
			}
		}

		// Build the request
		stdReq, err = customReq.Build()
		if err != nil {
			fmt.Printf("Error building request: %v\n", err)
			os.Exit(1)
		}
	}

	// Print verbose information if requested
	if *verbose {
		fmt.Printf("> %s %s\n", stdReq.Method, stdReq.URL)
		for key, values := range stdReq.Header {
			for _, value := range values {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		fmt.Println(">")
	}

	// Execute the request
	resp, err := client.Do(stdReq)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response, err := customhttp.NewResponse(resp)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	// Print verbose information if requested
	if *verbose {
		fmt.Printf("< %s\n", response.Status)
		for key, values := range response.Headers {
			for _, value := range values {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}
		fmt.Println("<")
	}

	// Determine output destination
	var output *os.File
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()

		// For file download, write response body directly to file
		if !*jsonOutput {
			// Create a reader from the response body
			bodyReader := bytes.NewReader(response.Body)
			
			// Download the content to the specified file
			err = file.Download(bodyReader, *outputFile)
			if err != nil {
				fmt.Printf("Error downloading content: %v\n", err)
				os.Exit(1)
			}
			
			if *verbose {
				fmt.Printf("Downloaded content to %s\n", *outputFile)
			}
			
			// Exit early as we've already written the file
			return
		}
	} else {
		output = os.Stdout
	}

	// Format and output the response
	if *jsonOutput {
		// For JSON output
		if *verbose {
			// For verbose JSON output, include headers and status
			data := map[string]interface{}{
				"status_code": response.StatusCode,
				"status":      response.Status,
				"headers":     response.Headers,
				"body":        string(response.Body),
			}
			jsonBytes, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintln(output, string(jsonBytes))
		} else {
			// Try to parse response as JSON first
			jsonData, jsonErr := response.JSON()
			if jsonErr == nil {
				jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
				if err != nil {
					fmt.Printf("Error formatting JSON: %v\n", err)
					os.Exit(1)
				}
				fmt.Fprintln(output, string(jsonBytes))
			} else {
				// If not valid JSON, output raw body as string
				fmt.Fprintln(output, string(response.Body))
			}
		}
	} else {
		// For plain text output
		fmt.Fprintln(output, string(response.Body))
	}
}
