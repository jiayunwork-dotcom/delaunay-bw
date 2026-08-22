package delaunay

import (
	"errors"
	"math"
	"math/rand"
	"testing"

	"delaunay-bw/internal/geom"
)

// jitteredGrid builds a cols x rows grid of points with a fixed-seed jitter
// large enough to avoid coincidences and cocircularities.
func jitteredGrid(cols, rows int, seed int64) []geom.Point {
	r := rand.New(rand.NewSource(seed))
	pts := make([]geom.Point, 0, cols*rows)
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			pts = append(pts, geom.Point{
				X: float64(i) + r.Float64()*0.3,
				Y: float64(j) + r.Float64()*0.3,
			})
		}
	}
	return pts
}

// regularHull builds the vertices of a regular n-gon centered at the origin.
// Every vertex lies on a common circle, which exercises the cocircular
// tolerance path of the kernel.
func regularHull(n int, radius float64) []geom.Point {
	pts := make([]geom.Point, n)
	for k := 0; k < n; k++ {
		a := 2 * math.Pi * float64(k) / float64(n)
		pts[k] = geom.Point{X: radius * math.Cos(a), Y: radius * math.Sin(a)}
	}
	return pts
}

// mustTriangulate runs the kernel and fails the test on any error.
func mustTriangulate(t *testing.T, pts []geom.Point) *Result {
	t.Helper()
	res, err := Triangulate(pts, Options{Tolerance: 1e-9})
	if err != nil {
		t.Fatalf("Triangulate(%d points): unexpected error: %v", len(pts), err)
	}
	return res
}

// TestTriangulateTooFewPoints pins the simplest failure rule: a triangle
// needs at least three points, so a set of one or two points is rejected.
func TestTriangulateTooFewPoints(t *testing.T) {
	cases := []struct {
		name string
		pts  []geom.Point
	}{
		{"zero points", nil},
		{"one point", []geom.Point{{X: 1, Y: 1}}},
		{"two points", []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 1}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Triangulate(tc.pts, Options{})
			if !errors.Is(err, geom.ErrTooFewPoints) {
				t.Errorf("Triangulate(%s): got err %v, want %v", tc.name, err, geom.ErrTooFewPoints)
			}
		})
	}
}

// TestTriangulateDuplicatePoint pins the duplicate-site rule: two points that
// coincide within tolerance collapse the cavity and are rejected.
func TestTriangulateDuplicatePoint(t *testing.T) {
	pts := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0.5, Y: 0.866},
		{X: 0.9999999999, Y: 0}, // within 1e-9 relative tolerance of (1,0)
	}
	_, err := Triangulate(pts, Options{})
	if !errors.Is(err, geom.ErrDuplicatePoint) {
		t.Errorf("got err %v, want duplicate-point error", err)
	}
}

// TestTriangulateCollinear pins the collinear-set rule: points that all lie
// on one line span no area and cannot be triangulated.
func TestTriangulateCollinear(t *testing.T) {
	pts := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0.001},
		{X: 2, Y: 0.002},
		{X: 3, Y: 0.003},
	}
	_, err := Triangulate(pts, Options{})
	if !errors.Is(err, geom.ErrCollinearSet) {
		t.Errorf("got err %v, want collinear-set error", err)
	}
}

// TestTriangulateNaN pins the non-finite rule: NaN coordinates make every
// geometric predicate meaningless and are rejected before insertion.
func TestTriangulateNaN(t *testing.T) {
	pts := []geom.Point{
		{X: 0, Y: 0},
		{X: math.NaN(), Y: 1},
		{X: 1, Y: 0},
	}
	_, err := Triangulate(pts, Options{})
	if !errors.Is(err, geom.ErrNonFinite) {
		t.Errorf("got err %v, want non-finite error", err)
	}
}

// TestEmptyCircleProperty verifies the core invariant: after triangulation no
// input point lies strictly inside the circumcircle of any triangle.
func TestEmptyCircleProperty(t *testing.T) {
	for _, seed := range []int64{1, 2, 3, 4} {
		pts := jitteredGrid(6, 5, seed)
		res := mustTriangulate(t, pts)
		if !IsDelaunay(pts, res.Triangles, 1e-9) {
			t.Errorf("seed %d: empty-circle property violated on %d triangles", seed, len(res.Triangles))
		}
	}
}

