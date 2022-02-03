package gfx

import (
	//	"fmt"
	"fmt"
	"skyfx/gfx/gltf"
	"skyfx/math"
	m44 "skyfx/math/m44"
	q4 "skyfx/math/q4"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"
	"skyfx/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// import "m44"
// import "q4"

// // Constants ...
var MODEL_GEOMETRY int32 = gltf.ASSET_GEOMETRY
var MODEL_ANIMATION int32 = gltf.ASSET_ANIMATION

// // Globals ...
var g_models []Model
var g_animations []Animation // ISSUE : can't be named g_animations as it clashes at runtime with src/gui/animation.cx::g_animations

// AnimationId ...
type AnimationId struct {
	animation int32
}

// AnimationSource ...
type AnimationSource struct {
	channels   []int32
	frames     []math.TRS
	transforms []math.M44
	anim       bool
	t          math.V3
	r          math.V4
	s          math.V3
	transform  math.M44
}

// Animation ...
type Animation struct {
	asset      gltf.AssetId
	sources    []AnimationSource
	time       float32
	direction  float32
	min        float32
	max        float32
	trs        []math.TRS
	transforms []math.M44
	keys       []float32
	lastIndex  int32
	runtime    bool
}

// AnimationIsValid ...
func AnimationIsValid(id AnimationId) (out bool) {
	out = id.animation >= 0 && id.animation < int32(len(g_animations))
	return
}

// AnimationInvalid ...
func AnimationInvalid() (out AnimationId) {
	out.animation = -1
	return
}

// AnimationCreate ...
func AnimationCreate() (out AnimationId) {
	out.animation = int32(len(g_animations))

	var animation Animation
	animation.direction = 1.0
	g_animations = append(g_animations, animation)
	return
}

type Skeleton struct {
	root int32
}

type Primitive struct {
	mesh      MeshId
	node      int32
	nodeIndex int32

	effect EffectId

	useSkin           int32
	baseTexture       TextureId
	metalRoughTexture TextureId
	emissiveTexture   TextureId
	normalTexture     TextureId
	occlusionTexture  TextureId

	gltfPrimitive gltf.Primitive
	gltfMaterial  gltf.Material

	min              math.V3
	max              math.V3
	vertices         []uint8
	indices          []uint8
	attributes       []VertexAttribute
	indexByteStride  int32
	vertexByteStride uint32
	indexType        uint32
	usePosition      bool
	useNormal        bool
	useColor         bool
	useTexcoord      bool
	useTangent       bool
	useWeight        bool
	useJoint         bool

	hash0 uint64
	hash1 uint64
}

func PrimitiveComputeHash(p Primitive) (hash0 uint64, hash1 uint64) { //
	hash0 = uint64(p.effect.effect) |
		uint64(p.baseTexture.texture)<<8 |
		uint64(p.normalTexture.texture)<<16 |
		uint64(p.metalRoughTexture.texture)<<24 |
		uint64(p.emissiveTexture.texture)<<32 |
		uint64(p.occlusionTexture.texture)<<40
	return
}

type PrimitiveBucket struct {
	hash0             uint64
	hash1             uint64
	primitives        []int32
	effect            EffectId
	baseTexture       TextureId
	normalTexture     TextureId
	metalRoughTexture TextureId
	emissiveTexture   TextureId
	occlusionTexture  TextureId
}

func PrimitiveSort1(primitives []Primitive, indices []int32, buckets []PrimitiveBucket) (out []PrimitiveBucket) {
	out = buckets
	var meshCount int = len(indices)
	var bucketCount int = len(buckets)
	for i := 0; i < bucketCount; i++ {
		out[i].primitives = out[i].primitives[:0]
	}
	for i := 0; i < meshCount; i++ {
		var mesh int32 = indices[i]
		var hash0 uint64 = primitives[mesh].hash0
		var hash1 uint64 = primitives[mesh].hash1

		var bucketIndex int32 = -1
		var outCount int = len(out)
		for bucket := 0; bucket < outCount; bucket++ {
			if out[bucket].hash0 == hash0 && out[bucket].hash1 == hash1 {
				bucketIndex = int32(bucket)
				bucket = outCount
			}
		}

		// TODO REMOVE var outPrimitives []int32
		if bucketIndex == -1 {
			if IsValidMesh(primitives[mesh].mesh) {
				var newBucket PrimitiveBucket
				newBucket.hash0 = hash0
				newBucket.hash1 = hash1
				newBucket.effect = primitives[mesh].effect
				newBucket.baseTexture = primitives[mesh].baseTexture
				newBucket.normalTexture = primitives[mesh].normalTexture
				newBucket.metalRoughTexture = primitives[mesh].metalRoughTexture
				newBucket.emissiveTexture = primitives[mesh].emissiveTexture
				newBucket.occlusionTexture = primitives[mesh].occlusionTexture
				newBucket.primitives = append(newBucket.primitives, mesh)
				out = append(out, newBucket)
			}
		} else {
			out[bucketIndex].primitives = append(out[bucketIndex].primitives, mesh)
		}
	}
	return
}

// func PrimitiveSort0(primitives []Primitive, indices []int32, in []int32) (out []int32) {
//     out = in
//     var meshCount int32 = len(indices)
//     for i := 0; i < meshCount; i++ {
//         var mesh int32 = indices[i]
//         var p Primitive = primitives[mesh]
//         var sortedCount int32 = len(out)
//         for si := 0; si < sortedCount; si++ {
//             var smesh int32 = out[si]
//             var sp Primitive = primitives[smesh]
//             var in bool
//             if p.effect.effect < sp.effect.effect {
//                 in = true
//             } else if p.effect.effect == sp.effect.effect {
//                 if p.baseTexture.texture < sp.baseTexture.texture {
//                     in = true
//                 } else if p.baseTexture.texture == sp.baseTexture.texture {
//                     if p.NormalTexture.texture < sp.NormalTexture.texture {
//                         in = true
//                     } else if p.NormalTexture.texture == sp.NormalTexture.texture {
//                         if p.metalRoughTexture.texture < sp.metalRoughTexture.texture {
//                             in = true
//                         } else if p.metalRoughTexture.texture == sp.metalRoughTexture.texture {
//                             if p.emissiveTexture.texture < sp.emissiveTexture.texture {
//                                 in = true
//                             } else if p.emissiveTexture.texture == sp.emissiveTexture.texture {
//                                 if p.occlusionTexture.texture < sp.occlusionTexture.texture {
//                                     in = true
//                                 }
//                             }
//                         }
//                     }
//                 }
//             }
//             if in {
//                 out = insert(out, si, mesh)
//                 si = sortedCount
//             }
//         }
//         if sortedCount == len(out) {
//             out = append(out, mesh)
//         }
//     }
// }

// ModelId ...
type ModelId struct {
	model int32
}

// Model ...
type Model struct {
	asset gltf.AssetId

	dataDir  string
	filename string

	animated bool
	skinned  bool

	min math.V3
	max math.V3

	nodes        []int32
	nodesMap     []int32
	determinants []float32
	transforms   []math.M44
	matrixStack  []math.M44

	meshes            []int32
	meshCacheKey      []int32
	meshCacheVal      []MeshId
	primitives        []Primitive
	opaqueMeshes      []int32
	transparentMeshes []int32

	skeleton   Skeleton
	animations []AnimationId

	worldIsIdentity bool
	world           math.M44 // TODO : optim inverse
	view            math.M44 // TODO : remove
	projection      math.M44 // TODO : remove

	cameraPosition      math.V4   // TODO : remove
	environmentSpecular TextureId // TODO : remove
	environmentDiffuse  TextureId // TODO : remove
	brdf                TextureId // TODO : remove

	exposure float32

	inverseBinds  []math.M44 // TODO : skeleton struct
	jointMatrices []math.M44 // TODO : runtime struct

	joints []int32 // TODO : skeleton struct
	stack  []int32 // TODO : remove
}

// ModelIsValid ...
func ModelIsValid(id ModelId) (out bool) {
	out = id.model >= 0 && id.model < int32(len(g_models))
	return
}

// ModelInvalid ...
func ModelInvalid() (out ModelId) {
	out.model = -1
	return
}

// ModelCreate ...
func ModelCreate() (out ModelId) {
	out.model = int32(len(g_models))

	var model Model
	model.skeleton.root = -1
	model.asset = gltf.AssetInvalid()
	g_models = append(g_models, model)
	ModelSetWorld(out, m44.IDENTITY)
	return
}

// ModelCreateFromFile ...
func ModelCreateFromFile(dataDir string, filename string, options int32) (out ModelId) { // TODO : only load animations
	out = ModelCreate()

	var modelIndex int32 = out.model
	g_models[modelIndex].dataDir = dataDir
	g_models[modelIndex].filename = filename

	var asset gltf.AssetId = gltf.AssetCreate(dataDir, filename, options)
	if gltf.AssetIsValid(asset) == false {
		fmt.Printf("gltf.AssetCreate failed\n")
		out = ModelInvalid()
		return
	}

	g_models[modelIndex].asset = asset
	//gltf.AssetPrint(asset)

	if modelInstance(out, options) == false {
		fmt.Printf("modelInstance failed\n")
		out = ModelInvalid()
		return
	}

	modelTransform(out, AnimationInvalid(), 0.0, false)
	return
}

// ModelGetPath ...
func ModelGetPath(id ModelId) (out string) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	out = g_models[id.model].dataDir
	return
}

