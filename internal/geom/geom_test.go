package geom

import (
	"errors"
	"math"
	"testing"
)

// TestConvexHullSquare pins the hull of four square corners: all four points
// are hull vertices in counter-clockwise order.
func TestConvexHullSquare(t *testing.T) {
	pts := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
	}
	hull := ConvexHull(pts)
	if len(hull) != 4 {
		t.Fatalf("square hull has %d vertices, want 4: %v", len(hull), hull)
	}
	if !IsCCW(pts[hull[0]], pts[hull[1]], pts[hull[2]], 1e-12) {
		t.Errorf("hull is not counter-clockwise")
	}
	if HullArea(pts, hull) != 1 {
		t.Errorf("square hull area = %g, want 1", HullArea(pts, hull))
	}
}

// TestConvexHullDropsCollinear pins that Andrew's chain drops points lying
// exactly on a hull edge, so the returned hull has only true vertices.
func TestConvexHullDropsCollinear(t *testing.T) {
	pts := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 1, Y: 0}, // on edge (0,0)-(2,0)
		{X: 0, Y: 2},
		{X: 2, Y: 2},
	}
	hull := ConvexHull(pts)
	if len(hull) != 4 {
		t.Fatalf("hull has %d vertices, want 4 (collinear mid-edge point dropped): %v", len(hull), hull)
	}
}

// TestCircumcenterRightTriangle pins the circumcenter of a right triangle at
// the hypotenuse midpoint.
func TestCircumcenterRightTriangle(t *testing.T) {
	o, err := Circumcenter(Point{X: 0, Y: 0}, Point{X: 2, Y: 0}, Point{X: 0, Y: 2})
	if err != nil {
		t.Fatalf("circumcenter: %v", err)
	}
	if o.Distance(Point{X: 1, Y: 1}) > 1e-12 {
		t.Errorf("circumcenter = %s, want (1, 1)", o.String())
	}
	r0 := Point{X: 0, Y: 0}.DistanceSq(o)
	r1 := Point{X: 2, Y: 0}.DistanceSq(o)
	r2 := Point{X: 0, Y: 2}.DistanceSq(o)
	if math.Abs(r0-r1) > 1e-12 || math.Abs(r0-r2) > 1e-12 {
		t.Errorf("circumcenter not equidistant (r2=%g %g %g)", r0, r1, r2)
	}
}

// TestCircumcenterEquidistant checks that the closed-form center keeps all
// three vertices at the same distance (within floating point noise).
func TestCircumcenterEquidistant(t *testing.T) {
	pts := []Point{{X: 0.3, Y: 1.2}, {X: 2.7, Y: -0.4}, {X: -1.1, Y: 0.6}}
	o, err := Circumcenter(pts[0], pts[1], pts[2])
	if err != nil {
		t.Fatalf("circumcenter: %v", err)
	}
	r0 := pts[0].DistanceSq(o)
	r1 := pts[1].DistanceSq(o)
	r2 := pts[2].DistanceSq(o)
	maxDev := math.Max(math.Abs(r0-r1), math.Abs(r0-r2))
	if maxDev > 1e-9 {
		t.Errorf("circumcenter deviation %g too large (r2=%g %g %g)", maxDev, r0, r1, r2)
	}
}

// TestCircumcenterCollinear pins that collinear input yields ErrCollinear.
func TestCircumcenterCollinear(t *testing.T) {
	_, err := Circumcenter(Point{X: 0, Y: 0}, Point{X: 1, Y: 1}, Point{X: 2, Y: 2})
	if !errors.Is(err, ErrCollinear) {
		t.Errorf("got %v, want ErrCollinear", err)
	}
}

// TestValidatePointsTooFew pins the minimum-count rule.
func TestValidatePointsTooFew(t *testing.T) {
	if err := ValidatePoints(nil, 1e-9); !errors.Is(err, ErrTooFewPoints) {
		t.Errorf("empty set: got %v, want ErrTooFewPoints", err)
	}
	if err := ValidatePoints([]Point{{X: 1, Y: 1}, {X: 2, Y: 2}}, 1e-9); !errors.Is(err, ErrTooFewPoints) {
		t.Errorf("pair: got %v, want ErrTooFewPoints", err)
	}
}

// TestValidatePointsDuplicate pins the duplicate detection.
func TestValidatePointsDuplicate(t *testing.T) {
	pts := []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0.5, Y: 0.9}, {X: 1e-12, Y: 0}}
	if err := ValidatePoints(pts, 1e-9); !errors.Is(err, ErrDuplicatePoint) {
		t.Errorf("got %v, want ErrDuplicatePoint", err)
	}
}

// TestValidatePointsNonFinite pins the NaN/Inf rule.
func TestValidatePointsNonFinite(t *testing.T) {
	pts := []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: math.NaN(), Y: 1}}
	if err := ValidatePoints(pts, 1e-9); !errors.Is(err, ErrNonFinite) {
		t.Errorf("NaN: got %v, want ErrNonFinite", err)
	}
	pts2 := []Point{{X: 0, Y: 0}, {X: math.Inf(1), Y: 1}, {X: 1, Y: 0}}
	if err := ValidatePoints(pts2, 1e-9); !errors.Is(err, ErrNonFinite) {
		t.Errorf("Inf: got %v, want ErrNonFinite", err)
	}
}

// TestValidatePointsCollinear pins the collinear-set rule.
func TestValidatePointsCollinear(t *testing.T) {
	pts := []Point{{X: 0, Y: 0}, {X: 1, Y: 0.001}, {X: 2, Y: 0.002}}
	if err := ValidatePoints(pts, 1e-9); !errors.Is(err, ErrCollinearSet) {
		t.Errorf("got %v, want ErrCollinearSet", err)
	}
}

// TestPolygonAreaSigned checks the shoelace formula on a square in both
// orientations.
func TestPolygonAreaSigned(t *testing.T) {
	ccw := []Point{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}
	if a := PolygonArea(ccw); math.Abs(a-4) > 1e-12 {
		t.Errorf("ccw area = %g, want 4", a)
	}
	cw := []Point{ccw[3], ccw[2], ccw[1], ccw[0]}
	if a := PolygonArea(cw); math.Abs(a+4) > 1e-12 {
		t.Errorf("cw area = %g, want -4", a)
	}
}

// TestAllCollinearDetection checks both positive and negative cases.
func TestAllCollinearDetection(t *testing.T) {
	onLine := []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 5, Y: 0}, {X: -3, Y: 0}}
	if !AllCollinear(onLine, 1e-9) {
		t.Error("collinear points not detected")
	}
	offLine := []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 5, Y: 0.01}, {X: -3, Y: 0}}
	if AllCollinear(offLine, 1e-9) {
		t.Error("non-collinear points flagged as collinear")
	}
}

// TestPointInConvexPolygon checks containment on a square.
func TestPointInConvexPolygon(t *testing.T) {
	square := []Point{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}
	if !PointInConvexPolygon(Point{X: 1, Y: 1}, square, 1e-12) {
		t.Error("center should be inside")
	}
	if PointInConvexPolygon(Point{X: 3, Y: 1}, square, 1e-12) {
		t.Error("outside point should be rejected")
	}
	if !PointInConvexPolygon(Point{X: 2, Y: 1}, square, 1e-12) {
		t.Error("boundary point should be accepted")
	}
}
