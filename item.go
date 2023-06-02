package cache

import "time"

type Item[T comparable] struct {
	init  func() (T, error)
	timer *time.Timer
	ttl   time.Duration
	value T
	zero  T
}

func NewItem[T comparable](ttl time.Duration, init func() (T, error)) Item[T] {
	i := Item[T]{
		init: init,
		ttl:  ttl,
	}

	return i
}

func (i *Item[T]) Clear() {
	if i.timer != nil {
		i.timer.Stop()
		i.timer = nil
	}

	i.value = i.zero
}

func (i *Item[T]) Value() (T, error) {
	if i.value != i.zero {
		return i.value, nil
	}

	v, err := i.init()
	if err != nil {
		return i.zero, err
	}

	i.value = v
	i.timer = time.AfterFunc(i.ttl, i.Clear)

	return v, nil
}
