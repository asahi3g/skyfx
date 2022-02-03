package gfx

import (
	"fmt"
	gomath "math"
	"skyfx/math"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"
	"skyfx/utils"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// import "fps"
// import "mat"
// import "v1"
// import "v3"
// import "v4"
// import "q4"
// import "m44"

// // TODO : stop using hardcoded vertex attributes in MeshAppend*
// // TODO : use DrawElement

// Globals ...
var g_vaos []uint32
var g_vbos []uint32
var g_meshes []Mesh

var g_freeMeshes []MeshId

// VertexAttribute ...
type VertexAttribute struct {
	componentType     uint32
	componentCount    uint32
	componentByteSize uint32
	componentOffset   uint32
	byteOffset        uint32
	binding           uint32
}

// VertexAttributeCreate ...
func VertexAttributeCreate(binding uint32, componentCount uint32, componentType uint32) (out VertexAttribute) {
	out.binding = binding
	out.componentCount = componentCount
	out.componentType = componentType
	if componentType == gl.FLOAT {
		out.componentByteSize = 4
	}
	return
}

// VertexAttributePrint ...
func VertexAttributePrint(a VertexAttribute) {
	fmt.Printf("componentCount %d, componentType %d, componentByteSize %d, componentOffset %d, byteOffset %d, binding %d\n",
		a.componentCount, a.componentType, a.componentByteSize, a.componentOffset, a.byteOffset, a.binding)
}

// MeshId ...
type MeshId struct {
	mesh int32
}

// Mesh ...
type Mesh struct {
	id MeshId

	primitive uint32
	index     int32

	indices              []uint8
	indexFormat          uint32
	indexByteStride      int32
	indexComponentStride int32
	indexByteCount       int32
	indexCount           int32

	vertices              []uint8
	attributes            []VertexAttribute
	vertexByteStride      uint32
	vertexComponentStride uint32
	vertexByteCount       uint32
	vertexCount           uint32

	channels  []channelInfo
	frontFace uint32
	cullFace  uint32

	min math.V3
	max math.V3

	vao         uint32
	ibo         uint32
	vbo         uint32
	usePosition bool
	useNormal   bool
	useColor    bool
	useTexcoord bool
	useTangent  bool
	useWeight   bool
	useJoint    bool
}

// InvalidMesh ...
func InvalidMesh() (out MeshId) {
	out.mesh = -1
	return
}

// IsValidMesh ...
func IsValidMesh(id MeshId) (out bool) {
	out = int(id.mesh) >= 0 && int(id.mesh) < len(g_meshes)
	return
}

// MeshUsePosition ...
func MeshUsePosition(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].usePosition
	return
}

// MeshUseNormal ...
func MeshUseNormal(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useNormal
	return
}

// MeshUseColor ...
func MeshUseColor(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useColor
	return
}

// MeshUseTexcoord ...
func MeshUseTexcoord(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useTexcoord
	return
}

// MeshUseTangent ...
func MeshUseTangent(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useTangent
	return
}

// MeshUseWeight ...
func MeshUseWeight(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useWeight
	return
}

// MeshUseJoint ...
func MeshUseJoint(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	out = g_meshes[id.mesh].useJoint
	return
}

// MeshSetCulling ...
func MeshSetCulling(id MeshId, frontFace uint32, cullFace uint32) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	g_meshes[id.mesh].frontFace = frontFace
	g_meshes[id.mesh].cullFace = cullFace
}

// MeshPrint ...
func MeshPrint(message string, id MeshId) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var index int32 = id.mesh
	fmt.Printf("id %d, indices %d, vertices %d, attributes %d, primitive %d, vao %d, vbo %d, ibo %d, vertexByteStride %d, vertexComponentStride %d, vertexCount %d\n",
		index,
		len(g_meshes[index].indices),
		len(g_meshes[index].vertices),
		len(g_meshes[index].attributes),
		g_meshes[index].primitive,
		g_meshes[index].vao,
		g_meshes[index].vbo,
		g_meshes[index].ibo,
		g_meshes[index].vertexByteStride,
		g_meshes[index].vertexComponentStride,
		g_meshes[index].vertexCount)
}

// MeshIsEmpty ...
func MeshIsEmpty(id MeshId) (out bool) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var mesh int32 = id.mesh
	var index int = len(g_meshes[mesh].vertices)
	out = index <= 0
	return
}

// MeshLock ...
func MeshLock(primitive uint32, indexFormat uint32, indexCount int32, attributes []VertexAttribute, vertexCount uint32) (out MeshId) {
	var freeMeshCount int = len(g_freeMeshes)
	if freeMeshCount > 0 {
		var newFreeMeshCount int = freeMeshCount - 1
		out.mesh = g_freeMeshes[newFreeMeshCount].mesh
		g_freeMeshes = g_freeMeshes[:newFreeMeshCount]
	} else {
		out = MeshCreate(primitive, indexFormat, indexCount, attributes, vertexCount)
	}
	utils.PanicIfNot(IsValidMesh(out), "invalid id")
	return
}

// MeshUnlock ...
func MeshUnlock(id MeshId) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	g_freeMeshes = append(g_freeMeshes, id)
}

func meshCreate(primitive uint32, indexType uint32, indexCount int32, indices []uint8, attributes []VertexAttribute, vertexCount uint32, vertices []uint8, usage uint32) (out MeshId) {
	out.mesh = int32(len(g_meshes))

	if primitive == gl.TRIANGLES {
		utils.PanicIfNot((indexCount%3) == 0, "(indexCount % 3) == 0")
	} else if primitive == gl.LINES {
		utils.PanicIfNot((indexCount%2) == 0, "(indexCount % 2) == 0")
	}

	var mesh Mesh
	mesh.attributes = attributes
	mesh.primitive = primitive
	mesh.indexFormat = indexType
	mesh.frontFace = gl.CCW
	mesh.cullFace = gl.BACK

	// stride
	var indexByteStride int32 = 0
	if indexType == gl.UNSIGNED_SHORT {
		indexByteStride = 2
	} else if indexType == gl.UNSIGNED_INT {
		indexByteStride = 4
	} else {
		utils.PanicIf(true, "invalid index format")
	}
	mesh.indexByteStride = indexByteStride

	//
	var attributeCount int = len(attributes)
	var vertexByteStride uint32
	var vertexComponentStride uint32
	for i := 0; i < attributeCount; i++ {
		var attribute VertexAttribute = attributes[i]
		var componentCount uint32 = attribute.componentCount
		var componentByteSize uint32 = attribute.componentByteSize
		attribute.byteOffset = vertexByteStride
		attribute.componentOffset = vertexComponentStride
		//VertexAttributePrint(attribute)
		vertexComponentStride = vertexComponentStride + componentCount
		vertexByteStride = vertexByteStride + componentCount*componentByteSize
		attributes[i] = attribute
	}
	mesh.vertexByteStride = vertexByteStride
	mesh.vertexComponentStride = vertexComponentStride

	var indexByteCount int = len(indices)
	if indexCount > 0 {
		indexByteCount = int(indexCount * indexByteStride)
	}
	mesh.indexByteCount = int32(indexByteCount)

	var vertexByteCount uint32 = uint32(len(vertices))
	if vertexCount > 0 {
		vertexByteCount = vertexCount * vertexByteStride
	}
	mesh.vertexByteCount = vertexByteCount

	utils.PanicIfNot((primitive == gl.TRIANGLES || primitive == gl.LINES), "(primitive == gl.TRIANGLES || primitive == gl.LINES)")
	utils.PanicIfNot(vertexByteCount > 0, "vertexByteCount > 0")
	utils.PanicIfNot(indexCount > 0, "indexCount > 0")

	utils.PanicIfNot(vertexByteStride > 0, "mesh.vertexByteStride > 0")
	utils.PanicIfNot((vertexByteCount%vertexByteStride) == 0, "(vertexByteCount % vertexByteStride) == 0")

	// ibo
	gl.GenBuffers(1, &mesh.ibo)
	utils.PanicIf(GlError(), "gl.GenBuffers")
	g_vbos = append(g_vbos, mesh.ibo)

	// vbo
	gl.GenBuffers(1, &mesh.vbo)
	utils.PanicIf(GlError(), "gl.GenBuffers")
	g_vbos = append(g_vbos, mesh.vbo)

	// vao
	gl.GenVertexArrays(1, &mesh.vao)
	utils.PanicIf(GlError(), "gl.GenVertexArraysCore")
	g_vaos = append(g_vaos, mesh.vao)

	gl.BindVertexArray(mesh.vao)
	utils.PanicIf(GlError(), "gl.BindVertexArrayCore")

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	utils.PanicIf(GlError(), "gl.BindBuffer")

	var unsafeVB unsafe.Pointer
	if len(vertices) > 0 {
		unsafeVB = unsafe.Pointer(&vertices[0])
	}
	gl.BufferData(gl.ARRAY_BUFFER, int(vertexByteCount), unsafeVB, usage)
	utils.PanicIf(GlError(), "gl.BufferData")

	for i := 0; i < attributeCount; i++ {
		gl.EnableVertexAttribArray(attributes[i].binding)
		utils.PanicIf(GlError(), "gl.EnableVertexAttribArray")

		gl.VertexAttribPointerWithOffset(attributes[i].binding, int32(attributes[i].componentCount), attributes[i].componentType, false, int32(vertexByteStride), uintptr(attributes[i].byteOffset))
		utils.PanicIf(GlError(), "gl.VertexAttribPointerI32")
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ibo)
	utils.PanicIf(GlError(), "gl.BindBuffer")

	var unsafeIB unsafe.Pointer
	if len(indices) > 0 {
		unsafeIB = unsafe.Pointer(&indices[0])
	}
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexByteCount, unsafeIB, usage)
	utils.PanicIf(GlError(), "gl.BufferData")

	g_meshes = append(g_meshes, mesh) // ISSUE : can't use global slice if declared in another file
	utils.PanicIfNot(IsValidMesh(out), "invalid id")
	return
}

