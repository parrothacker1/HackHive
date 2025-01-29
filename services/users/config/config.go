package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var ConnectionString string = "sqlite:///./sqlite.db"
var JWTSecret []byte = []byte("testing")
var EventLog string = "./events.log"

// TODO: take JWT secret from env
// TODO: also make sure that this takes data if it is only going for server .. or else stick to normal
func config_db() {
  port_str := os.Getenv("DB_PORT")
  if port_str == "" {
    logrus.Fatal("The DB_PORT is not specified")
  }
  port,err := strconv.ParseUint(port_str,10,32)
  if err != nil {
    logrus.Fatal("There was an error in parsing DB_PORT: %v\n",err)
  }
  host := os.Getenv("DB_HOST")
  if host == "" {
    logrus.Fatal("The DB_HOST is not specified")
  }
  user := os.Getenv("DB_USER")
  if  user == "" {
    logrus.Fatal("The DB_USER is not specified")
  }
  password := os.Getenv("DB_PASSWORD")
  if password == "" {
    logrus.Fatal("The DB_PASSWORD is not specified")
  }
  db := os.Getenv("DB_DATABASE")
  if db == "" {
    logrus.Fatal("The DB_DATABASE is not specified")
  }
  ConnectionString = fmt.Sprintf(
    "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    host,
    port,
    user,
    password,
    db,
  )
}

func config_jwt() {
  if jwt_secret := os.Getenv("JWT_SECRET"); jwt_secret == "" {
    logrus.Fatal("The JWT_SECRET is not specified")
  } else {
    JWTSecret = []byte(jwt_secret)
  }
}

func config_log() {
  if log_file := os.Getenv("EVENT_LOG"); log_file == "" {
    logrus.Fatal("The EVENT_LOG is not specified")
  } else {
    EventLog = log_file
  }
}

func init() {
  if os.Getenv("GO_ENV") != "test" {
    config_db()
    config_jwt()
    config_log()
  }
}