// TestGridJitterFormula pins the Euler formula T = 2n - 2 - h on the shipped
// example: a 5x4 jittered grid of 20 points whose convex hull has 14
// vertices must yield exactly 24 triangles, and the mesh must stay Delaunay.
func TestGridJitterFormula(t *testing.T) {
	pts, err := loadExamplePoints("../../example/grid-jitter.json")
	if err != nil {
		t.Fatalf("load example: %v", err)
	}
	if len(pts) != 20 {
		t.Fatalf("example has %d points, want 20", len(pts))
	}
	res := mustTriangulate(t, pts)
	if res.Stats.NumTriangles != 24 {
		t.Errorf("triangle count = %d, want 24 (2*20-2-14)", res.Stats.NumTriangles)
	}
	if res.Stats.HullSize != 14 {
		t.Errorf("hull size = %d, want 14", res.Stats.HullSize)
	}
	if !res.Stats.MatchesFormula {
		t.Errorf("formula mismatch: expected %d triangles by 2n-2-h", res.Stats.ExpectFormula)
	}
	if !IsDelaunay(pts, res.Triangles, 1e-9) {
		t.Error("example triangulation is not Delaunay")
	}
}

// TestInsertPointAddsTwoTriangles pins the medium rule: inserting one point
// strictly inside the convex hull increases the triangle count by exactly two
// (2n-2-h grows by 2 because the hull size is unchanged).
func TestInsertPointAddsTwoTriangles(t *testing.T) {
	base := regularHull(6, 1.0) // hexagon, hull size 6, 2*6-2-6 = 4 triangles
	baseRes := mustTriangulate(t, base)
	if baseRes.Stats.NumTriangles != 4 {
		t.Fatalf("hexagon triangulation has %d triangles, want 4", baseRes.Stats.NumTriangles)
	}

	interior := geom.Point{X: 0, Y: 0} // center is strictly inside the hull
	extended := append([]geom.Point{}, base...)
	extended = append(extended, interior)

	extRes, err := Triangulate(extended, Options{})
	if err != nil {
		t.Fatalf("Triangulate with interior point: %v", err)
	}
	if diff := CompareTriangleCounts(baseRes, extRes); diff != 2 {
		t.Errorf("triangle delta after inserting interior point = %d, want 2", diff)
	}
	if extRes.Stats.HullSize != 6 {
		t.Errorf("hull size changed to %d, want 6", extRes.Stats.HullSize)
	}
}

// TestRigidTranslationInvariance pins the invariance rule: a rigid translation
// of the whole point set leaves the combinatorial structure unchanged and
// translates every circumcenter by the same vector.
func TestRigidTranslationInvariance(t *testing.T) {
	orig := jitteredGrid(5, 4, 7)
	base := mustTriangulate(t, orig)

	shift := geom.Point{X: 3.7, Y: -2.1}
	translated := make([]geom.Point, len(orig))
	for i, p := range orig {
		translated[i] = p.Add(shift)
	}
	moved := mustTriangulate(t, translated)

	if len(base.Triangles) != len(moved.Triangles) {
		t.Fatalf("triangle count changed: %d vs %d", len(base.Triangles), len(moved.Triangles))
	}
	// The vertex-index triples must be identical as sets.
	a := triangleSet(base.Triangles)
	b := triangleSet(moved.Triangles)
	for key := range a {
		if !b[key] {
			t.Errorf("triangle %v missing after translation (topology changed)", key)
		}
	}

	// Circumcenters must follow the translation. Match triangles by their
	// vertex set (the mesh is built with map iteration, so triangle order
	// between two runs is not guaranteed).
	centers := make(map[[3]int]geom.Point, len(moved.Triangles))
	for i := range moved.Triangles {
		c, err := CircumcenterOf(translated, moved.Triangles[i])
		if err != nil {
			t.Fatalf("circumcenter: %v", err)
		}
		centers[moved.Triangles[i].key()] = c
	}
	for i := range base.Triangles {
		c0, err := CircumcenterOf(orig, base.Triangles[i])
		if err != nil {
			t.Fatalf("circumcenter: %v", err)
		}
		c1, ok := centers[base.Triangles[i].key()]
		if !ok {
			t.Fatalf("moved mesh is missing triangle %v", base.Triangles[i])
		}
		want := c0.Add(shift)
		if c1.Distance(want) > 1e-8 {
			t.Errorf("circumcenter of %v did not follow translation: got %s want %s",
				base.Triangles[i], c1.String(), want.String())
		}
	}
}