// MeshInstance ...
func MeshInstance(primitive uint32, indexFormat uint32, indices []uint8, attributes []VertexAttribute, vertices []uint8) (out MeshId) {
	out = meshCreate(primitive, indexFormat, 0, indices, attributes, 0, vertices, gl.STATIC_DRAW)
	return
}

// MeshCreate ...
func MeshCreate(primitive uint32, indexFormat uint32, indexCount int32, attributes []VertexAttribute, vertexCount uint32) (out MeshId) { // TODO: refactor
	var indices []uint8
	var vertices []uint8
	out = meshCreate(primitive, indexFormat, indexCount, indices, attributes, vertexCount, vertices, gl.STREAM_DRAW)
	return
}

// MeshBegin ...
func MeshBegin(id MeshId) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var mesh int32 = id.mesh
	// var lenIndices int = len(g_meshes[mesh].indices)
	// var lenVertices int = len(g_meshes[mesh].vertices)
	g_meshes[mesh].vertices = g_meshes[mesh].vertices[:0]
	g_meshes[mesh].indices = g_meshes[mesh].indices[:0]
}

// MeshEnd ...
func MeshEnd(id MeshId) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	//MeshPrint("MeshEnd : ", id)

	var mesh int32 = id.mesh
	var indexCount int32 = 0
	var vertexCount uint32 = 0
	var indexByteCount int32 = int32(len(g_meshes[mesh].indices))
	var vertexByteCount uint32 = uint32(len(g_meshes[mesh].vertices))
	if vertexByteCount > 0 && indexByteCount > 0 {

		var indexByteStride int32 = g_meshes[mesh].indexByteStride
		utils.PanicIfNot(indexByteStride > 0, "invalid index stride")
		utils.PanicIfNot((indexByteCount%indexByteStride) == 0, "(indexByteCount % indexByteStride) == 0)")
		indexCount = indexByteCount / indexByteStride

		var vertexByteStride uint32 = g_meshes[mesh].vertexByteStride
		utils.PanicIfNot(vertexByteStride > 0, "invalid vertex stride")
		utils.PanicIfNot((vertexByteCount%vertexByteStride) == 0, fmt.Sprintf("(vertexByteCount mod vertexByteStride) == 0 :: (%d mod %d) == %d", vertexByteCount, vertexByteStride, vertexByteCount%vertexByteStride))
		vertexCount = vertexByteCount / vertexByteStride

		gl.BindVertexArray(g_meshes[mesh].vao)
		utils.PanicIf(GlError(), "gl.BindVertexArrayCore")

		gl.BindBuffer(gl.ARRAY_BUFFER, g_meshes[mesh].vbo)
		utils.PanicIf(GlError(), "gl.BindBuffer")
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(vertexByteCount), unsafe.Pointer(&g_meshes[mesh].vertices[0]))
		utils.PanicIf(GlError(), "gl.BufferSubData")

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g_meshes[mesh].ibo)
		utils.PanicIf(GlError(), "gl.BindBuffer")
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, int(indexByteCount), unsafe.Pointer(&g_meshes[mesh].indices[0]))
		utils.PanicIf(GlError(), "gl.BufferSubData")
	}
	g_meshes[mesh].indexCount = indexCount
	g_meshes[mesh].vertexCount = vertexCount
	g_meshes[mesh].indexByteCount = indexByteCount
	g_meshes[mesh].vertexByteCount = vertexByteCount
}

func AppendUI16(indices []uint8, value uint16) []uint8 {
	indices = append(indices, byte(value))
	indices = append(indices, byte(value>>8))
	return indices
}

func AppendUI32(indices []uint8, value uint32) []uint8 {
	indices = append(indices, byte(value))
	indices = append(indices, byte(value>>8))
	indices = append(indices, byte(value>>16))
	indices = append(indices, byte(value>>24))
	return indices
}

func AppendF32(indices []uint8, value float32) []uint8 {
	v := gomath.Float32bits(value)
	indices = append(indices, byte(v))
	indices = append(indices, byte(v>>8))
	indices = append(indices, byte(v>>16))
	indices = append(indices, byte(v>>24))
	return indices
}

// MeshRender ...
func MeshRender(id MeshId) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var mesh int32 = id.mesh
	gl.BindVertexArray(g_meshes[mesh].vao)
	utils.PanicIf(GlError(), "gl.BindVertexArrayCore")

	var indexFormat uint32 = g_meshes[mesh].indexFormat
	var indexCount int32 = g_meshes[mesh].indexCount

	var cullFace uint32 = g_meshes[mesh].cullFace
	if cullFace == gl.NONE {
		DisableCulling()
	} else {
		EnableCulling(g_meshes[mesh].frontFace, cullFace)
	}

	gl.DrawElements(g_meshes[mesh].primitive, indexCount, indexFormat, unsafe.Pointer(nil))
	utils.PanicIf(GlError(), "gl.DrawArrays")
}

