package intersect

import (
	"skyfx/fps"
	"skyfx/math"
	v3 "skyfx/math/v3"
)

// ProjectPointsOnAxis ...
func ProjectPointsOnAxis(points []math.V3, axis math.V3) (min float32, max float32) {
	min = math.MAX_f32
	max = math.MIN_f32
	var count int = len(points)
	for i := 0; i < count; i++ {
		var dot float32 = v3.Dot(axis, points[i])
		if dot < min {
			min = dot
		}
		if dot > max {
			max = dot
		}
	}
	return
}

// TriangleIntersectsAABB ...
func TriangleIntersectsAABB(sizeMin math.V3, sizeMax math.V3, boxVertices []math.V3, triangleVertices []math.V3) (out bool) {
	var triangleMin float32
	var triangleMax float32
	var boxMin float32
	var boxMax float32

	var bn0 math.V3 = v3.Make(1.0, 0.0, 0.0)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, bn0)
	if (triangleMax < sizeMin.X) || (triangleMin > sizeMax.X) {
		out = false
		return
	}

	var bn1 math.V3 = v3.Make(0.0, 1.0, 0.0)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, bn1)
	if (triangleMax < sizeMin.Y) || (triangleMin > sizeMax.Y) {
		out = false
		return
	}

	var bn2 math.V3 = v3.Make(0.0, 0.0, 1.0)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, bn2)
	if (triangleMax < sizeMin.Z) || (triangleMin > sizeMax.Z) {
		out = false
		return
	}

	var t01 math.V3 = v3.Sub(triangleVertices[0], triangleVertices[1])
	var t12 math.V3 = v3.Sub(triangleVertices[1], triangleVertices[2])
	var t20 math.V3 = v3.Sub(triangleVertices[2], triangleVertices[2])
	var triangleNormal math.V3 = v3.Cross(t01, t20)
	var triangleOffset float32 = v3.Dot(t01, t20)
	boxMin, boxMax = ProjectPointsOnAxis(triangleVertices, triangleNormal)
	if (boxMax < triangleOffset) || (boxMin > triangleOffset) {
		out = false
		return
	}

	var axis math.V3

	axis = v3.Cross(t01, bn0)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t01, bn1)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t01, bn2)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t12, bn0)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t12, bn1)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t12, bn2)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t20, bn0)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t20, bn1)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	axis = v3.Cross(t20, bn2)
	boxMin, boxMax = ProjectPointsOnAxis(boxVertices, axis)
	triangleMin, triangleMax = ProjectPointsOnAxis(triangleVertices, axis)
	if (boxMax < triangleMin) || (boxMin > triangleMax) {
		out = false
		return
	}

	out = true
	return
}

// RayIntersectsPlane ...
func RayIntersectsPlane(r0 math.V3, r1 math.V3, p math.V3, n math.V3) (inter bool, time float32) {
	var ls math.V3 = v3.Sub(r1, r0)
	var ll float32 = v3.Length(ls)
	var l math.V3 = v3.Divf(ls, ll)
	var d float32 = v3.Dot(n, l)
	if math.Abs_f32(d) > 0.0000001 {
		var pr0 math.V3 = v3.Sub(p, r0)
		inter = true
		time = (v3.Dot(pr0, n) / d) / ll
		if d < 0.0 {
			d = 1.0 - d
		}
	} else {
	}
	return
}

// RayIntersectsAABB ...
func RayAABBIntersection(p0 math.V3, p1 math.V3, min math.V3, max math.V3) (time float32) {
	time = -1.0
	var dx float32 = p1.X - p0.X
	var dy float32 = p1.Y - p0.Y
	var dz float32 = p1.Z - p0.Z

	var tmp float32

	var txmin float32 = (min.X - p0.X) / dx
	var txmax float32 = (max.X - p0.X) / dx
	if txmin > txmax {
		tmp = txmin
		txmin = txmax
		txmax = tmp
	}

	var tymin float32 = (min.Y - p0.Y) / dy
	var tymax float32 = (max.Y - p0.Y) / dy
	if tymin > tymax {
		tmp = tymin
		tymin = tymax
		tymax = tmp
	}

	if (txmin > tymax) || (tymin > txmax) {
	} else {
		if tymin > txmin {
			txmin = tymin
		}
		if tymax < txmax {
			txmax = tymax
		}

		var tzmin float32 = (min.Z - p0.Z) / dz
		var tzmax float32 = (max.Z - p0.Z) / dz
		if tzmin > tzmax {
			tmp = tzmin
			tzmin = tzmax
			tzmax = tmp
		}

		if (txmin > tzmax) || (tzmin > txmax) {
		} else {
			if tzmin > txmin {
				txmin = tzmin
			}
			if tzmax < txmax {
				txmax = tzmax
			}
			time = txmax / math.Sqrt_f32(dx*dx+dy*dy+dz*dz)
		}
	}
	return
}

