package m44

import (
	gomath "math"
	"skyfx/math"
)

// import "mat"

// // TODO : add NAN checks

// // Constants ...
var IDENTITY math.M44 = ident()
var INVALID math.M44

// equal ...
func Equ(a math.M44, b math.M44) (out bool) {
	out = a.V00 == b.V00 &&
		a.V01 == b.V01 &&
		a.V02 == b.V02 &&
		a.V03 == b.V03 &&
		a.V10 == b.V10 &&
		a.V11 == b.V11 &&
		a.V12 == b.V12 &&
		a.V13 == b.V13 &&
		a.V20 == b.V20 &&
		a.V21 == b.V21 &&
		a.V22 == b.V22 &&
		a.V23 == b.V23 &&
		a.V30 == b.V30 &&
		a.V31 == b.V31 &&
		a.V32 == b.V32 &&
		a.V33 == b.V33
	return
}

func Nequ(a math.M44, b math.M44) (out bool) {
	out = a.V00 != b.V00 ||
		a.V01 != b.V01 ||
		a.V02 != b.V02 ||
		a.V03 != b.V03 ||
		a.V10 != b.V10 ||
		a.V11 != b.V11 ||
		a.V12 != b.V12 ||
		a.V13 != b.V13 ||
		a.V20 != b.V20 ||
		a.V21 != b.V21 ||
		a.V22 != b.V22 ||
		a.V23 != b.V23 ||
		a.V30 != b.V30 ||
		a.V31 != b.V31 ||
		a.V32 != b.V32 ||
		a.V33 != b.V33
	return
}

// isident ...
func Isident(m math.M44) (out bool) {
	out = m.V00 == 1.0 &&
		m.V01 == 0.0 &&
		m.V02 == 0.0 &&
		m.V03 == 0.0 &&
		m.V10 == 0.0 &&
		m.V11 == 1.0 &&
		m.V12 == 0.0 &&
		m.V13 == 0.0 &&
		m.V20 == 0.0 &&
		m.V21 == 0.0 &&
		m.V22 == 1.0 &&
		m.V23 == 0.0 &&
		m.V30 == 0.0 &&
		m.V31 == 0.0 &&
		m.V32 == 0.0 &&
		m.V33 == 1.0
	return
}

// // isnan ...
// func isnan(m m44) (out bool) {
// 	out = float32.isnan(m.V00) ||
// 		  float32.isnan(m.V01) ||
// 		  float32.isnan(m.V02) ||
// 		  float32.isnan(m.V03) ||
// 		  float32.isnan(m.V10) ||
// 		  float32.isnan(m.V11) ||
// 		  float32.isnan(m.V12) ||
// 		  float32.isnan(m.V13) ||
// 		  float32.isnan(m.V20) ||
// 		  float32.isnan(m.V21) ||
// 		  float32.isnan(m.V22) ||
// 		  float32.isnan(m.V23) ||
// 		  float32.isnan(m.V30) ||
// 		  float32.isnan(m.V31) ||
// 		  float32.isnan(m.V32) ||
// 		  float32.isnan(m.V33)
// }

// // to_str ...
// func to_str(m m44) (out string) {
// 	out = fmt.Sprintf("{ %f, %f, %f, %f | %f, %f, %f, %f | %f, %f, %f, %f | %f, %f, %f, %f }",
// 			m.V00, m.V01, m.V02, m.V03,
// 			m.V10, m.V11, m.V12, m.V13,
// 			m.V20, m.V21, m.V22, m.V23,
// 			m.V30, m.V31, m.V32, m.V33)
// }

// ident ...
func ident() (out math.M44) {
	out.V00 = 1.0
	out.V11 = 1.0
	out.V22 = 1.0
	out.V33 = 1.0
	return
}

// // transpose ...
// func transpose(a m44) (out m44) { // ISSUE : if named MatrixIdentity
// 	out.V00 = a.V00
// 	out.V01 = a.V10
// 	out.V02 = a.V20
// 	out.V03 = a.V30

// 	out.V10 = a.V01
// 	out.V11 = a.V11
// 	out.V12 = a.V21
// 	out.V13 = a.V31

// 	out.V20 = a.V02
// 	out.V21 = a.V12
// 	out.V22 = a.V22
// 	out.V23 = a.V32

