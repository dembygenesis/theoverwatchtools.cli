package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New(t *testing.T) {
	cfg, err := New()
	assert.NoError(t, err, "unexpected error initialising config")
	assert.NotNil(t, cfg, "unexpected nil config")
}
