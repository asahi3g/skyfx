package tuto4

import (
	"fmt"
	"skyfx/app"
	"skyfx/fps"
	"skyfx/gfx"
	"skyfx/math"

	m44 "skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Globals ...
var mesh gfx.MeshId
var texture gfx.TextureId

var dx float32 = 0.0
var dy float32 = 0.0
var dz float32 = -50.0
var ry float32 = 0.0
var world math.M44

var projection math.M44

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

	if key == app.KEYCODE_LEFT {
		dx = dx + 1.0
	}
	if key == app.KEYCODE_RIGHT {
		dx = dx - 1.0
	}
	if key == app.KEYCODE_DOWN {
		if mods > 0 {
			dy = dy + 1.0
		} else {
			dz = dz - 1.0
		}
	}
	if key == app.KEYCODE_UP {
		if mods > 0 {
			dy = dy - 1.0
		} else {
			dz = dz + 1.0
		}
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnUpdate() {
	projection = m44.Make_project(0.1, 10000.0, 0.5, gfx.ViewportSize.X, gfx.ViewportSize.Y)

	ry = ry + float32(fps.DeltaSecond())
	world = m44.Makef_AT(0.0, 1.0, 0.0, ry, dx, dy, dz)

	gfx.MeshBegin(mesh)
	gfx.MeshAppendBox(mesh, false, false,
		v3.ZERO,
		v3.Make(10.0, 0.0, 0.0),
		v3.Make(0.0, 10.0, 0.0),
		v3.Make(0.0, 0.0, 10.0),
		v4.ONE)
	gfx.MeshEnd(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnRender() {
	gfx.DisableBlending()
	gfx.DepthState(true, gl.LESS, true)
	gfx.Clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, v4.Make(0.3, 0.5, 0.6, 1.0), 1.0, 0)
	gfx.EffectUse(gfx.FxTexture3D)
	gfx.EffectAssignTexture(gfx.FxTexture3D, gfx.SAMPLER_COLOR_0, texture, gfx.SpLinearWrap)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_WORLD, world, false)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_VIEW, m44.IDENTITY, false)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_PROJECTION, projection, false)
	gfx.MeshRender(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
	mesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, 36, gfx.VertexLayout, 36)
	texture = gfx.TextureCreate(fmt.Sprintf("%s/textures/Skycoin-Cloud-BW-Vertical-on_black@2x.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {

	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Perspective Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(mesh)
		app.Destroy()
	}
}