// 	out.V30 = a.V03
// 	out.V31 = a.V13
// 	out.V32 = a.V23
// 	out.V33 = a.V33
// }

// makef_AT ...
func Makef_AT(ax float32, ay float32, az float32, aw float32, tx float32, ty float32, tz float32) (out math.M44) {
	var a2 float32 = 0.5 * aw
	sin, cos := gomath.Sincos(float64(a2))
	var sina float32 = float32(sin)
	var qx float32 = ax * sina
	var qy float32 = ay * sina
	var qz float32 = az * sina
	var qw float32 = float32(cos)

	var qxx float32 = 2.0 * qx * qx
	var qxy float32 = 2.0 * qx * qy
	var qxz float32 = 2.0 * qx * qz
	var qxw float32 = 2.0 * qx * qw
	var qyy float32 = 2.0 * qy * qy
	var qyz float32 = 2.0 * qy * qz
	var qyw float32 = 2.0 * qy * qw
	var qzz float32 = 2.0 * qz * qz
	var qzw float32 = 2.0 * qz * qw

	out.V00 = 1.0 - qyy - qzz
	out.V01 = qxy + qzw
	out.V02 = qxz - qyw

	out.V10 = qxy - qzw
	out.V11 = 1.0 - qxx - qzz
	out.V12 = qyz + qxw

	out.V20 = qxz + qyw
	out.V21 = qyz - qxw
	out.V22 = 1.0 - qxx - qyy

	out.V30 = tx
	out.V31 = ty
	out.V32 = tz
	out.V33 = 1.0
	return
}

// // makev_AT ...
// func makev_AT(aa math.V4, t math.V3) (out m44) {
// 	var a2 float32 = 0.5 * aa.W
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = aa.X * sina
// 	var qy float32 = aa.Y * sina
// 	var qz float32 = aa.Z * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V30 = t.X
// 	out.V31 = t.Y
// 	out.V32 = t.Z
// 	out.V33 = 1.0
// }

// // make_SAT ...
// func make_SAT(trs TRS) (out m44) {
// 	var a2 float32 = 0.5 * trs.r.W
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = trs.r.X * sina
// 	var qy float32 = trs.r.Y * sina
// 	var qz float32 = trs.r.Z * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = (1.0 - qyy - qzz) * trs.S.X
// 	out.V01 = (qxy + qzw) * trs.S.X
// 	out.V02 = (qxz - qyw) * trs.S.X

// 	out.V10 = (qxy - qzw) * trs.S.Y
// 	out.V11 = (1.0 - qxx - qzz) * trs.S.Y
// 	out.V12 = (qyz + qxw) * trs.S.Y

// 	out.V20 = (qxz + qyw) * trs.S.Z
// 	out.V21 = (qyz - qxw) * trs.S.Z
// 	out.V22 = (1.0 - qxx - qyy) * trs.S.Z

// 	out.V30 = trs.t.X
// 	out.V31 = trs.t.Y
// 	out.V32 = trs.t.Z
// 	out.V33 = 1.0
// }

// makef_SAT ...
func Makef_SAT(sx float32, sy float32, sz float32, ax float32, ay float32, az float32, aw float32, tx float32, ty float32, tz float32) (out math.M44) {
	var a2 float32 = 0.5 * aw
	sina, qw := math.Sincos_f32(a2)
	var qx float32 = ax * sina
	var qy float32 = ay * sina
	var qz float32 = az * sina

	var qxx float32 = 2.0 * qx * qx
	var qxy float32 = 2.0 * qx * qy
	var qxz float32 = 2.0 * qx * qz
	var qxw float32 = 2.0 * qx * qw
	var qyy float32 = 2.0 * qy * qy
	var qyz float32 = 2.0 * qy * qz
	var qyw float32 = 2.0 * qy * qw
	var qzz float32 = 2.0 * qz * qz
	var qzw float32 = 2.0 * qz * qw

	out.V00 = (1.0 - qyy - qzz) * sx
	out.V01 = (qxy + qzw) * sx
	out.V02 = (qxz - qyw) * sx

	out.V10 = (qxy - qzw) * sy
	out.V11 = (1.0 - qxx - qzz) * sy
	out.V12 = (qyz + qxw) * sy

	out.V20 = (qxz + qyw) * sz
	out.V21 = (qyz - qxw) * sz
	out.V22 = (1.0 - qxx - qyy) * sz

	out.V30 = tx
	out.V31 = ty
	out.V32 = tz
	out.V33 = 1.0
	return
}

