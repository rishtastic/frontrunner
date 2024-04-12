package runner

import (
	"time"
)

func raceTasks[T any](tasks []func() T, k int, rch chan<- []T) {
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
	rch <- res
}

func taskChannel[T any](ch chan<- T, t func() T) {
	defer func() { recover() }()
	r := t()
	ch <- r
}

func timeout(ch chan<- bool, d time.Duration) {
	time.AfterFunc(d, func() {
		ch <- true
	})
}
