package models

import (
	"gorm.io/gorm"
)

type User struct {
  gorm.Model
  UserID string `gorm:"column:user_id;primaryKey" validate:"uuid4" json:"userID"`
  Name string `gorm:"column:user_name;not null" json:"name"`
  Email string `gorm:"column:user_email" validate:"email" json:"email"`
  Password string `gorm:"column:user_password" json:"password"`
  Role string `gorm:"column:user_role;not null" validate:"oneof=user admin moderator" json:"role"` // admin | moderator | user 
  TeamID *string `gorm:"column:team_id" validate:"uuid4" json:"teamID"`
  Points *int64 `gorm:"column:user_points" validate:"gte=0" json:"points"`

  Team *Team `gorm:"foreignKey:TeamID;references:TeamID" json:"-"`
}

type Team struct {
  gorm.Model
  TeamID string `gorm:"column:team_id;primaryKey" json:"teamID"`
  Name string `gorm:"column:team_name;not null" json:"name"`
  Leader string `gorm:"column:team_leader;unique" json:"leader"`
  Points int `gorm:"column:team_points" json:"teamPoints"`

  Members []User `gorm:"foreignKey:TeamID;references:TeamID" json:"users"`
  LeaderUser User `gorm:"foreignKey:Leader;references:UserID" json:"-"`
}

