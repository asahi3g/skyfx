package gfx

import (
	"fmt"
	"skyfx/gfx/gltf"
	"skyfx/math"
	"skyfx/math/intersect"
	m44 "skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"
	"skyfx/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var OCTREE_GRAPHICS int32 = 1
var OCTREE_COLLISIONS int32 = 2

// Globals ...
var g_octrees []Octree
var g_octreeCells []math.V3

// OctreeCell ...
type OctreeCell struct {
	index                    int32
	level                    int32
	center                   math.V3
	min                      math.V3
	max                      math.V3
	size                     math.V3
	transparents             []int32
	opaques                  []int32
	transparentTriangleCount int32
	opaqueTriangleCount      int32
	triangleCount            int32
	transparentPrimitives    []int32
	opaquePrimitives         []int32
	positions                []float32
	normals                  []float32
	planesOrigin             []math.V3
	planesNormal             []math.V3
	points                   []math.V3
	children                 [8]int32
}

// OctreeId ...
type OctreeId struct {
	octree int32
}

// InvalidOctree ...
func InvalidOctree() (out OctreeId) {
	out.octree = -1
	return
}

// OctreeIsValid ...
func OctreeIsValid(id OctreeId) (out bool) {
	out = id.octree >= 0 && id.octree < int32(len(g_octrees))
	return
}

// Octree ...
type Octree struct {
	// runtime
	debugMesh          MeshId
	visibles           []int32
	visibleCount       int32
	offsets            []int32
	tmp                []int32
	visiblePrimitives  []int32
	opaqueBuckets      []PrimitiveBucket
	transparentBuckets []PrimitiveBucket

	// serialized
	center     math.V3
	size       math.V3
	min        math.V3
	max        math.V3
	maxLevel   int32
	world      math.M44
	transforms []math.M44
	cells      []OctreeCell
	mins       []math.V3
	maxs       []math.V3
	primitives []Primitive
}

func octreeCreate(debug bool) (out OctreeId) {
	var octree Octree
	octree.maxLevel = -1
	if debug {
		octree.debugMesh = MeshCreate(gl.LINES, gl.UNSIGNED_SHORT, 10*8192*3, VertexLayout, 10*8192*3)
	} else {
		octree.debugMesh = InvalidMesh()
	}
	out.octree = int32(len(g_octrees))
	g_octrees = append(g_octrees, octree)
	return
}

var scratchF32 math.Vector_f32
var scratchUI8 math.Vector_ui8

func readI64(handle int32) (out int64) {
	var success bool
	out, success = utils.ReadI64(handle)
	utils.PanicIf(success == false, "failed to read int64")
	return
}

func readUI64(handle int32) (out uint64) {
	var success bool
	out, success = utils.ReadUI64(handle)
	utils.PanicIf(success == false, "failed to read uint64")
	return
}

func readF32(handle int32) (out float32) {
	var success bool
	out, success = utils.ReadF32(handle)
	utils.PanicIf(success == false, "failed to read float32")
	return
}

func readI32(handle int32) (out int32) {
	var success bool
	out, success = utils.ReadI32(handle)
	utils.PanicIf(success == false, "failed to read int32")
	return
}

func readBOOL(handle int32) (out bool) {
	var success bool
	out, success = utils.ReadBOOL(handle)
	utils.PanicIf(success == false, "failed to read bool")
	return
}

func readV3F(handle int32, out *math.V3) {
	success := math.ReadV3(handle, out)
	utils.PanicIf(success == false, "failed to read math.V3")
	return
}

func readV4F(handle int32, out *math.V4) {
	success := math.ReadV4(handle, out)
	utils.PanicIf(success == false, "failed to read math.V3")
	return
}

func readM44F(handle int32, out *math.M44) {
	success := math.ReadM44(handle, out)
	utils.PanicIf(success == false, "failed to read math.V3")
	return
}

func readSTR(handle int32) (out string) {
	var count uint64 = readUI64(handle)
	if count > 0 {
		success := utils.ReadUI8Slice(handle, scratchUI8.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read string")
		out = string(scratchUI8.Get())
	}
	return
}

func readF32Slice(handle int32, vector *math.Vector_f32) {
	var success bool
	var count uint64 = readUI64(handle)
	if count > 0 {
		success = utils.ReadF32Slice(handle, vector.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read []float32")
	}
	return
}

func readI32Slice(handle int32, vector *math.Vector_i32) {
	var success bool
	var count uint64 = readUI64(handle)
	if count > 0 {
		success = utils.ReadI32Slice(handle, vector.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read []int32")
	}
	return
}

func readUI8Slice(handle int32, vector *math.Vector_ui8) {
	var success bool
	var count uint64 = readUI64(handle)
	if count > 0 {
		success = utils.ReadUI8Slice(handle, vector.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read []uint8")
	}
	return
}

func readV3FSlice(handle int32, vector *math.Vector_v3) {
	var success bool
	var count uint64 = readUI64(handle)
	if count > 0 {
		success = math.ReadV3Slice(handle, vector.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read []math.v3")
	}
	return
}

func readM44FSlice(handle int32, vector *math.Vector_m44) {
	var success bool
	var count uint64 = readUI64(handle)
	if count > 0 {
		success = math.ReadM44Slice(handle, vector.Resize(count, false), count)
		utils.PanicIf(success == false, "failed to read []math.m44")
	}
	return
}

func writeI64(handle int32, value int64) {
	var success bool = utils.WriteI64(handle, value)
	utils.PanicIf(success == false, "failed to write int64")
}

func writeUI64(handle int32, value uint64) {
	var success bool = utils.WriteUI64(handle, value)
	utils.PanicIf(success == false, "failed to write uint64")
}

func writeF32(handle int32, value float32) {
	var success bool = utils.WriteF32(handle, value)
	utils.PanicIf(success == false, "failed to write float32")
}

func writeI32(handle int32, value int32) {
	var success bool = utils.WriteI32(handle, value)
	utils.PanicIf(success == false, "failed to write int32")
}

func writeBOOL(handle int32, value bool) {
	var success bool = utils.WriteBool(handle, value)
	utils.PanicIf(success == false, "failed to write bool")
}

func writeV3F(handle int32, value math.V3) {
	var success bool = math.WriteV3(handle, &value)
	utils.PanicIf(success == false, "failed to write v3f")
}

func writeV4F(handle int32, value math.V4) {
	var success bool = math.WriteV4(handle, &value)
	utils.PanicIf(success == false, "failed to write v4f")
}

func writeM44F(handle int32, value math.M44) {
	var success bool = math.WriteM44(handle, &value)
	utils.PanicIf(success == false, "failed to write m44f")
}

func writeSTR(handle int32, value string) {
	var count uint64 = uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		var success bool = utils.WriteUI8Slice(handle, []uint8(value))
		utils.PanicIf(success == false, "failed to write []float32")
	}
	return
}

func writeF32Slice(handle int32, value []float32) {
	var count uint64 = uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		var success bool = utils.WriteF32Slice(handle, value)
		utils.PanicIf(success == false, "failed to write []float32")
	}
}

func writeI32Slice(handle int32, value []int32) {
	var count uint64 = uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		var success bool = utils.WriteI32Slice(handle, value)
		utils.PanicIf(success == false, "failed to write []int32")
	}
}

func writeUI8Slice(handle int32, value []uint8) {
	var count uint64 = uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		var success bool = utils.WriteUI8Slice(handle, value)
		utils.PanicIf(success == false, "failed to write []uint8")
	}
}

func writeV3FSlice(handle int32, value []math.V3) {
	count := uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		success := math.WriteV3Slice(handle, value)
		utils.PanicIf(success == false, "failed to write []math.v3")
	}
}

func writeM44FSlice(handle int32, value []math.M44) {
	count := uint64(len(value))
	writeUI64(handle, count)
	if count > 0 {
		success := math.WriteM44Slice(handle, value)
		utils.PanicIf(success == false, "failed to write []math.m44")
	}
}

func readOctreeCells(handle int32, value []OctreeCell) (out []OctreeCell) {
	out = value
	var count int32 = int32(readUI64(handle))
	out = out[:0]
	var i int32
	for i = 0; i < count; i++ {
		var oc OctreeCell
		oc.index = readI32(handle)
		oc.level = readI32(handle)
		readV3F(handle, &oc.center)
		readV3F(handle, &oc.min)
		readV3F(handle, &oc.max)
		readV3F(handle, &oc.size)

		{
			var v math.Vector_i32
			readI32Slice(handle, &v)
			oc.transparents = v.Get()
		}
		{
			var v math.Vector_i32
			readI32Slice(handle, &v)
			oc.opaques = v.Get()
		}
		oc.transparentTriangleCount = readI32(handle)
		oc.opaqueTriangleCount = readI32(handle)
		oc.triangleCount = readI32(handle)
		{
			var v math.Vector_i32
			readI32Slice(handle, &v)
			oc.transparentPrimitives = v.Get()
		}
		{
			var v math.Vector_i32
			readI32Slice(handle, &v)
			oc.opaquePrimitives = v.Get()
		}
		{
			var v math.Vector_f32
			readF32Slice(handle, &v)
			oc.positions = v.Get()
		}
		{
			var v math.Vector_f32
			readF32Slice(handle, &v)
			oc.normals = v.Get()
		}
		{
			var v math.Vector_v3
			readV3FSlice(handle, &v)
			oc.planesOrigin = v.Get()
		}
		{
			var v math.Vector_v3
			readV3FSlice(handle, &v)
			oc.planesNormal = v.Get()
		}
		{
			var v math.Vector_v3
			readV3FSlice(handle, &v)
			oc.points = v.Get()
		}
		for c := 0; c < 8; c++ {
			oc.children[c] = readI32(handle)
		}
		out = append(out, oc)
	}
	return
}

func writeOctreeCells(handle int32, value []OctreeCell, options int32) {
	var count int = len(value)
	writeUI64(handle, uint64(count))
	for i := 0; i < count; i++ {
		writeI32(handle, value[i].index)
		writeI32(handle, value[i].level)
		writeV3F(handle, value[i].center)
		writeV3F(handle, value[i].min)
		writeV3F(handle, value[i].max)
		writeV3F(handle, value[i].size)
		writeI32Slice(handle, value[i].transparents)
		writeI32Slice(handle, value[i].opaques)
		writeI32(handle, value[i].transparentTriangleCount)
		writeI32(handle, value[i].opaqueTriangleCount)
		writeI32(handle, value[i].triangleCount)
		writeI32Slice(handle, value[i].transparentPrimitives)
		writeI32Slice(handle, value[i].opaquePrimitives)

		if options == OCTREE_COLLISIONS {
			writeF32Slice(handle, value[i].positions)
			writeF32Slice(handle, value[i].normals)
			writeV3FSlice(handle, value[i].planesOrigin)
			writeV3FSlice(handle, value[i].planesNormal)
			writeV3FSlice(handle, value[i].points)
		} else {
			writeUI64(handle, 0)
			writeUI64(handle, 0)
			writeUI64(handle, 0)
			writeUI64(handle, 0)
			writeUI64(handle, 0)
		}
		for c := 0; c < 8; c++ {
			writeI32(handle, value[i].children[c])
		}
	}
}

func readEffect(handle int32) (out EffectId) {
	var key int64 = readI64(handle)
	out = TemplateInstanceFromKey(g_tfxPbr, key)
	return
}

func writeEffect(handle int32, effect EffectId) {
	writeI64(handle, EffectGetKey(effect))
}

func readGLTFMaterial(handle int32) (out gltf.Material) {
	readV4F(handle, &out.PbrMetallicRoughness.BaseColorFactor)
	out.PbrMetallicRoughness.BaseColorTexture.Scale = readF32(handle)
	out.PbrMetallicRoughness.MetallicFactor = readF32(handle)
	out.PbrMetallicRoughness.RoughnessFactor = readF32(handle)
	readV4F(handle, &out.PbrSpecularGlossiness.DiffuseFactor)
	out.PbrSpecularGlossiness.DiffuseTexture.Scale = readF32(handle)
	out.PbrSpecularGlossiness.GlossinessFactor = readF32(handle)
	readV3F(handle, &out.PbrSpecularGlossiness.SpecularFactor)
	out.PbrSpecularGlossiness.SpecularGlossinessTexture.Scale = readF32(handle)
	out.NormalTexture.Scale = readF32(handle)
	out.EmissiveTexture.Scale = readF32(handle)
	readV4F(handle, &out.EmissiveFactor)
	out.AlphaMode = readI32(handle)
	out.DoubleSided = readI32(handle)
	return
}

func writeGLTFMaterial(handle int32, material gltf.Material) {
	writeV4F(handle, material.PbrMetallicRoughness.BaseColorFactor)
	writeF32(handle, material.PbrMetallicRoughness.BaseColorTexture.Scale)
	writeF32(handle, material.PbrMetallicRoughness.MetallicFactor)
	writeF32(handle, material.PbrMetallicRoughness.RoughnessFactor)
	writeV4F(handle, material.PbrSpecularGlossiness.DiffuseFactor)
	writeF32(handle, material.PbrSpecularGlossiness.DiffuseTexture.Scale)
	writeF32(handle, material.PbrSpecularGlossiness.GlossinessFactor)
	writeV3F(handle, material.PbrSpecularGlossiness.SpecularFactor)
	writeF32(handle, material.PbrSpecularGlossiness.SpecularGlossinessTexture.Scale)
	writeF32(handle, material.NormalTexture.Scale)
	writeF32(handle, material.EmissiveTexture.Scale)
	writeV4F(handle, material.EmissiveFactor)
	writeI32(handle, material.AlphaMode)
	writeI32(handle, material.DoubleSided)
}

func readTexture(handle int32) (out TextureId) {
	out = InvalidTexture()
	var texturePath string = readSTR(handle)
	if len(texturePath) > 0 {
		out = TextureCreate(texturePath, FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)
	}
	return
}

func writeTexture(handle int32, texture TextureId) {
	var texturePath string
	if IsValidTexture(texture) {
		texturePath = TextureGetPath(texture)
	}
	writeSTR(handle, texturePath)
}

func readPrimitives(handle int32, value []Primitive) (out []Primitive) {
	out = value
	var count int32 = int32(readUI64(handle))
	out = out[:0]
	var i int32
	for i = int32(0); i < count; i++ {
		var primitive Primitive
		primitive.effect = readEffect(handle)
		primitive.baseTexture = readTexture(handle)
		primitive.metalRoughTexture = readTexture(handle)
		primitive.emissiveTexture = readTexture(handle)
		primitive.normalTexture = readTexture(handle)
		primitive.occlusionTexture = readTexture(handle)

		primitive.gltfMaterial = readGLTFMaterial(handle)

		readV3F(handle, &primitive.min)
		readV3F(handle, &primitive.max)

		{
			var v math.Vector_ui8
			readUI8Slice(handle, &v)
			primitive.vertices = v.Get()
		}

		{
			var v math.Vector_ui8
			readUI8Slice(handle, &v)
			primitive.indices = v.Get()
		}

		var attributes []VertexAttribute = primitive.attributes
		var attributeCount int32 = int32(readUI64(handle))
		for attributeIndex := int32(0); attributeIndex < attributeCount; attributeIndex++ {
			var attribute VertexAttribute
			attribute.componentCount = uint32(readI32(handle))
			attribute.componentType = uint32(readI32(handle))
			attribute.componentByteSize = uint32(readI32(handle))
			attribute.componentOffset = uint32(readI32(handle))
			attribute.byteOffset = uint32(readI32(handle))
			attribute.binding = uint32(readI32(handle))
			attributes = append(attributes, attribute)
		}
		primitive.attributes = attributes

		primitive.indexByteStride = readI32(handle)
		primitive.vertexByteStride = uint32(readI32(handle))
		primitive.useSkin = readI32(handle)
		primitive.usePosition = readBOOL(handle)
		primitive.useNormal = readBOOL(handle)
		primitive.useColor = readBOOL(handle)
		primitive.useTexcoord = readBOOL(handle)
		primitive.useTangent = readBOOL(handle)
		primitive.useWeight = readBOOL(handle)
		primitive.useJoint = readBOOL(handle)

		var indexType int32 = int32(readUI64(handle))
		var primitiveType int32 = int32(readUI64(handle))
		var mesh MeshId = InvalidMesh()
		var vertexCount int = len(primitive.vertices)
		var indexCount int = len(primitive.indices)
		if (vertexCount > 0 && indexCount > 0) &&
			((indexType == gl.UNSIGNED_SHORT) || (indexType == gl.UNSIGNED_INT)) &&
			((primitiveType == gl.TRIANGLES) || (primitiveType == gl.LINES)) {
			mesh = MeshCreate(uint32(primitiveType), uint32(indexType),
				int32(len(primitive.indices))/primitive.indexByteStride,
				attributes, uint32(len(primitive.vertices))/uint32(primitive.vertexByteStride))

			MeshBegin(mesh)
			g_meshes[mesh.mesh].vertices = primitive.vertices
			g_meshes[mesh.mesh].indices = primitive.indices
			g_meshes[mesh.mesh].usePosition = primitive.usePosition
			g_meshes[mesh.mesh].useNormal = primitive.useNormal
			g_meshes[mesh.mesh].useColor = primitive.useColor
			g_meshes[mesh.mesh].useTexcoord = primitive.useTexcoord
			g_meshes[mesh.mesh].useTangent = primitive.useTangent
			g_meshes[mesh.mesh].useWeight = primitive.useWeight
			g_meshes[mesh.mesh].useJoint = primitive.useJoint
			g_meshes[mesh.mesh].min = primitive.min
			g_meshes[mesh.mesh].max = primitive.max
			MeshEnd(mesh)
		}

		primitive.mesh = mesh

		out = append(out, primitive)
	}
	return
}

func writePrimitives(handle int32, value []Primitive, options int32) {
	var count int = len(value)
	writeUI64(handle, uint64(count))
	for i := 0; i < count; i++ {
		writeEffect(handle, value[i].effect)
		writeTexture(handle, value[i].baseTexture)
		writeTexture(handle, value[i].metalRoughTexture)
		writeTexture(handle, value[i].emissiveTexture)
		writeTexture(handle, value[i].normalTexture)
		writeTexture(handle, value[i].occlusionTexture)

		writeGLTFMaterial(handle, value[i].gltfMaterial)

		writeV3F(handle, value[i].min)
		writeV3F(handle, value[i].max)

		if options == OCTREE_GRAPHICS {
			writeUI8Slice(handle, value[i].vertices)
			writeUI8Slice(handle, value[i].indices)
		} else {
			writeUI64(handle, 0)
			writeUI64(handle, 0)
		}

		var attributes []VertexAttribute = value[i].attributes
		var attributeCount int = len(attributes)
		writeUI64(handle, uint64(attributeCount))
		for attributeIndex := 0; attributeIndex < attributeCount; attributeIndex++ {
			writeI32(handle, int32(attributes[attributeIndex].componentCount))
			writeI32(handle, int32(attributes[attributeIndex].componentType))
			writeI32(handle, int32(attributes[attributeIndex].componentByteSize))
			writeI32(handle, int32(attributes[attributeIndex].componentOffset))
			writeI32(handle, int32(attributes[attributeIndex].byteOffset))
			writeI32(handle, int32(attributes[attributeIndex].binding))
		}

		writeI32(handle, value[i].indexByteStride)
		writeI32(handle, int32(value[i].vertexByteStride))
		writeI32(handle, value[i].useSkin)
		writeBOOL(handle, value[i].usePosition)
		writeBOOL(handle, value[i].useNormal)
		writeBOOL(handle, value[i].useColor)
		writeBOOL(handle, value[i].useTexcoord)
		writeBOOL(handle, value[i].useTangent)
		writeBOOL(handle, value[i].useWeight)
		writeBOOL(handle, value[i].useJoint)

		var indexType uint64 = 0
		var primitiveType uint64 = 0
		var mesh MeshId = value[i].mesh
		if IsValidMesh(mesh) {
			indexType = uint64(g_meshes[mesh.mesh].indexFormat)
			primitiveType = uint64(g_meshes[mesh.mesh].primitive)
		}
		writeUI64(handle, indexType)
		writeUI64(handle, primitiveType)
	}
	//mesh MeshId
}

func octreeLoad(id OctreeId) {
	var octreePrimitives []Primitive = g_octrees[id.octree].primitives
	var octreePrimitiveCount int = len(octreePrimitives)
	for primitiveIndex := 0; primitiveIndex < octreePrimitiveCount; primitiveIndex++ {
		var hash0 uint64
		var hash1 uint64
		hash0, hash1 = PrimitiveComputeHash(octreePrimitives[primitiveIndex])
		octreePrimitives[primitiveIndex].hash0 = hash0
		octreePrimitives[primitiveIndex].hash1 = hash1
	}
}

var fingerPrint uint64 = 383838383838

func OctreeLoad(path string, debug bool) (out OctreeId) {
	out = octreeCreate(debug)
	var octreeIndex int32 = out.octree
	var handle int32 = utils.Open(path)
	if handle != -1 {
		readV3F(handle, &g_octrees[octreeIndex].center)
		readV3F(handle, &g_octrees[octreeIndex].size)
		readV3F(handle, &g_octrees[octreeIndex].min)
		readV3F(handle, &g_octrees[octreeIndex].max)
		g_octrees[octreeIndex].maxLevel = readI32(handle)
		readM44F(handle, &g_octrees[octreeIndex].world)
		{
			var v math.Vector_m44
			readM44FSlice(handle, &v)
			g_octrees[octreeIndex].transforms = v.Get()
		}
		g_octrees[octreeIndex].cells = readOctreeCells(handle, g_octrees[octreeIndex].cells)
		{
			var v math.Vector_v3
			readV3FSlice(handle, &v)
			g_octrees[octreeIndex].mins = v.Get()
		}
		{
			var v math.Vector_v3
			readV3FSlice(handle, &v)
			g_octrees[octreeIndex].maxs = v.Get()
		}
		g_octrees[octreeIndex].primitives = readPrimitives(handle, g_octrees[octreeIndex].primitives)
		var fp uint64 = readUI64(handle)
		utils.PanicIf(fp != fingerPrint, "failed to OctreeLoad : invalid fingerPrint")
		var success bool = utils.Close(handle)
		if success == false {
			out = InvalidOctree()
		}
	}

	octreeLoad(out)
	return
}

func OctreeSave(id OctreeId, path string, options int32) (out bool) {
	var handle int32 = utils.Create(path)
	if handle != -1 {
		writeV3F(handle, g_octrees[id.octree].center)
		writeV3F(handle, g_octrees[id.octree].size)
		writeV3F(handle, g_octrees[id.octree].min)
		writeV3F(handle, g_octrees[id.octree].max)
		writeI32(handle, g_octrees[id.octree].maxLevel)
		writeM44F(handle, g_octrees[id.octree].world)
		writeM44FSlice(handle, g_octrees[id.octree].transforms)
		writeOctreeCells(handle, g_octrees[id.octree].cells, options)
		writeV3FSlice(handle, g_octrees[id.octree].mins)
		writeV3FSlice(handle, g_octrees[id.octree].maxs)
		writePrimitives(handle, g_octrees[id.octree].primitives, options)
		writeUI64(handle, fingerPrint)
		out = utils.Close(handle)
	}
	return
}

func neqUI8Slice(left []uint8, right []uint8) (out bool) {
	out = true
	var leftCount int = len(left)
	var rightCount int = len(right)

	if leftCount == rightCount {
		out = false
		for i := 0; i < leftCount; i++ {
			if left[i] != right[i] {
				i = leftCount
				out = true
			}
		}
	}
	return
}

func neqI32Slice(left []int32, right []int32) (out bool) {
	out = true
	var leftCount int = len(left)
	var rightCount int = len(right)

	if leftCount == rightCount {
		out = false
		for i := 0; i < leftCount; i++ {
			if left[i] != right[i] {
				i = leftCount
				out = true
			}
		}
	}
	return
}

func neqF32Slice(left []float32, right []float32) (out bool) {
	out = true
	var leftCount int = len(left)
	var rightCount int = len(right)

	if leftCount == rightCount {
		out = false
		for i := 0; i < leftCount; i++ {
			if left[i] != right[i] {
				out = true
				i = leftCount
			}
		}
	}
	return
}

func neqV3FSlice(left []math.V3, right []math.V3) (out bool) {
	out = true
	var leftCount int = len(left)
	var rightCount int = len(right)

	if leftCount == rightCount {
		out = false
		for i := 0; i < leftCount; i++ {
			if v3.Nequ(left[i], right[i]) {
				out = true
				i = leftCount
			}
		}
	}
	return
}

func OctreeAssertEquals(left OctreeId, right OctreeId, options int32) {
	utils.PanicIf(v3.Nequ(g_octrees[left.octree].center, g_octrees[right.octree].center), "center")
	utils.PanicIf(v3.Nequ(g_octrees[left.octree].size, g_octrees[right.octree].size), "size")
	utils.PanicIf(v3.Nequ(g_octrees[left.octree].min, g_octrees[right.octree].min), "min")
	utils.PanicIf(v3.Nequ(g_octrees[left.octree].max, g_octrees[right.octree].max), "max")
	utils.PanicIf((g_octrees[left.octree].maxLevel != g_octrees[right.octree].maxLevel), "maxLevel")
	utils.PanicIf(m44.Nequ(g_octrees[left.octree].world, g_octrees[right.octree].world), "world")

	var leftTX []math.M44 = g_octrees[left.octree].transforms
	var rightTX []math.M44 = g_octrees[right.octree].transforms

	var leftTxCount int = len(leftTX)
	var rightTxCount int = len(rightTX)

	utils.PanicIf(leftTxCount != rightTxCount, fmt.Sprintf("transform count : %d vs %d", leftTxCount, rightTxCount))
	for i := 0; i < leftTxCount; i++ {
		utils.PanicIf(m44.Nequ(leftTX[i], rightTX[i]), "transforms")
	}

	var leftCells []OctreeCell = g_octrees[left.octree].cells
	var rightCells []OctreeCell = g_octrees[right.octree].cells

	var leftCellCount int = len(leftCells)
	var rightCellCount int = len(rightCells)

	utils.PanicIf(leftCellCount != rightCellCount, fmt.Sprintf("cell count : %d vs %d", leftCellCount, rightCellCount))
	for i := 0; i < leftCellCount; i++ {
		var lc OctreeCell = leftCells[i]
		var rc OctreeCell = rightCells[i]

		utils.PanicIf(lc.index != rc.index, "index")
		utils.PanicIf(lc.level != rc.level, "level")
		utils.PanicIf(v3.Nequ(lc.center, rc.center), "center")
		utils.PanicIf(v3.Nequ(lc.min, rc.min), "min")
		utils.PanicIf(v3.Nequ(lc.max, rc.max), "max")
		utils.PanicIf(v3.Nequ(lc.size, rc.size), "size")

		utils.PanicIf(neqI32Slice(lc.transparents, rc.transparents), "transparents")
		utils.PanicIf(neqI32Slice(lc.opaques, rc.opaques), "opaques")
		utils.PanicIf(lc.transparentTriangleCount != rc.transparentTriangleCount, "transparentTriangleCount")
		utils.PanicIf(lc.opaqueTriangleCount != rc.opaqueTriangleCount, "opaqueTriangleCount")
		utils.PanicIf(lc.triangleCount != rc.triangleCount, "triangleCount")
		utils.PanicIf(neqI32Slice(lc.transparentPrimitives, rc.transparentPrimitives), "transparentPrimitives")
		utils.PanicIf(neqI32Slice(lc.opaquePrimitives, rc.opaquePrimitives), "opaquePrimitives")

		if options == OCTREE_COLLISIONS {
			utils.PanicIf(neqF32Slice(lc.positions, rc.positions), "positions")
			utils.PanicIf(neqF32Slice(lc.normals, rc.normals), "normals")
			utils.PanicIf(neqV3FSlice(lc.planesOrigin, rc.planesOrigin), "planesOrigin")
			utils.PanicIf(neqV3FSlice(lc.planesNormal, rc.planesNormal), "planesNormal")
			utils.PanicIf(neqV3FSlice(lc.points, rc.points), "points")
		}

		for c := 0; c < 8; c++ {
			utils.PanicIf(lc.children[c] != rc.children[c], "children")
		}
	}

	var leftPrimitives []Primitive = g_octrees[left.octree].primitives
	var rightPrimitives []Primitive = g_octrees[right.octree].primitives

	var leftPrimitiveCount int = len(leftPrimitives)
	var rightPrimitiveCount int = len(rightPrimitives)

	utils.PanicIf(leftPrimitiveCount != rightPrimitiveCount, fmt.Sprintf("primitive count : %d vs %d", leftPrimitiveCount, rightPrimitiveCount))
	for i := 0; i < leftPrimitiveCount; i++ {
		var lp Primitive = leftPrimitives[i]
		var rp Primitive = rightPrimitives[i]

		utils.PanicIf(lp.effect.effect != rp.effect.effect, "effect")
		utils.PanicIf(lp.baseTexture.texture != rp.baseTexture.texture, "baseTexture")
		utils.PanicIf(lp.metalRoughTexture.texture != rp.metalRoughTexture.texture, "metalRoughTexture")
		utils.PanicIf(lp.emissiveTexture.texture != rp.emissiveTexture.texture, "emissiveTexture")
		utils.PanicIf(lp.normalTexture.texture != rp.normalTexture.texture, "normalTexture")
		utils.PanicIf(lp.occlusionTexture.texture != rp.occlusionTexture.texture, "occlusionTexture")

		var lm gltf.Material = lp.gltfMaterial
		var rm gltf.Material = rp.gltfMaterial
		utils.PanicIf(v4.Nequ(lm.PbrMetallicRoughness.BaseColorFactor, rm.PbrMetallicRoughness.BaseColorFactor), "baseColorFactor")
		utils.PanicIf(lm.PbrMetallicRoughness.BaseColorTexture.Scale != rm.PbrMetallicRoughness.BaseColorTexture.Scale, "baseColorTexture.scale")
		utils.PanicIf(lm.PbrMetallicRoughness.MetallicFactor != rm.PbrMetallicRoughness.MetallicFactor, "metallicFactor")
		utils.PanicIf(lm.PbrMetallicRoughness.RoughnessFactor != rm.PbrMetallicRoughness.RoughnessFactor, "roughnessFactor")
		utils.PanicIf(v4.Nequ(lm.PbrSpecularGlossiness.DiffuseFactor, rm.PbrSpecularGlossiness.DiffuseFactor), "diffuseFactor")
		utils.PanicIf(lm.PbrSpecularGlossiness.DiffuseTexture.Scale != rm.PbrSpecularGlossiness.DiffuseTexture.Scale, "diffuseTexture.scale")
		utils.PanicIf(lm.PbrSpecularGlossiness.GlossinessFactor != rm.PbrSpecularGlossiness.GlossinessFactor, "glossinessFactor")
		utils.PanicIf(v3.Nequ(lm.PbrSpecularGlossiness.SpecularFactor, rm.PbrSpecularGlossiness.SpecularFactor), "specularFactor")
		utils.PanicIf(lm.PbrSpecularGlossiness.SpecularGlossinessTexture.Scale != rm.PbrSpecularGlossiness.SpecularGlossinessTexture.Scale, "specularGlossinessTexture.scale")
		utils.PanicIf(lm.NormalTexture.Scale != rm.NormalTexture.Scale, "normalTexture.scale")
		utils.PanicIf(lm.EmissiveTexture.Scale != rm.EmissiveTexture.Scale, "emissiveTexture.scale")
		utils.PanicIf(v4.Nequ(lm.EmissiveFactor, rm.EmissiveFactor), "emissiveFactor")
		utils.PanicIf(lm.AlphaMode != rm.AlphaMode, "alphaMode")
		utils.PanicIf(lm.DoubleSided != rm.DoubleSided, "doubleSided")

		utils.PanicIf(v3.Nequ(lp.min, rp.min), "min")
		utils.PanicIf(v3.Nequ(lp.max, rp.max), "max")

		if options == OCTREE_GRAPHICS {
			utils.PanicIf(neqUI8Slice(lp.vertices, rp.vertices), "vertices")
			utils.PanicIf(neqUI8Slice(lp.indices, rp.indices), "indices")
		}

		var leftAttributes []VertexAttribute = lp.attributes
		var rightAttributes []VertexAttribute = rp.attributes

		var leftAttributeCount int = len(leftAttributes)
		var rightAttributeCount int = len(rightAttributes)

		utils.PanicIf(leftAttributeCount != rightAttributeCount, "attributes")
		for attributeIndex := 0; attributeIndex < leftAttributeCount; attributeIndex++ {
			utils.PanicIf(leftAttributes[attributeIndex].componentCount != rightAttributes[attributeIndex].componentCount, "componentCount")
			utils.PanicIf(leftAttributes[attributeIndex].componentType != rightAttributes[attributeIndex].componentType, "componentCount")
			utils.PanicIf(leftAttributes[attributeIndex].componentByteSize != rightAttributes[attributeIndex].componentByteSize, "componentByteSize")
			utils.PanicIf(leftAttributes[attributeIndex].componentOffset != rightAttributes[attributeIndex].componentOffset, "componentOffset")
			utils.PanicIf(leftAttributes[attributeIndex].byteOffset != rightAttributes[attributeIndex].byteOffset, "byteOffset")
			utils.PanicIf(leftAttributes[attributeIndex].binding != rightAttributes[attributeIndex].binding, "binding")
		}

		utils.PanicIf(lp.indexByteStride != rp.indexByteStride, "indexByteStride")
		utils.PanicIf(lp.vertexByteStride != rp.vertexByteStride, "vertexByteStride")

		utils.PanicIf(lp.useSkin != rp.useSkin, "useSkin")
		utils.PanicIf(lp.usePosition != rp.usePosition, "usePosition")
		utils.PanicIf(lp.useNormal != rp.useNormal, "useNormal")
		utils.PanicIf(lp.useColor != rp.useColor, "useColor")
		utils.PanicIf(lp.useTexcoord != rp.useTexcoord, "useTexcoord")
		utils.PanicIf(lp.useTangent != rp.useTangent, "useTangent")
		utils.PanicIf(lp.useWeight != rp.useWeight, "useWeight")
		utils.PanicIf(lp.useJoint != rp.useJoint, "useJoint")

		if options == OCTREE_GRAPHICS {
			var leftMesh MeshId = lp.mesh
			var rightMesh MeshId = rp.mesh
			utils.PanicIf(IsValidMesh(leftMesh) != IsValidMesh(rightMesh), "mesh")

			if IsValidMesh(leftMesh) {
				utils.PanicIf(g_meshes[leftMesh.mesh].indexFormat != g_meshes[rightMesh.mesh].indexFormat, "indexType")
				utils.PanicIf(g_meshes[leftMesh.mesh].primitive != g_meshes[rightMesh.mesh].primitive, "primitiveType")
			}
		}
	}

	utils.PanicIf(neqV3FSlice(g_octrees[left.octree].mins, g_octrees[right.octree].mins), "mins")
	utils.PanicIf(neqV3FSlice(g_octrees[left.octree].maxs, g_octrees[right.octree].maxs), "maxs")
}

func OctreeDestroy(id OctreeId) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
}

