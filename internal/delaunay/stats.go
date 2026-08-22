package delaunay

import (
	"delaunay-bw/internal/geom"
)

// Stats summarizes a triangulation and its relation to the input point set.
// The ExpectTriangles field encodes the formula that the user-facing console
// displays: for n points with a convex hull of h vertices, a planar
// triangulation always has exactly 2n-2-h triangles. The kernel exposes the
// check so tests and the API can pin the formula directly.
type Stats struct {
	NumPoints     int     `json:"num_points"`
	NumTriangles  int     `json:"num_triangles"`
	HullSize      int     `json:"hull_size"`
	HullArea      float64 `json:"hull_area"`
	TriangleArea  float64 `json:"triangle_area"`
	AreaDelta     float64 `json:"area_delta"`
	ExpectFormula int     `json:"expect_triangles_by_formula"`
	MatchesFormula bool   `json:"matches_formula"`
}

// ComputeStats fills a Stats from the point set, the mesh and its hull. The
// formula T = 2n - 2 - h is an Euler characteristic identity for any planar
// triangulation of a point set whose hull has h vertices and no collinear
// hull points; the hull returned by geom.ConvexHull has exactly that shape.
func ComputeStats(pts []geom.Point, tris []Triangle, hull []int) Stats {
	s := Stats{
		NumPoints:    len(pts),
		NumTriangles: len(tris),
		HullSize:     len(hull),
		HullArea:     geom.HullArea(pts, hull),
	}
	for _, t := range tris {
		s.TriangleArea += geom.TriangleArea(pts[t[0]], pts[t[1]], pts[t[2]])
	}
	d := s.HullArea - s.TriangleArea
	if d < 0 {
		d = -d
	}
	s.AreaDelta = d
	s.ExpectFormula = 2*len(pts) - 2 - len(hull)
	s.MatchesFormula = s.NumTriangles == s.ExpectFormula
	return s
}

// ExpectedTriangleCount returns the formula value 2n - 2 - h. It exists as a
// pure function so callers can compare against it without building a Stats.
func ExpectedTriangleCount(n, h int) int {
	return 2*n - 2 - h
}

// InsertedTriangleDelta is the number of triangles added when a point is
// inserted strictly inside the convex hull of an existing mesh: exactly two,
// by the same Euler identity (the point count rises by one while the hull
// vertex count is unchanged, so T increases by 2).
const InsertedTriangleDelta = 2
