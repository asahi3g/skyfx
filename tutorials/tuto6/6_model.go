package tuto6

import (
	"fmt"
	"skyfx/app"
	"skyfx/fps"
	"skyfx/gfx"
	"skyfx/math"
	"skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"

	"golang.org/x/mobile/gl"
)

// Globals ...
var ry float32 = 0.0

var modelPaths []string

var models []gfx.ModelId
var scales []float32
var rotations []float32
var speeds []float32

var currentModel int32

var skyMesh gfx.MeshId = gfx.InvalidMesh()
var groundMesh gfx.MeshId = gfx.InvalidMesh()
var texture gfx.TextureId = gfx.InvalidTexture()
var skyDiffuses []gfx.TextureId
var skySpeculars []gfx.TextureId
var currentSky int32
var currentSkybox int32
var brdf gfx.TextureId = gfx.InvalidTexture()

var world math.M44

var camera app.CameraId

var worldPosY float32 = -10.0

var worldSizeX float32 = 2000.0
var worldSizeZ float32 = 2000.0

var worldCellX int32 = 64
var worldCellZ int32 = 64

var groundScale float32 = 20.0
var groundRight float32 = worldSizeX / groundScale
var groundBack float32 = worldSizeZ / groundScale
var groundWidth uint32 = 16
var groundHeight uint32 = 16

