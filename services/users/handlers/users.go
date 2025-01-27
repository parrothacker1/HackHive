package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/parrothacker1/Solvelt/users/config"
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
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
      resp := response {
        Status: "error",
        Message: "Invalid body format",
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    if err := validate.Struct(user); err != nil {
      resp := response {
        Status: "error",
        Message: "Error with validation of data: "+err.Error(),
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    user.SetPassword(user.Password)
    if err := database.DB.Create(&user).Error; err != nil {
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
    type claims struct {
      UserID string `json:"user_id"`
      jwt.RegisteredClaims;
    }
    custom_claims := claims {
      user.UserID,
      jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(6*time.Hour)),
        Issuer: "Solvelt",
        IssuedAt: jwt.NewNumericDate(time.Now()),
      },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,custom_claims)
    if tokenString,err := token.SignedString(config.JWTSecret);err != nil {
      resp := response {
        Status: "fail",
        Message: "error with generating JWT",
      }
      w.WriteHeader(http.StatusInternalServerError)
      json.NewEncoder(w).Encode(resp)
    } else {
      type response_jwt struct {
        Status string
        Token string
      } 
      resp := response_jwt {
        Status: "success",
        Token: tokenString,
      }
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(resp)
    }
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

func ResetPassword() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("reset password of user"))
  }
}
