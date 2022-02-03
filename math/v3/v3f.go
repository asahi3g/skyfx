package v3

import (
	gomath "math"
	"skyfx/math"
)

// import "mat"
// import "v1"
// import "v2"
// import "m44"

// // Constants ...
var ZERO math.V3 = Makef(0.0)

var ONE math.V3 = Makef(1.0)

// var RED v3 = make(1.0, 0.0, 0.0)
var GREEN math.V3 = Make(0.0, 1.0, 0.0)

// var BLUE v3 = make(0.0, 0.0, 1.0)
// var PINK v3 = make(1.0, 0.0, 1.0)
// var YELLOW v3 = make(1.0, 1.0, 0.0)
// var SKY v3 = make(0.0, 1.0, 1.0)

var MIN math.V3 = Makef(-gomath.MaxFloat32)
var MAX math.V3 = Makef(gomath.MaxFloat32)

// var PI v3 = makef(math.PI_f32)

// // isnan ...
// func isnan(a v3) (out bool) {
// 	out = float32.isnan(a.X) ||
// 		  float32.isnan(a.Y) ||
// 		  float32.isnan(a.Z)
// }

// // to_str ...
// func to_str(a v3) (out string) {
// 	out = fmt.Sprintf("{ %f, %f, %f }", a.X, a.Y, a.Z)
// }

// min ...
func Min(a math.V3, b math.V3) (out math.V3) {
	//utils.PanicIf(isnan(a), "NAN min #0")
	//utils.PanicIf(isnan(b), "NAN min #1")
	out = a
	if b.X < a.X {
		out.X = b.X
	}
	if b.Y < a.Y {
		out.Y = b.Y
	}
	if b.Z < a.Z {
		out.Z = b.Z
	}
	//utils.PanicIf(isnan(out), "NAN min")
	return
}

// // minf ...
// func minf(a V3, b float32) (out v3) {
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
// 	//utils.PanicIf(isnan(out), "NAN min")
// }

// max ...
func Max(a math.V3, b math.V3) (out math.V3) {
	//utils.PanicIf(isnan(a), "NAN max #0")
	//utils.PanicIf(isnan(b), "NAN max #1")
	out = a
	if b.X > a.X {
		out.X = b.X
	}
	if b.Y > a.Y {
		out.Y = b.Y
	}
	if b.Z > a.Z {
		out.Z = b.Z
	}
	//utils.PanicIf(isnan(out), "NAN max")
	return
}

// // maxf ...
// func maxf(a V3, b float32) (out v3) {
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
// 	//utils.PanicIf(isnan(out), "NAN max")
// }

// // clamp ...
// func clamp(a V3, vmin V3, vmax v3) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(isnan(vmin), "NAN clamp #1")
// 	//utils.PanicIf(isnan(vmax), "NAN clamp #2")
// 	out.X = math.Max_f32(vmin.X, math.Min_f32(a.X, vmax.X))
// 	out.Y = math.Max_f32(vmin.Y, math.Min_f32(a.Y, vmax.Y))
// 	out.Z = math.Max_f32(vmin.Z, math.Min_f32(a.Z, vmax.Z))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // clampf ...
// func clampf(a V3, fmin float32, fmax float32) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN clamp #0")
// 	//utils.PanicIf(float32.isnan(fmin), "NAN clamp #1")
// 	//utils.PanicIf(float32.isnan(fmax), "NAN clamp #2")
// 	out.X = math.Max_f32(fmin, math.Min_f32(a.X, fmax))
// 	out.Y = math.Max_f32(fmin, math.Min_f32(a.Y, fmax))
// 	out.Z = math.Max_f32(fmin, math.Min_f32(a.Z, fmax))
// 	//utils.PanicIf(isnan(out), "NAN clamp")
// }

// // sat ...
// func sat(a v3) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN sat #0")
// 	out.X = math.Max_f32(0.0, math.Min_f32(1.0, a.X))
// 	out.Y = math.Max_f32(0.0, math.Min_f32(1.0, a.Y))
// 	out.Z = math.Max_f32(0.0, math.Min_f32(1.0, a.Z))
// 	//utils.PanicIf(isnan(out), "NAN sat")
// }

