package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
  gorm.Model
  UserID string `gorm:"column:user_id;primaryKey;unique" validate:"uuid4" json:"userID"`
  Name string `gorm:"column:user_name;not null" json:"name"`
  Email string `gorm:"column:user_email" validate:"email" json:"email"`
  Password string `gorm:"column:user_password;not null" json:"-"`
  Role string `gorm:"column:user_role;not null" validate:"oneof=user admin moderator" json:"role"` // admin | moderator | user 
  TeamID *string `gorm:"column:team_id" validate:"uuid4" json:"teamID"`
  Points *int64 `gorm:"column:user_points" validate:"gte=0" json:"points"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
  u.UserID = uuid.NewString()
  u.Password = u.Password // hash it 
  if u.Role == "admin" || u.Role == "moderator" {
    u.TeamID = nil
    u.Points = nil
  }
  return nil
}

func (u *User) AfterDelete(tx *gorm.DB) error { //need testing
  if u.Role == "user" {
    var team Team
    tx.First(&team,u.TeamID)
    if u.UserID == team.Leader {
      var next_leader User 
      tx.Where("team_id = ?",team.TeamID).Order("created_at asc").First(&next_leader)
      tx.Model(team).Update("team_leader",next_leader.UserID)
    }
  }
  return nil
}

type Team struct {
  gorm.Model
  TeamID string `gorm:"column:team_id;primaryKey;unique" json:"teamID"`
  Name string `gorm:"column:team_name;not null" json:"name"`
  Leader string `gorm:"column:team_leader;unique" json:"leader"`
  Points int `gorm:"column:team_points" json:"teamPoints"`
  Secret string `gorm:"column:team_secret" json:"-"`
  LeaderUser User `gorm:"foreignKey:Leader;references:UserID" json:"-"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error { // need testing
  t.TeamID = uuid.NewString()
  t.Secret = "testing" // this will be a random string
  return nil
}

func (t* Team) AfterDelete(tx *gorm.DB) error { // need testing
  tx.Model(&User{}).Where("team_id = ?",t.TeamID).Update("team_id",nil)
  return nil
}
