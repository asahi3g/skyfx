package app

import (
	"fmt"
	"os"
	"skyfx/fps"
	"skyfx/gfx"
	"skyfx/math"
	v2 "skyfx/math/v2"
	v4 "skyfx/math/v4"
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// TODO : handle position in options
// TODO : handle fullscreen in options

// Globals ...
var g_app Application
var fakeExit float64 // TODO : REMOVE
var g_Window *glfw.Window

// Constants ...
var APP_HINT_NONE int32 = 0
var APP_HINT_DEBUG int32 = 1
var APP_HINT_FULLSCREEN int32 = 2
var APP_HINT_RESIZABLE int32 = 4

const (
	APP_NONE EventType = iota
	APP_START
	APP_STOP
	APP_KEYBOARD
	APP_MOUSE
	APP_FOCUS_ON
	APP_FOCUS_OFF
	APP_FRAMEBUFFER_SIZE
	APP_WINDOW_SIZE
	APP_WINDOW_POSITION
	APP_PAINT
	APP_CURSOR_POS   // TODO : to deprecate
	APP_MOUSE_BUTTON // TODO : to deprecate
)

const (
	ACTION_RELEASE ActionType = 0 // TODO : break from glfw.Action
	ACTION_PRESS              = 1
	ACTION_REPEAT             = 2
	ACTION_MOVE               = 3
)

type ActionType uint32
type EventType uint32
type AppEvent struct {
	eventType EventType
	action    ActionType
	time      uint64
	x         float64
	y         float64
	key       int32
	scancode  int64
	mods      int32
}

var eventCount int

var events []AppEvent
var polled []AppEvent
var polledCount int

var mutex sync.Mutex

// Options ...
type Options struct {
	width     uint32
	height    uint32
	fps       int32
	glVersion string
	hints     int32
	dataDir   string
}

// Application ...
type Application struct {
	options    Options
	exit       bool
	visible    bool
	name       string
	window     string
	major      int32
	minor      int32
	glVersion  int32
	x          int32
	y          int32
	width      uint32
	height     uint32
	xscale     float32
	yscale     float32
	xy_scale   math.V2
	xyxy_scale math.V4
	fullscreen bool

	prfSwap     fps.ProfileId
	prfFrame    fps.ProfileId
	prfUpdate   fps.ProfileId
	prfSound    fps.ProfileId
	prfUIUpdate fps.ProfileId
	prfUIResize fps.ProfileId
	prfRender   fps.ProfileId

	currentArgIsValid bool
	currentArg        string
	userHelp          string

	onParse             onParseCallback
	onStart             onStartCallback
	onStop              onStopCallback
	onKeyboard          onKeyboardCallback
	onMouse             onMouseCallback
	onFramebufferResize onResizeCallback
	onWindowResize      onResizeCallback
	onWindowPosition    onPositionCallback
	onUpdate            onUpdateCallback
	onRender            onRenderCallback
}

type onParseCallback func()
type onStartCallback func()
type onStopCallback func()
type onKeyboardCallback func(*KeyboardEvent)
type onMouseCallback func(*MouseEvent)
type onResizeCallback func(float64, float64)
type onPositionCallback func(float64, float64)
type onUpdateCallback func()
type onRenderCallback func()

// Help ...
func Help(message string, exitCode int32) {
	fmt.Printf(message)
	fmt.Printf("\nUsage:\n")
	fmt.Printf("++help      : Prints this message.\n")
	fmt.Printf("++width     : Width in pixels\n")
	fmt.Printf("++height    : Height in pixels\n")
	fmt.Printf("++fps       : Target framerate, vsyng disabled when >= 0\n")
	fmt.Printf("++hints     : Window hint flags (fullscreen, resizable)\n")
	fmt.Printf("++data      : Resource path\n")
	fmt.Printf("++glVersion : Select opengl version\n")
	if len(g_app.userHelp) > 0 {
		fmt.Printf(g_app.userHelp)
	}
	os.Exit(int(exitCode))
}

func windowCreate(width int, height int, name string) {
	if g_Window != nil {
		panic("Windows is already created")
	}

	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err)
	}
	g_Window = window
}

func makeCurrentContext() {
	g_Window.MakeContextCurrent()
}

func Update() {
	glfw.PollEvents()
	PollEvents()
}

