package gfx

import (
	"skyfx/math"
	"skyfx/utils"
)

// Constants
const (
	ATTRIBUTE_POSITION = iota
	ATTRIBUTE_COLOR
	ATTRIBUTE_TEXCOORD
	ATTRIBUTE_NORMAL
	ATTRIBUTE_TANGENT
	ATTRIBUTE_WEIGHT
	ATTRIBUTE_JOINT
)

const (
	ATTRIBUTE_TEXCOORD_1 = ATTRIBUTE_TEXCOORD + 1
	ATTRIBUTE_TEXCOORD_2
	ATTRIBUTE_TEXCOORD_3
	ATTRIBUTE_TEXCOORD_4
	ATTRIBUTE_TEXCOORD_5
	ATTRIBUTE_TEXCOORD_6
)

const (
	SAMPLER_COLOR_0 = iota
	SAMPLER_NORMAL
	SAMPLER_METAL_ROUGH
	SAMPLER_ENV_SPECULAR
	SAMPLER_ENV_DIFFUSE
	SAMPLER_BRDF
	SAMPLER_EMISSIVE
	SAMPLER_OCCLUSION
	SAMPLER_COLOR_1
	SAMPLER_COLOR_2
	SAMPLER_COLOR_3
)

const (
	UNIFORM_DEBUG_0 = iota
	UNIFORM_TIME
	UNIFORM_WORLD
	UNIFORM_VIEW
	UNIFORM_PROJECTION
	UNIFORM_WORLD_INVERSE
	UNIFORM_CAMERA_POSITION
	UNIFORM_COLOR
	UNIFORM_METAL_ROUGH
	UNIFORM_EMISSIVE
	UNIFORM_SKELETON
	UNIFORM_PARTICLE
	UNIFORM_TARGET_SIZE
	UNIFORM_PBR
)

// Globals ...
var g_templates []Template // TODO : use map
var g_effects []Effect     // TODO : use map

// EffectId ...
type EffectId struct {
	effect int32
}

// Effect ...
type Effect struct {
	key        int64
	program    Program
	attributes []Location
	samplers   []Location
	uniforms   []Location
}

// EffectIsValid ...
func EffectIsValid(id EffectId) (out bool) {
	out = int(id.effect) >= 0 && int(id.effect) < len(g_effects)
	return
}

// EffectInvalid ...
func EffectInvalid() (out EffectId) {
	out.effect = -1
	return
}

// EffectCreate ...
func EffectCreate(key int64, attributes []Location, samplers []Location, uniforms []Location, vertexShader ShaderId, pixelShader ShaderId) (out EffectId) {
	var effect Effect
	effect.key = key

	var attr []Location = effect.attributes
	attr = appendLocations(attr, attributes)
	effect.attributes = attr

	var samp []Location = effect.samplers
	samp = appendLocations(samp, samplers)
	effect.samplers = samp

	var unif []Location = effect.uniforms
	unif = appendLocations(unif, uniforms)
	effect.uniforms = unif

	effect.program = ProgramCreate(attr, samp, unif, ShaderGetGlsl(vertexShader), ShaderGetGlsl(pixelShader))
	out.effect = int32(len(g_effects))
	g_effects = append(g_effects, effect)
	return
}

// EffectGetKey ...
func EffectGetKey(id EffectId) (out int64) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	out = g_effects[id.effect].key
	return
}

// EffectUse ...
func EffectUse(id EffectId) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	ProgramUse(g_effects[id.effect].program)
	return
}

// EffectIsValidSamplerSlot ...
func EffectIsValidSamplerSlot(id EffectId, slot uint32) (out bool) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	out = int(slot) >= 0 && int(slot) < len(g_effects[id.effect].samplers)
	return
}

// EffectIsValidUniformSlot ...
func EffectIsValidUniformSlot(id EffectId, slot int32) (out bool) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	out = int(slot) >= 0 && int(slot) < len(g_effects[id.effect].uniforms)
	return
}

// EffectIsValidUniformLocation ...
func EffectIsValidUniformLocation(id EffectId, slot int32) (out bool) {
	out = false
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	if EffectIsValidUniformSlot(id, slot) {
		out = g_effects[id.effect].uniforms[slot].IsValid()
	}
	return
}

// EffectTryAssignTexture ...
func EffectTryAssignTexture(id EffectId, slot uint32, texture TextureId, sampler SamplerState) {
	if EffectIsValidSamplerSlot(id, slot) {
		EffectAssignTexture(id, slot, texture, sampler)
	}
}

