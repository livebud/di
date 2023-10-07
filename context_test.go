package di_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/di"
	"github.com/matryer/is"
)

func TestContext(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader[*Env](in, loadEnv)
	env, err := di.Load[*Env](in)
	is.NoErr(err)
	is.True(env != nil)
	env, err = di.Load[*Env](in)
	is.NoErr(err)
	is.True(env != nil)
	ctx := di.WithInjector(context.Background(), in)
	env, err = di.LoadFrom[*Env](ctx)
	is.NoErr(err)
	is.True(env != nil)
}

type Stack struct {
	fns []func(http.Handler) http.Handler
}

func (s *Stack) Append(fn func(http.Handler) http.Handler) {
	s.fns = append(s.fns, fn)
}

func (s *Stack) Compose(bottom http.Handler) http.Handler {
	stack := func(next http.Handler) http.Handler {
		for _, m := range s.fns {
			next = m(next)
		}
		return next
	}
	return stack(bottom)
}

func (s *Stack) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Compose(http.NotFoundHandler()).ServeHTTP(w, r)
}

func TestMiddleware(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader[*Env](in, loadEnv)
	err := di.Loader[*Stack](in, func(in di.Injector) (*Stack, error) {
		return &Stack{}, nil
	})
	is.NoErr(err)
	called := 0
	// Attach the injector
	di.Append[*Stack](in, func(in di.Injector, s *Stack) error {
		s.Append(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				in := di.Clone(in)
				r = r.WithContext(di.WithInjector(r.Context(), in))
				next.ServeHTTP(w, r)
			})
		})
		called++
		return nil
	})
	stack, err := di.Load[*Stack](in)
	is.NoErr(err)
	h := stack.Compose(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env, err := di.LoadFrom[*Env](r.Context())
		is.NoErr(err)
		is.True(env != nil)
		called++
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	is.Equal(called, 2)
}
