package main

import (
	"testing"
)

func TestNeighbors(t *testing.T) {
	g := new(3, 3, 0)
	if len(g.closeCoord(0, 0)) != 3 {
		t.Errorf("wrong number for (0, 0)")
	}
	if len(g.closeCoord(1, 1)) != 8 {
		t.Errorf("wrong number for (1, 1)")
	}
	if len(g.closeCoord(0, 1)) != 5 {
		t.Errorf("wrong number for (1, 1)")
	}
	if len(g.closeCoord(2, 2)) != 3 {
		t.Errorf("wrong number for (2, 2)")
	}
	if len(g.closeCoord(2, 1)) != 5 {
		t.Errorf("wrong number for (2, 1)")
	}
	if len(g.closeCoord(0, 1)) != 5 {
		t.Errorf("wrong number for (0, 1)")
	}
	if len(g.closeCoord(1, 2)) != 5 {
		t.Errorf("wrong number for (1, 2)")
	}
	if len(g.closeCoord(2, 0)) != 3 {
		t.Errorf("wrong number for (2, 0)")
	}
}
