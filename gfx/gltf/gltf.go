package gltf

import (
	"fmt"
	"skyfx/math"
	"skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"
	"skyfx/utils"
	"skyfx/utils/json"

	"golang.org/x/mobile/gl"
)

// import "json"
// import "os"
// import "gl"
// import "mat"
// import "v1"
// import "v3"
// import "v4"
// import "m44"

// Constants ...
const (
	HEADER_SIZE       int32  = 12
	CHUNK_HEADER_SIZE int32  = 8
	VERSION           uint32 = 2
	MAGIC             uint32 = 1179937895 // 0x46546C67
	MAGIC_JSON_CHUNK  uint32 = 1313821514 // 0x4E4F534A
	MAGIC_BIN_CHUNK   uint32 = 5130562    // 0x004E4942
)

const (
	ASSET_GEOMETRY  int32 = 1
	ASSET_ANIMATION int32 = 2

// /*var FILE_INVALID int32 = 0
// var FILE_GLTF int32 = 1
// var FILE_GLB int32 = 2
)

// var BUFFER_VIEW_INVALID int32 = 0
// var BUFFER_VIEW_INDICES int32 = 1
// var BUFFER_VIEW_VERTIECS int32 = 2
// */
const (
	ATTRIBUTE_INVALID  int32 = 0
	ATTRIBUTE_POSITION int32 = 1
	ATTRIBUTE_NORMAL   int32 = 2
	ATTRIBUTE_TANGENT  int32 = 3
	ATTRIBUTE_TEXCOORD int32 = 4
	ATTRIBUTE_COLOR    int32 = 5
	ATTRIBUTE_JOINT    int32 = 6
	ATTRIBUTE_WEIGHT   int32 = 7
)

// /*
// var COMPONENT_INVALID int32 = 0
// var COMPONENT_R8 int32 = 1
// var COMPONENT_R8U int32 = 2
// var COMPONENT_R16 int32 = 3
// var COMPONENT_R16U int32 = 4
// var COMPONENT_R32U int32 = 5
// var COMPONENT_R32F int32 = 6*/

const (
	TYPE_INVALID int32 = 0
	TYPE_SCALAR  int32 = 1
	TYPE_VEC2    int32 = 2
	TYPE_VEC3    int32 = 3
	TYPE_VEC4    int32 = 4
	TYPE_MAT2    int32 = 5
	TYPE_MAT3    int32 = 6
	TYPE_MAT4    int32 = 7
)

const (
	PRIMITIVE_POINTS         int32 = 0
	PRIMITIVE_LINES          int32 = 1
	PRIMITIVE_LINE_LOOP      int32 = 2
	PRIMITIVE_LINE_STRIP     int32 = 3
	PRIMITIVE_TRIANGLES      int32 = 4
	PRIMITIVE_TRIANGLE_STRIP int32 = 5
	PRIMIITVE_TRIANGLE_FAN   int32 = 6
)

const (
	ALPHA_OPAQUE int32 = 0
	ALPHA_MASK   int32 = 1
	ALPHA_BLEND  int32 = 2
)

const (
	ANIMATION_PATH_INVALID     int32 = 0
	ANIMATION_PATH_TRANSLATION int32 = 1
	ANIMATION_PATH_ROTATION    int32 = 2
	ANIMATION_PATH_SCALE       int32 = 3
	ANIMATION_PATH_WEIGHTS     int32 = 4
)

const (
	INTERPOLATION_LINEAR       int32 = 0
	INTERPOLATION_STEP         int32 = 1
	INTERPOLATION_CUBIC_SPLINE int32 = 2
)

const (
	// var CAMERA_INVALID int32 = 0
	CAMERA_PERSPECTIVE int32 = 1

// var CAMERA_ORTHOGRAPHIC int32 = 2
// /*
// var LIGHT_INVALID int32 = 0
// var LIGHT_DIRECTIONAL int32 = 1
// var LIGHT_POINT int32 = 2
// var LIGHT_SPOT int32 = 3
// */
)

// // Globals ...
var g_assets []Asset
var CurrentAsset AssetId // TODO : remove

type Buffer struct {
	uri        string
	byteLength int32
	name       string
}

type BufferView struct {
	buffer     int32
	byteOffset int64
	byteLength int32
	byteStride int32
	target     int32
	name       string
}

// /*
// type AccessorSparse struct {
// 	count int64
// 	indices[]int32
// 	values []int32
// }
// */

type Accessor struct {
	bufferView     int32
	byteOffset     int64
	ComponentType  int32
	ComponentCount uint32
	//normalized bool
	AttributeTypeStr  string
	AttributeTypeEnum int32
	Count             int32
	Max               []float32
	Min               []float32
	//sparse int32
	name     string
	DataF32  []float32
	DataUI16 []uint16
	DataUI32 []uint32
	Loaded   bool
}

type Attribute struct {
	name          string
	AttributeType int32
	index         int32
	Accessor      int32
}

type Texture struct {
	sampler int32
	Source  int32
	name    string
}

type Image struct {
	Uri        string
	mimeType   string
	bufferView int32
	name       string
}

type Sampler struct {
	magFilter int32
	minFilter int32
	wrapS     int32
	wrapT     int32
	name      string
}

// /*
// type TextureTransform struct {
// 	offset [2]float32
// 	rotation float32
// 	scale [2]float32
// 	texcoord int32
// }
// */

type TextureView struct {
	Index    int32
	texCoord int32
	Scale    float32
	//hasTransform bool
	//transform int32
}

type MetallicRoughness struct {
	BaseColorFactor          math.V4
	BaseColorTexture         TextureView
	MetallicFactor           float32
	RoughnessFactor          float32
	MetallicRoughnessTexture TextureView
}

type SpecularGlossiness struct {
	DiffuseFactor             math.V4
	DiffuseTexture            TextureView
	GlossinessFactor          float32
	SpecularFactor            math.V3
	SpecularGlossinessTexture TextureView
}

type Material struct {
	name                  string
	PbrMetallicRoughness  MetallicRoughness
	PbrSpecularGlossiness SpecularGlossiness
	NormalTexture         TextureView
	OcclusionTexture      TextureView
	EmissiveTexture       TextureView
	EmissiveFactor        math.V4
	AlphaMode             int32
	//alphaCutoff float32
	DoubleSided int32 //bool
	//KHR_extension_unlit bool
}

// /*
// type MorphTarget struct {
// 	attributes []int32
// }
// */

type Primitive struct {
	Attributes []Attribute
	Indices    int32
	Material   int32
	Mode       int32
	//targets []int32
}

type Mesh struct {
	primitives []Primitive
	weights    []float32
	name       string
}

type Skin struct {
	InverseBindMatrices int32
	Skeleton            int32
	Joints              []int32
	name                string
}

type PerspectiveCamera struct {
	aspectRatio float32
	yfov        float32
	znear       float32
	zfar        float32
}

type Camera struct {
	name              string
	cameraTypeStr     string
	cameraTypeEnum    int32
	perspectiveCamera PerspectiveCamera
}

// /*
// type Light struct {
// 	lightType int32
// 	name string
// 	color [3]float32
// 	intensity float32
// 	range int32
// 	innerConeAngle float32
// 	outerConeAngle float32
// }
// */

type AnimationSampler struct {
	interpolation int32
	InputA        int32
	OutputA       int32
}

type AnimationChannel struct {
	Sampler int32
	Node    int32
	Path    int32
}

type Animation struct {
	Channels []AnimationChannel
	Samplers []AnimationSampler
	name     string
}

type Node struct {
	children    []int32
	name        string
	camera      int32
	skin        int32
	mesh        int32
	matrix      math.M44
	rotation    math.V4
	scale       math.V3
	translation math.V3
	//weights []float32
	//parent int32
	//light int32
	hasMatrix      bool
	hasScale       bool
	hasRotation    bool
	hasTranslation bool
}

// Scene
type Scene struct {
	name  string
	nodes []int32
}

// AssetId ...
type AssetId struct {
	asset int32
}

// Asset ...
type Asset struct {
	accessors   []Accessor
	animations  []Animation
	buffers     []Buffer
	bufferViews []BufferView
	cameras     []Camera
	images      []Image
	materials   []Material
	meshes      []Mesh
	nodes       []Node
	samplers    []Sampler
	scene       int32
	scenes      []Scene
	skins       []Skin
	textures    []Texture
	//extensionsUsed []string
	//extensionsRequired []string

	path string

	author     string
	license    string
	source     string
	title      string
	copyright  string
	generator  string
	version    string
	minVersion string

	//lights []Light

}

// AssetInvalid ...
func AssetInvalid() (out AssetId) {
	out.asset = -1
	return
}

// AssetIsValid ...
func AssetIsValid(id AssetId) (out bool) {
	out = id.asset >= 0 && id.asset <= int32(len(g_assets))
	return
}

func alwaysMore(file int32) (out bool) {
	var more bool
	var success bool
	more, success = json.TokenMore(file)
	out = more && success
	return
}

func readM44F(file int32, res *math.M44) (success bool) {
	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool
	var value float32

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V00 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V01 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V02 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V03 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V10 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V11 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V12 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V13 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V20 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V21 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V22 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V23 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V30 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V31 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V32 = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).V33 = value
	success = json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT)
	return
}

func readV4F(file int32, res *math.V4) (success bool) {
	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool
	var value float32

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).X = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).Y = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).Z = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).W = value
	success = json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT)
	return
}

func readV3F(file int32, res *math.V3) (success bool) {
	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool
	var value float32

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).X = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).Y = value

	more, success = json.TokenMore(file)
	if more == false || success == false {
		success = false
		return
	}
	if json.ReadF32(file, &value) == false {
		success = false
		return
	}
	(*res).Z = value
	success = json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT)
	return
}

// //------------------------------------------------------------------------------
// func printI32Array(message string, array []int32) {
// 	var count int32 = len(array)
// 	fmt.Printf("%s, count : %d [", message, count)
// 	for n := 0; n < count; n++ {
// 		fmt.Printf("%d, ", array[n])
// 	}
// 	fmt.Printf("]\n")
// }

// //------------------------------------------------------------------------------
// func printF32Array(message string, array []float32) {
// 	var count int32 = len(array)
// 	fmt.Printf("%s, count : %d [", message, count)
// 	for n := 0; n < count; n++ {
// 		fmt.Printf("%f, ", array[n])
// 	}
// 	fmt.Printf("]\n")
// }

// // AssetPrint ...
// /*func AssetPrint(id AssetId) { // TODO : format as json
// 	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
// 	fmt.Printf("asset : %s\n", g_assets[id.asset].path)
// 	fmt.Printf("---- copyright : %s\n", g_assets[id.asset].copyright)
// 	fmt.Printf("---- generator : %s\n", g_assets[id.asset].generator)
// 	fmt.Printf("---- version : %s\n", g_assets[id.asset].version)
// 	fmt.Printf("---- minVersion : %s\n", g_assets[id.asset].minVersion)
// 	fmt.Printf("---- scene : %d\n", g_assets[id.asset].scene.scene)
// 	var sceneCount int32 = len(g_assets[id.asset].scenes)
// 	fmt.Printf("scenes, count : %d [\n", sceneCount)
// 	for i := 0; i < sceneCount; i++ {
// 		fmt.Printf("---- scene %d, name '%s'\n", i, g_assets[id.asset].scenes[i].name)
// 		printI32Array("---- ---- nodes", g_assets[id.asset].scenes[i].nodes)
// 	}
// 	fmt.Printf("]\n")

// 	var nodes []Node = g_assets[id.asset].nodes
// 	var nodeCount int32 = len(nodes)
// 	fmt.Printf("nodes, count : %d [\n", nodeCount)
// 	for i := 0; i < nodeCount; i++ {
// 		fmt.Printf("--- node %d, name '%s'\n", i, nodes[i].name)
// 		fmt.Printf("---- ---- mesh %d\n", nodes[i].mesh)
// 		var children []int32 = nodes[i].children
// 		printI32Array("---- ---- children", g_assets[id.asset].nodes[i].children)
// 		var matrix []float32 = nodes[i].matrix
// 		printF32Array("---- ---- matrix", nodes[i].matrix)
// 		fmt.Printf("---- ---- camera %d\n", nodes[i].camera)
// 	}
// 	fmt.Printf("]\n")

// 	var meshes []Mesh = g_assets[id.asset].meshes
// 	var meshCount int32 = len(meshes)
// 	fmt.Printf("meshes, count : %d [\n", meshCount)
// 	for i := 0; i < meshCount; i++ {
// 		fmt.Printf("---- mesh %d, name '%s'\n", i, meshes[i].name)
// 		var primitives []Primitive = meshes[i].primitives
// 		var primitiveCount int32 = len(primitives)
// 		fmt.Printf("---- ---- primitive, count %d\n", primitiveCount)
// 		for p := 0; p < primitiveCount; p++ {
// 			fmt.Printf("---- ---- primitive %d\n", p)
// 			fmt.Printf("---- ---- ---- mode %d\n", primitives[p].mode)
// 			fmt.Printf("---- ---- ---- indices %d\n", primitives[p].indices)
// 			fmt.Printf("---- ---- ---- material %d\n", primitives[p].material)
// 			var attributes []Attribute = primitives[p].attributes
// 			var attributeCount int32 = len(attributes)
// 			fmt.Printf("---- ---- ---- attributes, count %d\n", attributeCount)
// 			for a := 0; a < attributeCount; a++ {
// 				fmt.Printf("---- ---- ---- attribute %d, name '%s'\n", a, attributes[a].name)
// 				fmt.Printf("---- ---- ---- ---- type %d\n", attributes[a].attributeType)
// 				fmt.Printf("---- ---- ---- ---- index %d\n", attributes[a].index)
// 				fmt.Printf("---- ---- ---- ---- accessor %d\n", attributes[a].accessor)
// 				//printI32Array("---- ---- ---- morphs", morphs)
// 			}
// 		}
// 	}
// 	fmt.Printf("]\n")

// 	var accessors []Accessor = g_assets[id.asset].accessors
// 	var accessorCount int32 = len(accessors)
// 	fmt.Printf("accessors, count : %d [\n", accessorCount)
// 	for i := 0; i < accessorCount; i++ {
// 		fmt.Printf("---- accessor %d, name '%s'\n", i, accessors[i].name)
// 		fmt.Printf("---- ---- bufferView %d\n", accessors[i].bufferView)
// 		fmt.Printf("---- ---- byteOffset %d\n", accessors[i].byteOffset)
// 		fmt.Printf("---- ---- componentType %d\n", accessors[i].componentType)
// 		fmt.Printf("---- ---- attributeType %d, %s\n", accessors[i].attributeTypeEnum, accessors[i].attributeTypeStr)
// 		fmt.Printf("---- ---- count %f\n", accessors[i].count)
// 		var min []float32 = accessors[i].min
// 		printF32Array("---- ---- ---- min", min)
// 		var max []float32 = accessors[i].max
// 		printF32Array("---- ---- ---- max", max)
// 	}
// 	fmt.Printf("]\n")

// 	var materials []Material = g_assets[id.asset].materials
// 	var materialCount int32 = len(materials)
// 	fmt.Printf("materials, count : %d [\n", materialCount)
// 	for i := 0; i < materialCount; i++ {
// 		fmt.Printf("---- material %d, name '%s'\n", i, materials[i].name)
// 		fmt.Printf("---- ---- metallicRoughness\n")
// 		fmt.Printf("---- ---- ---- baseColorFactor %f, %f, %f, %f\n",
// 			materials[i].PbrMetallicRoughness.baseColorFactor[0],
// 			materials[i].PbrMetallicRoughness.baseColorFactor[1],
// 			materials[i].PbrMetallicRoughness.baseColorFactor[2],
// 			materials[i].PbrMetallicRoughness.baseColorFactor[3])
// 		fmt.Printf("---- ---- ---- baseColorTexture %d\n", materials[i].PbrMetallicRoughness.baseColorTexture.index)
// 		fmt.Printf("---- ---- ---- metallicFactor %f\n", materials[i].PbrMetallicRoughness.metallicFactor)
// 		fmt.Printf("---- ---- ---- roughnessFactor %f\n", materials[i].PbrMetallicRoughness.roughnessFactor)
// 		fmt.Printf("---- ---- ---- metallicRoughnessTexture %d\n", materials[i].PbrMetallicRoughness.metallicRoughnessTexture.index)
// 	}
// 	fmt.Printf("]\n")