// OctreeCreate ...
func OctreeCreate(model ModelId, world math.M44, maxLevel int32, debug bool, options int32) (out OctreeId) {
	out = octreeCreate(debug)

	realTriangleCount = 0

	var min math.V3 = ModelGetMin(model)
	var max math.V3 = ModelGetMax(model)

	min = v3.Transform_point(min, world)
	max = v3.Transform_point(max, world)

	var newMin math.V3 = v3.Min(min, max)
	var newMax math.V3 = v3.Max(min, max)

	g_octrees[out.octree].size = v3.Mulf(v3.Sub(newMax, newMin), 0.5)
	g_octrees[out.octree].center = v3.Add(newMin, g_octrees[out.octree].size)
	g_octrees[out.octree].min = newMin
	g_octrees[out.octree].max = newMax
	g_octrees[out.octree].maxLevel = maxLevel
	g_octrees[out.octree].world = world

	var totalCount int32 = resetOffsets(out)
	var cells []OctreeCell
	for i := int32(0); i <= totalCount; i++ {
		var cell OctreeCell
		cells = append(cells, cell)
	}
	g_octrees[out.octree].cells = cells
	g_octreeCells = append(g_octreeCells, v3.Make(1.0, 1.0, 1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(1.0, 1.0, -1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(1.0, -1.0, 1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(1.0, -1.0, -1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(-1.0, 1.0, 1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(-1.0, 1.0, -1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(-1.0, -1.0, 1.0))
	g_octreeCells = append(g_octreeCells, v3.Make(-1.0, -1.0, -1.0))

	_ /*var root int32*/ = octreeUpdateCells(out, 0, g_octrees[out.octree].center, g_octrees[out.octree].size)

	var modelTransforms []math.M44 = g_models[model.model].transforms
	var transformCount int = len(modelTransforms)
	var transforms []math.M44
	for i := 0; i < transformCount; i++ {
		transforms = append(transforms, m44.MulISSUE(modelTransforms[i], world))
	}
	g_octrees[out.octree].transforms = transforms

	var cellCount int32 = resetOffsets(out)
	for i := int32(0); i < cellCount; i++ {
		cells[i].transparentPrimitives = cells[i].transparentPrimitives[:0]
		cells[i].opaquePrimitives = cells[i].opaquePrimitives[:0]
	}
	var primitives []Primitive = g_models[model.model].primitives

	var primitiveCount int = len(primitives)
	var mins []math.V3 = g_octrees[out.octree].mins
	var maxs []math.V3 = g_octrees[out.octree].maxs

	mins = mins[:0]
	maxs = mins[:0]
	for i := 0; i < primitiveCount; i++ {
		var nodeIndex int32 = primitives[i].nodeIndex
		var transform math.M44 = transforms[nodeIndex] // ISSUE : can't use transforms[primitives[i].nodeIndex]

		var primMin math.V3 = v3.Transform_point(primitives[i].min, transform)
		var primMax math.V3 = v3.Transform_point(primitives[i].max, transform)

		var newPrimMin math.V3 = v3.Min(primMin, primMax)
		var newPrimMax math.V3 = v3.Max(primMin, primMax)

		mins = append(mins, newPrimMin)
		maxs = append(maxs, newPrimMax)

		//mins = append(mins, v3.min(newMin, newMax))
		//maxs = append(maxs, v3.max(newMin, newMax))
	}

	g_octrees[out.octree].mins = mins
	g_octrees[out.octree].maxs = maxs

	octreeSplitModel(out, model, 0, g_models[model.model].opaqueMeshes, g_models[model.model].transparentMeshes)

	var asset gltf.AssetId = g_models[model.model].asset
	for i := int32(0); i < cellCount; i++ {
		if i >= 0 { //== 54 || i == 6 {
			cells[i].opaques = octreeSplitMeshes(out, cells[i].opaquePrimitives,
				asset, primitives, transforms, cells, i, cells[i].opaques, options)
			cells[i].transparents = octreeSplitMeshes(out, cells[i].transparentPrimitives,
				asset, primitives, transforms, cells, i, cells[i].transparents, options)
			var opcell []Primitive = g_octrees[out.octree].primitives
			var opaqueTriCount int32 = computeTriangleCount(opcell, cells[i].opaques)
			var transparentTriCount int32 = computeTriangleCount(opcell, cells[i].transparents)
			fmt.Printf("Splitting cell %d/%d... : Opaque %d + Transparent %d = %d \n", i, cellCount, opaqueTriCount, transparentTriCount, opaqueTriCount+transparentTriCount)
			var positions []float32 = cells[i].positions
			var positionCount int = len(positions)
			var triangleCount int = positionCount / (3 * 3)
			var normals []float32 = cells[i].normals
			if options == OCTREE_COLLISIONS {
				normals = normals[:0]
				for t := 0; t < triangleCount; t++ {
					var toffset int = t * 9
					var t0 math.V3 = v3.Make(positions[toffset+0], positions[toffset+1], positions[toffset+2])
					var t1 math.V3 = v3.Make(positions[toffset+3], positions[toffset+4], positions[toffset+5])
					var t2 math.V3 = v3.Make(positions[toffset+6], positions[toffset+7], positions[toffset+8])

					var t01 math.V3 = v3.Sub(t1, t0)
					var t02 math.V3 = v3.Sub(t2, t0)
					var t012 math.V3 = v3.Cross(t01, t02)
					var len012 float32 = v3.Length(t012)
					var normal math.V3 = v3.ZERO
					if len012 > 0.0 {
						normal = v3.Divf(t012, len012)
					}

					normals = append(normals, normal.X)
					normals = append(normals, normal.Y)
					normals = append(normals, normal.Z)
				}
			}
			cells[i].normals = normals
			cells[i].triangleCount = int32(triangleCount)
			cells[i].opaqueTriangleCount = opaqueTriCount
			cells[i].transparentTriangleCount = transparentTriCount
			utils.PanicIf((triangleCount != int(opaqueTriCount+transparentTriCount)),
				fmt.Sprintf("TRI %d, SUM %d, OPAQUE %d, ALPHA %d, POS %d, MOD %d", triangleCount,
					opaqueTriCount+transparentTriCount,
					opaqueTriCount, transparentTriCount, positionCount, positionCount%9))
			utils.PanicIf(((positionCount % 9) != 0), fmt.Sprintf("POS %d, MOD %d", positionCount, positionCount%9))
		}
	}

	octreeLoad(out)
	return
}

func resetOffsets(id OctreeId) (count int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	var totalCount int32
	var layerCount int32 = 1
	var maxLevel int32 = g_octrees[id.octree].maxLevel
	var offsets []int32 = g_octrees[id.octree].offsets
	offsets = offsets[:0]
	for i := int32(0); i <= maxLevel; i++ {
		offsets = append(offsets, totalCount)
		totalCount = totalCount + layerCount
		layerCount = layerCount * 8
	}

	g_octrees[id.octree].offsets = offsets
	count = totalCount

	g_octrees[id.octree].visibles = g_octrees[id.octree].visibles[:0]
	return
}

// OctreeUpdate ...
func OctreeUpdate(id OctreeId, frustum FrustumId, targetLevel int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")

	_ /*var totalCount int32*/ = resetOffsets(id)
	octreeUpdateLevel(id, 0, frustum, true, targetLevel)

	if IsValidMesh(g_octrees[id.octree].debugMesh) {
		MeshBegin(g_octrees[id.octree].debugMesh)
		var cells []OctreeCell = g_octrees[id.octree].cells
		var visibles []int32 = g_octrees[id.octree].visibles
		var visibleCount int = len(visibles)
		for i := 0; i < visibleCount; i++ {
			var offset int32 = visibles[i]
			var center math.V3 = cells[offset].center
			var size math.V3 = v3.Mulf(cells[offset].size, 0.99)
			MeshAppendBox(g_octrees[id.octree].debugMesh, true, false,
				center, v3.Make(size.X, 0.0, 0.0), v3.Make(0.0, size.Y, 0.0), v3.Make(0.0, 0.0, size.Z), v4.RED)
		}
		/*
		   var min math.V3
		   var max math.V3
		   var mins []math.V3 = g_octrees[id.octree].mins
		   var maxs []math.V3 = g_octrees[id.octree].maxs
		   var primitiveCount int32 = len(mins)
		   for i := 0; i < primitiveCount; i++ {
		       var size math.V3 = v3.Mulf(v3.Sub(maxs[i], mins[i]), 0.5)
		       var center math.V3 = v3.Add(mins[i], size)
		       MeshAppendBox(g_octrees[id.octree].debugMesh, true, false,
		           center, v3.Make(size.X, 0.0, 0.0), v3.Make(0.0, size.Y, 0.0), v3.Make(0.0, 0.0, size.Z),
		           v4.PINK)
		   }
		*/
		MeshEnd(g_octrees[id.octree].debugMesh)
	}
}

/*var P0 fps.ProfileId = fps.InvalidProfile()
var P1 fps.ProfileId = fps.InvalidProfile()
var P2 fps.ProfileId = fps.InvalidProfile()
var P3 fps.ProfileId = fps.InvalidProfile()
var P4 fps.ProfileId = fps.InvalidProfile()
var P5 fps.ProfileId = fps.InvalidProfile()
var P6 fps.ProfileId = fps.InvalidProfile()*/

// OctreeRender ...
func OctreeRender(id OctreeId, world math.M44, view math.M44, projection math.M44,
	envSpec TextureId, envDiff TextureId, brdf TextureId, cameraPosition math.V4, exposure float32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	DisableBlending()
	DepthState(true, gl.LESS, true)
	g_octrees[id.octree].opaqueBuckets = sortBuckets(id, g_octrees[id.octree].opaqueBuckets, true)
	octreeRender(id, world, view, projection, g_octrees[id.octree].opaqueBuckets,
		envSpec, envDiff, brdf, cameraPosition, exposure)

	EnableBlending(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	DepthState(true, gl.LESS, false)
	g_octrees[id.octree].transparentBuckets = sortBuckets(id, g_octrees[id.octree].transparentBuckets, false)
	octreeRender(id, world, view, projection, g_octrees[id.octree].transparentBuckets,
		envSpec, envDiff, brdf, cameraPosition, exposure)

	if IsValidMesh(g_octrees[id.octree].debugMesh) {
		DepthState(true, gl.LESS, false)
		EffectUse(g_fxVertexColor3D)
		EffectAssignM44(g_fxVertexColor3D, UNIFORM_WORLD, world, false)
		EffectAssignM44(g_fxVertexColor3D, UNIFORM_VIEW, view, false)
		EffectAssignM44(g_fxVertexColor3D, UNIFORM_PROJECTION, projection, false)
		MeshRender(g_octrees[id.octree].debugMesh)
	}
}

func sortBuckets(id OctreeId, in []PrimitiveBucket, opaque bool) (out []PrimitiveBucket) {
	out = in

	var cells []OctreeCell = g_octrees[id.octree].cells

	var visibles []int32 = g_octrees[id.octree].visibles
	var visibleCount int = len(visibles)

	var octreePrimitives []Primitive = g_octrees[id.octree].primitives

	var visiblePrimitives []int32 = g_octrees[id.octree].visiblePrimitives
	visiblePrimitives = visiblePrimitives[:0]

	//P2 = fps.CreateStartProfile(P2, "OctreeRender : sort visibles")
	for visible := 0; visible < visibleCount; visible++ {
		var cellIndex int32 = visibles[visible]
		var cellPrimitives []int32
		if opaque {
			cellPrimitives = cells[cellIndex].opaques
		} else {
			cellPrimitives = cells[cellIndex].transparents
		}
		var primitiveCount int = len(cellPrimitives)
		for primitive := 0; primitive < primitiveCount; primitive++ {
			var octreePrimitiveIndex int32 = cellPrimitives[primitive]
			if octreePrimitiveIndex >= 0 {
				visiblePrimitives = append(visiblePrimitives, cellPrimitives[primitive])
			}
		}
	}

	out = PrimitiveSort1(octreePrimitives, visiblePrimitives, out)

	//fps.StopProfile(P2)
	return
}

func octreeRender(id OctreeId, world math.M44, view math.M44, projection math.M44, buckets []PrimitiveBucket,
	envSpec TextureId, envDiff TextureId, brdf TextureId, cameraPosition math.V4, exposure float32) {
	//P0 = fps.CreateStartProfile(P0, "OctreeRender")
	//P1 = fps.CreateStartProfile(P1, "P1")
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")

	var lastEffect int32 = -1
	var lastBaseTexture int32 = -1
	var lastNormalTexture int32 = -1
	var lastMetalRoughTexture int32 = -1
	var lastEmissiveTexture int32 = -1
	var lastOcclusionTexture int32 = -1
	var lastBaseColorFactor math.V4 = v4.Makef(-1.0)
	var lastEmissiveFactor math.V4 = v4.Makef(-1.0)
	var lastMetalRoughFactor math.V4 = v4.Makef(-1.0)
	// TODO REMOVE var lastTransform math.M44 = m44.INVALID

	var octreePrimitives []Primitive = g_octrees[id.octree].primitives
	//fps.StopProfile(P1)
	//P3 = fps.CreateStartProfile(P3, "OctreeRender : render")
	var bucketCount int = len(buckets)
	for bucket := 0; bucket < bucketCount; bucket++ {
		var sortedPrimitives []int32 = buckets[bucket].primitives
		var sortedCount int = len(sortedPrimitives)
		var effect EffectId = buckets[bucket].effect
		if lastEffect != effect.effect {
			lastEffect = effect.effect

			EffectUse(effect)
			EffectTryAssignTexture(effect, SAMPLER_ENV_SPECULAR, envSpec, g_linearClamp)
			EffectTryAssignTexture(effect, SAMPLER_ENV_DIFFUSE, envDiff, SpLinear0Clamp)
			EffectTryAssignTexture(effect, SAMPLER_BRDF, brdf, SpLinear0Clamp) // TODO : use gltf sampler

			EffectAssignM44(effect, UNIFORM_WORLD, world, false)
			EffectAssignM44(effect, UNIFORM_VIEW, view, false)
			EffectAssignM44(effect, UNIFORM_PROJECTION, projection, false)

			EffectAssignV4(effect, UNIFORM_CAMERA_POSITION, cameraPosition)
			EffectAssignV4(effect, UNIFORM_DEBUG_0, DEBUG_0)
			EffectAssignV4(effect, UNIFORM_PBR, v4.Make(float32(TextureGetMipmapCount(envSpec)), exposure, 0.0, 0.0))

			lastBaseTexture = -1
			lastNormalTexture = -1
			lastMetalRoughTexture = -1
			lastEmissiveTexture = -1
			lastOcclusionTexture = -1
			lastBaseColorFactor = v4.Makef(-1.0)
			lastEmissiveFactor = v4.Makef(-1.0)
			lastMetalRoughFactor = v4.Makef(-1.0)
		}

		var baseTexture TextureId = buckets[bucket].baseTexture
		if lastBaseTexture != baseTexture.texture {
			lastBaseTexture = baseTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_COLOR_0, baseTexture, SpLinearWrap)
		}

		var normalTexture TextureId = buckets[bucket].normalTexture
		if lastNormalTexture != normalTexture.texture {
			lastNormalTexture = normalTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_NORMAL, normalTexture, SpLinearWrap)
		}

		var metalRoughTexture TextureId = buckets[bucket].metalRoughTexture
		if lastMetalRoughTexture != metalRoughTexture.texture {
			lastMetalRoughTexture = metalRoughTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_METAL_ROUGH, metalRoughTexture, SpLinearWrap)
		}

		var emissiveTexture TextureId = buckets[bucket].emissiveTexture
		if lastEmissiveTexture != emissiveTexture.texture {
			lastEmissiveTexture = emissiveTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_EMISSIVE, emissiveTexture, SpLinearWrap)
		}

		var occlusionTexture TextureId = buckets[bucket].occlusionTexture
		if lastOcclusionTexture != occlusionTexture.texture {
			lastOcclusionTexture = occlusionTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_OCCLUSION, occlusionTexture, SpLinearWrap)
		}

		for sorted := 0; sorted < sortedCount; sorted++ {
			var primitiveIndex int32 = sortedPrimitives[sorted]
			var primitive Primitive = octreePrimitives[primitiveIndex]
			var mesh MeshId = primitive.mesh

			var metallicRoughness gltf.MetallicRoughness = primitive.gltfMaterial.PbrMetallicRoughness
			if EffectIsValidUniformLocation(effect, UNIFORM_COLOR) {
				var baseColorFactor math.V4 = metallicRoughness.BaseColorFactor
				baseColorFactor.W = baseColorFactor.W // * alpha

				if baseColorFactor.X != lastBaseColorFactor.X ||
					baseColorFactor.Y != lastBaseColorFactor.Y ||
					baseColorFactor.Z != lastBaseColorFactor.Z ||
					baseColorFactor.W != lastBaseColorFactor.W {
					lastBaseColorFactor = baseColorFactor
					EffectAssignV4(effect, UNIFORM_COLOR, baseColorFactor)
				}
			}

			if EffectIsValidUniformLocation(effect, UNIFORM_METAL_ROUGH) {
				var metallicRoughnessFactor math.V4
				metallicRoughnessFactor.X = primitive.gltfMaterial.NormalTexture.Scale
				metallicRoughnessFactor.Y = metallicRoughness.RoughnessFactor
				metallicRoughnessFactor.Z = metallicRoughness.MetallicFactor

				if metallicRoughnessFactor.X != lastMetalRoughFactor.X ||
					metallicRoughnessFactor.Y != lastMetalRoughFactor.Y ||
					metallicRoughnessFactor.Z != lastMetalRoughFactor.Z ||
					metallicRoughnessFactor.W != lastMetalRoughFactor.W {
					lastMetalRoughFactor = metallicRoughnessFactor
					EffectAssignV4(effect, UNIFORM_METAL_ROUGH, metallicRoughnessFactor)
				}
			}

			if EffectIsValidUniformLocation(effect, UNIFORM_EMISSIVE) {
				var emissiveFactor math.V4 = primitive.gltfMaterial.EmissiveFactor
				if emissiveFactor.X != lastEmissiveFactor.X ||
					emissiveFactor.Y != lastEmissiveFactor.Y ||
					emissiveFactor.Z != lastEmissiveFactor.Z ||
					emissiveFactor.W != lastEmissiveFactor.W {
					lastEmissiveFactor = emissiveFactor
					EffectAssignV4(effect, UNIFORM_EMISSIVE, emissiveFactor)
				}
			}

			//fps.StopProfile(P4)
			//P5 = fps.CreateStartProfile(P5, "mesh_render")
			MeshSetCulling(mesh, gl.CCW, gl.BACK)
			MeshRender(mesh)
			//fps.StopProfile(P5)
		}
	}
	//fps.StopProfile(P3)
	//fps.StopProfile(P0)
}

func addVertexV2(positionAttribute int32, attributeCount int32, channels []channelInfo, index int32, in []uint8, transform math.M44, cellIndex int32, cells []OctreeCell) (out []uint8) {
	out = in
	var positions []float32 = cells[cellIndex].positions
	for c := int32(0); c < attributeCount; c++ {
		var dataType int32 = channels[c].dataType
		// TODO REMOVE var dataLen int32 = channels[c].dataLen
		var count uint32 = channels[c].count
		var offset uint32 = uint32(index) * count
		if dataType == gl.FLOAT {
			var dataF32 []float32 = channels[c].dataF32
			if positionAttribute == c {
				utils.PanicIf(count != 3, "invalid position attribute")
				if count == 3 {
					var pos math.V3 = v3.Make(dataF32[offset], dataF32[offset+1], dataF32[offset+2])
					pos = v3.Transform_point(pos, transform)
					out = AppendF32(out, pos.X)
					out = AppendF32(out, pos.Y)
					out = AppendF32(out, pos.Z)
					positions = append(positions, pos.X)
					positions = append(positions, pos.Y)
					positions = append(positions, pos.Z)
				}
			} else {
				for i := uint32(0); i < count; i++ {
					out = AppendF32(out, dataF32[offset+i])
				}
			}
		} else if dataType == gl.UNSIGNED_SHORT {
			var dataUI16 []uint16 = channels[c].dataUI16
			for i := uint32(0); i < count; i++ {
				out = AppendUI16(out, dataUI16[offset+i])
			}
		} else if dataType == gl.UNSIGNED_INT {
			var dataUI32 []uint32 = channels[c].dataUI32
			for i := uint32(0); i < count; i++ {
				out = AppendUI32(out, dataUI32[offset+i])
			}
		} else {
			utils.PanicIf(true, "unhandled dataType")
		}
	}
	cells[cellIndex].positions = positions
	return
}
func addVertex(attributeCount int32, positionAttribute int32, channels []channelInfo, position math.V3, index uint32, in []uint8, cellIndex int32, cells []OctreeCell) (out []uint8) {
	out = in
	var positions []float32 = cells[cellIndex].positions
	for c := int32(0); c < attributeCount; c++ {
		if c == positionAttribute {
			out = AppendF32(out, position.X)
			out = AppendF32(out, position.Y)
			out = AppendF32(out, position.Z)
			positions = append(positions, position.X)
			positions = append(positions, position.Y)
			positions = append(positions, position.Z)
		} else {
			var dataType int32 = channels[c].dataType
			// TODO REMOVE var dataLen int32 = channels[c].dataLen
			var count uint32 = channels[c].count
			var offset uint32 = index * count
			if dataType == gl.FLOAT {
				var dataF32 []float32 = channels[c].dataF32
				for i := uint32(0); i < count; i++ {
					out = AppendF32(out, dataF32[offset+i])
				}
			} else if dataType == gl.UNSIGNED_SHORT {
				var dataUI16 []uint16 = channels[c].dataUI16
				for i := uint32(0); i < count; i++ {
					out = AppendUI16(out, dataUI16[offset+i])
				}
			} else if dataType == gl.UNSIGNED_INT {
				var dataUI32 []uint32 = channels[c].dataUI32
				for i := uint32(0); i < count; i++ {
					out = AppendUI32(out, dataUI32[offset+i])
				}
			} else {
				utils.PanicIf(true, "unhandled dataType")
			}
		}
	}
	cells[cellIndex].positions = positions
	return
}

func lerpVertex(attributeCount int32, positionAttribute int32, channels []channelInfo, posA math.V3, iA uint32, posB math.V3, iB uint32, time float32, in []uint8) (out []uint8) {
	out = in
	for c := int32(0); c < attributeCount; c++ {
		if c == positionAttribute {
			out = AppendF32(out, math.Lerp_f32(posA.X, posB.X, time))
			out = AppendF32(out, math.Lerp_f32(posA.Y, posB.Y, time))
			out = AppendF32(out, math.Lerp_f32(posA.Z, posB.Z, time))
		} else {
			var dataType int32 = channels[c].dataType
			// TODO REMOVE var dataLen int32 = channels[c].dataLen
			var count uint32 = channels[c].count
			var oA uint32 = iA * count
			var oB uint32 = iB * count
			if dataType == gl.FLOAT {
				var dataF32 []float32 = channels[c].dataF32
				for i := uint32(0); i < count; i++ {
					out = AppendF32(out, math.Lerp_f32(dataF32[oA+i], dataF32[oB+i], time))
				}
			} else if dataType == gl.UNSIGNED_SHORT {
				var dataUI16 []uint16 = channels[c].dataUI16
				for i := uint32(0); i < count; i++ {
					out = AppendUI16(out, dataUI16[oB+i])
				}
			} else if dataType == gl.UNSIGNED_INT {
				var dataUI32 []uint32 = channels[c].dataUI32
				for i := uint32(0); i < count; i++ {
					out = AppendUI32(out, dataUI32[oB+i])
				}
			} else {
				utils.PanicIf(true, "unhandled dataType")
			}
		}
	}
	return
}

func getPrimitiveIndex(octreePrimitives []Primitive, attributes []VertexAttribute, cellPrimitives []int32, parent Primitive) (out int32) {
	out = -1
	var primitiveCount int = len(cellPrimitives)
	var attributeCount int = len(attributes)

	for primitiveIndex := 0; primitiveIndex < primitiveCount; primitiveIndex++ {
		var octreePrimitiveIndex int32 = cellPrimitives[primitiveIndex]
		var primitive Primitive = octreePrimitives[octreePrimitiveIndex]
		var primitiveAttributes []VertexAttribute = primitive.attributes
		if attributeCount == len(primitiveAttributes) {
			var same bool = true
			for i := 0; i < attributeCount; i++ {
				if attributes[i].componentCount != primitiveAttributes[i].componentCount ||
					attributes[i].componentType != primitiveAttributes[i].componentType ||
					attributes[i].componentByteSize != primitiveAttributes[i].componentByteSize ||
					attributes[i].componentOffset != primitiveAttributes[i].componentOffset ||
					attributes[i].byteOffset != primitiveAttributes[i].byteOffset ||
					attributes[i].binding != primitiveAttributes[i].binding {
					same = false
					i = attributeCount
				}
			}

			if same {
				if primitive.baseTexture.texture == parent.baseTexture.texture &&
					primitive.metalRoughTexture.texture == parent.metalRoughTexture.texture &&
					primitive.normalTexture.texture == parent.normalTexture.texture &&
					primitive.emissiveTexture.texture == parent.emissiveTexture.texture &&
					primitive.occlusionTexture.texture == parent.occlusionTexture.texture &&
					v4.Equ(primitive.gltfMaterial.PbrMetallicRoughness.BaseColorFactor, parent.gltfMaterial.PbrMetallicRoughness.BaseColorFactor) &&
					primitive.gltfMaterial.PbrMetallicRoughness.MetallicFactor == parent.gltfMaterial.PbrMetallicRoughness.MetallicFactor &&
					primitive.gltfMaterial.PbrMetallicRoughness.RoughnessFactor == parent.gltfMaterial.PbrMetallicRoughness.RoughnessFactor &&
					v4.Equ(primitive.gltfMaterial.EmissiveFactor, parent.gltfMaterial.EmissiveFactor) &&
					primitive.gltfMaterial.AlphaMode == parent.gltfMaterial.AlphaMode &&
					primitive.gltfMaterial.DoubleSided == parent.gltfMaterial.DoubleSided {

					out = int32(primitiveIndex)
					return
				}
			}
		}
	}
	return
}

func channelLerpVertex(to []channelInfo, from []channelInfo, i0 uint32, i1 uint32, time float32, cellIndex int32, cells []OctreeCell, positionAttribute int32) (out []channelInfo) {
	out = to
	var attributeCount int = len(from)
	for c := 0; c < attributeCount; c++ {
		var dataType int32 = from[c].dataType
		// TODO REMOVE var dataLen int32 = from[c].dataLen
		var count uint32 = from[c].count
		out[c].count = count
		var offset0 uint32 = i0 * count
		var offset1 uint32 = i1 * count
		if dataType == gl.FLOAT {
			var dataOut []float32 = out[c].dataF32
			var dataF32 []float32 = from[c].dataF32
			for i := uint32(0); i < count; i++ {
				var f0 float32 = dataF32[offset0+i]
				var f1 float32 = dataF32[offset1+i]
				var fl float32 = math.Lerpsat_f32(f0, f1, time)

				dataOut = append(dataOut, fl)
				//dataOut = append(dataOut, dataF32[offset0 + i])
			}
			out[c].dataF32 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else if dataType == gl.UNSIGNED_SHORT {
			var dataOut []uint16 = to[c].dataUI16
			var dataUI16 []uint16 = from[c].dataUI16
			for i := uint32(0); i < count; i++ {
				dataOut = append(dataOut, dataUI16[offset0+i])
			}
			out[c].dataUI16 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else if dataType == gl.UNSIGNED_INT {
			var dataOut []uint32 = to[c].dataUI32
			var dataUI32 []uint32 = from[c].dataUI32
			for i := uint32(0); i < count; i++ {
				dataOut = append(dataOut, dataUI32[offset0+i])
			}
			out[c].dataUI32 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else {
			utils.PanicIf(true, "unhandled dataType")
		}
	}
	return
}

func channelResize(from []channelInfo) (out []channelInfo) {
	out = from
	var channelCount int = len(from)
	for i := 0; i < channelCount; i++ {
		out[i].dataLen = 0
		out[i].dataF32 = out[i].dataF32[:0]
		out[i].dataUI16 = out[i].dataUI16[:0]
		out[i].dataUI32 = out[i].dataUI32[:0]
	}
	return
}

func channelGetPosition(from []channelInfo, channelIndex int32, vertexIndex uint32) (out math.V3) {
	var count uint32 = from[channelIndex].count
	var data []float32 = from[channelIndex].dataF32
	var offset uint32 = vertexIndex * count
	out.X = data[offset]
	out.Y = data[offset+1]
	out.Z = data[offset+2]
	return
}

func channelAppendVertex(to []channelInfo, from []channelInfo, index uint32, cellIndex int32, cells []OctreeCell, positionAttribute int32) (out []channelInfo) {
	out = to
	var attributeCount int = len(from)

	for c := 0; c < attributeCount; c++ {
		var dataType int32 = from[c].dataType
		// TODO REMOVE var dataLen int32 = from[c].dataLen
		var count uint32 = from[c].count
		var offset uint32 = index * count
		out[c].count = count
		if dataType == gl.FLOAT {
			var dataOut []float32 = out[c].dataF32
			var dataF32 []float32 = from[c].dataF32
			for i := uint32(0); i < count; i++ {
				var pf float32 = dataF32[offset+i]
				dataOut = append(dataOut, pf)
			}
			out[c].dataF32 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else if dataType == gl.UNSIGNED_SHORT {
			var dataOut []uint16 = to[c].dataUI16
			var dataUI16 []uint16 = from[c].dataUI16
			for i := uint32(0); i < count; i++ {
				dataOut = append(dataOut, dataUI16[offset+i])
			}
			out[c].dataUI16 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else if dataType == gl.UNSIGNED_INT {
			var dataOut []uint32 = to[c].dataUI32
			var dataUI32 []uint32 = from[c].dataUI32
			for i := uint32(0); i < count; i++ {
				dataOut = append(dataOut, dataUI32[offset+i])
			}
			out[c].dataUI32 = dataOut
			out[c].dataLen = uint32(len(dataOut))
		} else {
			utils.PanicIf(true, "unhandled dataType")
		}
	}
	return
}

func computeTriangleCount(primitives []Primitive, cellPrimitives []int32) (out int32) {
	var primitiveCount int = len(cellPrimitives)
	for i := 0; i < primitiveCount; i++ {
		var primitiveIndex int = int(cellPrimitives[i])
		if primitiveIndex >= 0 && primitiveIndex < len(primitives) {
			var indices []uint8 = primitives[primitiveIndex].indices
			var indexByteStride int32 = primitives[primitiveIndex].indexByteStride
			out = out + int32(len(indices))/int32(indexByteStride*3)
		}
	}
	return
}

var realTriangleCount int32

func octreeSplitMeshes(id OctreeId, meshes []int32, asset gltf.AssetId, primitives []Primitive, transforms []math.M44, cells []OctreeCell, cellIndex int32, cellPrimitives []int32, options int32) (out []int32) {
	var cellMax math.V3 = cells[cellIndex].max
	var cellMin math.V3 = cells[cellIndex].min
	var cellPlanesNormals []math.V3 = cells[cellIndex].planesNormal
	var cellPlanesOrigins []math.V3 = cells[cellIndex].planesOrigin
	var cellPoints []math.V3 = cells[cellIndex].points
	var vertPos []bool
	var meshCount int = len(meshes)
	var mode uint32
	var cullFace uint32
	var octreePrimitives []Primitive = g_octrees[id.octree].primitives
	var triangleVertices []math.V3 // TODO USE SCRATCH PAD ?
	triangleVertices = append(triangleVertices, v3.ZERO)
	triangleVertices = append(triangleVertices, v3.ZERO)
	triangleVertices = append(triangleVertices, v3.ZERO)

	for meshIndex := 0; meshIndex < meshCount; meshIndex++ {
		var renderable int32 = meshes[meshIndex]
		var primitive Primitive = primitives[renderable]
		var transform math.M44 = transforms[primitive.nodeIndex]
		var channels []channelInfo = g_meshes[primitive.mesh.mesh].channels
		mode = g_meshes[primitive.mesh.mesh].primitive
		cullFace = g_meshes[primitive.mesh.mesh].cullFace

		var attributes []VertexAttribute = g_meshes[primitive.mesh.mesh].attributes
		var attributeCount int32 = int32(len(attributes))
		var primitiveIndex int32 = getPrimitiveIndex(octreePrimitives, attributes, cellPrimitives, primitive)
		if primitiveIndex == -1 {
			var cellPrimitive Primitive
			cellPrimitive.attributes = attributes
			cellPrimitive.mesh = InvalidMesh()
			cellPrimitive.min = v3.MAX
			cellPrimitive.max = v3.MIN
			cellPrimitive.baseTexture = primitive.baseTexture
			cellPrimitive.metalRoughTexture = primitive.metalRoughTexture
			cellPrimitive.normalTexture = primitive.normalTexture
			cellPrimitive.emissiveTexture = primitive.emissiveTexture
			cellPrimitive.occlusionTexture = primitive.occlusionTexture
			cellPrimitive.gltfMaterial = primitive.gltfMaterial
			cellPrimitive.gltfPrimitive = primitive.gltfPrimitive
			cellPrimitive.useSkin = primitive.useSkin
			cellPrimitive.effect = primitive.effect
			cellPrimitive.indexByteStride = g_meshes[primitive.mesh.mesh].indexByteStride
			cellPrimitive.vertexByteStride = g_meshes[primitive.mesh.mesh].vertexByteStride
			cellPrimitive.usePosition = g_meshes[primitive.mesh.mesh].usePosition
			cellPrimitive.useNormal = g_meshes[primitive.mesh.mesh].useNormal
			cellPrimitive.useColor = g_meshes[primitive.mesh.mesh].useColor
			cellPrimitive.useTexcoord = g_meshes[primitive.mesh.mesh].useTexcoord
			cellPrimitive.useTangent = g_meshes[primitive.mesh.mesh].useTangent
			cellPrimitive.useWeight = g_meshes[primitive.mesh.mesh].useWeight
			cellPrimitive.useJoint = g_meshes[primitive.mesh.mesh].useJoint
			primitiveIndex = int32(len(cellPrimitives))
			cellPrimitives = append(cellPrimitives, int32(len(octreePrimitives)))
			octreePrimitives = append(octreePrimitives, cellPrimitive)
		}

		var octreePrimitiveIndex int32 = cellPrimitives[primitiveIndex]

		// TODO : bake
		var positionAttribute int32 = -1
		for attributeIndex := int32(0); attributeIndex < attributeCount; attributeIndex++ {
			var binding uint32 = attributes[attributeIndex].binding
			if binding == ATTRIBUTE_POSITION {
				positionAttribute = int32(attributeIndex)
			}
		}

		var vertices []uint8 = octreePrimitives[octreePrimitiveIndex].vertices
		var oldVertexCount int = len(vertices)

		var indices []uint8 = octreePrimitives[octreePrimitiveIndex].indices
		var oldIndexCount int = len(indices)

		var posMin math.V3 = octreePrimitives[octreePrimitiveIndex].min
		var posMax math.V3 = octreePrimitives[octreePrimitiveIndex].max

		utils.PanicIf(positionAttribute == -1, "invalid position attribute")
		var positionData []float32 = channels[positionAttribute].dataF32
		var positionComponentCount uint32 = attributes[positionAttribute].componentCount
		// TODO REMOVE var world math.M44 = transforms[primitive.nodeIndex]
		var indicesAccessor gltf.Accessor = gltf.AssetGetAccessor(asset, primitive.gltfPrimitive.Indices)
		var indicesType int32 = indicesAccessor.ComponentType
		if indicesAccessor.ComponentCount == 1 {
			if indicesAccessor.ComponentType == gl.UNSIGNED_SHORT {
				var data []uint16 = indicesAccessor.DataUI16
				var index uint16 = uint16(len(indices) / 2)
				var indexCount int = len(data)
				utils.PanicIf((indexCount%3) != 0, "invalid index count")
				var triangleCount int = indexCount / 3
				for triangle := 0; triangle < triangleCount; triangle++ {
					var ui0 uint16 = data[triangle*3+0]
					var ui1 uint16 = data[triangle*3+1]
					var ui2 uint16 = data[triangle*3+2]
					var i0 uint32 = uint32(ui0)
					var i1 uint32 = uint32(ui1)
					var i2 uint32 = uint32(ui2)
					var pi0 uint32 = i0 * positionComponentCount
					var pi1 uint32 = i1 * positionComponentCount
					var pi2 uint32 = i2 * positionComponentCount
					var p0 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi0], positionData[pi0+1], positionData[pi0+2]), transform)
					var p1 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi1], positionData[pi1+1], positionData[pi1+2]), transform)
					var p2 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi2], positionData[pi2+1], positionData[pi2+2]), transform)

					triangleVertices[0] = p0
					triangleVertices[1] = p1
					triangleVertices[2] = p2

					var x0 float32 = p0.X
					var y0 float32 = p0.Y
					var z0 float32 = p0.Z
					var x1 float32 = p1.X
					var y1 float32 = p1.Y
					var z1 float32 = p1.Z
					var x2 float32 = p2.X
					var y2 float32 = p2.Y
					var z2 float32 = p2.Z
					var add int32
					if x0 >= cellMin.X && x0 <= cellMax.X &&
						y0 >= cellMin.Y && y0 <= cellMax.Y &&
						z0 >= cellMin.Z && z0 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x0)
						posMin.Y = math.Min_f32(posMin.Y, y0)
						posMin.Z = math.Min_f32(posMin.Z, z0)
						posMax.X = math.Max_f32(posMax.X, x0)
						posMax.Y = math.Max_f32(posMax.Y, y0)
						posMax.Z = math.Max_f32(posMax.Z, z0)
						add = add | 1
					}
					if x1 >= cellMin.X && x1 <= cellMax.X &&
						y1 >= cellMin.Y && y1 <= cellMax.Y &&
						z1 >= cellMin.Z && z1 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x1)
						posMin.Y = math.Min_f32(posMin.Y, y1)
						posMin.Z = math.Min_f32(posMin.Z, z1)
						posMax.X = math.Max_f32(posMax.X, x1)
						posMax.Y = math.Max_f32(posMax.Y, y1)
						posMax.Z = math.Max_f32(posMax.Z, z1)
						add = add | 2
					}
					if x2 >= cellMin.X && x2 <= cellMax.X &&
						y2 >= cellMin.Y && y2 <= cellMax.Y &&
						z2 >= cellMin.Z && z2 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x2)
						posMin.Y = math.Min_f32(posMin.Y, y2)
						posMin.Z = math.Min_f32(posMin.Z, z2)
						posMax.X = math.Max_f32(posMax.X, x2)
						posMax.Y = math.Max_f32(posMax.Y, y2)
						posMax.Z = math.Max_f32(posMax.Z, z2)
						add = add | 4
					}

					if add == 7 {
						if cullFace == gl.BACK || cullFace == gl.NONE {
							vertices = addVertex(attributeCount, positionAttribute, channels, p0, i0, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p1, i1, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p2, i2, vertices, cellIndex, cells)
							indices = AppendUI16(indices, index)
							index = index + 1
							indices = AppendUI16(indices, index)
							index = index + 1
							indices = AppendUI16(indices, index)
							index = index + 1
						}

						if cullFace == gl.FRONT || cullFace == gl.NONE {
							vertices = addVertex(attributeCount, positionAttribute, channels, p2, i2, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p1, i1, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p0, i0, vertices, cellIndex, cells)
							indices = AppendUI16(indices, index)
							index = index + 1
							indices = AppendUI16(indices, index)
							index = index + 1
							indices = AppendUI16(indices, index)
							index = index + 1
						}
					} else {
						if intersect.TriangleIntersectsAABB(cellMin, cellMax, cellPoints, triangleVertices) {
							var newChannels []channelInfo // TODO : remove alloc
							var oldChannels []channelInfo // TODO : remove alloc
							for a := int32(0); a < attributeCount; a++ {
								var chaninfo channelInfo
								chaninfo.count = channels[a].count
								chaninfo.dataType = channels[a].dataType
								oldChannels = append(oldChannels, chaninfo)
								newChannels = append(newChannels, chaninfo)
							}

							oldChannels = channelAppendVertex(oldChannels, channels, i0, cellIndex, cells, positionAttribute)
							oldChannels = channelAppendVertex(oldChannels, channels, i1, cellIndex, cells, positionAttribute)
							oldChannels = channelAppendVertex(oldChannels, channels, i2, cellIndex, cells, positionAttribute)

							for p := 0; p < 6; p++ {
								var upFront int32 = -1
								var downFront int32 = -1
								var planeNormal math.V3 = cellPlanesNormals[p]
								var planeOrigin math.V3 = cellPlanesOrigins[p]
								// TODO REMOVE var step int32 = -1
								var vertCount uint32 = oldChannels[0].dataLen / oldChannels[0].count
								vertPos = vertPos[:0]
								var belowCount int32
								for v := uint32(0); v < vertCount; v++ {
									var currentPos math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v), transform)
									var pdp float32 = v3.Dot(v3.Sub(currentPos, planeOrigin), planeNormal)
									var pointBelowPlane bool
									if pdp < 0.0 {
										pointBelowPlane = true
										belowCount = belowCount + 1
									}
									vertPos = append(vertPos, pointBelowPlane)
								}

								for v := uint32(0); v < vertCount; v++ {
									var nextPos uint32 = (v + 1) % vertCount
									if vertPos[v] == true && vertPos[nextPos] == false {
										upFront = int32(v)
									} else if vertPos[v] == false && vertPos[nextPos] == true {
										downFront = int32(v)
									}
								}

								newChannels = channelResize(newChannels)
								if upFront >= 0 && downFront >= 0 {
									utils.PanicIf(upFront == downFront, "invalid fronts")
									var vup math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(upFront)), transform)
									var upNext int32 = (upFront + 1) % int32(vertCount)
									var vupNext math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(upNext)), transform)
									var itup bool
									var tup float32
									itup, tup = intersect.RayIntersectsPlane(vup, vupNext, planeOrigin, planeNormal)
									if itup {
										newChannels = channelLerpVertex(newChannels, oldChannels, uint32(upFront), uint32(upNext), tup, cellIndex, cells, positionAttribute)
									} else {
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(upFront), cellIndex, cells, positionAttribute)
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(upNext), cellIndex, cells, positionAttribute)
									}

									var vdown math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(downFront)), transform)
									var downNext int32 = (downFront + 1) % int32(vertCount)
									var vdownNext math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(downNext)), transform)
									var itdown bool
									var tdown float32
									itdown, tdown = intersect.RayIntersectsPlane(vdown, vdownNext, planeOrigin, planeNormal)
									if itdown {
										newChannels = channelLerpVertex(newChannels, oldChannels, uint32(downFront), uint32(downNext), tdown, cellIndex, cells, positionAttribute)
									} else {
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(downFront), cellIndex, cells, positionAttribute)
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(downNext), cellIndex, cells, positionAttribute)
									}

									var cfront bool = true
									var nv int32 = (downFront + 1) % int32(vertCount)
									for cfront == true {
										if nv == upFront {
											cfront = false
										} else {
											newChannels = channelAppendVertex(newChannels, oldChannels, uint32(nv), cellIndex, cells, positionAttribute)
											nv = (nv + 1) % int32(vertCount)
										}
									}
									newChannels = channelAppendVertex(newChannels, oldChannels, uint32(nv), cellIndex, cells, positionAttribute)
									var tmpChannels []channelInfo = oldChannels
									oldChannels = newChannels
									newChannels = tmpChannels
								}
							}

							var vertexCount uint32 = oldChannels[0].dataLen / oldChannels[0].count
							for v := uint32(0); v < vertexCount-1; v++ {
								triangleVertices[0] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, 0), transform)
								triangleVertices[1] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v), transform)
								triangleVertices[2] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v+1), transform)

								if intersect.TriangleIntersectsAABB(cellMin, cellMax, cellPoints, triangleVertices) {
									if cullFace == gl.BACK || cullFace == gl.NONE {
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, 0, vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v+1), vertices, transform, cellIndex, cells)
										indices = AppendUI16(indices, index)
										index = index + 1
										indices = AppendUI16(indices, index)
										index = index + 1
										indices = AppendUI16(indices, index)
										index = index + 1
									}

									if cullFace == gl.FRONT || cullFace == gl.NONE {
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v+1), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, 0, vertices, transform, cellIndex, cells)
										indices = AppendUI16(indices, index)
										index = index + 1
										indices = AppendUI16(indices, index)
										index = index + 1
										indices = AppendUI16(indices, index)
										index = index + 1
									}
								} else {
								}
							}
						}
					}
					octreePrimitives[octreePrimitiveIndex].indices = indices
					octreePrimitives[octreePrimitiveIndex].vertices = vertices
					octreePrimitives[octreePrimitiveIndex].min = posMin
					octreePrimitives[octreePrimitiveIndex].max = posMax
					octreePrimitives[octreePrimitiveIndex].indexType = uint32(indicesType)
				}
			} else if indicesAccessor.ComponentType == gl.UNSIGNED_INT {
				var data []uint32 = indicesAccessor.DataUI32
				var index uint32 = uint32(len(indices) / 4)
				var indexCount int = len(data)
				utils.PanicIf((indexCount%3) != 0, "invalid index count")
				var triangleCount int = indexCount / 3
				for triangle := 0; triangle < triangleCount; triangle++ {
					var ui0 uint32 = data[triangle*3+0]
					var ui1 uint32 = data[triangle*3+1]
					var ui2 uint32 = data[triangle*3+2]
					// var i0 int32 = int32(ui0)
					// var i1 int32 = int32(ui1)
					// var i2 int32 = int32(ui2)

					var pi0 uint32 = ui0 * positionComponentCount
					var pi1 uint32 = ui1 * positionComponentCount
					var pi2 uint32 = ui2 * positionComponentCount
					var p0 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi0], positionData[pi0+1], positionData[pi0+2]), transform)
					var p1 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi1], positionData[pi1+1], positionData[pi1+2]), transform)
					var p2 math.V3 = v3.Transform_point(
						v3.Make(positionData[pi2], positionData[pi2+1], positionData[pi2+2]), transform)

					triangleVertices[0] = p0
					triangleVertices[1] = p1
					triangleVertices[2] = p2

					var x0 float32 = p0.X
					var y0 float32 = p0.Y
					var z0 float32 = p0.Z
					var x1 float32 = p1.X
					var y1 float32 = p1.Y
					var z1 float32 = p1.Z
					var x2 float32 = p2.X
					var y2 float32 = p2.Y
					var z2 float32 = p2.Z
					var add int32
					if x0 >= cellMin.X && x0 <= cellMax.X &&
						y0 >= cellMin.Y && y0 <= cellMax.Y &&
						z0 >= cellMin.Z && z0 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x0)
						posMin.Y = math.Min_f32(posMin.Y, y0)
						posMin.Z = math.Min_f32(posMin.Z, z0)
						posMax.X = math.Max_f32(posMax.X, x0)
						posMax.Y = math.Max_f32(posMax.Y, y0)
						posMax.Z = math.Max_f32(posMax.Z, z0)
						add = add | 1
					}
					if x1 >= cellMin.X && x1 <= cellMax.X &&
						y1 >= cellMin.Y && y1 <= cellMax.Y &&
						z1 >= cellMin.Z && z1 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x1)
						posMin.Y = math.Min_f32(posMin.Y, y1)
						posMin.Z = math.Min_f32(posMin.Z, z1)
						posMax.X = math.Max_f32(posMax.X, x1)
						posMax.Y = math.Max_f32(posMax.Y, y1)
						posMax.Z = math.Max_f32(posMax.Z, z1)
						add = add | 2
					}
					if x2 >= cellMin.X && x2 <= cellMax.X &&
						y2 >= cellMin.Y && y2 <= cellMax.Y &&
						z2 >= cellMin.Z && z2 <= cellMax.Z {
						posMin.X = math.Min_f32(posMin.X, x2)
						posMin.Y = math.Min_f32(posMin.Y, y2)
						posMin.Z = math.Min_f32(posMin.Z, z2)
						posMax.X = math.Max_f32(posMax.X, x2)
						posMax.Y = math.Max_f32(posMax.Y, y2)
						posMax.Z = math.Max_f32(posMax.Z, z2)
						add = add | 4
					}

					if add == 7 {
						if cullFace == gl.BACK || cullFace == gl.NONE {
							//if realTriangleCount >= minTriangle && realTriangleCount <= maxTriangle {
							vertices = addVertex(attributeCount, positionAttribute, channels, p0, ui0, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p1, ui1, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p2, ui2, vertices, cellIndex, cells)
							indices = AppendUI32(indices, index)
							index = index + 1
							indices = AppendUI32(indices, index)
							index = index + 1
							indices = AppendUI32(indices, index)
							index = index + 1
							//}
							realTriangleCount = realTriangleCount + 1
						}
						if cullFace == gl.FRONT || cullFace == gl.NONE {
							//if realTriangleCount >= minTriangle && realTriangleCount <= maxTriangle {
							vertices = addVertex(attributeCount, positionAttribute, channels, p2, ui2, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p1, ui1, vertices, cellIndex, cells)
							vertices = addVertex(attributeCount, positionAttribute, channels, p0, ui0, vertices, cellIndex, cells)
							indices = AppendUI32(indices, index)
							index = index + 1
							indices = AppendUI32(indices, index)
							index = index + 1
							indices = AppendUI32(indices, index)
							index = index + 1
							//}
							realTriangleCount = realTriangleCount + 1
						}
					} else {
						if intersect.TriangleIntersectsAABB(cellMin, cellMax, cellPoints, triangleVertices) {
							var newChannels []channelInfo // TODO : remove alloc
							var oldChannels []channelInfo // TODO : remove alloc
							for a := int32(0); a < attributeCount; a++ {
								var chaninfo channelInfo
								chaninfo.count = channels[a].count
								chaninfo.dataType = channels[a].dataType
								oldChannels = append(oldChannels, chaninfo)
								newChannels = append(newChannels, chaninfo)
							}

							oldChannels = channelAppendVertex(oldChannels, channels, ui0, cellIndex, cells, positionAttribute)
							oldChannels = channelAppendVertex(oldChannels, channels, ui1, cellIndex, cells, positionAttribute)
							oldChannels = channelAppendVertex(oldChannels, channels, ui2, cellIndex, cells, positionAttribute)

							for p := 0; p < 6; p++ {
								var upFront int32 = -1
								var downFront int32 = -1
								var planeNormal math.V3 = cellPlanesNormals[p]
								var planeOrigin math.V3 = cellPlanesOrigins[p]
								// TODO REMOVE var step int32 = -1
								var vertCount uint32 = oldChannels[0].dataLen / oldChannels[0].count
								vertPos = vertPos[:0]
								var belowCount int32
								for v := uint32(0); v < vertCount; v++ {
									var currentPos math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v), transform)
									var pdp float32 = v3.Dot(v3.Sub(currentPos, planeOrigin), planeNormal)
									var pointBelowPlane bool
									if pdp < 0.0 {
										pointBelowPlane = true
										belowCount = belowCount + 1
									}
									vertPos = append(vertPos, pointBelowPlane)
								}

								for v := uint32(0); v < vertCount; v++ {
									var nextPos uint32 = (v + 1) % vertCount
									if vertPos[v] == true && vertPos[nextPos] == false {
										upFront = int32(v)
									} else if vertPos[v] == false && vertPos[nextPos] == true {
										downFront = int32(v)
									}
								}

								newChannels = channelResize(newChannels)
								if upFront >= 0 && downFront >= 0 {
									utils.PanicIf(upFront == downFront, "invalid fronts")
									var vup math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(upFront)), transform)
									var upNext int32 = (upFront + 1) % int32(vertCount)
									var vupNext math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(upNext)), transform)
									var itup bool
									var tup float32
									itup, tup = intersect.RayIntersectsPlane(vup, vupNext, planeOrigin, planeNormal)
									if itup {
										newChannels = channelLerpVertex(newChannels, oldChannels, uint32(upFront), uint32(upNext), tup, cellIndex, cells, positionAttribute)
									} else {
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(upFront), cellIndex, cells, positionAttribute)
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(upNext), cellIndex, cells, positionAttribute)
									}

									var vdown math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(downFront)), transform)
									var downNext int32 = (downFront + 1) % int32(vertCount)
									var vdownNext math.V3 = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, uint32(downNext)), transform)
									var itdown bool
									var tdown float32
									itdown, tdown = intersect.RayIntersectsPlane(vdown, vdownNext, planeOrigin, planeNormal)
									if itdown {
										newChannels = channelLerpVertex(newChannels, oldChannels, uint32(downFront), uint32(downNext), tdown, cellIndex, cells, positionAttribute)
									} else {
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(downFront), cellIndex, cells, positionAttribute)
										newChannels = channelAppendVertex(newChannels, oldChannels, uint32(downNext), cellIndex, cells, positionAttribute)
									}

									var cfront bool = true
									var nv int32 = (downFront + 1) % int32(vertCount)
									for cfront == true {
										if nv == upFront {
											cfront = false
										} else {
											newChannels = channelAppendVertex(newChannels, oldChannels, uint32(nv), cellIndex, cells, positionAttribute)
											nv = (nv + 1) % int32(vertCount)
										}
									}
									newChannels = channelAppendVertex(newChannels, oldChannels, uint32(nv), cellIndex, cells, positionAttribute)
									var tmpChannels []channelInfo = oldChannels
									oldChannels = newChannels
									newChannels = tmpChannels
								}
							}

							var vertexCount uint32 = oldChannels[0].dataLen / oldChannels[0].count
							for v := uint32(0); v < vertexCount-1; v++ {
								triangleVertices[0] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, 0), transform)
								triangleVertices[1] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v), transform)
								triangleVertices[2] = v3.Transform_point(channelGetPosition(oldChannels, positionAttribute, v+1), transform)

								if intersect.TriangleIntersectsAABB(cellMin, cellMax, cellPoints, triangleVertices) {
									if cullFace == gl.BACK || cullFace == gl.NONE {
										//if realTriangleCount >= minTriangle && realTriangleCount <= maxTriangle {
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, 0, vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v+1), vertices, transform, cellIndex, cells)
										indices = AppendUI32(indices, index)
										index = index + 1
										indices = AppendUI32(indices, index)
										index = index + 1
										indices = AppendUI32(indices, index)
										index = index + 1
										//}
										realTriangleCount = realTriangleCount + 1
									}
									if cullFace == gl.FRONT || cullFace == gl.NONE {
										//if realTriangleCount >= minTriangle && realTriangleCount <= maxTriangle {
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v+1), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, int32(v), vertices, transform, cellIndex, cells)
										vertices = addVertexV2(positionAttribute, attributeCount, oldChannels, 0, vertices, transform, cellIndex, cells)
										indices = AppendUI32(indices, index)
										index = index + 1
										indices = AppendUI32(indices, index)
										index = index + 1
										indices = AppendUI32(indices, index)
										index = index + 1
										//}
										realTriangleCount = realTriangleCount + 1
									}
								} else {
								}
							}
						}
					}
					octreePrimitives[octreePrimitiveIndex].indices = indices
					octreePrimitives[octreePrimitiveIndex].vertices = vertices
					octreePrimitives[octreePrimitiveIndex].min = posMin
					octreePrimitives[octreePrimitiveIndex].max = posMax
					octreePrimitives[octreePrimitiveIndex].indexType = uint32(indicesType)
				}
			}
		}
		if len(indices) == oldIndexCount || len(vertices) == oldVertexCount {
			//fmt.Printf("MESH NOT DISPATCHED IDX %d, %d, VTX %d, %d\n",
			//      oldIndexCount, len(indices), oldVertexCount, len(vertices))
		}
	}

	g_octrees[id.octree].primitives = octreePrimitives

	if options == OCTREE_GRAPHICS {
		var primitiveCount int = len(cellPrimitives)
		for primitiveIndex := 0; primitiveIndex < primitiveCount; primitiveIndex++ {
			var octreePrimitiveIndex int32 = cellPrimitives[primitiveIndex]
			var cellPrimitive Primitive = octreePrimitives[octreePrimitiveIndex]
			if len(cellPrimitive.vertices) > 0 && len(cellPrimitive.indices) > 0 {
				var attributes []VertexAttribute = cellPrimitive.attributes
				var mesh MeshId = MeshCreate(mode, cellPrimitive.indexType,
					int32(len(cellPrimitive.indices))/cellPrimitive.indexByteStride,
					attributes, uint32(len(cellPrimitive.vertices))/cellPrimitive.vertexByteStride)

				MeshBegin(mesh)
				g_meshes[mesh.mesh].vertices = cellPrimitive.vertices
				g_meshes[mesh.mesh].indices = cellPrimitive.indices
				g_meshes[mesh.mesh].usePosition = cellPrimitive.usePosition
				g_meshes[mesh.mesh].useNormal = cellPrimitive.useNormal
				g_meshes[mesh.mesh].useColor = cellPrimitive.useColor
				g_meshes[mesh.mesh].useTexcoord = cellPrimitive.useTexcoord
				g_meshes[mesh.mesh].useTangent = cellPrimitive.useTangent
				g_meshes[mesh.mesh].useWeight = cellPrimitive.useWeight
				g_meshes[mesh.mesh].useJoint = cellPrimitive.useJoint
				g_meshes[mesh.mesh].min = cellPrimitive.min
				g_meshes[mesh.mesh].max = cellPrimitive.max
				MeshEnd(mesh)

				octreePrimitives[octreePrimitiveIndex].mesh = mesh
			} else {
				//fmt.Printf("Empty primitive... looks doubious\n")
				cellPrimitives[primitiveIndex] = -1
			}
		}
	}
	out = cellPrimitives
	return
}

