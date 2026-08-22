package voronoi

import (
	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/geom"
)

// ValidateVertexCenters verifies the defining property of the Voronoi
// construction: every Voronoi vertex equals the circumcenter of its triangle.
// The check recomputes each circumcenter and compares it with the stored
// vertex within the relative tolerance.
func ValidateVertexCenters(pts []geom.Point, tris []delaunay.Triangle, verts []geom.Point, tol float64) error {
	if len(verts) != len(tris) {
		return &MismatchError{Want: len(tris), Got: len(verts), What: "vertex count must equal triangle count"}
	}
	scale := geom.ScaleForTolerance(pts)
	thr := tol * scale
	for i, t := range tris {
		c, err := delaunay.CircumcenterOf(pts, t)
		if err != nil {
			return err
		}
		if !c.Equal(verts[i], thr) {
			return &MismatchError{
				Want: len(tris), Got: len(verts),
				What: "vertex " + itoa(i) + " is not the circumcenter of its triangle",
			}
		}
	}
	return nil
}

// ValidateDualEdgeCount checks that the number of dual edges equals the
// number of interior mesh edges, (3T - h)/2. A missing or duplicated dual
// edge breaks the duality and is reported here.
func ValidateDualEdgeCount(numTriangles, hullSize, numEdges int) error {
	want := InteriorEdgeCount(numTriangles, hullSize)
	if numEdges != want {
		return &MismatchError{
			Want: want,
			Got:  numEdges,
			What: "dual edge count must equal interior mesh edge count (3T-h)/2",
		}
	}
	return nil
}

// MismatchError describes a structural mismatch in a Voronoi diagram.
type MismatchError struct {
	Want int
	Got  int
	What string
}

// Error implements the error interface.
func (e *MismatchError) Error() string {
	return "voronoi: " + e.What + " (want " + itoa(e.Want) + ", got " + itoa(e.Got) + ")"
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [24]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