// MeshAppendOrthoLine ...
// TODO : pixel coordinates and clipping shouldn't be done in this function
func MeshAppendOrthoLine(id MeshId, line math.V4, color math.V4, clip math.V4, depth float32) {
	if clip.Z > 0.0 && clip.W > 0.0 {
		utils.PanicIfNot(IsValidMesh(id), "invalid id")
		var mesh int32 = id.mesh

		var clipLeft float32 = clip.X
		var clipRight float32 = clip.X + clip.Z
		var clipBottom float32 = clip.Y
		var clipTop float32 = clip.Y + clip.W

		var x0 float32 = line.X
		var y0 float32 = line.Y
		var x1 float32 = line.Z
		var y1 float32 = line.W

		// TODO REMOVE var cRight bool = x0 > clipRight && x1 > clipRight // ##0 AABB clipping...
		var cLeft bool = x0 < clipLeft && x1 < clipLeft
		var cBottom bool = y0 < clipBottom && y1 < clipBottom
		var cTop bool = y0 > clipTop && y1 > clipTop
		if cLeft == false && cBottom == false && //cRight == false &&
			cTop == false {

			var r float32 = color.X
			var g float32 = color.Y
			var b float32 = color.Z
			var a float32 = color.W

			if x0 < clipLeft {
				x0 = clipLeft
			} else if x0 > clipRight {
				x0 = clipRight
			}

			if x1 < clipLeft {
				x1 = clipLeft
			} else if x1 > clipRight {
				x1 = clipRight
			}

			if y0 < clipBottom {
				y0 = clipBottom
			} else if y0 > clipTop {
				y0 = clipTop
			}

			if y1 < clipBottom {
				y1 = clipBottom
			} else if y1 > clipTop {
				y1 = clipTop
			}

			x0 = 2.0*x0/gfx_width - 1.0
			y0 = 2.0*y0/gfx_height - 1.0
			x1 = 2.0*x1/gfx_width - 1.0
			y1 = 2.0*y1/gfx_height - 1.0

			var indices []uint8 = g_meshes[mesh].indices
			var vertices []uint8 = g_meshes[mesh].vertices
			var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride

			var indexFormat uint32 = g_meshes[mesh].indexFormat

			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+1))
			} else if indexFormat == gl.UNSIGNED_INT {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+1))
			} else {
				utils.PanicIf(true, "invalid index format")
			}

			vertices = AppendF32(vertices, x0)
			vertices = AppendF32(vertices, y0)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, 0.0)
			vertices = AppendF32(vertices, 0.0)

			vertices = AppendF32(vertices, x1)
			vertices = AppendF32(vertices, y1)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, 0.0)
			vertices = AppendF32(vertices, 0.0)

			g_meshes[mesh].vertices = vertices
			g_meshes[mesh].indices = indices
		}
	}
}

// MeshAppendOrtoRect ...
func MeshAppendOrthoRect(id MeshId, rect math.V4, left math.V4, right math.V4, bottom math.V4, top math.V4, clip math.V4, depth float32) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var x0 float32 = rect.X
	var y0 float32 = rect.Y
	var x1 float32 = x0 + rect.Z
	var y1 float32 = y0 + rect.W

	// TODO : use DrawElements : same code as MeshAppendOrthoQuad
	MeshAppendOrthoLine(id, v4.Make(x0, y0, x0, y1), left, clip, depth)
	MeshAppendOrthoLine(id, v4.Make(x1, y0, x1, y1), right, clip, depth)
	MeshAppendOrthoLine(id, v4.Make(x0, y0, x1, y0), bottom, clip, depth)
	MeshAppendOrthoLine(id, v4.Make(x0, y1, x1, y1), top, clip, depth)
}

// MeshAppendOrthoQuad ...
// TODO : pixel coordinate and clipping shouldn't be done in this function
// RECATOR :
func MeshAppendOrthoQuad(id MeshId, rect math.V4, uv math.V4, color math.V4, clip math.V4, depth float32) {
	if clip.Z > 0.0 && clip.W > 0.0 {
		utils.PanicIfNot(IsValidMesh(id), "invalid id")

		var clipLeft float32 = clip.X
		var clipRight float32 = clip.X + clip.Z
		var clipBottom float32 = clip.Y
		var clipTop float32 = clip.Y + clip.W

		var mesh int32 = id.mesh

		var primitive uint32 = g_meshes[mesh].primitive
		utils.PanicIfNot(primitive == gl.TRIANGLES, "(*mesh).primitive == gl.TRIANGLES")

		var x float32 = rect.X
		var y float32 = rect.Y
		var w float32 = rect.Z
		var h float32 = rect.W

		var x1 float32 = x + w
		var y1 float32 = y + h

		var cLeft bool = x < clipLeft && x1 < clipLeft
		var cBottom bool = y < clipBottom && y1 < clipBottom
		var cRight bool = x > clipRight && x1 > clipRight // ##0 AABB clipping...
		var cTop bool = y > clipTop && y1 > clipTop

		if cLeft == false && cBottom == false && cRight == false && cTop == false {
			var r float32 = color.X
			var g float32 = color.Y
			var b float32 = color.Z
			var a float32 = color.W

			var u0 float32 = uv.X
			var v0 float32 = uv.Y
			var u1 float32 = uv.Z
			var V1 float32 = uv.W

			var deltaLeft float32 = math.Max_f32(0.0, clipLeft-x)
			var deltaRight float32 = math.Min_f32(0.0, clipRight-x1)

			var deltaBottom float32 = math.Max_f32(0.0, clipBottom-y)
			var deltaTop float32 = math.Min_f32(0.0, clipTop-y1)

			x = math.Max_f32(x, clipLeft)
			y = math.Max_f32(y, clipBottom)

			x1 = math.Min_f32(x1, clipRight)
			y1 = math.Min_f32(y1, clipTop)

			// var nw float32 = x1 - x
			// var nh float32 = y1 - y

			var du float32 = u1 - u0
			var dv float32 = v0 - V1

			// TODO : inverted uv, fix texture loading ??
			u0 = u0 + du*deltaLeft/w
			u1 = u1 + du*deltaRight/w
			V1 = V1 + dv*deltaBottom/h
			v0 = v0 + dv*deltaTop/h

			x = 2.0*float32(int32(x))/gfx_width - 1.0
			y = 2.0*float32(int32(y))/gfx_height - 1.0
			x1 = 2.0*float32(int32(x1))/gfx_width - 1.0
			y1 = 2.0*float32(int32(y1))/gfx_height - 1.0

			var vertices []uint8 = g_meshes[mesh].vertices
			var indices []uint8 = g_meshes[mesh].indices

			var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride // TODO : uint64 offset
			var indexFormat uint32 = g_meshes[mesh].indexFormat
			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+3))
			} else if indexFormat == gl.UNSIGNED_INT {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+3))
			} else {
				utils.PanicIf(true, "invalid vertex format")
			}

			vertices = AppendF32(vertices, x)
			vertices = AppendF32(vertices, y)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, u0)
			vertices = AppendF32(vertices, V1)

			vertices = AppendF32(vertices, x1)
			vertices = AppendF32(vertices, y)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, u1)
			vertices = AppendF32(vertices, V1)

			vertices = AppendF32(vertices, x)
			vertices = AppendF32(vertices, y1)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, u0)
			vertices = AppendF32(vertices, v0)

			vertices = AppendF32(vertices, x1)
			vertices = AppendF32(vertices, y1)
			vertices = AppendF32(vertices, depth)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, u1)
			vertices = AppendF32(vertices, v0)

			g_meshes[mesh].vertices = vertices
			g_meshes[mesh].indices = indices
		}
	}
}

