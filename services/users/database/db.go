package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/parrothacker1/Solvelt/users/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnnectToDatabase() {
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
  connection_str := fmt.Sprintf(
    "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    host,
    port,
    user,
    password,
    db,
  )
  DB, err = gorm.Open(postgres.Open(connection_str), &gorm.Config{
    SkipDefaultTransaction: true,
    TranslateError: true,
  })
  if err != nil {
    logrus.Fatal("Failed to connect to the database")
  }
  logrus.Debug("Successfully connected to the database")
  DB.AutoMigrate(&models.User{},&models.Team{})
  logrus.Debug("Successfully migrated models")
}
