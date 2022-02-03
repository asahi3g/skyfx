package app

import (
	"skyfx/math"
	"skyfx/math/m44"
	"skyfx/math/q4"
	v2 "skyfx/math/v2"
	v3 "skyfx/math/v3"
	v4 "skyfx/math/v4"
	"skyfx/utils"
)

// TODO : rotation based on mouse
// TODO : translation based on at vector
// TODO : apply on gfx.State
// FIX : matrix ordering
//

// Globals ...
var g_cameras []Camera

// CameraId ...
type CameraId struct {
	camera int32
}

// CameraState ...
type CameraState struct {
	near     float32
	far      float32
	fov      float32
	width    float32
	height   float32
	position math.V3
	yaw      float32
	pitch    float32
	//orientation v4
}

// Camera ...
type Camera struct {
	id CameraId

	previous CameraState
	current  CameraState
	next     CameraState

	nearTime     float32
	farTime      float32
	fovTime      float32
	widthTime    float32
	heightTime   float32
	positionTime float32
	yawTime      float32
	pitchTime    float32
	//orientationTime float32

	pitchQuaternion math.V4
	yawQuaternion   math.V4

	transform   math.M44
	view        math.M44
	projection  math.M44
	viewProj    math.M44
	invViewProj math.M44
	invProj     math.M44

	xdir             float32
	zdir             float32
	leftDrag         float32
	leftDragVelocity float32
	dragPosition     math.V2
	dragDelta        math.V2
	previousTarget   math.V3

	translationSpeed float32
	//rotationSpeed float32
	yawSpeed   float32
	pitchSpeed float32
}

// CameraIsValid ...
func CameraIsValid(id CameraId) (out bool) {
	out = id.camera >= 0 && id.camera < int32(len(g_cameras))
	return
}

// CameraGetPosition ...
func CameraGetPosition(id CameraId) (out math.V3) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].current.position
	return
}

// CameraGetPitchQuaternion ...
func CameraGetPitchQuaternion(id CameraId) (out math.V4) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].pitchQuaternion
	return
}

// CameraGetYawQuaternion ...
func CameraGetYawQuaternion(id CameraId) (out math.V4) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].yawQuaternion
	return
}

// CameraGetTransform. ...
func CameraGetTransform(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].transform
	return
}

// CameraGetView ...
func CameraGetView(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].view
	return
}

// CameraGetViewProj ...
func CameraGetViewProj(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].viewProj
	return
}

// CameraGetInvProj ...
func CameraGetInvProj(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].invProj
	return
}

// CameraGetInvViewProj ...
func CameraGetInvViewProj(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].invViewProj
	return
}

// CameraGetNear ...
func CameraGetNear(id CameraId) (out float32) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].current.near
	return
}

// CameraGetFar ...
func CameraGetFar(id CameraId) (out float32) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].current.far
	return
}

// CameraGetRight ...
func CameraGetRight(id CameraId) (out math.V3) {
	utils.PanicIfNot(CameraIsValid(id), "")
	var view math.M44 = g_cameras[id.camera].view
	out = v3.Make(view.V00, view.V10, view.V20)
	return
}

// CameraGetUp ...
func CameraGetUp(id CameraId) (out math.V3) {
	utils.PanicIfNot(CameraIsValid(id), "")
	var view math.M44 = g_cameras[id.camera].view
	out = v3.Make(view.V01, view.V11, view.V21)
	return
}

// CameraGetAt ...
func CameraGetAt(id CameraId) (out math.V3) {
	utils.PanicIfNot(CameraIsValid(id), "")
	var view math.M44 = g_cameras[id.camera].view
	out = v3.Make(view.V02, view.V12, view.V22)
	return
}

// CameraGetProjection ...
func CameraGetProjection(id CameraId) (out math.M44) {
	utils.PanicIfNot(CameraIsValid(id), "")
	out = g_cameras[id.camera].projection
	return
}