// // makev_SAT ...
// func makev_SAT(s math.V3, aa math.V4, t math.V3) (out m44) {
// 	var a2 float32 = 0.5 * aa.W
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = aa.X * sina
// 	var qy float32 = aa.Y * sina
// 	var qz float32 = aa.Z * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = (1.0 - qyy - qzz) * s.X
// 	out.V01 = (qxy + qzw) * s.X
// 	out.V02 = (qxz - qyw) * s.X

// 	out.V10 = (qxy - qzw) * s.Y
// 	out.V11 = (1.0 - qxx - qzz) * s.Y
// 	out.V12 = (qyz + qxw) * s.Y

// 	out.V20 = (qxz + qyw) * s.Z
// 	out.V21 = (qyz - qxw) * s.Z
// 	out.V22 = (1.0 - qxx - qyy) * s.Z

// 	out.V30 = t.X
// 	out.V31 = t.Y
// 	out.V32 = t.Z
// 	out.V33 = 1.0
// }

// makev_QT ...
func Makev_QT(q math.V4, t math.V3) (out math.M44) {
	var qxx float32 = 2.0 * q.X * q.X
	var qxy float32 = 2.0 * q.X * q.Y
	var qxz float32 = 2.0 * q.X * q.Z
	var qxw float32 = 2.0 * q.X * q.W
	var qyy float32 = 2.0 * q.Y * q.Y
	var qyz float32 = 2.0 * q.Y * q.Z
	var qyw float32 = 2.0 * q.Y * q.W
	var qzz float32 = 2.0 * q.Z * q.Z
	var qzw float32 = 2.0 * q.Z * q.W

	out.V00 = 1.0 - qyy - qzz
	out.V01 = qxy + qzw
	out.V02 = qxz - qyw

	out.V10 = qxy - qzw
	out.V11 = 1.0 - qxx - qzz
	out.V12 = qyz + qxw

	out.V20 = qxz + qyw
	out.V21 = qyz - qxw
	out.V22 = 1.0 - qxx - qyy

	out.V30 = t.X
	out.V31 = t.Y
	out.V32 = t.Z
	out.V33 = 1.0
	return
}

// make_SQT ...
func Make_SQT(trs math.TRS) (out math.M44) {
	var qxx float32 = 2.0 * trs.R.X * trs.R.X
	var qxy float32 = 2.0 * trs.R.X * trs.R.Y
	var qxz float32 = 2.0 * trs.R.X * trs.R.Z
	var qxw float32 = 2.0 * trs.R.X * trs.R.W
	var qyy float32 = 2.0 * trs.R.Y * trs.R.Y
	var qyz float32 = 2.0 * trs.R.Y * trs.R.Z
	var qyw float32 = 2.0 * trs.R.Y * trs.R.W
	var qzz float32 = 2.0 * trs.R.Z * trs.R.Z
	var qzw float32 = 2.0 * trs.R.Z * trs.R.W

	out.V00 = (1.0 - qyy - qzz) * trs.S.X
	out.V01 = (qxy + qzw) * trs.S.X
	out.V02 = (qxz - qyw) * trs.S.X

	out.V10 = (qxy - qzw) * trs.S.Y
	out.V11 = (1.0 - qxx - qzz) * trs.S.Y
	out.V12 = (qyz + qxw) * trs.S.Y

	out.V20 = (qxz + qyw) * trs.S.Z
	out.V21 = (qyz - qxw) * trs.S.Z
	out.V22 = (1.0 - qxx - qyy) * trs.S.Z

	out.V30 = trs.T.X
	out.V31 = trs.T.Y
	out.V32 = trs.T.Z
	out.V33 = 1.0
	return
}

