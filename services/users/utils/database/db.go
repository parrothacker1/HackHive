package database

import (
	"github.com/parrothacker1/Solvelt/users/config"
	"github.com/parrothacker1/Solvelt/users/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnnectToDatabase() { 
  DB, err := gorm.Open(postgres.Open(config.ConnectionString), &gorm.Config{
    SkipDefaultTransaction: true,
    TranslateError: true,
    Logger: logger.Default.LogMode(logger.Silent),
  })
  if err != nil {
    logrus.Fatal("Failed to connect to the database")
  }
  logrus.Debug("Successfully connected to the database")
  DB.AutoMigrate(&models.User{},&models.Team{})
  logrus.Debug("Successfully migrated models")
}
