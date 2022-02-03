package gfx

import (
	"fmt"
	"skyfx/math"
	v4 "skyfx/math/v4"
	"skyfx/utils"

	"github.com/go-gl/gl/v2.1/gl"
)

// // Globals ...
var g_targets []Target

// // TargetId ...
type TargetId struct {
	target int32
}

// // Target ...
type Target struct {
	id          TargetId
	scale       uint32
	width       uint32
	height      uint32
	fwidth      float32
	fheight     float32
	bounds      math.V4
	size        math.V4
	framebuffer uint32
	colors      []TextureId
	depth       TextureId
	stencil     TextureId
	age         int32
	resizable   bool
}

// // InvalidTarget ...
func InvalidTarget() (out TargetId) {
	out.target = -1
	return
}

// // IsValidTarget ...
func IsValidTarget(id TargetId) (out bool) {
	out = id.target >= 0 && int(id.target) < len(g_targets)
	return
}

// TargetPrint ...
func TargetPrint(message string, id TargetId) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	fmt.Printf("%s : id %d, framebuffer %d\n", message, id.target, g_targets[id.target].framebuffer)
	var colors []TextureId = g_targets[id.target].colors
	var colorCount int = len(colors)
	for i := 0; i < colorCount; i++ {
		if IsValidTexture(colors[i]) {
			TexturePrint(message, colors[i])
		}
	}
	var depth TextureId = g_targets[id.target].depth
	if IsValidTexture(depth) {
		TexturePrint(message, depth)
	}
	var stencil TextureId = g_targets[id.target].stencil
	if IsValidTexture(stencil) {
		TexturePrint(message, stencil)
	}
}

// TargetCreate ...
func TargetCreate(colors int32, width uint32, height uint32, scale uint32, resizable bool) (out TargetId) {
	out.target = int32(len(g_targets))

	var target Target
	target.id = out
	var framebuffer uint32 = 0
	gl.GenFramebuffers(1, &framebuffer)
	utils.PanicIf(GlError(), "gl.GenFramebuffers")
	target.framebuffer = framebuffer
	for i := 0; i < int(colors); i++ {
		target.colors = append(target.colors, InvalidTexture())
	}
	target.depth = InvalidTexture()
	target.stencil = InvalidTexture()

	g_targets = append(g_targets, target)

	targetResize(out, width, height, scale, resizable)
	utils.PanicIfNot(IsValidTarget(out), "")
	return
}

func targetResize(id TargetId, width uint32, height uint32, scale uint32, resizable bool) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	var swidth uint32 = width / scale
	var sheight uint32 = height / scale
	var fwidth float32 = float32(swidth)
	var fheight float32 = float32(sheight)

	g_targets[id.target].scale = scale
	g_targets[id.target].width = swidth
	g_targets[id.target].height = sheight
	g_targets[id.target].fwidth = fwidth
	g_targets[id.target].fheight = fheight
	g_targets[id.target].bounds = v4.Make(0.0, 0.0, fwidth, fheight)
	g_targets[id.target].size = v4.Make(fwidth, fheight, 1.0/fwidth, 1.0/fheight)
	g_targets[id.target].resizable = resizable
}

// TargetDestroy ...
func TargetDestroy(target Target) {
}

// TargetBind ...
func TargetBind(id TargetId) {
	utils.PanicIfNot(IsValidTarget(id), "")
	//TargetPrint("TargetBind", id)
	bindFramebuffer(g_targets[id.target].framebuffer)

	var bounds math.V4 = g_targets[id.target].bounds
	SetViewport(bounds)
	SetScissor(bounds)

	var colors []TextureId = g_targets[id.target].colors
	DrawBuffers(int32(len(colors)))
}

// TargetAttachColor ...
func TargetAttachColor(id TargetId, texture TextureId, i uint32) {
	//fmt.Printf("TargetAttachColor\n")
	utils.PanicIfNot(IsValidTarget(id), "")
	targetAttachTexture(id, gl.COLOR_ATTACHMENT0+i, texture)
	var colors []TextureId = g_targets[id.target].colors
	colors[i] = texture
	//fmt.Printf("TARGET_ATTACH_COLOR : TARGET %d, COLOR %d, TEXTURE %d\n", id.target, i, texture.texture)
}

// TargetAttachDepth ...
func TargetAttachDepth(id TargetId, depth TextureId) {
	//fmt.Printf("TargetAttachDepth\n")
	utils.PanicIfNot(IsValidTarget(id), "")
	targetAttachTexture(id, gl.DEPTH_ATTACHMENT, depth)
	g_targets[id.target].depth = depth
}

