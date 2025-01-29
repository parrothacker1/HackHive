package handlers

import (
	"net/http"
)

var CreateTeam http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("creating team"))
}

var UpdateTeam http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("updating team"))
}

var DeleteTeam http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("deleting team"))
}

var GetTeam http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("get details of team"))
}

var JoinTeam http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("user join a team"))
}
