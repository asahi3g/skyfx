package math

import (
	"encoding/binary"
	gomath "math"
	"math/rand"
	"skyfx/utils"
)

// Constants ...
const (
	MAX_f32 float32 = gomath.MaxFloat32
	MIN_f32 float32 = -MAX_f32

	PI_f32   float32 = 3.14159265358979323846264338327950288
	M2PI_f32 float32 = 2.0 * PI_f32

	MAX_ui32 uint32 = ^uint32(0)
)

// Types...
type V2 struct {
	X float32
	Y float32
}

type V3 struct {
	X float32
	Y float32
	Z float32
}

type V4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

type M44 struct {
	V00 float32
	V01 float32
	V02 float32
	V03 float32

	V10 float32
	V11 float32
	V12 float32
	V13 float32

	V20 float32
	V21 float32
	V22 float32
	V23 float32

	V30 float32
	V31 float32
	V32 float32
	V33 float32
}

type TRS struct {
	T V3
	R V4
	S V3
}

// Funcs...

// IsNan
func IsNan_f32(f float32) bool {
	return f != f
}

func IsNan_f64(f float64) bool {
	return f != f
}

// Abs
func Abs_i32(i int32) int32 {
	if i < 0 {
		return -i
	}
	return i
}

// Min ...
func Min_i32(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func Min_f32(a float32, b float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN Min_f32 input #0")
	panicIf(IsNan_f32(b), "NaN Min_f32 input #1")
	out = a
	if b < a {
		out = b
	}
	panicIf(IsNan_f32(out), "NaN Min_f32 output #0")
	return
}

// max ...
func Max_f32(a float32, b float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN Max_f32 input #0")
	panicIf(IsNan_f32(b), "NaN Max_f32 input #1")
	out = a
	if b > a {
		out = b
	}
	panicIf(IsNan_f32(out), "NaN Max_f32 output #0")
	return
}

// Sincos
func Sincos_f32(a float32) (s, c float32) {
	panicIf(IsNan_f32(a), "NaN Sincos_f32 input #0")

	sin, cos := gomath.Sincos(float64(a))
	panicIf(IsNan_f64(sin), "NaN Sincos_f32 output #0")
	panicIf(IsNan_f64(cos), "NaN Sincos_f32 output #1")

	s = float32(sin)
	c = float32(cos)
	panicIf(IsNan_f32(s), "NaN Sincos_f32 output #0")
	panicIf(IsNan_f32(c), "NaN Sincos_f32 output #1")
	return
}

// Acos
func Acos_f32(value float32) (out float32) {
	panicIf(IsNan_f32(value), "NaN Acos_f32 input #0")

	acos := gomath.Acos(float64(value))
	panicIf(IsNan_f64(acos), "NaN Acos_f32 output #0")

	out = float32(acos)
	panicIf(IsNan_f32(out), "NaN Sincos_f32 output #0")
	return
}

// Sqrt
func Sqrt_f32(a float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN Sqrt_f32 input #0")

	out = float32(gomath.Sqrt(float64(a)))
	panicIf(IsNan_f32(out), "NaN Sqrt_f32 output #0")
	return
}

// Abs...
func Abs_f32(a float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN AbsF32 input #0")
	out = a
	if a < 0.0 {
		out = -a
	}
	panicIf(IsNan_f32(out), "NaN AbsF32 output #0")
	return
}

// clamp ...
func Clamp_f32(a float32, fmin float32, fmax float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN Clamp_f32 input #0")
	panicIf(IsNan_f32(fmin), "NaN Clamp_f32 input #1")
	panicIf(IsNan_f32(fmax), "NaN Clamp_f32 input #2")

	out = Max_f32(fmin, Min_f32(a, fmax))
	panicIf(IsNan_f32(out), "NaN Clamp_f32 output #0")
	return
}

// sat ...
func Sat_f32(a float32) (out float32) {
	panicIf(IsNan_f32(a), "NaN Sat_f32 input #0")

	out = Max_f32(0.0, Min_f32(1.0, a))
	panicIf(IsNan_f32(out), "NaN Sat_f32 output #0")
	return
}

// lerp ...
func Lerp_f32(a float32, b float32, t float32) (out float32) {
	panicIf(IsNan_f32(out), "NaN Lerp_f32 input #0")
	panicIf(IsNan_f32(out), "NaN Lerp_f32 input #1")
	panicIf(IsNan_f32(out), "NaN Lerp_f32 input #2")

	out = a*(1.0-t) + b*t
	panicIf(IsNan_f32(out), "NaN Lerp_f32 output #0")
	return
}

// lerpsat ...
func Lerpsat_f32(a float32, b float32, t float32) (out float32) {
	panicIf(IsNan_f32(out), "NaN Lerpsat_f32 input #0")
	panicIf(IsNan_f32(out), "NaN Lerpsat_f32 input #1")
	panicIf(IsNan_f32(out), "NaN Lerpsat_f32 input #2")

	t = Max_f32(0.0, Min_f32(1.0, t))
	var nt float32 = Max_f32(0.0, Min_f32(1.0, 1.0-t))
	out = a*nt + b*t
	utils.PanicIf(IsNan_f32(out), "NaN Lerpsat_f32 output #0")
	return
}

// srand ...
func Srand_f32() (out float32) {
	out = rand.Float32()*2.0 - 1.0
	utils.PanicIf(IsNan_f32(out), "NaN Srand_f32 output #0")
	return
}

type Vector_f32 struct {
	data  []float32
	count uint64
}

func (this *Vector_f32) Resize(count uint64, copy bool) []float32 {
	if this.data == nil && count > 0 {
		this.data = make([]float32, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]float32, count-this.count)...)
		} else {
			this.data = make([]float32, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_f32) Get() []float32 {
	return this.data[:this.count]
}

func (this *Vector_f32) Count() uint64 {
	return this.count
}

type Vector_i32 struct {
	data  []int32
	count uint64
}

func (this *Vector_i32) Resize(count uint64, copy bool) []int32 {
	if this.data == nil && count > 0 {
		this.data = make([]int32, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]int32, count-this.count)...)
		} else {
			this.data = make([]int32, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_i32) Get() []int32 {
	return this.data[:this.count]
}

func (this *Vector_i32) Count() uint64 {
	return this.count
}

type Vector_ui32 struct {
	data  []uint32
	count uint64
}

func (this *Vector_ui32) Resize(count uint64, copy bool) []uint32 {
	if this.data == nil && count > 0 {
		this.data = make([]uint32, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]uint32, count-this.count)...)
		} else {
			this.data = make([]uint32, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_ui32) Get() []uint32 {
	return this.data[:this.count]
}

func (this *Vector_ui32) Count() uint64 {
	return this.count
}

type Vector_ui16 struct {
	data  []uint16
	count uint64
}

func (this *Vector_ui16) Resize(count uint64, copy bool) []uint16 {
	if this.data == nil && count > 0 {
		this.data = make([]uint16, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]uint16, count-this.count)...)
		} else {
			this.data = make([]uint16, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_ui16) Get() []uint16 {
	return this.data[:this.count]
}

func (this *Vector_ui16) Count() uint64 {
	return this.count
}

type Vector_ui8 struct {
	data  []uint8
	count uint64
}

func (this *Vector_ui8) Resize(count uint64, copy bool) []uint8 {
	if this.data == nil && count > 0 {
		this.data = make([]uint8, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]uint8, count-this.count)...)
		} else {
			this.data = make([]uint8, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_ui8) Get() []uint8 {
	return this.data[:this.count]
}

func (this *Vector_ui8) Count() uint64 {
	return this.count
}

func ReadV3(handle int32, out *V3) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &out.X); err == nil {
			if err := binary.Read(file, binary.LittleEndian, &out.Y); err == nil {
				if err := binary.Read(file, binary.LittleEndian, &out.Z); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func ReadV3Slice(handle int32, dest []V3, count uint64) (success bool) {
	if count > 0 {
		if file := utils.ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i].X); err == nil {
					if err := binary.Read(file, binary.LittleEndian, &dest[i].Y); err == nil {
						if err := binary.Read(file, binary.LittleEndian, &dest[i].Z); err == nil {
							success = true
						}
					}
				}
			}
		}
	}

	return
}

