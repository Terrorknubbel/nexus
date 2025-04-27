package main

import (
	"context"
	"nexus/internal/process"
	"nexus/pkg/models"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAllProcesses() []models.Process {
	return process.Collect()
}
