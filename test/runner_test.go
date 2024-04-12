package test

import (
	"testing"
	"time"

	"github.com/rishtastic/frontrunner/runner"
)

func TestFrontRunner_First(t *testing.T) {
	t.Run("Should return error when no tasks are provided", func(t *testing.T) {
		r := runner.NewRunner[int]()
		_, err := r.First()
		if err == nil {
			t.Fatal("error should not be nil")
		}
	})

	t.Run("Should return first result when execution is successful", func(t *testing.T) {
		exp := 1
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
			t.Fatal("error should not be nil")
		}
	})

	t.Run("Should return error when insufficient tasks are provided", func(t *testing.T) {
		task := func() bool {
			return true
		}
		r := runner.NewRunner(task)
		_, err := r.FirstK(k)
		if err == nil {
			t.Fatal("error should not be nil")
		}
	})

	t.Run("Should return exactly k items when tasks run successfully", func(t *testing.T) {
		r := runner.NewRunner[int]()
		for i := range k + 2 {
			exp := i
			r.Add(func() int { return exp })
		}
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

func TestFrontRunner_FirstWithTimout(t *testing.T) {
	timeOut := time.Second
	t.Run("Should return error when no tasks are provided", func(t *testing.T) {
		r := runner.NewRunner[int]()
		_, ok, err := r.FirstWithTimeout(timeOut)
		if err == nil {
			t.Fatal("error should not be nil")
		}
		if ok {
			t.Fatal("completion should be false")
		}
	})

	t.Run("Should return completion false if timed out", func(t *testing.T) {
		r := runner.NewRunner[int]()
		ch := make(chan int)
		defer func() { ch <- 1 }()
		r.Add(func() int {
			return <-ch
		})
		_, ok, err := r.FirstWithTimeout(timeOut)
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}
		if ok {
			t.Fatal("completion should be false")
		}
	})

	t.Run("Should return first result and completed when execution is successful", func(t *testing.T) {
		exp := 1
		task := func() int {
			return exp
		}
		r := runner.NewRunner(task)
		act, ok, err := r.FirstWithTimeout(timeOut)
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}
		if !ok {
			t.Fatalf("completion should be true, was %v", ok)
		}

		if act != exp {
			t.Fatalf("res should be %v, was %v", exp, act)
		}
	})
}

func TestFrontRunner_FirstKWithTimeout(t *testing.T) {
	k := 3
	timeOut := time.Second

	t.Run("Should return error when no tasks are provided", func(t *testing.T) {
		r := runner.NewRunner[int]()
		_, ok, err := r.FirstKWithTimeout(k, timeOut)
		if err == nil {
			t.Fatal("error should not be nil")
		}
		if ok {
			t.Fatal("completion should be false")
		}
	})

	t.Run("Should return error when insufficient tasks are provided", func(t *testing.T) {
		task := func() bool {
			return true
		}
		r := runner.NewRunner(task)
		_, ok, err := r.FirstKWithTimeout(k, timeOut)
		if err == nil {
			t.Fatal("error should not be nil")
		}
		if ok {
			t.Fatal("completion should be false")
		}
	})

	t.Run("Should return completion false if timed out", func(t *testing.T) {
		r := runner.NewRunner[int]()
		ch := make(chan int)
		defer func() { ch <- 1 }()
		r.Add(func() int {
			return <-ch
		})
		for i := range k - 1 {
			exp := i
			r.Add(func() int { return exp })
		}
		_, ok, err := r.FirstKWithTimeout(k, timeOut)
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}
		if ok {
			t.Fatal("completion should be false")
		}
	})

	t.Run("Should return exactly k items when k tasks run successfully", func(t *testing.T) {
		r := runner.NewRunner[int]()
		ch := make(chan int)
		defer func() { ch <- 1 }()
		r.Add(func() int {
			return <-ch
		})
		for i := range k {
			exp := i
			r.Add(func() int { return exp })
		}
		res, ok, err := r.FirstKWithTimeout(k, timeOut)
		if err != nil {
			t.Fatalf("error should be nil, was %v", err)
		}
		if !ok {
			t.Fatalf("completion should be true, was %v", ok)
		}
		l := len(res)
		if l != k {
			t.Fatalf("should have exactly %d items, got %d", k, l)
		}
	})
}
