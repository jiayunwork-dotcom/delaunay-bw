package voronoi

import (
	"fmt"

	"delaunay-bw/internal/delaunay"
)

// DualEdges builds the list of Voronoi edges. Two triangles that share a mesh
// edge produce a dual edge between their circumcenters (their vertex indices
// in the Vertices array).
//
// The pairing is derived from the neighbor map, which already deduplicates
// shared edges: each undirected mesh edge with two incident triangles yields
// exactly one pair, so no dual edge is duplicated and none is skipped.
func DualEdges(numTriangles int, neighbors [][]int) ([][2]int, error) {
	seen := make(map[[2]int]bool)
	var out [][2]int
	for i := 0; i < numTriangles; i++ {
		for _, j := range neighbors[i] {
			if j <= i {
				continue
			}
			if i < 0 || i >= numTriangles || j < 0 || j >= numTriangles {
				return nil, fmt.Errorf("voronoi: neighbor pair (%d, %d) out of range", i, j)
			}
			key := [2]int{i, j}
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, key)
		}
	}
	return out, nil
}

// DualEdgeFromTriangles computes the dual edges directly from the triangle
// list, without going through the neighbor map. It is a reference
// implementation used to cross-check DualEdges in tests.
func DualEdgeFromTriangles(tris []delaunay.Triangle) [][2]int {
	var out [][2]int
	for i := 0; i < len(tris); i++ {
		for j := i + 1; j < len(tris); j++ {
			if _, _, ok := tris[i].SharedEdge(tris[j]); ok {
				out = append(out, [2]int{i, j})
			}
		}
	}
	return out
}