func WriteV3(handle int32, out *V3) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, &out.X); err == nil {
			if err := binary.Write(file, binary.LittleEndian, &out.Y); err == nil {
				if err := binary.Write(file, binary.LittleEndian, &out.Z); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteV3Slice(handle int32, value []V3) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if value != nil {
			for i := 0; i < len(value); i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i].X); err == nil {
					if err := binary.Write(file, binary.LittleEndian, value[i].Y); err == nil {
						if err := binary.Write(file, binary.LittleEndian, value[i].Z); err == nil {
							success = true
						}
					}
				}
			}
		}
	}
	return
}

type Vector_v3 struct {
	data  []V3
	count uint64
}

func (this *Vector_v3) Resize(count uint64, copy bool) []V3 {
	if this.data == nil && count > 0 {
		this.data = make([]V3, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]V3, count-this.count)...)
		} else {
			this.data = make([]V3, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_v3) Get() []V3 {
	return this.data[:this.count]
}

func (this *Vector_v3) Count() uint64 {
	return this.count
}

func ReadV4(handle int32, out *V4) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &out.X); err == nil {
			if err := binary.Read(file, binary.LittleEndian, &out.Y); err == nil {
				if err := binary.Read(file, binary.LittleEndian, &out.Z); err == nil {
					if err := binary.Read(file, binary.LittleEndian, &out.W); err == nil {
						success = true
					}
				}
			}
		}
	}
	return
}

