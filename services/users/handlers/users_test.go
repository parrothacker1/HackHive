package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parrothacker1/Solvelt/users/config"
	"github.com/stretchr/testify/require"
)

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
  var response struct {
    Status string;
    Token string;
  }
  json.NewDecoder(rr.Body).Decode(&response)
  _,err = jwt.Parse(response.Token,func(t *jwt.Token) (interface{}, error) {
    if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil,fmt.Errorf("Wrong algorithm")
    }
    return config.JWTSecret,nil
  })
  require.NoError(t,err,"JWT cannot be verified which it should.")
  body = `{"name":"Tester_Solvelt1","email":"test@gmail.com","password":"tester","role":"user"}`
  req,err = http.NewRequest("POST","/api/users",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  rr1 := httptest.NewRecorder()
  CreateUser().ServeHTTP(rr1,req)
  require.Equal(t,http.StatusConflict,rr1.Code,"Creating the user with same email is working which should not.")
}