func appendV9(vertices []uint8, x float32, y float32, z float32, r float32, g float32, b float32, a float32, u float32, v float32) (out []uint8) {
	out = vertices
	out = AppendF32(out, x)
	out = AppendF32(out, y)
	out = AppendF32(out, z)
	out = AppendF32(out, r)
	out = AppendF32(out, g)
	out = AppendF32(out, b)
	out = AppendF32(out, a)
	out = AppendF32(out, u)
	out = AppendF32(out, v)
	return
}

// // MeshAppendParticle ...
// func MeshAppendParticle(id MeshId, particle Particle) {
// 	utils.PanicIfNot(IsValidMesh(id), "invalid id")

// 	var mesh int32 = id.mesh

// 	var primitive int32 = g_meshes[mesh].primitive
// 	utils.PanicIfNot(primitive == gl.TRIANGLES, "invalid primitive")

// 	var indices []uint8 = g_meshes[mesh].indices
// 	var vertices []uint8 = g_meshes[mesh].vertices
// 	var offset int32 = len(vertices) / g_meshes[mesh].vertexByteStride
// 	var indexFormat int32 = g_meshes[mesh].indexFormat
// 	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT, "invalid index format")

// 	var index uint16 = uint16(offset)
// 	indices = AppendUI16(indices, index)
// 	indices = AppendUI16(indices, index + 1)
// 	indices = AppendUI16(indices, index + 2)
// 	indices = AppendUI16(indices, index + 2)
// 	indices = AppendUI16(indices, index + 1)
// 	indices = AppendUI16(indices, index + 3)

// 	var px float32 = particle.position.X
// 	var py float32 = particle.position.Y
// 	var pz float32 = particle.position.Z

// 	var vx float32 = particle.velocity.X
// 	var vy float32 = particle.velocity.Y
// 	var vz float32 = particle.velocity.Z

// 	var ox float32 = particle.orientation.X
// 	var oy float32 = particle.orientation.Y
// 	var oz float32 = particle.orientation.Z
// 	var ow float32 = particle.orientation.W

// 	var ovx float32 = particle.angularVelocity.X
// 	var ovy float32 = particle.angularVelocity.Y
// 	var ovz float32 = particle.angularVelocity.Z
// 	var ovw float32 = particle.angularVelocity.W

// 	var sx float32 = particle.scale.X
// 	var sy float32 = particle.scale.Y
// 	var sz float32 = particle.scale.Z

// 	var svx float32 = particle.scaleVelocity.X
// 	var svy float32 = particle.scaleVelocity.Y
// 	var svz float32 = particle.scaleVelocity.Z

// 	var r float32 = particle.color.X
// 	var g float32 = particle.color.Y
// 	var b float32 = particle.color.Z
// 	var a float32 = particle.color.W

// 	vertices = AppendF32(vertices, px)
// 	vertices = AppendF32(vertices, py)
// 	vertices = AppendF32(vertices, pz)
// 	vertices = AppendF32(vertices, r)
// 	vertices = AppendF32(vertices, g)
// 	vertices = AppendF32(vertices, b)
// 	vertices = AppendF32(vertices, a)
// 	vertices = AppendF32(vertices, 1.0)
// 	vertices = AppendF32(vertices, 1.0)
// 	vertices = AppendF32(vertices, vx)
// 	vertices = AppendF32(vertices, vy)
// 	vertices = AppendF32(vertices, vz)
// 	vertices = AppendF32(vertices, ox)
// 	vertices = AppendF32(vertices, oy)
// 	vertices = AppendF32(vertices, oz)
// 	vertices = AppendF32(vertices, ow)
// 	vertices = AppendF32(vertices, ovx)
// 	vertices = AppendF32(vertices, ovy)
// 	vertices = AppendF32(vertices, ovz)
// 	vertices = AppendF32(vertices, ovw)
// 	vertices = AppendF32(vertices, sx)
// 	vertices = AppendF32(vertices, sy)
// 	vertices = AppendF32(vertices, sz)
// 	vertices = AppendF32(vertices, svx)
// 	vertices = AppendF32(vertices, svy)
// 	vertices = AppendF32(vertices, svz)
// 	vertices = AppendF32(vertices, particle.time)
// 	vertices = AppendF32(vertices, particle.life)
// 	vertices = AppendF32(vertices, particle.fadeIn)
// 	vertices = AppendF32(vertices, particle.fadeOut)

// 	vertices = AppendF32(vertices, px)
// 	vertices = AppendF32(vertices, py)
// 	vertices = AppendF32(vertices, pz)
// 	vertices = AppendF32(vertices, r)
// 	vertices = AppendF32(vertices, g)
// 	vertices = AppendF32(vertices, b)
// 	vertices = AppendF32(vertices, a)
// 	vertices = AppendF32(vertices, 0.0)
// 	vertices = AppendF32(vertices, 1.0)
// 	vertices = AppendF32(vertices, vx)
// 	vertices = AppendF32(vertices, vy)
// 	vertices = AppendF32(vertices, vz)
// 	vertices = AppendF32(vertices, ox)
// 	vertices = AppendF32(vertices, oy)
// 	vertices = AppendF32(vertices, oz)
// 	vertices = AppendF32(vertices, ow)
// 	vertices = AppendF32(vertices, ovx)
// 	vertices = AppendF32(vertices, ovy)
// 	vertices = AppendF32(vertices, ovz)
// 	vertices = AppendF32(vertices, ovw)
// 	vertices = AppendF32(vertices, sx)
// 	vertices = AppendF32(vertices, sy)
// 	vertices = AppendF32(vertices, sz)
// 	vertices = AppendF32(vertices, svx)
// 	vertices = AppendF32(vertices, svy)
// 	vertices = AppendF32(vertices, svz)
// 	vertices = AppendF32(vertices, particle.time)
// 	vertices = AppendF32(vertices, particle.life)
// 	vertices = AppendF32(vertices, particle.fadeIn)
// 	vertices = AppendF32(vertices, particle.fadeOut)

// 	vertices = AppendF32(vertices, px)
// 	vertices = AppendF32(vertices, py)
// 	vertices = AppendF32(vertices, pz)
// 	vertices = AppendF32(vertices, r)
// 	vertices = AppendF32(vertices, g)
// 	vertices = AppendF32(vertices, b)
// 	vertices = AppendF32(vertices, a)
// 	vertices = AppendF32(vertices, 1.0)
// 	vertices = AppendF32(vertices, 0.0)
// 	vertices = AppendF32(vertices, vx)
// 	vertices = AppendF32(vertices, vy)
// 	vertices = AppendF32(vertices, vz)
// 	vertices = AppendF32(vertices, ox)
// 	vertices = AppendF32(vertices, oy)
// 	vertices = AppendF32(vertices, oz)
// 	vertices = AppendF32(vertices, ow)
// 	vertices = AppendF32(vertices, ovx)
// 	vertices = AppendF32(vertices, ovy)
// 	vertices = AppendF32(vertices, ovz)
// 	vertices = AppendF32(vertices, ovw)
// 	vertices = AppendF32(vertices, sx)
// 	vertices = AppendF32(vertices, sy)
// 	vertices = AppendF32(vertices, sz)
// 	vertices = AppendF32(vertices, svx)
// 	vertices = AppendF32(vertices, svy)
// 	vertices = AppendF32(vertices, svz)
// 	vertices = AppendF32(vertices, particle.time)
// 	vertices = AppendF32(vertices, particle.life)
// 	vertices = AppendF32(vertices, particle.fadeIn)
// 	vertices = AppendF32(vertices, particle.fadeOut)