// AABBIntersectsPlane ...
func AABBIntersectsPlane(min math.V3, max math.V3, plane math.V4) (out int32) {
	var pos math.V3
	var neg math.V3

	if plane.X >= 0.0 {
		pos.X = max.X
		neg.X = min.X
	} else {
		pos.X = min.X
		neg.X = max.X
	}

	if plane.Y >= 0.0 {
		pos.Y = max.Y
		neg.Y = min.Y
	} else {
		pos.Y = min.Y
		neg.Y = max.Y
	}

	if plane.Z >= 0.0 {
		pos.Z = max.Z
		neg.Z = min.Z
	} else {
		pos.Z = min.Z
		neg.Z = max.Z
	}

	if (plane.X*neg.X + plane.Y*neg.Y + plane.Z*neg.Z + plane.W) > 0.0 {
		out = -1
		//return
	} else if (plane.X*pos.X + plane.Y*pos.Y + plane.Z*pos.Z + plane.W) < 0.0 {
		out = 1
		//return
	} else {
		out = 0
	}
	return
}

// AABBIntersectsPlanes ...
func AABBIntersectsPlanes(min math.V3, max math.V3, planes []math.V4) (out int32) {
	out = -1
	// TODO REMOVE var intersect bool
	var planeCount int = len(planes)
	for i := 0; i < planeCount; i++ {
		var inter int32 = AABBIntersectsPlane(min, max, planes[i])
		if inter == 1 {
			out = 1
			return
		} else if inter == 0 {
			out = 0
		} else if inter == -1 {
		}
	}
	return
}

func triangleSign(p0 math.V3, p1 math.V3, p2 math.V3) (sign float32) {
	sign = (p0.X-p2.X)*(p1.Y-p2.Y) - (p1.X-p2.X)*(p0.Y-p2.Y)
	return
}

func pointInTriangle(p math.V3, t0 math.V3, t1 math.V3, t2 math.V3) (collision bool) {
	var pt0 math.V3 = v3.Normalize(v3.Sub(t0, p))
	var pt1 math.V3 = v3.Normalize(v3.Sub(t1, p))
	var pt2 math.V3 = v3.Normalize(v3.Sub(t2, p))

	var a float32 = v3.Dot(pt0, pt1)
	var b float32 = v3.Dot(pt1, pt2)
	var c float32 = v3.Dot(pt2, pt0)

	var angle float32 = math.Acos_f32(a) + math.Acos_f32(b) + math.Acos_f32(c)

	//printf("POINT_IN_TRIANGLE p %s, t0 %s, t1 %s, t2 %s, a %f, b %f, c %f, %f\n", v3.to_str(p), v3.to_str(t0), v3.to_str(t1), v3.to_str(t2), a, b, c, math.Abs_f32(angle - (math.M2PI_f32)))
	collision = math.Abs_f32(angle-(math.M2PI_f32)) < 0.01
	return
}

func pointInTriangle2(p math.V3, t0 math.V3, t1 math.V3, t2 math.V3) (collision bool) {
	// u=P2−P1
	var u math.V3 = v3.Sub(t1, t0)
	var v math.V3 = v3.Sub(t2, t0)
	var n math.V3 = v3.Cross(u, v)
	var ndn float32 = v3.Dot(n, n)
	var gamma float32 = v3.Dot(v3.Cross(u, p), n)
	gamma = gamma / ndn
	var beta float32 = v3.Dot(v3.Cross(p, v), n)
	beta = beta / ndn
	var alpha float32 = 1.0 - gamma - beta
	// The point P′ lies inside T if:
	return ((0.0 <= alpha) && (alpha <= 1.0) &&
		(0.0 <= beta) && (beta <= 1.0) &&
		(0.0 <= gamma) && (gamma <= 1.0))
}