// CameraCreate ...
func CameraCreate() (out CameraId) {
	out.camera = int32(len(g_cameras))

	var camera Camera
	camera.id = out

	camera.transform = m44.IDENTITY
	camera.view = m44.IDENTITY
	camera.projection = m44.IDENTITY
	camera.viewProj = m44.IDENTITY
	camera.invViewProj = m44.IDENTITY
	camera.invProj = m44.IDENTITY

	g_cameras = append(g_cameras, camera)

	utils.PanicIfNot(CameraIsValid(out), "")
	return
}

// CameraSetProjection ...
func CameraSetProjection(id CameraId, near float32, far float32, fov float32, width float32, height float32, force bool) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var cameraIndex int32 = id.camera

	var previous CameraState = g_cameras[cameraIndex].previous
	var current CameraState = g_cameras[cameraIndex].current
	var next CameraState = g_cameras[cameraIndex].next

	var time float32 = 0.0
	if force == true {
		time = 1.0
	}

	if next.near != near {
		next.near = near
		previous.near = current.near
		g_cameras[cameraIndex].nearTime = time
	}

	if next.far != far {
		next.far = far
		previous.far = current.far
		g_cameras[cameraIndex].farTime = time
	}

	if next.fov != fov {
		next.fov = fov
		previous.fov = current.fov
		g_cameras[cameraIndex].fovTime = time
	}

	if next.width != width {
		next.width = width
		previous.width = current.width
		g_cameras[cameraIndex].widthTime = time
	}

	if next.height != height {
		next.height = height
		previous.height = current.height
		g_cameras[cameraIndex].heightTime = time
	}

	g_cameras[cameraIndex].previous = previous
	g_cameras[cameraIndex].next = next
}

// CameraSetPosition...
func CameraSetPosition(id CameraId, position math.V3, force bool) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var cameraIndex int32 = id.camera

	var previous CameraState = g_cameras[cameraIndex].previous
	var current CameraState = g_cameras[cameraIndex].current
	var next CameraState = g_cameras[cameraIndex].next

	if force == true {
		previous.position = current.position
		next.position = position
		current.position = position
		g_cameras[cameraIndex].positionTime = 1.0
	} else {
		if v3.Equ(next.position, position) == false {
			next.position = position
			previous.position = current.position
			g_cameras[cameraIndex].positionTime = 0.0
		}
	}

	g_cameras[cameraIndex].previous = previous
	g_cameras[cameraIndex].next = next
	g_cameras[cameraIndex].current = current
}

// CameraSetYawPitch ...
func CameraSetYawPitch(id CameraId, yaw float32, pitch float32, force bool) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var cameraIndex int32 = id.camera

	var previous CameraState = g_cameras[cameraIndex].previous
	var current CameraState = g_cameras[cameraIndex].current
	var next CameraState = g_cameras[cameraIndex].next

	if force == true {
		previous.yaw = current.yaw
		next.yaw = yaw
		current.yaw = yaw

		previous.pitch = current.pitch
		next.pitch = pitch
		current.pitch = pitch
		g_cameras[cameraIndex].yawTime = 1.0
		g_cameras[cameraIndex].pitchTime = 1.0
	} else {

		if next.yaw != yaw {
			next.yaw = yaw
			previous.yaw = current.yaw
			g_cameras[cameraIndex].yawTime = 0.0
		}

		if next.pitch != pitch {
			next.pitch = pitch
			previous.pitch = current.pitch
			g_cameras[cameraIndex].pitchTime = 0.0
		}
	}

	g_cameras[cameraIndex].previous = previous
	g_cameras[cameraIndex].next = next
	g_cameras[cameraIndex].current = current
}

// CameraSetPositionAndTarget ...
/*func CameraSetPositionAndTarget(id CameraId, position math.V3, target math.V3, force bool) {
    utils.PanicIfNot(CameraIsValid(id), "")

    var direction math.V3 = v3.Normalize(v3.Sub(position, target))
    CameraSetPositionAndDirection(id, position, direction, force)
}*/

