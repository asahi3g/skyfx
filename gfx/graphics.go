package gfx

import (
	"fmt"
	v2 "skyfx/math/v2"
	v4 "skyfx/math/v4"
	"skyfx/utils"
	math "skyfx/math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	GL_MAX_COLOR_ATTACHMENTS = 4
)

var PINK math.V4 = math.V4{X: 1.0, Y: 0.0, Z: 1.0, W: 1.0}

var g_clearColor math.V4 = PINK

var g_glVersion int32
var g_drawBuffers [GL_MAX_COLOR_ATTACHMENTS]uint32 = [GL_MAX_COLOR_ATTACHMENTS]uint32{gl.NONE, gl.NONE, gl.NONE, gl.NONE}

var VertexLayout []VertexAttribute = []VertexAttribute{
	VertexAttributeCreate(ATTRIBUTE_POSITION, 3, gl.FLOAT),
	VertexAttributeCreate(ATTRIBUTE_COLOR, 4, gl.FLOAT),
	VertexAttributeCreate(ATTRIBUTE_TEXCOORD, 2, gl.FLOAT),
}

// import "gl"
// import "mat"
// import "os"
// import "v1"
// import "v2"

// // TODO : print gl extensions on init
// // TODO : refactor transition shader
// //
// // Constants ...
// var ALWAYS int32 = gl.ALWAYS
// var BACK int32 = gl.BACK
// var COLOR int32 = gl.COLOR
// var COLOR_BUFFER_BIT int32 = gl.COLOR_BUFFER_BIT
// var CLAMP_TO_EDGE int32 = gl.CLAMP_TO_EDGE
// var CCW int32 = gl.CCW
// var CW int32 = gl.CW
// var DEPTH int32 = gl.DEPTH
// var DEPTH_BUFFER_BIT int32 = gl.DEPTH_BUFFER_BIT
// var EQUAL int32 = gl.EQUAL
// var FLOAT int32 = gl.FLOAT
// var KEEP int32 = gl.KEEP
// var LESS int32 = gl.LESS
// var LINEAR int32 = gl.LINEAR
// var LINES int32 = gl.LINES
// var NONE int32 = gl.NONE
// var ONE int32 = gl.ONE
// var ONE_MINUS_SRC_ALPHA int32 = gl.ONE_MINUS_SRC_ALPHA
// var REPEAT int32 = gl.REPEAT
// var REPLACE int32 = gl.REPLACE
// var RGBA int32 = gl.RGBA
// var SRC_ALPHA int32 = gl.SRC_ALPHA
// var STENCIL int32 = gl.STENCIL
// var STENCIL_BUFFER_BIT int32 = gl.STENCIL_BUFFER_BIT
// var TEXTURE_2D int32 = gl.TEXTURE_2D
// var TEXTURE_CUBE_MAP int32 = gl.TEXTURE_CUBE_MAP
// var TRIANGLES int32 = gl.TRIANGLES
// var UNSIGNED_SHORT int32 = gl.UNSIGNED_SHORT
// var ZERO int32 = gl.ZERO

// Globals ...
var DEBUG_0 math.V4

const (
	GL_VERSION_NONE   int32 = 0
	GL_VERSION_DS_3_2 int32 = 1
	GL_VERSION_ES_3_1 int32 = 2
)

var g_sizeofF32 uint32 = 4
var g_sizeofI32 uint32 = 4
var g_sizeofUI32 uint32 = 4
var g_sizeofUI16 uint32 = 2

var g_colorMaskRed bool = true
var g_colorMaskGreen bool = true
var g_colorMaskBlue bool = true
var g_colorMaskAlpha bool = true

var g_cull bool = false
var g_frontFace uint32 = gl.CCW
var g_cullFace uint32 = gl.BACK

var g_blend bool = false
var g_srcColor uint32 = 0
var g_srcAlpha uint32 = 0
var g_dstColor uint32 = 0
var g_dstAlpha uint32 = 0

var g_depthTest bool = false
var g_depthFunc uint32 = gl.LESS
var g_depthMask bool = true

var g_stencilTest bool = false
var g_stencilWrite uint32 = 256
var g_stencilFunc uint32 = gl.ALWAYS
var g_stencilFuncRef int32 = 0
var g_stencilFuncMask uint32 = 256
var g_stencilOpFail uint32 = gl.KEEP
var g_stencilOpDepthFail uint32 = gl.KEEP
var g_stencilOpDepthPass uint32 = gl.KEEP

var g_clearDepth float64
var g_clearStencil int32

var g_textureCUBE uint32 = 0
var g_texture2D uint32 = 0

var g_framebuffer uint32 = 0

var gfx_width float32 = 0.0
var gfx_height float32 = 0.0
var ViewportSize math.V2 = v2.ZERO
var ViewportBounds math.V4 = v4.ZERO

var gfx_ratio_x float32 = 1.0
var gfx_ratio_y float32 = 1.0

