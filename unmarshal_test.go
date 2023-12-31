package di_test

import (
	"testing"

	"github.com/livebud/di"
	"github.com/matryer/is"
)

type Context struct {
	Env *Env
	Log *Log
}

func TestUnmarshal(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader(in, loadEnv)
	di.Loader(in, loadLog)
	var ctx Context
	err := di.Unmarshal(in, &ctx)
	is.NoErr(err)
	is.Equal(ctx.Env.name, "production")
	is.Equal(ctx.Log.lvl, "info")
}

func TestUnmarshalPointer(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader(in, loadEnv)
	di.Loader(in, loadLog)
	var ctx Context
	err := di.Unmarshal(in, &ctx)
	is.NoErr(err)
	is.Equal(ctx.Env.name, "production")
	is.Equal(ctx.Log.lvl, "info")
}
