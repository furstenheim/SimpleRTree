package SimpleRTree

type computeDistanceItem struct {
	node *Node
	mind, maxd float64
}

// insertion sort (max possible size = 9)
func insertNode (cs *[MAX_POSSIBLE_SIZE]computeDistanceItem, c computeDistanceItem, index int8) {
	cs[index] = c
	maxd := c.maxd
	for i := index - 1; i >= 0; i-- {
		if cs[i].maxd > maxd {
			cs[i], cs[i + 1] = cs[i + 1], cs[i]
		} else {
			break
		}
	}
}