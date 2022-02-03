package tuto1

import (
	"skyfx/app"
	"skyfx/gfx"
	v4 "skyfx/math/v4"
)

// Globals ...
var mesh gfx.MeshId

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
	var w float32 = 200.0
	var h float32 = 200.0
	var x float32 = (gfx.ViewportSize.X - w) / 2.0
	var y float32 = (gfx.ViewportSize.Y - h) / 2.0

	gfx.MeshBegin(mesh)
	gfx.MeshAppendOrthoQuad(mesh,
		v4.Make(x, y, w, h),
		v4.BLUE,
		v4.GREEN,
		v4.Make(0.0, 0.0, gfx.ViewportSize.X, gfx.ViewportSize.Y),
		0.0)
	gfx.MeshEnd(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnRender() {
	gfx.EffectUse(gfx.FxVertexColor2D)
	gfx.MeshRender(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {

	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Colored Quad Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(mesh)
		app.Destroy()
	}
}
