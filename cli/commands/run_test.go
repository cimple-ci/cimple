package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun_Settings(t *testing.T) {
	assert := assert.New(t)

	command := Run()

	assert.Equal("run", command.Name)
	assert.Equal([]string{"r"}, command.Aliases)
}
