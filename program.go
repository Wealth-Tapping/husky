package husky

import (
	"os"
	"os/signal"
	"syscall"
)

type Program struct {
	chs []chan os.Signal
}

func (s *Program) Add(chs ...chan os.Signal) {
	s.chs = append(s.chs, chs...)
}

func (s *Program) Waiting() os.Signal {
	sign := s.waiting()
	for index := range s.chs {
		ch := s.chs[index]
		ch <- sign
	}
	for index := range s.chs {
		ch := s.chs[index]
		<-ch
	}
	return sign
}
func (s *Program) waiting() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	sign := <-ch
	return sign
}

func NewProgram() *Program {
	return &Program{chs: []chan os.Signal{}}
}