func ReadV4Slice(handle int32, dest []V4, count uint64) (success bool) {
	if count > 0 {
		if file := utils.ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i].X); err == nil {
					if err := binary.Read(file, binary.LittleEndian, &dest[i].Y); err == nil {
						if err := binary.Read(file, binary.LittleEndian, &dest[i].Z); err == nil {
							if err := binary.Read(file, binary.LittleEndian, &dest[i].W); err == nil {
								success = true
							}
						}
					}
				}
			}
		}
	}

	return
}

func WriteV4(handle int32, out *V4) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, out.X); err == nil {
			if err := binary.Write(file, binary.LittleEndian, out.Y); err == nil {
				if err := binary.Write(file, binary.LittleEndian, out.Z); err == nil {
					if err := binary.Write(file, binary.LittleEndian, out.W); err == nil {
						success = true
					}
				}
			}
		}
	}
	return
}

func WriteV4Slice(handle int32, value []V4) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if value != nil {
			for i := 0; i < len(value); i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i].X); err == nil {
					if err := binary.Write(file, binary.LittleEndian, value[i].Y); err == nil {
						if err := binary.Write(file, binary.LittleEndian, value[i].Z); err == nil {
							if err := binary.Write(file, binary.LittleEndian, value[i].W); err == nil {
								success = true
							}
						}
					}
				}
			}
		}
	}
	return
}

type Vector_v4 struct {
	data  []V4
	count uint64
}

