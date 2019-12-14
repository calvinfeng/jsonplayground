package main

// Edge defines connection between two nodes.
type Edge struct {
	ID          uint64
	Source      *Node
	Destination *Node
}

// Node has weight, slightly different from conventional approach.
type Node struct {
	ID       uint64
	Incoming *Edge
	Outgoing *Edge
	Weight   uint64
}

// QueueMap is used for the lack of a priority queue map.
type QueueMap map[*Node]uint64

func (q QueueMap) DequeueMin() *Node {
	var selected *Node
	var minCost uint64
	for n, cost := range q {
		if selected == nil && minCost == 0 {
			selected = n
			cost = minCost
			continue
		}

		if minCost < cost {
			selected = n
			cost = minCost
		}
	}

	delete(q, selected)
	return selected
}

// Dijkstra is used to find shortest path but can be inverted to find the longest path between
// two nodes.
func Dijkstra(start *Node) {

}