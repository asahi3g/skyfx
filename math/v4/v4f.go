package v4

import (
	"skyfx/math"
)

// import "mat"
// import "v1"
// import "v2"
// import "v3"

var ZERO math.V4 = Makef(0.0)
var ONE math.V4 = Makef(1.0)

var RED math.V4 = Make(1.0, 0.0, 0.0, 1.0)
var GREEN math.V4 = Make(0.0, 1.0, 0.0, 1.0)
var BLUE math.V4 = Make(0.0, 0.0, 1.0, 1.0)

// var YELLOW v4 = make(1.0, 1.0, 0.0, 1.0)
// var SKY    v4 = make(0.0, 1.0, 1.0, 1.0)
// var BLACK  v4 = make(0.0, 0.0, 0.0, 1.0)
var ALPHA math.V4 = Make(0.0, 0.0, 0.0, 1.0)

// var MIN v4 = makef(math.MIN_f32)
// var MAX v4 = makef(math.MAX_f32)
// var PI v4 = makef(math.PI_f32)

// // isnan ...
// func isnan(a v4) (out bool) {
// 	out = float32.isnan(a.X) ||
// 		  float32.isnan(a.Y) ||
// 		  float32.isnan(a.Z) ||
// 		  float32.isnan(a.W)
// }

// // to_str ...
// func to_str(a v4) (out string) {
// 	out = fmt.Sprintf("{ %f, %f, %f, %f }", a.X, a.Y, a.Z, a.W)
// }

// // min ...
// func min(a v4, b v4) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN min #0")
// 	//utils.PanicIf(isnan(b), "NAN min #1")
// 	out = a
// 	if b.X < a.X {
// 		out.X = b.X
// 	}
// 	if b.Y < a.Y {
// 		out.Y = b.Y
// 	}
// 	if b.Z < a.Z {
// 		out.Z = b.Z
// 	}
// 	if b.W < a.W {
// 		out.W = b.W
// 	}
// 	//utils.PanicIf(isnan(out), "NAN min")
// }

// // minf ...
// func minf(a v4, b float32) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN min #0")
// 	//utils.PanicIf(float32.isnan(b), "NAN min #1")
// 	out = a
// 	if b < a.X {
// 		out.X = b
// 	}
// 	if b < a.Y {
// 		out.Y = b
// 	}
// 	if b < a.Z {
// 		out.Z = b
// 	}
// 	if b < a.W {
// 		out.W = b
// 	}
// 	//utils.PanicIf(isnan(out), "NAN min")
// }

// // max ...
// func max(a v4, b v4) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN max #0")
// 	//utils.PanicIf(isnan(b), "NAN max #1")
// 	out = a
// 	if b.X > a.X {
// 		out.X = b.X
// 	}
// 	if b.Y > a.Y {
// 		out.Y = b.Y
// 	}
// 	if b.Z > a.Z {
// 		out.Z = b.Z
// 	}
// 	if b.W > a.W {
// 		out.W = b.W
// 	}
// 	//utils.PanicIf(isnan(out), "NAN max")
// }

// // maxf ...
// func maxf(a v4, b float32) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN max #0")
// 	//utils.PanicIf(float32.isnan(b), "NAN max #1")
// 	out = a
// 	if b > a.X {
// 		out.X = b
// 	}
// 	if b > a.Y {
// 		out.Y = b
// 	}
// 	if b > a.Z {
// 		out.Z = b
// 	}
// 	if b > a.W {
// 		out.W = b
// 	}
// 	//utils.PanicIf(isnan(out), "NAN max")
// }