func SameSide(p1 math.V3, p2 math.V3, a math.V3, b math.V3) (out bool) {
	var cp1 math.V3 = v3.Cross(v3.Sub(b, a), v3.Sub(p1, a))
	var cp2 math.V3 = v3.Cross(v3.Sub(b, a), v3.Sub(p2, a))
	if v3.Dot(cp1, cp2) >= 0.0 {
		out = true
	}
	return
}

func pointInTriangle3(p math.V3, a math.V3, b math.V3, c math.V3) (out bool) {
	if SameSide(p, a, b, c) {
		if SameSide(p, b, a, c) {
			if SameSide(p, c, a, b) {
				out = true
			}
		}
	}
	return
}

func pointInTriangle4(P math.V3, A math.V3, B math.V3, C math.V3) (out bool) {
	var v0 math.V3 = v3.Sub(C, A)
	var v1 math.V3 = v3.Sub(B, A)
	var v2 math.V3 = v3.Sub(P, A)

	var dot00 float32 = v3.Dot(v0, v0)
	var dot01 float32 = v3.Dot(v0, v1)
	var dot02 float32 = v3.Dot(v0, v2)
	var dot11 float32 = v3.Dot(v1, v1)
	var dot12 float32 = v3.Dot(v1, v2)

	var invDenom float32 = 1.0 / (dot00*dot11 - dot01*dot01)
	var u float32 = (dot11*dot02 - dot01*dot12) * invDenom
	var v float32 = (dot00*dot12 - dot01*dot02) * invDenom

	out = (u >= 0.0) && (v >= 0.0) && (u+v < 1.0)
	return
}

/*func pointInTriangle (pt math.V3, p1 math.V3, p2 math.V3, p3 v3) (inside bool) {
    var d0 float32 = triangleSign(pt, p1, p2)
    var d1 float32 = triangleSign(pt, p2, p3)
    var d2 float32 = triangleSign(pt, p3, p1)
    inside = (((d0 < 0.0) || (d1 < 0.0) || (d2 < 0.0)) && ((d0 > 0.0) || (d1 > 0.0) || (d2 > 0.0))) == false
}*/

func getLowestRoot(a float32, b float32, c float32) (out float32) {
	out = math.MAX_f32
	var det float32 = b*b - 4.0*a*c
	if det >= 0.0 {
		var sqrtDet float32 = math.Sqrt_f32(det)
		var r1 float32 = (-b - sqrtDet) / (2.0 * a)
		var r2 float32 = (-b + sqrtDet) / (2.0 * a)

		if r1 > r2 {
			var tmp float32 = r2
			r2 = r1
			r1 = tmp
		}

		if r1 > 0.0 {
			out = r1
		} else if r2 > 0.0 {
			out = r2
		}
	}
	return
}

func testVertex(p math.V3, velSqrLen float32, start math.V3, vel math.V3) (out float32) {
	out = math.MAX_f32
	var v math.V3 = v3.Sub(start, p)
	var b float32 = 2.0 * v3.Dot(vel, v)
	var c float32 = v3.Sqlength(v) - 1.0
	var newT float32 = getLowestRoot(velSqrLen, b, c)
	if newT <= 1.0 {
		out = newT
	}
	return
}

