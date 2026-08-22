package voronoi

import (
	"sort"

	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/geom"
)

// EdgeFromSitePair builds the Voronoi edge between the two sites whose
// Delaunay edge crosses it: for a mesh edge (a, b) shared by triangles i and
// j, the Voronoi edge connects vertex i to vertex j. The returned geometry is
// the segment between the two circumcenters.
func EdgeFromSitePair(verts []geom.Point, a, b int) [2]geom.Point {
	return [2]geom.Point{verts[a], verts[b]}
}

// EdgeMidpoint returns the midpoint of a Voronoi edge, which lies on the
// shared Delaunay edge of the two incident sites.
func EdgeMidpoint(verts []geom.Point, e [2]int) geom.Point {
	p := verts[e[0]]
	q := verts[e[1]]
	return geom.Point{X: 0.5 * (p.X + q.X), Y: 0.5 * (p.Y + q.Y)}
}

// SortEdges returns a copy of the edges sorted first by the first endpoint
// and then by the second. It gives callers a deterministic order for
// comparing diagrams between runs or implementations.
func SortEdges(edges [][2]int) [][2]int {
	out := make([][2]int, len(edges))
	copy(out, edges)
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

// BuildDelaunayEdges enumerates the unique undirected edges of the mesh,
// which are exactly the edges crossed by the Voronoi edges. It is exposed so
// tests can verify the crossing relationship (each Voronoi edge crosses the
// corresponding Delaunay edge).
func BuildDelaunayEdges(tris []delaunay.Triangle) []delaunay.Edge {
	seen := make(map[delaunay.Edge]bool)
	var out []delaunay.Edge
	for _, t := range tris {
		for k := 0; k < 3; k++ {
			e := delaunay.EdgeOf(t, k)
			if !seen[e] {
				seen[e] = true
				out = append(out, e)
			}
		}
	}
	return out
}

// Crosses reports whether the segment between two Voronoi vertices and the
// segment between the two input points share a point, within the tolerance.
// It is a helper for tests that verify the duality geometrically.
func Crosses(a, b geom.Point, p, q geom.Point, tol float64) bool {
	// Orientations of the four points around each segment.
	d1 := geom.CrossFrom(p, q, a)
	d2 := geom.CrossFrom(p, q, b)
	d3 := geom.CrossFrom(a, b, p)
	d4 := geom.CrossFrom(a, b, q)
	opp1 := (d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)
	opp2 := (d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)
	if opp1 && opp2 {
		return true
	}
	// Degenerate cases: endpoints touching. Approximate by distance.
	if a.Distance(p) <= tol || a.Distance(q) <= tol || b.Distance(p) <= tol || b.Distance(q) <= tol {
		return true
	}
	return false
}