// 	var bufferViews []BufferView = g_assets[id.asset].bufferViews
// 	var bufferViewCount int32 = len(bufferViews)
// 	fmt.Printf("bufferViews, count : %d [\n", bufferViewCount)
// 	for i := 0; i < bufferViewCount; i++ {
// 		fmt.Printf("---- bufferView %d, name '%s'\n", i, bufferViews[i].name)
// 		fmt.Printf("---- ---- buffer %d\n", bufferViews[i].buffer)
// 		fmt.Printf("---- ---- byteOffset %d\n", bufferViews[i].byteOffset)
// 		fmt.Printf("---- ---- byteLength %d\n", bufferViews[i].byteLength)
// 		fmt.Printf("---- ---- byteStride %d\n", bufferViews[i].byteStride)
// 		fmt.Printf("---- ---- target %d\n", bufferViews[i].target)
// 	}
// 	fmt.Printf("]\n")

// 	var buffers []Buffer = g_assets[id.asset].buffers
// 	var bufferCount int32 = len(buffers)
// 	fmt.Printf("buffers, count : %d [\n", bufferCount)
// 	for i := 0; i < bufferCount; i++ {
// 		fmt.Printf("---- buffer %d, name '%s'\n", i, buffers[i].name)
// 		fmt.Printf("---- ---- %s\n", buffers[i].uri)
// 		fmt.Printf("---- ---- %d\n", buffers[i].byteLength)
// 	}
// 	fmt.Printf("]\n")

