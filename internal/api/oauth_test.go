package api

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv_returns_value_of_env_var(t *testing.T) {
	key := "foo"
	value := "bar"
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	assert.Equal(t, getEnv(key, ""), value)
}

func TestGetEnv_returns_fallback_value_if_env_var_is_missing(t *testing.T) {
	assert.Equal(t, getEnv("foo", "fallback"), "fallback")
}
