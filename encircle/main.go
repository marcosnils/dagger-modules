package main

import (
	"context"
)

const CONFIG = "./.circleci/config.yml"

type Encircle struct{}

func (m *Encircle) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}

func (d *Directory) EncircleJob(ctx context.Context, job string) error {
	cfg, executor, err := setup(ctx)
	if err != nil {
		return err
	}

	err = executor.ExecuteJob(job, cfg.Jobs[job])
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) EncircleWorkflow(ctx context.Context, workflow string) error {
	cfg, executor, err := setup(ctx)
	if err != nil {
		return err
	}

	err = executor.ExecuteWorkflow(workflow, cfg.Workflows[workflow], cfg.Jobs)
	if err != nil {
		return err
	}

	return nil
}

func setup(ctx context.Context) (*Config, *Executor, error) {
	executor := NewExecutor(ctx)
	cfg, err := readConfig(CONFIG)

	return cfg, executor, err
}