// CameraSetPositionAndDirection ...
/*func CameraSetPositionAndDirection(id CameraId, position math.V3, direction math.V3, force bool) {
    utils.PanicIfNot(CameraIsValid(id), "")

    var at math.V3 = direction
    //at.X = 0.0 - direction.X
    //at.Y = 0.0 - direction.Y
    //at.Z = 0.0 - direction.Z
    at = v3.Normalize(at)

    var up math.V3 = v3.Make(0.0, 1.0, 0.0)
    if math.Abs_f32(v3.Dot(at, up)) > 0.9999 {
        up = v3.UP
    }

    var right math.V3 = v3.Normalize(v3.Cross(up, at))
    up = v3.Normalize(v3.Cross(at, right))

    var m00 float32 = right.X
    var m01 float32 = up.X
    var m02 float32 = at.X

    var m10 float32 = right.Y
    var m11 float32 = up.Y
    var m12 float32 = at.Y

    var m20 float32 = right.Z
    var m21 float32 = up.Z
    var m22 float32 = at.Z

    var ox float32
    var oy float32
    var oz float32
    var ow float32

    var num8 float32 = m00 + m11 + m22
    if num8 > 0.0 {
        var num float32 = math.Sqrt_f32(num8 + 1.0)
        ow = num * 0.5
        num = 0.5 / num
        ox = (m12 - m21) * num
        oy = (m20 - m02) * num
        oz = (m01 - m10) * num
    } else if (m00 >= m11) && (m00 >= m22) {
        var num7 float32 = math.Sqrt_f32(1.0 + m00 - m11 - m22)
        var num4 float32 = 0.5 / num7
        ox = 0.5 * num7
        oy = (m01 + m10) * num4
        oz = (m02 + m20) * num4
        ow = (m12 - m21) * num4
    } else if (m11 > m22) {
        var num6 float32 = math.Sqrt_f32(1.0 + m11 - m00 - m22)
        var num3 float32 = 0.5 / num6
        ox = (m10 + m01) * num3
        oy = 0.5 * num6
        oz = (m21 + m12) * num3
        ow = (m20 - m02) * num3
    } else {
        var num5 float32 = math.Sqrt_f32(1.0 + m22 - m00 - m11)
        var num2 float32 = 0.5 / num5
        ox = (m20 + m02) * num2
        oy = (m21 + m12) * num2
        oz = 0.5 * num5
        ow = (m01 - m10) * num2
    }

    CameraSetPositionAndOrientation(id, position, v4.Make(0.0-ox, 0.0-oy, 0.0-oz, ow), force)
}*/

// CameraSetPositionAndOrientation ...
/*func CameraSetPositionAndOrientation(id CameraId, position math.V3, orientation math.V4, force bool) {
    utils.PanicIfNot(CameraIsValid(id), "")

    var cameraIndex int32 = id.camera

    var previous CameraState = g_cameras[cameraIndex].previous
    var current CameraState = g_cameras[cameraIndex].current
    var next CameraState = g_cameras[cameraIndex].next

    var time float32 = 0.0
    if force == true {
        time = 1.0
    }

    if v3.equ(next.position, position) == false {
        next.position = position
        previous.position = current.position
        g_cameras[cameraIndex].positionTime = time
    }

    var diff float32 = v4.sqlength(v4.sub(next.orientation, orientation))
    if diff > 0.0001 {
        next.orientation = orientation
        previous.orientation = current.orientation
        g_cameras[cameraIndex].orientationTime = time
    }

    g_cameras[cameraIndex].previous = previous
    g_cameras[cameraIndex].next = next
}*/

// TpsCameraProcess ...
func TpsCameraProcessMouse(id CameraId, event *MouseEvent) {
	utils.PanicIfNot(CameraIsValid(id), "")
	cameraProcessMouse(id, event)
}

