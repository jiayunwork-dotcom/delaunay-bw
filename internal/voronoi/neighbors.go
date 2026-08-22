package voronoi

import "delaunay-bw/internal/delaunay"

// neighborMap returns, for every triangle index, the indices of the triangles
// that share an edge with it. A triangle on the convex hull boundary has
// fewer neighbors than edges; the interior ones always have three.
func neighborMap(tris []delaunay.Triangle) [][]int {
	edgeToTris := make(map[delaunay.Edge][]int)
	for i := range tris {
		for k := 0; k < 3; k++ {
			e := delaunay.EdgeOf(tris[i], k)
			edgeToTris[e] = append(edgeToTris[e], i)
		}
	}

	neighbors := make([][]int, len(tris))
	for _, list := range edgeToTris {
		if len(list) < 2 {
			continue
		}
		for _, a := range list {
			for _, b := range list {
				if a != b {
					neighbors[a] = append(neighbors[a], b)
				}
			}
		}
	}
	return neighbors
}

// NeighborCounts returns, for each triangle, how many neighbors it has.
// The Voronoi vertex of a boundary triangle is only connected to its interior
// neighbors; the dual edge count is the sum of interior shared edges.
func NeighborCounts(neighbors [][]int) []int {
	out := make([]int, len(neighbors))
	for i, ns := range neighbors {
		out[i] = len(ns)
	}
	return out
}

// InteriorEdgeCount returns how many edges of the mesh are shared by exactly
// two triangles. Each such edge corresponds to one Voronoi dual edge. The
// count is (3T - h)/2, where T is the triangle count and h the hull size,
// because the hull contributes exactly h boundary edges.
func InteriorEdgeCount(T, h int) int {
	return (3*T - h) / 2
}