// makef_SQT ...
func Makef_SQT(sx float32, sy float32, sz float32, qx float32, qy float32, qz float32, qw float32, tx float32, ty float32, tz float32) (out math.M44) {
	var qxx float32 = 2.0 * qx * qx
	var qxy float32 = 2.0 * qx * qy
	var qxz float32 = 2.0 * qx * qz
	var qxw float32 = 2.0 * qx * qw
	var qyy float32 = 2.0 * qy * qy
	var qyz float32 = 2.0 * qy * qz
	var qyw float32 = 2.0 * qy * qw
	var qzz float32 = 2.0 * qz * qz
	var qzw float32 = 2.0 * qz * qw

	out.V00 = (1.0 - qyy - qzz) * sx
	out.V01 = (qxy + qzw) * sx
	out.V02 = (qxz - qyw) * sx

	out.V10 = (qxy - qzw) * sy
	out.V11 = (1.0 - qxx - qzz) * sy
	out.V12 = (qyz + qxw) * sy

	out.V20 = (qxz + qyw) * sz
	out.V21 = (qyz - qxw) * sz
	out.V22 = (1.0 - qxx - qyy) * sz

	out.V30 = tx
	out.V31 = ty
	out.V32 = tz
	out.V33 = 1.0
	return
}

// makev_SQT ...
func Makev_SQT(s math.V3, q math.V4, t math.V3) (out math.M44) {
	var qxx float32 = 2.0 * q.X * q.X
	var qxy float32 = 2.0 * q.X * q.Y
	var qxz float32 = 2.0 * q.X * q.Z
	var qxw float32 = 2.0 * q.X * q.W
	var qyy float32 = 2.0 * q.Y * q.Y
	var qyz float32 = 2.0 * q.Y * q.Z
	var qyw float32 = 2.0 * q.Y * q.W
	var qzz float32 = 2.0 * q.Z * q.Z
	var qzw float32 = 2.0 * q.Z * q.W

	out.V00 = (1.0 - qyy - qzz) * s.X
	out.V01 = (qxy + qzw) * s.X
	out.V02 = (qxz - qyw) * s.X

	out.V10 = (qxy - qzw) * s.Y
	out.V11 = (1.0 - qxx - qzz) * s.Y
	out.V12 = (qyz + qxw) * s.Y

	out.V20 = (qxz + qyw) * s.Z
	out.V21 = (qyz - qxw) * s.Z
	out.V22 = (1.0 - qxx - qyy) * s.Z

	out.V30 = t.X
	out.V31 = t.Y
	out.V32 = t.Z
	out.V33 = 1.0
	return
}

// make_project ...
func Make_project(near float32, far float32, fov float32, width float32, height float32) (out math.M44) {
	sin, cos := gomath.Sincos(float64(fov))
	var m00 float32 = float32(cos) / float32(sin)

	prange := far - near

	out.V00 = m00
	out.V11 = m00 / (height / width)
	out.V22 = (0.0 - far - near) / prange
	out.V23 = -1.0
	out.V32 = (0.0 - 2.0*near*far) / prange
	return
}

// // make_translate ...
// func make_translate(x float32, y float32, z float32) (out m44) {

// 	out.V00 = 1.0
// 	out.V11 = 1.0
// 	out.V22 = 1.0
// 	out.V30 = x
// 	out.V31 = y
// 	out.V32 = z
// 	out.V33 = 1.0
// }

// // from_axis_angle ...
// func from_axis_angle(ax float32, ay float32, az float32, a float32) (out m44) {
// 	var a2 float32 = 0.5 * a
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = ax * sina
// 	var qy float32 = ay * sina
// 	var qz float32 = az * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V33 = 1.0
// }

// // from_axis_angle_v31 ...
// func from_axis_angle_v31(axis math.V3, angle float32) (out m44) {
// 	var a2 float32 = 0.5 * angle
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = axis.X * sina
// 	var qy float32 = axis.Y * sina
// 	var qz float32 = axis.Z * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V33 = 1.0
// }

// // from_axis_angle_v4 ...
// func from_axis_angle_v4(aa math.V4) (out m44) {
// 	var a2 float32 = 0.5 * aa.W
// 	var sina float32 = float32.sin(a2)
// 	var qx float32 = aa.X * sina
// 	var qy float32 = aa.Y * sina
// 	var qz float32 = aa.Z * sina
// 	var qw float32 = float32.cos(a2)

// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V33 = 1.0
// }

// // from_quat ...
// func from_quat(qx float32, qy float32, qz float32, qw float32) (out m44) {
// 	var qxx float32 = 2.0 * qx * qx
// 	var qxy float32 = 2.0 * qx * qy
// 	var qxz float32 = 2.0 * qx * qz
// 	var qxw float32 = 2.0 * qx * qw
// 	var qyy float32 = 2.0 * qy * qy
// 	var qyz float32 = 2.0 * qy * qz
// 	var qyw float32 = 2.0 * qy * qw
// 	var qzz float32 = 2.0 * qz * qz
// 	var qzw float32 = 2.0 * qz * qw

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V33 = 1.0
// }

