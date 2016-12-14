package project

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublishParser_NoDestination(t *testing.T) {
	assert := assert.New(t)

	publishHcl := `
publish example {
	files = ["/path/to/files"]
	skip = true
	env {
		VAL = "1"
	}
}
`
	ast, err := extractObject(publishHcl)
	if assert.Nil(err) {
		parser := &PublishParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		publishStep := step.(PublishStep)
		assert.Equal("example", publishStep.GetName())
		assert.True(publishStep.GetSkip())
		assert.Equal([]string{"/path/to/files"}, publishStep.Files)
		assert.Equal(1, len(publishStep.GetEnv()))
		assert.Equal("1", publishStep.GetEnv()["VAL"])

		assert.Empty(publishStep.Destinations)
	}
}

func TestPublishParser_BintrayDestination(t *testing.T) {
	assert := assert.New(t)

	publishHcl := `
publish example {
	destination bintray {
		subject = "my-subject"
		repository = "my-repo"
		package = "my-package"
	}
	files = ["/path/to/files"]
}
`
	ast, err := extractObject(publishHcl)
	if assert.Nil(err) {
		parser := &PublishParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		publishStep := step.(PublishStep)
		assert.Equal([]string{"/path/to/files"}, publishStep.Files)
		assert.Equal(1, len(publishStep.Destinations))

		destination := publishStep.Destinations[0].(*bintrayPublishDestination)
		assert.Equal("my-subject", destination.Subject)
		assert.Equal("my-repo", destination.Repository)
		assert.Equal("my-package", destination.Package)
	}
}

func TestPublishParser_MultipleDestinations(t *testing.T) {
	assert := assert.New(t)

	publishHcl := `
publish example {
	destination bintray {
		subject = "my-subject"
		repository = "my-repo"
		package = "my-package"
	}
	destination bintray {
		subject = "my-other-subject"
		repository = "my-other-repo"
		package = "my-other-package"
	}
	files = ["/path/to/files"]
}
`
	ast, err := extractObject(publishHcl)
	if assert.Nil(err) {
		parser := &PublishParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		publishStep := step.(PublishStep)
		assert.Equal(2, len(publishStep.Destinations))
	}
}

func extractObject(h string) (*ast.ObjectItem, error) {
	file, err := hcl.Parse(h)
	if err != nil {
		return nil, err
	}

	list, ok := file.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("Failed to turn node into ObjectList")
	}

	matches := list.Filter("publish")

	return matches.Items[0], nil
}