func (this *Vector_v4) Resize(count uint64, copy bool) []V4 {
	if this.data == nil && count > 0 {
		this.data = make([]V4, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]V4, count-this.count)...)
		} else {
			this.data = make([]V4, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_v4) Get() []V4 {
	return this.data[:this.count]
}

func (this *Vector_v4) Count() uint64 {
	return this.count
}

func ReadM44(handle int32, out *M44) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &out.V00); err == nil {
			if err := binary.Read(file, binary.LittleEndian, &out.V01); err == nil {
				if err := binary.Read(file, binary.LittleEndian, &out.V02); err == nil {
					if err := binary.Read(file, binary.LittleEndian, &out.V03); err == nil {
						if err := binary.Read(file, binary.LittleEndian, &out.V10); err == nil {
							if err := binary.Read(file, binary.LittleEndian, &out.V11); err == nil {
								if err := binary.Read(file, binary.LittleEndian, &out.V12); err == nil {
									if err := binary.Read(file, binary.LittleEndian, &out.V13); err == nil {
										if err := binary.Read(file, binary.LittleEndian, &out.V20); err == nil {
											if err := binary.Read(file, binary.LittleEndian, &out.V21); err == nil {
												if err := binary.Read(file, binary.LittleEndian, &out.V22); err == nil {
													if err := binary.Read(file, binary.LittleEndian, &out.V23); err == nil {
														if err := binary.Read(file, binary.LittleEndian, &out.V30); err == nil {
															if err := binary.Read(file, binary.LittleEndian, &out.V31); err == nil {
																if err := binary.Read(file, binary.LittleEndian, &out.V32); err == nil {
																	if err := binary.Read(file, binary.LittleEndian, &out.V33); err == nil {
																		success = true
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func ReadM44Slice(handle int32, dest []M44, count uint64) (success bool) {
	if count > 0 {
		if file := utils.ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i].V00); err == nil {
					if err := binary.Read(file, binary.LittleEndian, &dest[i].V01); err == nil {
						if err := binary.Read(file, binary.LittleEndian, &dest[i].V02); err == nil {
							if err := binary.Read(file, binary.LittleEndian, &dest[i].V03); err == nil {
								if err := binary.Read(file, binary.LittleEndian, &dest[i].V10); err == nil {
									if err := binary.Read(file, binary.LittleEndian, &dest[i].V11); err == nil {
										if err := binary.Read(file, binary.LittleEndian, &dest[i].V12); err == nil {
											if err := binary.Read(file, binary.LittleEndian, &dest[i].V13); err == nil {
												if err := binary.Read(file, binary.LittleEndian, &dest[i].V20); err == nil {
													if err := binary.Read(file, binary.LittleEndian, &dest[i].V21); err == nil {
														if err := binary.Read(file, binary.LittleEndian, &dest[i].V22); err == nil {
															if err := binary.Read(file, binary.LittleEndian, &dest[i].V23); err == nil {
																if err := binary.Read(file, binary.LittleEndian, &dest[i].V30); err == nil {
																	if err := binary.Read(file, binary.LittleEndian, &dest[i].V31); err == nil {
																		if err := binary.Read(file, binary.LittleEndian, &dest[i].V32); err == nil {
																			if err := binary.Read(file, binary.LittleEndian, &dest[i].V33); err == nil {
																				success = true
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return
}

func WriteM44(handle int32, out *M44) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, out.V00); err == nil {
			if err := binary.Write(file, binary.LittleEndian, out.V01); err == nil {
				if err := binary.Write(file, binary.LittleEndian, out.V02); err == nil {
					if err := binary.Write(file, binary.LittleEndian, out.V03); err == nil {
						if err := binary.Write(file, binary.LittleEndian, out.V10); err == nil {
							if err := binary.Write(file, binary.LittleEndian, out.V11); err == nil {
								if err := binary.Write(file, binary.LittleEndian, out.V12); err == nil {
									if err := binary.Write(file, binary.LittleEndian, out.V13); err == nil {
										if err := binary.Write(file, binary.LittleEndian, out.V20); err == nil {
											if err := binary.Write(file, binary.LittleEndian, out.V21); err == nil {
												if err := binary.Write(file, binary.LittleEndian, out.V22); err == nil {
													if err := binary.Write(file, binary.LittleEndian, out.V23); err == nil {
														if err := binary.Write(file, binary.LittleEndian, out.V30); err == nil {
															if err := binary.Write(file, binary.LittleEndian, out.V31); err == nil {
																if err := binary.Write(file, binary.LittleEndian, out.V32); err == nil {
																	if err := binary.Write(file, binary.LittleEndian, out.V33); err == nil {
																		success = true
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func WriteM44Slice(handle int32, value []M44) (success bool) {
	if file := utils.ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i].V00); err == nil {
					if err := binary.Write(file, binary.LittleEndian, value[i].V01); err == nil {
						if err := binary.Write(file, binary.LittleEndian, value[i].V02); err == nil {
							if err := binary.Write(file, binary.LittleEndian, value[i].V03); err == nil {
								if err := binary.Write(file, binary.LittleEndian, value[i].V10); err == nil {
									if err := binary.Write(file, binary.LittleEndian, value[i].V11); err == nil {
										if err := binary.Write(file, binary.LittleEndian, value[i].V12); err == nil {
											if err := binary.Write(file, binary.LittleEndian, value[i].V13); err == nil {
												if err := binary.Write(file, binary.LittleEndian, value[i].V20); err == nil {
													if err := binary.Write(file, binary.LittleEndian, value[i].V21); err == nil {
														if err := binary.Write(file, binary.LittleEndian, value[i].V22); err == nil {
															if err := binary.Write(file, binary.LittleEndian, value[i].V23); err == nil {
																if err := binary.Write(file, binary.LittleEndian, value[i].V30); err == nil {
																	if err := binary.Write(file, binary.LittleEndian, value[i].V31); err == nil {
																		if err := binary.Write(file, binary.LittleEndian, value[i].V32); err == nil {
																			if err := binary.Write(file, binary.LittleEndian, value[i].V33); err == nil {
																				success = true
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

type Vector_m44 struct {
	data  []M44
	count uint64
}

func (this *Vector_m44) Resize(count uint64, copy bool) []M44 {
	if this.data == nil && count > 0 {
		this.data = make([]M44, count)
	} else if count <= this.count {
	} else if count > this.count {
		if copy {
			this.data = append(this.data, make([]M44, count-this.count)...)
		} else {
			this.data = make([]M44, count)
		}
	}
	this.count = count
	return this.data
}

func (this *Vector_m44) Get() []M44 {
	return this.data[:this.count]
}

func (this *Vector_m44) Count() uint64 {
	return this.count
}