// ModelGetName ...
func ModelGetName(id ModelId) (out string) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	out = g_models[id.model].filename
	return
}

// ModelGetMin ...
func ModelGetMin(id ModelId) (out math.V3) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	out = g_models[id.model].min
	return
}

// ModelGetMax ...
func ModelGetMax(id ModelId) (out math.V3) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	out = g_models[id.model].max
	return
}

// ModelGetAnimation ...
func ModelGetAnimation(id ModelId, animationIndex int32) (out AnimationId) {
	out = AnimationInvalid()
	utils.PanicIfNot(ModelIsValid(id), "invalid model")

	var modelIndex int32 = id.model
	var animationCount int32 = int32(len(g_models[modelIndex].animations))
	if animationIndex >= 0 && animationIndex < animationCount {
		out = g_models[modelIndex].animations[animationIndex]
	}
	return
}

// // ModelGetNodeByName ...
// func ModelGetNodeByName(id ModelId, name string) (out int32) {
// 	out = -1
// 	utils.PanicIfNot(ModelIsValid(id), "invalid model")

// 	var modelIndex int32 = id.model
// 	var asset gltf.AssetId = g_models[modelIndex].asset
// 	var nodeCount int32 = gltf.AssetGetNodeCount(asset)
// 	for i := 0; i < nodeCount; i++ {
// 		var nodeName string = gltf.NodeGetName(asset, i)
// 		if nodeName == name {
// 			out = i
// 			return
// 		}
// 	}
// }

type channelInfo struct {
	count    uint32
	dataType int32
	dataLen  uint32
	dataF32  []float32
	dataUI8  []uint8
	dataUI16 []uint16
	dataUI32 []uint32
}

func effectInstance(primitive Primitive) (out EffectId) {
	var mesh MeshId = primitive.mesh

	var useColorVtx bool = MeshUseColor(mesh)
	var useColorMap bool = IsValidTexture(primitive.baseTexture)
	var useColorUni bool = true

	var usePbrMap bool = IsValidTexture(primitive.metalRoughTexture)
	var usePbrUni bool = usePbrMap == false

	var useEmissiveMap bool = IsValidTexture(primitive.emissiveTexture)
	var useEmissiveUni bool = useEmissiveMap == false

	var useNormalMap bool = IsValidTexture(primitive.normalTexture)
	var useNormalVtx bool = MeshUseNormal(mesh)
	var useTangentVtx bool = MeshUseTangent(mesh)

	var useOcclusionTexture bool = IsValidTexture(primitive.occlusionTexture)

	var useWeights bool = MeshUseWeight(mesh)
	var useJoints bool = MeshUseJoint(mesh)
	var useSkin bool = useWeights && useJoints

	/*fmt.Printf("EFFECT_INSTANCE\n")
	if useColorVtx {
		fmt.Printf("USE_COLOR_VTX\n")
	}
	if useColorMap {
		fmt.Printf("USE_COLOR_MAP\n")
	}
	if useColorUni {
		fmt.Printf("USE_COLOR_UNI\n")
	}
	if usePbrMap {
		fmt.Printf("USE_PBR_MAP\n")
	}
	if usePbrUni {
		fmt.Printf("USE_PBR_UNI\n")
	}
	if useEmissiveMap {
		fmt.Printf("USE_EMIS_MAP\n")
	}
	if useEmissiveUni {
		fmt.Printf("USE_EMIS_UNI\n")
	}
	if useNormalMap {
		fmt.Printf("USE_NORMAL_MAP\n")
	}
	if useNormalVtx {
		fmt.Printf("USE_NORMAL_VTX\n")
	}
	if useTangentVtx {
		fmt.Printf("USE_TANGENT_VTX\n")
	}
	if useOcclusionTexture {
		fmt.Printf("USE_OCCLUSION_TEXTURE\n")
	}
	if useWeights {
		fmt.Printf("USE_WEIGHTS\n")
	}
	if useJoints {
		fmt.Printf("USE_JOINTS\n")
	}
	if useSkin {
		fmt.Printf("USE_SKIN\n")
	}*/
	TemplateSetKey(g_tfxPbr, USE_COLOR_UNI, useColorUni)
	TemplateSetKey(g_tfxPbr, USE_COLOR_VTX, useColorVtx)
	TemplateSetKey(g_tfxPbr, USE_COLOR_MAP, useColorMap)

	TemplateSetKey(g_tfxPbr, USE_PBR_UNI, usePbrUni)
	TemplateSetKey(g_tfxPbr, USE_PBR_MAP, usePbrMap)

	TemplateSetKey(g_tfxPbr, USE_EMISSIVE_UNI, useEmissiveUni)
	TemplateSetKey(g_tfxPbr, USE_EMISSIVE_MAP, useEmissiveMap)

	TemplateSetKey(g_tfxPbr, USE_NORMAL_VTX, useNormalVtx)
	TemplateSetKey(g_tfxPbr, USE_NORMAL_MAP, useNormalMap)
	TemplateSetKey(g_tfxPbr, USE_TANGENT_VTX, useTangentVtx)

	TemplateSetKey(g_tfxPbr, USE_OCCLUSION_MAP, useOcclusionTexture)

	TemplateSetKey(g_tfxPbr, USE_SKIN, useSkin)

	TemplateSetKey(g_tfxPbr, USE_DEBUG_A, true)

	out = TemplateInstance(g_tfxPbr)
	return
}