func testEdge(index int32, pa math.V3, pb math.V3, velSqrLen float32, start math.V3, vel math.V3) (out float32, intersection math.V3) {
	out = math.MAX_f32
	var edge math.V3 = v3.Sub(pb, pa)
	var v math.V3 = v3.Sub(pa, start)

	var edgeSqrLen float32 = v3.Sqlength(edge)
	var edgeDotVel float32 = v3.Dot(edge, vel)
	var edgeDotSphereVert float32 = v3.Dot(edge, v)

	var a float32 = edgeSqrLen*(0.0-velSqrLen) + edgeDotVel*edgeDotVel
	var b float32 = edgeSqrLen*(2.0*v3.Dot(vel, v)) - 2.0*edgeDotVel*edgeDotSphereVert
	var c float32 = edgeSqrLen*(1.0-v3.Sqlength(v)) + edgeDotSphereVert*edgeDotSphereVert

	var newT float32 = getLowestRoot(a, b, c)
	if newT <= 1.0 {
		var f float32 = (edgeDotVel*newT - edgeDotSphereVert) / edgeSqrLen
		if f >= 0.0 && f <= 1.0 {
			v = v3.Mulf(edge, f)
			v = v3.Add(v, pa)
			out = newT
			intersection = v
		}
	}
	return
}

type TriangleIntersection struct {
	distance float32
	point    math.V3
	normal   math.V3
	base     math.V3
}

type SphereCollision struct {
	center                     math.V3
	velocity                   math.V3
	destination                math.V3
	radius                     math.V3
	invRadius                  math.V3
	velocityLength             float32
	unitVelocity               math.V3
	scaledVelocity             math.V3
	scaledVelocitySquareLength float32
	scaledCenter               math.V3
	minTime                    float32
	count                      int32
	tag0                       int32
	tag1                       int32
	intersections              []TriangleIntersection
}

func SphereCollisionCreate(center math.V3, velocity math.V3, radius math.V3) {
	g_SphereCollision.center = center
	g_SphereCollision.destination = v3.Add(center, velocity)
	g_SphereCollision.velocity = velocity
	g_SphereCollision.radius = radius
	g_SphereCollision.invRadius = v3.Div(v3.ONE, radius)
	g_SphereCollision.velocityLength = v3.Length(velocity)
	if g_SphereCollision.velocityLength > 0.0 {
		g_SphereCollision.unitVelocity = v3.Divf(velocity, g_SphereCollision.velocityLength)
	} else {
		g_SphereCollision.unitVelocity = v3.ZERO
	}
	g_SphereCollision.scaledVelocity = v3.Div(g_SphereCollision.velocity, radius)
	g_SphereCollision.scaledVelocitySquareLength = v3.Sqlength(g_SphereCollision.scaledVelocity)
	g_SphereCollision.scaledCenter = v3.Div(center, radius)
	g_SphereCollision.intersections = g_SphereCollision.intersections[:0]
	g_SphereCollision.minTime = math.MAX_f32
	g_SphereCollision.count = 0
}

func SphereCollisionAppendIntersection(distance float32, point math.V3, normal math.V3, base math.V3) (out bool) {
	var inter TriangleIntersection

	inter.distance = distance
	inter.point = point
	inter.normal = normal
	inter.base = base
	var intersections []TriangleIntersection = g_SphereCollision.intersections

	var intersectionCount int = len(intersections)
	if intersectionCount == 0 {
		intersections = append(intersections, inter)
		out = true
		g_SphereCollision.count = g_SphereCollision.count + 1
	} else if distance < intersections[0].distance {
		intersections[0] = inter
		out = true
		g_SphereCollision.count = g_SphereCollision.count + 1
	}
	g_SphereCollision.intersections = intersections
	g_SphereCollision.minTime = distance
	return
}

var g_SphereCollision SphereCollision

