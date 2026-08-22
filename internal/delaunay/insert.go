package delaunay

import "delaunay-bw/internal/geom"

// InsertPoint performs one Bowyer–Watson step for the point p (stored at
// index pidx in pts):
//
//  1. mark every triangle whose circumcircle contains p (the "bad" set),
//  2. extract the cavity boundary from the union of the bad triangles,
//  3. delete the bad triangles,
//  4. reconnect every boundary edge to p.
//
// The cavity boundary is directed so that reconnecting never flips the hole
// winding; as a defensive measure each new triangle is still reoriented by
// its signed area, so a numerically marginal edge cannot create a clockwise
// mesh.
//
// The function always returns a slice with the surviving triangles plus the
// new ones; the caller keeps responsibility for the final strip of
// super-triangle-touching cells.
func InsertPoint(tris []Triangle, pts []geom.Point, p geom.Point, pidx int, tol float64) []Triangle {
	bad := collectBadTriangles(tris, pts, p, tol)
	if len(bad) == 0 {
		// This can only happen for points outside the super triangle, which
		// the construction prevents. Keep the mesh unchanged defensively.
		return tris
	}

	boundary := cavityBoundary(tris, bad)

	out := make([]Triangle, 0, len(tris)-len(bad)+len(boundary))
	for i := range tris {
		if !bad[i] {
			out = append(out, tris[i])
		}
	}
	for _, e := range boundary {
		out = append(out, orientCCW(e.A, e.B, pidx, pts))
	}
	return out
}

// orientCCW returns the triangle (a, b, c) ordered counter-clockwise. When
// the signed area is negative the order is swapped; when it is (numerically)
// zero the original order is kept, because the caller's tolerance already
// rejected truly degenerate configurations.
func orientCCW(a, b, c int, pts []geom.Point) Triangle {
	if geom.SignedAreaTwice(pts[a], pts[b], pts[c]) < 0 {
		return Triangle{a, c, b}
	}
	return Triangle{a, b, c}
}

// InsertAll inserts every point of pts into the mesh one at a time, starting
// from the single super triangle. This is the sequential Bowyer–Watson loop.
// The initial mesh is the triangle of the three appended super vertices.
func InsertAll(pts []geom.Point, super []geom.Point, n int, tol float64) []Triangle {
	all := make([]geom.Point, 0, len(pts)+3)
	all = append(all, pts...)
	all = append(all, super...)

	initial := Triangle{n, n + 1, n + 2}
	tris := []Triangle{initial}

	for i := range pts {
		tris = InsertPoint(tris, all, pts[i], i, tol)
	}
	return tris
}
