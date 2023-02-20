package repository

import (
	"sync"
)

type message struct {
	sync.Once
	done chan struct{}
	feed <-chan string
	str  string
}

func (m *message) aggregate() {
	fn := func() {
		m.done = make(chan struct{})
		for {
			s, ok := <-m.feed
			m.str += s
			if ok {
				continue
			}

			m.done <- struct{}{}
			return
		}
	}

	m.Once.Do(func() { go fn() })
}