func cellDispatchPrimitives(id OctreeId, level int32, cellMin math.V3, cellMax math.V3, parents []int32, children []int32) (out []int32) {

	out = children
	var mins []math.V3 = g_octrees[id.octree].mins
	var maxs []math.V3 = g_octrees[id.octree].maxs
	var parentCount int = len(parents)

	for parentIndex := 0; parentIndex < parentCount; parentIndex++ {
		var primitiveIndex int32 = parents[parentIndex]
		var primMin math.V3 = mins[primitiveIndex]
		var primMax math.V3 = maxs[primitiveIndex]

		if primMin.X <= cellMax.X && primMax.X >= cellMin.X &&
			primMin.Y <= cellMax.Y && primMax.Y >= cellMin.Y &&
			primMin.Z <= cellMax.Z && primMax.Z >= cellMin.Z {
			out = append(out, primitiveIndex)
		}
	}
	return
}

func octreeSplitModel(id OctreeId, model ModelId, level int32, opaques []int32, transparents []int32) {
	var maxLevel int32 = g_octrees[id.octree].maxLevel
	if level <= maxLevel { //&& ((len(opaques) > 0) || (len(transparents) > 0)) {
		var offset int32 = g_octrees[id.octree].offsets[level]
		g_octrees[id.octree].offsets[level] = offset + 1
		var cells []OctreeCell = g_octrees[id.octree].cells
		var cellMin math.V3 = cells[offset].min
		var cellMax math.V3 = cells[offset].max
		cells[offset].opaquePrimitives = cellDispatchPrimitives(
			id, level, cellMin, cellMax, opaques, cells[offset].opaquePrimitives)
		cells[offset].transparentPrimitives = cellDispatchPrimitives(
			id, level, cellMin, cellMax, transparents, cells[offset].transparentPrimitives)

		if level < maxLevel {
			var childLevel int32 = level + 1
			for i := 0; i < 8; i++ {
				octreeSplitModel(id, model, childLevel, cells[offset].opaquePrimitives, cells[offset].transparentPrimitives)
			}
		}
	}
}