var gfx_viewportX int32 = 0
var gfx_viewportY int32 = 0

var gfx_viewportWidth int32 = 0
var gfx_viewportHeight int32 = 0

var gfx_scissor bool = false
var gfx_scissorX int32 = 0
var gfx_scissorY int32 = 0
var gfx_scissorWidth int32 = 0
var gfx_scissorHeight int32 = 0

var g_tfxDefault TemplateId
var g_tfxFade TemplateId
var g_tfxMerge TemplateId
var g_tfxSky TemplateId
var g_tfxPbr TemplateId

var FxVertexColor2D EffectId = EffectInvalid()
var FxTexture2D EffectId = EffectInvalid()
var g_fxUniformColor3D EffectId = EffectInvalid()
var g_fxVertexColor3D EffectId = EffectInvalid()
var FxTexture3D EffectId = EffectInvalid()
var g_fxParticles EffectId = EffectInvalid()
var g_fxFade_0 EffectId = EffectInvalid()
var g_fxFade_1 EffectId = EffectInvalid()
var g_fxMerge EffectId = EffectInvalid()
var g_fxStencil EffectId = EffectInvalid()
var g_fxEdge EffectId = EffectInvalid()
var FxSky EffectId = EffectInvalid()

var SpLinearWrap SamplerState = SamplerStateCreate(gl.LINEAR_MIPMAP_LINEAR, gl.LINEAR, gl.REPEAT, gl.REPEAT, gl.REPEAT)
var g_linearClamp SamplerState = SamplerStateCreate(gl.LINEAR_MIPMAP_LINEAR, gl.LINEAR, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
var g_nearestClamp SamplerState = SamplerStateCreate(gl.NEAREST_MIPMAP_NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

var g_linear0Wrap SamplerState = SamplerStateCreate(gl.LINEAR, gl.LINEAR, gl.REPEAT, gl.REPEAT, gl.REPEAT)
var SpLinear0Clamp SamplerState = SamplerStateCreate(gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
var g_nearest0Clamp SamplerState = SamplerStateCreate(gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

var g_particleLayout []VertexAttribute

var g_fullscreenQuad MeshId

// Debug ...
//var g_debugMesh MeshId
var g_debugColorCount int32 = 8

//var g_debugColors [8]math.V4 // ISSUE ../..//src/gfx/graphics.cx:138: syntax error: unexpected PERIOD
var g_useDebugColors bool

// GlError ...
func GlError() bool {
	error := false
	for true {
		glerr := gl.GetError()
		if glerr != 0 {
			error = true
			if glerr == gl.INVALID_ENUM {
				fmt.Print("GL_INVALID_ENUM\n")
			} else if glerr == gl.INVALID_VALUE {
				fmt.Print("GL_INVALID_VALUE\n")
			} else if glerr == gl.INVALID_OPERATION {
				fmt.Print("GL_INVALID_OPERATION\n")
			} else if glerr == gl.STACK_OVERFLOW {
				fmt.Print("GL_STACK_OVERFLOW\n")
			} else if glerr == gl.STACK_UNDERFLOW {
				fmt.Print("GL_STACK_UNDERFLOW\n")
			} else if glerr == gl.OUT_OF_MEMORY {
				fmt.Print("GL_OUT_OF_MEMORY\n")
			} else if glerr != 0 {
				utils.PanicIfNot(false, fmt.Sprintf("invalid glError %d\n", glerr))
			}
		} else {
			break
		}
	}

	return error
}

// Init ...
func Init(width uint32, height uint32, dataDir string, glVersion int32) {
	glfw.Init()

	g_glVersion = glVersion

	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_POSITION, 3, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_COLOR, 4, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD, 2, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_1, 3, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_2, 4, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_3, 4, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_4, 3, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_5, 3, FLOAT))
	// g_particleLayout = append(g_particleLayout, VertexAttributeCreate(ATTRIBUTE_TEXCOORD_6, 4, FLOAT))

	// /*g_debugColors[0] = v4.RED
	// g_debugColors[1] = v4.GREEN
	// g_debugColors[2] = v4.BLUE
	// g_debugColors[3] = v4.PINK
	// g_debugColors[4] = v4.YELLOW
	// g_debugColors[5] = v4.SKY
	// g_debugColors[6] = v4.ONE
	// g_debugColors[7] = v4.BLACK*/

	gl.Init()
	utils.PanicIf(GlError(), "gl.Init")

	gl.Disable(gl.BLEND)
	utils.PanicIf(GlError(), "gl.Disable(gl.BLEND)")

	gl.Disable(gl.CULL_FACE)
	utils.PanicIf(GlError(), "gl.Disable(gl.CULL_FACE)")

	gl.Disable(gl.DEPTH_TEST)
	utils.PanicIf(GlError(), "gl.Disable(gl.DEPTH_TEST)")

	gl.ClearColor(g_clearColor.X, g_clearColor.Y, g_clearColor.Z, g_clearColor.W)
	utils.PanicIf(GlError(), "gl.ClearColor")

	resizeViewport(width, height)

	shaderInit(dataDir)

	g_tfxDefault = TemplateCreate("fxTemplateDefault", "default.vsh", "default.psh")
	g_tfxFade = TemplateCreate("fxTemplateFade", "default.vsh", "fade.psh")
	g_tfxMerge = TemplateCreate("fxTemplateMerge", "default.vsh", "merge.psh")
	g_tfxSky = TemplateCreate("fxTemplateSky", "default.vsh", "sky.psh")
	g_tfxPbr = TemplateCreate("fxTemplatePbr", "default.vsh", "default.psh")

	// default
	TemplateReset(g_tfxDefault)
	TemplateSetKey(g_tfxDefault, USE_COLOR_VTX, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_COLOR, "i_v4Color")
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	FxVertexColor2D = TemplateInstance(g_tfxDefault)

	TemplateReset(g_tfxDefault)
	TemplateSetKey(g_tfxDefault, USE_COLOR_VTX, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_MAP, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_COLOR, "i_v4Color")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxDefault, SAMPLER_COLOR_0, "u_t2Color", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	FxTexture2D = TemplateInstance(g_tfxDefault)

	TemplateReset(g_tfxDefault)
	TemplateSetVertexKey(g_tfxDefault, USE_TRANSFORM, true)
	TemplateSetPixelKey(g_tfxDefault, USE_COLOR_UNI, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxDefault, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxDefault, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxDefault, UNIFORM_PROJECTION, "u_m44Projection")
	TemplateBindUniform(g_tfxDefault, UNIFORM_COLOR, "u_v4Color")
	g_fxUniformColor3D = TemplateInstance(g_tfxDefault)

	TemplateReset(g_tfxDefault)
	TemplateSetVertexKey(g_tfxDefault, USE_TRANSFORM, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_VTX, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_COLOR, "i_v4Color")
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxDefault, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxDefault, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxDefault, UNIFORM_PROJECTION, "u_m44Projection")
	g_fxVertexColor3D = TemplateInstance(g_tfxDefault)

	TemplateReset(g_tfxDefault)
	TemplateSetVertexKey(g_tfxDefault, USE_TRANSFORM, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_VTX, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_MAP, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_COLOR, "i_v4Color")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxDefault, SAMPLER_COLOR_0, "u_t2Color", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxDefault, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxDefault, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxDefault, UNIFORM_PROJECTION, "u_m44Projection")
	FxTexture3D = TemplateInstance(g_tfxDefault)

	// particle
	TemplateReset(g_tfxDefault)
	TemplateSetVertexKey(g_tfxDefault, USE_TRANSFORM, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_UNI, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_VTX, true)
	TemplateSetKey(g_tfxDefault, USE_COLOR_MAP, true)
	TemplateSetKey(g_tfxDefault, USE_PARTICLE, true)
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_COLOR, "i_v4Color")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_1, "i_v3Velocity")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_2, "i_v4Orientation")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_3, "i_v4AngularVelocity")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_4, "i_v3Scale")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_5, "i_v3ScaleVelocity")
	TemplateBindAttribute(g_tfxDefault, ATTRIBUTE_TEXCOORD_6, "i_v4Particle")
	TemplateBindSampler(g_tfxDefault, SAMPLER_COLOR_0, "u_t2Color", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxDefault, SAMPLER_COLOR_1, "u_t2Depth", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxDefault, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxDefault, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxDefault, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxDefault, UNIFORM_COLOR, "u_v4Color")
	TemplateBindUniform(g_tfxDefault, UNIFORM_PARTICLE, "u_v4Particle")
	TemplateBindUniform(g_tfxDefault, UNIFORM_PROJECTION, "u_m44Projection")
	TemplateBindUniform(g_tfxDefault, UNIFORM_TARGET_SIZE, "u_v4TargetSize")
	g_fxParticles = TemplateInstance(g_tfxDefault)

	// fade
	TemplateReset(g_tfxFade)
	TemplateSetKey(g_tfxFade, USE_COLOR_MAP, true)
	TemplateSetPixelKey(g_tfxFade, USE_PASS_0, true)
	TemplateBindAttribute(g_tfxFade, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxFade, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxFade, SAMPLER_COLOR_0, "u_t2Src", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxFade, SAMPLER_COLOR_1, "u_t2Dst", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxFade, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxFade, UNIFORM_TIME, "u_fTime")
	g_fxFade_0 = TemplateInstance(g_tfxFade)

	TemplateReset(g_tfxFade)
	TemplateSetKey(g_tfxFade, USE_COLOR_MAP, true)
	TemplateSetPixelKey(g_tfxFade, USE_PASS_1, true)
	TemplateBindAttribute(g_tfxFade, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxFade, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxFade, SAMPLER_COLOR_0, "u_t2Src", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxFade, SAMPLER_COLOR_1, "u_t2Dst", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxFade, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxFade, UNIFORM_TIME, "u_fTime")
	g_fxFade_1 = TemplateInstance(g_tfxFade)

	// merge
	TemplateReset(g_tfxMerge)
	TemplateSetKey(g_tfxMerge, USE_COLOR_MAP, true)
	TemplateSetKey(g_tfxMerge, USE_MERGE, true)
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_0, "u_t2Opaque", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_1, "u_t2Alpha0", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_2, "u_t2Alpha1", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_3, "u_t2Depth", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxMerge, UNIFORM_DEBUG_0, "DEBUG_0")
	g_fxMerge = TemplateInstance(g_tfxMerge)

	// edge
	TemplateReset(g_tfxMerge)
	TemplateSetKey(g_tfxMerge, USE_COLOR_MAP, true)
	TemplateSetKey(g_tfxMerge, USE_EDGE, true)
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_3, "u_t2Depth", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxMerge, UNIFORM_TARGET_SIZE, "u_v4TargetSize")
	TemplateBindUniform(g_tfxMerge, UNIFORM_DEBUG_0, "DEBUG_0")
	g_fxEdge = TemplateInstance(g_tfxMerge)

	// stencil
	TemplateReset(g_tfxMerge)
	TemplateSetKey(g_tfxMerge, USE_COLOR_MAP, true)
	TemplateSetKey(g_tfxMerge, USE_STENCIL, true)
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxMerge, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindSampler(g_tfxMerge, SAMPLER_COLOR_3, "u_t2Depth", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxMerge, UNIFORM_DEBUG_0, "DEBUG_0")
	g_fxStencil = TemplateInstance(g_tfxMerge)

	// sky
	TemplateReset(g_tfxSky)
	TemplateSetVertexKey(g_tfxSky, USE_TRANSFORM, true)
	TemplateSetPixelKey(g_tfxSky, USE_LINEAR_ALBEDO, true)
	TemplateBindAttribute(g_tfxSky, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindSampler(g_tfxSky, SAMPLER_ENV_DIFFUSE, "u_tcEnvironmentDiffuse", gl.TEXTURE_CUBE_MAP)
	TemplateBindUniform(g_tfxSky, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxSky, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxSky, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxSky, UNIFORM_PROJECTION, "u_m44Projection")
	TemplateBindUniform(g_tfxSky, UNIFORM_PBR, "u_v4PBR")
	TemplateBindUniform(g_tfxSky, UNIFORM_CAMERA_POSITION, "u_v4CameraPosition")
	FxSky = TemplateInstance(g_tfxSky)

	// pbr
	TemplateReset(g_tfxPbr)
	TemplateSetVertexKey(g_tfxPbr, USE_TRANSFORM, true)
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_POSITION, "i_v3Position")
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_NORMAL, "i_v3Normal")
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_TEXCOORD, "i_v2Texcoord")
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_TANGENT, "i_v4Tangent")
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_WEIGHT, "i_v4Weight")
	TemplateBindAttribute(g_tfxPbr, ATTRIBUTE_JOINT, "i_v4Joint")
	TemplateBindSampler(g_tfxPbr, SAMPLER_COLOR_0, "u_t2Color", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxPbr, SAMPLER_NORMAL, "u_t2Normal", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxPbr, SAMPLER_METAL_ROUGH, "u_t2MetalRough", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxPbr, SAMPLER_BRDF, "u_t2BRDF", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxPbr, SAMPLER_EMISSIVE, "u_t2Emissive", gl.TEXTURE_2D)
	TemplateBindSampler(g_tfxPbr, SAMPLER_ENV_SPECULAR, "u_tcEnvironmentSpecular", gl.TEXTURE_CUBE_MAP)
	TemplateBindSampler(g_tfxPbr, SAMPLER_ENV_DIFFUSE, "u_tcEnvironmentDiffuse", gl.TEXTURE_CUBE_MAP)
	TemplateBindSampler(g_tfxPbr, SAMPLER_OCCLUSION, "u_t2Occlusion", gl.TEXTURE_2D)
	TemplateBindUniform(g_tfxPbr, UNIFORM_DEBUG_0, "DEBUG_0")
	TemplateBindUniform(g_tfxPbr, UNIFORM_WORLD, "u_m44World")
	TemplateBindUniform(g_tfxPbr, UNIFORM_VIEW, "u_m44View")
	TemplateBindUniform(g_tfxPbr, UNIFORM_PROJECTION, "u_m44Projection")
	TemplateBindUniform(g_tfxPbr, UNIFORM_WORLD_INVERSE, "u_m44WorldInverse")
	TemplateBindUniform(g_tfxPbr, UNIFORM_CAMERA_POSITION, "u_v4CameraPosition")
	TemplateBindUniform(g_tfxPbr, UNIFORM_COLOR, "u_v4Color")
	TemplateBindUniform(g_tfxPbr, UNIFORM_METAL_ROUGH, "u_v4MetalRough")
	TemplateBindUniform(g_tfxPbr, UNIFORM_EMISSIVE, "u_v4Emissive")
	TemplateBindUniform(g_tfxPbr, UNIFORM_SKELETON, "u_m44Skeleton")
	TemplateBindUniform(g_tfxPbr, UNIFORM_PBR, "u_v4PBR")

	g_fullscreenQuad = MeshCreate(gl.TRIANGLES, gl.UNSIGNED_SHORT, 6*3, VertexLayout, 6*3)
	resizeFullscreenQuad(width, height)
}

