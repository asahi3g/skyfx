package app

import (
	"math"
)

// Constants ...
var BUTTON_LEFT int32 = 0
var BUTTON_RIGHT int32 = 1

var MOUSE_RELEASE int32 = 0
var MOUSE_PRESS int32 = 1
var MOUSE_REPEATE int32 = 2
var MOUSE_MOVE int32 = 3

const (
	KEY_RELEASE ActionType = 0
	KEY_PRESS              = 1
	KEY_REPEAT             = 2
)

var KEYCODE_SPACE int32 = 32
var KEYCODE_ESCAPE int32 = 256
var KEYCODE_ENTER int32 = 257
var KEYCODE_TAB int32 = 258
var KEYCODE_RIGHT int32 = 262
var KEYCODE_LEFT int32 = 263
var KEYCODE_DOWN int32 = 264
var KEYCODE_UP int32 = 265
var KEYCODE_LEFT_CTRL int32 = 341
var KEYCODE_SHIFT int32 = 340
var KEYCODE_MENU int32 = 343

var KEYCODE_1 int32 = 49
var KEYCODE_2 int32 = 50
var KEYCODE_3 int32 = 51
var KEYCODE_4 int32 = 52
var KEYCODE_5 int32 = 53
var KEYCODE_6 int32 = 54
var KEYCODE_7 int32 = 55
var KEYCODE_8 int32 = 56
var KEYCODE_9 int32 = 57
var KEYCODE_0 int32 = 48

var KEYCODE_A int32 = 65
var KEYCODE_B int32 = 66
var KEYCODE_D int32 = 68
var KEYCODE_S int32 = 83
var KEYCODE_V int32 = 86
var KEYCODE_W int32 = 87

var MOD_NONE int32 = 0
var MOD_SHIFT int32 = 1
var MOD_CTRL int32 = 2
var MOD_ALT int32 = 4
var MOD_MENU int32 = 8

var EVENT_ERROR int32 = 0
var EVENT_UNUSED int32 = 1
var EVENT_CONSUMED int32 = 2

// KeyboardEvent ...
type KeyboardEvent struct {
	Key      int32
	Scancode int64
	Action   ActionType
	Mods     int32
}

// MouseEvent ...
type MouseEvent struct {
	X      float64
	Y      float64
	Button int32
	Touch  int64
	Action ActionType
	Mods   int32
}

// Event ...
type Event struct {
	keyboard KeyboardEvent
	mouse    MouseEvent
}

// InvalidKeyboardEvent ...
func InvalidKeyboardEvent() (out KeyboardEvent) {
	out.Key = -1
	out.Scancode = -1
	out.Action = 0
	out.Mods = 0
	return
}

// InvalidMouseEvent ...
func InvalidMouseEvent() (out MouseEvent) {
	out.X = math.NaN()
	out.Y = math.NaN()
	out.Button = -1
	out.Action = 0
	out.Mods = 0
	return
}

// InvalidEvent ...
func InvalidEvent() (out Event) {
	out.keyboard = InvalidKeyboardEvent()
	out.mouse = InvalidMouseEvent()
	return
}
