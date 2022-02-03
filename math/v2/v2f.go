package v2

import "skyfx/math"

// import "mat"
// import "v1"

// Constants ...
var ZERO math.V2 = Makef(0.0)
var ONE math.V2 = Makef(1.0)

// var RED math.V2 = make(1.0, 0.0)
// var GREEN math.V2 = make(0.0, 1.0)
// var MIN math.V2 = makef(math.MIN_f32)
// var MAX math.V2 = makef(math.MAX_f32)
// var PI math.V2 = makef(math.PI_f32)

// // isnan ...
// func isnan(a math.V2) (out bool) {
// 	out = float32.isnan(a.X) ||
// 		  float32.isnan(a.Y)
// }

// // to_str ...
// func to_str(a math.V2) (out string) {
// 	out = fmt.Sprintf("{ %f, %f }", a.X, a.Y)
// }

// // min ...
// func min(a math.V2, b math.V2) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN min #0")
// 	//utils.PanicIf(isnan(b), "NAN min #1")
// 	out = a
// 	if b.X < a.X {
// 		out.X = b.X
// 	}
// 	if b.Y < a.Y {
// 		out.Y = b.Y
// 	}
// 	//utils.PanicIf(isnan(out), "NAN min")
// }

// // minf ...
// func minf(a math.V2, b float32) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN min #0")
// 	//utils.PanicIf(float32.isnan(b), "NAN min #1")
// 	out = a
// 	if b < a.X {
// 		out.X = b
// 	}
// 	if b < a.Y {
// 		out.Y = b
// 	}
// 	//utils.PanicIf(isnan(out), "NAN min")
// }

// // max ...
// func max(a math.V2, b math.V2) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN max #0")
// 	//utils.PanicIf(isnan(b), "NAN max #1")
// 	out = a
// 	if b.X > a.X {
// 		out.X = b.X
// 	}
// 	if b.Y > a.Y {
// 		out.Y = b.Y
// 	}
// 	//utils.PanicIf(isnan(out), "NAN max")
// }

// // maxf ...
// func maxf(a math.V2, b float32) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN max #0")
// 	//utils.PanicIf(float32.isnan(b), "NAN max #1")
// 	out = a
// 	if b > a.X {
// 		out.X = b
// 	}
// 	if b > a.Y {
// 		out.Y = b
// 	}
// 	//utils.PanicIf(isnan(out), "NAN max")
// }

// // clamp ...
// func clamp(a math.V2, vmin math.V2, vmax math.V2) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(isnan(vmin), "NAN clamp #1")
// 	//utils.PanicIf(isnan(vmax), "NAN clamp #2")
// 	out.X = math.Max_f32(vmin.X, math.Min_f32(a.X, vmax.X))
// 	out.Y = math.Max_f32(vmin.Y, math.Min_f32(a.Y, vmax.Y))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // clampf ...
// func clampf(a math.V2, fmin float32, fmax float32) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(float32.isnan(fmin), "NAN clamp #1")
// 	//utils.PanicIf(float32.isnan(fmax), "NAN clamp #2")
// 	out.X = math.Max_f32(fmin, math.Min_f32(a.X, fmax))
// 	out.Y = math.Max_f32(fmin, math.Min_f32(a.Y, fmax))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // sat ...
// func sat(a math.V2) (out math.V2) {
// 	//utils.PanicIf(isnan(a), "NAN sat #0")
// 	out.X = math.Max_f32(0.0, math.Min_f32(1.0, a.X))
// 	out.Y = math.Max_f32(0.0, math.Min_f32(1.0, a.Y))
// 	//utils.PanicIf(isnan(out), "NAN sat")
// }

// make ...
func Make(x float32, y float32) (out math.V2) {
	out.X = x
	out.Y = y
	//utils.PanicIf(isnan(out), "NAN make")
	return
}

// makef ...
func Makef(a float32) (out math.V2) {
	out.X = a
	out.Y = a
	//utils.PanicIf(isnan(out), "NAN makef")
	return
}

// // equ ...
// func equ(a math.V2, b math.V2) (out bool) {
// 	//utils.PanicIf(isnan(a), "NAN equ #0")
// 	//utils.PanicIf(isnan(b), "NAN equ #1")
// 	out = a.X == b.X &&
// 		  a.Y == b.Y
// }

// // nequ ...
// func nequ(a math.V2, b math.V2) (out bool) {
// 	//utils.PanicIf(isnan(a), "NAN nequ")
// 	//utils.PanicIf(isnan(b), "NAN nequ")
// 	out = a.X != b.X ||
// 		  a.Y != b.Y
// }

// // add ...
// func add(a math.V2, b math.V2) (out math.V2) {
// 	out.X = a.X + b.X
// 	out.Y = a.Y + b.Y
// 	//utils.PanicIf(isnan(out), "NAN add")
// }

