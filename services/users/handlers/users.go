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
      Status string   `json:"status"`
      Token string    `json:"token"`
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
  validate := validator.New()
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
  type body struct {
    Name    string    `json:"name" validate:"omitempty,min=3"`
    Email   string    `json:"email" validate:"omitempty,email"` 
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
  if err := validate.Struct(request); err != nil || (request.Name == "" && request.Email == "") {
    if err != nil {
      resp := response{
        Status:  "error",
        Message: "Error with validation of data: " + err.Error(),
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    } else {
      resp := response{
        Status: "error",
        Message: "Both fiels cannot be empty",
      }
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(resp)
      return
    }
  }
  update_data := map[string]interface{}{}
  if request.Name != "" {
    update_data["user_name"]=request.Name
  }
  if request.Email != "" {
    update_data["user_email"]=request.Email
  }
  if result := database.DB.Model(&models.User{}).Where("user_id = ?",user_id).Updates(update_data); result.Error != nil || result.RowsAffected == 0 {
    if result.Error == gorm.ErrRecordNotFound {
      resp := response {
        Status: "fail",
        Message: "The user with this ID does not exist",
      }
      w.WriteHeader(http.StatusNotFound)
      json.NewEncoder(w).Encode(resp)
      return
    } else if result.Error == gorm.ErrDuplicatedKey {
      resp := response {
        Status: "fail",
        Message: "The email is already occupied",
      }
      w.WriteHeader(http.StatusConflict)
      json.NewEncoder(w).Encode(resp)
      return
    } else {
      resp := response {
        Status: "error",
        Message: "Failed to update the user data",
      }
      w.WriteHeader(http.StatusInternalServerError)
      json.NewEncoder(w).Encode(resp)
      return
    }
  }
  resp := response {
    Status: "success",
    Message: "User data has been updated successfully",
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(resp)
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
  user_id, ok := r.Context().Value("user_id").(string)
  if user_id == "" || !ok {
    resp := response {
      Status: "error",
      Message: "failed to get user_id",
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(resp)
    return
  }
  var user_data models.User
  if err := database.DB.Table("users").Where("user_id = ?").Find(&user_data).Error; err != nil {
    if err == gorm.ErrRecordNotFound {
      resp := response {
        Status: "fail",
        Message: "The user with this ID does not exist",
      }
      w.WriteHeader(http.StatusNotFound)
      json.NewEncoder(w).Encode(resp)
      return
    } else {
      resp := response {
        Status: "error",
        Message: "Failed to retrive data from server",
      }
      w.WriteHeader(http.StatusInternalServerError)
      json.NewEncoder(w).Encode(resp)
      return
    }
  }
  type response_success struct {
    Status    string                  `json:"status"`;
    Data      map[string]interface{} `json:"details"`;
  }
  resp := response_success {
    Status: "success",
    Data: map[string]interface{}{
      "name": user_data.Name,
      "email": user_data.Email,
      "role": user_data.Role,
      "team_id": user_data.TeamID,
      "points": user_data.Points,
    },
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(resp)
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
  user_id,ok := r.Context().Value("user_id").(string)
  validate := validator.New()
  type body struct {
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
  if user_id == "" || !ok {
    resp := response {
      Status: "error",
      Message: "Failed to get the user ID",
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(resp)
    return
  }
  var user_data models.User
  if err := database.DB.Table("users").Where("user_id = ?",user_id).First(&user_data).Error; err != nil {
    if err == gorm.ErrRecordNotFound {
      resp := response{
        Status: "fail",
        Message: "User with this ID does not exist",
      }
      w.WriteHeader(http.StatusNotFound)
      json.NewEncoder(w).Encode(resp)
      return
    } else {
      resp := response{
        Status: "error",
        Message: "Failed to retrive data",
      }
      w.WriteHeader(http.StatusInternalServerError)
      json.NewEncoder(w).Encode(resp)
      return
    }
  }
  user_data.SetPassword(request.Password)
  if result := database.DB.Table("users").Where("user_id = ?",user_id).Update("user_password",user_data.Password); result.Error != nil || result.RowsAffected == 0 {
    resp := response{
      Status: "error",
      Message: "Failed to update the password",
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(resp)
    return
  }
  resp := response{
    Status: "success",
    Message: "Password reset successfully",
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(resp)
}
