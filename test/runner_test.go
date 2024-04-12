package test

import (
	"testing"

	"github.com/rishtastic/frontrunner/runner"
)

func TestFrontRunner_First(t *testing.T) {
	t.Run("Should return error when no tasks are provided", func(t *testing.T) {
		r := runner.NewRunner[int]()
		_, err := r.First()
		if err == nil {
			t.Fatalf("error should not be nil")
		}
	})

	t.Run("Should return first result when execution is successful", func(t *testing.T) {
		exp := 0
		task := func() int {
			return exp
		}
		r := runner.NewRunner(task)
		act, err := r.First()
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}

		if act != exp {
			t.Fatalf("res should be %v, was %v", exp, act)
		}
	})
}

func TestFrontRunner_FirstK(t *testing.T) {
	k := 3

	t.Run("Should return error when no tasks are provided", func(t *testing.T) {
		r := runner.NewRunner[int]()
		_, err := r.FirstK(k)
		if err == nil {
			t.Fatalf("error should not be nil")
		}
	})

	t.Run("Should return error when insufficient tasks are provided", func(t *testing.T) {
		task := func() bool {
			return true
		}
		r := runner.NewRunner(task)
		_, err := r.FirstK(k)
		if err == nil {
			t.Fatalf("error should not be nil")
		}
	})

	t.Run("Should return exactly k items when tasks run successfully", func(t *testing.T) {
		tasks := []func() int{}
		for i := range k + 2 {
			exp := i
			tasks = append(tasks, func() int { return exp })
		}
		print(len(tasks))
		r := runner.NewRunner(tasks...)
		res, err := r.FirstK(k)
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}
		l := len(res)
		if l != k {
			t.Fatalf("should have exactly %d items, got %d", k, l)
		}
	})
}
