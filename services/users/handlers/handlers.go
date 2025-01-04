package handlers

import "net/http"

func CreateUser() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Die"))
  }
}