// Run ...
func Run(name string, width uint32, height uint32, framePerSecond int32, dataDir string, glVersion int32, fpsUpdate uint32) (success bool) {
	g_app.options.width = width
	g_app.options.height = height
	g_app.options.fps = framePerSecond
	g_app.options.dataDir = dataDir
	g_app.options.glVersion = gfx.GetGLVersionToString(glVersion)

	path, _ /*err*/ := os.Getwd()
	fmt.Printf("Current Directory: %s\n", path)

	parse()

	g_app.exit = false
	g_app.window = name
	g_app.name = name
	g_app.glVersion = gfx.GetGLVersionFromString(g_app.options.glVersion)
	fmt.Printf("OpenGL Version :  %d, '%s'\n", g_app.glVersion, g_app.options.glVersion)

	if g_app.glVersion == gfx.GL_VERSION_DS_3_2 {
		g_app.major = 3
		g_app.minor = 2
	} else if g_app.glVersion == gfx.GL_VERSION_ES_3_1 {
		g_app.major = 3
		g_app.minor = 1
	} else {
		success = false
		return
	}
	fps.Init(uint32(Fps()))

	g_app.prfSwap = fps.CreateProfile("swap")
	g_app.prfFrame = fps.CreateProfile("frame")
	g_app.prfUpdate = fps.CreateProfile("update")
	g_app.prfSound = fps.CreateProfile("sound")
	g_app.prfUIUpdate = fps.CreateProfile("ui update")
	g_app.prfUIResize = fps.CreateProfile("ui resize")
	g_app.prfRender = fps.CreateProfile("render")

	Init()

	for Running() && (fakeExit < 60.0) {
		fps.BeginUpdate(fpsUpdate)
		Update()
		if Visible() {
			fps.StartProfile(g_app.prfSwap)
			{
				makeCurrentContext()
				fps.StartProfile(g_app.prfFrame)
				{
					// Update ...
					fps.StartProfile(g_app.prfUpdate)
					{
						fps.StartProfile(g_app.prfSound)
						// TODO SND snd.Update()
						fps.StopProfile(g_app.prfSound)

						fps.StartProfile(g_app.prfUIUpdate)
						//fakeExit = fakeExit + fps.DeltaSecond()
						if g_app.onUpdate != nil {
							g_app.onUpdate()
						}
						// TODO GUI gui.Update(fps.DeltaSecond())
						fps.StopProfile(g_app.prfUIUpdate)

						fps.StartProfile(g_app.prfUIResize)
						// TODO GUI gui.Resize()
						fps.StopProfile(g_app.prfUIResize)
					}
					fps.StopProfile(g_app.prfUpdate)

					// Render ...
					fps.StartProfile(g_app.prfRender)
					if g_app.onRender != nil {
						g_app.onRender()
					}
					// TODO GUI gui.Render()
					fps.StopProfile(g_app.prfRender)
				}
				fps.StopProfile(g_app.prfFrame)
				g_Window.SwapBuffers()
			}
			fps.StopProfile(g_app.prfSwap)
		}
		fps.EndUpdate()
	}
	return
}

// Width ...
func Width() (out uint32) {
	out = g_app.width
	return
}

// Height ...
func Height() (out uint32) {
	out = g_app.height
	return
}

// ScaleX ...
func ScaleX() (out float32) {
	out = g_app.xscale
	return
}

// ScaleY ...
func ScaleY() (out float32) {
	out = g_app.yscale
	return
}

// ScaleXY ...
func ScaleXY() (out math.V2) {
	out = g_app.xy_scale
	return
}

// ScaleXYXY ...
func ScaleXYXY() (out math.V4) {
	out = g_app.xyxy_scale
	return
}

// DataDir ...
func DataDir() (out string) {
	out = g_app.options.dataDir
	return
}

// Fps ...
func Fps() (out int32) {
	out = g_app.options.fps
	return
}

// GLVersion ...
func GLVersion() (out int32) {
	out = g_app.glVersion
	return
}

// Name ...
func Name() (out string) {
	out = g_app.name
	return
}

// WindowName ...
func WindowName() (out string) {
	out = g_app.window
	return
}

// Exit ...
func Exit() {
	fmt.Printf("exiting %s...\n", g_app.name)
	g_app.exit = true
}

