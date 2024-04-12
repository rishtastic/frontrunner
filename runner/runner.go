package runner

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Runner is a struct that handles running multiple tasks concurrently.
type Runner[T any] struct {
	tasks []func() T
	mu    sync.Mutex
}

// NewRunner creates a variable of type Runner,
// accepts variadic tasks to be executed,
// return returns pointer to Runner variable.
func NewRunner[T any](tasks ...func() T) *Runner[T] {
	return &Runner[T]{
		tasks: tasks,
		mu:    sync.Mutex{},
	}
}

// Add adds new tasks to the runner,
// accepts variadic tasks to be added.
func (r *Runner[T]) Add(t ...func() T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, t...)
}

// Runs the tasks concurrently and returns result of the first one that completes,
// returns error if no tasks are provided.
func (r *Runner[T]) First() (T, error) {
	var res T
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.tasks) <= 0 {
		return res, errors.New("runner: no tasks provided")
	}
	rch := make(chan []T)
	go raceTasks(r.tasks, 1, rch)
	resArr := <-rch
	return resArr[0], nil
}

// Runs the tasks concurrently and returns results of the first k completed tasks,
// returns error if k <= 0,
// returns error if not enough tasks are provided.
func (r *Runner[T]) FirstK(k int) ([]T, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if k <= 0 {
		return nil, fmt.Errorf("runner: k should be greater than 0")
	}
	if len(r.tasks) < k {
		return nil, fmt.Errorf("runner: not enough tasks: need %d, have %d", k, len(r.tasks))
	}
	res := make(chan []T)
	go raceTasks(r.tasks, k, res)
	return <-res, nil
}

// Runs the tasks concurrently and returns result of the first one that completes,
// also returns the result of the operation, true if run successfully, false if timed out or error,
// the operaton is timed out if atleast 1 task doesn't complete within the duration (d),
// returns error if no tasks are provided.
func (r *Runner[T]) FirstWithTimeout(d time.Duration) (T, bool, error) {
	var res T
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.tasks) <= 0 {
		return res, false, errors.New("runner: no tasks provided")
	}

	timed := make(chan bool, 1)
	go timeout(timed, d)

	rch := make(chan []T, 1)
	go raceTasks(r.tasks, 1, rch)

	select {
	case <-timed:
		return res, false, nil
	case resArr := <-rch:
		return resArr[0], true, nil
	}
}

// Runs the tasks concurrently and returns results of the first k completed tasks,
// also returns the result of the operation, true if run successfully, false if timed out or error,
// the operaton is timed out if atleast k tasks don't complete within the duration (d),
// returns error if k <= 0,
// returns error if not enough tasks are provided.
func (r *Runner[T]) FirstKWithTimeout(k int, d time.Duration) ([]T, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if k <= 0 {
		return nil, false, fmt.Errorf("runner: k should be greater than 0")
	}
	if len(r.tasks) < k {
		return nil, false, fmt.Errorf("runner: not enough tasks: need %d, have %d", k, len(r.tasks))
	}

	timed := make(chan bool, 1)
	go timeout(timed, d)

	rch := make(chan []T)
	go raceTasks(r.tasks, k, rch)

	select {
	case <-timed:
		return nil, false, nil
	case res := <-rch:
		return res, true, nil
	}
}
