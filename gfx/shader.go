package gfx

import (
	"fmt"
	"io/ioutil"
	"skyfx/utils"
)

// // TODO : use glShaderSource char**string param instead of string concatenation

// Constants ...
const (
	SHADER_VERTEX   int64 = 0
	SHADER_PIXEL    int64 = 1
	SHADER_VARIANT  int64 = 8
	USE_TRANSFORM   int64 = 1 << 0
	USE_COLOR_UNI   int64 = 1 << 1
	USE_COLOR_VTX   int64 = 1 << 2
	USE_COLOR_MAP   int64 = 1 << 3
	USE_PBR_UNI     int64 = 1 << 4
	USE_PBR_MAP     int64 = 1 << 5
	USE_TANGENT     int64 = 1 << 6
	USE_NORMAL_VTX  int64 = 1 << 7
	USE_NORMAL_MAP  int64 = 1 << 8
	USE_TANGENT_VTX int64 = 1 << 9
	// TODO FIX LATER USE_OCCLUSION_MAP int64 = 1 << 10
	USE_EMISSIVE_UNI  int64 = 1 << 11
	USE_EMISSIVE_MAP  int64 = 1 << 12
	USE_OCCLUSION_MAP int64 = 1 << 13
	USE_DEBUG_A       int64 = 1 << 14
	USE_PASS_0        int64 = 1 << 15
	USE_PASS_1        int64 = 1 << 16
	USE_SKIN          int64 = 1 << 17
	USE_PARTICLE      int64 = 1 << 18
	USE_MERGE         int64 = 1 << 19
	USE_EDGE          int64 = 1 << 20
	USE_STENCIL       int64 = 1 << 21
	USE_LINEAR_ALBEDO int64 = 1 << 22
)

// Globals ...
var g_shaderPath string

var g_shaders []Shader
var g_sources []ShaderSource
var g_variants []string
var g_defines []ShaderDefine

// ShaderDefine ...
type ShaderDefine struct {
	flag  int64
	token string
}

// ShaderKey ...
type ShaderKey struct {
	value int64
}

// ShaderId ...
type ShaderId struct {
	shader int32
}

// ShaderSource ...
type ShaderSource struct {
	path     string
	filename string
	code     string
}

// Shader ...
type Shader struct {
	shaderType int64
	variant    int64
	flags      int64
	key        int64
	defines    string
	source     int32
	glsl       string
}

// ShaderDefineAdd ...
func ShaderDefineAdd(flag int64, token string) {
	var define ShaderDefine
	define.flag = flag
	define.token = token
	g_defines = append(g_defines, define)
}

// ShaderDefineGetStr ...
func ShaderDefineGetStr(flag int64) (out string) {
	var count int = len(g_defines)
	for i := 0; i < count; i++ { // TODO : use map
		if g_defines[i].flag == flag {
			out = g_defines[i].token
			return
		}
	}

	utils.PanicIf(true, "invalid flag")
	return
}

// ShaderDefineGetHeader ...
func ShaderDefineGetHeader(flags int64) (out string) {
	var glVersion int32 = GetGLVersion()

	if glVersion == GL_VERSION_DS_3_2 {
		out = "#version 330\n"
		out = out + "#define texture2D texture\n"
		out = out + "#define textureCube texture\n"
		out = out + "#define textureCubeLod textureLod\n"
	} else if glVersion == GL_VERSION_ES_3_1 {
		out = "#version 310 es\n"
		out = out + "precision highp float;\n"
		out = out + "#define texture2D texture\n"
		out = out + "#define textureCube texture\n"
		out = out + "#define textureCubeLod textureLod\n"
	} else {
		panic("Invalid gl version")
	}

	for i := 0; i < 32; i++ {
		var flag int64 = 1 << int64(i)
		if (flags & flag) == flag {
			out = out + ShaderDefineGetStr(flag)
		}
	}
	return
}

func shaderInit(dataDir string) {
	g_shaderPath = fmt.Sprintf("%sshaders/", dataDir)
	fmt.Printf(">> SHADER_INIT DATADIR %s, SHADER_PATH %s\n", dataDir, g_shaderPath)

	ShaderDefineAdd(USE_TRANSFORM, "#define USE_TRANSFORM\n")
	ShaderDefineAdd(USE_COLOR_UNI, "#define USE_COLOR_UNI\n")
	ShaderDefineAdd(USE_COLOR_VTX, "#define USE_COLOR_VTX\n")
	ShaderDefineAdd(USE_COLOR_MAP, "#define USE_COLOR_MAP\n")
	ShaderDefineAdd(USE_PBR_UNI, "#define USE_PBR_UNI\n")
	ShaderDefineAdd(USE_PBR_MAP, "#define USE_PBR_MAP\n")
	ShaderDefineAdd(USE_TANGENT, "#define USE_TANGENT\n")
	ShaderDefineAdd(USE_NORMAL_VTX, "#define USE_NORMAL_VTX\n")
	ShaderDefineAdd(USE_NORMAL_MAP, "#define USE_NORMAL_MAP\n")
	ShaderDefineAdd(USE_TANGENT_VTX, "#define USE_TANGENT_VTX\n")
	ShaderDefineAdd(USE_OCCLUSION_MAP, "#define USE_OCCLUSION_MAP\n")
	ShaderDefineAdd(USE_EMISSIVE_UNI, "#define USE_EMISSIVE_UNI\n")
	ShaderDefineAdd(USE_EMISSIVE_MAP, "#define USE_EMISSIVE_MAP\n")
	ShaderDefineAdd(USE_OCCLUSION_MAP, "#define USE_OCCLUSION_MAP\n")
	ShaderDefineAdd(USE_DEBUG_A, "#define USE_DEBUG_A\n")
	ShaderDefineAdd(USE_PASS_0, "#define USE_PASS_0\n")
	ShaderDefineAdd(USE_PASS_1, "#define USE_PASS_1\n")
	ShaderDefineAdd(USE_SKIN, "#define USE_SKIN\n")
	ShaderDefineAdd(USE_PARTICLE, "#define USE_PARTICLE\n")
	ShaderDefineAdd(USE_MERGE, "#define USE_MERGE\n")
	ShaderDefineAdd(USE_STENCIL, "#define USE_STENCIL\n")
	ShaderDefineAdd(USE_EDGE, "#define USE_EDGE\n")
	ShaderDefineAdd(USE_LINEAR_ALBEDO, "#define USE_LINEAR_ALBEDO\n")
}

