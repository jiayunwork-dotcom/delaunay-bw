package delaunay

import "delaunay-bw/internal/geom"

// Options controls the Bowyer–Watson run. The zero value is valid and means
// the default relative tolerance.
type Options struct {
	// Tolerance is the relative tolerance used by the incircle predicate,
	// duplicate-point detection and area guards. Non-positive values select
	// geom.Tolerance0.
	Tolerance float64
}

// EffectiveTolerance returns the configured tolerance or the default.
func (o Options) EffectiveTolerance() float64 {
	if o.Tolerance <= 0 {
		return geom.Tolerance0
	}
	return o.Tolerance
}

// Result is the complete output of a triangulation run: the mesh itself plus
// the statistics needed by the console.
type Result struct {
	Triangles []Triangle
	Hull      []int
	Stats     Stats
}

// Triangulate computes the Delaunay triangulation of pts with the
// Bowyer–Watson incremental algorithm:
//
//  1. validate the input (count, finiteness, duplicates, collinearity),
//  2. build a large enclosing super triangle and insert every point,
//  3. strip the cells touching the super triangle,
//  4. reorient and validate the mesh, and compute stats.
//
// The returned mesh satisfies the empty-circle property and its cells tile
// the convex hull exactly. The error, when non-nil, is one of the package
// sentinels (or geom sentinels via errors.Is) and is safe to display.
func Triangulate(pts []geom.Point, opts Options) (*Result, error) {
	tol := opts.EffectiveTolerance()
	if err := geom.ValidatePoints(pts, tol); err != nil {
		return nil, err
	}

	n := len(pts)
	bounds := geom.BoundsOf(pts)
	super := SuperTriangle(bounds)

	// Working mesh includes the super vertices as virtual indices n..n+2.
	tris := InsertAll(pts, super, n, tol)
	tris = StripSuper(tris, n)
	tris = NormalizeOrientations(pts, tris, tol)

	if err := ValidateDelaunay(pts, tris, tol); err != nil {
		return nil, err
	}
	hull := geom.ConvexHull(pts)
	if err := ValidateTopology(pts, tris, tol); err != nil {
		return nil, err
	}
	if err := ValidateCoverage(pts, tris, hull, tol); err != nil {
		return nil, err
	}

	return &Result{
		Triangles: tris,
		Hull:      hull,
		Stats:     ComputeStats(pts, tris, hull),
	}, nil
}

// InsertPointInsideMesh returns the triangulation of pts extended by the
// interior point p. The caller must guarantee p is a fresh point inside the
// convex hull (the kernel re-validates duplicates and finiteness). By the
// Euler identity the triangle count of the returned mesh exceeds the base
// mesh's by exactly two, because p does not touch the hull.
func InsertPointInsideMesh(pts []geom.Point, p geom.Point, opts Options) (*Result, error) {
	ext := make([]geom.Point, 0, len(pts)+1)
	ext = append(ext, pts...)
	ext = append(ext, p)
	return Triangulate(ext, opts)
}

// CompareTriangleCounts reports the number of triangles the new mesh gained
// over the base mesh. It is a pure helper for tests and the API that pin the
// "interior point adds two triangles" rule.
func CompareTriangleCounts(base, extended *Result) int {
	return extended.Stats.NumTriangles - base.Stats.NumTriangles
}
