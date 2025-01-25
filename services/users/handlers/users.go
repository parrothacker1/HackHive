package handlers

import "net/http"

func CreateUser() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("creating user"))
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
