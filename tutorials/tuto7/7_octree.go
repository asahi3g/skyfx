package tuto7

import (
	"fmt"
	"skyfx/app"
	"skyfx/fps"
	"skyfx/gfx"
	"skyfx/math"
	"skyfx/math/m44"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"

	"time"

	"golang.org/x/mobile/gl"
)

// TODO : can't navigate with keyboard at startup
// TODO : dialog disapear if exit button is pushed with keyboard

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

var octree gfx.OctreeId = gfx.InvalidOctree()
var frustum gfx.FrustumId = gfx.InvalidFrustum()

var world math.M44

var camera app.CameraId
var snapFrustum bool = true

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
		} else if key == app.KEYCODE_1 && action == app.KEY_PRESS {
			if snapFrustum {
				snapFrustum = false
			} else {
				snapFrustum = true
			}
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
			currentModel = (currentModel + 1) % int32(len(models))
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
	if snapFrustum {
		gfx.FrustumUpdate(frustum, app.CameraGetInvViewProj(camera))
	}
	var dts float64 = fps.DeltaSecond()
	var dt float32 = float32(dts)

	var s float32 = scales[currentModel]

	ry = ry + dt*rotations[currentModel]
	world = m44.Makef_SAT(s, s, s, 0.0, 1.0, 0.0, ry, 0.0, 0.0, 0.0)

	if currentCamera == CAMERA_FREE {
		app.FreeCameraUpdate(camera, dts, 1.0, 1.0)
	} else if currentCamera == CAMERA_FPS {
		app.FpsCameraUpdate(camera, dts, v3.ZERO)
	} else if currentCamera == CAMERA_TPS {
		app.TpsCameraUpdate(camera, dts, v3.ZERO, 30.0)
	}
	gfx.OctreeUpdate(octree, frustum, 4)
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

	// octree
	gfx.OctreeRender(octree, m44.IDENTITY, view, projection,
		skySpeculars[currentSky], skyDiffuses[currentSky], brdf, v4.Make_v31(cameraPosition, 1.0), 1.0)
	gfx.FrustumRender(frustum, m44.IDENTITY, view, projection)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func loadModel(path string, filename string, s float32, r float32, v float32) (out int32) {
	var modelPath string = fmt.Sprintf("%smodels/%s/", app.DataDir(), path)

	var model gfx.ModelId = gfx.ModelInvalid()
	var modelCount int = len(models)
	for modelIndex := 0; modelIndex < modelCount; modelIndex++ {
		var m gfx.ModelId = models[modelIndex]
		if (gfx.ModelGetPath(m) == modelPath) && (gfx.ModelGetName(m) == filename) {
			model = m
			out = int32(modelIndex)
			modelIndex = modelCount
		}
	}

	if gfx.ModelIsValid(model) == false {
		out = int32(len(models))
		model = gfx.ModelCreateFromFile(modelPath, filename, gfx.MODEL_GEOMETRY|gfx.MODEL_ANIMATION)
		models = append(models, model)
		scales = append(scales, s)
		rotations = append(rotations, r)
		speeds = append(speeds, v)
	}
	return
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func loadOctree(path string, name string, maxLevel int32, options int32, save bool, world math.M44, lods bool, reload bool) (out gfx.OctreeId) {
	var t0 int64

	var minLevel int32 = 1
	if lods == false {
		minLevel = maxLevel
	}

	for level := minLevel; level <= maxLevel; level++ {
		var octreePath string = fmt.Sprintf("%s/models/%s/octree_lod_%d.oct", app.DataDir(), path, level)
		if reload == false || save == true {
			t0 = time.Now().UnixNano()
			var modelIndex int32 = loadModel(path, name, 1.0, 0.0, 1.0)
			var octreeModel gfx.ModelId = models[modelIndex]
			gfx.ModelUpdate(octreeModel, gfx.AnimationInvalid(), 0.0, true)
			fmt.Printf("model loaded in %f seconds\n", fps.NanoToSecond(time.Now().UnixNano()-t0))

			t0 = time.Now().UnixNano()
			var left gfx.OctreeId = gfx.OctreeCreate(octreeModel, world, level, true, options)
			fmt.Printf("Octree model split in %f seconds\n", fps.NanoToSecond(time.Now().UnixNano()-t0))
			out = left

			if save {
				t0 = time.Now().UnixNano()
				_ /*var octreeSaved bool*/ = gfx.OctreeSave(left, octreePath, options)
				fmt.Printf("Octree model saved in %f seconds\n", fps.NanoToSecond(time.Now().UnixNano()-t0))
			}
		}

		t0 = time.Now().UnixNano()
		if reload {
			var right gfx.OctreeId = gfx.OctreeLoad(octreePath, true)
			fmt.Printf("Octree model loded in %f seconds\n", fps.NanoToSecond(time.Now().UnixNano()-t0))
			//gfx.OctreeAssertEquals(left, right, options)
			out = right
		}
	}
	return
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
func appOnStart() {
	skyMesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, 36, gfx.VertexLayout, 24)
	gfx.MeshBegin(skyMesh)
	gfx.MeshAppendBox(skyMesh, false, true,
		v3.ZERO,
		v3.Make(5000.0, 0.0, 0.0), v3.Make(0.0, 5000.0, 0.0), v3.Make(0.0, 0.0, 5000.0),
		v4.ONE)
	gfx.MeshEnd(skyMesh)

	var octreeWorld math.M44 = m44.Makef_SAT(1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0)
	var gfxOctree gfx.OctreeId = loadOctree("skylight/skyminer/high", "scene.gltf", 4, gfx.OCTREE_GRAPHICS, false, octreeWorld, false, false)

	//var octreeWorld math.M44 = m44.makev_SAT(v3.Makef(100.0), v4.Make(0.0, 0.0, 1.0, 3.0 * math.PI_f32/2.0), v3.Make(0.0, 1500.0, 600.0))
	//var gfxOctree gfx.OctreeId = loadOctree("skylight/skyminer/high", "scene.gltf", 2, gfx.OCTREE_GRAPHICS, false, octreeWorld, false, true)
	//var phxOctree gfx.OctreeId = loadOctree("skylight/skyminer/collision", "scene.gltf", 4, gfx.OCTREE_COLLISIONS, false, octreeWorld, false, true)

	//var octreeWorld math.M44 = m44.makev_SAT(v3.Makef(100.0), v4.Make(0.0, 0.0, 1.0, 0.0), v3.Make(0.0, 1500.0, 600.0))
	//var gfxOctree gfx.OctreeId = loadOctree("skylight/skyantenna/high", "scene.gltf", 2, gfx.OCTREE_GRAPHICS, false, octreeWorld, false, true)
	//var phxOctree gfx.OctreeId = loadOctree("skylight/skyantenna/collision", "scene.gltf", 4, gfx.OCTREE_COLLISIONS, false, octreeWorld, false, true)

	//var octreeWorld math.M44 = m44.makev_SAT(v3.Makef(20.0), v4.Make(0.0, 0.0, 1.0, 0.0), v3.Make(0.0, 100.0, 0.0))
	//var phxOctree gfx.OctreeId = loadOctree("glTF-Sample-Models/Sponza/glTF", "Sponza.gltf", 4, gfx.OCTREE_COLLISIONS, false, octreeWorld, false, true)
	//var gfxOctree gfx.OctreeId = loadOctree("glTF-Sample-Models/Sponza/glTF", "Sponza.gltf", 1, gfx.OCTREE_GRAPHICS, false, octreeWorld, false, true)

	octree = gfxOctree
	frustum = gfx.FrustumCreate()

	groundMesh = gfx.MeshLock(gl.TRIANGLES, gl.UNSIGNED_SHORT, int32(6*groundWidth*groundHeight), gfx.VertexLayout, uint32(6*groundWidth*groundHeight))
	// TODO REMOVE var groundScale float32 = 20.0
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
	/*skies = append(skies, "ennis")
	  skies = append(skies, "field")
	  skies = append(skies, "footprint_court")
	  skies = append(skies, "helipad")
	  skies = append(skies, "papermill")
	  skies = append(skies, "pisa")
	  skies = append(skies, "studio_grey")
	  skies = append(skies, "studio_red_green")*/
	for i := 0; i < len(skies); i++ {
		skySpeculars = append(skySpeculars, gfx.TextureCreateCube(fmt.Sprintf("%s/textures/environments/%s/specular/specular_.hdr", app.DataDir(), skies[i]), gfx.FORMAT_RGB_16F, 0, 0, 32, false))
		skyDiffuses = append(skyDiffuses, gfx.TextureCreateCube(fmt.Sprintf("%s/textures/environments/%s/diffuse/diffuse_.hdr", app.DataDir(), skies[i]), gfx.FORMAT_RGB_16F, 0, 0, 0, false))
	}

	brdf = gfx.TextureCreate(fmt.Sprintf("%s/textures/environments/brdf.png", app.DataDir()), gfx.FORMAT_R8_G8_B8_A8, 0, 0, 0, false, false)

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

	if app.Run("skyfx: Octree Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		gfx.MeshUnlock(skyMesh)
		gfx.MeshUnlock(groundMesh)
		app.Destroy()
	}
}