var CAMERA_TPS int32 = 0
var CAMERA_FREE int32 = 1
var CAMERA_FPS int32 = 2
var CAMERA_COUNT int32 = 3
var currentCamera int32 = CAMERA_TPS

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
	} else if mods == app.MOD_ALT && action == app.KEY_PRESS {
		if key == app.KEYCODE_LEFT {
			if currentSky <= 0 {
				currentSky = int32(len(skyDiffuses))
			}
			currentSky = currentSky - 1
		} else if key == app.KEYCODE_RIGHT {
			currentSky = int32((currentSky + 1) % int32(len(skyDiffuses)))
		} else if key == app.KEYCODE_UP {
			if currentSkybox <= 0 {
				currentSkybox = 1
			}
			currentSkybox = currentSkybox - 1
		} else if key == app.KEYCODE_DOWN {
			currentSkybox = (currentSkybox + 1) % 2
		}
	} else if mods == app.MOD_CTRL && action == app.KEY_PRESS {
		if key == app.KEYCODE_LEFT {
			if currentModel <= 0 {
				currentModel = int32(len(models))
			}
			currentModel = currentModel - 1
			ry = 0.0
		} else if key == app.KEYCODE_RIGHT {
			currentModel = int32((currentModel + 1) % int32(len(models)))
			ry = 0.0
		} else if key == app.KEYCODE_UP {
			rotations[currentModel] = rotations[currentModel] + 1.0
		} else if key == app.KEYCODE_DOWN {
			rotations[currentModel] = rotations[currentModel] - 1.0
		}
	} else if mods == app.MOD_SHIFT && action == app.KEY_PRESS {
		if key == app.KEYCODE_LEFT {
			currentCamera = currentCamera - 1
			if currentCamera < 0 {
				currentCamera = CAMERA_COUNT - 1
			}
		} else if key == app.KEYCODE_RIGHT {
			currentCamera = currentCamera + 1
			if currentCamera >= CAMERA_COUNT {
				currentCamera = 0
			}
		} else if key == app.KEYCODE_UP {
		} else if key == app.KEYCODE_DOWN {
		}
	} /* else if mods == app.MOD_SHIFT && action == app.KEY_PRESS {
		if key == app.KEYCODE_UP {
			scales[currentModel] = scales[currentModel] + 0.05
		} else if key == app.KEYCODE_DOWN {
			scales[currentModel] = scales[currentModel] - 0.05
		}
	}*/

	if currentCamera == CAMERA_FREE {
		app.FreeCameraProcessKeyboard(camera, event)
		// } else if currentCamera == CAMERA_FPS {
		// 	app.FpsCameraProcessKeyboard(camera, event)
		// } else if currentCamera == CAMERA_TPS {
		// 	app.TpsCameraProcessKeyboard(camera, event)
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnMouse(event *app.MouseEvent) {
	if currentCamera == CAMERA_FREE {
		app.FreeCameraProcessMouse(camera, event)
	} else if currentCamera == CAMERA_FPS {
		app.FpsCameraProcessMouse(camera, event)
	} else if currentCamera == CAMERA_TPS {
		app.TpsCameraProcessMouse(camera, event)
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnUpdate() {
	var dts float64 = fps.DeltaSecond()
	var dt float32 = float32(dts)

	var s float32 = scales[currentModel]

	ry = ry + dt*rotations[currentModel]
	world = m44.Makef_SAT(s, s, s, 0.0, 1.0, 0.0, ry, 0.0, 0.0, 0.0)

	var model gfx.ModelId = models[currentModel]
	var anim gfx.AnimationId = gfx.ModelGetAnimation(model, 0)
	if currentCamera == CAMERA_FREE {
		app.FreeCameraUpdate(camera, dts, 1.0, 1.0)
	} else if currentCamera == CAMERA_FPS {
		app.FpsCameraUpdate(camera, dts, v3.ZERO)
	} else if currentCamera == CAMERA_TPS {
		app.TpsCameraUpdate(camera, dts, v3.ZERO, 30.0)
	}
	var step float32 = dt * speeds[currentModel]
	if gfx.AnimationIsValid(anim) {
		gfx.AnimationUpdate(anim, step, true, -1.0)
	}
	gfx.ModelUpdate(model, anim, step, true)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnRender() {
	var view math.M44 = app.CameraGetView(camera)
	var projection math.M44 = app.CameraGetProjection(camera)
	var cameraPosition math.V3 = app.CameraGetPosition(camera)

	gfx.DisableBlending()
	gfx.DepthState(true, gl.LESS, true)
	gfx.Clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, v4.Make(0.3, 0.5, 0.6, 1.0), 1.0, 0)

	//sky
	gfx.EffectUse(gfx.FxSky)
	var skybox gfx.TextureId
	if currentSkybox == 0 {
		skybox = skySpeculars[currentSky]
	} else {
		skybox = skyDiffuses[currentSky]
	}
	gfx.EffectAssignTexture(gfx.FxSky, gfx.SAMPLER_ENV_DIFFUSE, skybox, gfx.SpLinear0Clamp)
	gfx.EffectAssignM44(gfx.FxSky, gfx.UNIFORM_WORLD, m44.IDENTITY, false)
	gfx.EffectAssignM44(gfx.FxSky, gfx.UNIFORM_VIEW, view, false)
	gfx.EffectAssignM44(gfx.FxSky, gfx.UNIFORM_PROJECTION, projection, false)
	gfx.MeshRender(skyMesh)

	// ground
	gfx.EffectUse(gfx.FxTexture3D)
	gfx.EffectAssignTexture(gfx.FxTexture3D, gfx.SAMPLER_COLOR_0, texture, gfx.SpLinearWrap)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_WORLD, m44.IDENTITY, false)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_VIEW, view, false)
	gfx.EffectAssignM44(gfx.FxTexture3D, gfx.UNIFORM_PROJECTION, projection, false)
	gfx.MeshRender(groundMesh)

	// model
	var model gfx.ModelId = models[currentModel]
	gfx.ModelSetWorld(model, world)
	gfx.ModelSetView(model, view)
	gfx.ModelSetProjection(model, projection)
	gfx.ModelSetEnvironmentSpecular(model, skySpeculars[currentSky])
	gfx.ModelSetEnvironmentDiffuse(model, skyDiffuses[currentSky])
	gfx.ModelSetBRDF(model, brdf)
	gfx.ModelSetCameraPosition(model, cameraPosition)
	gfx.ModelSetExposure(model, 1.0)
	gfx.ModelRender(model)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func loadModel(path string, filename string, s float32, r float32, v float32) {
	var model gfx.ModelId = gfx.ModelCreateFromFile(fmt.Sprintf("%smodels/%s/", app.DataDir(), path), filename, gfx.MODEL_GEOMETRY|gfx.MODEL_ANIMATION)
	models = append(models, model)

	scales = append(scales, s)
	rotations = append(rotations, r)
	speeds = append(speeds, v)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
	//StartCPUProfile("6_model_init", 100)

	skyMesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, 36, gfx.VertexLayout, 24)
	gfx.MeshBegin(skyMesh)
	gfx.MeshAppendBox(skyMesh, false, true,
		v3.ZERO,
		v3.Make(5000.0, 0.0, 0.0), v3.Make(0.0, 5000.0, 0.0), v3.Make(0.0, 0.0, 5000.0),
		v4.ONE)
	gfx.MeshEnd(skyMesh)

	// loadModel("glTF-Sample-Models/Box/glTF", "Box.gltf", 1.0, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/BoxAnimated/glTF", "BoxAnimated.gltf", 1.0, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/NormalTangentTest/glTF", "NormalTangentTest.gltf", 3.0, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/NormalTangentMirrorTest/glTF", "NormalTangentMirrorTest.gltf", 3.0, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/MetalRoughSpheres/glTF", "MetalRoughSpheres.gltf", 1.0, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/CesiumMilkTruck/glTF", "CesiumMilkTruck.gltf", 0.5, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/2CylinderEngine/glTF", "2CylinderEngine.gltf", 0.005, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/ReciprocatingSaw/glTF", "ReciprocatingSaw.gltf", 0.01, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/GearboxAssy/glTF", "GearboxAssy.gltf", 1.0, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/Buggy/glTF", "Buggy.gltf", 0.05, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/RiggedSimple/glTF", "RiggedSimple.gltf", 3.0, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/RiggedFigure/glTF", "RiggedFigure.gltf", 3.0, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/CesiumMan/glTF", "CesiumMan.gltf", 3.0, 0.0, 2.0)
	// loadModel("glTF-Sample-Models/Monster/glTF", "Monster.gltf", 0.05, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/BrainStem/glTF", "BrainStem.gltf", 1.0, 0.0, 1.0)
	// loadModel("glTF-Sample-Models/Duck/glTF", "Duck.gltf", 1.0, 1.0, 1.0)
	loadModel("glTF-Sample-Models/DamagedHelmet/glTF", "DamagedHelmet.gltf", 3.0, 1.0, 1.0)
	// loadModel("glTF-Sample-Models/BoomBox/glTF", "BoomBox.gltf", 200.0, 1.0, 1.0)
	loadModel("glTF-Sample-Models/VC/glTF", "VC.gltf", 1.0, 0.0, 1.0)
	loadModel("glTF-Sample-Models/Sponza/glTF", "Sponza.gltf", 10.0, 0.0, 1.0)
	// loadModel("skylight/skycoin/low", "scene.gltf", 1000.0, 0.0, 1.0)
	// loadModel("skylight/skywatch/medium", "scene.gltf", 1.0, 0.0, 1.0)
	// loadModel("skylight/orangepi/high", "scene.gltf", 10.0, 0.0, 1.0)
	loadModel("skylight/skyminer/high", "scene.gltf", 1.0, 0.0, 1.0)

	groundMesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, int32(6*groundWidth*groundHeight), gfx.VertexLayout, 6*groundWidth*groundHeight)
	// var groundScale float32 = 20.0
	gfx.MeshBegin(groundMesh)
	gfx.MeshAppendPlane(groundMesh, false, true,
		v3.Make(0.0, worldPosY, 0.0),
		v3.Make(groundRight, 0.0, 0.0),
		v3.GREEN,
		v3.Make(0.0, 0.0, groundBack),
		groundWidth, groundHeight,
		v4.Make(0.0, 0.0, 10.0, 10.0),
		v4.ONE)
	gfx.MeshEnd(groundMesh)

	texture = gfx.TextureCreate(fmt.Sprintf("%s/textures/Skycoin-Cloud-BW-Vertical-on_black@2x.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, -32, false, false)

	var skies []string
	skies = append(skies, "doge2")
	// skies = append(skies, "ennis")
	// skies = append(skies, "field")
	// skies = append(skies, "footprint_court")
	// skies = append(skies, "helipad")
	// skies = append(skies, "papermill")
	// skies = append(skies, "pisa")
	// skies = append(skies, "studio_grey")
	// skies = append(skies, "studio_red_green")
	for i := 0; i < len(skies); i++ {
		skySpeculars = append(skySpeculars, gfx.TextureCreateCube(fmt.Sprintf("%s/textures/environments/%s/specular/specular_.hdr", app.DataDir(), skies[i]), gfx.FORMAT_RGB_16F, 0, 0, 32, false))
		skyDiffuses = append(skyDiffuses, gfx.TextureCreateCube(fmt.Sprintf("%s/textures/environments/%s/diffuse/diffuse_.hdr", app.DataDir(), skies[i]), gfx.FORMAT_RGB_16F, 0, 0, 0, false))
	}

	brdf = gfx.TextureCreate(fmt.Sprintf("%s/textures/environments/brdf.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, 0, false, false)

	camera = app.CameraCreate()
	app.CameraSetProjection(camera, 1.0, 10000.0, 0.5, gfx.ViewportSize.X, gfx.ViewportSize.Y, true)
	app.CameraSetPosition(camera, v3.ZERO, true)
	app.CameraSetYawPitch(camera, 0.0, 0.0, true)

	//StopCPUProfile("6_model_init")
	//StartCPUProfile("6_model_run", 100)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func Run() {

	app.SetOnStart(appOnStart)
	app.SetOnKeyboard(appOnKeyboard)
	app.SetOnMouse(appOnMouse)
	app.SetOnUpdate(appOnUpdate)
	app.SetOnRender(appOnRender)

	if app.Run("skyfx: Model Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(skyMesh)
		gfx.MeshUnlock(groundMesh)
		app.Destroy()
	}
}
