package models

import (
	"fmt"
	"log"
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
    testUser := User {
      Name: "Admin",
      Email: "admin@hackhive.com",
      Role: "admin",
    }
    testUser.SetPassword("testing")
    result := testDB.Create(&testUser)
    require.NoError(t,result.Error,"Failed to create an admin user")
    var testAdmin User 
    testDB.First(&testAdmin)
    require.Equal(t,(*string)(nil),testAdmin.TeamID,"The teamID is not null, which is supposed to be null")
    require.Equal(t,(*int64)(nil),testAdmin.Points,"The points is not null, which is supposed to be null")
    check,err := testAdmin.ComparePassword("testing")
    require.NoError(t,err,"There is an error when doing hash comparison")
    require.True(t,check,"The password and the hash are not equal")
    testDB.Delete(&testAdmin)
  })
}

func TestTeam(t *testing.T) {
  t.Run("Creating a Normal Team",func(t *testing.T) {
    var leader,normal User
    var testTeam Team
    leader = User {
      Name: "LeaderTestTeam",
      Email: "leader@hackhive.com",
      Role: "user",
    }
    leader.SetPassword("leader")
    normal = User{
      Name: "NormalTestTeam",
      Email: "normal@hackhive.com",
      Role: "user",
    }
    normal.SetPassword("normal")
    require.NoError(t,testDB.Create(&leader).Error,"Error in creating leader user")
    require.NoError(t,testDB.Create(&normal).Error,"Error in creating normal user")
    testTeam = Team {
      Name: "TestTeam",
      Points: 0,
      Leader: &leader.UserID,
    }
   require.NoError(t,testDB.Create(&testTeam).Error,"Error in creating team")
   require.NoError(t,testDB.Model(&leader).Update("team_id",testTeam.TeamID).Error,"Error in Updating Leader User")
   require.NoError(t,testDB.Model(&normal).Update("team_id",testTeam.TeamID).Error,"Error in Updating Normal User")
  })
  t.Run("Checking Team if leader is deleted",func(t *testing.T) {
    var testTeam Team 
    testDB.Find(&testTeam).Where("team_name = ?","TestTeam")
    fmt.Println(*testTeam.Leader)
    testDB.Where("user_id = ?",testTeam.Leader).Unscoped().Delete(&User{})
    testDB.Find(&testTeam).Where("team_name = ?","TestTeam")
    fmt.Println(*testTeam.Leader)
  })
  t.Run("Checking User if Team is deleted",func(t *testing.T) {
    var testUser User
    var testTeam Team 
    testDB.Find(&testTeam).Where("team_name = ?","TestTeam")
    testDB.Delete(&testTeam)
    testDB.First(&testUser)
    require.Equal(t,(*string)(nil),testUser.TeamID,"The TeamID is not null.The AfterDelete Hook is not working")
  })
}
