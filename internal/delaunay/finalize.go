package delaunay

import "delaunay-bw/internal/geom"

// StripSuper removes every triangle that references one of the virtual
// super-triangle vertices (n, n+1, n+2). After this step the mesh contains
// only cells whose vertices are real input points.
//
// The input points keep indices [0, n); the super vertices are appended at
// the end of the working point array. Because the super vertices never
// collide with real indices, the strip is exact: no real cell is ever removed
// by mistake, and no ghost triangle can survive outside the convex hull.
func StripSuper(tris []Triangle, n int) []Triangle {
	out := make([]Triangle, 0, len(tris))
	for _, t := range tris {
		if t.HasAnyVertex(n, n+1, n+2) {
			continue
		}
		out = append(out, t)
	}
	return out
}

// NormalizeOrientations walks the mesh and reorders any triangle that ended
// up clockwise (a defensive sweep after reconnection, useful when clients
// hand the kernel externally built meshes). Triangles are kept only if they
// are non-degenerate within tolerance.
func NormalizeOrientations(pts []geom.Point, tris []Triangle, tol float64) []Triangle {
	out := make([]Triangle, 0, len(tris))
	scale := geom.ScaleForTolerance(pts)
	areaTol := tol * scale * scale
	for _, t := range tris {
		area := geom.SignedAreaTwice(pts[t[0]], pts[t[1]], pts[t[2]])
		if area < 0 {
			t = Triangle{t[0], t[2], t[1]}
			area = -area
		}
		if area <= areaTol {
			// A degenerate cell adds no area and confuses every downstream
			// statistic; drop it.
			continue
		}
		out = append(out, t)
	}
	return out
}