func resizeFullscreenQuad(width uint32, height uint32) {
	MeshBegin(g_fullscreenQuad)
	MeshAppendOrthoQuad(g_fullscreenQuad, ViewportBounds, v4.Make(0.0, 1.0, 1.0, 0.0), v4.ONE, ViewportBounds, 0.0)
	MeshEnd(g_fullscreenQuad)
}

// NoGlError ...
func NoGlErrror() (noerror bool) {
	var error bool = GlError()
	noerror = error == false // panic if noerror = (GlError() == false)
	return
}

// GlStatus ...
func GlStatus() (error bool) {
	error = false
	var glstat uint32 = gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	utils.PanicIf(GlError(), "CheckFramebufferStatus")
	if glstat != gl.FRAMEBUFFER_COMPLETE {
		error = true
		if glstat == gl.FRAMEBUFFER_UNDEFINED {
			fmt.Printf("GL_FRAMEBUFFER_UNDEFINED\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_READ_BUFFER\n")
		} else if glstat == gl.FRAMEBUFFER_UNSUPPORTED {
			fmt.Printf("GL_FRAMEBUFFER_UNSUPPORTED\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_MULTISAMPLE\n")
		} else if glstat == gl.FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS {
			fmt.Printf("GL_FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS\n")
		} else {
			utils.PanicIfNot(false, "invalid glFramebufferStatus")
		}
	}
	return
}

