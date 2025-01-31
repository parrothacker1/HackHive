package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/parrothacker1/Solvelt/users/config"
	"github.com/parrothacker1/Solvelt/users/models"
	"github.com/parrothacker1/Solvelt/users/utils/database"
	"gorm.io/gorm"
)

//TODO: make sure that admin accout has following restrictions
//      a seperate SSL certificate for Admin (403 for normal folks)
//      a seperate SSL certificate for Moderator (403 for normal folks)
//      normal for Users
var CreateUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  validate := validator.New()
  var user models.User
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
  custom_claims := jwtclaims {
    user.UserID,
    "",
    user.Role,
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

var UpdateUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("updating user"))
}

var DeleteUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  user_id,ok := r.Context().Value("user_id").(string)
  if user_id == "" || !ok {
    resp := response {
      Status: "error",
      Message: "failed to get user_id",
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(resp)
    return
  }
  result := database.DB.Where("user_id = ?",user_id).Delete(&models.User{})
  if result.RowsAffected == 0 {
    resp := response{
      Status: "fail",
      Message: "User with this ID does not exists",
    }
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(resp)
    return
  }
  if result.Error != nil {
    resp := response {
      Status: "error",
      Message: "Failed to delete user",
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(resp)
    return
  }
  resp := response {
    Status: "success",
    Message: "user deleted successfully",
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(resp)
}

var GetUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("get user details"))
}


var LoginUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
    validate := validator.New()
    type body struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=4"`
      }
    var request body
    if r.Body == nil {
        resp := response{
            Status:  "error",
            Message: "Invalid body",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(resp)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      resp := response{
        Status:  "error",
        Message: "Invalid body format",
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    if err := validate.Struct(request); err != nil {
      resp := response{
        Status:  "error",
        Message: "Error with validation of data: " + err.Error(),
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
    var user models.User
    if err := database.DB.Table("users").Where("user_email = ?", request.Email).First(&user).Error; err != nil {
      if err == gorm.ErrRecordNotFound { 
        resp := response{
          Status:  "fail",
          Message: "The email or password is invalid",
        }
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(resp)
        return
      } else {
        resp := response{
          Status: "error",
          Message: "Error in fetching data",
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(resp)
        return
      }
    }
    var teamid string
    if (user.TeamID == nil) { teamid = ""} else { teamid = *user.TeamID}
    if ok, _ := user.ComparePassword(request.Password); ok {
      claims := jwtclaims{
        user.UserID,
        teamid,
        user.Role,
        jwt.RegisteredClaims{
          ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
          Issuer:    "Solvelt",
          IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
      }
      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      tokenString, err := token.SignedString(config.JWTSecret)
      if err != nil {
        resp := response{
          Status:  "fail",
          Message: "Error with generating JWT",
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(resp)
        return
      }
      type response_jwt struct {
        Status string
        Token string
      }
      resp := response_jwt{
        Status: "success",
        Token:  tokenString,
      }
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(resp)
    } else {
      resp := response{
        Status:  "fail",
        Message: "The email or password is invalid",
      }
      w.WriteHeader(http.StatusUnauthorized)
      json.NewEncoder(w).Encode(resp)
    }
}


var ResetPassword http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("reset password of user"))
}