// make ...
func Make(x float32, y float32, z float32) (out math.V3) {
	out.X = x
	out.Y = y
	out.Z = z
	//utils.PanicIf(isnan(out), "NAN make")
	return
}

// makef ...
func Makef(a float32) (out math.V3) {
	out.X = a
	out.Y = a
	out.Z = a
	//utils.PanicIf(isnan(out), "NAN makef")
	return
}

// // make_v12 ...
// func make_v12(x float32, a v2) (out v3) {
// 	out.X = x
// 	out.Y = a.X
// 	out.Z = a.Y
// 	//utils.PanicIf(isnan(out), "NAN make_v12")
// }

// // make_v21 ...
// func make_v21(a v2, z float32) (out v3) {
// 	out.X = a.X
// 	out.Y = a.Y
// 	out.Z = z
// 	//utils.PanicIf(isnan(out), "NAN make_v21")
// }

// // x ...
// func make_x(a v3) (out v3) {
// 	out.X = a.X
// }

// // y ...
// func make_y(a v3) (out v3) {
// 	out.Y = a.Y
// }

// // z ...
// func make_z(a v3) (out v3) {
// 	out.Z = a.Z
// }

// // xy ...
// func xy(a v3) (out v2) {
// 	//utils.PanicIf(isnan(a), "NAN xy #0")
// 	out.X = a.X
// 	out.Y = a.Y
// 	//utils.PanicIf(v2.isnan(out), "NAN xy")
// }

// // zw ...
// func yz(a v3) (out v2) {
// 	//utils.PanicIf(isnan(a), "NAN zw #0")
// 	out.X = a.Y
// 	out.Y = a.Z
// 	//utils.PanicIf(v2.isnan(out), "NAN yz")
// }

// equ ...
func Equ(a math.V3, b math.V3) (out bool) {
	//utils.PanicIf(isnan(a), "NAN equ a#0")
	//utils.PanicIf(isnan(b), "NAN equ b#0")
	out = a.X == b.X &&
		a.Y == b.Y &&
		a.Z == b.Z
	return
}

// nequ ...
func Nequ(a math.V3, b math.V3) (out bool) {
	//utils.PanicIf(isnan(a), "NAN nequ a#0")
	//utils.PanicIf(isnan(b), "NAN nequ b#0")
	out = a.X != b.X ||
		a.Y != b.Y ||
		a.Z != b.Z
	return
}

// add ...
func Add(a math.V3, b math.V3) (out math.V3) {
	out.X = a.X + b.X
	out.Y = a.Y + b.Y
	out.Z = a.Z + b.Z
	//utils.PanicIf(isnan(out), "NAN add")
	return
}

// sub ...
func Sub(a math.V3, b math.V3) (out math.V3) {
	out.X = a.X - b.X
	out.Y = a.Y - b.Y
	out.Z = a.Z - b.Z
	//utils.PanicIf(isnan(out), "NAN sub")
	return
}

// // neg ...
// func neg(a v3) (out v3) {
// 	out.X = 0.0 - a.X
// 	out.Y = 0.0 - a.Y
// 	out.Z = 0.0 - a.Z
// 	//utils.PanicIf(isnan(out), "NAN neg")
// }

// mul ...
func Mul(a math.V3, b math.V3) (out math.V3) {
	out.X = a.X * b.X
	out.Y = a.Y * b.Y
	out.Z = a.Z * b.Z
	//utils.PanicIf(isnan(out), "NAN mul")
	return
}

// mulf ...
func Mulf(a math.V3, f float32) (out math.V3) {
	out.X = a.X * f
	out.Y = a.Y * f
	out.Z = a.Z * f
	//utils.PanicIf(isnan(out), "NAN mulf")
	return
}

// // madd ...
// func madd(a V3, b V3, c v3) (out v3) {
// 	out.X = a.X * b.X + c.X
// 	out.Y = a.Y * b.Y + c.Y
// 	out.Z = a.Z * b.Z + c.Z
// 	//utils.PanicIf(isnan(out), "NAN madd")
// }