func meshInstance(id gltf.AssetId, primitive gltf.Primitive) (out MeshId) {
	var mode uint32
	//fmt.Printf("instancing : mode %d, indices %d, material %d\n", primitive.mode, primitive.indices, primitive.material)
	if primitive.Mode == gltf.PRIMITIVE_LINES {
		mode = gl.LINES
	} else if primitive.Mode == gltf.PRIMITIVE_TRIANGLES {
		mode = gl.TRIANGLES
	} else {
		fmt.Printf("primitive type not implemented\n")
		out = InvalidMesh()
		return
	}

	var min math.V3
	var max math.V3

	var cpuMin math.V3 = v3.MAX
	var cpuMax math.V3 = v3.MIN

	var usePosition bool
	var useNormal bool
	var useColor bool
	var useTexcoord bool
	var useTangent bool
	var useWeight bool
	var useJoint bool

	// attributes
	var channels []channelInfo
	var attributes []VertexAttribute
	var vertexCount uint32 = 0
	var attributeCount int = len(primitive.Attributes)
	var byteStride uint32 = 0
	var componentStride uint32 = 0
	for a := 0; a < attributeCount; a++ {
		var accessor gltf.Accessor = gltf.AssetGetAccessor(id, primitive.Attributes[a].Accessor)
		utils.PanicIfNot(accessor.Loaded, "attributes accessor is not loaded")
		if a > 0 && uint32(accessor.Count) != vertexCount {
			fmt.Printf("wrong number of vectices\n")
			out = InvalidMesh()
			return
		}

		vertexCount = uint32(accessor.Count)

		var channel channelInfo
		var attribute VertexAttribute

		if accessor.ComponentType == gl.FLOAT { // TODO : how to pass []struct to BufferData ?
			channel.dataF32 = accessor.DataF32
			channel.dataLen = uint32(len(channel.dataF32))
			channel.dataType = gl.FLOAT
			attribute.componentType = gl.FLOAT
			attribute.componentByteSize = g_sizeofF32
		} else if accessor.ComponentType == gl.UNSIGNED_SHORT {
			channel.dataUI16 = accessor.DataUI16
			channel.dataLen = uint32(len(channel.dataUI16))
			channel.dataType = gl.UNSIGNED_SHORT
			attribute.componentType = gl.UNSIGNED_SHORT
			attribute.componentByteSize = g_sizeofUI16
		} else if accessor.ComponentType == gl.UNSIGNED_INT {
			channel.dataUI32 = accessor.DataUI32
			channel.dataLen = uint32(len(channel.dataUI32))
			channel.dataType = gl.UNSIGNED_INT
			attribute.componentType = gl.UNSIGNED_INT
			attribute.componentByteSize = g_sizeofUI32
		} else {
			fmt.Printf("buffer type not implemented %d\n", accessor.ComponentType)
			out = InvalidMesh()
			return
		}

		attribute.componentCount = accessor.ComponentCount
		attribute.componentOffset = componentStride
		attribute.byteOffset = byteStride

		var attributeType int32 = primitive.Attributes[a].AttributeType
		var binding int32 = -1
		if attributeType == gltf.ATTRIBUTE_POSITION {
			binding = ATTRIBUTE_POSITION
			min.X = accessor.Min[0]
			min.Y = accessor.Min[1]
			min.Z = accessor.Min[2]
			max.X = accessor.Max[0]
			max.Y = accessor.Max[1]
			max.Z = accessor.Max[2]
			usePosition = true
		} else if attributeType == gltf.ATTRIBUTE_NORMAL {
			binding = ATTRIBUTE_NORMAL
			useNormal = true
		} else if attributeType == gltf.ATTRIBUTE_COLOR {
			binding = ATTRIBUTE_COLOR
			useColor = true
		} else if attributeType == gltf.ATTRIBUTE_TEXCOORD {
			binding = ATTRIBUTE_TEXCOORD
			useTexcoord = true
		} else if attributeType == gltf.ATTRIBUTE_TANGENT {
			binding = ATTRIBUTE_TANGENT
			useTangent = true
		} else if attributeType == gltf.ATTRIBUTE_WEIGHT {
			binding = ATTRIBUTE_WEIGHT
			useWeight = true
		} else if attributeType == gltf.ATTRIBUTE_JOINT {
			binding = ATTRIBUTE_JOINT
			useJoint = true
		}

		attribute.binding = uint32(binding)
		channel.count = attribute.componentCount

		attributes = append(attributes, attribute)
		channels = append(channels, channel)

		byteStride = byteStride + attribute.componentCount*attribute.componentByteSize
		componentStride = componentStride + attribute.componentCount
	}

	var vertexLen uint32 = componentStride * vertexCount
	if vertexLen <= 0 {
		fmt.Printf("invalid vertex buffer\n")
		out = InvalidMesh()
		return
	}

	// vertices // TODO : generate buffers from go data ?
	var vertices []uint8 // TODO : test perf with loop reordering + resize(vertices, vertexLen)
	for v := uint32(0); v < vertexCount; v++ {
		for c := 0; c < attributeCount; c++ {
			var dataType int32 = channels[c].dataType
			//var dataLen int32 = channels[c].dataLen
			var count uint32 = channels[c].count
			var offset uint32 = v * count
			if dataType == gl.FLOAT {
				var data []float32 = channels[c].dataF32
				for i := uint32(0); i < count; i++ {
					vertices = AppendF32(vertices, data[offset+i])
				}
				if attributes[c].binding == ATTRIBUTE_POSITION {
					cpuMin.X = math.Min_f32(cpuMin.X, data[offset+0])
					cpuMin.Y = math.Min_f32(cpuMin.Y, data[offset+1])
					cpuMin.Z = math.Min_f32(cpuMin.Z, data[offset+2])
					cpuMax.X = math.Max_f32(cpuMax.X, data[offset+0])
					cpuMax.Y = math.Max_f32(cpuMax.Y, data[offset+1])
					cpuMax.Z = math.Max_f32(cpuMax.Z, data[offset+2])
				}
			} else if dataType == gl.UNSIGNED_SHORT {
				var data []uint16 = channels[c].dataUI16
				for i := uint32(0); i < count; i++ {
					vertices = AppendUI16(vertices, data[offset+i])
				}
			} else if dataType == gl.UNSIGNED_INT {
				var data []uint32 = channels[c].dataUI32
				for i := uint32(0); i < count; i++ {
					vertices = AppendUI32(vertices, data[offset+i])
				}
			} else {
				utils.PanicIf(true, "unhandled dataType")
			}
		}
	}

	// indices
	var indices []uint8
	var indicesAccessor gltf.Accessor = gltf.AssetGetAccessor(id, primitive.Indices)
	utils.PanicIfNot(indicesAccessor.Loaded, "indices accessor is not loaded")
	var indicesType uint32 = uint32(indicesAccessor.ComponentType)
	if indicesAccessor.ComponentCount != 1 ||
		(indicesAccessor.ComponentType != gl.UNSIGNED_SHORT && indicesAccessor.ComponentType != gl.UNSIGNED_INT) ||
		indicesAccessor.AttributeTypeEnum != gltf.TYPE_SCALAR {
		fmt.Printf("invalid index buffer format : componentCount %d, componentType %d, attributeType %d\n",
			indicesAccessor.ComponentCount, indicesAccessor.ComponentType, indicesAccessor.AttributeTypeEnum)
		out = InvalidMesh()
		return
	}

	var indexByteCount int32 = 0
	if indicesType == gl.UNSIGNED_SHORT {
		var data []uint16 = indicesAccessor.DataUI16
		var count int = len(data)
		indexByteCount = 2
		for i := 0; i < count; i++ {
			indices = AppendUI16(indices, data[i])
		}
	} else if indicesType == gl.UNSIGNED_INT {
		var data []uint32 = indicesAccessor.DataUI32
		var count int = len(data)
		indexByteCount = 4
		for i := 0; i < count; i++ {
			indices = AppendUI32(indices, data[i])
		}
	} else {
		fmt.Printf("invalid index buffer type\n")
		out = InvalidMesh()
		return
	}

	if len(attributes) > 0 && len(indices) > 0 && len(vertices) > 0 {
		out = MeshCreate(mode, indicesType, int32(len(indices))/indexByteCount, attributes, uint32(len(vertices))/byteStride) // TODO : use real index/vertex count instead of bytecount
		MeshBegin(out)                                                                                                        // TODO : remove MeshBegin/MeshEnd, data should be uploaded with BufferData in MeshInstance.
		g_meshes[out.mesh].channels = channels
		g_meshes[out.mesh].vertices = vertices
		g_meshes[out.mesh].indices = indices
		g_meshes[out.mesh].usePosition = usePosition
		g_meshes[out.mesh].useNormal = useNormal
		g_meshes[out.mesh].useColor = useColor
		g_meshes[out.mesh].useTexcoord = useTexcoord
		g_meshes[out.mesh].useTangent = useTangent
		g_meshes[out.mesh].useWeight = useWeight
		g_meshes[out.mesh].useJoint = useJoint
		g_meshes[out.mesh].min = min
		g_meshes[out.mesh].max = max
		MeshEnd(out)
		/*fmt.Printf("DELTA_MIN %f, %f, %f, DELTA_MAX %f, %f, %f\n",
		  cpuMin.X - min.X, cpuMin.Y - min.Y, cpuMin.Z - min.Z,
		  cpuMax.X - max.X, cpuMax.Y - max.Y, cpuMax.Z - max.Z)*/

	}
	return
}

