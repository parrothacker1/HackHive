package models

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/parrothacker1/HackHive/testutils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func init() {
  container,err := testutils.SetupTestDB()
  if err != nil {
    log.Fatalf("Failed to create Test Postgres Container: %v",err)
  }
  container_host,err := container.Host(testutils.TestContainerContext)
  if err != nil {
    log.Fatalf("Failed in getting host from container: %v",err)
  }
  container_port,err := container.MappedPort(testutils.TestContainerContext,"5432/tcp")
  if err != nil {
    log.Fatalf("Failed in getting port from container: %v",err)
  }
  endpoint := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable","tester_hackhive","testing_hackhive",container_host,container_port.Port(),"testdb_hackhive")
  testDB,err = gorm.Open(postgres.Open(endpoint),&gorm.Config{})
  if err != nil {
    log.Fatalf("Failed to connect to the postgres container: %v",err)
  }
  testDB.AutoMigrate(&User{},&Team{})
}

func TestUsers(t *testing.T) {
  t.Run("Creating a Admin User",func(t *testing.T) {
    result := testDB.Create(&User{
      Name: "Admin",
      Email: "admin@hackhive.com",
      Password: "testing",
      Role: "admin",
    })
    require.NoError(t,result.Error,"Failed to create an admin user")
    var testAdmin User 
    testDB.First(&testAdmin)
    fmt.Println(reflect.TypeOf(nil))
    require.Equal(t,(*string)(nil),testAdmin.TeamID,"The teamID is not null, which is supposed to be null")
    require.Equal(t,(*int64)(nil),testAdmin.Points,"The points is not null, which is supposed to be null")
//    require.Equal()
  })

}
