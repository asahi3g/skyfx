package gfx

// Globals ...
// var State State

// State ...
type State struct {
	View       []float32
	Projection []float32
}

// SamplerState ...
type SamplerState struct {
	min int32
	mag int32
	s   int32
	t   int32
	r   int32
}

func SamplerStateCreate(min int32, mag int32, s int32, t int32, r int32) (out SamplerState) {
	out.min = min
	out.mag = mag
	out.s = s
	out.t = t
	out.r = r
	return
}
