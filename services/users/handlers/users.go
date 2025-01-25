package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/parrothacker1/Solvelt/users/database"
	"github.com/parrothacker1/Solvelt/users/models"
	"gorm.io/gorm"
)

//TODO: make sure that admin accout has following restrictions
//      a seperate SSL certificate for Admin (403 for normal folks)
//      a seperate SSL certificate for Moderator (403 for normal folks)
//      normal for Users
func CreateUser() http.HandlerFunc {
  validate := validator.New()
  return func(w http.ResponseWriter, r *http.Request) {
    var user models.User
    w.Header().Set("Content-Type","application/json")
    if r.Body == nil {
      resp := response {
        Status: "error",
        Message: "Invalid body",
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
      resp := response {
        Status: "error",
        Message: "Invalid body format",
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    err = validate.Struct(user)
    if err != nil {
      resp := response {
        Status: "error",
        Message: "Error with validation of data: "+err.Error(),
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    user.SetPassword(user.Password)
    err = database.DB.Create(&user).Error
    fmt.Println(err)
    if err != nil {
      resp := response {
        Status: "fail",
        Message: "Failed to create the user: "+err.Error(),
      }
      if err == gorm.ErrDuplicatedKey {
        w.WriteHeader(http.StatusConflict)
      } else {
        w.WriteHeader(http.StatusInternalServerError)
      }
      json.NewEncoder(w).Encode(resp)
      return
    }
    resp := response {
      Status: "success",
      Message: "User created successfully",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp)
  }
}

func UpdateUser() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("updating user"))
  }
}

func DeleteUser() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("deleting user"))
  }
}

func GetUser() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("get user details"))
  }
}

func Login() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("login the user"))
  }
}
