package build

import (
	"github.com/gyuho/goraph"
	"sort"
)

type TaskNode interface {
	GetID() string
	GetDependencies() []string
}

type BuildStrategy interface {
	Build(runner func(taskName string) error) error
}

type buildStrategy struct {
	graph      goraph.Graph
	builtTasks map[goraph.Node]bool
}

type nodeSorter struct {
	nodes []goraph.Node
}

func (ns nodeSorter) Len() int {
	return len(ns.nodes)
}

func (ns nodeSorter) Swap(i, j int) {
	ns.nodes[i], ns.nodes[j] = ns.nodes[j], ns.nodes[i]
}

func (ns nodeSorter) Less(i, j int) bool {
	return ns.nodes[i].String() < ns.nodes[j].String()
}

func PopulateGraph(tasks []TaskNode) goraph.Graph {
	graph := goraph.NewGraph()

	for _, t := range tasks {
		nd1 := goraph.NewNode(t.GetID())
		graph.AddNode(nd1)
	}

	for _, t := range tasks {
		for _, d := range t.GetDependencies() {
			graph.ReplaceEdge(goraph.StringID(d), goraph.StringID(t.GetID()), 1)
		}
	}

	return graph
}

func Entrypoints(graph goraph.Graph) []goraph.Node {
	sp := make([]goraph.Node, 0)

	for _, nd := range graph.GetNodes() {
		s, _ := graph.GetSources(nd.ID())
		if len(s) == 0 {
			sp = append(sp, nd)
		}
	}

	ns := &nodeSorter{
		nodes: sp,
	}
	sort.Sort(ns)

	return ns.nodes
}

func NewBuildStrategy(tasks []TaskNode) BuildStrategy {
	return &buildStrategy{
		graph:      PopulateGraph(tasks),
		builtTasks: make(map[goraph.Node]bool),
	}
}

func (bs buildStrategy) Build(runner func(taskName string) error) error {
	rc := &runnerContext{
		processor: runner,
	}

	for _, entryNode := range Entrypoints(bs.graph) {
		err := bs.processNode(rc, entryNode)
		if err != nil {
			return err
		}
	}

	return nil
}

type runnerContext struct {
	processor func(name string) error
}

func (bs buildStrategy) processNode(rc *runnerContext, node goraph.Node) error {
	sources, err := bs.graph.GetSources(node.ID())
	if err != nil {
		return err
	}

	for _, s := range sources {
		err := bs.processNode(rc, s)
		if err != nil {
			return err
		}
	}

	if _, ok := bs.builtTasks[node]; ok {
		return nil
	}

	err = rc.processor(node.String())
	if err != nil {
		return err
	}
	bs.builtTasks[node] = true

	targets, err := bs.graph.GetTargets(node.ID())
	if err != nil {
		return err
	}

	for _, targetNode := range targets {
		err := bs.processNode(rc, targetNode)
		if err != nil {
			return err
		}
	}

	return nil
}
