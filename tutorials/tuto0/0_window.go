package tuto0

import (
	"skyfx/app"
	"skyfx/gfx"
)

func Run() {
	if app.Run("skyfx: Window Tutorial", 1024, 768, 60, "assets/", gfx.GL_VERSION_DS_3_2, 2) {
		app.Destroy()
	}
}