// sub ...
func Sub(a math.V2, b math.V2) (out math.V2) {
	out.X = a.X - b.X
	out.Y = a.Y - b.Y
	//utils.PanicIf(isnan(out), "NAN sub")
	return
}

// // neg ...
// func neg(a math.V2) (out math.V2) {
// 	out.X = 0.0 - a.X
// 	out.Y = 0.0 - a.Y
// 	//utils.PanicIf(isnan(out), "NAN neg")
// }

// // mul ...
// func mul(a math.V2, b math.V2) (out math.V2) {
// 	out.X = a.X * b.X
// 	out.Y = a.Y * b.Y
// 	//utils.PanicIf(isnan(out), "NAN mul")
// }

// // mulf ...
// func mulf(a math.V2, f float32) (out math.V2) {
// 	out.X = a.X * f
// 	out.Y = a.Y * f
// 	//utils.PanicIf(isnan(out), "NAN mulf")
// }

// // madd ...
// func madd(a math.V2, b math.V2, c math.V2) (out math.V2) {
// 	out.X = a.X * b.X + c.X
// 	out.Y = a.Y * b.Y + c.Y
// 	//utils.PanicIf(isnan(out), "NAN madd")
// }

// // maddf ...
// func maddf(a math.V2, f float32, c math.V2) (out math.V2) {
// 	out.X = a.X * f + c.X
// 	out.Y = a.Y * f + c.Y
// 	//utils.PanicIf(isnan(out), "NAN maddf")
// }

// // div ...
// func div(a math.V2, b math.V2) (out math.V2) {
// 	out.X = a.X / b.X
// 	out.Y = a.Y / b.Y
// 	//utils.PanicIf(isnan(out), "NAN div")
// }

// // divf ...
// func divf(a math.V2, f float32) (out math.V2) {
// 	out.X = a.X / f
// 	out.Y = a.Y / f
// 	//utils.PanicIf(isnan(out), "NAN divf")
// }

// // lerp ...
// func lerp(a math.V2, b math.V2, t math.V2) (out math.V2) {
// 	out.X = a.X * (1.0 - t.X) + b.X * t.X
// 	out.Y = a.Y * (1.0 - t.Y) + b.Y * t.Y
// 	//utils.PanicIf(isnan(out), "NAN lerpf")
// }

// // lerpsat ...
// func lerpsat(a math.V2, b math.V2, t math.V2) (out math.V2) {
// 	out.X = a.X * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.X)) +
// 			b.X * math.Max_f32(0.0, math.Min_f32(1.0, t.X))
// 	out.Y = a.Y * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.Y)) +
// 			b.Y * math.Max_f32(0.0, math.Min_f32(1.0, t.Y))
// 	//utils.PanicIf(isnan(out), "NAN lerpsat")
// }

// // lerpf ...
// func lerpf(a math.V2, b math.V2, t float32) (out math.V2) {
// 	var nt float32 = 1.0 - t
// 	out.X = a.X * nt + b.X * t
// 	out.Y = a.Y * nt + b.Y * t
// 	//utils.PanicIf(isnan(out), "NAN lerp")
// }

// // lerpsatf ...
// func lerpsatf(a math.V2, b math.V2, t float32) (out math.V2) {
// 	t = math.Max_f32(0.0, math.Min_f32(1.0, t))
// 	var nt float32 = math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t))
// 	out.X = a.X * nt + b.X * t
// 	out.Y = a.Y * nt + b.Y * t
// 	//utils.PanicIf(isnan(out), "NAN lerpsat")
// }

// // rand ...
// func rand() (out math.V2) {
// 	out.X = float32.rand()
// 	out.Y = float32.rand()
// 	//utils.PanicIf(isnan(out), "NAN rand")
// }

// // srand ...
// func srand() (out math.V2) {
// 	out.X = float32.rand() * 2.0 - 1.0
// 	out.Y = float32.rand() * 2.0 - 1.0
// 	//utils.PanicIf(isnan(out), "NAN srand")
// }

// // dot ...
// func dot(a math.V2, b math.V2) (out float32) {
// 	out = a.X*b.X + a.Y*b.Y
// 	//utils.PanicIf(float32.isnan(out), "NAN dot")
// }

// // sqlength ...
// func sqlength(a math.V2) (out float32) {
// 	out = a.X*a.X + a.Y*a.Y
// 	//utils.PanicIf(float32.isnan(out), "NAN sqlength")
// }

// // length ...
// func length(a math.V2) (out float32) {
// 	out = math.Sqrt_f32(a.X*a.X + a.Y*a.Y)
// 	//utils.PanicIf(float32.isnan(out), "NAN length")
// }

// // normalize ...
// func normalize(a math.V2) (out math.V2) {
// 	var l float32 = math.Sqrt_f32(a.X*a.X + a.Y*a.Y)
// 	out.X = a.X / l
// 	out.Y = a.Y / l
// 	//utils.PanicIf(isnan(out), "NAN normalize")
// }