// EffectAssignTexture ...
func EffectAssignTexture(id EffectId, slot uint32, texture TextureId, sampler SamplerState) {
	utils.PanicIfNot(EffectIsValid(id), "invalid effect")
	utils.PanicIfNot(EffectIsValidSamplerSlot(id, slot), "invalid sampler slot")
	//utils.PanicIfNot(IsValidTexture(texture), "invalid texture")
	var location uint32 = g_effects[id.effect].samplers[slot].location
	var textureName uint32 = 0
	var validTexture bool = IsValidTexture(texture)
	if validTexture {
		textureName = TextureName(texture)
	}
	assignTexture(g_effects[id.effect].samplers[slot].Type, slot, location, textureName)
	if validTexture {
		TextureSamplerState(texture, sampler)
	}
}

// EffectAssignFloat ...
func EffectAssignFloat(id EffectId, slot int32, value float32) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	utils.PanicIfNot(EffectIsValidUniformSlot(id, slot), "invalid uniform slot")
	assignFloat(g_effects[id.effect].uniforms[slot].location, value)
}

// EffectAssignV4 ...
func EffectAssignV4(id EffectId, slot int32, value math.V4) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	utils.PanicIfNot(EffectIsValidUniformSlot(id, slot), "invalid uniform slot")
	assignVector4(g_effects[id.effect].uniforms[slot].location, value)
}

// EffectAssignM44 ...
func EffectAssignM44(id EffectId, slot int32, value math.M44, transpose bool) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	utils.PanicIfNot(EffectIsValidUniformSlot(id, slot), "invalid uniform slot")
	assignMatrix4(g_effects[id.effect].uniforms[slot].location, value, transpose)
}

// EffectAssignM44V ...
func EffectAssignM44V(id EffectId, slot int32, value []math.M44, transpose bool) {
	utils.PanicIfNot(EffectIsValid(id), "invalid id")
	utils.PanicIfNot(EffectIsValidUniformSlot(id, slot), "invalid uniform slot")
	assignMatrix4V(g_effects[id.effect].uniforms[slot].location, value, transpose)
}

// TemplateId ...
type TemplateId struct {
	template int32
}

// Template ...
type Template struct {
	name string

	variantAttributes []Location
	variant           int64

	vertexVariant int64
	pixelVariant  int64

	vertexKey ShaderKey
	pixelKey  ShaderKey

	attributes []Location
	samplers   []Location
	uniforms   []Location
	effect     EffectId
}

func appendLocations(dst []Location, src []Location) (out []Location) {
	var count int = len(src)
	dst = dst[:0]
	for i := 0; i < count; i++ {
		dst = append(dst, src[i])
	}

	//dst = resize(dst, count)
	//var i int32 = copy(dst, src) // ISSUE : memory corruption

	out = dst
	return
}

// TemplateInvalid ...
func TemplateInvalid() (out TemplateId) {
	out.template = -1
	return
}

// TemplateIsValid ...
func TemplateIsValid(id TemplateId) (out bool) {
	out = id.template >= 0 && int(id.template) < len(g_templates) // && is_valid_program(id.program) ##2 program id ??
	return
}

// TemplateCreate ...
func TemplateCreate(name string, vertexFilename string, pixelFilename string) (out TemplateId) {
	var vertexVariant int64 = ShaderAddVariant(vertexFilename)
	var pixelVariant int64 = ShaderAddVariant(pixelFilename)
	var templateVariant int64 = (vertexVariant << SHADER_VARIANT) | pixelVariant

	/*var templateCount int32 = len(g_templates)
	for i := 0; i < templateCount; i++ { // TODO : use map
		if g_templates[i].variant == templateVariant {
			out.template = i
			return
		}
	}*/

	out.template = int32(len(g_templates))

	var template Template
	template.name = name
	template.variant = templateVariant
	template.vertexVariant = vertexVariant
	template.pixelVariant = pixelVariant
	g_templates = append(g_templates, template)

	utils.PanicIfNot(TemplateIsValid(out), "invalid id")
	return
}

// TemplateReset ...
func TemplateReset(id TemplateId) {
	TemplateClearKey(id)
	g_templates[id.template].attributes = g_templates[id.template].attributes[:0]
	g_templates[id.template].samplers = g_templates[id.template].samplers[:0]
	g_templates[id.template].uniforms = g_templates[id.template].uniforms[:0]
}

// TemplateClearKey ...
func TemplateClearKey(id TemplateId) {
	TemplateClearVertexKey(id)
	TemplateClearPixelKey(id)
}

// TemplateSetKey ...
func TemplateSetKey(id TemplateId, flag int64, value bool) {
	TemplateSetVertexKey(id, flag, value)
	TemplateSetPixelKey(id, flag, value)
}