func octreeUpdateCells(id OctreeId, level int32, center math.V3, size math.V3) (out int32) {
	out = -1
	var maxLevel int32 = g_octrees[id.octree].maxLevel
	if level <= maxLevel {
		var min math.V3 = v3.Sub(center, size)
		var max math.V3 = v3.Add(center, size)

		var offset int32 = g_octrees[id.octree].offsets[level]
		out = offset
		var cells []OctreeCell = g_octrees[id.octree].cells

		cells[offset].index = offset
		cells[offset].level = level
		cells[offset].center = center
		cells[offset].size = size
		cells[offset].min = min
		cells[offset].max = max

		var points []math.V3
		var p0 math.V3 = min
		var p1 math.V3 = v3.Make(max.X, min.Y, min.Z)
		var p2 math.V3 = v3.Make(min.X, max.Y, min.Z)
		var p3 math.V3 = v3.Make(max.X, max.Y, min.Z)
		var p4 math.V3 = v3.Make(min.X, min.Y, max.Z)
		var p5 math.V3 = v3.Make(max.X, min.Y, max.Z)
		var p6 math.V3 = v3.Make(min.X, max.Y, max.Z)
		var p7 math.V3 = v3.Make(max.X, max.Y, max.Z)

		points = append(points, p0)
		points = append(points, p1)
		points = append(points, p2)
		points = append(points, p3)
		points = append(points, p4)
		points = append(points, p5)
		points = append(points, p6)
		points = append(points, p7)

		var origins []math.V3
		var normals []math.V3

		var p40 math.V3 = v3.Normalize(v3.Sub(p4, p0))
		var p10 math.V3 = v3.Normalize(v3.Sub(p1, p0))
		var p20 math.V3 = v3.Normalize(v3.Sub(p2, p0))
		var p31 math.V3 = v3.Normalize(v3.Sub(p3, p1))
		var p51 math.V3 = v3.Normalize(v3.Sub(p5, p1))
		var p32 math.V3 = v3.Normalize(v3.Sub(p3, p2))
		var p62 math.V3 = v3.Normalize(v3.Sub(p6, p2))
		var p64 math.V3 = v3.Normalize(v3.Sub(p6, p4))
		var p54 math.V3 = v3.Normalize(v3.Sub(p5, p4))

		normals = append(normals, v3.Normalize(v3.Cross(p40, p20)))
		normals = append(normals, v3.Normalize(v3.Cross(p31, p51)))
		normals = append(normals, v3.Normalize(v3.Cross(p10, p40)))
		normals = append(normals, v3.Normalize(v3.Cross(p62, p32)))
		normals = append(normals, v3.Normalize(v3.Cross(p20, p10)))
		normals = append(normals, v3.Normalize(v3.Cross(p54, p64)))

		origins = append(origins, p0)
		origins = append(origins, p1)
		origins = append(origins, p0)
		origins = append(origins, p2)
		origins = append(origins, p0)
		origins = append(origins, p4)

		cells[offset].planesNormal = normals
		cells[offset].planesOrigin = origins
		cells[offset].points = points

		g_octrees[id.octree].offsets[level] = offset + 1

		var childLevel int32 = level + 1
		var childSize math.V3 = v3.Mulf(size, 0.5)
		for i := 0; i < 8; i++ {
			cells[offset].children[i] = octreeUpdateCells(id, childLevel, v3.Add(center, v3.Mul(g_octreeCells[i], childSize)), childSize)
		}
	}
	return
}

