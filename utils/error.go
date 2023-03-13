package utils

func PanicError(msg string, err error) {
	if err != nil {
		panic(msg)
	}
}
