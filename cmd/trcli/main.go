// cmd/trcli はCI/CDパイプラインに組み込むためのコマンドラインクライアント。
// タスク種別・必要capability・品質下限をフラグで渡すと、選定されたmodel_idを標準出力へ、
// 理由を標準エラーへ出力する。該当モデルがなければexit 1、実行時エラーならexit 2。
//
//	go run ./cmd/trcli -addr http://127.0.0.1:8427 -key $TR_API_KEY -task summarization -require vision -quality economy
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type routeInput struct {
	Source               string   `json:"source"`
	TaskType             string   `json:"task_type"`
	RequiredCapabilities []string `json:"required_capabilities"`
	MinQualityTier       string   `json:"min_quality_tier"`
}

type modelSpec struct {
	Provider string `json:"provider"`
	ModelID  string `json:"model_id"`
}

type routingLog struct {
	ChosenModelID string `json:"chosen_model_id"`
	Reasoning     string `json:"reasoning"`
}

type routeResult struct {
	Log         routingLog `json:"log"`
	ChosenModel *modelSpec `json:"chosen_model,omitempty"`
}

func main() {
	addr := flag.String("addr", envOr("TR_ADDR", "http://127.0.0.1:8427"), "AI Task Router API base URL")
	key := flag.String("key", os.Getenv("TR_API_KEY"), "APIキー（またはTR_API_KEY環境変数）")
	task := flag.String("task", "", "タスク種別（自由記述）")
	require := flag.String("require", "", "必要なcapability（カンマ区切り、例: vision,function_calling）")
	quality := flag.String("quality", "economy", "品質下限（economy/standard/premium）")
	source := flag.String("source", "cli", "呼び出し元の識別子")
	flag.Parse()

	if *key == "" {
		fmt.Fprintln(os.Stderr, "error: -key or TR_API_KEY is required")
		os.Exit(2)
	}

	var caps []string
	if *require != "" {
		caps = strings.Split(*require, ",")
	}

	in := routeInput{
		Source:               *source,
		TaskType:             *task,
		RequiredCapabilities: caps,
		MinQualityTier:       *quality,
	}
	body, err := json.Marshal(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building request: %s\n", err)
		os.Exit(2)
	}

	req, err := http.NewRequest(http.MethodPost, *addr+"/api/v1/route", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building request: %s\n", err)
		os.Exit(2)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling server: %s\n", err)
		os.Exit(2)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %s\n", err)
		os.Exit(2)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "server returned %d: %s\n", resp.StatusCode, string(respBody))
		os.Exit(2)
	}

	var result routeResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing server response: %s\n", err)
		os.Exit(2)
	}

	fmt.Fprintln(os.Stderr, result.Log.Reasoning)

	if result.ChosenModel == nil {
		os.Exit(1)
	}
	fmt.Println(result.ChosenModel.Provider + "/" + result.ChosenModel.ModelID)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