// OctreeIntersectsSphere ...
func OctreeIntersectsSphere(id OctreeId, position math.V3, radius float32, in []int32) (out []int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	_ /*var totalCount int32*/ = resetOffsets(id)
	g_octrees[id.octree].tmp = in
	_ /*var inter int32*/ = octreeIntersectsSphere(id, position, radius, 0)
	out = g_octrees[id.octree].tmp
	return
}

//OctreeGetCellOpaqueTriangleCount ...
func OctreeGetCellOpaqueTriangleCount(id OctreeId, cell int32) (out int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].opaqueTriangleCount
	return
}

//OctreeGetCellTransparentTriangleCount ...
func OctreeGetCellTransparentTriangleCount(id OctreeId, cell int32) (out int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].transparentTriangleCount
	return
}

//OctreeGetCellTriangleCount ...
func OctreeGetCellTriangleCount(id OctreeId, cell int32) (out int32) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].triangleCount
	return
}

// OctreeGetCellPositions ...
func OctreeGetCellPositions(id OctreeId, cell int32) (out []float32) {
	//utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].positions
	return
}

// OctreeGetCellNormals ...
func OctreeGetCellNormals(id OctreeId, cell int32) (out []float32) {
	out = g_octrees[id.octree].cells[cell].normals
	return
}