// TpsCameraUpdate ...
func TpsCameraUpdate(id CameraId, deltaTime float64, target math.V3, distance float32) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var cameraIndex int32 = id.camera

	var deltaPosition math.V3 = v3.Sub(target, g_cameras[cameraIndex].previousTarget)
	g_cameras[cameraIndex].previousTarget = target

	var yawPitch math.V2 = cameraUpdateYawPitch(id, deltaTime)
	var qyaw math.V4 = q4.From_yaw_pitch_roll(yawPitch.X, 0.0, 0.0)
	var qpitch math.V4 = q4.From_yaw_pitch_roll(0.0, yawPitch.Y, 0.0)
	var orientation math.V4 = q4.Mul(qyaw, qpitch)
	orientation = v4.Normalize(orientation)
	//var orientation math.V4 = q4.from_yaw_pitch_roll(yawPitch.X, yawPitch.Y, 0.0)

	var direction math.V3
	direction.X = (2.0*orientation.X*orientation.Z + 2.0*orientation.Y*orientation.W)
	direction.Y = (2.0*orientation.Y*orientation.Z - 2.0*orientation.X*orientation.W)
	direction.Z = 1.0 - 2.0*orientation.X*orientation.X - 2.0*orientation.Y*orientation.Y
	direction = v3.Mulf(v3.Normalize(direction), distance)

	var currentPosition math.V3 = g_cameras[cameraIndex].current.position
	currentPosition = v3.Add(currentPosition, deltaPosition)
	CameraSetPosition(id, currentPosition, true)

	var position math.V3 = v3.Add(target, direction)
	CameraSetPosition(id, position, false)
	CameraSetYawPitch(id, yawPitch.X, yawPitch.Y, false)
	cameraUpdate(id, deltaTime)
}

// FpsCameraProcess ...
func FpsCameraProcessMouse(id CameraId, event *MouseEvent) {
	utils.PanicIfNot(CameraIsValid(id), "")
	cameraProcessMouse(id, event)
}

// FpsCameraUpdate ...
func FpsCameraUpdate(id CameraId, deltaTime float64, position math.V3) {
	utils.PanicIfNot(CameraIsValid(id), "")
	var yawPitch math.V2 = cameraUpdateYawPitch(id, deltaTime)
	CameraSetPosition(id, position, true)
	CameraSetYawPitch(id, yawPitch.X, yawPitch.Y, false)
	cameraUpdate(id, deltaTime)
}

// FreeCameraProcess ...
func FreeCameraProcessMouse(id CameraId, event *MouseEvent) {
	utils.PanicIfNot(CameraIsValid(id), "")
	cameraProcessMouse(id, event)
}
func FreeCameraProcessKeyboard(id CameraId, event *KeyboardEvent) {
	utils.PanicIfNot(CameraIsValid(id), "")
	cameraProcessKeyboard(id, event)
}

// FreeCameraUpdate ...
func FreeCameraUpdate(id CameraId, deltaTime float64, rotationSpeed float32, translationSpeed float32) {
	utils.PanicIfNot(CameraIsValid(id), "")

	rotationSpeed = 0.5 * rotationSpeed
	translationSpeed = 30.0 * translationSpeed

	var dt float32 = float32(deltaTime)

	var cameraIndex int32 = id.camera
	var view math.M44 = g_cameras[cameraIndex].view

	var right math.V3 = v3.Make(view.V00, view.V10, view.V20)
	// var up math.V3 = v3.Make(view.V01, view.V11, view.V21)
	var at math.V3 = v3.Make(view.V02, view.V12, view.V22)

	var dx float32 = g_cameras[cameraIndex].xdir * dt * translationSpeed
	var dz float32 = g_cameras[cameraIndex].zdir * dt * translationSpeed

	var position math.V3 = g_cameras[cameraIndex].next.position
	position.X = position.X + dx*right.X + dz*at.X
	position.Y = position.Y + dx*right.Y + dz*at.Y
	position.Z = position.Z + dx*right.Z + dz*at.Z

	var yawPitch math.V2 = cameraUpdateYawPitch(id, deltaTime)
	CameraSetPosition(id, position, false)
	CameraSetYawPitch(id, yawPitch.X, yawPitch.Y, false)
	cameraUpdate(id, deltaTime)
}

