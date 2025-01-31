package middlewares

import (
	"log"
	"testing"

	"github.com/parrothacker1/Solvelt/users/utils/tests"
)

func TestMain(m *testing.M) {
  container,err := tests.SetupTestDB()
  if err != nil {
    log.Fatal(container)
  }
  m.Run()
}

func TestAuth(t *testing.T) {

}

func TestTeam(t *testing.T) {

}

func TestCleanUp(t *testing.T) {

}
