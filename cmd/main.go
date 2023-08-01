package main

import (
	"fmt"

	"github.com/livebud/di"
)

type Env struct {
	DatabaseURL string
}

func provideEnv(in di.Injector) (*Env, error) {
	return &Env{
		DatabaseURL: "postgres://localhost:5432/db",
	}, nil
}

type Log struct {
	env *Env
}

func provideLog(in di.Injector) (*Log, error) {
	env, err := di.Load[*Env](in)
	if err != nil {
		return nil, err
	}
	return &Log{env}, nil
}

func main() {
	in := di.New()
	di.Provide(in, provideEnv)
	di.Provide(in, provideLog)
	log, err := di.Load[*Log](in)
	if err != nil {
		panic(err)
	}
	fmt.Println(log.env.DatabaseURL)
}