func cameraProcessKeyboard(id CameraId, event *KeyboardEvent) {
	var key int32 = event.Key
	var action ActionType = event.Action
	var mods int32 = event.Mods

	var cameraIndex int32 = id.camera
	var xdir float32 = g_cameras[cameraIndex].xdir
	var zdir float32 = g_cameras[cameraIndex].zdir

	if mods <= 0 {
		if key == KEYCODE_LEFT || key == KEYCODE_A {
			if action == KEY_PRESS {
				xdir = xdir - 1.0
			} else if action == KEY_RELEASE {
				xdir = xdir + 1.0
			}
		} else if key == KEYCODE_RIGHT || key == KEYCODE_D {
			if action == KEY_PRESS {
				xdir = xdir + 1.0
			} else if action == KEY_RELEASE {
				xdir = xdir - 1.0
			}
		} else if key == KEYCODE_DOWN || key == KEYCODE_S {
			if action == KEY_PRESS {
				zdir = zdir + 1.0
			} else if action == KEY_RELEASE {
				zdir = zdir - 1.0
			}
		} else if key == KEYCODE_UP || key == KEYCODE_W {
			if action == KEY_PRESS {
				zdir = zdir - 1.0
			} else if action == KEY_RELEASE {
				zdir = zdir + 1.0
			}
		}
	}

	g_cameras[cameraIndex].xdir = xdir
	g_cameras[cameraIndex].zdir = zdir
}

func cameraProcessMouse(id CameraId, event *MouseEvent) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var mods int32 = event.Mods
	var button int32 = event.Button
	var state int32 = int32(event.Action)

	var cameraIndex int32 = id.camera

	var leftDrag float32 = g_cameras[cameraIndex].leftDrag
	var leftDragVelocity float32 = g_cameras[cameraIndex].leftDragVelocity
	var dragDelta math.V2

	var position math.V2 = v2.Make(float32(event.X), float32(event.Y))
	if mods <= 0 {
		if state == MOUSE_PRESS {
			if button == BUTTON_LEFT {
				leftDrag = 1.0
				leftDragVelocity = 1.0
				g_cameras[cameraIndex].dragPosition = position
			}
		} else if state == MOUSE_RELEASE {
			if button == BUTTON_LEFT {
				leftDrag = 0.0
			}
		} else if state == MOUSE_MOVE {
			if leftDrag > 0.0 {
				var dragPosition math.V2 = g_cameras[cameraIndex].dragPosition
				dragDelta = v2.Sub(position, dragPosition)
				g_cameras[cameraIndex].dragPosition = position
				leftDragVelocity = 1.0
			}
		}
	}

	leftDragVelocity = leftDragVelocity * 0.98
	g_cameras[cameraIndex].leftDrag = leftDrag
	g_cameras[cameraIndex].leftDragVelocity = leftDragVelocity
	g_cameras[cameraIndex].dragDelta = dragDelta
}

/*func cameraYawPitchToQuaternion(yaw float32, pitch float32) (out math.V4) {
}*/