// OctreeGetCellCenter ...
func OctreeGetCellCenter(id OctreeId, cell int32) (out math.V3) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].center
	return
}

// OctreeGetCellMin ...
func OctreeGetCellMin(id OctreeId, cell int32) (out math.V3) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].min
	return
}

// OctreeGetCellMax ...
func OctreeGetCellMax(id OctreeId, cell int32) (out math.V3) {
	utils.PanicIfNot(OctreeIsValid(id), "invalid id")
	out = g_octrees[id.octree].cells[cell].max
	return
}

func octreeIntersectsSphere(id OctreeId, position math.V3, radius float32, level int32) (inter int32) {
	var out []int32 = g_octrees[id.octree].tmp
	var maxLevel int32 = g_octrees[id.octree].maxLevel
	if level <= maxLevel {
		var offsets []int32 = g_octrees[id.octree].offsets
		var offset int32 = offsets[level]
		offsets[level] = offset + 1

		var cells []OctreeCell = g_octrees[id.octree].cells
		inter = intersect.SphereIntersectsAABB(position, radius, cells[offset].min, cells[offset].max)
		var childLevel int32 = level + 1
		var recurse int32 = 8
		if inter <= 0 {
			if level == maxLevel {
				out = append(out, offset)
				var cells []OctreeCell = g_octrees[id.octree].cells
				if cells[offset].triangleCount > 0 {
					g_octrees[id.octree].tmp = out
				}
			} else {
				recurse = 0
				for i := 0; i < 8; i++ {
					var childInter int32
					childInter = octreeIntersectsSphere(id, position, radius, childLevel)
					if childInter == -1 {
						recurse = int32(8 - i - 1)
						i = 8
					}
				}
			}
		}
		if recurse > 0 {
			var layerCount int32 = recurse
			for i := childLevel; i <= maxLevel; i++ {
				offsets[i] = offsets[i] + layerCount
				layerCount = layerCount * 8
			}
		}
	}
	return
}