// // clamp ...
// func clamp(a v4, vmin v4, vmax v4) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(isnan(vmin), "NAN clamp #1")
// 	//utils.PanicIf(isnan(vmax), "NAN clamp #2")
// 	out.X = math.Max_f32(vmin.X, math.Min_f32(a.X, vmax.X))
// 	out.Y = math.Max_f32(vmin.Y, math.Min_f32(a.Y, vmax.Y))
// 	out.Z = math.Max_f32(vmin.Z, math.Min_f32(a.Z, vmax.Z))
// 	out.W = math.Max_f32(vmin.W, math.Min_f32(a.W, vmax.W))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // clampf ...
// func clampf(a v4, fmin float32, fmax float32) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(float32.isnan(fmin), "NAN clamp #1")
// 	//utils.PanicIf(float32.isnan(fmax), "NAN clamp #2")
// 	out.X = math.Max_f32(fmin, math.Min_f32(a.X, fmax))
// 	out.Y = math.Max_f32(fmin, math.Min_f32(a.Y, fmax))
// 	out.Z = math.Max_f32(fmin, math.Min_f32(a.Z, fmax))
// 	out.W = math.Max_f32(fmin, math.Min_f32(a.W, fmax))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // sat ...
// func sat(a v4) (out v4) {
// 	//utils.PanicIf(isnan(a), "NAN sat #0")
// 	out.X = math.Max_f32(0.0, math.Min_f32(1.0, a.X))
// 	out.Y = math.Max_f32(0.0, math.Min_f32(1.0, a.Y))
// 	out.Z = math.Max_f32(0.0, math.Min_f32(1.0, a.Z))
// 	out.W = math.Max_f32(0.0, math.Min_f32(1.0, a.W))
// 	//utils.PanicIf(isnan(out), "NAN sat")
// }

// make ...
func Make(x float32, y float32, z float32, w float32) (out math.V4) {
	out.X = x
	out.Y = y
	out.Z = z
	out.W = w
	//utils.PanicIf(isnan(out), "NAN make")
	return
}

// makef ...
func Makef(a float32) (out math.V4) {
	out.X = a
	out.Y = a
	out.Z = a
	out.W = a
	//utils.PanicIf(isnan(out), "NAN makef")
	return
}

// make v31 ...
func Make_v31(a math.V3, w float32) (out math.V4) {
	out.X = a.X
	out.Y = a.Y
	out.Z = a.Z
	out.W = w
	//utils.PanicIf(isnan(out), "NAN make_v31")
	return
}

// // make_v22 ...
// func make_v22(a v2, b v2) (out v4) {
// 	out.X = a.X
// 	out.Y = a.Y
// 	out.Z = b.X
// 	out.W = b.Y
// 	//utils.PanicIf(isnan(out), "NAN make_v22")
// }

// // make_v13 ...
// func make_v13(x float32, a v3) (out v4) {
// 	out.X = x
// 	out.Y = a.X
// 	out.Z = a.Y
// 	out.W = a.Z
// 	//utils.PanicIf(isnan(out), "NAN make_v13")
// }

// // xyz ...
// func xyz(a v4) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN xyz #0")
// 	out.X = a.X
// 	out.Y = a.Y
// 	out.Z = a.Z
// 	//utils.PanicIf(v3.isnan(out), "NAN xyz")
// }

// // yzw ...
// func yzw(a v4) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN yzw #0")
// 	out.X = a.Y
// 	out.Y = a.Z
// 	out.Z = a.W
// 	//utils.PanicIf(v3.isnan(out), "NAN yzw")
// }

// // xy ...
// func xy(a v4) (out v2) {
// 	//utils.PanicIf(isnan(a), "NAN xy #0")
// 	out.X = a.X
// 	out.Y = a.Y
// 	//utils.PanicIf(v2.isnan(out), "NAN xy")
// }

// // zw ...
// func zw(a v4) (out v2) {
// 	//utils.PanicIf(isnan(a), "NAN zw #0")
// 	out.X = a.Z
// 	out.Y = a.W
// 	//utils.PanicIf(v2.isnan(out), "NAN zw")
// }

// equ ...
func Equ(a math.V4, b math.V4) (out bool) {
	//utils.PanicIf(isnan(a), "NAN equ a#0")
	//utils.PanicIf(isnan(b), "NAN equ b#1")
	out = a.X == b.X &&
		a.Y == b.Y &&
		a.Z == b.Z &&
		a.W == b.W
	return
}

