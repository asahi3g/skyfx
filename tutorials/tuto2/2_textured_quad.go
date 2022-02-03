package tuto2

import (
	"fmt"
	"skyfx/app"
	"skyfx/gfx"
	v4 "skyfx/math/v4"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Globals ...
var mesh gfx.MeshId
var texture gfx.TextureId

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnKeyboard(event *app.KeyboardEvent) {
	var key int32 = event.Key
	var action app.ActionType = event.Action
	var mods int32 = event.Mods
	if mods == app.MOD_NONE {
		if key == app.KEYCODE_ESCAPE && action == app.KEY_PRESS {
			app.Exit()
		}
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnUpdate() {
	var w float32 = gfx.TextureWidthF32(texture) / 10.0
	var h float32 = gfx.TextureHeightF32(texture) / 10.0
	var x float32 = (gfx.ViewportSize.X - w) / 2.0
	var y float32 = (gfx.ViewportSize.Y - h) / 2.0

	gfx.MeshBegin(mesh)
	gfx.MeshAppendOrthoQuad(mesh,
		v4.Make(x, y, w, h),
		v4.BLUE,
		v4.ONE,
		v4.Make(0.0, 0.0, gfx.ViewportSize.X, gfx.ViewportSize.Y),
		0.0)
	gfx.MeshEnd(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnRender() {
	gfx.EffectUse(gfx.FxTexture2D)
	gfx.EffectAssignTexture(gfx.FxTexture2D, gfx.SAMPLER_COLOR_0, texture, gfx.SpLinearWrap)
	gfx.MeshRender(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
	mesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, 6, gfx.VertexLayout, 4)
	texture = gfx.TextureCreate(fmt.Sprintf("%s/textures/Skycoin-Cloud-BW-Vertical-on_black@2x.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {
	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Textured Quad Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(mesh)
		app.Destroy()
	}
}