// ToggleFullscreen ...
func ToggleFullscreen() {
	var fullscreen bool = g_app.fullscreen == false

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if fullscreen {
		g_Window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		g_Window.SetMonitor(nil, int(g_app.x), int(g_app.y), int(g_app.options.width), int(g_app.options.height), mode.RefreshRate)
	}
	g_Window.MakeContextCurrent()
	g_app.fullscreen = fullscreen
}

// Running ...
func Running() (out bool) { // ISSUE : member function does not contain expressions
	if Visible() {
		out = (g_Window.ShouldClose() == false && g_app.exit == false)
	} else {
		out = true
	}
	return
}

// Visible ...
func Visible() (out bool) {
	out = g_app.visible
	return
}

// Resize ...
func Resize(width uint32, height uint32) {
	/*if g_app.fullscreen == false {
		g_app.options.width = width
		g_app.options.height = height
	}*/

	g_app.width = width
	g_app.height = height
}

// Destroy ...
func Destroy() {
	makeCurrentContext()
	// TODO SND snd.Destroy()
	gfx.Destroy()
}

// SetUserHelp ...
func SetUserHelp(message string) {
	g_app.userHelp = message
}

// GetCurrentArg ...
func GetCurrentArg() (out string) {
	out = g_app.currentArg
	return
}

// SetCurrentArgIsValid ...
func SetCurrentArgIsValid(value bool) {
	g_app.currentArgIsValid = value
}

// SetOnParse ...
func SetOnParse(onParse func()) {
	g_app.onParse = onParse
}

// SetOnStart ...
func SetOnStart(callback onStartCallback) {
	g_app.onStart = callback
}

// SetOnStop ...
func SetOnStop(callback onStopCallback) {
	g_app.onStop = callback
}

// SetOnKeyboard ...
func SetOnKeyboard(callback onKeyboardCallback) {
	g_app.onKeyboard = callback
}

// SetOnMouse ...
func SetOnMouse(callback onMouseCallback) {
	g_app.onMouse = callback
}

// SetOnFramebufferResize ...
func SetOnFramebufferResize(callback onResizeCallback) {
	g_app.onFramebufferResize = callback
}

// SetOnWindowResize ...
func SetOnWindowResize(callback onResizeCallback) {
	g_app.onWindowResize = callback
}

// SetOnWindowPosition ...
func SetOnWindowPosition(callback onPositionCallback) {
	g_app.onWindowPosition = callback
}

// SetOnUpdate ...
func SetOnUpdate(callback onUpdateCallback) {
	g_app.onUpdate = callback
}

// SetOnRender ...
func SetOnRender(callback onRenderCallback) {
	g_app.onRender = callback
}

func PushEvent(eventType EventType) {
	mutex.Lock()
	defer mutex.Unlock()

	event := AppEvent{
		eventType,
		ACTION_RELEASE, // TODO : ACTION_NONE
		uint64(time.Now().UnixNano()),
		0.0,
		0.0,
		0,
		0,
		0,
	}
	if eventCount < len(events) {
		events[eventCount] = event
	} else {
		events = append(events, event)
	}
	eventCount++
}

func PushWindowSizeEvent(width float64, height float64) {
	index := eventCount
	PushEvent(APP_WINDOW_SIZE)
	events[index].x = width
	events[index].y = height
}

func PushWindowPositionEvent(x float64, y float64) {
	index := eventCount
	PushEvent(APP_WINDOW_POSITION)
	events[index].x = x
	events[index].y = y
}

func PushFramebufferSizeEvent(width float64, height float64) {
	index := eventCount
	PushEvent(APP_FRAMEBUFFER_SIZE)
	events[index].x = width
	events[index].y = height
}

func PushKeyboardEvent(action ActionType, key int32, scancode int64, mods int32) {
	index := eventCount
	PushEvent(APP_KEYBOARD)
	events[index].action = action
	events[index].key = key
	events[index].scancode = scancode
	events[index].mods = mods
}

func PushMouseEvent(eventType EventType, action ActionType, button int32, touch int64, mods int32, x float64, y float64) {
	index := eventCount
	PushEvent(eventType)
	events[index].action = action
	events[index].key = button
	events[index].scancode = touch
	events[index].mods = mods
	events[index].x = x
	events[index].y = y
}