// nequ ...
func Nequ(a math.V4, b math.V4) (out bool) {
	//utils.PanicIf(isnan(a), "NAN nequ a#0")
	//utils.PanicIf(isnan(b), "NAN nequ b#1")
	out = a.X != b.X ||
		a.Y != b.Y ||
		a.Z != b.Z ||
		a.W != b.W
	return
}

// // add ...
// func add(a v4, b v4) (out v4) {
// 	out.X = a.X + b.X
// 	out.Y = a.Y + b.Y
// 	out.Z = a.Z + b.Z
// 	out.W = a.W + b.W
// 	//utils.PanicIf(isnan(out), "NAN add")
// }

// // sub ...
// func sub(a v4, b v4) (out v4) {
// 	out.X = a.X - b.X
// 	out.Y = a.Y - b.Y
// 	out.Z = a.Z - b.Z
// 	out.W = a.W - b.W
// 	//utils.PanicIf(isnan(out), "NAN sub")
// }

// // neg ...
// func neg(a v4) (out v4) { // ISSUE : can't use unary -
// 	out.X = 0.0-a.X
// 	out.Y = 0.0-a.Y
// 	out.Z = 0.0-a.Z
// 	out.W = 0.0-a.W
// 	//utils.PanicIf(isnan(out), "NAN neg")
// }

// // mul ...
// func mul(a v4, b v4) (out v4) {
// 	out.X = a.X * b.X
// 	out.Y = a.Y * b.Y
// 	out.Z = a.Z * b.Z
// 	out.W = a.W * b.W
// 	//utils.PanicIf(isnan(out), "NAN mul")
// }

// // mulf ...
// func mulf(a v4, f float32) (out v4) {
// 	out.X = a.X * f
// 	out.Y = a.Y * f
// 	out.Z = a.Z * f
// 	out.W = a.W * f
// 	//utils.PanicIf(isnan(out), "NAN mulf")
// }

// // madd ...
// func madd(a v4, b v4, c v4) (out v4) {
// 	out.X = a.X * b.X + c.X
// 	out.Y = a.Y * b.Y + c.Y
// 	out.Z = a.Z * b.Z * c.Z
// 	out.W = a.W * b.W * c.W
// 	//utils.PanicIf(isnan(out), "NAN madd")
// }

// // maddf ...
// func maddf(a v4, f float32, c v4) (out v4) {
// 	out.X = a.X * f + c.X
// 	out.Y = a.Y * f + c.Y
// 	out.Z = a.Z * f + c.Z
// 	out.W = a.W * f + c.W
// 	//utils.PanicIf(isnan(out), "NAN maddf")
// }

// // div ...
// func div(a v4, b v4) (out v4) {
// 	out.X = a.X / b.X
// 	out.Y = a.Y / b.Y
// 	out.Z = a.Z / b.Z
// 	out.W = a.W / b.W
// 	//utils.PanicIf(isnan(out), "NAN div")
// }

// // divf ...
// func divf(a v4, f float32) (out v4) {
// 	out.X = a.X / f
// 	out.Y = a.Y / f
// 	out.Z = a.Z / f
// 	out.W = a.W / f
// 	//utils.PanicIf(isnan(out), "NAN divf")
// }

// // lerp ...
// func lerp(a v4, b v4, t v4) (out v4) {
// 	out.X = a.X * (1.0 - t.X) + b.X * t.X
// 	out.Y = a.Y * (1.0 - t.Y) + b.Y * t.Y
// 	out.Z = a.Z * (1.0 - t.Z) + b.Z * t.Z
// 	out.W = a.W * (1.0 - t.W) + b.W * t.W
// 	//utils.PanicIf(isnan(out), "NAN lerp")
// }