// GlString ...
func GlString(value uint32) (out string) {
	if value == gl.RGBA {
		out = "RGBA"
	} else if value == gl.DEPTH_COMPONENT {
		out = "DEPTH_COMPONENT"
	} else if value == gl.DEPTH_COMPONENT24 {
		out = "DEPTH_COMPONENT24"
	} else if value == gl.UNSIGNED_BYTE {
		out = "UNSIGNED_BYTE"
	} else if value == gl.TEXTURE_2D {
		out = "TEXTURE_2D"
	} else if value == gl.TEXTURE_CUBE_MAP {
		out = "TEXTURE_CUBE_MAP"
	} else if value == gl.RGBA8 {
		out = "RGBA8"
	} else {
		out = "unknown gl enum"
	}
	return
}

// NoGlStatus ...
func NoGlStatus() (nostatus bool) { // ISSUE : typo in function name GlStatus instaed of NoGlStatus : error: /home/reash/go/src/github.com/skycoin/cx/lib/cxfx/src/graphics.cx:123 identifier 'error' does not exist
	var status bool = GlStatus()
	nostatus = status == false // panic if nostatus = (GlStatus() == false)
	return
}

// Location ...
type Location struct {
	location uint32
	Type     uint32
	name     string
	bound    bool
}

func (this *Location) IsValid() bool {
	return this.location != math.MAX_ui32
}

