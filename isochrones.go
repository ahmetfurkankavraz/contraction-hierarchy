package ch

import (
	"container/heap"
	"fmt"
)

// Isochrones Returns set of vertices and corresponding distances restricted by maximum travel cost for source vertex
// source - source vertex (user defined label)
// maxCost - restriction on travel cost for breadth search
// See ref. https://wiki.openstreetmap.org/wiki/Isochrone and https://en.wikipedia.org/wiki/Isochrone_map
// Note: implemented breadth-first searching path algorithm does not guarantee shortest pathes to reachable vertices (until all edges have cost 1.0). See ref: https://en.wikipedia.org/wiki/Breadth-first_search
// Note: result for estimated costs could be also inconsistent due nature of data structure
func (graph *Graph) Isochrones(source int64, maxCost float64) (map[int64]float64, error) {
	var ok bool
	if source, ok = graph.mapping[source]; !ok {
		return nil, fmt.Errorf("no such source")
	}
	Q := &minHeap{}
	heap.Init(Q)
	distance := make(map[int64]float64, len(graph.Vertices))
	Q.Push(&minHeapVertex{id: source, distance: 0})
	visit := make(map[int64]bool)
	for Q.Len() != 0 {
		next := heap.Pop(Q).(*minHeapVertex)
		visit[next.id] = true
		if next.distance <= maxCost {
			distance[graph.Vertices[next.id].Label] = next.distance
			vertexList := graph.Vertices[next.id].outIncidentEdges
			for i := range vertexList {
				neighbor := vertexList[i].vertexID
				if v1, ok1 := graph.shortcuts[next.id]; ok1 {
					if _, ok2 := v1[neighbor]; ok2 {
						// Ignore shortcut
						continue
					}
				}
				target := vertexList[i].vertexID
				cost := vertexList[i].weight
				alt := distance[graph.Vertices[next.id].Label] + cost
				if visit[target] {
					continue
				}
				Q.Push(&minHeapVertex{id: target, distance: alt})
			}
		}
	}
	return distance, nil
}