func SphereIntersectsPlane1(pt0 math.V3, normal math.V3) {

	var nx float32 = normal.X
	var ny float32 = normal.Y
	var nz float32 = normal.Z

	var ndn float32 = nx*nx + ny*ny + nz*nz
	//printf("POSI %d, NDN %f\n", posI, ndn)
	if ndn > 0.0 {
		if (nx*g_SphereCollision.unitVelocity.X + ny*g_SphereCollision.unitVelocity.Y + nz*g_SphereCollision.unitVelocity.Z) < 0.0 {
			//printf("NORMAL %s, VELOCITY %s, DOT %f\n", v3.to_str(normal), v3.to_str(g_SphereCollision.unitVelocity), v3.Dot(normal, g_SphereCollision.unitVelocity))
			var rx float32 = g_SphereCollision.invRadius.X
			var ry float32 = g_SphereCollision.invRadius.Y
			var rz float32 = g_SphereCollision.invRadius.Z

			var scaledP0 math.V3
			scaledP0.X = pt0.X * rx
			scaledP0.Y = pt0.Y * ry
			scaledP0.Z = pt0.Z * rz

			var pD float32 = 0.0 - nx*scaledP0.X - ny*scaledP0.Y - nz*scaledP0.Z

			var embedded bool = false

			var distToPlane float32 = nx*g_SphereCollision.scaledCenter.X + ny*g_SphereCollision.scaledCenter.Y + nz*g_SphereCollision.scaledCenter.Z + pD
			var normDotVel float32 = nx*g_SphereCollision.scaledVelocity.X + ny*g_SphereCollision.scaledVelocity.Y + nz*g_SphereCollision.scaledVelocity.Z

			var t0 float32 = 0.0
			var t1 float32 = 0.0
			var tbreak bool = false
			if normDotVel >= -0.000001 && normDotVel <= 0.000001 {
				if distToPlane <= -1.0 || distToPlane >= 1.0 {
					tbreak = true
				} else {
					embedded = true
					t0 = 0.0
					t1 = 1.0
				}
			} else {
				t0 = (1.0 - distToPlane) / normDotVel
				t1 = (-1.0 - distToPlane) / normDotVel
				if t0 > t1 {
					var temp float32 = t1
					t1 = t0
					t0 = temp
				}

				if t0 > 1.0 || t1 < 0.0 {
					tbreak = true
				}

				if t0 < 0.0 {
					t0 = 0.0
				}
				if t1 > 1.0 {
					t1 = 1.0
				}
			}

			if tbreak == false {
				if embedded == false {
					var svx float32 = g_SphereCollision.scaledVelocity.X
					var svy float32 = g_SphereCollision.scaledVelocity.Y
					var svz float32 = g_SphereCollision.scaledVelocity.Z

					var planeIntersect math.V3
					planeIntersect.X = g_SphereCollision.scaledCenter.X - nx + svx*t0
					planeIntersect.Y = g_SphereCollision.scaledCenter.Y - ny + svy*t0
					planeIntersect.Z = g_SphereCollision.scaledCenter.Z - nz + svz*t0

					_ /*var collided bool*/ = SphereCollisionAppendIntersection(t0, planeIntersect, normal, scaledP0)
				}
			}
		}
	}
}