// ShaderKeyClear ...
func ShaderKeyClear() (out ShaderKey) {
	out.value = 0
	return
}

// ShaderKeySet ...
func ShaderKeySet(key ShaderKey, flag int64, value bool) (out ShaderKey) {
	if value {
		out = ShaderKeyAdd(key, flag)
	} else {
		out = ShaderKeyRem(key, flag)
	}
	return
}

// ShaderKeyAdd ...
func ShaderKeyAdd(key ShaderKey, flag int64) (out ShaderKey) {
	key.value = key.value | flag
	out = key
	return
}

// ShaderKeyRem ...
func ShaderKeyRem(key ShaderKey, flag int64) (out ShaderKey) {
	key.value = key.value & (int64(-1) ^ flag)
	out = key
	return
}

// ShaderAddVariant ...
func ShaderAddVariant(filename string) (out int64) {
	var count int = len(g_variants)
	for i := 0; i < count; i++ {
		if g_variants[i] == filename {
			out = int64(i)
			return
		}
	}

	out = int64(count)
	g_variants = append(g_variants, filename)
	return
}

// ShaderGetVariant ...
func ShaderGetVariant(variant int64) (out string) {
	utils.PanicIfNot(variant >= 0 && variant < int64(len(g_variants)), "invalid variant")
	out = g_variants[variant]
	return
}

// ShaderSourceIsValid ...
func ShaderSourceIsValid(id int32) (out bool) {
	out = int(id) >= 0 && int(id) < len(g_sources)
	return
}

// ShaderSourceCreate ...
func ShaderSourceCreate(path string, filename string) (out int32) {
	out = -1
	var source ShaderSource
	source.path = path
	source.filename = filename

	var shaderPath string = fmt.Sprintf("%s%s", path, filename)

	bytes, err := ioutil.ReadFile(shaderPath)
	if err == nil {
		source.code = string(bytes)
		out = int32(len(g_sources))
		g_sources = append(g_sources, source)
	}
	return
}

// ShaderIsValid ...
func ShaderIsValid(id ShaderId) (out bool) {
	out = int(id.shader) >= 0 && int(id.shader) < len(g_shaders)
	return
}

// ShaderCreate ...
func ShaderCreate(shaderType int64, shaderVariant int64, shaderFlags int64, shaderKey int64) (out ShaderId) {
	var source int32 = -1

	var count int32 = int32(len(g_shaders))
	for i := int32(0); i < count; i++ { // TODO : use map
		if g_shaders[i].key == shaderKey {
			out.shader = i
			return
		} else if g_shaders[i].variant == shaderVariant {
			source = g_shaders[i].source
		}
	}

	out.shader = int32(len(g_shaders))
	var shader Shader
	shader.shaderType = shaderType
	shader.variant = shaderVariant
	shader.flags = shaderFlags
	shader.key = shaderKey

	shader.defines = ShaderDefineGetHeader(shaderFlags)

	if ShaderSourceIsValid(source) {
		shader.source = source
	} else {
		shader.source = ShaderSourceCreate(g_shaderPath, ShaderGetVariant(shaderVariant))
	}

	var common int32 = ShaderSourceCreate(g_shaderPath, "common.glsl")
	shader.glsl = shader.defines + g_sources[common].code
	shader.glsl = shader.glsl + g_sources[shader.source].code // ISSUE : can't use + operator
	/*fmt.Printf("//--------------------------------------------------------------------------------\n")
	fmt.Printf("// filename : %s\n", ShaderGetVariant(shaderVariant))
	fmt.Printf("// type : %d\n", shaderType)
	fmt.Printf("// variant : %d\n", shaderVariant)
	fmt.Printf("// flags : %d\n", shaderFlags)
	fmt.Printf("// key : %d\n", shaderKey)
	fmt.Printf("// glsl :\n%s\n", shader.glsl)*/

	g_shaders = append(g_shaders, shader)

	utils.PanicIfNot(ShaderIsValid(out), "invalid id")
	return
}

// ShaderGetGlsl ...
func ShaderGetGlsl(id ShaderId) (out string) {
	utils.PanicIfNot(ShaderIsValid(id), "invalid id")
	out = g_shaders[id.shader].glsl
	return
}