// LocationBind ...
func LocationBind(locations []Location, location uint32, slot uint32, name string, Type uint32) (out []Location) {
	var count int32 = int32(len(locations))
	for i := count; i <= int32(slot); i++ {
		locations = append(locations, LocationCreate(math.MAX_ui32, "", math.MAX_ui32))
	}

	locations[slot] = LocationCreate(location, name, Type)
	out = locations
	return
}

// LocationCreate ...
func LocationCreate(location uint32, name string, Type uint32) (out Location) {
	out.location = location
	out.name = name
	out.Type = Type
	out.bound = false
	return
}

func assignTexture(target uint32, slot uint32, sampler uint32, name uint32) {
	gl.Uniform1i(int32(sampler), int32(slot))
	utils.PanicIf(GlError(), "gl.Uniform1i")

	gl.ActiveTexture(gl.TEXTURE0 + slot)
	utils.PanicIf(GlError(), "gl.ActiveTexture")

	bindTexture(target, name)
}

func assignFloat(slot uint32, value float32) {
	gl.Uniform1f(int32(slot), value)
	utils.PanicIf(GlError(), "gl.Uniform1f")
}

func assignVector4(slot uint32, value math.V4) {
	gl.Uniform4fv(int32(slot), 1, &value.X)
	utils.PanicIf(GlError(), "gl.UniformV4F")
}

func assignMatrix4(slot uint32, value math.M44, transpose bool) {
	gl.UniformMatrix4fv(int32(slot), 1, transpose, &value.V00)
	utils.PanicIf(GlError(), "gl.UniformM44F")
}

