package fps

import (
	"fmt"
	"time"
)

// Globals ...
var g_framerate Framerate

// Framerate ...
type Framerate struct {
	currentTime  int64
	previousTime int64
	deltaTime    int64
	fpsTime      int64
	frameTime    int64
	fps          int64

	deltaSecond float64
}

// DeltaSecond ...
func DeltaSecond() (out float64) {
	out = g_framerate.deltaSecond
	return
}

// DeltaNano ...
func DeltaNano() (out int64) {
	out = g_framerate.deltaTime
	return
}

// SecondToNano ...
func SecondToNano(second float64) (nano int64) {
	nano = int64(second * 1000000000.0)
	return
}

// NanoToSecond ...
func NanoToSecond(nano int64) (second float64) {
	second = float64(nano) / 1000000000.0
	return
}

// NanoToMilli ...
func NanoToMilli(nano int64) (milli float64) {
	milli = float64(nano) / 1000000.0
	return
}

// Init ...
func Init(targetFps uint32) {
	g_framerate.currentTime = time.Now().UnixNano()
	g_framerate.previousTime = g_framerate.currentTime
	g_framerate.deltaTime = 0
	g_framerate.fpsTime = 0
	g_framerate.frameTime = 0
	var targetFpsI64 int64 = int64(targetFps)
	if targetFpsI64 > 0 {
		g_framerate.frameTime = SecondToNano(1.0) / targetFpsI64
	}
	g_framerate.fps = 0
}

// BeginUpdate ...
func BeginUpdate(updateStepSecond uint32) {
	g_framerate.currentTime = time.Now().UnixNano()
	g_framerate.deltaTime = g_framerate.currentTime - g_framerate.previousTime
	if NanoToSecond(g_framerate.deltaTime) > 0.5 {
		g_framerate.deltaTime = SecondToNano(0.5)
		fmt.Printf("Discrading frames %f\n", float32(NanoToSecond(g_framerate.deltaTime)))
	}
	g_framerate.deltaSecond = NanoToSecond(g_framerate.deltaTime)

	if g_framerate.deltaTime == g_framerate.currentTime { // TODO : remove ? fixed steps ?
		g_framerate.deltaTime = g_framerate.frameTime
	}

	//var updateStepNano int64 = SecondToNano(float64(updateStepSecond))

	var full bool = false
	if (g_framerate.currentTime - g_framerate.fpsTime) > SecondToNano(float64(updateStepSecond)) {
		g_framerate.fpsTime = g_framerate.currentTime
		var dt int64 = int64(updateStepSecond)
		var frameCount int64 = 0
		if dt > 0 {
			frameCount = (g_framerate.fps + 1) / dt
		}
		fmt.Printf("fps : %d, dt : %fms\n", frameCount, NanoToMilli(g_framerate.deltaTime))
		g_framerate.fps = 0
		full = true
		PrintProfiles() // ISSUE : garbage value when : len(g_profiles)
	} else {
		g_framerate.fps = g_framerate.fps + 1 // ISSUE : g_framerate.fps++ : framerate.cx:78 function 'int64.add' expected input argument of type 'int64'; 'int32' was provided
	}

	ClearProfiles(full) // ISSUE : garbage value when : len(g_profiles)
}

// EndUpdate ...
func EndUpdate() {
	if g_framerate.frameTime > 0 {
		var updateTime int64 = time.Now().UnixNano()
		var deltaUpdateTime int64 = updateTime - g_framerate.currentTime

		if deltaUpdateTime < g_framerate.frameTime {
			var sleepTimeMilli int32 = int32(NanoToMilli(g_framerate.frameTime - deltaUpdateTime))
			time.Sleep(time.Duration(sleepTimeMilli))
		}
	}

	g_framerate.previousTime = g_framerate.currentTime
}
