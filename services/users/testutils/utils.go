package testutils

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var TestContainerContext context.Context

func SetupTestDB() (testcontainers.Container,error) {
  TestContainerContext = context.Background()
  req := testcontainers.ContainerRequest{
    Image: "postgres:17-alpine",
    ExposedPorts: []string{"5432/tcp"},
    Env: map[string]string{
      "POSTGRES_USER": "tester_hackhive",
      "POSTGRES_PASSWORD": "testing_hackhive",
      "POSTGRES_DB": "testdb_hackhive",
    },
    WaitingFor: wait.ForListeningPort("5432/tcp"),
  }
  postgresC,err := testcontainers.GenericContainer(TestContainerContext,testcontainers.GenericContainerRequest{
    ContainerRequest: req,
    Started: true,
  })
  if err != nil {
    return nil,err
  }
  return postgresC,nil
}

func CleanUpDB(t *testing.T,container testcontainers.Container) {
  testcontainers.CleanupContainer(t,container)
} 
