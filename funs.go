package husky

import (
	"context"
	"runtime/debug"
	"slices"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func ToPoint[P any](p P) *P {
	return &p
}

func If[R any](logic bool, r1, r2 R) R {
	if logic {
		return r1
	} else {
		return r2
	}
}

func NilDefault[P any](data *P, defaultVal P) P {
	if data == nil {
		return defaultVal
	} else {
		return *data
	}
}

func Filter[S ~[]E, E any](datas S, f func(item E) bool) S {
	result := make(S, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		if r {
			result = append(result, data)
		}
	}
	return result
}

func Match[T comparable](v T, matchers ...T) bool {
	for i := range matchers {
		if v == matchers[i] {
			return true
		}
	}
	return false
}

func Map[P any, R any](datas []P, f func(item P) R) []R {
	result := make([]R, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		result = append(result, r)
	}
	return result
}

func Go(call func()) {
	go func() {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			Error(context.WithValue(context.Background(), ContextKeyTraceId, ""), err)
			debug.PrintStack()
		}()
		call()
	}()
}

func Reverse[S ~[]E, E any](s S) {
	slices.Reverse(s)
}

func NewContext() context.Context {
	return context.WithValue(context.Background(), ContextKeyTraceId, strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
}

func NewContextWithParent(parent context.Context) context.Context {
	return context.WithValue(parent, ContextKeyTraceId, strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
}