// // lerpsat ...
// func lerpsat(a v4, b v4, t v4) (out v4) {
// 	out.X = a.X * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.X)) +
// 			b.X * math.Max_f32(0.0, math.Min_f32(1.0, t.X))
// 	out.Y = a.Y * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.Y)) +
// 			b.Y * math.Max_f32(0.0, math.Min_f32(1.0, t.Y))
// 	out.Z = a.Z * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.Z)) +
// 			b.Z * math.Max_f32(0.0, math.Min_f32(1.0, t.Z))
// 	out.W = a.W * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.W)) +
// 			b.W * math.Max_f32(0.0, math.Min_f32(1.0, t.W))
// 	//utils.PanicIf(isnan(out), "NAN lerpsat")
// }

// // lerpf ...
// func lerpf(a v4, b v4, t float32) (out v4) {
// 	var nt float32 = 1.0 - t
// 	out.X = a.X * nt + b.X * t
// 	out.Y = a.Y * nt + b.Y * t
// 	out.Z = a.Z * nt + b.Z * t
// 	out.W = a.W * nt + b.W * t
// 	//utils.PanicIf(isnan(out), "NAN lerp")
// }

// // lerpsatf ...
// func lerpsatf(a v4, b v4, t float32) (out v4) {
// 	t = math.Max_f32(0.0, math.Min_f32(1.0, t))
// 	var nt float32 = math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t))
// 	out.X = a.X*nt + b.X*t
// 	out.Y = a.Y*nt + b.Y*t
// 	out.Z = a.Z*nt + b.Z*t
// 	out.W = a.W*nt + b.W*t
// 	//utils.PanicIf(isnan(out), "NAN lerpsat")
// }

// // rand ...
// func rand() (out v4) {
// 	out.X = float32.rand()
// 	out.Y = float32.rand()
// 	out.Z = float32.rand()
// 	out.W = float32.rand()
// 	//utils.PanicIf(isnan(out), "NAN rand")
// }

// // srand ...
// func srand() (out v4) {
// 	out.X = float32.rand() * 2.0 - 1.0
// 	out.Y = float32.rand() * 2.0 - 1.0
// 	out.Z = float32.rand() * 2.0 - 1.0
// 	out.W = float32.rand() * 2.0 - 1.0
// 	//utils.PanicIf(isnan(out), "NAN srand")
// }

// // dot ...
// func dot(a v4, b v4) (out float32) {
// 	out = a.X * b.X + a.Y * b.Y + a.Z * b.Z + a.W * b.W
// 	//utils.PanicIf(float32.isnan(out), "NAN dot")
// }

// // sqlength ...
// func sqlength(a v4) (out float32) {
// 	out = a.X * a.X + a.Y * a.Y + a.Z * a.Z + a.W * a.W
// 	//utils.PanicIf(float32.isnan(out), "NAN sqlength")
// }

// // length ...
// func length(a v4) (out float32) {
// 	out = math.Sqrt_f32(a.X * a.X + a.Y * a.Y + a.Z * a.Z + a.W * a.W)
// 	//utils.PanicIf(float32.isnan(out), "NAN length")
// }

// normalize ...
func Normalize(a math.V4) (out math.V4) {
	var l float32 = math.Sqrt_f32(a.X*a.X + a.Y*a.Y + a.Z*a.Z + a.W*a.W)
	out.X = a.X / l
	out.Y = a.Y / l
	out.Z = a.Z / l
	out.W = a.W / l
	//utils.PanicIf(isnan(out), "NAN normalize")
	return
}

// transform ...
func Transform(a math.V4, m math.M44) (out math.V4) {
	out.X = a.X*m.V00 + a.Y*m.V10 + a.Z*m.V20 + a.W*m.V30
	out.Y = a.X*m.V01 + a.Y*m.V11 + a.Z*m.V21 + a.W*m.V31
	out.Z = a.X*m.V02 + a.Y*m.V12 + a.Z*m.V22 + a.W*m.V32
	out.W = a.X*m.V03 + a.Y*m.V13 + a.Z*m.V23 + a.W*m.V33
	//utils.PanicIf(isnan(out), "NAN transform")
	return
}
