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

func Test_makeCliSecretStore(t *testing.T) {
	assert := assert.New(t)

	values := []string{"type:key:password", "type:key2:p2", "type2:key3:p3"}
	store, err := makeCliSecretStore(values)

	if assert.Nil(err) {
		p, err := store.Get("type", "key")
		if assert.Nil(err) {
			assert.Equal("password", p)
		}

		p, err = store.Get("type", "key2")
		if assert.Nil(err) {
			assert.Equal("p2", p)
		}

		p, err = store.Get("type2", "key3")
		if assert.Nil(err) {
			assert.Equal("p3", p)
		}
	}
}

func Test_makeCliSecretStore_invalid_input(t *testing.T) {
	assert := assert.New(t)

	values := []string{"type:key"}
	store, err := makeCliSecretStore(values)

	assert.Nil(store)
	assert.NotNil(err)
}
