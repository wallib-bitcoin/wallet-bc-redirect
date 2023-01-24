package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

	// make a request to the redirectURL based on the method of the original request
	var resp *http.Response
	var err error
	switch request.Method {

	case http.MethodGet:
		req, err := http.NewRequest(http.MethodGet, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header = request.Header
		resp, err = http.DefaultClient.Do(req)
		log.Println(fmt.Sprintf("Request from remote: %v", req))
	case http.MethodPost:
		req, err := http.NewRequest(http.MethodPost, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header = request.Header
		resp, err = http.DefaultClient.Do(req)
		log.Println(fmt.Sprintf("Request from redirect: %v", req))
	case http.MethodPut:
		req, err := http.NewRequest(http.MethodPut, redirectURL, request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header = request.Header
		resp, err = http.DefaultClient.Do(req)
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