// ModelSetCameraPosition ...
func ModelSetCameraPosition(id ModelId, position math.V3) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	var cameraPosition math.V4
	cameraPosition.X = position.X
	cameraPosition.Y = position.Y
	cameraPosition.Z = position.Z
	cameraPosition.W = 1.0
	g_models[id.model].cameraPosition = cameraPosition
}

// ModelSetEnvironmentSpecular ...
func ModelSetEnvironmentSpecular(id ModelId, environment TextureId) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].environmentSpecular = environment
}

// ModelSetEnvironmentDiffuse...
func ModelSetEnvironmentDiffuse(id ModelId, environment TextureId) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].environmentDiffuse = environment
}

// ModelSetExposure...
func ModelSetExposure(id ModelId, exposure float32) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].exposure = exposure
}

// ModelSetBRDF ...
func ModelSetBRDF(id ModelId, brdf TextureId) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].brdf = brdf
}

// ModelSetWorld ...
func ModelSetWorld(id ModelId, world math.M44) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].world = world
	g_models[id.model].worldIsIdentity = m44.Isident(world)
}

// ModelSetView ...
func ModelSetView(id ModelId, view math.M44) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].view = view
}

// ModelSetProjection ...
func ModelSetProjection(id ModelId, projection math.M44) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	g_models[id.model].projection = projection
}

func loadTexture(asset gltf.AssetId, dataDir string, index int32) (out TextureId) {
	out = InvalidTexture()
	//if gltf.AssetIsValidTexture(asset, index) {
	var texture gltf.Texture = gltf.AssetGetTexture(asset, index)
	var image gltf.Image = gltf.AssetGetImage(asset, texture.Source)
	if image.Uri != "" {
		out = TextureCreate(fmt.Sprintf("%s%s", dataDir, image.Uri), FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)
	}
	//}
	return
}

// ModelUpdate ...
func ModelUpdate(id ModelId, animId AnimationId, deltaTime float32, loop bool) { // TODO : bound check
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	if g_models[id.model].animated || g_models[id.model].skinned {
		modelTransform(id, animId, deltaTime, loop)
	}
}

