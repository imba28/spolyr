package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew__correct_credentials(t *testing.T) {
	_, err := New("root", "example", testDatabaseName)
	assert.Nil(t, err)
}

func TestNew__wrong_credentials(t *testing.T) {
	_, err := New("", "", testDatabaseName)
	assert.Error(t, err, "should return error if database credentials are wrong")
}
