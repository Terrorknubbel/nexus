package main

import (
	"context"
	"nexus/internal/config"
	"nexus/internal/process"
	"nexus/pkg/models"
)

type App struct {
	ctx         context.Context
	procManager *process.Manager
}

func NewApp() *App {
	cfg := config.NewDefault()
	procMgr := process.NewManager(cfg)
	return &App{procManager: procMgr}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAllProcesses() []models.Process {
	procs, err := a.procManager.CollectAll()
	if err != nil {
		panic(err)
	}
	return procs
}