func parse() {
	// TODO FIX LATER
	//var argc int = len(os.Args) // ISSUE : doubious compilation error : function 'len' expected receiving variable of type 'int32'; 'string' was provided

	// >> tmp issue #214
	// var width int32
	// var height int32
	// var fps int32
	// var hints int32
	// var dataDir string
	// var help bool
	// var glVersion string
	// // << tmp issue #214

	// var helpMatch bool
	// var widthMatch bool
	// var heightMatch bool
	// var fpsMatch bool
	// var hintsMatch bool
	// var dataDirMatch bool
	// var glVersionMatch bool

	// var hintNames []string
	// hintNames = []string{"fullscreen", "resizable"}
	// var hintValues []int32
	// hintValues = []int32{APP_HINT_FULLSCREEN, APP_HINT_RESIZABLE}

	// for a := 0; a < argc; a++ {
	// 	arg := os.Args[a]
	// 	if args.Bool(arg, "help", &help, &helpMatch) {
	// 		Help("", 0)
	// 	}

	// 	if args.I32(arg, "width", &width, &widthMatch) {
	// 		if width < 0 || width > 65536 {
	// 			Help(fmt.Sprintf("invalid value %s\n", arg), cx.PANIC)
	// 		}
	// 		g_app.options.width = width
	// 		continue
	// 	}

	// 	if args.I32(arg, "height", &height, &heightMatch) {
	// 		if height < 0 || height > 65536 {
	// 			Help(fmt.Sprintf("invalid value %s\n", arg), cx.PANIC)
	// 		}
	// 		g_app.options.height = height
	// 		continue
	// 	}

	// 	if args.I32(arg, "fps", &fps, &fpsMatch) {
	// 		if fps < 0 || fps > 65536 {
	// 			Help(fmt.Sprintf("invalid value %s\n", arg), cx.PANIC)
	// 		}
	// 		g_app.options.fps = fps
	// 		continue
	// 	}

	// 	if args.Str(arg, "data", &dataDir, &dataDirMatch) {
	// 		g_app.options.dataDir = dataDir
	// 		continue
	// 	}

	// 	if args.Str(arg, "glVersion", &glVersion, &glVersionMatch) {
	// 		g_app.options.glVersion = glVersion
	// 		if gfx.GetGLVersionFromString(glVersion) == gfx.GL_VERSION_NONE {
	// 			Help(fmt.Sprintf("invalid value %s\n", arg), cx.PANIC)
	// 		}
	// 		continue
	// 	}

	// 	if args.Flags(arg, "hints", &hints, &hintsMatch, hintNames, hintValues) {
	// 		g_app.options.hints = hints
	// 		continue
	// 	}

	// 	if g_app.hasOnParse == 1 {
	// 		g_app.currentArg = arg
	// 		call_i32_i32(g_app.onParse, a, argc)
	// 		if g_app.currentArgIsValid == 1 {
	// 			continue
	// 		}
	// 	}

	// 	Help(fmt.Sprintf("invalid argument %s", arg), cx.ASSERT)
	// }
}

