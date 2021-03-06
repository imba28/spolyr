package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	EnvUsername = "DATABASE_USER"
	EnvPassword = "DATABASE_PASSWORD"
	EnvHost     = "DATABASE_HOST"

	testDatabaseName = "tests"
)

func TestNew__correct_credentials(t *testing.T) {
	r, err := New(os.Getenv(EnvUsername), os.Getenv(EnvPassword), testDatabaseName, os.Getenv(EnvHost), 3)
	defer tearDown(r)
	assert.Nil(t, err)
}

func TestNew__wrong_credentials(t *testing.T) {
	r, err := New("", "", testDatabaseName, "localhost", 3)
	defer tearDown(r)

	assert.Error(t, err, "should return error if database credentials are wrong")
}
