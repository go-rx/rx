package rx

func buildOptions[T any, O ~func(*T)](options []O) *T {
	opts := new(T)
	for _, option := range options {
		option(opts)
	}
	return opts
}
