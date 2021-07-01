package vector

import "testing"

func TestLess(t *testing.T) {
	var tests = []struct {
		v1   *Clock
		v2   *Clock
		less bool
	}{
		{
			&Clock{"1", map[string]int{"1": 1, "2": 1, "3": 1}},
			&Clock{"2", map[string]int{"1": 2, "2": 2, "3": 3}},
			true,
		},
		{
			&Clock{"2", map[string]int{"1": 2, "2": 2, "3": 3}},
			&Clock{"1", map[string]int{"1": 1, "2": 1, "3": 1}},
			false,
		},
		{
			&Clock{"1", map[string]int{"1": 2, "2": 2, "3": 3}},
			&Clock{"2", map[string]int{"1": 3, "2": 4, "3": 1}},
			false,
		},
		{
			&Clock{"1", map[string]int{"1": 1, "2": 2, "3": 3}},
			&Clock{"2", map[string]int{"1": 1, "2": 2, "3": 3}},
			false,
		},
	}

	for _, test := range tests {
		if test.less != test.v1.Less(test.v2) {
			t.Errorf("less is not correct for %+v, %+v", test.v1, test.v2)
		}
	}
}

func TestConcurrent(t *testing.T) {
	var tests = []struct {
		v1   *Clock
		v2   *Clock
		conc bool
	}{
		{
			&Clock{"1", map[string]int{"1": 2, "2": 2, "3": 3}},
			&Clock{"2", map[string]int{"1": 3, "2": 4, "3": 1}},
			true,
		},
		{
			&Clock{"1", map[string]int{"1": 1, "2": 1, "3": 1}},
			&Clock{"2", map[string]int{"1": 2, "2": 2, "3": 3}},
			false,
		},
	}

	for _, test := range tests {
		if test.conc != test.v1.Concurrent(test.v2) {
			t.Errorf("concurrent is not correct for %+v, %+v", test.v1, test.v2)
		}
	}
}

func TestEqual(t *testing.T) {
	var tests = []struct {
		v1 *Clock
		v2 *Clock
		eq bool
	}{
		{
			&Clock{"1", map[string]int{}},
			&Clock{"2", map[string]int{}},
			true,
		},
		{
			&Clock{"1", map[string]int{"1": 1, "2": 1, "3": 1}},
			&Clock{"1", map[string]int{"1": 1, "2": 1, "3": 1}},
			true,
		},
		{
			&Clock{"1", map[string]int{"1": 0}},
			&Clock{"2", map[string]int{"2": 0}},
			true,
		},
		{
			&Clock{"1", map[string]int{"1": 1}},
			&Clock{"2", map[string]int{"2": 1}},
			false,
		},
	}

	for _, test := range tests {
		if test.eq != test.v1.Equal(test.v2) {
			t.Errorf("equal is not correct for %+v, %+v", test.v1, test.v2)
		}
	}
}
