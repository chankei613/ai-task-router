package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/ai-task-router/internal/api"
	"github.com/chankei613/ai-task-router/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const apiAddr = "127.0.0.1:8427"

// App はWailsのバインディング。実処理は internal/api.Server が持っている。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "ai-task-router.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "AI Task Router ready (api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) GetAPIURL() string {
	return "http://" + apiAddr
}

func (a *App) ListModels() ([]db.ModelSpec, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListModels()
}

func (a *App) CreateModel(provider, modelID, qualityTier string, inputPrice, outputPrice float64, capabilities []string) (db.ModelSpec, error) {
	if !a.ready {
		return db.ModelSpec{}, errNotReady
	}
	return a.server.UpsertModel(api.UpsertModelInput{
		Provider: provider, ModelID: modelID, QualityTier: qualityTier,
		InputPricePer1M: inputPrice, OutputPricePer1M: outputPrice, Capabilities: capabilities,
	})
}

func (a *App) SetModelEnabled(id string, enabled bool) (db.ModelSpec, error) {
	if !a.ready {
		return db.ModelSpec{}, errNotReady
	}
	return a.server.SetModelEnabled(id, enabled)
}

func (a *App) DeleteModel(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeleteModel(id)
}

func (a *App) Route(source, taskType string, requiredCapabilities []string, minQualityTier string) (api.RouteResult, error) {
	if !a.ready {
		return api.RouteResult{}, errNotReady
	}
	return a.server.Route(api.RouteInput{
		Source: source, TaskType: taskType,
		RequiredCapabilities: requiredCapabilities, MinQualityTier: minQualityTier,
	})
}

func (a *App) ListRoutingLogs() ([]db.RoutingLog, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListRoutingLogs()
}

func (a *App) GetRoutingLog(id string) (api.RouteResult, error) {
	if !a.ready {
		return api.RouteResult{}, errNotReady
	}
	return a.server.GetRoutingLog(id)
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".ai-task-router")
}
