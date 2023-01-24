package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
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
	maxRedirections, _ := strconv.Atoi(os.Getenv("MAX_REDIR"))

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

	queryParams := request.URL.Query()
	if queryParams != nil {
		if err := validateQueryParameters(queryParams); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		redirectURL += "?" + queryParams.Encode()
	}

	// Validate the headers
	for key := range request.Header {
		key = http.CanonicalHeaderKey(key)
		if err := validateInput(key); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid header key: %s", key), http.StatusBadRequest)
			return
		}
	}

	// Parse the form values
	if err := request.ParseForm(); err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing form values: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the form values
	for key, values := range request.Form {
		if err := validateInput(key); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid form key: %s", key), http.StatusBadRequest)
			return
		}
		for _, value := range values {
			if err := validateInput(value); err != nil {
				http.Error(writer, fmt.Sprintf("Invalid form value: %s", value), http.StatusBadRequest)
				return
			}
		}
	}

	// set its timeout
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	// make a request to the redirectURL based on the method of the original request
	var resp *http.Response
	var err error
	var count = 0
	switch request.Method {

	case http.MethodGet:
		req, err := http.NewRequest(http.MethodGet, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("Request from redirect: %v", req))

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		req = req.WithContext(ctx)
		req.Header = request.Header
		resp, err = client.Do(req)
		cancel()

		count, err = countRedirections(count, maxRedirections)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Max %v", err), http.StatusBadRequest)
			return
		}
	case http.MethodPost:
		req, err := http.NewRequest(http.MethodPost, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}

		log.Println(fmt.Sprintf("Request from redirect: %v", req))

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		req = req.WithContext(ctx)
		req.Header = request.Header
		resp, err = client.Do(req)
		cancel()

		count, err = countRedirections(count, maxRedirections)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Max %v", err), http.StatusBadRequest)
			return
		}
	case http.MethodPut:
		req, err := http.NewRequest(http.MethodPut, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header = request.Header
		resp, err = client.Do(req)
		log.Println(fmt.Sprintf("Request from redirect: %v", req))
		count, err = countRedirections(count, maxRedirections)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Max %v", err), http.StatusBadRequest)
			return
		}
	default:
		http.Error(writer, "Invalid request method", http.StatusBadRequest)
		return
	}

	log.Println(fmt.Sprintf("Response from remote: %v", resp))

	if err != nil {
		http.Error(writer, fmt.Sprintf("Error making request to %s: %v", redirectURL, err), http.StatusInternalServerError)
		return
	}
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

	// Validate the query parameters
	validKey := regexp.MustCompile(pattern)
	for key := range queryParams {
		if !validKey.MatchString(key) {
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

func countRedirections(redirectCount int, maxRedirections int) (int, error) {
	if redirectCount > maxRedirections {
		return 0, fmt.Errorf("redirect has reahed the maximun number of redirections permitted: %v", redirectCount)
	}

	return redirectCount + 1, nil
}
