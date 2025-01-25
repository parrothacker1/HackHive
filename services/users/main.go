package main

import (
	"net/http"

	"github.com/parrothacker1/Solvelt/users/database"
	"github.com/parrothacker1/Solvelt/users/routers"
)

func main() {
  userRouter := routers.AddUserRoutes()
  teamRouter := routers.AddTeamRouter()
  mux := http.NewServeMux()
  mux.Handle("/api/users/",http.StripPrefix("/api/users",&userRouter))
  mux.Handle("/api/teams/",http.StripPrefix("/api/teams",&teamRouter))
  database.ConnnectToDatabase()
  http.ListenAndServe("0.0.0.0:3000",mux)
}
