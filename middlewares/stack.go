package middlewares

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Stack struct {
	mws []Middleware
}

func NewStack() *Stack { return &Stack{} }

func (s *Stack) Use(mw Middleware) *Stack {
	s.mws = append(s.mws, mw)
	return s
}

func (s *Stack) Then(h http.Handler) http.Handler {
	for i := len(s.mws) - 1; i >= 0; i-- {
		h = s.mws[i](h)
	}
	return h
}
