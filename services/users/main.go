package main

import (
	"net/http"
	"os"

	"github.com/parrothacker1/Solvelt/users/config"
	"github.com/parrothacker1/Solvelt/users/utils/database"
	"github.com/parrothacker1/Solvelt/users/utils/loggers"
	"github.com/parrothacker1/Solvelt/users/routers"
	"github.com/sirupsen/logrus"
)

func main() {
  file,err := os.OpenFile(config.EventLog,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666);if err != nil {
    logrus.Fatalf("Failed to open events.log file: %v",err);
  }
  defer file.Close()
  loggers.EventLogger.SetOutput(file)
  userRouter := routers.UserRouter
  teamRouter := routers.TeamRouter
  mux := http.NewServeMux()
  mux.Handle("/api/users/",http.StripPrefix("/api/users",userRouter))
  mux.Handle("/api/teams/",http.StripPrefix("/api/teams",teamRouter))
  database.ConnnectToDatabase()
  http.ListenAndServe("0.0.0.0:3000",mux)
}
