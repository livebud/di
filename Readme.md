# DI

[![Go Reference](https://pkg.go.dev/badge/github.com/livebud/di.svg)](https://pkg.go.dev/github.com/livebud/di)

Dependency injection using Generics.

## Features

- Provide dependencies in a natural and type-safe way
- Register dependencies with other dependencies (e.g. register a controller with the router)
- Swap out dependencies during testing
- Middleware support
- Unmarshal dependencies into a struct

## Install

```
go get github.com/livebud/di
```

## Example

In the following example, we load the logger which loads the environment:

```go
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
  fmt.Println(log.env.DatabaseURL)
}
```

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT
