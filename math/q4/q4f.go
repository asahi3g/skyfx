package q4

import (
	gomath "math"
	"skyfx/math"
	// 	v1 "skyfx/math/v1"
	// 	v3 "skyfx/math/v3"
	// 	v4 "skyfx/math/v4"
)

// var IDENTITY math.V4 = v4.Make(0.0, 0.0, 0.0, 1.0)

// slerp ...
func Slerp(a math.V4, b math.V4, t float32) (out math.V4) {
	//utils.PanicIf(v4.isnan(a), "NAN slerp #0")
	//utils.PanicIf(v4.isnan(b), "NAN slerp #1")
	//utils.PanicIf(float32.isnan(t), "NAN slerp #2")
	var dot float32 = a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
	if dot < 0.0 {
		b.X = 0.0 - b.X
		b.Y = 0.0 - b.Y
		b.Z = 0.0 - b.Z
		b.W = 0.0 - b.W
		dot = -dot
	}

	if dot > 0.9995 {
		out.X = a.X + t*(b.X-a.X)
		out.Y = a.Y + t*(b.Y-a.Y)
		out.Z = a.Z + t*(b.Z-a.Z)
		out.W = a.W + t*(b.W-a.W)
	} else {
		var theta0 float32 = float32(gomath.Acos(float64(dot)))
		var theta float32 = theta0 * t

		sint, cost := gomath.Sincos(float64(theta))
		var sinTheta float32 = float32(sint)
		var sinTheta0 float32 = float32(gomath.Sin(float64(theta0)))
		var s0 float32 = float32(cost) - dot*sinTheta/sinTheta0
		var s1 float32 = sinTheta / sinTheta0

		out.X = a.X*s0 + b.X*s1
		out.Y = a.Y*s0 + b.Y*s1
		out.Z = a.Z*s0 + b.Z*s1
		out.W = a.W*s0 + b.W*s1
	}

	var l float32 = float32(gomath.Sqrt(float64(out.X*out.X + out.Y*out.Y + out.Z*out.Z + out.W*out.W)))
	out.X = out.X / l
	out.Y = out.Y / l
	out.Z = out.Z / l
	out.W = out.W / l
	//utils.PanicIf(v4.isnan(out), "NAN slerp")
	return
}

// // lerp ....
// func lerp(a math.V4, b math.V4, t float32) (out math.V4) {
// 	//utils.PanicIf(v4.isnan(a), "NAN lerp #0")
// 	//utils.PanicIf(v4.isnan(b), "NAN lerp #2")
// 	//utils.PanicIf(float32.isnan(t), "NAN lerp #3")
// 	if (a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W) < 0.0 {
// 		t = 0.0 - t
// 	}
// 	var mt float32 = 1.0 - t

// 	out.X = mt*a.X + t*b.X
// 	out.Y = mt*a.Y + t*b.Y
// 	out.Z = mt*a.Z + t*b.Z
// 	out.W = mt*a.W + t*b.W

// 	var l float32 = float32(gomath.Sqrt(float64(out.X*out.X + out.Y*out.Y + out.Z*out.Z + out.W*out.W)))
// 	out.X = out.X / l
// 	out.Y = out.Y / l
// 	out.Z = out.Z / l
// 	out.W = out.W / l
// 	//utils.PanicIf(v4.isnan(out), "NAN lerp")
// 	return
// }

// // from_vectors ...
// func from_vectors(a math.V3, b math.V3, f math.V3) (out math.V4) {
// 	//utils.PanicIf(v3.isnan(a), "NAN from_vectors #0")
// 	//utils.PanicIf(v3.isnan(b), "NAN from_vectors #1")
// 	//utils.PanicIf(v3.isnan(f), "NAN from_vectors #2")

// 	var d float32 = a.X*b.X + a.Y*b.Y + a.Z*b.Z
// 	if d >= 1.0 {
// 		out = IDENTITY
// 	}
// 	if d < (1e-6 - 1.0) {
// 		if v3.Nequ(f, v3.ZERO) {
// 			out = from_axis_angle(f.X, f.Y, f.Z, math.PI_f32)
// 		} else {
// 			var x math.V3 = v3.Cross(v3.RED, a)
// 			if v3.Sqlength(x) == 0.0 {
// 				x = v3.Cross(v3.GREEN, a)
// 			}
// 			x = v3.Normalize(a)
// 			out = from_axis_angle(x.X, x.Y, x.Z, math.PI_f32)
// 		}
// 	} else {
// 		var s float32 = math.Sqrt_f32(1.0+d) * 2.0
// 		var is float32 = 1.0 / s
// 		var c math.V3 = v3.Cross(a, b)