// TargetAttachStencil ...
func TargetAttachStencil(id TargetId, stencil TextureId) {
	fmt.Printf("TargetAttachStencil\n")
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	targetAttachTexture(id, gl.STENCIL_ATTACHMENT, stencil)
	g_targets[id.target].stencil = stencil
}

// TargetGetAge ...
func TargetGetAge(id TargetId) (out int32) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	out = g_targets[id.target].age
	return
}

// TargetSetAge ...
func TargetSetAge(id TargetId, age int32) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	g_targets[id.target].age = age
	return
}

// TargetGetSize ...
func TargetGetSize(id TargetId) (out math.V4) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	out = g_targets[id.target].size
	return
}

// TargetGetWidth ...
func TargetGetWidth(id TargetId) (out float32) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	out = g_targets[id.target].fwidth
	return
}

// TargetGetHeight ...
func TargetGetHeight(id TargetId) (out float32) {
	utils.PanicIfNot(IsValidTarget(id), "invalid target")
	out = g_targets[id.target].fheight
	return
}

// TargetGetColor ...
func TargetGetColor(id TargetId, i int32) (out TextureId) {
	out = InvalidTexture()
	if IsValidTarget(id) {
		var colors []TextureId = g_targets[id.target].colors
		if i >= 0 && i < int32(len(colors)) {
			out = colors[i]
			//fmt.Printf("TARGET_GET_COLOR : TARGET %d, COLOR %d, TEXTURE %d\n", id.target, i, out.texture)
		}
	}
	return
}

// TargetGetDepth ...
func TargetGetDepth(id TargetId) (out TextureId) {
	if IsValidTarget(id) {
		out = g_targets[id.target].depth
	} else {
		out = InvalidTexture()
	}
	return
}

// DestroyTargets ...
func DestroyTargets() {
	var count int = len(g_targets)
	for i := 0; i < count; i++ {
		gl.DeleteFramebuffers(1, &g_targets[i].framebuffer)
		utils.PanicIf(GlError(), "g.DeleteFramebuffers")
	}
}

// ResizeTargets ...
func ResizeTargets(width uint32, height uint32) {
	var count int = len(g_targets)
	for i := 0; i < count; i++ {
		var id TargetId = g_targets[i].id
		var resizable bool = g_targets[i].resizable
		if resizable {
			var colors []TextureId = g_targets[i].colors
			var colorCount int = len(colors)
			var targetScale uint32 = g_targets[i].scale
			var targetWidth uint32 = width / targetScale
			var targetHeight uint32 = height / targetScale
			if targetWidth != g_targets[i].width || targetHeight != g_targets[i].height {
				targetResize(id, width, height, targetScale, resizable)
				for c := 0; c < colorCount; c++ {
					if IsValidTexture(colors[c]) {
						TextureResize(colors[c], targetWidth, targetHeight)
						TextureSamplerState(colors[c], g_nearest0Clamp)
						fmt.Printf("ResizeTargets\n")
						TargetAttachColor(id, colors[c], uint32(c))
					}
				}

				var depth TextureId = g_targets[i].depth
				if IsValidTexture(depth) {
					TextureResize(depth, targetWidth, targetHeight)
					TextureSamplerState(depth, g_nearest0Clamp)
					TargetAttachDepth(id, depth)
				}

				var stencil TextureId = g_targets[i].stencil
				if IsValidTexture(stencil) && depth.texture != stencil.texture {
					TextureResize(stencil, targetWidth, targetHeight)
					TextureSamplerState(stencil, g_nearest0Clamp)
					TargetAttachStencil(id, stencil)
				}
			}
		}
	}
}

func targetAttachTexture(id TargetId, attachment uint32, texture TextureId) {
	utils.PanicIfNot(IsValidTarget(id), "")
	//TexturePrint("targetAttachTexture", texture)
	bindFramebuffer(g_targets[id.target].framebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, attachment, gl.TEXTURE_2D, TextureName(texture), 0)
	utils.PanicIf(GlError(), "gl.FramebufferTexture2D")
	utils.PanicIf(GlStatus(), "gl.CheckFramebufferStatus")
}

// TODO renderbuffers
// TargetAttachBuffer ...
//	renderbuffer = gl.GenRenderbuffers(1, renderbuffer)
//	gl.BindRenderbuffer(gl.RENDERBUFFER, renderbuffer)
//	gl.RenderbufferStorage(gl.RENDERBUFFER, format, w, h)
//	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthbuffer)
//	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.STENCIL_ATTACHMENT, gl.RENDERBUFFER, stencilbuffer)
