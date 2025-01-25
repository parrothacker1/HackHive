package models

import (
	"fmt"
	"log"
	"testing"

	"github.com/parrothacker1/Solvelt/users/testutils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	container, err := testutils.SetupTestDB()
	if err != nil {
		log.Fatalf("Failed to create Test Postgres Container: %v", err)
	}

	containerHost, err := container.Host(testutils.TestContainerContext)
	if err != nil {
		log.Fatalf("Failed to get host from container: %v", err)
	}

	containerPort, err := container.MappedPort(testutils.TestContainerContext, "5432/tcp")
	if err != nil {
		log.Fatalf("Failed to get port from container: %v", err)
	}

	endpoint := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
    "tester_solvelt", "testing_solvelt", containerHost, containerPort.Port(), "testdb_solvelt")

	testDB, err = gorm.Open(postgres.Open(endpoint), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the postgres container: %v", err)
	}

	if err := testDB.AutoMigrate(&User{}, &Team{}); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	m.Run()
}

func TestCreateAdminUser(t *testing.T) {
	t.Run("Creating an Admin User", func(t *testing.T) {
		testUser := User{
			Name:  "Admin",
			Email: "admin@hackhive.com",
			Role:  "admin",
		}
		testUser.SetPassword("testing")

		result := testDB.Create(&testUser)
		require.NoError(t, result.Error, "Failed to create an admin user")

		var testAdmin User
		testDB.First(&testAdmin)
		require.Equal(t, (*string)(nil), testAdmin.TeamID, "The teamID is not null, it should be null")
		require.Equal(t, (*int64)(nil), testAdmin.Points, "The points is not null, it should be null")

		check, err := testAdmin.ComparePassword("testing")
		require.NoError(t, err, "There is an error when comparing the password hash")
		require.True(t, check, "The password and the hash are not equal")

		// Cleanup
		require.NoError(t, testDB.Delete(&testAdmin).Error, "Failed to delete admin user")
	})
}

func TestCreateAndAssignTeam(t *testing.T) {
	t.Run("Creating a Normal Team and Assigning Users", func(t *testing.T) {
		var leader, normal User
		var testTeam Team

		leader = User{
			Name:  "LeaderTestTeam",
			Email: "leader@hackhive.com",
			Role:  "user",
		}
		leader.SetPassword("leader")

		normal = User{
			Name:  "NormalTestTeam",
			Email: "normal@hackhive.com",
			Role:  "user",
		}
		normal.SetPassword("normal")

		require.NoError(t, testDB.Create(&leader).Error, "Error creating leader user")
		require.NoError(t, testDB.Create(&normal).Error, "Error creating normal user")

		testTeam = Team{
			Name:  "TestTeam",
			Points: 0,
			Leader: leader.UserID,
		}
		require.NoError(t, testDB.Create(&testTeam).Error, "Error creating team")

		require.NoError(t, testDB.Model(&leader).Update("team_id", testTeam.TeamID).Error, "Error updating leader user with team")
		require.NoError(t, testDB.Model(&normal).Update("team_id", testTeam.TeamID).Error, "Error updating normal user with team")
	})
}

func TestDeleteLeaderAndCheckTeam(t *testing.T) {
	t.Run("Leader Deletion Updates Team Leader", func(t *testing.T) {
		var testTeam Team
		require.NoError(t, testDB.Where("team_name = ?", "TestTeam").First(&testTeam).Error, "Error getting team")

		var leaderUser User
		require.NoError(t, testDB.Where("user_id = ?", testTeam.Leader).First(&leaderUser).Error, "Error getting the leader user")

		require.NoError(t, testDB.Delete(&leaderUser).Error, "Error deleting leader")

		var normalUser User
		require.NoError(t, testDB.Where("user_name = ?", "NormalTestTeam").First(&normalUser).Error, "Error getting the normal user")

		require.NoError(t, testDB.Where("team_name = ?", "TestTeam").First(&testTeam).Error, "Error getting team after deletion")
		require.Equal(t, normalUser.UserID, testTeam.Leader, "The leader was not updated after deletion")
	})
}

func TestDeleteTeamAndCheckUsers(t *testing.T) {
	t.Run("Team Deletion Nullifies User's TeamID", func(t *testing.T) {
		var testTeam Team
		require.NoError(t, testDB.Where("team_name = ?", "TestTeam").First(&testTeam).Error, "Error getting team")

		require.NoError(t, testDB.Delete(&testTeam).Error, "Error deleting team")

		var testUser User
		require.NoError(t, testDB.First(&testUser).Error, "Error getting user after team deletion")
		require.Equal(t, (*string)(nil), testUser.TeamID, "The TeamID is not null. AfterDelete hook might not be working")
	})
}

func TestCleanup(t *testing.T) {
	testutils.CleanUpDB(t, testutils.TestContainer)
}