// 		out.X = c.X * is
// 		out.Y = c.Y * is
// 		out.Z = c.Z * is
// 		out.W = s * 0.5

// 		out = v4.Normalize(out)
// 	}

// 	/*out.X = a.X*b.Y + a.Y*b.Y + a.Z*b.Z
// 	out.Y = a.Y*b.Z - b.Y*a.Z
// 	out.Z = b.X*a.Z - a.X*b.Z
// 	out.W = a.X*b.Y - b.X*a.Y
// 	var l float32 = math.Sqrt_f32(out.X*out.X + out.Y*out.Y + out.Z*out.Z + out.W*out.W)
// 	out.W = out.W + l
// 	out.X = out.X / l
// 	out.Y = out.Y / l
// 	out.Z = out.Z / l
// 	out.W = out.W / l*/
// 	//utils.PanicIf(v4.isnan(out), "NAN from_vectors")
// }

// // from_axis_angle ...
// func from_axis_angle(x float32, y float32, z float32, a float32) (out v4) {
// 	//utils.PanicIf(float32.isnan(x), "NAN from_axis_angle #0")
// 	//utils.PanicIf(float32.isnan(y), "NAN from_axis_angle #1")
// 	//utils.PanicIf(float32.isnan(z), "NAN from_axis_angle #2")
// 	//utils.PanicIf(float32.isnan(a), "NAN from_axis_angle #3")
// 	var a2 float32 = 0.5 * a
// 	var sina float32 = float32.sin(a2)
// 	out.X = x * sina
// 	out.Y = y * sina
// 	out.Z = z * sina
// 	out.W = float32.cos(a2)
// 	//utils.PanicIf(v4.isnan(out), "NAN from_axis_angle")
// }

// // from_axis_angle_v31 ...
// func from_axis_angle_v31(axis math.V3, angle float32) (out v4) {
// 	//utils.PanicIf(v3.isnan(axis), "NAN from_axis_angle_v31 #0")
// 	//utils.PanicIf(float32.isnan(angle), "NAN from_axis_angle_v31 #1")
// 	var a2 float32 = 0.5 * angle
// 	var sina float32 = float32.sin(a2)
// 	out.X = axis.X * sina
// 	out.Y = axis.Y * sina
// 	out.Z = axis.Z * sina
// 	out.W = float32.cos(a2)
// 	//utils.PanicIf(v4.isnan(out), "NAN from_axis_angle_v31")
// }

// // from_axis_angle_v4 ...
// func from_axis_angle_v4(aa mat.v4) (out v4) {
// 	//utils.PanicIf(v4.isnan(aa), "NAN from_axis_angle_v4 #0")
// 	var a2 float32 = 0.5 * aa.W
// 	var sina float32 = float32.sin(a2)
// 	out.X = aa.X * sina
// 	out.Y = aa.Y * sina
// 	out.Z = aa.Z * sina
// 	out.W = float32.cos(a2)
// 	//utils.PanicIf(v4.isnan(out), "NAN from_axis_angle")
// }

// from_yaw_pitch_roll ...
func From_yaw_pitch_roll(y float32, p float32, r float32) (out math.V4) {
	//utils.PanicIf(float32.isnan(y), "NAN from_yaw_pitch_roll #0")
	//utils.PanicIf(float32.isnan(p), "NAN from_yaw_pitch_roll #1")
	//utils.PanicIf(float32.isnan(r), "NAN from_yaw_pitch_roll #2")

	var y2 float32 = y * 0.5
	siny, cosy := math.Sincos_f32(y2)

	var p2 float32 = p * 0.5
	sinp, cosp := math.Sincos_f32(p2)

	var r2 float32 = r * 0.5
	sinr, cosr := math.Sincos_f32(r2)

	out.X = cosy*sinp*cosr + siny*cosp*sinr
	out.Y = siny*cosp*cosr + cosy*sinp*sinr
	out.Z = cosy*cosp*sinr + siny*sinp*cosr
	out.W = cosy*cosp*cosr + siny*sinp*sinr
	//utils.PanicIf(v4.isnan(out), "NAN from_yaw_pitch_roll")
	return
}

