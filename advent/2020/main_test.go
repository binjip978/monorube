package main

import (
	"fmt"
	"testing"
)

func TestBagParser(t *testing.T) {
	var cases = []struct {
		input string
		from  string
		to    []bagEdge
	}{
		{
			"light red bags contain 1 bright white bag, 2 muted yellow bags.",
			"light red",
			[]bagEdge{{"bright white", 1}, {"muted yellow", 2}},
		},
		{
			"dark orange bags contain 3 bright white bags, 4 muted yellow bags.",
			"dark orange",
			[]bagEdge{{"bright white", 3}, {"muted yellow", 4}},
		},
		{
			"bright white bags contain 1 shiny gold bag.",
			"bright white",
			[]bagEdge{{"shiny gold", 1}},
		},
		{
			"muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.",
			"muted yellow",
			[]bagEdge{{"shiny gold", 2}, {"faded blue bags", 9}},
		},
		{
			"shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.",
			"shiny gold",
			[]bagEdge{{"dark olive", 2}, {"vibrant plum", 2}},
		},
		{
			"dark olive bags contain 3 faded blue bags, 4 dotted black bags.",
			"dark olive",
			[]bagEdge{{"faded blue", 4}, {"dotted black", 2}},
		},
		{
			"vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.",
			"vibrant plum",
			[]bagEdge{{"faded blue", 5}, {"dotted black", 6}},
		},
		{
			"faded blue bags contain no other bags.",
			"faded blue",
			[]bagEdge{},
		},
		{
			"dotted black bags contain no other bags.",
			"dotted black",
			[]bagEdge{},
		},
	}

	for _, tt := range cases {
		fmt.Println(tt.input)
	}
}
