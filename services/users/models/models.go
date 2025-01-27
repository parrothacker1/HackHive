package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tabler interface {
  TableName() string
}

// TODO: abstract the password so that the object is not able to access it
type User struct {
  gorm.Model
  UserID    string  `gorm:"column:user_id;primaryKey;unique" validate:"omitempty,uuid4" json:"userID"`
  Name      string  `gorm:"column:user_name;not null" validate:"required,min=3,max=50" json:"name"`
  Email     string  `gorm:"column:user_email;unique" validate:"required,email" json:"email"`
  Password  string  `gorm:"column:user_password;not null" validate:"required" json:"password"`
  Role      string  `gorm:"column:user_role;not null" validate:"oneof=user admin moderator,required" json:"role"` // admin | moderator | user 
  TeamID    *string `gorm:"column:team_id" validate:"omitempty,uuid4" json:"teamID"`
  Points    *int64  `gorm:"column:user_points" validate:"omitempty,gte=0" json:"points"`
}

func (User) TableName() string {
  return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
  u.UserID = uuid.NewString()
  return nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
  if u.Role == "admin" || u.Role == "moderator" {
    u.TeamID = nil
    u.Points = nil
  }
  if (u.Password == "") {
    return fmt.Errorf("The password is not set")
  }
  return nil
}

func (u *User) AfterDelete(tx *gorm.DB) error {
  if u.Role == "user" && u.TeamID != nil {
    var team Team
    tx.Table("teams").Where("team_id = ?",u.TeamID).First(&team)
    if u.UserID == team.Leader {
      var next_leader User 
      tx.Where("team_id = ?",team.TeamID).Order("created_at asc").First(&next_leader)
      tx.Model(team).Update("team_leader",next_leader.UserID)
    }
  }
  return nil
}

func (u *User) SetPassword(password string) error {
  hashed,err := argon2id.CreateHash(password,argon2id.DefaultParams)
  if err != nil {
    return err
  }
  u.Password = hashed
  return nil
}

func (u *User) ComparePassword(password string) (bool,error) {
  if (u.Password == "") {
    return false,fmt.Errorf("The password is empty.")
  }
  return argon2id.ComparePasswordAndHash(password,u.Password)
}

type Team struct {
  gorm.Model
  TeamID      string  `gorm:"column:team_id;primaryKey;unique" validate:"omitempty,uuid4" json:"teamID"`
  Name        string  `gorm:"column:team_name;not null" validate:"required,min=3" json:"name"`
  Leader      string  `gorm:"column:team_leader;unique" validate:"required,uuid4" json:"leader"`
  Points      int     `gorm:"column:team_points" validate:"gte=0" json:"teamPoints"`
  Secret      string  `gorm:"column:team_secret" json:"-"`
  LeaderUser  *User   `gorm:"foreignKey:Leader;references:UserID" json:"-"`
}

func (Team) TableName() string {
  return "teams"
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
  t.TeamID = uuid.NewString()
  var err error
  t.Secret,err = func() (string,error) {
    bytes := make([]byte,32)
    _, err := rand.Read(bytes)
    if err != nil {
      return "",err
    }
    return hex.EncodeToString(bytes),nil
  }()
  return err
}

func (t* Team) AfterDelete(tx *gorm.DB) error {
  result := tx.Table("users").Where("team_id = ?",t.TeamID).Updates(map[string]interface{}{"team_id":nil})
  return result.Error
}