// // from_quat_v4 ...
// func from_quat_v4(q math.V4) (out m44) {
// 	var qxx float32 = 2.0 * q.X * q.X
// 	var qxy float32 = 2.0 * q.X * q.Y
// 	var qxz float32 = 2.0 * q.X * q.Z
// 	var qxw float32 = 2.0 * q.X * q.W
// 	var qyy float32 = 2.0 * q.Y * q.Y
// 	var qyz float32 = 2.0 * q.Y * q.Z
// 	var qyw float32 = 2.0 * q.Y * q.W
// 	var qzz float32 = 2.0 * q.Z * q.Z
// 	var qzw float32 = 2.0 * q.Z * q.W

// 	out.V00 = 1.0 - qyy - qzz
// 	out.V01 = qxy + qzw
// 	out.V02 = qxz - qyw

// 	out.V10 = qxy - qzw
// 	out.V11 = 1.0 - qxx - qzz
// 	out.V12 = qyz + qxw

// 	out.V20 = qxz + qyw
// 	out.V21 = qyz - qxw
// 	out.V22 = 1.0 - qxx - qyy

// 	out.V33 = 1.0
// }

// // make_rotate_x ...
// func make_rotate_x(alpha float32) (out m44) {
// 	var cosA float32 = float32.cos(alpha)
// 	var sinA float32 = float32.sin(alpha)

// 	out.V00 = 1.0
// 	out.V11 = cosA
// 	out.V12 = sinA
// 	out.V21 = -sinA
// 	out.V22 = cosA
// 	out.V33 = 1.0
// }

// // make_rotate_y ...
// func make_rotate_y(alpha float32) (out m44) {
// 	var cosA float32 = float32.cos(alpha)
// 	var sinA float32 = float32.sin(alpha)

// 	out.V00 = cosA
// 	out.V02 = sinA
// 	out.V11 = 1.0
// 	out.V20 = -sinA
// 	out.V22 = cosA
// 	out.V33 = 1.0
// }

// // make_scale ...
// func make_scale(x float32, y float32, z float32) (out m44) {
// 	out.V00 = x
// 	out.V11 = y
// 	out.V22 = z
// 	out.V33 = 1.0
// }

// determinant ...
func Determinant(m math.M44) (out float32) {
	var a float32 = m.V22*m.V33 - m.V23*m.V32
	var b float32 = m.V21*m.V33 - m.V23*m.V31
	var c float32 = m.V21*m.V32 - m.V22*m.V31
	var d float32 = m.V20*m.V33 - m.V23*m.V30
	var e float32 = m.V20*m.V32 - m.V22*m.V30
	var f float32 = m.V20*m.V31 - m.V21*m.V30

	out = m.V00*(m.V11*a-m.V12*b+m.V13*c) -
		m.V01*(m.V10*a-m.V12*d+m.V13*e) +
		m.V02*(m.V10*b-m.V11*d+m.V13*f) -
		m.V03*(m.V10*c-m.V11*e+m.V12*f)
	return
}