func SphereIntersectsTriangle(positions []float32, normals []float32, debugTriangles []math.V3) (out []math.V3) {
	out = debugTriangles
	var rx float32 = g_SphereCollision.invRadius.X
	var ry float32 = g_SphereCollision.invRadius.Y
	var rz float32 = g_SphereCollision.invRadius.Z

	var svx float32 = g_SphereCollision.scaledVelocity.X
	var svy float32 = g_SphereCollision.scaledVelocity.Y
	var svz float32 = g_SphereCollision.scaledVelocity.Z

	var nx float32
	var ny float32
	var nz float32
	var ndn float32

	var pD float32

	var embedded bool

	var scaledP0 math.V3
	var scaledP1 math.V3
	var scaledP2 math.V3

	var tbreak bool
	var planeIntersect math.V3

	var posCount int32 = int32((len(positions) / 9) * 9)
	var posI int32
	var posN int32

	var disp float32
	var inter math.V3

	var t0 float32
	var t1 float32

	var distToPlane float32
	var normDotVel float32

	var temp float32

	// var v02x float32
	// var v02y float32
	// var v02z float32

	// var v01x float32
	// var v01y float32
	// var v01z float32

	// var v0Ix float32
	// var v0Iy float32
	// var v0Iz float32

	// var dot00 float32
	// var dot01 float32
	// var dot02 float32
	// var dot11 float32
	// var dot12 float32

	// var invDenom float32
	// var u float32
	// var v float32

	// var f float32
	// var b float32
	// var c float32

	// var v0x float32
	// var v0y float32
	// var v0z float32
	// var svv0 float32
	// var sqv0 float32
	// var newT float32
	var normal math.V3
	var debugNormal math.V3
	//printf("POS_COUNT %d\n", posCount)
	for posI < posCount {
		nx = normals[posN]
		ny = normals[posN+1]
		nz = normals[posN+2]

		normal = v3.Make(nx, ny, nz)
		debugNormal = v3.Mulf(normal, 0.03)

		ndn = nx*nx + ny*ny + nz*nz
		//printf("POSI %d, NDN %f\n", posI, ndn)
		if ndn > 0.0 {
			if (nx*g_SphereCollision.unitVelocity.X + ny*g_SphereCollision.unitVelocity.Y + nz*g_SphereCollision.unitVelocity.Z) < 0.0 {
				//printf("NORMAL %s, VELOCITY %s, DOT %f\n", v3.to_str(normal), v3.to_str(g_SphereCollision.unitVelocity), v3.Dot(normal, g_SphereCollision.unitVelocity))
				scaledP0.X = positions[posI] * rx
				scaledP0.Y = positions[posI+1] * ry
				scaledP0.Z = positions[posI+2] * rz

				pD = 0.0 - nx*scaledP0.X - ny*scaledP0.Y - nz*scaledP0.Z

				embedded = false

				distToPlane = nx*g_SphereCollision.scaledCenter.X + ny*g_SphereCollision.scaledCenter.Y + nz*g_SphereCollision.scaledCenter.Z + pD
				normDotVel = nx*g_SphereCollision.scaledVelocity.X + ny*g_SphereCollision.scaledVelocity.Y + nz*g_SphereCollision.scaledVelocity.Z

				t0 = 0.0
				t1 = 0.0
				tbreak = false
				if normDotVel >= -0.000001 && normDotVel <= 0.000001 {
					if distToPlane <= -1.0 || distToPlane >= 1.0 {
						tbreak = true
					} else {
						embedded = true
						t0 = 0.0
						t1 = 1.0
					}
				} else {
					t0 = (1.0 - distToPlane) / normDotVel
					t1 = (-1.0 - distToPlane) / normDotVel
					if t0 > t1 {
						temp = t1
						t1 = t0
						t0 = temp
					}

					if t0 > 1.0 || t1 < 0.0 {
						tbreak = true
					}

					if t0 < 0.0 {
						t0 = 0.0
					}
					if t1 > 1.0 {
						t1 = 1.0
					}
				}

				if tbreak == false {
					scaledP1.X = positions[posI+3] * rx
					scaledP1.Y = positions[posI+4] * ry
					scaledP1.Z = positions[posI+5] * rz

					scaledP2.X = positions[posI+6] * rx
					scaledP2.Y = positions[posI+7] * ry
					scaledP2.Z = positions[posI+8] * rz

					if embedded == false {
						planeIntersect.X = g_SphereCollision.scaledCenter.X - nx + svx*t0
						planeIntersect.Y = g_SphereCollision.scaledCenter.Y - ny + svy*t0
						planeIntersect.Z = g_SphereCollision.scaledCenter.Z - nz + svz*t0

						if pointInTriangle4(planeIntersect, scaledP0, scaledP1, scaledP2) {
							if SphereCollisionAppendIntersection(t0, planeIntersect, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d TRIANGLE DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, t0, v3.to_str(planeIntersect), v3.to_str(v3.Make(nx, ny, nz)))
								tbreak = true
							}
						}
					}

					if tbreak == false {
						disp = testVertex(scaledP0, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, scaledP0, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d VERTEX 0 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(scaledP0), v3.to_str(v3.Make(nx, ny, nz)))
							}

						}

						disp = testVertex(scaledP1, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, scaledP1, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d VERTEX 1 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(scaledP1), v3.to_str(v3.Make(nx, ny, nz)))
							}
						}

						disp = testVertex(scaledP2, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, scaledP2, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								// printf("%d VERTEX 2 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(scaledP2), v3.to_str(v3.Make(nx, ny, nz)))
							}
						}

						disp, inter = testEdge(0, scaledP0, scaledP1, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, inter, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d EDGE 0 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(inter), v3.to_str(v3.Make(nx, ny, nz)))
							}
						}

						disp, inter = testEdge(1, scaledP1, scaledP2, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, inter, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d EDGE 1 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(inter), v3.to_str(v3.Make(nx, ny, nz)))
							}
						}

						disp, inter = testEdge(2, scaledP2, scaledP0, g_SphereCollision.scaledVelocitySquareLength,
							g_SphereCollision.scaledCenter, g_SphereCollision.scaledVelocity)
						if disp != math.MAX_f32 {
							if SphereCollisionAppendIntersection(disp, inter, normal, scaledP0) {
								out = append(out, v3.Mul(v3.Add(scaledP0, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP1, debugNormal), g_SphereCollision.radius))
								out = append(out, v3.Mul(v3.Add(scaledP2, debugNormal), g_SphereCollision.radius))
								//printf("%d EDGE 2 DISTANCE %f, INTER %s NORMAL %s\n", posI / 9, disp, v3.to_str(inter), v3.to_str(v3.Make(nx, ny, nz)))
							}
						}
					}
				}
			}
		}

		posI = posI + 9
		posN = posN + 3
	}
	return
}

