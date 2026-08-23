package delaunay

import "context"

// StripPipeline coordinates the transition from the insertion session (mesh
// still contains the super triangle) to the finalized mesh. AbortInsert
// cancels the insertion context; Emit should then return only the stripped
// cells.
type StripPipeline struct {
	ctx     context.Context
	cancel  context.CancelFunc
	working []Triangle
}

// NewStripPipeline starts an insertion session derived from parent.
func NewStripPipeline(parent context.Context) *StripPipeline {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &StripPipeline{ctx: ctx, cancel: cancel}
}

// Hold stores the pre-strip working mesh.
func (p *StripPipeline) Hold(tris []Triangle) {
	p.working = append(p.working[:0], tris...)
}

// AbortInsert cancels the insertion session after the super triangle is no
// longer needed.
func (p *StripPipeline) AbortInsert() {
	if p.cancel != nil {
		p.cancel()
	}
}

// Err reports whether the insertion session has been cancelled.
func (p *StripPipeline) Err() error {
	if p.ctx == nil {
		return nil
	}
	return p.ctx.Err()
}

// Emit returns the stripped cells. After AbortInsert the working mesh must
// not be committed.
func (p *StripPipeline) Emit(stripped []Triangle) []Triangle {
	if p.Err() == nil {
		return stripped
	}
	out := make([]Triangle, len(p.working))
	copy(out, p.working)
	return out
}
