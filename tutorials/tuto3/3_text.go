package tuto3

import (
	"fmt"
	"skyfx/app"
	"skyfx/gfx"
	v2 "skyfx/math/v2"
	v4 "skyfx/math/v4"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/skycoin/gltext"
)

// Globals ...
var mesh gfx.MeshId
var texture gfx.TextureId
var font string = "skycoinRegular"

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
	var textWidth int32
	var textHeight int32
	textWidth, textHeight = gfx.MeasureText(font, app.Name())

	var w float32 = float32(textWidth)
	var h float32 = float32(textHeight)
	var x float32 = (gfx.ViewportSize.X - w) / 2.0
	var y float32 = (gfx.ViewportSize.Y - h) / 2.0

	gfx.MeshBegin(mesh)
	gfx.MeshAppendText(mesh, texture, font,
		v2.Make(x, y), v2.ONE,
		v4.RED,
		app.Name(),
		false, v4.ZERO, v4.ZERO, // TODO remove debug arguments
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
	var nameLen int32 = int32(len(app.Name()))
	mesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, nameLen*6, gfx.VertexLayout, uint32(nameLen*4))
	texture = gfx.TextureCreateFont(font, fmt.Sprintf("%s/fonts/Skycoin-Regular.ttf", app.DataDir()), 64, 32, 127, int32(gltext.LeftToRight), -1)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {

	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Text Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(mesh)
		app.Destroy()
	}
}