// inverse ...
func Inverse(a math.M44) (out math.M44) {
	var d1 float32 = a.V00*a.V11 - a.V01*a.V10
	var d2 float32 = a.V00*a.V12 - a.V02*a.V10
	var d3 float32 = a.V00*a.V13 - a.V03*a.V10
	var d4 float32 = a.V01*a.V12 - a.V02*a.V11
	var d5 float32 = a.V01*a.V13 - a.V03*a.V11
	var d6 float32 = a.V02*a.V13 - a.V03*a.V12
	var d7 float32 = a.V20*a.V31 - a.V21*a.V30
	var d8 float32 = a.V20*a.V32 - a.V22*a.V30
	var d9 float32 = a.V20*a.V33 - a.V23*a.V30
	var d10 float32 = a.V21*a.V32 - a.V22*a.V31
	var d11 float32 = a.V21*a.V33 - a.V23*a.V31
	var d12 float32 = a.V22*a.V33 - a.V23*a.V32
	var md float32 = d1*d12 - d2*d11 + d3*d10 + d4*d9 - d5*d8 + d6*d7

	out.V00 = (a.V11*d12 - a.V12*d11 + a.V13*d10) * md
	out.V01 = (0.0 - a.V01*d12 + a.V02*d11 - a.V03*d10) * md
	out.V02 = (a.V31*d6 - a.V32*d5 + a.V33*d4) * md
	out.V03 = (0.0 - a.V21*d6 + a.V22*d5 - a.V23*d4) * md
	out.V10 = (0.0 - a.V10*d12 + a.V12*d9 - a.V13*d8) * md
	out.V11 = (a.V00*d12 - a.V02*d9 + a.V03*d8) * md
	out.V12 = (0.0 - a.V30*d6 + a.V32*d3 - a.V33*d2) * md
	out.V13 = (a.V20*d6 - a.V22*d3 + a.V23*d2) * md
	out.V20 = (a.V10*d11 - a.V11*d9 + a.V13*d7) * md
	out.V21 = (0.0 - a.V00*d11 + a.V01*d9 - a.V03*d7) * md
	out.V22 = (a.V30*d5 - a.V31*d3 + a.V33*d1) * md
	out.V23 = (0.0 - a.V20*d5 + a.V21*d3 - a.V23*d1) * md
	out.V30 = (0.0 - a.V10*d10 + a.V11*d8 - a.V12*d7) * md
	out.V31 = (a.V00*d10 - a.V01*d8 + a.V02*d7) * md
	out.V32 = (0.0 - a.V30*d4 + a.V31*d2 - a.V32*d1) * md
	out.V33 = (a.V20*d4 - a.V21*d2 + a.V22*d1) * md
	return
}

// mulISSUE ...
func MulISSUE(a math.M44, b math.M44) (out math.M44) {
	out.V00 = a.V00*b.V00 + a.V01*b.V10 + a.V02*b.V20 + a.V03*b.V30
	out.V01 = a.V00*b.V01 + a.V01*b.V11 + a.V02*b.V21 + a.V03*b.V31
	out.V02 = a.V00*b.V02 + a.V01*b.V12 + a.V02*b.V22 + a.V03*b.V32
	out.V03 = a.V00*b.V03 + a.V01*b.V13 + a.V02*b.V23 + a.V03*b.V33
	out.V10 = a.V10*b.V00 + a.V11*b.V10 + a.V12*b.V20 + a.V13*b.V30
	out.V11 = a.V10*b.V01 + a.V11*b.V11 + a.V12*b.V21 + a.V13*b.V31
	out.V12 = a.V10*b.V02 + a.V11*b.V12 + a.V12*b.V22 + a.V13*b.V32
	out.V13 = a.V10*b.V03 + a.V11*b.V13 + a.V12*b.V23 + a.V13*b.V33
	out.V20 = a.V20*b.V00 + a.V21*b.V10 + a.V22*b.V20 + a.V23*b.V30
	out.V21 = a.V20*b.V01 + a.V21*b.V11 + a.V22*b.V21 + a.V23*b.V31
	out.V22 = a.V20*b.V02 + a.V21*b.V12 + a.V22*b.V22 + a.V23*b.V32
	out.V23 = a.V20*b.V03 + a.V21*b.V13 + a.V22*b.V23 + a.V23*b.V33
	out.V30 = a.V30*b.V00 + a.V31*b.V10 + a.V32*b.V20 + a.V33*b.V30
	out.V31 = a.V30*b.V01 + a.V31*b.V11 + a.V32*b.V21 + a.V33*b.V31
	out.V32 = a.V30*b.V02 + a.V31*b.V12 + a.V32*b.V22 + a.V33*b.V32
	out.V33 = a.V30*b.V03 + a.V31*b.V13 + a.V32*b.V23 + a.V33*b.V33
	return
}

// // last ...
// func last(stack []m44) (out m44) {
// 	out = stack[len(stack) - 1]
// }

// push ...
func Push(stack []math.M44, matrix math.M44) (out []math.M44) {
	out = stack
	var stackLen int32 = int32(len(out))
	if stackLen > 0 {
		out = append(out, MulISSUE(matrix, out[stackLen-1]))
	} else {
		out = append(out, matrix)
	}
	return
}

// pop ...
func Pop(stack []math.M44, count int32) (out []math.M44) {
	out = stack
	out = out[:len(out)-int(count)]
	return
}
