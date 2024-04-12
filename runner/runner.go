package runner

type Runner[T any] struct {
	tasks []func() T
	ch    chan T
}

func NewRunner[T any](tasks ...func() T) *Runner[T] {
	return &Runner[T]{
		tasks: tasks,
		ch:    make(chan T, len(tasks)),
	}
}
