package main

import (
	"fmt"
	"context"
)

type Golang struct {}

func (m *Golang) Base(ctx context.Context, version string) (*Container) {
	cache := dag.CacheVolume("gomodcache")
	image := fmt.Sprintf("golang:%s", version)
	return dag.Container().
	From(image).
	WithMountedCache("/gomodcache", cache).
	WithEnvironmentVariable("GOMODCACHE", "/gomodcache")
}

func (c *Container) GoBuild(ctx context.Context, args []string) (*Container, error) {
	command := append([]string{"go", "build"}, args)
	return c.WithExec(command).Sync(ctx)
}

func (c *Container) GoTest(ctx context.Context, ctr *Container, args []string) (*Container, error) {
	command := append([]string{"go", "test"}, args)
	return c.WithExec(command).Sync(ctx)
}

func (d *Directory) GoLint(ctx context.Context) (string, error) {
	return dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", d).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stdout(ctx)
}

