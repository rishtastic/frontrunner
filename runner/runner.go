package runner

import (
	"errors"
	"fmt"
	"sync"
)

type Runner[T any] struct {
	tasks []func() T
	mu    sync.Mutex
}

func NewRunner[T any](tasks ...func() T) *Runner[T] {
	return &Runner[T]{
		tasks: tasks,
		mu:    sync.Mutex{},
	}
}

func (r *Runner[T]) Add(t ...func() T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, t...)
}

func (r *Runner[T]) First() (T, error) {
	var res T
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.tasks) <= 0 {
		return res, errors.New("runner: no tasks provided")
	}
	resArr := raceTasks(r.tasks, 1)
	return resArr[0], nil
}

func (r *Runner[T]) FirstK(k int) ([]T, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if k <= 0 {
		return nil, fmt.Errorf("k should be greater than 0")
	}
	if len(r.tasks) < k {
		return nil, fmt.Errorf("runner: not enough tasks: need %d, have %d", k, len(r.tasks))
	}
	res := raceTasks(r.tasks, k)
	return res, nil
}

func raceTasks[T any](tasks []func() T, k int) []T {
	ch := make(chan T, len(tasks))
	for _, task := range tasks {
		t := task
		go taskChannel(ch, t)
	}
	res := []T{}
	for range k {
		r := <-ch
		res = append(res, r)
	}
	return res
}

func taskChannel[T any](ch chan<- T, t func() T) {
	r := t()
	ch <- r
}
