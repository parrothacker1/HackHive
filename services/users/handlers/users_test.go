package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
  fmt.Println("testing user")
}

func TestCreateUser(t *testing.T) {
  body := `{"name":"Tester_Solvelt","email":"test@gmail.com","password":"tester","role":"admin"}`
  req,err := http.NewRequest("POST","/api/users",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  rr := httptest.NewRecorder()
  CreateUser().ServeHTTP(rr,req)
  require.Equal(t,http.StatusOK,rr.Code,"The CreateUser handler is not working")
  body = `{"name":"Tester_Solvelt1","email":"test@gmail.com","password":"tester","role":"user"}`
  req,err = http.NewRequest("POST","/api/users",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  CreateUser().ServeHTTP(rr,req)
  require.Equal(t,http.StatusConflict,rr.Code,"Creating the same user is working which should not (conflict)")
}

/*
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the HelloHandler
func TestHelloHandler(t *testing.T) {
	// Create a new HTTP request to pass to the handler
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the HelloHandler directly with the mock request and response recorder
	handler := http.HandlerFunc(HelloHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, status)
	}

	// Check the body of the response
	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("expected body %v, got %v", expected, rr.Body.String())
	}
}
*/
