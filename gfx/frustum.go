package gfx

import (
	"skyfx/math"
	"skyfx/math/intersect"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Globals ...
var g_frustums []Frustum

// Frustum ...
type Frustum struct {
	cornersV4 []math.V4 // TODO : issue when using array
	corners   []math.V3 // TODO : issue when using array
	planes    []math.V4 // TODO : issue when using array
	points    []math.V3
	debugMesh MeshId
}

// FrustumId
type FrustumId struct {
	frustum int32
}

// InvalidFrustum ...
func InvalidFrustum() (out FrustumId) {
	out.frustum = -1
	return
}

// FrustumCreate ...
func FrustumCreate() (out FrustumId) {
	out.frustum = int32(len(g_frustums))
	var frustum Frustum
	frustum.debugMesh = MeshCreate(gl.LINES, gl.UNSIGNED_SHORT, 1024*3, VertexLayout, 1024*3)

	frustum.planes = append(frustum.planes, v4.ZERO)
	frustum.planes = append(frustum.planes, v4.ZERO)
	frustum.planes = append(frustum.planes, v4.ZERO)
	frustum.planes = append(frustum.planes, v4.ZERO)
	frustum.planes = append(frustum.planes, v4.ZERO)
	frustum.planes = append(frustum.planes, v4.ZERO)

	frustum.points = append(frustum.points, v3.ZERO)
	frustum.points = append(frustum.points, v3.ZERO)
	frustum.points = append(frustum.points, v3.ZERO)
	frustum.points = append(frustum.points, v3.ZERO)
	frustum.points = append(frustum.points, v3.ZERO)
	frustum.points = append(frustum.points, v3.ZERO)

	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)
	frustum.cornersV4 = append(frustum.cornersV4, v4.ZERO)

	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)
	frustum.corners = append(frustum.corners, v3.ZERO)

	g_frustums = append(g_frustums, frustum)
	return
}

func planeFromPoints(a math.V3, b math.V3, c math.V3) (out math.V4) {
	var ab math.V3 = v3.Sub(b, a)
	var ac math.V3 = v3.Sub(c, a)
	var cross math.V3 = v3.Cross(ab, ac)
	var normal math.V3 = v3.Normalize(cross)
	out.X = normal.X
	out.Y = normal.Y
	out.Z = normal.Z
	out.W = 0.0 - v3.Dot(normal, a)
	return
}

// FrustumUpdate ...
func FrustumUpdate(id FrustumId, invViewProj math.M44) {
	var cornersV4 []math.V4 = g_frustums[id.frustum].cornersV4
	cornersV4[0] = v4.Make(-1.0, -1.0, -1.0, 1.0)
	cornersV4[1] = v4.Make(1.0, -1.0, -1.0, 1.0)
	cornersV4[2] = v4.Make(-1.0, -1.0, 1.0, 1.0)
	cornersV4[3] = v4.Make(1.0, -1.0, 1.0, 1.0)
	cornersV4[4] = v4.Make(-1.0, 1.0, -1.0, 1.0)
	cornersV4[5] = v4.Make(1.0, 1.0, -1.0, 1.0)
	cornersV4[6] = v4.Make(-1.0, 1.0, 1.0, 1.0)
	cornersV4[7] = v4.Make(1.0, 1.0, 1.0, 1.0)

	var corners []math.V3 = g_frustums[id.frustum].corners
	for i := 0; i < 8; i++ {
		cornersV4[i] = v4.Transform(cornersV4[i], invViewProj)
		corners[i] = v3.Divf(v3.Make(cornersV4[i].X, cornersV4[i].Y, cornersV4[i].Z), cornersV4[i].W)
	}

	var planes []math.V4 = g_frustums[id.frustum].planes
	planes[0] = planeFromPoints(corners[0], corners[1], corners[2])
	planes[1] = planeFromPoints(corners[4], corners[6], corners[5])
	planes[2] = planeFromPoints(corners[4], corners[0], corners[6])
	planes[3] = planeFromPoints(corners[5], corners[7], corners[1])
	planes[4] = planeFromPoints(corners[0], corners[4], corners[5])
	planes[5] = planeFromPoints(corners[2], corners[7], corners[6])

	var points []math.V3 = g_frustums[id.frustum].points
	points[0] = corners[0]
	points[1] = corners[4]
	points[2] = corners[4]
	points[3] = corners[5]
	points[4] = corners[0]
	points[5] = corners[2]

	var debugMesh MeshId = g_frustums[id.frustum].debugMesh
	MeshBegin(debugMesh)

	MeshAppendLine(debugMesh, corners[0], corners[1], v4.BLUE)
	MeshAppendLine(debugMesh, corners[1], corners[3], v4.BLUE)
	MeshAppendLine(debugMesh, corners[3], corners[2], v4.BLUE)
	MeshAppendLine(debugMesh, corners[2], corners[0], v4.BLUE)

	MeshAppendLine(debugMesh, corners[4], corners[5], v4.BLUE)
	MeshAppendLine(debugMesh, corners[5], corners[7], v4.BLUE)
	MeshAppendLine(debugMesh, corners[7], corners[6], v4.BLUE)
	MeshAppendLine(debugMesh, corners[6], corners[4], v4.BLUE)

	MeshAppendLine(debugMesh, corners[0], corners[4], v4.BLUE)
	MeshAppendLine(debugMesh, corners[1], corners[5], v4.BLUE)
	MeshAppendLine(debugMesh, corners[2], corners[6], v4.BLUE)
	MeshAppendLine(debugMesh, corners[3], corners[7], v4.BLUE)
	MeshEnd(debugMesh)
}

// FrustumIntersectsAABB ...
func FrustumIntersectsAABB(id FrustumId, min math.V3, max math.V3) (out int32) {
	out = intersect.AABBIntersectsPlanes(min, max, g_frustums[id.frustum].planes)
	return
}

// // FrustumIntersectsSphere ...
// func FrustumIntersectsSphere(id FrustumId, center math.V3, radius float32) (out int32) {
//     out = mat.SphereIntersectsPlanes(center, radius, g_frustums[id.frustum].planes, g_frustums[id.frustum].points)
// }

// FrustumRender ...
func FrustumRender(id FrustumId, world math.M44, view math.M44, projection math.M44) {
	DepthState(true, gl.LESS, false)
	EffectUse(g_fxVertexColor3D)
	EffectAssignM44(g_fxVertexColor3D, UNIFORM_WORLD, world, false)
	EffectAssignM44(g_fxVertexColor3D, UNIFORM_VIEW, view, false)
	EffectAssignM44(g_fxVertexColor3D, UNIFORM_PROJECTION, projection, false)
	MeshRender(g_frustums[id.frustum].debugMesh)
}