// TemplateClearVertexKey ...
func TemplateClearVertexKey(id TemplateId) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].vertexKey = ShaderKeyClear()
}

// TemplateSetVertexKey ...
func TemplateSetVertexKey(id TemplateId, flag int64, value bool) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].vertexKey = ShaderKeySet(g_templates[id.template].vertexKey, flag, value)
}

// TemplateClearPixelKey ...
func TemplateClearPixelKey(id TemplateId) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].pixelKey = ShaderKeyClear()
}

// TemplateSetPixelKey ...
func TemplateSetPixelKey(id TemplateId, flag int64, value bool) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].pixelKey = ShaderKeySet(g_templates[id.template].pixelKey, flag, value)
}

// TemplateBindAttribute ...
func TemplateBindAttribute(id TemplateId, location uint32, name string) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].attributes = LocationBind(g_templates[id.template].attributes, location, location, name, math.MAX_ui32)
}

// TemplateBindSampler ...
func TemplateBindSampler(id TemplateId, slot uint32, name string, Type uint32) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].samplers = LocationBind(g_templates[id.template].samplers, math.MAX_ui32, slot, name, Type)
}

// TemplateBindUniform ...
func TemplateBindUniform(id TemplateId, slot uint32, name string) {
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	g_templates[id.template].uniforms = LocationBind(g_templates[id.template].uniforms, math.MAX_ui32, slot, name, math.MAX_ui32)
}

func templateInstanceFromKey(id TemplateId, effectKey int64, vertexVariant int64, vertexFlags int64, vertexKey int64, pixelVariant int64, pixelFlags int64, pixelKey int64) (out EffectId) {
	out = EffectInvalid()
	utils.PanicIfNot(TemplateIsValid(id), "invalid id")
	// TODO UNUSED var templateVariant int64 = g_templates[id.template].variant

	var vertexShader ShaderId = ShaderCreate(SHADER_VERTEX, vertexVariant, vertexFlags, vertexKey)
	var pixelShader ShaderId = ShaderCreate(SHADER_PIXEL, pixelVariant, pixelFlags, pixelKey)

	var effectId EffectId = g_templates[id.template].effect
	var validEffect bool = false
	if EffectIsValid(effectId) == true { // ISSUE : no short-circuit evaluation
		if EffectGetKey(effectId) == effectKey {
			validEffect = true
		}
	}

	if validEffect == false {
		var effectCount int = len(g_effects)
		for i := 0; i < effectCount; i++ {
			if g_effects[i].key == effectKey {
				effectId.effect = int32(i)
				i = effectCount
				validEffect = true
			}
		}
	}

	if validEffect == false {
		effectId = EffectCreate(effectKey,
			g_templates[id.template].attributes,
			g_templates[id.template].samplers,
			g_templates[id.template].uniforms,
			vertexShader,
			pixelShader)
	}

	g_templates[id.template].effect = effectId
	EffectUse(effectId)
	out = effectId
	return
}

// TemplateInstanceFromKey ...
func TemplateInstanceFromKey(id TemplateId, effectKey int64) (out EffectId) {
	var vertexVariant int64 = g_templates[id.template].vertexVariant
	var vertexKey int64 = (effectKey >> 32) & ((1 << 32) - 1)
	var vertexFlags int64 = (vertexKey >> SHADER_VARIANT) & ((1 << (32 - SHADER_VARIANT)) - 1)
	var pixelVariant int64 = g_templates[id.template].pixelVariant
	var pixelKey int64 = effectKey & ((1 << 32) - 1)
	var pixelFlags int64 = (pixelKey >> SHADER_VARIANT) & ((1 << (32 - SHADER_VARIANT)) - 1)
	out = templateInstanceFromKey(id, effectKey, vertexVariant, vertexFlags, vertexKey, pixelVariant, pixelFlags, pixelKey)
	return
}

// TemplateInstance ...
func TemplateInstance(id TemplateId) (out EffectId) {
	var vertexVariant int64 = g_templates[id.template].vertexVariant
	var vertexFlags int64 = g_templates[id.template].vertexKey.value
	var vertexKey int64 = (vertexFlags << SHADER_VARIANT) | vertexVariant
	var pixelVariant int64 = g_templates[id.template].pixelVariant
	var pixelFlags int64 = g_templates[id.template].pixelKey.value
	var pixelKey int64 = (pixelFlags << SHADER_VARIANT) | pixelVariant
	var effectKey int64 = vertexKey<<32 | pixelKey
	out = templateInstanceFromKey(id, effectKey, vertexVariant, vertexFlags, vertexKey, pixelVariant, pixelFlags, pixelKey)
	return
}