func assignMatrix4V(slot uint32, value []math.M44, transpose bool) {
	gl.UniformMatrix4fv(int32(slot), int32(len(value)), transpose, &value[0].V00)
	utils.PanicIf(GlError(), "gl.UniformM44FV")
}

func bindTexture(target uint32, name uint32) {
	if target == gl.TEXTURE_2D {
		//if g_texture2D != name  { // TODO : fix cache
		gl.BindTexture(gl.TEXTURE_2D, name)
		utils.PanicIf(GlError(), "gl.BindTexture")
		g_texture2D = name
		//}
	} else if target == gl.TEXTURE_CUBE_MAP {
		//if g_textureCUBE != name  { // TODO : fix cache
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, name)
		utils.PanicIf(GlError(), "gl.BindTexture")
		g_textureCUBE = name
		//}
	}
}

func bindFramebuffer(name uint32) {
	if g_framebuffer != name {
		gl.BindFramebuffer(gl.FRAMEBUFFER, name)
		utils.PanicIf(GlError(), "gl.BindFramebuffer")
	}
}

// DrawBuffers ...
func DrawBuffers(count int32) {
	for i := 0; i < len(g_drawBuffers); i++ {
		if i < int(count) {
			g_drawBuffers[i] = uint32(gl.COLOR_ATTACHMENT0 + i)
		} else {
			g_drawBuffers[i] = uint32(gl.NONE)
		}
	}
	gl.DrawBuffers(count, &g_drawBuffers[0])
	utils.PanicIf(GlError(), "gl.DrawBuffers")
}

// EnableCulling ...
func EnableCulling(frontFace uint32, cullFace uint32) {
	if g_cull == false {
		g_cull = true
		gl.Enable(gl.CULL_FACE)
		utils.PanicIf(GlError(), "gl.Enable(gl.CULL_FACE)")
	}

	if g_frontFace != frontFace {
		g_frontFace = frontFace
		gl.FrontFace(g_frontFace)
		utils.PanicIf(GlError(), "gl.FrontFace")
	}

	if g_cullFace != cullFace {
		g_cullFace = cullFace
		gl.CullFace(g_cullFace)
		utils.PanicIf(GlError(), "gl.CullFace")
	}
}

// DisableCulling ...
func DisableCulling() {
	if g_cull == true {
		g_cull = false
		gl.Disable(gl.CULL_FACE)
		utils.PanicIf(GlError(), "gl.Disable(gl.CULL_FACE")
	}
}

// EnableBlending ...
func EnableBlending(srcBlend uint32, dstBlend uint32) {
	if g_blend == false {
		g_blend = true
		gl.Enable(gl.BLEND)
		utils.PanicIf(GlError(), "gl.Enable(gl.BLEND)")
	}

	if (g_srcColor != srcBlend) || (g_dstColor != dstBlend) || (g_srcAlpha != srcBlend) || (g_dstAlpha != dstBlend) {
		g_srcColor = srcBlend
		g_srcAlpha = srcBlend
		g_dstColor = dstBlend
		g_dstAlpha = dstBlend
		gl.BlendFunc(srcBlend, dstBlend)
		utils.PanicIf(GlError(), "gl.BlendFunc")
	}
}

// EnableBlendingSeparate ...
func EnableBlendingSeparate(srcColor uint32, dstColor uint32, srcAlpha uint32, dstAlpha uint32) {
	if g_blend == false {
		g_blend = true
		gl.Enable(gl.BLEND)
		utils.PanicIf(GlError(), "gl.Enable(gl.BLEND)")
	}

	if (g_srcColor != srcColor) || (g_dstColor != dstColor) || (g_srcAlpha != srcAlpha) || (g_dstAlpha != dstAlpha) {
		g_srcColor = srcColor
		g_dstColor = dstColor
		g_srcAlpha = srcAlpha
		g_dstAlpha = dstAlpha
		gl.BlendFuncSeparate(g_srcColor, g_dstColor, g_srcAlpha, g_dstAlpha)
		utils.PanicIf(GlError(), "gl.BlendFuncSeparate")
	}
}

// DisableBlending ...
func DisableBlending() {
	if g_blend == true {
		g_blend = false
		gl.Disable(gl.BLEND)
		utils.PanicIf(GlError(), "gl.Disable(gl.BLEND)")
	}
}

// DepthState ...
func DepthState(test bool, compare uint32, write bool) {
	DepthTest(test)
	DepthFunc(compare)
	DepthWrite(write)
}

// DepthTest ...
func DepthTest(value bool) {
	if value != g_depthTest {
		if value {
			gl.Enable(gl.DEPTH_TEST)
			utils.PanicIf(GlError(), "gl.Enable(gl.DEPTH_TEST)")
		} else {
			gl.Disable(gl.DEPTH_TEST)
			utils.PanicIf(GlError(), "gl.Disable(gl.DEPTH_TEST)")
		}
		g_depthTest = value
	}
}