func Init() {
	g_app.visible = true
	fmt.Printf("app start : %s...\n", g_app.name)

	glfw.Init() // ##0 terminate

	var resizable int = glfw.False
	if (g_app.options.hints & APP_HINT_RESIZABLE) == APP_HINT_RESIZABLE {
		resizable = glfw.True
	}

	glfw.WindowHint(glfw.Resizable, resizable)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.ContextVersionMajor, int(g_app.major))
	glfw.WindowHint(glfw.ContextVersionMinor, int(g_app.minor))
	//glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.True)
	//glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)

	var xscale float32
	var yscale float32
	xscale, yscale = glfw.GetPrimaryMonitor().GetContentScale()

	g_app.xscale = xscale
	g_app.yscale = yscale
	g_app.xy_scale = v2.Make(xscale, yscale)
	g_app.xyxy_scale = v4.Make(xscale, yscale, xscale, yscale)

	var windowWidth int = int(float32(g_app.options.width) / xscale)
	var windowHeight int = int(float32(g_app.options.height) / yscale)

	windowCreate(windowWidth, windowHeight, g_app.name)
	makeCurrentContext()

	g_Window.SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyboardEvent(ActionType(action), int32(key), int64(scancode), int32(mods))
		})

	g_Window.SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushMouseEvent(APP_MOUSE, ACTION_MOVE, 0, -1, 0, xpos, ypos)
		})

	g_Window.SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			x, y := w.GetCursorPos()
			PushMouseEvent(APP_MOUSE, ActionType(action), int32(key), -1, int32(mods), x, y)
		})

	g_Window.SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})

	g_Window.SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})

	g_Window.SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(float64(x), float64(y)) // TODO to deprecate, use float64
		})

	PushEvent(APP_START)

	x, y := g_Window.GetPos()
	g_app.x = int32(x)
	g_app.y = int32(y)

	w, h := g_Window.GetFramebufferSize()
	g_app.width = uint32(w)
	g_app.height = uint32(h)
	g_app.width = g_app.width / 2
	g_app.height = g_app.height / 2
	g_app.width = g_app.options.width
	g_app.height = g_app.options.height

	var swapInterval int = 0
	if g_app.options.fps > 0 {
		swapInterval = 1 // should be based on targetFps
	}

	glfw.SwapInterval(swapInterval)

	if (g_app.options.hints & APP_HINT_FULLSCREEN) == APP_HINT_FULLSCREEN {
		ToggleFullscreen()
	}

	SetOnFramebufferResize(OnFramebufferResize)
	SetOnWindowResize(OnWindowResize)
	SetOnWindowPosition(OnWindowPosition)

	gfx.Init(Width(), Height(), DataDir(), GLVersion())
	// TODO GUI gui.Init(Width(), Height(), DataDir())
	// TODO SND snd.Init(64, 64)

	g_app.visible = true
}

func purgeEvents() {
	mutex.Lock()
	defer mutex.Unlock()

	if eventCount > 0 {
		polled = append(polled[0:0], events[0:eventCount]...)
		eventCount = 0
	}
}

func PollEvents() {

	var mouseEvent MouseEvent
	var keyboardEvent KeyboardEvent

	purgeEvents()
	eventCount := len(polled)
	for i := 0; i < eventCount; i++ {
		e := polled[i]
		switch e.eventType {
		case APP_START:
			if g_app.onStart != nil {
				g_app.onStart()
			}
		case APP_STOP:
			if g_app.onStop != nil {
				g_app.onStop()
			}
		case APP_KEYBOARD:
			if g_app.onKeyboard != nil {
				keyboardEvent.Key = e.key
				keyboardEvent.Scancode = e.scancode
				keyboardEvent.Action = e.action
				keyboardEvent.Mods = e.mods
				g_app.onKeyboard(&keyboardEvent)
			}
		case APP_MOUSE:
			if g_app.onMouse != nil {
				mouseEvent.Button = e.key
				mouseEvent.Touch = e.scancode
				mouseEvent.Action = e.action
				mouseEvent.Mods = e.mods
				mouseEvent.X = e.x
				mouseEvent.Y = e.y
				g_app.onMouse(&mouseEvent)
			}
		case APP_FRAMEBUFFER_SIZE:
			if g_app.onFramebufferResize != nil {
				g_app.onFramebufferResize(e.x, e.y)
			}
		case APP_WINDOW_SIZE:
			if g_app.onWindowResize != nil {
				g_app.onWindowResize(e.x, e.y)
			}
		case APP_WINDOW_POSITION:
			if g_app.onWindowPosition != nil {
				g_app.onWindowPosition(e.x, e.y)
			}
		}
	}
	polled = polled[0:0]
}

func OnFramebufferResize(width float64, height float64) {
	fmt.Printf("------------------------->>>>> FramebufferResizeCallback %f, %f\n", width, height)
	gfx.Resize(uint32(width), uint32(height)) // TODO remove cast
}

func OnWindowResize(width float64, height float64) {
	fmt.Printf("------------------------->>>>> WindowResizeCallback %f, %f\n", width, height)
	/*app.Resize(width, height)
	gfx.Resize(width * 2, height * 2)*/
	if g_app.fullscreen == false {
		g_app.options.width = uint32(width)   // TODO remove cast
		g_app.options.height = uint32(height) // TODO remove cast
	}
}

func OnWindowPosition(x float64, y float64) {
	if g_app.fullscreen == false {
		g_app.x = int32(x) // TODO remove cast
		g_app.y = int32(y) // TODO remove cast
	}
}