// 	vertices = AppendF32(vertices, px)
// 	vertices = AppendF32(vertices, py)
// 	vertices = AppendF32(vertices, pz)
// 	vertices = AppendF32(vertices, r)
// 	vertices = AppendF32(vertices, g)
// 	vertices = AppendF32(vertices, b)
// 	vertices = AppendF32(vertices, a)
// 	vertices = AppendF32(vertices, 0.0)
// 	vertices = AppendF32(vertices, 0.0)
// 	vertices = AppendF32(vertices, vx)
// 	vertices = AppendF32(vertices, vy)
// 	vertices = AppendF32(vertices, vz)
// 	vertices = AppendF32(vertices, ox)
// 	vertices = AppendF32(vertices, oy)
// 	vertices = AppendF32(vertices, oz)
// 	vertices = AppendF32(vertices, ow)
// 	vertices = AppendF32(vertices, ovx)
// 	vertices = AppendF32(vertices, ovy)
// 	vertices = AppendF32(vertices, ovz)
// 	vertices = AppendF32(vertices, ovw)
// 	vertices = AppendF32(vertices, sx)
// 	vertices = AppendF32(vertices, sy)
// 	vertices = AppendF32(vertices, sz)
// 	vertices = AppendF32(vertices, svx)
// 	vertices = AppendF32(vertices, svy)
// 	vertices = AppendF32(vertices, svz)
// 	vertices = AppendF32(vertices, particle.time)
// 	vertices = AppendF32(vertices, particle.life)
// 	vertices = AppendF32(vertices, particle.fadeIn)
// 	vertices = AppendF32(vertices, particle.fadeOut)

// 	g_meshes[mesh].vertices = vertices
// 	g_meshes[mesh].indices = indices

// }

// MeshAppendTriangle ...
// TODO remove wire branching
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO remove ccw branching
// TODO support arbitrary vertex layout
// TODO implement cw
// TODO strips
func MeshAppendTriangle(id MeshId, wire bool, ccw bool, p0 math.V3, p1 math.V3, p2 math.V3, color math.V4) { // TODO : use line strip
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot((wire == true && primitive == gl.LINES) || primitive == gl.TRIANGLES, "invalid primitive")

	var indices []uint8 = g_meshes[mesh].indices
	var vertices []uint8 = g_meshes[mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	if wire {
		if indexFormat == gl.UNSIGNED_SHORT {
			indices = AppendUI16(indices, uint16(offset+0))
			indices = AppendUI16(indices, uint16(offset+1))
			indices = AppendUI16(indices, uint16(offset+1))
			indices = AppendUI16(indices, uint16(offset+2))
			indices = AppendUI16(indices, uint16(offset+2))
			indices = AppendUI16(indices, uint16(offset+0))
			offset = offset + 3
		} else {
			indices = AppendUI32(indices, uint32(offset+0))
			indices = AppendUI32(indices, uint32(offset+1))
			indices = AppendUI32(indices, uint32(offset+1))
			indices = AppendUI32(indices, uint32(offset+2))
			indices = AppendUI32(indices, uint32(offset+2))
			indices = AppendUI32(indices, uint32(offset+0))
			offset = offset + 3
		}
	} else {
		if ccw {
			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+2))
				offset = offset + 3
			} else {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+2))
				offset = offset + 3
			}
		} else {
			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+1))
				offset = offset + 3
			} else {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+1))
				offset = offset + 3
			}
		}
	}

	var r float32 = color.X
	var g float32 = color.Y
	var b float32 = color.Z
	var a float32 = color.W

	var u0 float32 = 1.0
	var u1 float32 = 0.0
	var v0 float32 = 1.0
	var v1 float32 = 0.0

	vertices = AppendF32(vertices, p0.X)
	vertices = AppendF32(vertices, p0.Y)
	vertices = AppendF32(vertices, p0.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u0)
	vertices = AppendF32(vertices, v0)

	vertices = AppendF32(vertices, p1.X)
	vertices = AppendF32(vertices, p1.Y)
	vertices = AppendF32(vertices, p1.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u1)
	vertices = AppendF32(vertices, v0)

	vertices = AppendF32(vertices, p2.X)
	vertices = AppendF32(vertices, p2.Y)
	vertices = AppendF32(vertices, p2.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u0)
	vertices = AppendF32(vertices, v1)

	g_meshes[mesh].vertices = vertices
	g_meshes[mesh].indices = indices
}

// MeshAppendQuad ...
// TODO remove wire branching
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO remove ccw branching
// TODO support arbitrary vertex layout
// TODO implement cw
// TODO strips
func MeshAppendQuad(id MeshId, wire bool, ccw bool, position math.V3, right math.V3, top math.V3, back math.V3, color math.V4, p math.V4) { // TODO : use line strip
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot((wire == true && primitive == gl.LINES) || primitive == gl.TRIANGLES, "invalid primitive")

	var indices []uint8 = g_meshes[mesh].indices
	var vertices []uint8 = g_meshes[mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	if wire {
		if indexFormat == gl.UNSIGNED_SHORT {
			indices = AppendUI16(indices, uint16(offset+0))
			indices = AppendUI16(indices, uint16(offset+1))
			indices = AppendUI16(indices, uint16(offset+1))
			indices = AppendUI16(indices, uint16(offset+3))
			indices = AppendUI16(indices, uint16(offset+3))
			indices = AppendUI16(indices, uint16(offset+2))
			indices = AppendUI16(indices, uint16(offset+2))
			indices = AppendUI16(indices, uint16(offset+0))
			offset = offset + 4
		} else {
			indices = AppendUI32(indices, uint32(offset+0))
			indices = AppendUI32(indices, uint32(offset+1))
			indices = AppendUI32(indices, uint32(offset+1))
			indices = AppendUI32(indices, uint32(offset+3))
			indices = AppendUI32(indices, uint32(offset+3))
			indices = AppendUI32(indices, uint32(offset+2))
			indices = AppendUI32(indices, uint32(offset+2))
			indices = AppendUI32(indices, uint32(offset+0))
			offset = offset + 4
		}
	} else {
		if ccw {
			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+3))
				offset = offset + 4
			} else {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+3))
				offset = offset + 4
			}
		} else {
			if indexFormat == gl.UNSIGNED_SHORT {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+3))
				offset = offset + 4
			} else {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+3))
				offset = offset + 4
			}
		}
	}

	var x float32 = position.X
	var y float32 = position.Y
	var z float32 = position.Z

	var rx float32 = right.X
	var ry float32 = right.Y
	var rz float32 = right.Z

	var tx float32 = top.X
	var ty float32 = top.Y
	var tz float32 = top.Z

	var bx float32 = back.X
	var by float32 = back.Y
	var bz float32 = back.Z

	var fbl math.V3 = v3.Make(x-rx-tx-bx, y-ry-ty-by, z-rz-tz-bz)
	var fbr math.V3 = v3.Make(x+rx-tx-bx, y+ry-ty-by, z+rz-tz-bz)
	var ftl math.V3 = v3.Make(x-rx+tx-bx, y-ry+ty-by, z-rz+tz-bz)
	var ftr math.V3 = v3.Make(x+rx+tx-bx, y+ry+ty-by, z+rz+tz-bz)

	var r float32 = color.X
	var g float32 = color.Y
	var b float32 = color.Z
	var a float32 = color.W

	var u0 float32 = 1.0
	var u1 float32 = 0.0
	var v0 float32 = 1.0
	var v1 float32 = 0.0

	vertices = AppendF32(vertices, fbl.X)
	vertices = AppendF32(vertices, fbl.Y)
	vertices = AppendF32(vertices, fbl.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u0)
	vertices = AppendF32(vertices, v0)
	vertices = AppendF32(vertices, p.X)
	vertices = AppendF32(vertices, p.Y)
	vertices = AppendF32(vertices, p.Z)
	vertices = AppendF32(vertices, p.W)

	vertices = AppendF32(vertices, fbr.X)
	vertices = AppendF32(vertices, fbr.Y)
	vertices = AppendF32(vertices, fbr.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u1)
	vertices = AppendF32(vertices, v0)
	vertices = AppendF32(vertices, p.X)
	vertices = AppendF32(vertices, p.Y)
	vertices = AppendF32(vertices, p.Z)
	vertices = AppendF32(vertices, p.W)

	vertices = AppendF32(vertices, ftl.X)
	vertices = AppendF32(vertices, ftl.Y)
	vertices = AppendF32(vertices, ftl.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u0)
	vertices = AppendF32(vertices, v1)
	vertices = AppendF32(vertices, p.X)
	vertices = AppendF32(vertices, p.Y)
	vertices = AppendF32(vertices, p.Z)
	vertices = AppendF32(vertices, p.W)

	vertices = AppendF32(vertices, ftr.X)
	vertices = AppendF32(vertices, ftr.Y)
	vertices = AppendF32(vertices, ftr.Z)
	vertices = AppendF32(vertices, r)
	vertices = AppendF32(vertices, g)
	vertices = AppendF32(vertices, b)
	vertices = AppendF32(vertices, a)
	vertices = AppendF32(vertices, u1)
	vertices = AppendF32(vertices, v1)
	vertices = AppendF32(vertices, p.X)
	vertices = AppendF32(vertices, p.Y)
	vertices = AppendF32(vertices, p.Z)
	vertices = AppendF32(vertices, p.W)

	g_meshes[mesh].vertices = vertices
	g_meshes[mesh].indices = indices
}

