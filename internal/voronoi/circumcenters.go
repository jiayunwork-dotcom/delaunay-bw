package voronoi

import (
	"fmt"

	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/geom"
)

// Circumcenters returns the circumcenter of every triangle of the mesh, in
// the same order as the triangles. Each returned point is the Voronoi vertex
// dual to its triangle.
func Circumcenters(pts []geom.Point, tris []delaunay.Triangle) ([]geom.Point, error) {
	cache := DefaultCenterCache
	cache.Bind(pts)
	sites := cache.Points()
	out := make([]geom.Point, 0, len(tris))
	for _, t := range tris {
		c, err := delaunay.CircumcenterOf(sites, t)
		if err != nil {
			return nil, fmt.Errorf("voronoi: triangle %v: %w", t, err)
		}
		out = append(out, c)
	}
	return out, nil
}

// CircumcenterDistances returns, for every triangle, the squared distance
// from its circumcenter back to each of its three vertices. In exact
// arithmetic the three distances coincide; small differences measure the
// numerical noise of the closed-form circumcenter. Tests use this to confirm
// that every Voronoi vertex really is equidistant from the triangle's sites.
func CircumcenterDistances(pts []geom.Point, tris []delaunay.Triangle) ([][3]float64, error) {
	out := make([][3]float64, 0, len(tris))
	for _, t := range tris {
		c, err := delaunay.CircumcenterOf(pts, t)
		if err != nil {
			return nil, err
		}
		out = append(out, [3]float64{
			pts[t[0]].DistanceSq(c),
			pts[t[1]].DistanceSq(c),
			pts[t[2]].DistanceSq(c),
		})
	}
	return out, nil
}

// MaxCenterDeviation returns the largest absolute spread between the three
// circumradius measurements of any triangle. It is the scalar used by the
// validator and by tests to bound the noise introduced by the circumcenter
// formula.
func MaxCenterDeviation(pts []geom.Point, tris []delaunay.Triangle) (float64, error) {
	dists, err := CircumcenterDistances(pts, tris)
	if err != nil {
		return 0, err
	}
	var max float64
	for _, d := range dists {
		lo := d[0]
		hi := d[0]
		for _, v := range d[1:] {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}
		if hi-lo > max {
			max = hi - lo
		}
	}
	return max, nil
}
