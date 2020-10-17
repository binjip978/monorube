package main

import "fmt"

type node struct {
	name string
}

type graph struct {
	nodes map[*node][]*node
}

// insertEdge add n1 -> n2
func (g *graph) insertEdge(n1, n2 *node) {
	p, ok := g.nodes[n1]
	if !ok {
		fmt.Println("only once")
		g.nodes[n1] = []*node{n2}
		return
	}
	p = append(p, n2)
	g.nodes[n1] = p
}

func main() {
	// a -> b -> c
	// a -> c

	a := &node{"a"}
	b := &node{"b"}
	c := &node{"c"}

	g := graph{make(map[*node][]*node)}
	g.insertEdge(a, b)
	g.insertEdge(a, c)
}