// MeshAppendEllipseGizmoUI16 ...
func MeshAppendEllipseGizmoUI16(id MeshId, transform math.M44, position math.V3, radius math.V3, steps int) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var indices []uint8 = g_meshes[id.mesh].indices
	var vertices []uint8 = g_meshes[id.mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[id.mesh].vertexByteStride
	utils.PanicIf(g_meshes[id.mesh].indexFormat != gl.UNSIGNED_SHORT, "invalid index format")

	// TODO REMOVEvar q math.V4
	var fstep float32 = math.M2PI_f32 / float32(steps-1)
	var p math.V4
	var index uint16 = uint16(offset)

	for i := 0; i < steps; i++ {
		indices = AppendUI16(indices, index)
		index = index + 1
		indices = AppendUI16(indices, index)

		var t float32 = float32(uint32(i)+offset) * fstep

		sin, cos := gomath.Sincos(float64(t))
		p.X = position.X
		p.Y = position.Y + radius.Y*float32(cos)
		p.Z = position.Z + radius.Z*float32(sin)
		p.W = 1.0
		p = v4.Transform(p, transform)

		vertices = AppendF32(vertices, p.X)
		vertices = AppendF32(vertices, p.Y)
		vertices = AppendF32(vertices, p.Z)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 0.0)
	}
	indices = indices[:len(indices)-4]
	for i := 0; i < steps; i++ {
		indices = AppendUI16(indices, index)
		index = index + 1
		indices = AppendUI16(indices, index)

		var t float32 = float32(uint32(i)+offset) * fstep

		sin, cos := gomath.Sincos(float64(t))
		p.X = position.X + radius.X*float32(cos)
		p.Y = position.Y
		p.Z = position.Z + radius.Z*float32(sin)
		p.W = 1.0
		p = v4.Transform(p, transform)

		vertices = AppendF32(vertices, p.X)
		vertices = AppendF32(vertices, p.Y)
		vertices = AppendF32(vertices, p.Z)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 0.0)
	}
	indices = indices[:len(indices)-4]

	for i := 0; i < steps; i++ {
		indices = AppendUI16(indices, index)
		index = index + 1
		indices = AppendUI16(indices, index)

		var t float32 = float32(uint32(i)+offset) * fstep

		sin, cos := gomath.Sincos(float64(t))
		p.X = position.X + radius.X*float32(cos)
		p.Y = position.Y + radius.Y*float32(sin)
		p.Z = position.Z
		p.W = 1.0
		p = v4.Transform(p, transform)

		vertices = AppendF32(vertices, p.X)
		vertices = AppendF32(vertices, p.Y)
		vertices = AppendF32(vertices, p.Z)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 1.0)
		vertices = AppendF32(vertices, 0.0)
		vertices = AppendF32(vertices, 0.0)
	}
	indices = indices[:len(indices)-4]

	g_meshes[id.mesh].vertices = vertices
	g_meshes[id.mesh].indices = indices
}