// TestCocircularStability pins the stability rule for four cocircular points:
// the kernel must terminate with a legal mesh whose triangles are all valid
// (no flipping oscillation), and the empty-circle property must hold within
// tolerance.
func TestCocircularStability(t *testing.T) {
	// Square corners are exactly cocircular.
	square := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
	}
	res := mustTriangulate(t, square)
	if len(res.Triangles) != 2 {
		t.Errorf("square triangulation has %d triangles, want 2", len(res.Triangles))
	}
	if err := ValidateTopology(square, res.Triangles, 1e-9); err != nil {
		t.Errorf("square mesh topology: %v", err)
	}
	if !IsDelaunay(square, res.Triangles, 1e-9) {
		t.Error("square mesh is not Delaunay within tolerance")
	}

	// Add the center point: also on the original circumcircle, still stable.
	withCenter := append([]geom.Point{}, square...)
	withCenter = append(withCenter, geom.Point{X: 0.5, Y: 0.5})
	res2 := mustTriangulate(t, withCenter)
	if len(res2.Triangles) != 4 {
		t.Errorf("square+center triangulation has %d triangles, want 4", len(res2.Triangles))
	}
	if !IsDelaunay(withCenter, res2.Triangles, 1e-9) {
		t.Error("square+center mesh is not Delaunay within tolerance")
	}

	// A regular hexagon is fully cocircular; inserting its center exercises
	// the all-bad-triangles cavity path.
	hex := regularHull(6, 1.0)
	res3 := mustTriangulate(t, hex)
	if len(res3.Triangles) != 4 {
		t.Errorf("hexagon has %d triangles, want 4", len(res3.Triangles))
	}
}

// TestNoGhostTrianglesOutsideConvexHull is the hardest invariant: after the
// super triangle is stripped, no ghost cell may survive outside the convex
// hull. This is verified three ways: every vertex index refers to a real
// input point, the mesh tiles the hull exactly (area conservation), and the
// topology is a closed manifold patch.
func TestNoGhostTrianglesOutsideConvexHull(t *testing.T) {
	pts := jitteredGrid(7, 6, 99) // 42 points with a nontrivial hull
	res := mustTriangulate(t, pts)

	for _, tri := range res.Triangles {
		for _, v := range tri {
			if v < 0 || v >= len(pts) {
				t.Fatalf("ghost triangle %v references virtual vertex %d", tri, v)
			}
		}
	}

	if err := ValidateTopology(pts, res.Triangles, 1e-9); err != nil {
		t.Fatalf("mesh topology invalid: %v", err)
	}
	if err := ValidateCoverage(pts, res.Triangles, res.Hull, 1e-9); err != nil {
		t.Fatalf("area coverage invalid: %v", err)
	}

	hullArea := geom.HullArea(pts, res.Hull)
	if math.Abs(hullArea-res.Stats.TriangleArea) > 1e-6 {
		t.Errorf("triangle area sum %.6g does not match hull area %.6g",
			res.Stats.TriangleArea, hullArea)
	}
}

// TestStripSuperTriangle pins the stripping step directly: cells touching the
// appended virtual vertices must be removed and only real cells must remain.
func TestStripSuperTriangle(t *testing.T) {
	tris := []Triangle{
		{0, 1, 2},   // real cell
		{1, 2, 3},   // real cell
		{0, 1, 20},  // touches super vertex 20
		{4, 21, 22}, // touches super vertices 21, 22
		{2, 3, 21},  // touches super vertex 21
	}
	got := StripSuper(tris, 20)
	if len(got) != 2 {
		t.Fatalf("StripSuper kept %d cells, want 2", len(got))
	}
	for _, tri := range got {
		if tri.HasAnyVertex(20, 21, 22) {
			t.Errorf("super vertex leaked into %v", tri)
		}
	}
}

// TestMeshTopology pins the manifold checks on a nontrivial mesh.
func TestMeshTopology(t *testing.T) {
	pts := jitteredGrid(6, 6, 11)
	res := mustTriangulate(t, pts)
	if err := ValidateTopology(pts, res.Triangles, 1e-9); err != nil {
		t.Fatalf("topology: %v", err)
	}
	if err := ValidateDelaunay(pts, res.Triangles, 1e-9); err != nil {
		t.Fatalf("delaunay: %v", err)
	}
}