func cameraUpdateYawPitch(id CameraId, deltaTime float64) (out math.V2) {
	utils.PanicIfNot(CameraIsValid(id), "")
	var rotationSpeed float32 = 0.5

	var dt float32 = float32(deltaTime)
	var cameraIndex int32 = id.camera

	out.X = g_cameras[cameraIndex].next.yaw
	out.Y = g_cameras[cameraIndex].next.pitch

	var leftDrag float32 = g_cameras[cameraIndex].leftDrag
	var leftDragVelocity float32 = g_cameras[cameraIndex].leftDragVelocity
	if (leftDrag * leftDragVelocity) > 0.0 {
		var dragDelta math.V2 = g_cameras[cameraIndex].dragDelta

		out.X = out.X - leftDrag*dragDelta.X*dt*rotationSpeed
		out.Y = out.Y - leftDrag*dragDelta.Y*dt*rotationSpeed
		g_cameras[cameraIndex].dragDelta = v2.ZERO

		/*var qyaw math.V4 = q4.from_yaw_pitch_roll(yaw, 0.0, 0.0)
		  var qpitch math.V4 = q4.from_yaw_pitch_roll(0.0, pitch, 0.0)
		  orientation = q4.mul(qyaw, orientation)
		  orientation = q4.mul(orientation, qpitch)
		  orientation = v4.normalize(orientation)*/
	}
	return
}