// MeshAppendLine ...
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO support arbitrary vertex layout
// TODO strips
func MeshAppendLine(id MeshId, pa math.V3, pb math.V3, color math.V4) { // TODO : use line strip
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot(primitive == gl.LINES, "invalid primitive")

	var indices []uint8 = g_meshes[mesh].indices
	var vertices []uint8 = g_meshes[mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	if indexFormat == gl.UNSIGNED_SHORT {
		indices = AppendUI16(indices, uint16(offset+0))
		indices = AppendUI16(indices, uint16(offset+1))
	} else {
		indices = AppendUI32(indices, uint32(offset+0))
		indices = AppendUI32(indices, uint32(offset+1))
	}

	vertices = AppendF32(vertices, pa.X)
	vertices = AppendF32(vertices, pa.Y)
	vertices = AppendF32(vertices, pa.Z)
	vertices = AppendF32(vertices, color.X)
	vertices = AppendF32(vertices, color.Y)
	vertices = AppendF32(vertices, color.Z)
	vertices = AppendF32(vertices, color.W)
	vertices = AppendF32(vertices, 0.0)
	vertices = AppendF32(vertices, 0.0)

	vertices = AppendF32(vertices, pb.X)
	vertices = AppendF32(vertices, pb.Y)
	vertices = AppendF32(vertices, pb.Z)
	vertices = AppendF32(vertices, color.X)
	vertices = AppendF32(vertices, color.Y)
	vertices = AppendF32(vertices, color.Z)
	vertices = AppendF32(vertices, color.W)
	vertices = AppendF32(vertices, 1.0)
	vertices = AppendF32(vertices, 1.0)

	g_meshes[mesh].vertices = vertices
	g_meshes[mesh].indices = indices
}

// MeshAppendBox ...
// TODO remove wire branching
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO remove ccw branching
// TODO support arbitrary vertex layout
// TODO implement cw
// TODO strips
func MeshAppendBox(id MeshId, wire bool, ccw bool, position math.V3, right math.V3, top math.V3, back math.V3, color math.V4) { // TODO : use line strip
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot((wire == true && primitive == gl.LINES) || primitive == gl.TRIANGLES, "invalid primitive")

	var indices []uint8 = g_meshes[mesh].indices
	var vertices []uint8 = g_meshes[mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	if wire {
		if indexFormat == gl.UNSIGNED_SHORT {
			for i := 0; i < 6; i++ {
				indices = AppendUI16(indices, uint16(offset+0))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+1))
				indices = AppendUI16(indices, uint16(offset+3))
				indices = AppendUI16(indices, uint16(offset+3))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+2))
				indices = AppendUI16(indices, uint16(offset+0))
				offset = offset + 4
			}
		} else {
			for i := 0; i < 6; i++ {
				indices = AppendUI32(indices, uint32(offset+0))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+1))
				indices = AppendUI32(indices, uint32(offset+3))
				indices = AppendUI32(indices, uint32(offset+3))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+2))
				indices = AppendUI32(indices, uint32(offset+0))
				offset = offset + 4
			}
		}
	} else {
		if ccw {
			if indexFormat == gl.UNSIGNED_SHORT {
				for i := 0; i < 6; i++ {
					indices = AppendUI16(indices, uint16(offset+0))
					indices = AppendUI16(indices, uint16(offset+1))
					indices = AppendUI16(indices, uint16(offset+2))
					indices = AppendUI16(indices, uint16(offset+2))
					indices = AppendUI16(indices, uint16(offset+1))
					indices = AppendUI16(indices, uint16(offset+3))
					offset = offset + 4
				}
			} else {
				for i := 0; i < 6; i++ {
					indices = AppendUI32(indices, uint32(offset+0))
					indices = AppendUI32(indices, uint32(offset+1))
					indices = AppendUI32(indices, uint32(offset+2))
					indices = AppendUI32(indices, uint32(offset+2))
					indices = AppendUI32(indices, uint32(offset+1))
					indices = AppendUI32(indices, uint32(offset+3))
					offset = offset + 4
				}
			}
		} else {
			if indexFormat == gl.UNSIGNED_SHORT {
				for i := 0; i < 6; i++ {
					indices = AppendUI16(indices, uint16(offset+0))
					indices = AppendUI16(indices, uint16(offset+2))
					indices = AppendUI16(indices, uint16(offset+1))
					indices = AppendUI16(indices, uint16(offset+1))
					indices = AppendUI16(indices, uint16(offset+2))
					indices = AppendUI16(indices, uint16(offset+3))
					offset = offset + 4
				}
			} else {
				for i := 0; i < 6; i++ {
					indices = AppendUI32(indices, uint32(offset+0))
					indices = AppendUI32(indices, uint32(offset+2))
					indices = AppendUI32(indices, uint32(offset+1))
					indices = AppendUI32(indices, uint32(offset+1))
					indices = AppendUI32(indices, uint32(offset+2))
					indices = AppendUI32(indices, uint32(offset+3))
					offset = offset + 4
				}
			}
		}
	}

	var x float32 = position.X
	var y float32 = position.Y
	var z float32 = position.Z

	var rx float32 = right.X
	var ry float32 = right.Y
	var rz float32 = right.Z

	var tx float32 = top.X
	var ty float32 = top.Y
	var tz float32 = top.Z

	var bx float32 = back.X
	var by float32 = back.Y
	var bz float32 = back.Z

	var fbl math.V3 = v3.Make(x-rx-tx-bx, y-ry-ty-by, z-rz-tz-bz)
	var fbr math.V3 = v3.Make(x+rx-tx-bx, y+ry-ty-by, z+rz-tz-bz)
	var ftl math.V3 = v3.Make(x-rx+tx-bx, y-ry+ty-by, z-rz+tz-bz)
	var ftr math.V3 = v3.Make(x+rx+tx-bx, y+ry+ty-by, z+rz+tz-bz)
	var bbl math.V3 = v3.Make(x-rx-tx+bx, y-ry-ty+by, z-rz-tz+bz)
	var bbr math.V3 = v3.Make(x+rx-tx+bx, y+ry-ty+by, z+rz-tz+bz)
	var btl math.V3 = v3.Make(x-rx+tx+bx, y-ry+ty+by, z-rz+tz+bz)
	var btr math.V3 = v3.Make(x+rx+tx+bx, y+ry+ty+by, z+rz+tz+bz)

	var r float32 = color.X
	var g float32 = color.Y
	var b float32 = color.Z
	var a float32 = color.W

	var u0 float32 = 1.0
	var u1 float32 = 0.0
	var v0 float32 = 1.0
	var v1 float32 = 0.0

	vertices = appendV9(vertices, fbl.X, fbl.Y, fbl.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, fbr.X, fbr.Y, fbr.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, ftl.X, ftl.Y, ftl.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, ftr.X, ftr.Y, ftr.Z, r, g, b, a, u1, v1)

	vertices = appendV9(vertices, fbr.X, fbr.Y, fbr.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, bbr.X, bbr.Y, bbr.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, ftr.X, ftr.Y, ftr.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, btr.X, btr.Y, btr.Z, r, g, b, a, u1, v1)

	vertices = appendV9(vertices, bbr.X, bbr.Y, bbr.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, bbl.X, bbl.Y, bbl.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, btr.X, btr.Y, btr.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, btl.X, btl.Y, btl.Z, r, g, b, a, u1, v1)

	vertices = appendV9(vertices, bbl.X, bbl.Y, bbl.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, fbl.X, fbl.Y, fbl.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, btl.X, btl.Y, btl.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, ftl.X, ftl.Y, ftl.Z, r, g, b, a, u1, v1)

	vertices = appendV9(vertices, ftl.X, ftl.Y, ftl.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, ftr.X, ftr.Y, ftr.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, btl.X, btl.Y, btl.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, btr.X, btr.Y, btr.Z, r, g, b, a, u1, v1)

	vertices = appendV9(vertices, fbr.X, fbr.Y, fbr.Z, r, g, b, a, u0, v0)
	vertices = appendV9(vertices, fbl.X, fbl.Y, fbl.Z, r, g, b, a, u1, v0)
	vertices = appendV9(vertices, bbr.X, bbr.Y, bbr.Z, r, g, b, a, u0, v1)
	vertices = appendV9(vertices, bbl.X, bbl.Y, bbl.Z, r, g, b, a, u1, v1)

	g_meshes[mesh].vertices = vertices
	g_meshes[mesh].indices = indices
}

// MeshAppendEllipse ...
// TODO remove wire branching
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO remove ccw branching
// TODO support arbitrary vertex layout
// TODO implement cw
// TODO strips
func MeshAppendEllipse(id MeshId, wire bool, ccw bool, transform math.M44, resx uint32, resy uint32, uv math.V4, color math.V4) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	utils.PanicIf(ccw == false, "cw not implemented")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot((wire == true && primitive == gl.LINES) || primitive == gl.TRIANGLES, "invalid primitive")

	var vertices []uint8 = g_meshes[mesh].vertices
	var indices []uint8 = g_meshes[mesh].indices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	var k1 uint32
	var k2 uint32
	for x := uint32(0); x < resx; x++ {
		k1 = offset + x*(resy+1)
		k2 = k1 + resy + 1
		for y := uint32(0); y < resy; y++ {
			if x != 0 {
				if wire {
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(k1))
						indices = AppendUI16(indices, uint16(k2))
						indices = AppendUI16(indices, uint16(k2))
						indices = AppendUI16(indices, uint16(k1+1))
						indices = AppendUI16(indices, uint16(k1+1))
						indices = AppendUI16(indices, uint16(k1))
					} else {
						indices = AppendUI32(indices, uint32(k1))
						indices = AppendUI32(indices, uint32(k2))
						indices = AppendUI32(indices, uint32(k2))
						indices = AppendUI32(indices, uint32(k1+1))
						indices = AppendUI32(indices, uint32(k1+1))
						indices = AppendUI32(indices, uint32(k1))
					}
				} else {
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(k1))
						indices = AppendUI16(indices, uint16(k1+1))
						indices = AppendUI16(indices, uint16(k2))
					} else {
						indices = AppendUI32(indices, uint32(k1))
						indices = AppendUI32(indices, uint32(k1+1))
						indices = AppendUI32(indices, uint32(k2))
					}
				}
			}

			if x != (x - 1) {
				if wire {
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(k1+1))
						indices = AppendUI16(indices, uint16(k2))
						indices = AppendUI16(indices, uint16(k2))
						indices = AppendUI16(indices, uint16(k2+1))
						indices = AppendUI16(indices, uint16(k2+1))
						indices = AppendUI16(indices, uint16(k1+1))
					} else {
						indices = AppendUI32(indices, uint32(k1+1))
						indices = AppendUI32(indices, uint32(k2))
						indices = AppendUI32(indices, uint32(k2))
						indices = AppendUI32(indices, uint32(k2+1))
						indices = AppendUI32(indices, uint32(k2+1))
						indices = AppendUI32(indices, uint32(k1+1))
					}
				} else {
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(k1+1))
						indices = AppendUI16(indices, uint16(k2+1))
						indices = AppendUI16(indices, uint16(k2))
					} else {
						indices = AppendUI32(indices, uint32(k1+1))
						indices = AppendUI32(indices, uint32(k2+1))
						indices = AppendUI32(indices, uint32(k2))
					}
				}
			}
			k1++
			k2++
		}
	}

	var ystep float32 = 2.0 * math.PI_f32 / float32(resy)
	var xstep float32 = math.PI_f32 / float32(resx)

	var r float32 = color.X
	var g float32 = color.Y
	var b float32 = color.Z
	var a float32 = color.W

	for x := uint32(0); x <= resx; x++ {
		var xangle float32 = math.PI_f32/2.0 - float32(x)*xstep
		sinx, cosx := gomath.Sincos(float64(xangle))
		var posxy float32 = float32(cosx)
		var posz float32 = float32(sinx)
		var v0 float32 = float32(x) / float32(resx)

		for y := uint32(0); y <= resy; y++ {
			var yangle float32 = float32(y) * ystep
			siny, cosy := gomath.Sincos(float64(yangle))
			var posx float32 = posxy * float32(cosy)
			var posy float32 = posxy * float32(siny)
			var u0 float32 = float32(y) / float32(resy)

			var position math.V4 = v4.Transform(v4.Make(posx, posy, posz, 1.0), transform)
			vertices = appendV9(vertices, position.X, position.Y, position.Z, r, g, b, a, u0, v0)
		}
	}
	g_meshes[id.mesh].vertices = vertices
	g_meshes[id.mesh].indices = indices
}

