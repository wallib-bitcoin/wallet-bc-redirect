package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRedirectGet(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", ts.URL+"/redirect?key=1&key=2", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the header api-key
	requestApiKey := req.Header.Get("X-Api-Key")
	responseApiKey := rr.Header().Get("X-Api-Key")
	if requestApiKey != responseApiKey {
		t.Errorf("handler returned unexpected header: got %v want %v", requestApiKey, responseApiKey)
	}

}

func TestRedirectPost(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("POST", "/redirect", bytes.NewBuffer([]byte(`{"key": "value"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the header api-key
	requestApiKey := req.Header.Get("X-Api-Key")
	responseApiKey := rr.Header().Get("X-Api-Key")
	if requestApiKey != responseApiKey {
		t.Errorf("handler returned unexpected header: got %v want %v", requestApiKey, responseApiKey)
	}
}

func TestRedirectPut(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("PUT", "/redirect", bytes.NewBuffer([]byte(`{"key": "value"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the header api-key
	requestApiKey := req.Header.Get("X-Api-Key")
	responseApiKey := rr.Header().Get("X-Api-Key")
	if requestApiKey != responseApiKey {
		t.Errorf("handler returned unexpected header: got %v want %v", requestApiKey, responseApiKey)
	}
}

func TestRedirectGetWithWrongQueryParameters(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/redirect?key()=1&key()=2'", nil)
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectGetWithWrongHeaders(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/redirect?key=1&key=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key()", "api-key")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectGetWithEmptyURL(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", "")
	if err != nil {
		return
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/redirect?key=1&key=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key()", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}

func TestRedirectGetWithWrongURL(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", "http¡¡¡")
	if err != nil {
		return
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/redirect?key=1&key=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key()", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectGetWithWrongURLAbsolute(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", "//localhost")
	if err != nil {
		return
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/redirect?key=1&key=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key()", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectPostWithEmptyBody(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("POST", "/redirect", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestRedirectPutWithEmptyBody(t *testing.T) {
	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("PUT", "/redirect", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "api-key")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if strings.TrimSpace(rr.Body.String()) != "Request body is empty" {
		t.Errorf("handler returned unexpected error message: got (%v) want (%v)",
			"Request body is empty", rr.Body.String())
	}

}

type LimitedReader struct {
	R io.Reader
	N int64
}

func (r *LimitedReader) Read(p []byte) (n int, err error) {
	if r.N <= 0 {
		return 0, io.ErrShortWrite
	}
	if int64(len(p)) > r.N {
		p = p[0:r.N]
	}
	n, err = r.R.Read(p)
	r.N -= int64(n)
	return
}

func TestRedirectPostWithInvalidBody(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("POST", "/redirect", &LimitedReader{R: bytes.NewReader([]byte("body")), N: 0})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	if strings.TrimSpace(rr.Body.String()) != "Error reading request body: short write" {
		t.Errorf("Error: got (%v) want (%v)",
			rr.Body.String(), "Error reading request body")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectPutWithInvalidBody(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("PUT", "/redirect", &LimitedReader{R: bytes.NewReader([]byte("body")), N: 0})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	if strings.TrimSpace(rr.Body.String()) != "Error reading request body: short write" {
		t.Errorf("Error: got (%v) want (%v)",
			rr.Body.String(), "Error reading request body")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRedirectBadHttpMethod(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET1", ts.URL+"/redirect?key=1&key=2", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "api-key")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// check the header api-key
	requestApiKey := req.Header.Get("X-Api-Key")
	responseApiKey := rr.Header().Get("X-Api-Key")
	if requestApiKey != responseApiKey {
		t.Errorf("handler returned unexpected header: got %v want %v", requestApiKey, responseApiKey)
	}

}

func TestRedirectGetWithoutApiKey(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	err = os.Setenv("X_API_KEY", "api-key")
	err = os.Setenv("TOKEN", "token")
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", ts.URL+"/redirect?key=1&key=2&api-key=94a08da1fecbb6e8b46990538c7b50b2", nil)
	req.Header.Add("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the header api-key
	requestApiKey := req.Header.Get("X-Api-Key")
	responseApiKey := rr.Header().Get("X-Api-Key")
	if requestApiKey != responseApiKey {
		t.Errorf("handler returned unexpected header: got %v want %v", requestApiKey, responseApiKey)
	}

}

func TestRedirectGetWithoutAndInvalidApiKey(t *testing.T) {

	// Create a test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	err := os.Setenv("REDIRECT_URL", ts.URL)
	err = os.Setenv("X_API_KEY", "api-key")
	err = os.Setenv("TOKEN", "token1")
	if err != nil {
		return
	}
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", ts.URL+"/redirect?key=1&key=2&api-key=78b1e6d775cec5260001af137a79dbd51", nil)
	req.Header.Add("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Api-Key", "api-key")
	handler := http.HandlerFunc(redirect)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
