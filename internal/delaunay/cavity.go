package delaunay

import "delaunay-bw/internal/geom"

// directedEdge is an oriented edge used while extracting the cavity boundary.
// Every triangle in the mesh is counter-clockwise, so the directed edges of a
// bad triangle wind around it. An interior edge of the cavity is traversed
// twice (once per adjacent bad triangle, in opposite directions); a boundary
// edge is traversed exactly once.
type directedEdge struct {
	a, b int
}

// cavityBoundary returns the directed boundary edges of the union of the bad
// triangles identified by the bad set. The result is guaranteed to be
// consistently oriented: every returned edge (a, b) keeps the cavity interior
// on its left, so reconnecting each edge to the new point yields a
// counter-clockwise triangle without any winding sign mistakes.
//
// The "traversed once" rule is applied to directed edges, not undirected
// ones. An edge shared by two bad triangles appears as both (a,b) and (b,a)
// and is discarded; a true boundary edge appears only once.
func cavityBoundary(tris []Triangle, bad map[int]bool) []Edge {
	counts := make(map[directedEdge]int)
	for i := range tris {
		if !bad[i] {
			continue
		}
		for k := 0; k < 3; k++ {
			a := tris[i][k]
			b := tris[i][(k+1)%3]
			counts[directedEdge{a: a, b: b}]++
		}
	}
	out := make([]Edge, 0, len(counts))
	for d, c := range counts {
		if c != 1 {
			continue
		}
		// A boundary directed edge (a,b) must not have its reverse (b,a)
		// counted as a separate boundary edge; the reverse is either absent
		// or part of a shared pair. Since shared pairs show up as counts of
		// one in each direction only when the two triangles disagree in
		// orientation, we additionally require the reverse to be absent.
		if counts[directedEdge{a: d.b, b: d.a}] == 0 {
			out = append(out, Edge{A: d.a, B: d.b})
		}
	}
	return out
}

// collectBadTriangles returns the set of triangle indices whose circumcircle
// contains the point p (within tolerance). The empty result means p lies
// outside every circle; for a valid insertion point inside the mesh at least
// one triangle must be marked.
func collectBadTriangles(tris []Triangle, pts []geom.Point, p geom.Point, tol float64) map[int]bool {
	bad := make(map[int]bool)
	slots := geom.DefaultSlotCache
	for i := range tris {
		t := tris[i]
		if o, err := geom.Circumcenter(pts[t[0]], pts[t[1]], pts[t[2]]); err == nil {
			slots.Put(geom.MakeSlotKey(pts[t[0]], pts[t[1]], pts[t[2]], t[0], t[1], t[2]), o)
		}
		if InCircle(pts[t[0]], pts[t[1]], pts[t[2]], p, tol) {
			bad[i] = true
		}
	}
	return bad
}