// MeshAppendPlane ...
// TODO remove wire branching
// TODO remove short indices branching
// TODO remove debug colors branching
// TODO remove ccw branching
// TODO support arbitrary vertex layout
// TODO strips
func MeshAppendPlane(id MeshId, wire bool, ccw bool, position math.V3, right math.V3, top math.V3, back math.V3, resx uint32, resy uint32, uv math.V4, color math.V4) {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")

	var mesh int32 = id.mesh

	var primitive uint32 = g_meshes[mesh].primitive
	utils.PanicIfNot((wire == true && primitive == gl.LINES) || primitive == gl.TRIANGLES, "invalid primitive")

	var r float32 = color.X
	var g float32 = color.Y
	var b float32 = color.Z
	var a float32 = color.W

	var fw float32 = float32(resx) - 1.0
	var fh float32 = float32(resy) - 1.0

	var indices []uint8 = g_meshes[mesh].indices
	var vertices []uint8 = g_meshes[mesh].vertices
	var offset uint32 = uint32(len(vertices)) / g_meshes[mesh].vertexByteStride
	var indexFormat uint32 = g_meshes[mesh].indexFormat
	utils.PanicIf(indexFormat != gl.UNSIGNED_SHORT && indexFormat != gl.UNSIGNED_INT, "invalid index format")

	for x := uint32(0); x < resx; x++ {
		var fx float32 = float32(x)
		var dx float32 = fx/fw - 0.5
		//var dx1 float32 = (fx + 1.0) / fw
		var dw math.V3 = v3.Mulf(right, dx)
		//var dw1 math.V3 = v3.Mulf(right, dx1)
		var s0 float32 = uv.X + uv.Z*dx
		// var s1 float32 = uv.X + uv.Z*dx1
		// var s2 float32 = uv.X + uv.Z*dx
		// var s3 float32 = uv.X + uv.Z*dx1

		for y := uint32(0); y < resy; y++ {
			var fy float32 = float32(y)
			var dy float32 = fy/fh - 0.5
			//var dy1 float32 = (fy + 1.0) / fh
			var dh math.V3 = v3.Mulf(back, dy)
			//var dh1 math.V3 = v3.Mulf(back, dy1)

			var v0 math.V3 = v3.Add(position, v3.Add(dw, dh))
			// var v1 math.V3 = v3.Add(position, v3.Add(dw1, dh))
			// var v2 math.V3 = v3.Add(position, v3.Add(dw, dh1))
			// var v3 math.V3 = v3.Add(position, v3.Add(dw1, dh1))

			var t0 float32 = uv.Y + uv.W*dy
			// var t1 float32 = uv.Y + uv.W*dy
			// var t2 float32 = uv.Y + uv.W*dy1
			// var t3 float32 = uv.Y + uv.W*dy1 // ISSUE #247 : can redeclare v3 with another type

			if x < resx-1 &&
				y < resy-1 {
				var i0 uint32 = offset + y*resy + x
				var i1 uint32 = i0 + 1
				var i2 uint32 = i0 + resy
				var i3 uint32 = i2 + 1
				if wire {
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(i0))
						indices = AppendUI16(indices, uint16(i1))
						indices = AppendUI16(indices, uint16(i1))
						indices = AppendUI16(indices, uint16(i3))
						indices = AppendUI16(indices, uint16(i3))
						indices = AppendUI16(indices, uint16(i2))
						indices = AppendUI16(indices, uint16(i2))
						indices = AppendUI16(indices, uint16(i0))
						indices = AppendUI16(indices, uint16(i1))
						indices = AppendUI16(indices, uint16(i2))
					} else {
						indices = AppendUI32(indices, uint32(i0))
						indices = AppendUI32(indices, uint32(i1))
						indices = AppendUI32(indices, uint32(i1))
						indices = AppendUI32(indices, uint32(i3))
						indices = AppendUI32(indices, uint32(i3))
						indices = AppendUI32(indices, uint32(i2))
						indices = AppendUI32(indices, uint32(i2))
						indices = AppendUI32(indices, uint32(i0))
						indices = AppendUI32(indices, uint32(i1))
						indices = AppendUI32(indices, uint32(i2))
					}
				} else {
					if ccw == false {
						i1 = i0 + resy
						i2 = i0 + 1
					}
					if indexFormat == gl.UNSIGNED_SHORT {
						indices = AppendUI16(indices, uint16(i0))
						indices = AppendUI16(indices, uint16(i1))
						indices = AppendUI16(indices, uint16(i2))
						indices = AppendUI16(indices, uint16(i2))
						indices = AppendUI16(indices, uint16(i1))
						indices = AppendUI16(indices, uint16(i3))
					} else {
						indices = AppendUI32(indices, uint32(i0))
						indices = AppendUI32(indices, uint32(i1))
						indices = AppendUI32(indices, uint32(i2))
						indices = AppendUI32(indices, uint32(i2))
						indices = AppendUI32(indices, uint32(i1))
						indices = AppendUI32(indices, uint32(i3))
					}
				}
			}

			vertices = AppendF32(vertices, v0.X)
			vertices = AppendF32(vertices, v0.Y)
			vertices = AppendF32(vertices, v0.Z)
			vertices = AppendF32(vertices, r)
			vertices = AppendF32(vertices, g)
			vertices = AppendF32(vertices, b)
			vertices = AppendF32(vertices, a)
			vertices = AppendF32(vertices, s0)
			vertices = AppendF32(vertices, t0)
		}
	}
	g_meshes[mesh].vertices = vertices
	g_meshes[mesh].indices = indices
}