// EnableDepth ...
func EnableDepth() {
	if g_depthTest == false {
		g_depthTest = true
		gl.Enable(gl.DEPTH_TEST)
	}
}

// DisableDepth ...
func DisableDepth() {
	DepthTest(false)
	DepthWrite(false)
}

// DepthFunc ...
func DepthFunc(value uint32) {
	if g_depthFunc != value {
		g_depthFunc = value
		gl.DepthFunc(g_depthFunc)
		utils.PanicIf(GlError(), "gl.DepthFunc")
	}
}

// DepthWrite ...
func DepthWrite(value bool) {
	if g_depthMask != value {
		g_depthMask = value
		gl.DepthMask(g_depthMask)
		utils.PanicIf(GlError(), "gl.DepthWrite")
	}
}

// StencilState ...
func StencilState(test bool, compare uint32, ref int32, mask uint32, sfail uint32, dpfail uint32, dppass uint32, write uint32) {
	StencilTest(test)
	StencilFunc(compare, ref, mask)
	StencilOp(sfail, dpfail, dppass)
	StencilWrite(write)
}

// StencilTest ...
func StencilTest(value bool) {
	if value != g_stencilTest {
		if value {
			gl.Enable(gl.STENCIL_TEST)
			utils.PanicIf(GlError(), "gl.Enable(gl.STENCIL_TEST)")
		} else {
			gl.Disable(gl.STENCIL_TEST)
			utils.PanicIf(GlError(), "gl.Disable(gl.STENCIL_TEST")
		}
		g_stencilTest = value
	}
}

// DisableStencil ...
func DisableStencil() {
	StencilTest(false)
	StencilWrite(0)
}

// StencilFunc ...
func StencilFunc(compare uint32, ref int32, mask uint32) {
	if g_stencilFunc != compare ||
		g_stencilFuncRef != ref ||
		g_stencilFuncMask != mask {
		g_stencilFunc = compare
		g_stencilFuncRef = ref
		g_stencilFuncMask = mask
		gl.StencilFunc(g_stencilFunc, g_stencilFuncRef, g_stencilFuncMask)
		utils.PanicIf(GlError(), "gl.StencilFunc")
	}
}

// StencilOp ...
func StencilOp(sfail uint32, dpfail uint32, dppass uint32) {
	if g_stencilOpFail != sfail ||
		g_stencilOpDepthFail != dpfail ||
		g_stencilOpDepthPass != dppass {
		g_stencilOpFail = sfail
		g_stencilOpDepthFail = dpfail
		g_stencilOpDepthPass = dppass
		gl.StencilOp(g_stencilOpFail, g_stencilOpDepthFail, g_stencilOpDepthPass)
		utils.PanicIf(GlError(), "gl.StencilOp")
	}
}

// StencilWrite ...
func StencilWrite(write uint32) {
	if g_stencilWrite != write {
		g_stencilWrite = write
		gl.StencilMask(g_stencilWrite)
		utils.PanicIf(GlError(), "gl.StencilMask")
	}
}

// ClearBufferI ...
func ClearBufferI(buffer uint32, index int32, r int32, g int32, b int32, a int32) {
	var colors [4]int32 = [4]int32{r, g, b, a}
	gl.ClearBufferiv(buffer, index, &colors[0])
	utils.PanicIf(GlError(), "gl.ClearBufferI")
}

// ClearBufferI ...
func ClearBufferUI(buffer uint32, index int32, r uint32, g uint32, b uint32, a uint32) {
	var colors [4]uint32 = [4]uint32{r, g, b, a}
	gl.ClearBufferuiv(buffer, index, &colors[0])
	utils.PanicIf(GlError(), "gl.ClearBufferUI")
}

// ClearBufferF ...
func ClearBufferF(buffer uint32, index int32, r float32, g float32, b float32, a float32) {
	var colors [4]float32 = [4]float32{r, g, b, a}
	gl.ClearBufferfv(buffer, index, &colors[0])
	utils.PanicIf(GlError(), "gl.ClearBufferF")
}

// Clear ...
func Clear(buffers uint32, color math.V4, depth float64, stencil int32) {
	if (g_clearColor.X != color.X) ||
		(g_clearColor.Y != color.Y) ||
		(g_clearColor.Z != color.Z) ||
		(g_clearColor.W != color.W) {
		g_clearColor = color
		gl.ClearColor(g_clearColor.X, g_clearColor.Y, g_clearColor.Z, g_clearColor.W)
		utils.PanicIf(GlError(), "gl.ClearColor")
	}
	if g_clearDepth != depth {
		g_clearDepth = depth
		gl.ClearDepth(g_clearDepth)
	}
	if g_clearStencil != stencil {
		g_clearStencil = stencil
		gl.ClearStencil(g_clearStencil)
	}
	gl.Clear(buffers)
	utils.PanicIf(GlError(), "gl.Clear")
}

