package math

import "skyfx/utils"

// TODO build flag panicMath.

func panicIf(value bool, message string) {
	utils.PanicIf(value, message)
}
