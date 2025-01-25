package handlers

import (
	"net/http"
)

func CreateTeam() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("creating team"))
  }
}

func UpdateTeam() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("updating team"))
  }
}

func DeleteTeam() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("deleting team"))
  }
}

func GetTeam() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("get details of team team"))
  }
}