// ColorWrite ...
func ColorWrite(red bool, green bool, blue bool, alpha bool) {
	if g_colorMaskRed != red ||
		g_colorMaskGreen != green ||
		g_colorMaskBlue != blue ||
		g_colorMaskAlpha != alpha {
		g_colorMaskRed = red
		g_colorMaskGreen = green
		g_colorMaskBlue = blue
		g_colorMaskAlpha = alpha
		gl.ColorMask(g_colorMaskRed, g_colorMaskGreen, g_colorMaskBlue, g_colorMaskAlpha)
		utils.PanicIf(GlError(), "gl.ColorMask")
	}
}

// SetViewport ...
func SetViewport(bounds math.V4) {

	var x int32 = int32(bounds.X)
	var y int32 = int32(bounds.Y)
	var width int32 = int32(bounds.Z)
	var height int32 = int32(bounds.W)

	//fmt.Printf("SET_VIEWPORT %f, %f, %f, %f\n", bounds.X, bounds.Y, bounds.Z, bounds.W)
	if gfx_viewportX != x ||
		gfx_viewportY != y ||
		gfx_viewportWidth != width ||
		gfx_viewportHeight != height {
		gfx_viewportX = x
		gfx_viewportY = y
		gfx_viewportWidth = width
		gfx_viewportHeight = height
		gl.Viewport(x, y, width, height)
		utils.PanicIf(GlError(), "gl.Viewport")
	}
}

// SetScissor ...
func SetScissor(bounds math.V4) {
	if gfx_scissor == false {
		gfx_scissor = true
		gl.Enable(gl.SCISSOR_TEST)
		utils.PanicIf(GlError(), "gl.Scissor")
	}

	var x int32 = int32(bounds.X)
	var y int32 = int32(bounds.Y)
	var width int32 = int32(bounds.Z)
	var height int32 = int32(bounds.W)

	if gfx_scissorX != x ||
		gfx_scissorY != y ||
		gfx_scissorWidth != width ||
		gfx_scissorHeight != height {
		gfx_scissorX = x
		gfx_scissorY = y
		gfx_scissorWidth = width
		gfx_scissorHeight = height
		gl.Scissor(x, y, width, height)
		utils.PanicIf(GlError(), "gl.Scissor")
	}
}

func GetGLVersion() (out int32) {
	out = g_glVersion
	return
}

func GetGLVersionFromString(glVersion string) (out int32) {
	if glVersion == "gl32" {
		out = GL_VERSION_DS_3_2
	} else if glVersion == "gles31" {
		out = GL_VERSION_ES_3_1
	}
	return
}

func GetGLVersionToString(glVersion int32) (out string) {
	if glVersion == GL_VERSION_DS_3_2 {
		out = "gl32"
	} else if glVersion == GL_VERSION_ES_3_1 {
		out = "gles31"
	} else {
		out = "unsupported"
	}
	return
}

// Resize ...
func Resize(width uint32, height uint32) {
	resizeFullscreenQuad(width, height)
	resizeViewport(width, height)
	ResizeTargets(width, height)
	//ResizeTextures(width, height)
}

func resizeViewport(width uint32, height uint32) {

	gfx_width = float32(width)
	gfx_height = float32(height)
	ViewportSize = v2.Make(gfx_width, gfx_height)
	ViewportBounds = v4.Make(0.0, 0.0, gfx_width, gfx_height)
	gfx_ratio_y = gfx_width / gfx_height
	gfx_ratio_x = gfx_height / gfx_width
	SetViewport(ViewportBounds)
	SetScissor(ViewportBounds)
}

// Destroy ...
func Destroy() {

	//PopScissor()

	ProgramUse(NullProgram())

	// TODO FIX LATER
	// var i int32
	// for i = 0; i < len(g_vbos); i = i + 1 {
	// 	//fmt.Printf("gl.DeleteBuffers(%d)\n", g_vbos[i])
	// 	gl.DeleteBuffers(1, g_vbos[i])
	// 	utils.PanicIf(GlError(), "glDeleteBuffers")
	// }

	// for i = 0; i < len(g_vaos); i = i + 1 {
	// 	//fmt.Printf("gl.DeleteVertexArrays(%d)\n", g_vaos[i])
	// 	gl.DeleteVertexArrays(1, g_vaos[i])
	// 	utils.PanicIf(GlError(), "g_DeleteVertexArrays")
	// }

	// DestroyTargets()
	// DestroyTextures()

	// for i = 0; i < len(g_programs); i = i + 1 {
	// 	fmt.Printf("gl.DeleteProgram(%d)\n", g_programs[i])
	// 	gl.DeleteProgram(g_programs[i]) // ##0 crash
	// 	utils.PanicIf(GlError(), "gl.DeleteProgram")
	// }

	// gl.Destroy()
}
