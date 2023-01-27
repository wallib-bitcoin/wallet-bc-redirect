package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	http.HandleFunc("/", redirect)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		return
	}

}

func redirect(writer http.ResponseWriter, request *http.Request) {

	env := godotenv.Load()
	if env != nil {
		log.Fatalln("Error loading .env file")
	}

	redirectURL := os.Getenv("REDIRECT_URL")

	// if the redirectURL isn't set, return an error
	if redirectURL == "" {
		http.Error(writer, "REDIRECT_URL environment variable not set", http.StatusInternalServerError)
		return
	}

	// append the path of the original request to the redirectURL
	redirectURL += request.URL.Path

	// Validate the redirectURL
	if err := validateUrl(redirectURL); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(fmt.Sprintf("Headers from redirect: %v", request.Header))
	// Validate the headers
	for key := range request.Header {
		key = http.CanonicalHeaderKey(key)
		log.Println(fmt.Sprintf("Header from redirect: %v", key))
		if err := validateInput(key); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid header key: %s", key), http.StatusBadRequest)
			return
		}
	}

	queryParams := request.URL.Query()
	if queryParams != nil {
		if err := validateQueryParameters(queryParams); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		redirectURL += "?" + queryParams.Encode()
	}

	// set its timeout
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	// make a request to the redirectURL based on the method of the original request
	var resp *http.Response
	switch request.Method {

	case http.MethodGet:

		req, _ := http.NewRequest(http.MethodGet, redirectURL, nil)
		log.Println(fmt.Sprintf("Request from redirect: %v", req))

		ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
		req = req.WithContext(ctx)
		req.Header = request.Header
		header := validateApiKey(request.URL.Query().Get("api-key"), request.Header.Get("x-api-key"))
		req.Header.Set("x-api-key", header)

		resp, _ = client.Do(req)

		cancel()
	case http.MethodPost:

		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		req, _ := http.NewRequest(http.MethodPost, redirectURL, buf)
		log.Println(fmt.Sprintf("Request from redirect: %v", req))

		ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
		req = req.WithContext(ctx)
		req.Header = request.Header
		header := validateApiKey(request.URL.Query().Get("api-key"), request.Header.Get("x-api-key"))
		log.Println(fmt.Printf("Header: %v\n", request.URL.Query().Get("api-key")))
		req.Header.Set("x-api-key", header)

		resp, err = client.Do(req)
		cancel()
	case http.MethodPut:

		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		if buf.Len() == 0 {
			http.Error(writer, "Request body is empty", http.StatusBadRequest)
			return
		}

		req, _ := http.NewRequest(http.MethodPut, redirectURL, buf)
		log.Println(fmt.Sprintf("Request from redirect: %v", req))

		ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
		req = req.WithContext(ctx)
		req.Header = request.Header
		header := validateApiKey(request.URL.Query().Get("api-key"), request.Header.Get("x-api-key"))
		req.Header.Set("x-api-key", header)

		resp, err = client.Do(req)
		cancel()
	default:
		http.Error(writer, "Invalid request method", http.StatusBadRequest)
		return
	}

	log.Println(fmt.Sprintf("Response from remote: %v", resp))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error catch body from request: %s", err), http.StatusInternalServerError)
		}
	}(resp.Body)

	// copy the headers of the response
	for k, v := range resp.Header {
		writer.Header().Set(k, v[0])
	}

	// copy the status code of the response
	writer.WriteHeader(resp.StatusCode)

	// copy the response body to the client
	_, _ = io.Copy(writer, resp.Body)

}

func validateUrl(redirectUrl string) error {
	u, err := url.ParseRequestURI(redirectUrl)
	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	if !u.IsAbs() {
		return fmt.Errorf("url must be absolute")
	}
	return nil
}

func validateQueryParameters(queryParams url.Values) error {
	// Define a regular expression
	pattern := `^[a-zA-Z0-9:_/?-]+$`
	valid := regexp.MustCompile(pattern)

	// Validate the key parameters
	for key := range queryParams {
		if !valid.MatchString(key) {
			return fmt.Errorf("invalid query parameters: %v", key)
		}
	}

	return nil
}

func validateInput(input string) error {
	// Define a regular expression
	pattern := `^[a-zA-Z0-9:_/?-]+$`
	r := regexp.MustCompile(pattern)

	// Validate the input
	if !r.MatchString(input) {
		return fmt.Errorf("invalid input")
	}

	return nil
}

func validateApiKey(apiKey string, xApiKey string) string {

	apiKeyURL := os.Getenv("X_API_KEY")

	if apiKey != "" {
		token := os.Getenv("TOKEN")
		data := []byte(token)
		sum := md5.Sum(data)

		if fmt.Sprintf("%x", sum) == apiKey {
			return apiKeyURL
		}
		return ""

	} else {
		return xApiKey
	}

}
