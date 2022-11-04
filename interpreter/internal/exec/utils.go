package exec

func ignoreError[T any](fn func() (T, error)) T {
	x, _ := fn()
	return x
}