// 	var cameras []Camera = g_assets[id.asset].cameras
// 	var cameraCount int32 = len(cameras)
// 	fmt.Printf("cameras, count : %d [\n", cameraCount)
// 	for i := 0; i < cameraCount; i++ {
// 		fmt.Printf("---- camera %d, name '%s'\n", i, cameras[i].name)
// 		fmt.Printf("---- ---- type %d, %s\n", cameras[i].cameraTypeEnum, cameras[i].cameraTypeStr)
// 		fmt.Printf("---- ---- ---- perspectiveCamera\n")
// 		fmt.Printf("---- ---- ---- ---- aspectRatio %f\n", cameras[i].perspectiveCamera.aspectRatio)
// 		fmt.Printf("---- ---- ---- ---- yfov %f\n", cameras[i].perspectiveCamera.yfov)
// 		fmt.Printf("---- ---- ---- ---- znear %f\n", cameras[i].perspectiveCamera.znear)
// 		fmt.Printf("---- ---- ---- ---- zfar %f\n", cameras[i].perspectiveCamera.zfar)
// 	}
// 	fmt.Printf("]\n")
// }
// */
//------------------------------------------------------------------------------
func parseAsset(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		if key == "extras" {
			if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
				out = AssetInvalid()
				return
			}

			for alwaysMore(file) {
				var extraKey string
				if json.ReadStr(file, &extraKey) == false {
					out = AssetInvalid()
					return
				}

				var extraValue string
				if json.ReadStr(file, &extraValue) == false {
					out = AssetInvalid()
					return
				}

				if extraKey == "author" {
					g_assets[id.asset].author = extraValue
				} else if extraKey == "license" {
					g_assets[id.asset].license = extraValue
				} else if extraKey == "source" {
					g_assets[id.asset].source = extraValue
				} else if extraKey == "title" {
					g_assets[id.asset].title = extraValue
				} else {
					out = AssetInvalid() // TODO : handle other extras.
					return
				}
			}

			if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
				out = AssetInvalid()
				return
			}
		} else {
			var value string
			if json.ReadStr(file, &value) == false {
				out = AssetInvalid()
				return
			}

			if key == "copyright" {
				g_assets[id.asset].copyright = value
			} else if key == "generator" {
				g_assets[id.asset].generator = value
			} else if key == "version" {
				g_assets[id.asset].version = value
			} else if key == "minVersion" {
				g_assets[id.asset].minVersion = value
			} else {
				out = AssetInvalid()
				return
			}
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseScene(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	var value int32
	if json.ReadI32(file, &value) == false {
		out = AssetInvalid()
		return
	}

	g_assets[id.asset].scene = value
	return
}

//------------------------------------------------------------------------------
func parseScenes(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var scene Scene

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "nodes" {
				if json.ReadI32Slice(file, &scene.nodes) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				scene.name = name
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}
		g_assets[id.asset].scenes = append(g_assets[id.asset].scenes, scene)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parsePerspectiveCamera(file int32, id AssetId) (out AssetId, perspectiveCamera PerspectiveCamera) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		if key == "aspectRatio" {
			var aspectRatio float32
			if json.ReadF32(file, &aspectRatio) == false {
				out = AssetInvalid()
				return
			}
			perspectiveCamera.aspectRatio = aspectRatio
		} else if key == "yfov" {
			var yfov float32
			if json.ReadF32(file, &yfov) == false {
				out = AssetInvalid()
				return
			}
			perspectiveCamera.yfov = yfov
		} else if key == "zfar" {
			var zfar float32
			if json.ReadF32(file, &zfar) == false {
				out = AssetInvalid()
				return
			}
			perspectiveCamera.zfar = zfar
		} else if key == "znear" {
			var znear float32
			if json.ReadF32(file, &znear) == false {
				out = AssetInvalid()
				return
			}
			perspectiveCamera.znear = znear
		} else {
			fmt.Printf("invalid key %s\n", key)
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseCameras(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var camera Camera

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				camera.name = name
			} else if key == "perspective" {
				var perspectiveCamera PerspectiveCamera
				out, perspectiveCamera = parsePerspectiveCamera(file, id)
				if AssetIsValid(out) == false {
					return
				}
				camera.perspectiveCamera = perspectiveCamera
			} else if key == "type" {
				var cameraType string
				if json.ReadStr(file, &cameraType) == false {
					out = AssetInvalid()
					return
				}
				camera.cameraTypeStr = cameraType
				if cameraType == "perspective" {
					camera.cameraTypeEnum = CAMERA_PERSPECTIVE
				} else {
					fmt.Printf("invalid camera type %s", cameraType)
					out = AssetInvalid()
					return
				}
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		g_assets[id.asset].cameras = append(g_assets[id.asset].cameras, camera)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseNodes(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var node Node
		node.mesh = -1
		node.matrix = m44.IDENTITY
		node.scale = v3.ONE
		node.rotation = v4.ALPHA
		node.translation = v3.ZERO

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				node.name = name
			} else if key == "mesh" {
				var mesh int32
				if json.ReadI32(file, &mesh) == false {
					out = AssetInvalid()
					return
				}
				node.mesh = mesh
			} else if key == "children" {
				if json.ReadI32Slice(file, &node.children) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "matrix" {
				node.hasMatrix = true
				if readM44F(file, &node.matrix) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "rotation" {
				node.hasRotation = true
				if readV4F(file, &node.rotation) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "translation" {
				node.hasTranslation = true
				if readV3F(file, &node.translation) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "scale" {
				node.hasScale = true
				if readV3F(file, &node.scale) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "camera" {
				var camera int32
				if json.ReadI32(file, &camera) == false {
					out = AssetInvalid()
					return
				}
				node.camera = camera
			} else if key == "skin" {
				var skin int32
				if json.ReadI32(file, &skin) == false {
					out = AssetInvalid()
					return
				}
				node.skin = skin
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if node.hasMatrix == false {
			if node.hasScale == false {
				node.scale = v3.ONE
			}

			if node.hasRotation == false {
				node.rotation = v4.ALPHA
			}

			if node.hasTranslation == false {
				node.translation = v3.ZERO
			}

			node.matrix = m44.Makev_SQT(node.scale, node.rotation, node.translation)
		}

		g_assets[id.asset].nodes = append(g_assets[id.asset].nodes, node)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseAttributes(file int32, id AssetId) (out AssetId, attributes []Attribute) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		var attributeType int32 = ATTRIBUTE_INVALID
		var index int32
		if key == "POSITION" {
			attributeType = ATTRIBUTE_POSITION
		} else if key == "NORMAL" {
			attributeType = ATTRIBUTE_NORMAL
		} else if key == "TANGENT" {
			attributeType = ATTRIBUTE_TANGENT
		} else if key == "TEXCOORD_0" {
			attributeType = ATTRIBUTE_TEXCOORD
		} else if key == "TEXCOORD_1" {
			attributeType = ATTRIBUTE_TEXCOORD
			index = 1
		} else if key == "TEXCOORD_2" {
			attributeType = ATTRIBUTE_TEXCOORD
			index = 2
		} else if key == "TEXCOORD_3" {
			attributeType = ATTRIBUTE_TEXCOORD
			index = 3
		} else if key == "COLOR_0" {
			attributeType = ATTRIBUTE_COLOR
		} else if key == "JOINTS_0" {
			attributeType = ATTRIBUTE_JOINT
		} else if key == "WEIGHTS_0" {
			attributeType = ATTRIBUTE_WEIGHT
		} else {
			fmt.Printf("invalid key %s\n", key)
			out = AssetInvalid()
			return
		}

		var accessor int32
		if json.ReadI32(file, &accessor) == false {
			out = AssetInvalid()
			return
		}

		var attribute Attribute
		attribute.name = key
		attribute.AttributeType = attributeType
		attribute.index = index
		attribute.Accessor = accessor
		attributes = append(attributes, attribute)
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parsePrimitives(file int32, id AssetId) (out AssetId, primitives []Primitive) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var primitive Primitive
		primitive.Mode = gl.TRIANGLES
		primitive.Material = -1
		primitive.Indices = -1

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "attributes" {
				var attributes []Attribute
				out, attributes = parseAttributes(file, id)
				if AssetIsValid(out) == false {
					return
				}
				primitive.Attributes = attributes
			} else if key == "indices" {
				var indices int32
				if json.ReadI32(file, &indices) == false {
					out = AssetInvalid()
					return
				}
				primitive.Indices = indices
			} else if key == "mode" {
				var mode int32
				if json.ReadI32(file, &mode) == false {
					out = AssetInvalid()
					return
				}
				primitive.Mode = mode
			} else if key == "material" {
				var material int32
				if json.ReadI32(file, &material) == false {
					out = AssetInvalid()
					return
				}
				primitive.Material = material
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		primitives = append(primitives, primitive)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseMeshes(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var mesh Mesh

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				mesh.name = name
			} else if key == "primitives" {
				var primitives []Primitive
				out, primitives = parsePrimitives(file, id)
				if AssetIsValid(out) == false {
					return
				}
				mesh.primitives = primitives
			}
		}

		g_assets[id.asset].meshes = append(g_assets[id.asset].meshes, mesh)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseAccessors(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var accessor Accessor

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "bufferView" {
				var bufferView int32
				if json.ReadI32(file, &bufferView) == false {
					out = AssetInvalid()
					return
				}
				accessor.bufferView = bufferView
			} else if key == "byteOffset" {
				var byteOffset int64
				if json.ReadI64(file, &byteOffset) == false {
					out = AssetInvalid()
					return
				}
				accessor.byteOffset = byteOffset
			} else if key == "componentType" {
				var componentType int32
				if json.ReadI32(file, &componentType) == false {
					out = AssetInvalid()
					return
				}
				accessor.ComponentType = componentType
			} else if key == "count" {
				var count int32
				if json.ReadI32(file, &count) == false {
					out = AssetInvalid()
					return
				}
				accessor.Count = count
			} else if key == "max" {
				if json.ReadF32Slice(file, &accessor.Max) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "min" {
				if json.ReadF32Slice(file, &accessor.Min) == false {
					out = AssetInvalid()
					return
				}
			} else if key == "type" {
				var value string
				if json.ReadStr(file, &value) == false {
					out = AssetInvalid()
					return
				}

				var attributeTypeEnum int32 = TYPE_INVALID
				var componentCount int32 = 0
				if value == "SCALAR" {
					attributeTypeEnum = TYPE_SCALAR
					componentCount = 1
				} else if value == "VEC2" {
					attributeTypeEnum = TYPE_VEC2
					componentCount = 2
				} else if value == "VEC3" {
					attributeTypeEnum = TYPE_VEC3
					componentCount = 3
				} else if value == "VEC4" {
					attributeTypeEnum = TYPE_VEC4
					componentCount = 4
				} else if value == "MAT2" {
					attributeTypeEnum = TYPE_MAT2
					componentCount = 4
				} else if value == "MAT3" {
					attributeTypeEnum = TYPE_MAT3
					componentCount = 9
				} else if value == "MAT4" {
					attributeTypeEnum = TYPE_MAT4
					componentCount = 16
				} else {
					fmt.Printf("invalid type %s\n", value)
					out = AssetInvalid()
					return
				}

				accessor.AttributeTypeEnum = attributeTypeEnum
				accessor.AttributeTypeStr = value
				accessor.ComponentCount = uint32(componentCount)
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		g_assets[id.asset].accessors = append(g_assets[id.asset].accessors, accessor)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseTexture(file int32, id AssetId) (out AssetId, texture TextureView) {
	out = id
	texture = InvalidTextureView()
	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		if key == "index" {
			var index int32
			if json.ReadI32(file, &index) == false {
				out = AssetInvalid()
				return
			}
			texture.Index = index
		} else if key == "scale" {
			var scale float32
			if json.ReadF32(file, &scale) == false {
				out = AssetInvalid()
				return
			}
			texture.Scale = scale
		} else if key == "texCoord" {
			var texCoord int32
			if json.ReadI32(file, &texCoord) == false {
				out = AssetInvalid()
				return
			}
			texture.texCoord = texCoord
		} else {
			fmt.Printf("unhandled texture param")
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
		out = AssetInvalid()
		return
	}
	return
}

func parsePbrSpecularGlossiness(file int32, id AssetId) (out AssetId, specularGlossiness SpecularGlossiness) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	specularGlossiness.DiffuseFactor = v4.ONE
	specularGlossiness.GlossinessFactor = 1.0
	specularGlossiness.SpecularFactor = v3.ONE

	specularGlossiness.DiffuseTexture = InvalidTextureView()
	specularGlossiness.SpecularGlossinessTexture = InvalidTextureView()

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		if key == "diffuseFactor" {
			var diffuseFactor math.V4
			if readV4F(file, &diffuseFactor) == false {
				out = AssetInvalid()
				return
			}
			specularGlossiness.DiffuseFactor = diffuseFactor
		} else if key == "glossinessFactor" {
			var glossinessFactor float32
			if json.ReadF32(file, &glossinessFactor) == false {
				out = AssetInvalid()
				return
			}
			specularGlossiness.GlossinessFactor = glossinessFactor
		} else if key == "specularFactor" {
			var specularFactor math.V3
			if readV3F(file, &specularFactor) == false {
				out = AssetInvalid()
				return
			}
			specularGlossiness.SpecularFactor = specularFactor
		} else if key == "diffuseTexture" {
			var texture TextureView
			out, texture = parseTexture(file, id)
			if AssetIsValid(out) == false {
				return
			}
			specularGlossiness.DiffuseTexture = texture
		} else if key == "specularGlossinessTexture" {
			var texture TextureView
			out, texture = parseTexture(file, id)
			if AssetIsValid(out) == false {
				return
			}
			specularGlossiness.SpecularGlossinessTexture = texture
		} else {
			fmt.Printf("invalid key %s\n", key)
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parsePbrMetallicRoughness(file int32, id AssetId) (out AssetId, metallicRoughness MetallicRoughness) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	metallicRoughness.BaseColorFactor = v4.ONE
	metallicRoughness.MetallicFactor = 1.0
	metallicRoughness.RoughnessFactor = 1.0

	metallicRoughness.BaseColorTexture = InvalidTextureView()
	metallicRoughness.MetallicRoughnessTexture = InvalidTextureView()

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var key string
		if json.ReadStr(file, &key) == false {
			out = AssetInvalid()
			return
		}

		if key == "baseColorFactor" {
			var baseColorFactor math.V4
			if readV4F(file, &baseColorFactor) == false {
				out = AssetInvalid()
				return
			}
			metallicRoughness.BaseColorFactor = baseColorFactor
		} else if key == "baseColorTexture" {
			var texture TextureView
			out, texture = parseTexture(file, id)
			if AssetIsValid(out) == false {
				return
			}
			metallicRoughness.BaseColorTexture = texture
		} else if key == "roughnessFactor" {
			var roughnessFactor float32
			if json.ReadF32(file, &roughnessFactor) == false {
				out = AssetInvalid()
				fmt.Printf("fail to parse roughnessFactor\n")
				return
			}
			metallicRoughness.RoughnessFactor = roughnessFactor
		} else if key == "metallicFactor" {
			var metallicFactor float32
			if json.ReadF32(file, &metallicFactor) == false {
				out = AssetInvalid()
				fmt.Printf("fail to parse metallicFactor\n")
				return
			}
			metallicRoughness.MetallicFactor = metallicFactor
		} else if key == "metallicRoughnessTexture" {
			var texture TextureView
			out, texture = parseTexture(file, id)
			if AssetIsValid(out) == false {
				return
			}
			metallicRoughness.MetallicRoughnessTexture = texture
		} else {
			fmt.Printf("invalid key %s\n", key)
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseMaterials(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		fmt.Printf("failed to match delim square left\n")
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			fmt.Printf("failed to match delim curly left\n")
			out = AssetInvalid()
			return
		}

		var material Material
		material.EmissiveTexture = InvalidTextureView()
		material.NormalTexture = InvalidTextureView()
		material.OcclusionTexture = InvalidTextureView()

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				fmt.Printf("failed to read material key\n")
				out = AssetInvalid()
				return
			}

			if key == "pbrMetallicRoughness" {
				var pbrMetallicRoughness MetallicRoughness
				out, pbrMetallicRoughness = parsePbrMetallicRoughness(file, id)
				if AssetIsValid(out) == false {
					fmt.Printf("failed to read pbr metallic roughness\n")
					return
				}
				material.PbrMetallicRoughness = pbrMetallicRoughness
			} else if key == "alphaMode" {
				var alphaMode string
				if json.ReadStr(file, &alphaMode) == false {
					fmt.Printf("failed to read alpha mode\n")
					out = AssetInvalid()
					return
				}

				if alphaMode == "OPAQUE" {
					material.AlphaMode = ALPHA_OPAQUE
				} else if alphaMode == "MASK" {
					material.AlphaMode = ALPHA_BLEND // TODO : imploment ALPHA_MASK
					//fmt.Printf("unhandled alphaMode %s\n", alphaMode)
					//out = AssetInvalid()
					//return
				} else if alphaMode == "BLEND" {
					material.AlphaMode = ALPHA_BLEND
				} else {
					fmt.Printf("invalid alphaMode %s\n", alphaMode)
					out = AssetInvalid()
					return
				}
			} else if key == "normalTexture" {
				var texture TextureView
				out, texture = parseTexture(file, id)
				if AssetIsValid(out) == false {
					fmt.Printf("failed to parse normal texture\n")
					return
				}
				material.NormalTexture = texture
			} else if key == "occlusionTexture" {
				var texture TextureView
				out, texture = parseTexture(file, id)
				if AssetIsValid(out) == false {
					fmt.Printf("failed to parse occlusion texture\n")
					return
				}
				material.OcclusionTexture = texture
			} else if key == "emissiveTexture" {
				var texture TextureView
				out, texture = parseTexture(file, id)
				if AssetIsValid(out) == false {
					fmt.Printf("failed to parse emissive texture\n")
					return
				}
				material.EmissiveTexture = texture
			} else if key == "emissiveFactor" {
				var emissiveFactor math.V3
				if readV3F(file, &emissiveFactor) == false {
					out = AssetInvalid()
					fmt.Printf("failed to parse emissive factor\n")
					return
				}
				material.EmissiveFactor = v4.Make_v31(emissiveFactor, 1.0)
			} else if key == "doubleSided" {
				var doubleSided bool
				if json.ReadBool(file, &doubleSided) == false {
					out = AssetInvalid()
					fmt.Printf("failed to parse double sided\n")
					return
				}
				if doubleSided == true {
					material.DoubleSided = 1
				}
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					fmt.Printf("failed to parse name\n")
					return
				}
				material.name = name
			} else if key == "extensions" {
				if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
					out = AssetInvalid()
					fmt.Printf("failed to parse extensions\n")
					return
				}

				for alwaysMore(file) {
					var extension string
					if json.ReadStr(file, &extension) == false {
						out = AssetInvalid()
						fmt.Printf("failed to parse extension\n")
						return
					}

					//material.KHR_extension_unlit = true
					if extension == "KHR_materials_unlit" {
						if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
							out = AssetInvalid()
							fmt.Printf("failed to parse KHR_extension_unlit\n")
							return
						}

						if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
							out = AssetInvalid()
							fmt.Printf("failed to parse KHR_extension_unlit\n")
							return
						}
					} else if extension == "KHR_materials_pbrSpecularGlossiness" {
						var pbrSpecularGlossiness SpecularGlossiness
						out, pbrSpecularGlossiness = parsePbrSpecularGlossiness(file, id)
						if AssetIsValid(out) == false {
							fmt.Printf("failed to read pbr specular glossiness\n")
							return
						}
						material.PbrSpecularGlossiness = pbrSpecularGlossiness
					} else {
						out = AssetInvalid()
						fmt.Printf("unhandled extensions %s\n", extension)
						return
					}
				}
				if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
					out = AssetInvalid()
					fmt.Printf("failed to parse extensions 2\n")
					return
				}
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			fmt.Printf("failed to match delim curly right\n")
			return
		}

		g_assets[id.asset].materials = append(g_assets[id.asset].materials, material)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	fmt.Printf("failed to match delim square right\n")
	out = AssetInvalid()
	return
}

// InvalidTextureView ...
func InvalidTextureView() (out TextureView) {
	out.Index = -1
	out.Scale = 1.0
	return
}

//------------------------------------------------------------------------------
func parseTextures(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var texture Texture
		texture.Source = -1
		texture.sampler = -1

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "sampler" {
				var sampler int32
				if json.ReadI32(file, &sampler) == false {
					out = AssetInvalid()
					return
				}
				texture.sampler = sampler
			} else if key == "source" {
				var source int32
				if json.ReadI32(file, &source) == false {
					out = AssetInvalid()
					return
				}
				texture.Source = source
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				texture.name = name
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].textures = append(g_assets[id.asset].textures, texture)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseImages(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var image Image

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "uri" {
				var uri string
				if json.ReadStr(file, &uri) == false {
					out = AssetInvalid()
					return
				}
				image.Uri = uri
			} else if key == "mimeType" {
				var mimeType string
				if json.ReadStr(file, &mimeType) == false {
					out = AssetInvalid()
					return
				}
				image.mimeType = mimeType
			} else if key == "bufferView" {
				var bufferView int32
				if json.ReadI32(file, &bufferView) == false {
					out = AssetInvalid()
					return
				}
				image.bufferView = bufferView
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				image.name = name
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].images = append(g_assets[id.asset].images, image)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseSamplers(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var sampler Sampler

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "magFilter" {
				var magFilter int32
				if json.ReadI32(file, &magFilter) == false {
					out = AssetInvalid()
					return
				}
				sampler.magFilter = magFilter
			} else if key == "minFilter" {
				var minFilter int32
				if json.ReadI32(file, &minFilter) == false {
					out = AssetInvalid()
					return
				}
				sampler.minFilter = minFilter
			} else if key == "wrapS" {
				var wrapS int32
				if json.ReadI32(file, &wrapS) == false {
					out = AssetInvalid()
					return
				}
				sampler.wrapS = wrapS
			} else if key == "wrapT" {
				var wrapT int32
				if json.ReadI32(file, &wrapT) == false {
					out = AssetInvalid()
					return
				}
				sampler.wrapT = wrapT
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				sampler.name = name
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].samplers = append(g_assets[id.asset].samplers, sampler)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseBufferViews(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var bufferView BufferView

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "buffer" {
				var buffer int32
				if json.ReadI32(file, &buffer) == false {
					out = AssetInvalid()
					return
				}
				bufferView.buffer = buffer
			} else if key == "byteOffset" {
				var byteOffset int64
				if json.ReadI64(file, &byteOffset) == false {
					out = AssetInvalid()
					return
				}
				bufferView.byteOffset = byteOffset
			} else if key == "byteLength" {
				var byteLength int32
				if json.ReadI32(file, &byteLength) == false {
					out = AssetInvalid()
					return
				}
				bufferView.byteLength = byteLength
			} else if key == "byteStride" {
				var byteStride int32
				if json.ReadI32(file, &byteStride) == false {
					out = AssetInvalid()
					return
				}
				bufferView.byteStride = byteStride
			} else if key == "target" {
				var target int32
				if json.ReadI32(file, &target) == false {
					out = AssetInvalid()
					return
				}
				bufferView.target = target
			} else if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				bufferView.name = name
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		g_assets[id.asset].bufferViews = append(g_assets[id.asset].bufferViews, bufferView)

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseBuffers(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var buffer Buffer

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "byteLength" {
				var byteLength int32
				if json.ReadI32(file, &byteLength) == false {
					out = AssetInvalid()
					return
				}
				buffer.byteLength = byteLength
			} else if key == "uri" {
				var uri string
				if json.ReadStr(file, &uri) == false {
					out = AssetInvalid()
					return
				}
				buffer.uri = uri
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].buffers = append(g_assets[id.asset].buffers, buffer)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseAnimationSamplers(file int32, id AssetId) (out AssetId, samplers []AnimationSampler) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var sampler AnimationSampler

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "interpolation" {
				var interpolation string
				if json.ReadStr(file, &interpolation) == false {
					out = AssetInvalid()
					return
				}
				if interpolation == "LINEAR" {
					sampler.interpolation = INTERPOLATION_LINEAR
				} else if interpolation == "STEP" {
					sampler.interpolation = INTERPOLATION_STEP
				} else if interpolation == "CUBICSPLINE" {
					sampler.interpolation = INTERPOLATION_CUBIC_SPLINE
				} else {
					fmt.Printf("invalid interpolation %s\n", interpolation)
					out = AssetInvalid()
					return
				}
			} else if key == "input" {
				var inp int32
				if json.ReadI32(file, &inp) == false {
					fmt.Printf("invalid input\n")
					out = AssetInvalid()
					return
				}
				sampler.InputA = inp
			} else if key == "output" {
				var outs int32 // ISSUE : variable can be name out without compilation error
				if json.ReadI32(file, &outs) == false {
					fmt.Printf("invalid output\n")
					out = AssetInvalid()
					return
				}
				sampler.OutputA = outs
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		samplers = append(samplers, sampler)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseAnimationChannels(file int32, id AssetId) (out AssetId, channels []AnimationChannel) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var channel AnimationChannel

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "sampler" {
				var index int32
				if json.ReadI32(file, &index) == false {
					out = AssetInvalid()
					return
				}
				channel.Sampler = index
			} else if key == "target" {
				if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
					out = AssetInvalid()
					return
				}

				for alwaysMore(file) {
					var targetKey string
					if json.ReadStr(file, &targetKey) == false {
						out = AssetInvalid()
						return
					}
					if targetKey == "node" {
						var node int32
						if json.ReadI32(file, &node) == false {
							out = AssetInvalid()
							return
						}
						channel.Node = node
					} else if targetKey == "path" {
						var path string
						if json.ReadStr(file, &path) == false {
							out = AssetInvalid()
							return
						}
						if path == "translation" {
							channel.Path = ANIMATION_PATH_TRANSLATION
						} else if path == "rotation" {
							channel.Path = ANIMATION_PATH_ROTATION
						} else if path == "scale" {
							channel.Path = ANIMATION_PATH_SCALE
						} else if path == "weights" {
							channel.Path = ANIMATION_PATH_WEIGHTS
						} else {
							out = AssetInvalid()
							return
						}
					} else {
						fmt.Printf("invalid key %s\n", targetKey)
					}
				}

				if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
					out = AssetInvalid()
					return
				}
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		channels = append(channels, channel)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseAnimations(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var animation Animation

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				animation.name = name
			} else if key == "channels" {
				var channels []AnimationChannel
				out, channels = parseAnimationChannels(file, id)
				if AssetIsValid(out) == false {
					return
				}
				animation.Channels = channels
			} else if key == "samplers" {
				var samplers []AnimationSampler
				out, samplers = parseAnimationSamplers(file, id)
				if AssetIsValid(out) == false {
					return
				}
				animation.Samplers = samplers
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].animations = append(g_assets[id.asset].animations, animation)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseSkins(file int32, id AssetId) (out AssetId) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
			out = AssetInvalid()
			return
		}

		var skin Skin

		for alwaysMore(file) {
			var key string
			if json.ReadStr(file, &key) == false {
				out = AssetInvalid()
				return
			}

			if key == "name" {
				var name string
				if json.ReadStr(file, &name) == false {
					out = AssetInvalid()
					return
				}
				skin.name = name
			} else if key == "inverseBindMatrices" {
				var inverseBindMatrices int32
				if json.ReadI32(file, &inverseBindMatrices) == false {
					out = AssetInvalid()
					return
				}
				skin.InverseBindMatrices = inverseBindMatrices
			} else if key == "skeleton" {
				var skeleton int32
				if json.ReadI32(file, &skeleton) == false {
					out = AssetInvalid()
					return
				}
				skin.Skeleton = skeleton
			} else if key == "joints" {
				if json.ReadI32Slice(file, &skin.Joints) == false {
					out = AssetInvalid()
					return
				}
			} else {
				fmt.Printf("invalid key %s\n", key)
				out = AssetInvalid()
				return
			}
		}

		if json.MatchDelim(file, json.JSON_DELIM_CURLY_RIGHT) == false {
			out = AssetInvalid()
			return
		}

		g_assets[id.asset].skins = append(g_assets[id.asset].skins, skin)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

//------------------------------------------------------------------------------
func parseExtensions(file int32, id AssetId) (out AssetId, extensions []string) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out = id

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_LEFT) == false {
		out = AssetInvalid()
		return
	}

	for alwaysMore(file) {
		var extension string
		if json.ReadStr(file, &extension) == false {
			out = AssetInvalid()
			return
		}

		//g_assets[id.asset].extensionsUsed = append(g_assets[id.asset].extensionsUsed, extension)
	}

	if json.MatchDelim(file, json.JSON_DELIM_SQUARE_RIGHT) {
		return
	}

	out = AssetInvalid()
	return
}

// AssetGetNodeCount ...
func AssetGetNodeCount(id AssetId) (out int32) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	var nodes []Node = g_assets[id.asset].nodes
	out = int32(len(nodes))
	return
}

// AssetGetRootNodes ...
func AssetGetRootNodes(id AssetId) (out []int32) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	var scenes []Scene = g_assets[id.asset].scenes
	var scene int = int(g_assets[id.asset].scene)

	if scene >= 0 && scene < len(scenes) {
		out = scenes[scene].nodes
	}
	return
}

// AssetGetMaterial ...
func AssetGetMaterial(id AssetId, material int32) (out Material) { // TODO : handle errors
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	var materials []Material = g_assets[id.asset].materials

	if material >= 0 && material < int32(len(materials)) {
		out = materials[material]
	}
	return
}

// AssetIsValidTexture ...
func AssetIsValidTexture(id AssetId, texture int32) (out bool) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var textures []Texture = g_assets[id.asset].textures
	out = texture >= 0 && texture < int32(len(textures))
	return
}

// AssetGetTexture ...
func AssetGetTexture(id AssetId, texture int32) (out Texture) { // TODO :handle errors
	out.sampler = -1
	out.Source = -1
	if AssetIsValidTexture(id, texture) {
		out = g_assets[id.asset].textures[texture]
	}
	return
}

// AssetGetImage ...
func AssetGetImage(id AssetId, image int32) (out Image) { // TODO : handle errors
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	out.bufferView = -1

	var images []Image = g_assets[id.asset].images

	if image >= 0 && image < int32(len(images)) {
		out = images[image]
	}
	return
}

// AssetGetAccessor ...
func AssetGetAccessor(id AssetId, accessor int32) (out Accessor) { // TODO : handle errors
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	var accessors []Accessor = g_assets[id.asset].accessors
	var accessorCount int32 = int32(len(accessors))
	if accessor >= 0 && accessor < accessorCount {
		out = accessors[accessor]
	}
	return
}

// AssetGetAnimations ...
func AssetGetAnimations(id AssetId) (out []Animation) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	out = g_assets[id.asset].animations
	return
}

// // NodeHasMatrix ...
// func NodeHasMatrix(id AssetId, node int32) (out bool) {
// 	out = false
// 	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
// 	var nodes []Node = g_assets[id.asset].nodes
// 	if node >= 0 && node < len(nodes) {
// 		out = nodes[node].hasMatrix
// 	}
// }

// NodeGetMesh ...
func NodeGetMesh(id AssetId, node int32) (out int32) {
	out = -1
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	if node >= 0 && node < int32(len(nodes)) {
		out = nodes[node].mesh
	}
	return
}

// // NodeGetName ...
// func NodeGetName(id AssetId, node int32) (out string) {
// 	out = ""
// 	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
// 	var nodes []Node = g_assets[id.asset].nodes
// 	if node >= 0 && node < len(nodes) {
// 		out = nodes[node].name
// 	}
// }

// NodeGetSkin ...
func NodeGetSkin(id AssetId, node int32) (out Skin) { // TODO : handle errors
	out.InverseBindMatrices = -1
	out.Skeleton = -1
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	if node >= 0 && node < int32(len(nodes)) {
		var skin int32 = nodes[node].skin
		var skins []Skin = g_assets[id.asset].skins
		if skin >= 0 && skin < int32(len(skins)) {
			out = skins[skin]
		}
	}
	return
}

// NodeGetMatrix ...
func NodeGetMatrix(id AssetId, node int32) (out math.M44) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	out = nodes[node].matrix
	return
}

// NodeGetRotation...
func NodeGetRotation(id AssetId, node int32) (out math.V4) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	out = nodes[node].rotation
	return
}

// NodeGetTranslation ...
func NodeGetTranslation(id AssetId, node int32) (out math.V3) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	out = nodes[node].translation
	return
}

// NodeGetScale ...
func NodeGetScale(id AssetId, node int32) (out math.V3) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	out = nodes[node].scale
	return
}

// NodeGetChildren ...
func NodeGetChildren(id AssetId, node int32) (out []int32) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var nodes []Node = g_assets[id.asset].nodes
	if node >= 0 && node < int32(len(nodes)) {
		out = nodes[node].children
	}
	return
}

// AnimationGetLength ...
func AnimationGetLength(id AssetId, animation int32) (min float32, max float32) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")
	var animations []Animation = g_assets[id.asset].animations
	var count int32 = int32(len(animations))
	min = math.MAX_f32
	max = math.MIN_f32
	if animation >= 0 && animation < count {
		var accessors []Accessor = g_assets[id.asset].accessors
		var accessorCount int = len(accessors)
		var samplers []AnimationSampler = animations[animation].Samplers
		var samplerCount int = len(samplers)
		for i := 0; i < samplerCount; i++ {
			var accessor int = int(samplers[i].InputA)
			if accessor >= 0 && accessor < accessorCount {
				if accessors[accessor].AttributeTypeEnum == TYPE_SCALAR {
					if accessors[accessor].Min[0] < min {
						min = accessors[accessor].Min[0]
					}
					if accessors[accessor].Max[0] > max {
						max = accessors[accessor].Max[0]
					}
				}
			}
		}
	}
	return
}