// ModelRender ...
func ModelRender(id ModelId) {
	DisableBlending()
	DepthState(true, gl.LESS, true)
	modelRender(id, g_models[id.model].opaqueMeshes, 1.0)

	EnableBlending(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	DepthState(true, gl.LESS, false)
	modelRender(id, g_models[id.model].transparentMeshes, 1.0)
}

func modelInstance(id ModelId, options int32) (success bool) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")
	var modelIndex int32 = id.model

	var stack []int32 = g_models[modelIndex].stack

	var asset gltf.AssetId = g_models[modelIndex].asset
	var roots []int32 = gltf.AssetGetRootNodes(asset)
	var rootCount int = len(roots)

	var animations []AnimationId = g_models[modelIndex].animations

	var nodes []int32 = g_models[modelIndex].nodes
	var nodesMap []int32 = g_models[modelIndex].nodesMap

	var determinants []float32 = g_models[modelIndex].determinants
	var meshCacheKey []int32 = g_models[modelIndex].meshCacheKey
	var meshCacheVal []MeshId = g_models[modelIndex].meshCacheVal

	var meshes []int32 = g_models[modelIndex].meshes

	var dataDir string = g_models[modelIndex].dataDir
	var primitives []Primitive = g_models[modelIndex].primitives

	var transparentMeshes []int32 = g_models[modelIndex].transparentMeshes
	var opaqueMeshes []int32 = g_models[modelIndex].opaqueMeshes

	var matrixStack []math.M44 = g_models[modelIndex].matrixStack
	matrixStack = matrixStack[:0]
	matrixStack = m44.Push(matrixStack, m44.IDENTITY)

	var transforms []math.M44 = g_models[modelIndex].transforms
	transforms = transforms[:0]

	var min math.V3 = v3.MAX
	var max math.V3 = v3.MIN

	stack = make([]int32, rootCount)
	for i := 0; i < rootCount; i++ {
		stack[i] = roots[i]
	}

	var nodeCount int32 = gltf.AssetGetNodeCount(asset)

	if (options & MODEL_ANIMATION) != 0 {
		var gltfAnimations []gltf.Animation = gltf.AssetGetAnimations(asset)
		var animationCount int = len(gltfAnimations)
		g_models[modelIndex].animated = animationCount > 0 // TODO : multiple models per scene
		for animIndex := 0; animIndex < animationCount; animIndex++ {
			var channels []gltf.AnimationChannel = gltfAnimations[animIndex].Channels
			var samplers []gltf.AnimationSampler = gltfAnimations[animIndex].Samplers

			var animation Animation
			animation.min, animation.max = gltf.AnimationGetLength(asset, int32(animIndex))

			var sources []AnimationSource
			for nodeIndex := int32(0); nodeIndex < nodeCount; nodeIndex++ {
				var source AnimationSource
				source.t = gltf.NodeGetTranslation(asset, nodeIndex)
				source.r = gltf.NodeGetRotation(asset, nodeIndex)
				source.s = gltf.NodeGetScale(asset, nodeIndex)
				sources = append(sources, source)
			}

			var channelCount int = len(channels)
			for channelIndex := 0; channelIndex < channelCount; channelIndex++ {
				var channel gltf.AnimationChannel = channels[channelIndex]
				var nodeIndex int32 = channel.Node
				sources[nodeIndex].channels = append(sources[nodeIndex].channels, int32(channelIndex))
			}

			var keys []float32 = animation.keys
			for nodeIndex := int32(0); nodeIndex < nodeCount; nodeIndex++ {
				var nodeChannels []int32 = sources[nodeIndex].channels
				var nodeChannelCount int = len(nodeChannels)
				for nodeChannelIndex := 0; nodeChannelIndex < nodeChannelCount; nodeChannelIndex++ {
					var channelIndex int32 = nodeChannels[nodeChannelIndex]
					var channel gltf.AnimationChannel = channels[channelIndex]
					var inputAccessor gltf.Accessor = gltf.AssetGetAccessor(asset, samplers[channel.Sampler].InputA)
					utils.PanicIfNot(inputAccessor.Loaded, "inputs accessor is not loaded")
					var inputs []float32 = inputAccessor.DataF32
					var inputCount int = len(inputs)
					for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
						var input float32 = inputs[inputIndex]
						var keyCount int = len(keys)
						var keyAdded bool
						for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
							var key float32 = keys[keyIndex]
							if input < key {
								keys = append(keys[:keyIndex+1], keys[keyIndex:]...)
								keys[keyIndex] = input
								keyIndex = len(keys) // ISSUE break
								keyAdded = true
							} else if input == key {
								keyIndex = len(keys) // ISSUE break
								keyAdded = true
							}
						}
						if keyAdded == false {
							keys = append(keys, input)
						}
					}
				}
			}

			var keyCount int = len(keys)
			for nodeIndex := 0; nodeIndex < int(nodeCount); nodeIndex++ {
				// TODO REMOVE var frames []math.TRS = sources[nodeIndex].frames
				var nodeChannels []int32 = sources[nodeIndex].channels
				var nodeChannelCount int = len(nodeChannels)
				for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
					var key float32 = keys[keyIndex]

					var animTranslation math.V3 = sources[nodeIndex].t
					var animRotation math.V4 = sources[nodeIndex].r
					var animScale math.V3 = sources[nodeIndex].s

					for nodeChannelIndex := 0; nodeChannelIndex < nodeChannelCount; nodeChannelIndex++ {
						var channelIndex int32 = nodeChannels[nodeChannelIndex]
						var channel gltf.AnimationChannel = channels[channelIndex]
						var inputAccessor gltf.Accessor = gltf.AssetGetAccessor(asset, samplers[channel.Sampler].InputA)
						utils.PanicIfNot(inputAccessor.Loaded, "inputs accessor is not loaded")
						var inputs []float32 = inputAccessor.DataF32
						var inputCount int = len(inputs)
						var index int = 0
						for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
							if key <= inputs[inputIndex] {
								index = inputIndex
								inputIndex = inputCount // ISSUE : break
							}
						}

						var prev int = index - 1
						if prev < 0 {
							prev = 0
						}

						var delta float32 = key - inputs[prev]
						var step float32 = inputs[index] - inputs[prev]
						var st float32 = 1.0
						if step > 0.0 {
							st = delta / step
						}
						var t float32 = st
						if t < 0.0 {
							t = 0.0
						} else if t > 1.0 {
							t = 1.0
						}

						var outputAccessor gltf.Accessor = gltf.AssetGetAccessor(asset, samplers[channel.Sampler].OutputA)
						utils.PanicIfNot(outputAccessor.Loaded, "outputs accessor is not loaded")
						var outputs []float32 = outputAccessor.DataF32
						var path int32 = channel.Path
						if path == gltf.ANIMATION_PATH_TRANSLATION {
							utils.PanicIf(outputAccessor.AttributeTypeEnum != gltf.TYPE_VEC3, "invalid attribute type")
							utils.PanicIf(outputAccessor.ComponentType != gl.FLOAT, "unhandled component type")
							var offset int = prev * 3
							var x0 float32 = outputs[offset]
							var y0 float32 = outputs[offset+1]
							var z0 float32 = outputs[offset+2]

							offset = index * 3
							var x1 float32 = outputs[offset]
							var y1 float32 = outputs[offset+1]
							var z1 float32 = outputs[offset+2]

							animTranslation.X = x0 + (x1-x0)*t
							animTranslation.Y = y0 + (y1-y0)*t
							animTranslation.Z = z0 + (z1-z0)*t
						} else if path == gltf.ANIMATION_PATH_ROTATION {
							utils.PanicIf(outputAccessor.AttributeTypeEnum != gltf.TYPE_VEC4, "invalid attribute type")
							utils.PanicIf(outputAccessor.ComponentType != gl.FLOAT, "unhandled component type")
							var offset int = prev * 4
							var x0 float32 = outputs[offset]
							var y0 float32 = outputs[offset+1]
							var z0 float32 = outputs[offset+2]
							var w0 float32 = outputs[offset+3]

							offset = index * 4
							var x1 float32 = outputs[offset]
							var y1 float32 = outputs[offset+1]
							var z1 float32 = outputs[offset+2]
							var w1 float32 = outputs[offset+3]

							animRotation.X = x0 + (x1-x0)*t
							animRotation.Y = y0 + (y1-y0)*t
							animRotation.Z = z0 + (z1-z0)*t
							animRotation.W = w0 + (w1-w0)*t
							animRotation = v4.Normalize(animRotation)
						} else if path == gltf.ANIMATION_PATH_SCALE {
							utils.PanicIf(outputAccessor.AttributeTypeEnum != gltf.TYPE_VEC3, "invalid attribute type")
							utils.PanicIf(outputAccessor.ComponentType != gl.FLOAT, "unhandled component type")
							var offset int = prev * 3
							var x0 float32 = outputs[offset]
							var y0 float32 = outputs[offset+1]
							var z0 float32 = outputs[offset+2]

							offset = index * 3
							var x1 float32 = outputs[offset]
							var y1 float32 = outputs[offset+1]
							var z1 float32 = outputs[offset+2]

							animScale.X = x0 + (x1-x0)*t
							animScale.Y = y0 + (y1-y0)*t
							animScale.Z = z0 + (z1-z0)*t
						} else if path == gltf.ANIMATION_PATH_WEIGHTS {
						} else {
							utils.PanicIf(true, "invalid animation path")
						}
					}

					var frame math.TRS
					frame.T = animTranslation
					frame.R = animRotation
					frame.S = animScale
					if nodeChannelCount > 0 {
						sources[nodeIndex].anim = true
					}
					sources[nodeIndex].frames = append(sources[nodeIndex].frames, frame)
					sources[nodeIndex].transforms = append(sources[nodeIndex].transforms, m44.Make_SQT(frame))
				}
				sources[nodeIndex].transform = gltf.NodeGetMatrix(asset, int32(nodeIndex))
			}
			animation.keys = keys
			animation.sources = sources
			animation.direction = 1.0
			animation.asset = asset

			var animId AnimationId = AnimationCreate()
			animations = append(animations, animId)
			g_animations[animId.animation] = animation
		}
	}

	nodesMap = make([]int32, nodeCount)
	for i := 0; i < int(nodeCount); i++ {
		nodesMap[i] = -1
	}

	for len(stack) > 0 {
		var stackLen int = len(stack)

		stackLen--
		var node int32 = stack[stackLen]
		stack = stack[:stackLen] //resize(stack, stackLen)

		if node < 0 {
			matrixStack = m44.Pop(matrixStack, 1)
		} else {
			var nodeIndex int = len(nodes)
			nodes = append(nodes, node)
			nodesMap[node] = int32(nodeIndex)
			//if (options & MODEL_ANIMATION) != 0 {
			var skeletonRoot int32 = g_models[modelIndex].skeleton.root
			if skeletonRoot == -1 {
				var skin gltf.Skin = gltf.NodeGetSkin(asset, node)
				if skin.InverseBindMatrices != -1 {
					var accessor gltf.Accessor = gltf.AssetGetAccessor(asset, skin.InverseBindMatrices)
					utils.PanicIfNot(accessor.Loaded, "inverse bind mantrines accessor is not loaded")
					var inverseBinds []math.M44 = g_models[modelIndex].inverseBinds
					utils.PanicIf(len(accessor.DataF32)%16 != 0, "invalid matrix")

					var acc []float32 = accessor.DataF32
					var acci int32 = 0
					var inverseBindCount int = len(accessor.DataF32) / 16
					for i := 0; i < inverseBindCount; i++ {
						var ib math.M44
						ib.V00 = acc[acci]
						acci++
						ib.V01 = acc[acci]
						acci++
						ib.V02 = acc[acci]
						acci++
						ib.V03 = acc[acci]
						acci++
						ib.V10 = acc[acci]
						acci++
						ib.V11 = acc[acci]
						acci++
						ib.V12 = acc[acci]
						acci++
						ib.V13 = acc[acci]
						acci++
						ib.V20 = acc[acci]
						acci++
						ib.V21 = acc[acci]
						acci++
						ib.V22 = acc[acci]
						acci++
						ib.V23 = acc[acci]
						acci++
						ib.V30 = acc[acci]
						acci++
						ib.V31 = acc[acci]
						acci++
						ib.V32 = acc[acci]
						acci++
						ib.V33 = acc[acci]
						acci++
						inverseBinds = append(inverseBinds, ib)

					}
					g_models[modelIndex].inverseBinds = inverseBinds
					g_models[modelIndex].jointMatrices = make([]math.M44, inverseBindCount)
				}
				g_models[modelIndex].skeleton.root = skin.Skeleton
				g_models[modelIndex].skinned = skin.Skeleton >= 0 // TODO : multiple skeleton per scene
				g_models[modelIndex].joints = skin.Joints
			}
			//}

			var nodeMatrix math.M44 = gltf.NodeGetMatrix(asset, node)
			var determinant float32 = m44.Determinant(nodeMatrix)
			determinants = append(determinants, determinant)
			// var cullFace uint32 = gl.BACK
			// if determinant < 0.0 {
			// 	cullFace = gl.FRONT
			// }
			matrixStack = m44.Push(matrixStack, nodeMatrix)

			var last math.M44 = matrixStack[len(matrixStack)-1]
			transforms = append(transforms, last)

			if (options & MODEL_GEOMETRY) != 0 {
				var mesh int32 = gltf.NodeGetMesh(asset, node)
				var gltfPrimitives []gltf.Primitive = gltf.MeshGetPrimitives(asset, mesh)
				var gltfPrimitiveCount int = len(gltfPrimitives)

				for p := 0; p < gltfPrimitiveCount; p++ {

					var render MeshId = InvalidMesh()

					var gltfPrimitive gltf.Primitive = gltfPrimitives[p]
					var gltfMaterial gltf.Material = gltf.AssetGetMaterial(asset, gltfPrimitive.Material)
					var metallicRoughness gltf.MetallicRoughness = gltfMaterial.PbrMetallicRoughness

					var indicesIndex int32 = gltfPrimitive.Indices
					var meshCacheCount int = len(meshCacheKey)
					var pp int
					for pp = 0; pp < meshCacheCount; pp++ {
						if meshCacheKey[pp] == indicesIndex {
							render = meshCacheVal[pp]
							pp = meshCacheCount
						}
					}

					var cullFace uint32 = gl.BACK
					if gltfMaterial.DoubleSided == 1 {
						cullFace = gl.NONE
					} else if determinant < 0.0 {
						cullFace = gl.FRONT
					}

					if IsValidMesh(render) == false {
						render = meshInstance(asset, gltfPrimitive)
						MeshSetCulling(render, gl.CCW, cullFace)
						meshCacheKey = append(meshCacheKey, indicesIndex)
						meshCacheVal = append(meshCacheVal, render)
					}

					if IsValidMesh(render) == false {
						success = false
						return
					}

					var primitive Primitive
					var meshMin math.V3 = g_meshes[render.mesh].min
					var meshMax math.V3 = g_meshes[render.mesh].max

					primitive.min = meshMin
					primitive.max = meshMax

					meshMin = v3.Transform_point(meshMin, last)
					meshMax = v3.Transform_point(meshMax, last)

					max = v3.Max(meshMax, max)
					min = v3.Min(meshMin, min)

					primitive.mesh = render
					primitive.node = node
					primitive.nodeIndex = int32(nodeIndex)
					primitive.baseTexture = loadTexture(asset, dataDir, metallicRoughness.BaseColorTexture.Index)
					primitive.metalRoughTexture = loadTexture(asset, dataDir, metallicRoughness.MetallicRoughnessTexture.Index)
					primitive.emissiveTexture = loadTexture(asset, dataDir, gltfMaterial.EmissiveTexture.Index)
					primitive.normalTexture = loadTexture(asset, dataDir, gltfMaterial.NormalTexture.Index)
					primitive.occlusionTexture = loadTexture(asset, dataDir, gltfMaterial.OcclusionTexture.Index)
					primitive.gltfMaterial = gltfMaterial
					primitive.gltfPrimitive = gltfPrimitive
					if MeshUseWeight(render) && MeshUseJoint(render) {
						primitive.useSkin = 1
					}

					var effect EffectId = effectInstance(primitive)
					if EffectIsValid(effect) == false {
						fmt.Printf("effectInstance failed\n")
						success = false
						return
					}
					primitive.effect = effect

					var primitiveIndex int32 = int32(len(primitives))
					primitives = append(primitives, primitive)

					meshes = append(meshes, primitiveIndex)
					var alphaMode int32 = gltfMaterial.AlphaMode
					if alphaMode == gltf.ALPHA_BLEND {
						transparentMeshes = append(transparentMeshes, primitiveIndex)
					} else if alphaMode == gltf.ALPHA_OPAQUE {
						opaqueMeshes = append(opaqueMeshes, primitiveIndex)
					}
				}
			}

			var children []int32 = gltf.NodeGetChildren(asset, node)
			var childCount int = len(children)
			stack = append(stack, -1)
			if childCount > 0 {
				for i := 0; i < childCount; i++ {
					stack = append(stack, children[i])
				}
			}
		}
	}

	stack = stack[:0]
	g_models[modelIndex].stack = stack
	g_models[modelIndex].nodes = nodes
	g_models[modelIndex].nodesMap = nodesMap
	g_models[modelIndex].determinants = determinants
	g_models[modelIndex].meshes = meshes
	g_models[modelIndex].primitives = primitives
	g_models[modelIndex].opaqueMeshes = sortMeshes(primitives, opaqueMeshes, g_models[modelIndex].opaqueMeshes)
	g_models[modelIndex].transparentMeshes = sortMeshes(primitives, transparentMeshes, g_models[modelIndex].transparentMeshes)
	g_models[modelIndex].meshCacheKey = meshCacheKey
	g_models[modelIndex].meshCacheVal = meshCacheVal
	g_models[modelIndex].animations = animations
	g_models[modelIndex].matrixStack = matrixStack[:0]
	g_models[modelIndex].transforms = transforms
	g_models[modelIndex].max = max
	g_models[modelIndex].min = min

	success = true
	return
}