// // to_axis ...
// func to_axis(x float32, y float32, z float32, w float32) (out v3) {
// 	//utils.PanicIf(float32.isnan(x), "NAN to_axis #0")
// 	//utils.PanicIf(float32.isnan(y), "NAN to_axis #1")
// 	//utils.PanicIf(float32.isnan(z), "NAN to_axis #2")
// 	//utils.PanicIf(float32.isnan(w), "NAN to_axis #3")
// 	var w2 float32 = 1.0 - w*w
// 	if w2 <= 0.0 {
// 		out.X = 0.0
// 		out.Y = 0.0
// 		out.Z = 1.0
// 	} else {
// 		var rtw2 float32 = 1.0 / math.Sqrt_f32(w2)
// 		out.X = x * rtw2
// 		out.Y = y * rtw2
// 		out.Z = z * rtw2
// 	}
// 	//utils.PanicIf(v3.isnan(out), "NAN to_axis")
// }

// // to_axis_angle ...
// func to_axis_angle(x float32, y float32, z float32, w float32) (out v4) {
// 	//utils.PanicIf(float32.isnan(x), "NAN to_axis_angle #0")
// 	//utils.PanicIf(float32.isnan(y), "NAN to_axis_angle #1")
// 	//utils.PanicIf(float32.isnan(z), "NAN to_axis_angle #2")
// 	//utils.PanicIf(float32.isnan(w), "NAN to_axis_angle #3")
// 	out.W = 2.0 * math.Acos_f32(w)
// 	var s float32 = math.Sqrt_f32(1.0 - w*w)
// 	if s < 0.001 {
// 		out.X = x
// 		out.Y = y
// 		out.Z = z
// 	} else {
// 		out.X = x / s
// 		out.Y = y / s
// 		out.Z = z / s
// 	}
// 	//utils.PanicIf(v4.isnan(out), "NAN to_axis_angle")
// }

// func to_axis_anglev(axisAngle mat.v4) (out v4) {
// 	//utils.PanicIf(v4.isnan(axisAngle), "NAN to_axis_anglev #0")
// 	out.W = 2.0 * math.Acos_f32(axisAngle.W)
// 	var s float32 = math.Sqrt_f32(1.0 - axisAngle.W*axisAngle.W)
// 	if s < 0.001 {
// 		out.X = axisAngle.X
// 		out.Y = axisAngle.Y
// 		out.Z = axisAngle.Z
// 	} else {
// 		out.X = axisAngle.X / s
// 		out.Y = axisAngle.Y / s
// 		out.Z = axisAngle.Z / s
// 	}
// 	//utils.PanicIf(v4.isnan(out), "NAN to_axis_anglev")
// }

// mul ...
func Mul(a math.V4, b math.V4) (out math.V4) {
	//utils.PanicIf(v4.isnan(a), "NAN mul #0")
	//utils.PanicIf(v4.isnan(b), "NAN mul #1")
	out.X = a.X*b.W + b.X*a.W + a.Y*b.Z - a.Z*b.Y
	out.Y = a.Y*b.W + b.Y*a.W + a.Z*b.X - a.X*b.Z
	out.Z = a.Z*b.W + b.Z*a.W + a.X*b.Y - a.Y*b.X
	out.W = a.W*b.W - a.X*b.X - a.Y*b.Y - a.Z*b.Z
	//utils.PanicIf(v4.isnan(out), "NAN mul")
	return
}

// // transform_vector ...
// func transform_vector(q v4, v v3) (out v3) {
// 	//utils.PanicIf(v4.isnan(q), "NAN transform_vector #0")
// 	//utils.PanicIf(v3.isnan(v), "NAN transform_vector #1")
// 	var a float32 = q.X * 2.0
// 	var b float32 = q.Y * 2.0
// 	var c float32 = q.Z * 2.0
// 	var d float32 = q.X * a
// 	var e float32 = q.Y * b
// 	var f float32 = q.Z * c
// 	var g float32 = q.X * b
// 	var h float32 = q.X * c
// 	var i float32 = q.Y * c
// 	var j float32 = q.W * a
// 	var k float32 = q.W * b
// 	var l float32 = q.W * c
// 	out.X = (1.0-e-f)*v.X + (g-l)*v.Y + (h+k)*v.Z
// 	out.Y = (g+l)*v.X + (1.0-d-f)*v.Y + (i-j)*v.Z
// 	out.Z = (h-k)*v.X + (i+j)*v.Y + (1.0-d-e)*v.Z
// 	//utils.PanicIf(v3.isnan(out), "NAN transform_vector")
// }