// // maddf ...
// func maddf(a V3, f float32, c v3) (out v3) {
// 	out.X = a.X * f + c.X
// 	out.Y = a.Y * f + c.Y
// 	out.Z = a.Z * f + c.Z
// 	//utils.PanicIf(isnan(out), "NAN maddf")
// }

// div ...
func Div(a math.V3, b math.V3) (out math.V3) {
	out.X = a.X / b.X
	out.Y = a.Y / b.Y
	out.Z = a.Z / b.Z
	//utils.PanicIf(isnan(out), "NAN div")
	return
}

// divf ...
func Divf(a math.V3, f float32) (out math.V3) {
	out.X = a.X / f
	out.Y = a.Y / f
	out.Z = a.Z / f
	//utils.PanicIf(isnan(out), "NAN divf")
	return
}

// // lerp ...
// func lerp(a V3, b V3, t v3) (out v3) {
// 	out.X = a.X * (1.0 - t.X) + b.X * t.X
// 	out.Y = a.Y * (1.0 - t.Y) + b.Y * t.Y
// 	out.Z = a.Z * (1.0 - t.Z) + b.Z * t.Z
// 	//utils.PanicIf(isnan(out), "NAN lerp")
// }

// // lerpsat ...
// func lerpsat(a V3, b V3, t v3) (out v3) {
// 	out.X = a.X * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.X)) +
// 			b.X * math.Max_f32(0.0, math.Min_f32(1.0, t.X))
// 	out.Y = a.Y * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.Y)) +
// 			b.Y * math.Max_f32(0.0, math.Min_f32(1.0, t.Y))
// 	out.Z = a.Z * math.Max_f32(0.0, math.Min_f32(1.0, 1.0 - t.Z)) +
// 			b.Z * math.Max_f32(0.0, math.Min_f32(1.0, t.Z))
// 	//utils.PanicIf(isnan(out), "NAN lerpsat")
// }

// lerpf ...
func Lerpf(a math.V3, b math.V3, t float32) (out math.V3) {
	var nt float32 = 1.0 - t
	out.X = a.X*nt + b.X*t
	out.Y = a.Y*nt + b.Y*t
	out.Z = a.Z*nt + b.Z*t
	//utils.PanicIf(isnan(out), "NAN lerp")
	return
}

// lerpsatf ...
func Lerpsatf(a math.V3, b math.V3, t float32) (out math.V3) {
	t = math.Max_f32(0.0, math.Min_f32(1.0, t))
	var nt float32 = math.Max_f32(0.0, math.Min_f32(1.0, 1.0-t))
	out.X = a.X*nt + b.X*t
	out.Y = a.Y*nt + b.Y*t
	out.Z = a.Z*nt + b.Z*t
	//utils.PanicIf(isnan(out), "NAN lerpsat")
	return
}

// // rand ...
// func rand() (out v3) {
// 	out.X = float32.rand()
// 	out.Y = float32.rand()
// 	out.Z = float32.rand()
// 	//utils.PanicIf(isnan(out), "NAN rand")
// }

// // srand ...
// func srand() (out v3) {
// 	out.X = float32.rand() * 2.0 - 1.0
// 	out.Y = float32.rand() * 2.0 - 1.0
// 	out.Z = float32.rand() * 2.0 - 1.0
// 	//utils.PanicIf(isnan(out), "NAN srand")
// }

// cross ...
func Cross(a math.V3, b math.V3) (out math.V3) {
	out.X = a.Y*b.Z - b.Y*a.Z
	out.Y = b.X*a.Z - a.X*b.Z
	out.Z = a.X*b.Y - b.X*a.Y
	//utils.PanicIf(isnan(out), "NAN cross")
	return
}

// dot ...
func Dot(a math.V3, b math.V3) (out float32) {
	out = a.X*b.X + a.Y*b.Y + a.Z*b.Z
	//utils.PanicIf(float32.isnan(out), "NAN dot")
	return
}

// sqlength ...
func Sqlength(a math.V3) (out float32) {
	out = a.X*a.X + a.Y*a.Y + a.Z*a.Z
	//utils.PanicIf(float32.isnan(out), "NAN sqlength")
	return
}

