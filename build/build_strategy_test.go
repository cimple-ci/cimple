package build

import (
	"github.com/gyuho/goraph"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeTaskNode struct {
	id           string
	dependencies []string
}

func (tn fakeTaskNode) GetID() string {
	return tn.id
}

func (tn fakeTaskNode) GetDependencies() []string {
	return tn.dependencies
}

func Test_PopulateGraph(t *testing.T) {
	tasks := []TaskNode{
		&fakeTaskNode{
			id:           "1",
			dependencies: []string{},
		},
		&fakeTaskNode{
			id:           "2",
			dependencies: []string{"1"},
		},
	}

	graph := PopulateGraph(tasks)

	assert := assert.New(t)
	assert.NotNil(graph.GetNode(goraph.StringID("1")))
	assert.NotNil(graph.GetNode(goraph.StringID("2")))

	sources, _ := graph.GetSources(goraph.StringID("1"))
	assert.Equal(0, len(sources))

	sources, _ = graph.GetSources(goraph.StringID("2"))
	assert.Equal(1, len(sources))
	assert.NotNil(sources[goraph.StringID("1")])

	targets, _ := graph.GetTargets(goraph.StringID("1"))
	assert.Equal(1, len(targets))
	assert.NotNil(targets[goraph.StringID("2")])

	targets, _ = graph.GetTargets(goraph.StringID("2"))
	assert.Equal(0, len(targets))
}

func Test_EntryPoints(t *testing.T) {
	tasks := []TaskNode{
		&fakeTaskNode{
			id:           "1",
			dependencies: []string{},
		},
		&fakeTaskNode{
			id:           "2",
			dependencies: []string{"1"},
		},
	}

	graph := PopulateGraph(tasks)
	nodes := Entrypoints(graph)

	assert := assert.New(t)
	assert.Equal(1, len(nodes))
	assert.Equal(graph.GetNode(goraph.StringID("1")), nodes[0])
}

func Test_Build(t *testing.T) {
	tasks := []TaskNode{
		&fakeTaskNode{
			id:           "1",
			dependencies: []string{},
		},
		&fakeTaskNode{
			id:           "2",
			dependencies: []string{},
		},
		&fakeTaskNode{
			id:           "3",
			dependencies: []string{"1"},
		},
		&fakeTaskNode{
			id:           "4",
			dependencies: []string{"2", "3"},
		},
	}

	builder := NewBuildStrategy(tasks)

	processOrder := []string{}
	err := builder.Build(func(name string) error {
		processOrder = append(processOrder, name)
		return nil
	})

	assert := assert.New(t)

	if assert.Nil(err) {
		assert.Equal([]string{"1", "3", "2", "4"}, processOrder)
	}
}
