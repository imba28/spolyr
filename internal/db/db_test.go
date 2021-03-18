package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	EnvUsername = "DATABASE_USER"
	EnvPassword = "DATABASE_PASSWORD"
	EnvHost     = "DATABASE_HOST"

	testDatabaseName = "tests"
)

func TestNew__correct_credentials(t *testing.T) {
	_, err := New(GetEnv(EnvUsername, ""), GetEnv(EnvPassword, ""), testDatabaseName, GetEnv(EnvHost, ""))
	assert.Nil(t, err)
}

func TestNew__wrong_credentials(t *testing.T) {
	_, err := New("", "", testDatabaseName, "localhost")
	assert.Error(t, err, "should return error if database credentials are wrong")
}