// length ...
func Length(a math.V3) (out float32) {
	out = math.Sqrt_f32(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	//utils.PanicIf(float32.isnan(out), "NAN length")
	return
}

// normalize ...
func Normalize(a math.V3) (out math.V3) {
	var l float32 = math.Sqrt_f32(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	out.X = a.X / l
	out.Y = a.Y / l
	out.Z = a.Z / l
	//utils.PanicIf(isnan(out), "NAN normalize")
	return
}

// // transform_vector_x ...
// func transform_vector_x(x float32, m m44) (out v3) {
// 	//utils.PanicIf(float32.isnan(x), "NAN transform_vector_x #0")
// 	//utils.PanicIf(m44.isnan(m), "NAN transform_vector_x #1")
// 	out.X = x * m.V00
// 	out.Y = x * m.V01
// 	out.Z = x * m.V02
// 	//utils.PanicIf(isnan(out), "NAN transform_vector_x")
// }

// // transform_vector_y ...
// func transform_vector_y(y float32, m m44) (out v3) {
// 	//utils.PanicIf(float32.isnan(y), "NAN transform_vector_y #0")
// 	//utils.PanicIf(m44.isnan(m), "NAN transform_vector_y #1")
// 	out.X = y * m.V10
// 	out.Y = y * m.V11
// 	out.Z = y * m.V12
// 	//utils.PanicIf(isnan(out), "NAN transform_vector_y")
// }

// // transform_vector_x ...
// func transform_vector_z(z float32, m m44) (out v3) {
// 	//utils.PanicIf(float32.isnan(z), "NAN transform_vector_z #0")
// 	//utils.PanicIf(m44.isnan(m), "NAN transform_vector_z #1")
// 	out.X = z * m.V20
// 	out.Y = z * m.V21
// 	out.Z = z * m.V22
// 	//utils.PanicIf(isnan(out), "NAN transform_vector_z")
// }

// // transform_vectorf ...
// func transform_vectorf(x float32, y float32, z float32, m m44) (out v3) {
// 	//utils.PanicIf(float32.isnan(x), "NAN transform_vectorf #0")
// 	//utils.PanicIf(float32.isnan(y), "NAN transform_vectorf #1")
// 	//utils.PanicIf(float32.isnan(z), "NAN transform_vectorf #2")
// 	//utils.PanicIf(m44.isnan(m), "NAN transform_vector #1")
// 	out.X = x * m.V00 + y * m.V10 + z * m.V20
// 	out.Y = x * m.V01 + y * m.V11 + z * m.V21
// 	out.Z = x * m.V02 + y * m.V12 + z * m.V22
// 	//utils.PanicIf(isnan(out), "NAN transform_vector")
// }

// // transform_vector ...
// func transform_vector(a V3, m m44) (out v3) {
// 	//utils.PanicIf(isnan(a), "NAN transform_vector #0")
// 	//utils.PanicIf(m44.isnan(m), "NAN transform_vector #1")
// 	out.X = a.X * m.V00 + a.Y * m.V10 + a.Z * m.V20
// 	out.Y = a.X * m.V01 + a.Y * m.V11 + a.Z * m.V21
// 	out.Z = a.X * m.V02 + a.Y * m.V12 + a.Z * m.V22
// 	//utils.PanicIf(isnan(out), "NAN transform_vector")
// }

// transform_point ...
func Transform_point(a math.V3, m math.M44) (out math.V3) {
	//utils.PanicIf(isnan(a), "NAN transform_vector #0")
	//utils.PanicIf(m44.isnan(m), "NAN transform_vector #1")
	out.X = a.X*m.V00 + a.Y*m.V10 + a.Z*m.V20 + m.V30
	out.Y = a.X*m.V01 + a.Y*m.V11 + a.Z*m.V21 + m.V31
	out.Z = a.X*m.V02 + a.Y*m.V12 + a.Z*m.V22 + m.V32
	//utils.PanicIf(isnan(out), "NAN transform_vector")
	return
}
