package utils

func PanicIf(condition bool, message string) {
	if condition {
		panic(message)
	}
}

func PanicIfNot(condition bool, message string) {
	PanicIf(!condition, message)
}