func AnimationMorph(left AnimationId, right AnimationId, out AnimationId, time float32) {
	utils.PanicIfNot(AnimationIsValid(left), "invalid animation")
	utils.PanicIfNot(AnimationIsValid(right), "invalid animation")
	utils.PanicIfNot(AnimationIsValid(out), "invalid animation")

	var leftIndex int32 = left.animation
	var trsLeft []math.TRS = g_animations[leftIndex].trs
	var leftSources []AnimationSource = g_animations[leftIndex].sources

	var rightIndex int32 = right.animation
	var trsRight []math.TRS = g_animations[rightIndex].trs
	var rightSources []AnimationSource = g_animations[rightIndex].sources

	var outIndex int32 = out.animation
	var trsOut []math.TRS = g_animations[outIndex].trs

	var count int = len(trsLeft)
	utils.PanicIf(count != len(trsRight), "incompatible animations")

	trsOut = trsOut[:0]
	var tt math.TRS
	for i := 0; i < count; i++ {
		tt.T = v3.Lerpf(trsLeft[i].T, trsRight[i].T, time)
		tt.R = q4.Slerp(trsLeft[i].R, trsRight[i].R, time)
		tt.S = v3.Lerpf(trsLeft[i].S, trsRight[i].S, time)
		trsOut = append(trsOut, tt)
	}

	g_animations[outIndex].min = math.Min_f32(g_animations[leftIndex].min, g_animations[rightIndex].min)
	g_animations[outIndex].max = math.Max_f32(g_animations[leftIndex].max, g_animations[rightIndex].max)
	g_animations[outIndex].trs = trsOut
	if len(rightSources) > 0 {
		g_animations[outIndex].sources = rightSources
	} else if len(leftSources) > 0 {
		g_animations[outIndex].sources = leftSources
	}
}

// // AnimationGetMax ...
// func AnimationGetMax(id AnimationId) (out float32) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid model")
// 	out = g_animations[id.animation].max
// }

// // AnimationSetTime ...
// func AnimationSetTime(id AnimationId, time float32) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	g_animations[id.animation].time = time
// 	g_animations[id.animation].lastIndex = 0
// }

// // AnimationGetRuntime ...
// func AnimationGetRuntime(id AnimationId) (out bool) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	out = g_animations[id.animation].runtime
// }

// // AnimationSetRuntime ...
// func AnimationSetRuntime(id AnimationId, runtime bool) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	g_animations[id.animation].runtime = runtime
// }

// // AnimationGetJointTRS ...
// func AnimationGetJointTRS(id AnimationId, joint int32) (out TRS) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	var trs []math.TRS = g_animations[id.animation].trs
// 	out = trs[joint]
// }

// // AnimationSetJointTRS ...
// func AnimationSetJointTRS(id AnimationId, joint int32, t math.V3, r math.V4, s math.V3) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	var trs[]math.TRS = g_animations[id.animation].trs
// 	trs[joint].t = t
// 	trs[joint].r = r
// 	trs[joint].S = s
// }

// // AnimationGetJointRotation ...
// func AnimationGetJointRotation(id AnimationId, joint int32) (out math.V4) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	var trs []math.TRS = g_animations[id.animation].trs
// 	out = trs[joint].r
// }

// // AnimationSetJointRotation ...
// func AnimationSetJointRotation(id AnimationId, joint int32, q math.V4) {
// 	utils.PanicIfNot(AnimationIsValid(id), "invalid animation")
// 	var trs[]math.TRS = g_animations[id.animation].trs
// 	trs[joint].r = q
// }