// // MeshGetName ...
// func MeshGetName(id AssetId, mesh int32) (out string) {
// 	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

// 	var meshes []Mesh = g_assets[id.asset].meshes
// 	if mesh >= 0 && mesh < len(meshes) {
// 		out = meshes[mesh].name
// 	}
// }

// MeshGetPrimitives ...
func MeshGetPrimitives(id AssetId, mesh int32) (out []Primitive) {
	utils.PanicIfNot(AssetIsValid(id), "invalid asset")

	var meshes []Mesh = g_assets[id.asset].meshes
	if mesh >= 0 && mesh < int32(len(meshes)) {
		out = meshes[mesh].primitives
	}
	return
}

func assetFailure(file int32, error string) (out AssetId) {
	fmt.Printf("gltf error : %s", error)
	_ /*var success bool */ = json.Close(file)
	out = AssetInvalid()
	return
}

func assetSuccess(file int32, id AssetId) (out AssetId) {
	if json.Close(file) {
		out = id
		return
	}

	out = AssetInvalid()
	return
}

// AssetCreate ...
func AssetCreate(dataDir string, filename string, options int32) (out AssetId) {
	var path string = fmt.Sprintf("%s%s", dataDir, filename)
	out.asset = int32(len(g_assets))

	var asset Asset
	CurrentAsset = out
	asset.path = path
	g_assets = append(g_assets, asset)

	var file int32 = json.Open(path)
	if alwaysMore(file) == false {
		out = assetFailure(file, "invalid gltf : failed to more\n")
		return
	}

	if json.MatchDelim(file, json.JSON_DELIM_CURLY_LEFT) == false {
		out = assetFailure(file, "invalid gltf : failed to read curly left\n")
		return
	}

	var stack int32 = 1
	var cont bool = true
	for cont == true {
		if alwaysMore(file) == false && stack <= 0 {
			cont = false
		} else {
			var tokenType int32
			var tokenSuccess bool
			tokenType, tokenSuccess = json.TokenNext(file)
			if tokenSuccess == false {
				out = assetFailure(file, "invalid gltf : failed to read token\n")
				return
			}

			var delimType int32
			var delimSuccess bool
			delimType, delimSuccess = json.TokenDelim(file)
			if delimSuccess && delimType == json.JSON_DELIM_CURLY_RIGHT {
				stack--
			} else {
				var key string
				var keySuccess bool
				key, keySuccess = json.TokenStr(file)

				if keySuccess == false {
					json.DebugToken(file, tokenType)
					out = assetFailure(file, "invalid gltf : failed to read string\n")
					return
				}
				//fmt.Printf(">>>>>>>>>>>>> parse %s\n", key)
				if key == "asset" {
					//fmt.Printf(">>>>>>>>>>>>> parse ASSET\n")
					out = parseAsset(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse ASSET\n")
				} else if key == "scene" {
					//fmt.Printf(">>>>>>>>>>>>> parse SCENE\n")
					out = parseScene(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse SCENE\n")
				} else if key == "scenes" {
					//fmt.Printf(">>>>>>>>>>>>> parse SCENES\n")
					out = parseScenes(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse SCENES\n")
				} else if key == "nodes" {
					//fmt.Printf(">>>>>>>>>>>>> parse NODES\n")
					out = parseNodes(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse NODES\n")
				} else if key == "cameras" {
					//fmt.Printf(">>>>>>>>>>>>> parse CAMERAS\n")
					out = parseCameras(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse CAMERAS\n")
				} else if key == "meshes" {
					//fmt.Printf(">>>>>>>>>>>>> parse MESHES\n")
					out = parseMeshes(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse MESHES\n")
				} else if key == "accessors" {
					//fmt.Printf(">>>>>>>>>>>>> parse ACCESSORS\n")
					out = parseAccessors(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse ACCESSORS\n")
				} else if key == "materials" {
					//fmt.Printf(">>>>>>>>>>>>> parse MATERIALS\n")
					out = parseMaterials(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse MATERIALS\n")
				} else if key == "textures" {
					//fmt.Printf(">>>>>>>>>>>>> parse TEXTURES\n")
					out = parseTextures(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse TEXTURES\n")
				} else if key == "images" {
					//fmt.Printf(">>>>>>>>>>>>> parse IMAGES\n")
					out = parseImages(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse IMAGES\n")
				} else if key == "samplers" {
					//fmt.Printf(">>>>>>>>>>>>> parse SAMPLERS\n")
					out = parseSamplers(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse SAMPLERS\n")
				} else if key == "bufferViews" {
					//fmt.Printf(">>>>>>>>>>>>> parse BUFFER_VIEWS\n")
					out = parseBufferViews(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse BUFFER_VIEWS\n")
				} else if key == "buffers" {
					//fmt.Printf(">>>>>>>>>>>>> parse BUFFERS\n")
					out = parseBuffers(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse BUFFERS\n")
				} else if key == "animations" {
					//fmt.Printf(">>>>>>>>>>>>> parse ANIMATIONS\n")
					out = parseAnimations(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse ANIMATIONS\n")
				} else if key == "skins" {
					//fmt.Printf(">>>>>>>>>>>>> parse SKINS\n")
					out = parseSkins(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse SKINS\n")
				} else if key == "extensionsUsed" {
					//fmt.Printf(">>>>>>>>>>>>> parse EXTENSIONS_USED\n")
					//var extensions []string
					out, _ /*extensions*/ = parseExtensions(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse EXTENSIONS_USED\n")
				} else if key == "extensionsRequired" {
					//fmt.Printf(">>>>>>>>>>>>> parse EXTENSIONS_REQUIRED\n")
					//var extensions []string
					out, _ /*extensions*/ = parseExtensions(file, out)
					//fmt.Printf("<<<<<<<<<<<<< parse EXTENSIONS_REQUIRED\n")
				} else {
					out = assetFailure(file, fmt.Sprintf("invalid gltf : parsing not implemented %s\n", key))
					return
					//fmt.Printf("<<<<<<<<<<<<<  parse %s\n", key)
				}
			}
		}
	}

	if AssetIsValid(out) == false {
		out = assetFailure(file, "invalid gltf\n")
		return
	}

	var accessors []Accessor = g_assets[out.asset].accessors
	var accessorCount int32 = int32(len(accessors))

	if (options & ASSET_GEOMETRY) != 0 {
		var meshes []Mesh = g_assets[out.asset].meshes
		var meshCount int = len(meshes)
		for i := 0; i < meshCount; i++ {
			var primitives []Primitive = meshes[i].primitives
			var primitiveCount int = len(primitives)
			for k := 0; k < primitiveCount; k++ {
				var attributes []Attribute = primitives[k].Attributes
				var attributeCount int = len(attributes)
				for l := 0; l < attributeCount; l++ {
					var accessor int32 = attributes[l].Accessor
					accessors[accessor].Loaded = true
				}
				var indicesAccessor int32 = primitives[k].Indices
				accessors[indicesAccessor].Loaded = true
			}
		}
	}

	if (options & ASSET_ANIMATION) != 0 {
		var animations []Animation = g_assets[out.asset].animations
		var animationCount int = len(animations)
		for i := 0; i < animationCount; i++ {
			var samplers []AnimationSampler = animations[i].Samplers
			var samplerCount int = len(samplers)
			for k := 0; k < samplerCount; k++ {
				var inputAccessor int32 = samplers[k].InputA
				var outputAccessor int32 = samplers[k].OutputA
				accessors[inputAccessor].Loaded = true
				accessors[outputAccessor].Loaded = true
			}
		}
	}

	var skins []Skin = g_assets[out.asset].skins
	var skinCount int = len(skins)
	for i := 0; i < skinCount; i++ {
		var ibmAccessor int32 = skins[i].InverseBindMatrices
		accessors[ibmAccessor].Loaded = true
	}

	var bufferViews []BufferView = g_assets[out.asset].bufferViews
	var bufferViewCount int = len(bufferViews)

	var buffers []Buffer = g_assets[out.asset].buffers
	var bufferCount int = len(buffers)

	var handle int32 = -1
	var previousUri string
	for a := 0; a < int(accessorCount); a++ {
		if accessors[a].Loaded {
			var bufferViewIndex int = int(accessors[a].bufferView)
			if bufferViewIndex < 0 || bufferViewIndex >= bufferViewCount {
				out = assetFailure(file, "invalid gltf : invalid buffer view index\n")
				return
			}

			var bufferView BufferView = bufferViews[bufferViewIndex]
			var bufferIndex int = int(bufferView.buffer)
			if bufferIndex < 0 || bufferIndex >= bufferCount {
				out = assetFailure(file, "invalid gltf : invalid buffer index\n")
				return
			}

			var dataF32 math.Vector_f32
			var dataUI16 math.Vector_ui16
			var dataUI32 math.Vector_ui32

			var buffer Buffer = buffers[bufferIndex]

			var uri string = fmt.Sprintf("%s%s", dataDir, buffer.uri)
			if uri != previousUri {
				if handle != -1 {
					if utils.Close(handle) == false {
						out = assetFailure(file, "invalid gltf : failed to close .bin\n")
						return
					}
				}
				handle = utils.Open(uri)
				previousUri = uri
			} else {
				var startOffset int64 = utils.Seek(handle, 0, utils.OS_SEEK_SET)
				if startOffset < 0 {
					out = assetFailure(file, "invalid gltf : failed to rewind file\n")
					return
				}
			}
			if handle >= 0 {
				var byteOffset int64 = accessors[a].byteOffset
				var bufferOffset int64 = bufferView.byteOffset + byteOffset

				if bufferOffset > 0 {
					bufferOffset = utils.Seek(handle, bufferOffset, utils.OS_SEEK_SET)
				}

				if bufferOffset < 0 {
					out = assetFailure(file, "invalid gltf : failed to seek\n")
					return
				}

				var count int32 = accessors[a].Count
				var componentCount int32 = int32(accessors[a].ComponentCount)
				var componentType int32 = accessors[a].ComponentType
				if componentType == gl.FLOAT {
					var success bool
					accessors[a].DataF32 = dataF32.Resize(uint64(count*componentCount), false)
					success = utils.ReadF32Slice(handle, accessors[a].DataF32, dataF32.Count())
					if success == false {
						out = assetFailure(file, "invalid gltf : failed to read float32 slice\n")
						return
					}
				} else if componentType == gl.UNSIGNED_SHORT {
					var success bool
					accessors[a].DataUI16 = dataUI16.Resize(uint64(count*componentCount), false)
					success = utils.ReadUI16Slice(handle, accessors[a].DataUI16, dataUI16.Count())
					if success == false {
						out = assetFailure(file, "invalid gltf : failed to read uint16 slice\n")
						return
					}
				} else if componentType == gl.UNSIGNED_INT {
					var success bool
					accessors[a].DataUI32 = dataUI32.Resize(uint64(count*componentCount), false)
					success = utils.ReadUI32Slice(handle, accessors[a].DataUI32, dataUI32.Count())
					if success == false {
						out = assetFailure(file, "invalid gltf : failed to read uint32 slice\n")
						return
					}
				} else {
					out = assetFailure(file, fmt.Sprintf("unhandled buffer format %d\n", componentType))
					_ /*var success bool*/ = utils.Close(handle)
					return
				}
			}
		}
	}

	if handle != -1 {
		if utils.Close(handle) == false {
			out = assetFailure(file, "invalid gltf : failed close .bin\n")
			return
		}
	}

	out = assetSuccess(file, out)
	return
}
