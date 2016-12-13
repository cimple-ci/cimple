package project

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArtifactParser_NoDestination(t *testing.T) {
	assert := assert.New(t)

	artifactHcl := `
artifact example {
	file = "/path/to/files"
	skip = true
	env {
		VAL = "1"
	}
}
`
	ast, err := extractObject(artifactHcl)
	if assert.Nil(err) {
		parser := &ArtifactParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		artifact := step.(ArtifactStep)
		assert.Equal("example", artifact.GetName())
		assert.True(artifact.GetSkip())
		assert.Equal("/path/to/files", artifact.File)
		assert.Equal(1, len(artifact.GetEnv()))
		assert.Equal("1", artifact.GetEnv()["VAL"])

		assert.Empty(artifact.Destinations)
	}
}

func TestArtifactParser_BintrayDestination(t *testing.T) {
	assert := assert.New(t)

	artifactHcl := `
artifact example {
	destination bintray {
		subject = "my-subject"
		repository = "my-repo"
		package = "my-package"
	}
	file = "/path/to/files"
}
`
	ast, err := extractObject(artifactHcl)
	if assert.Nil(err) {
		parser := &ArtifactParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		artifact := step.(ArtifactStep)
		assert.Equal("/path/to/files", artifact.File)
		assert.Equal(1, len(artifact.Destinations))

		destination := artifact.Destinations[0].(*bintrayArtifactDestination)
		assert.Equal("my-subject", destination.Subject)
		assert.Equal("my-repo", destination.Repository)
		assert.Equal("my-package", destination.Package)
	}
}

func TestArtifactParser_MultipleDestinations(t *testing.T) {
	assert := assert.New(t)

	artifactHcl := `
artifact example {
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
	file = "/path/to/files"
}
`
	ast, err := extractObject(artifactHcl)
	if assert.Nil(err) {
		parser := &ArtifactParser{}
		step, err := parser.Parse(ast)
		assert.Nil(err)

		artifact := step.(ArtifactStep)
		assert.Equal(2, len(artifact.Destinations))
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

	matches := list.Filter("artifact")

	return matches.Items[0], nil
}
