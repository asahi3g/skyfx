package tuto5

import (
	"fmt"
	"skyfx/app"
	"skyfx/fps"
	"skyfx/gfx"
	"skyfx/math"

	"skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Globals ...
var mesh gfx.MeshId
var texture gfx.TextureId

var ry float32 = 0.0
var world math.M44

var camera app.CameraId

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

	app.FreeCameraProcessKeyboard(camera, event)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnMouse(event *app.MouseEvent) {
	app.FreeCameraProcessMouse(camera, event)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnUpdate() {
	app.FreeCameraUpdate(camera, fps.DeltaSecond(), 1.0, 1.0)

	ry = ry + float32(fps.DeltaSecond())
	ry = 0
	world = m44.Makef_AT(0.0, 1.0, 0.0, ry, 0.0, 0.0, -50.0)

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
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_VIEW, app.CameraGetView(camera), false)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_PROJECTION, app.CameraGetProjection(camera), false)
	gfx.MeshRender(mesh)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
	mesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, 36, gfx.VertexLayout, 36)
	texture = gfx.TextureCreate(fmt.Sprintf("%s/textures/Skycoin-Cloud-BW-Vertical-on_black@2x.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)

	camera = app.CameraCreate()
	app.CameraSetProjection(camera, 1.0, 10000.0, 0.5, gfx.ViewportSize.X, gfx.ViewportSize.Y, true)
	app.CameraSetPosition(camera, v3.ZERO, true)
	app.CameraSetYawPitch(camera, 0.0, 0.0, true)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {

	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnMouse(appOnMouse)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Camera Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(mesh)
		app.Destroy()
	}
}
