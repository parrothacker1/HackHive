package handlers

import (
	"bytes"
	"context"
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
  t.Run("Creating a normal user as admin",func(t *testing.T) {
    body := `{"name":"Tester_Solvelt","email":"test@gmail.com","password":"tester","role":"admin"}`
    req,err := http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    CreateUser.ServeHTTP(rr,req)
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
  })
  t.Run("Creating a normal user but with same email id",func(t *testing.T) {
    body := `{"name":"Tester_Solvelt1","email":"test@gmail.com","password":"tester","role":"user"}`
    req,err := http.NewRequest("POST","/api/users",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    CreateUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusConflict,rr.Code,"Creating the user with same email is working which should not.")
  })
}

func TestLogin(t *testing.T) {
  t.Run("Logging in as admin user",func(t *testing.T) {
    body := `{"email":"test@gmail.com","password":"tester"}`
    req,err := http.NewRequest("POST","/api/users/login",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    LoginUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"The LoginUser is not working")
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
  })
  t.Run("Logging in with wrong email ID",func(t *testing.T) {
    body := `{"email":"testing2@gmail.com","password":"tester"}`
    req,err := http.NewRequest("POST","/api/users/login",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    LoginUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusUnauthorized,rr.Code,"The email which does not exists is not returning 401")
  })
  t.Run("Logging in with wrong password",func(t *testing.T) {
    body := `{"email":"test@gmail.com","password":"testeridk"}`
    req,err := http.NewRequest("POST","/api/users/login",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    LoginUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusUnauthorized,rr.Code,"The email which exists but wrong password is not returning 401")
  })
}
func TestUpdateUser(t *testing.T) {
  body := `{"name":"Tester_Solvelt","email":"dummy@gmail.com","password":"tester","role":"admin"}`
  req,err := http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  rr := httptest.NewRecorder()
  CreateUser.ServeHTTP(rr,req)
  require.Equal(t,http.StatusOK,rr.Code,"The CreateUser handler is not working")
  var response struct {
    Status string;
    Token string;
  }
  json.NewDecoder(rr.Body).Decode(&response)
  token,err := jwt.Parse(response.Token,func(t *jwt.Token) (interface{}, error) {
    if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil,fmt.Errorf("Wrong algorithm")
    }
    return config.JWTSecret,nil
  })
  claims := token.Claims.(jwt.MapClaims)
  user_id := claims["user_id"].(string)
  t.Run("Updating the dummy user",func(t *testing.T) {
    body := `{"name":"Testing_Solvelt"}`
    req,err :=  http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    UpdateUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"The UpdateUser is not working which it should")
  })
  t.Run("Updating dummy with a occupied email",func(t *testing.T) {
    body := `{"email":"test@gmail.com"}`
    req,err :=  http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    UpdateUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusConflict,rr.Code,"The UpdateUser should return 409 as it conflicts.")
  })
  t.Run("Updating dummy with password/role field",func(t *testing.T) {
    body := `{"password":"testing","role":"admin"}`
    req,err :=  http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    UpdateUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusBadRequest,rr.Code,"The UpdateUser should return 400 as the body is wrong")
  })
  t.Run("Cleanup dummy user",func(t *testing.T) {
    req,err := http.NewRequest("DELETE","/api/users/delete",nil)
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    DeleteUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"Failed to delete dummy user")
  })
}

func TestResetPassword(t *testing.T) {
  body := `{"email":"test@gmail.com","password":"tester"}`
  req,err := http.NewRequest("POST","/api/users/login",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  rr := httptest.NewRecorder()
  LoginUser.ServeHTTP(rr,req)
  require.Equal(t,http.StatusOK,rr.Code,"The LoginUser is not working")
  var response struct {
    Status string;
    Token string;
  }
  json.NewDecoder(rr.Body).Decode(&response)
  token,err := jwt.Parse(response.Token,func(t *jwt.Token) (interface{}, error) {
    if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil,fmt.Errorf("Wrong algorithm")
    }
    return config.JWTSecret,nil
  })
  claims := token.Claims.(jwt.MapClaims)
  user_id := claims["user_id"].(string)
  t.Run("Reseting password of admin user",func(t *testing.T) {
    body := `{"password":"testing"}`
    req,err := http.NewRequest("POST","/api/reset",bytes.NewBuffer([]byte(body)))
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    ResetPassword.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"The password reset function failed.")
  })
  t.Run("Logging in after reset",func(t *testing.T) {
    body := `{"email":"test@gmail.com","password":"testing"}`
    req,err := http.NewRequest("POST","/api/users/login",bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type","application/json")
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    rr := httptest.NewRecorder()
    LoginUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"The Password didn't change after reset")
  })
}

func TestDeleteUser(t *testing.T) {
  body := `{"name":"Tester_Solvelt","email":"dummy1@gmail.com","password":"tester","role":"admin"}`
  req,err := http.NewRequest("POST","/api/users/create",bytes.NewBuffer([]byte(body)))
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    t.Fatalf("Error in creating request: %v",err)
  }
  rr := httptest.NewRecorder()
  CreateUser.ServeHTTP(rr,req)
  require.Equal(t,http.StatusOK,rr.Code,"The CreateUser handler is not working")
  var response struct {
    Status string;
    Token string;
  }
  json.NewDecoder(rr.Body).Decode(&response)
  token,err := jwt.Parse(response.Token,func(t *jwt.Token) (interface{}, error) {
    if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil,fmt.Errorf("Wrong algorithm")
    }
    return config.JWTSecret,nil
  })
  claims := token.Claims.(jwt.MapClaims)
  user_id := claims["user_id"].(string)
  t.Run("Deleting a dummy user",func(t *testing.T) {
    req,err := http.NewRequest("DELETE","/api/users/delete",nil)
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    DeleteUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"Failed to delete dummy user")
  })
  t.Run("Deleting a non existing user",func(t *testing.T) {
   req,err := http.NewRequest("DELETE","/api/users/delete",nil)
    if err != nil {
      t.Fatalf("Error in creating request: %v",err)
    }
    ctx := context.WithValue(req.Context(),"user_id",user_id)
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    DeleteUser.ServeHTTP(rr,req)
    require.Equal(t,http.StatusNotFound,rr.Code,"The non existence user should'nt be deleted")
  })
}