func cameraUpdate(id CameraId, deltaTime float64) {
	utils.PanicIfNot(CameraIsValid(id), "")

	var cameraIndex int32 = id.camera

	var previous CameraState = g_cameras[cameraIndex].previous
	var next CameraState = g_cameras[cameraIndex].next
	var current CameraState = g_cameras[cameraIndex].current

	var nearTime float32 = g_cameras[cameraIndex].nearTime
	var farTime float32 = g_cameras[cameraIndex].farTime
	var fovTime float32 = g_cameras[cameraIndex].fovTime
	var widthTime float32 = g_cameras[cameraIndex].widthTime
	var heightTime float32 = g_cameras[cameraIndex].heightTime
	var positionTime float32 = g_cameras[cameraIndex].positionTime
	//var orientationTime float32 = g_cameras[cameraIndex].orientationTime
	var yawTime float32 = g_cameras[cameraIndex].yawTime
	var pitchTime float32 = g_cameras[cameraIndex].pitchTime

	var dt float32 = float32(deltaTime)
	nearTime = nearTime + dt
	farTime = farTime + dt
	fovTime = fovTime + dt
	widthTime = widthTime + dt
	heightTime = heightTime + dt

	var translationSpeed float32 = g_cameras[cameraIndex].translationSpeed
	//var rotationSpeed float32 = g_cameras[cameraIndex].rotationSpeed
	var yawSpeed float32 = g_cameras[cameraIndex].yawSpeed
	var pitchSpeed float32 = g_cameras[cameraIndex].pitchSpeed

	var currentPosition math.V3 = current.position
	var nextPosition math.V3 = next.position
	if v3.Equ(currentPosition, nextPosition) == false {
		translationSpeed = 6.0
	}

	/*var currentOrientation math.V4 = current.orientation
	  var nextOrientation math.V4 = next.orientation
	  if v4.equ(currentOrientation, nextOrientation) == false {
	      rotationSpeed = 5.0
	  }*/

	//var maxYawLerp float32 = 1.0
	//var maxPitchLerp float32 = 1.0

	var currentYaw float32 = current.yaw
	var nextYaw float32 = next.yaw
	if currentYaw != nextYaw {
		yawSpeed = 6.0
	}

	var currentPitch float32 = current.pitch
	var nextPitch float32 = next.pitch
	if currentPitch != nextPitch {
		pitchSpeed = 6.0
	}

	positionTime = positionTime + dt*translationSpeed
	//orientationTime = orientationTime + dt * rotationSpeed
	yawTime = yawTime + dt*yawSpeed
	pitchTime = pitchTime + dt*pitchSpeed

	translationSpeed = math.Max_f32(0.0, translationSpeed-12.0*dt)
	//rotationSpeed = math.Max_f32(0.0, rotationSpeed - 50.0 * dt)
	yawSpeed = math.Max_f32(0.0, yawSpeed-12.0*dt)
	pitchSpeed = math.Max_f32(0.0, pitchSpeed-12.0*dt)

	// TODO : don't lerp if not changes
	current.near = math.Lerpsat_f32(previous.near, next.near, nearTime)
	current.far = math.Lerpsat_f32(previous.far, next.far, farTime)
	current.fov = math.Lerpsat_f32(previous.fov, next.fov, fovTime)
	current.width = math.Lerpsat_f32(previous.width, next.width, widthTime)
	current.height = math.Lerpsat_f32(previous.height, next.height, heightTime)
	current.position = v3.Lerpsatf(previous.position, nextPosition, positionTime)
	//current.orientation = v4.lerpsatf(previous.orientation, nextOrientation, orientationTime)
	//current.orientation = q4.slerp(previous.orientation, next.orientation, math.Sat_f32(orientationTime))
	//current.orientation = v4.normalize(current.orientation)

	//g_cameras[cameraIndex].previous.pitch = current.pitch
	//g_cameras[cameraIndex].previous.yaw = current.yaw

	current.yaw = math.Lerpsat_f32(previous.yaw, nextYaw, yawTime)         // v1.clamp(yawTime, 0.0, maxYawLerp))
	current.pitch = math.Lerpsat_f32(previous.pitch, nextPitch, pitchTime) // v1.clamp(pitchTime, 0.0, maxPitchLerp))

	g_cameras[cameraIndex].nearTime = math.Sat_f32(nearTime)
	g_cameras[cameraIndex].farTime = math.Sat_f32(farTime)
	g_cameras[cameraIndex].fovTime = math.Sat_f32(fovTime)
	g_cameras[cameraIndex].widthTime = math.Sat_f32(widthTime)
	g_cameras[cameraIndex].heightTime = math.Sat_f32(heightTime)
	g_cameras[cameraIndex].positionTime = math.Sat_f32(positionTime)
	g_cameras[cameraIndex].yawTime = math.Sat_f32(yawTime)
	g_cameras[cameraIndex].pitchTime = math.Sat_f32(pitchTime)

	//g_cameras[cameraIndex].orientationTime = math.Sat_f32(orientationTime)
	g_cameras[cameraIndex].translationSpeed = translationSpeed
	//g_cameras[cameraIndex].rotationSpeed = rotationSpeed
	g_cameras[cameraIndex].yawSpeed = yawSpeed
	g_cameras[cameraIndex].pitchSpeed = pitchSpeed

	g_cameras[cameraIndex].current = current

	var projection math.M44 = m44.Make_project(current.near, current.far, current.fov, current.width, current.height)
	g_cameras[cameraIndex].projection = projection

	var position math.V3 = current.position

	var qyaw math.V4 = q4.From_yaw_pitch_roll(current.yaw, 0.0, 0.0)
	var qpitch math.V4 = q4.From_yaw_pitch_roll(0.0, current.pitch, 0.0)
	var orientation math.V4 = q4.Mul(qyaw, qpitch)
	orientation = v4.Normalize(orientation)
	g_cameras[cameraIndex].yawQuaternion = qyaw
	g_cameras[cameraIndex].pitchQuaternion = qpitch
	/*orientation = q4.mul(qyaw, orientation)
	  orientation = q4.mul(orientation, qpitch)
	  orientation = v4.normalize(orientation)*/

	var transform math.M44 = m44.Makev_QT(orientation, position)
	g_cameras[cameraIndex].transform = transform

	var view math.M44 = m44.Inverse(transform) // TODO : compute view matrix without MatrixInvert
	g_cameras[cameraIndex].view = view

	var invProj math.M44 = m44.Inverse(projection)
	g_cameras[cameraIndex].invProj = invProj

	var viewProj math.M44 = m44.MulISSUE(view, projection)
	g_cameras[cameraIndex].viewProj = viewProj

	var invViewProj math.M44 = m44.Inverse(viewProj)
	g_cameras[cameraIndex].invViewProj = invViewProj
}