// AnimationUpdate ...
func AnimationUpdate(id AnimationId, deltaTime float32, loop bool, max float32) {
	utils.PanicIfNot(AnimationIsValid(id), "invalid model")
	var animIndex int32 = id.animation
	// var animationCount int32 = int32(len(g_animations))

	var animation Animation = g_animations[animIndex]
	var lastIndex int32 = animation.lastIndex
	var animMin float32 = animation.min
	var animMax float32 = animation.max
	if max > 0.0 {
		animMax = max
	}

	var keys []float32 = animation.keys
	var keyCount int32 = int32(len(keys))

	var animDir float32 = animation.direction
	var animTime float32 = animation.time + animDir*deltaTime
	if animTime < animMin {
		if loop == true {
			animTime = animMax
			lastIndex = keyCount - 1
		} else {
			animTime = animMin
			animDir = 1.0
		}
	}
	if animTime > animMax {
		if loop == true {
			animTime = 0.0 //animMin
			lastIndex = 0
		} else {
			animTime = animMax
			animDir = -1.0
		}
	}

	var index int32 = -1
	for keyIndex := lastIndex; keyIndex < keyCount; keyIndex++ { // TODO : track time
		if animTime <= keys[keyIndex] {
			index = keyIndex
			keyIndex = keyCount // ISSUE : break
		}
	}

	if index >= keyCount {
		index = keyCount - 1
	}
	if index < 0 {
		index = 0
	}

	var prev int32 = index - 1
	if prev < 0 {
		prev = 0
	}

	var delta float32 = animTime - keys[prev]
	var step float32 = keys[index] - keys[prev]
	var st float32 = 1.0
	if step > 0.0 {
		st = delta / step
	}
	var t float32 = st
	if t < 0.0 {
		t = 0.0
	} else if t > 1.0 {
		t = 1.0
	}

	animation.lastIndex = prev
	var sources []AnimationSource = animation.sources
	var sourceCount int32 = int32(len(sources))
	if sourceCount > 0 {
		animation.time = animTime
		animation.direction = animDir

		var animTrs []math.TRS = animation.trs
		animTrs = animTrs[:0]

		var frame math.TRS
		for sourceIndex := int32(0); sourceIndex < sourceCount; sourceIndex++ {
			var source AnimationSource = sources[sourceIndex]
			if source.anim {
				var prevFrame math.TRS = source.frames[prev]
				var curFrame math.TRS = source.frames[index]
				frame.T = v3.Lerpsatf(prevFrame.T, curFrame.T, t)
				frame.R = q4.Slerp(prevFrame.R, curFrame.R, t)
				frame.S = v3.Lerpsatf(prevFrame.S, curFrame.S, t)
			} else {
				frame.T = v3.ZERO
				frame.R = v4.ALPHA
				frame.S = v3.ONE
			}
			animTrs = append(animTrs, frame)
		}

		animation.trs = animTrs
	}
	g_animations[animIndex] = animation
}

// func trs_to_str(t TRS) (out string) {
// 	out = fmt.Sprintf("%s, %s, %s", v3.to_str(t.t), v4.to_str(t.r), v3.to_str(t.S))
// }

// func debugNode(stackCount int32, asset gltf.AssetId, node int32) {
// 	for kk := 0; kk < stackCount; kk++ {
// 		fmt.Printf("--")
// 	}

// 	fmt.Printf("%s :\n", gltf.NodeGetName(asset, node))
// }

func modelTransform(id ModelId, animId AnimationId, deltaTime float32, loop bool) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")

	var modelIndex int32 = id.model
	var animIndex int32 = animId.animation

	var animationCount int32 = int32(len(g_animations))
	if animIndex >= 0 && animIndex < animationCount {
		var animated bool = g_models[modelIndex].animated
		var skinned bool = g_models[modelIndex].skinned
		if animated || skinned {

			var animation Animation = g_animations[animIndex]
			var animTrs []math.TRS = animation.trs
			var animTransforms []math.M44 = animation.transforms
			animTransforms = animTransforms[:0]
			var trsCount int32 = int32(len(animTrs))
			var sources []AnimationSource = animation.sources

			if trsCount > 0 {
				var modelAsset gltf.AssetId = g_models[modelIndex].asset
				for trsIndex := int32(0); trsIndex < trsCount; trsIndex++ {
					if sources[trsIndex].anim {
						animTransforms = append(animTransforms, m44.Make_SQT(animTrs[trsIndex]))
					} else {
						animTransforms = append(animTransforms, sources[trsIndex].transform)
					}
				}
				var animTransformCount int32 = int32(len(animTransforms))

				g_animations[animIndex].transforms = animTransforms

				var transforms []math.M44 = g_models[modelIndex].transforms
				transforms = transforms[:0]

				var matrixStack []math.M44 = g_models[modelIndex].matrixStack
				matrixStack = matrixStack[:0]
				matrixStack = m44.Push(matrixStack, m44.IDENTITY)

				var stack []int32 = g_models[modelIndex].stack

				var roots []int32 = gltf.AssetGetRootNodes(modelAsset)
				var rootCount int32 = int32(len(roots))
				stack = make([]int32, rootCount)
				for i := 0; i < int(rootCount); i++ {
					stack[i] = roots[i]
				}

				for len(stack) > 0 {
					var stackLen int32 = int32(len(stack))

					stackLen--
					var node int32 = stack[stackLen]

					stack = stack[:stackLen]
					if node < 0 {
						matrixStack = m44.Pop(matrixStack, 1)
					} else {
						if node < animTransformCount {
							matrixStack = m44.Push(matrixStack, animTransforms[node])
						} else {
							matrixStack = m44.Push(matrixStack, sources[node].transform)
						}

						transforms = append(transforms, matrixStack[len(matrixStack)-1])

						var children []int32 = gltf.NodeGetChildren(modelAsset, node)
						var childCount int32 = int32(len(children))
						stack = append(stack, -1)
						if childCount > 0 {
							for i := 0; i < int(childCount); i++ {
								stack = append(stack, children[i])
							}
						}
					}
				}

				g_models[modelIndex].stack = stack[:0]
				g_models[modelIndex].matrixStack = matrixStack[:0]
				g_models[modelIndex].transforms = transforms
			}
		}
	}

	var joints []int32 = g_models[modelIndex].joints
	var jointCount int32 = int32(len(joints))
	if jointCount > 0 {
		var jointMatrices []math.M44 = g_models[modelIndex].jointMatrices
		var inverseBinds []math.M44 = g_models[modelIndex].inverseBinds
		var nodesMap []int32 = g_models[modelIndex].nodesMap
		var transforms []math.M44 = g_models[modelIndex].transforms
		for j := 0; j < int(jointCount); j++ {
			var node int32 = joints[j]
			var nodeIndex int32 = nodesMap[node]
			jointMatrices[j] = m44.MulISSUE(inverseBinds[j], transforms[nodeIndex])
		}
	}
}

// // ModelGetMeshes ...
// func ModelGetMeshes(id ModelId) (out []int32) {
// 	utils.PanicIfNot(ModelIsValid(id), "invalid model")
// 	out = g_models[id.model].meshes
// }

// // ModelGetNodePosition ...
// func ModelGetNodePosition(id ModelId, node int32) (out math.V3) {
// 	utils.PanicIfNot(ModelIsValid(id), "invalid model")
// 	var modelIndex int32 = id.model
// 	var nodesMap []int32 = g_models[modelIndex].nodesMap
// 	var nodeIndex int32 = nodesMap[node]
// 	var transforms []m44 = g_models[modelIndex].transforms
// 	out.X = transforms[nodeIndex].V30
// 	out.Y = transforms[nodeIndex].V31
// 	out.Z = transforms[nodeIndex].V32
// }