var it0 fps.ProfileId = fps.InvalidProfile()
var it1 fps.ProfileId = fps.InvalidProfile()
var it2 fps.ProfileId = fps.InvalidProfile()
var it3 fps.ProfileId = fps.InvalidProfile()
var it4 fps.ProfileId = fps.InvalidProfile()
var it5 fps.ProfileId = fps.InvalidProfile()
var it6 fps.ProfileId = fps.InvalidProfile()
var it7 fps.ProfileId = fps.InvalidProfile()
var it8 fps.ProfileId = fps.InvalidProfile()
var it9 fps.ProfileId = fps.InvalidProfile()

// SphereIntersectsPlane ...
func SphereIntersectsPlane(center math.V3, radius float32, plane math.V4, point math.V3) (out int32) {
	out = -1

	var dx float32 = center.X - point.X
	var dy float32 = center.Y - point.Y
	var dz float32 = center.Z - point.Z

	var d float32 = -dx*plane.X + -dy*plane.Y + -dz*plane.Z
	if d > radius {
		out = 1
	} else if d < -radius {
		out = -1
	} else {
		out = 0
	}
	return
}

// SphereIntersectsPlanes ...
func SphereIntersectsPlanes(center math.V3, radius float32, planes []math.V4, points []math.V3) (out int32) {
	out = -1
	// TODO REMOVE var intersect bool
	var planeCount int = len(planes)
	for i := 0; i < planeCount; i++ {
		var inter int32 = SphereIntersectsPlane(center, radius, planes[i], points[i])
		if inter == 1 {
			out = 1
			return
		} else if inter == 0 {
			out = 0
		} else if inter == -1 {
		}
	}
	return
}

// SphereIntersectsAABB ...
func SphereIntersectsAABB(center math.V3, radius float32, min math.V3, max math.V3) (out int32) {
	out = 1
	var dmin float32
	var dmax float32
	var r2 float32 = radius * radius

	var minx float32 = center.X - min.X
	var miny float32 = center.Y - min.Y
	var minz float32 = center.Z - min.Z

	var maxx float32 = center.X - max.X
	var maxy float32 = center.Y - max.Y
	var maxz float32 = center.Z - max.Z

	var sqminx float32 = minx * minx
	var sqminy float32 = miny * miny
	var sqminz float32 = minz * minz

	var sqmaxx float32 = maxx * maxx
	var sqmaxy float32 = maxy * maxy
	var sqmaxz float32 = maxz * maxz

	if minx < 0.0 {
		dmin = dmin + sqminx
	} else if maxx > 0.0 {
		dmin = dmin + sqmaxx
	}
	dmax = dmax + math.Max_f32(sqminx, sqmaxx)

	if miny < 0.0 {
		dmin = dmin + sqminy
	} else if maxy > 0.0 {
		dmin = dmin + sqmaxy
	}
	dmax = dmax + math.Max_f32(sqminy, sqmaxy)

	if minz < 0.0 {
		dmin = dmin + sqminz
	} else if maxz > 0.0 {
		dmin = dmin + sqmaxz
	}
	dmax = dmax + math.Max_f32(sqminz, sqmaxz)

	if dmax <= r2 {
		out = -1
	} else if dmin <= r2 {
		out = 0
	}
	return
}
