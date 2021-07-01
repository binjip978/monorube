package vector

type Clock struct {
	node   string
	vector map[string]int
}

func (c *Clock) Merge(other *Clock) {
	for k, v := range other.vector {
		cv := c.vector[k]
		if cv > v {
			c.vector[k] = cv
		} else {
			c.vector[k] = v
		}
	}

	c.vector[c.node] += 1
}

func (c *Clock) Less(other *Clock) bool {
	var notEqual bool

	for k, v := range c.vector {
		if v > other.vector[k] {
			return false
		}

		if v != other.vector[k] {
			notEqual = true
		}
	}

	return notEqual
}

func (c *Clock) Equal(other *Clock) bool {
	for k, v := range c.vector {
		if v != other.vector[k] {
			return false
		}
	}

	for k, v := range other.vector {
		if v != c.vector[k] {
			return false
		}
	}

	return true
}

func (c *Clock) Concurrent(other *Clock) bool {
	return !c.Less(other) && !other.Less(c)
}