// // ModelGetNodeTransform ...
// func ModelGetNodeTransform(id ModelId, node int32) (out m44) {
// 	utils.PanicIfNot(ModelIsValid(id), "invalid model")
// 	var modelIndex int32 = id.model
// 	var nodesMap []int32 = g_models[modelIndex].nodesMap
// 	var nodeIndex int32 = nodesMap[node]
// 	var transforms []m44 = g_models[modelIndex].transforms
// 	out = transforms[nodeIndex]
// }

func sortMeshes(primitives []Primitive, meshes []int32, in []int32) (out []int32) {
	out = in
	var meshCount int = len(meshes)
	for i := 0; i < meshCount; i++ {
		var mesh int32 = meshes[i]
		var p Primitive = primitives[mesh]
		var sortedCount int = len(out)
		for si := 0; si < sortedCount; si++ {
			var smesh int32 = out[si]
			var sp Primitive = primitives[smesh]
			var in bool
			if p.effect.effect < sp.effect.effect {
				in = true
			} else if p.effect.effect == sp.effect.effect {
				if p.baseTexture.texture < sp.baseTexture.texture {
					in = true
				} else if p.baseTexture.texture == sp.baseTexture.texture {
					if p.normalTexture.texture < sp.normalTexture.texture {
						in = true
					} else if p.normalTexture.texture == sp.normalTexture.texture {
						if p.metalRoughTexture.texture < sp.metalRoughTexture.texture {
							in = true
						} else if p.metalRoughTexture.texture == sp.metalRoughTexture.texture {
							if p.emissiveTexture.texture < sp.emissiveTexture.texture {
								in = true
							} else if p.emissiveTexture.texture == sp.emissiveTexture.texture {
								if p.occlusionTexture.texture < sp.occlusionTexture.texture {
									in = true
								}
							}
						}
					}
				}
			}
			if in {
				out = append(out[:si+1], out[si:]...)
				out[si] = mesh
				si = sortedCount
			}
		}
		if sortedCount == len(out) {
			out = append(out, mesh)
		}
	}
	return
}

func modelRender(id ModelId, renderables []int32, alpha float32) {
	utils.PanicIfNot(ModelIsValid(id), "invalid model")

	// TODO REMOVE var meshRenderCount int32
	var modelIndex int32 = id.model
	var transforms []math.M44 = g_models[modelIndex].transforms
	var primitives []Primitive = g_models[modelIndex].primitives

	var environmentSpecular TextureId = g_models[modelIndex].environmentSpecular // TODO : renderState
	var environmentDiffuse TextureId = g_models[modelIndex].environmentDiffuse   // TODO : renderState
	var brdf TextureId = g_models[modelIndex].brdf                               // TODO : renderState
	var view math.M44 = g_models[modelIndex].view                                // TODO : renderState
	var projection math.M44 = g_models[modelIndex].projection                    // TODO : renderState
	var cameraPosition math.V4 = g_models[modelIndex].cameraPosition             // TODO : renderState
	var jointMatrices []math.M44 = g_models[modelIndex].jointMatrices
	var skinned bool = g_models[modelIndex].skinned
	var world math.M44 = g_models[modelIndex].world
	var worldIsIdentity bool = g_models[modelIndex].worldIsIdentity
	var exposure float32 = g_models[modelIndex].exposure

	var renderableCount int32 = int32(len(renderables))

	var lastEffect int32 = -1
	var lastBaseTexture int32 = -1
	var lastNormalTexture int32 = -1
	var lastMetalRoughTexture int32 = -1
	var lastEmissiveTexture int32 = -1
	var lastOcclusionTexture int32 = -1
	var lastBaseColorFactor math.V4 = v4.Makef(-1.0)
	var lastEmissiveFactor math.V4 = v4.Makef(-1.0)
	var lastMetalRoughFactor math.V4 = v4.Makef(-1.0)
	var lastTransform math.M44 = m44.INVALID

	for p := 0; p < int(renderableCount); p++ {

		var renderable int32 = renderables[p]
		var primitive Primitive = primitives[renderable]
		var effect EffectId = primitive.effect
		EffectUse(effect)

		if lastEffect != effect.effect {
			lastEffect = effect.effect

			EffectTryAssignTexture(effect, SAMPLER_ENV_SPECULAR, environmentSpecular, g_linearClamp) // TODO : use gltf sampler
			EffectTryAssignTexture(effect, SAMPLER_ENV_DIFFUSE, environmentDiffuse, SpLinear0Clamp)  // TODO : use gltf sampler
			EffectTryAssignTexture(effect, SAMPLER_BRDF, brdf, SpLinear0Clamp)                       // TODO : use gltf sampler

			EffectAssignM44(effect, UNIFORM_VIEW, view, false)
			EffectAssignM44(effect, UNIFORM_PROJECTION, projection, false)

			EffectAssignV4(effect, UNIFORM_CAMERA_POSITION, cameraPosition)
			EffectAssignV4(effect, UNIFORM_DEBUG_0, DEBUG_0)
			EffectAssignV4(effect, UNIFORM_PBR, v4.Make(float32(TextureGetMipmapCount(environmentSpecular)),
				exposure, 0.0, 0.0))
			if primitive.useSkin == 1 {
				EffectAssignM44V(effect, UNIFORM_SKELETON, jointMatrices, false)
			}

			lastBaseTexture = -1
			lastNormalTexture = -1
			lastMetalRoughTexture = -1
			lastEmissiveTexture = -1
			lastOcclusionTexture = -1
			lastBaseColorFactor = v4.Makef(-1.0)
			lastEmissiveFactor = v4.Makef(-1.0)
			lastMetalRoughFactor = v4.Makef(-1.0)
			lastTransform = m44.INVALID
		}

		if lastBaseTexture != primitive.baseTexture.texture {
			lastBaseTexture = primitive.baseTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_COLOR_0, primitive.baseTexture, SpLinearWrap) // TODO : use gltf sampler
		}

		if lastNormalTexture != primitive.normalTexture.texture {
			lastNormalTexture = primitive.normalTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_NORMAL, primitive.normalTexture, SpLinearWrap) // TODO : use gltf sampler
		}

		if lastMetalRoughTexture != primitive.metalRoughTexture.texture {
			lastMetalRoughTexture = primitive.metalRoughTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_METAL_ROUGH, primitive.metalRoughTexture, SpLinearWrap) // TODO : use gltf sampler
		}

		if lastEmissiveTexture != primitive.emissiveTexture.texture {
			lastEmissiveTexture = primitive.emissiveTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_EMISSIVE, primitive.emissiveTexture, SpLinearWrap) // TODO : use gltf sampler
		}

		if lastOcclusionTexture != primitive.occlusionTexture.texture {
			lastOcclusionTexture = primitive.occlusionTexture.texture
			EffectTryAssignTexture(effect, SAMPLER_OCCLUSION, primitive.occlusionTexture, SpLinearWrap) // TODO : use gltf sampler
		}

		var tmpWorld math.M44 = world
		if skinned == false {
			if worldIsIdentity == false {
				tmpWorld = m44.MulISSUE(transforms[primitive.nodeIndex], world)
			} else {
				tmpWorld = transforms[primitive.nodeIndex]
			}
		}

		if m44.Equ(tmpWorld, lastTransform) == false {
			lastTransform = tmpWorld
			EffectAssignM44(effect, UNIFORM_WORLD, tmpWorld, false)
		}
		//EffectAssignM44(effect, UNIFORM_WORLD_INVERSE, m44.inverse(tmpWorld), false) // TODO : optim inverse sooner in callstack

		var metallicRoughness gltf.MetallicRoughness = primitive.gltfMaterial.PbrMetallicRoughness
		if EffectIsValidUniformLocation(effect, UNIFORM_COLOR) {
			var baseColorFactor math.V4 = metallicRoughness.BaseColorFactor
			baseColorFactor.W = baseColorFactor.W * alpha

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
				/*fmt.Printf("METAL_ROUGH %f, %f, %f, %f\n",
				  metallicRoughnessFactor.X,
				  metallicRoughnessFactor.Y,
				  metallicRoughnessFactor.Z,
				  metallicRoughnessFactor.W)*/

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

		MeshRender(primitive.mesh)
	}
}
