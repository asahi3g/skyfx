package main

import (
	"runtime"
	"skyfx/skycad"
	"skyfx/tutorials/tuto0"
	"skyfx/tutorials/tuto1"
	"skyfx/tutorials/tuto2"
	"skyfx/tutorials/tuto3"
	"skyfx/tutorials/tuto4"
	"skyfx/tutorials/tuto5"
	"skyfx/tutorials/tuto6"
	"skyfx/tutorials/tuto7"
)

func init() {
	runtime.LockOSThread()
}

var tuto = 6

func main() {
	switch tuto {
	case 0:
		tuto0.Run()
	case 1:
		tuto1.Run()
	case 2:
		tuto2.Run()
	case 3:
		tuto3.Run()
	case 4:
		tuto4.Run()
	case 5:
		tuto5.Run()
	case 6:
		tuto6.Run()
	case 7:
		tuto7.Run()
	case 8:
		skycad.Run()
	}
}