func octreeUpdateLevel(id OctreeId, level int32, frustum FrustumId, debug bool, targetLevel int32) {
	var maxLevel int32 = g_octrees[id.octree].maxLevel
	targetLevel = math.Min_i32(maxLevel, targetLevel)
	if level <= targetLevel {
		var offsets []int32 = g_octrees[id.octree].offsets
		var offset int32 = offsets[level]
		offsets[level] = offset + 1

		var cells []OctreeCell = g_octrees[id.octree].cells

		var inter int32 = FrustumIntersectsAABB(frustum, cells[offset].min, cells[offset].max)

		var childLevel int32 = level + 1
		var recurse bool
		if inter == 0 {
			if level == targetLevel {
				g_octrees[id.octree].visibles = append(g_octrees[id.octree].visibles, offset)
			} else {
				recurse = true
				for i := 0; i < 8; i++ {
					octreeUpdateLevel(id, childLevel, frustum, debug, targetLevel)
				}
			}
		} else if inter == -1 {
			g_octrees[id.octree].visibles = append(g_octrees[id.octree].visibles, offset)
		} else if inter == 1 {
		}

		if recurse == false {
			var layerCount int32 = 8
			for i := childLevel; i <= maxLevel; i++ {
				offsets[i] = offsets[i] + layerCount
				layerCount = layerCount * 8
			}
		}
	}
}
