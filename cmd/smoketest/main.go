// cmd/smoketest はAI Task RouterのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → デフォルトモデルのシード確認 → 制約なしルーティング →
// capability制約ありルーティング → 該当なしケース → モデル無効化後の再ルーティング →
// カスタムモデル追加、の一連が通しで動くことを確認する。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/chankei613/ai-task-router/internal/api"
	"github.com/chankei613/ai-task-router/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	authed := func(method, path string, body []byte) *http.Response {
		req, _ := http.NewRequest(method, srv.URL+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+issued.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	// default models seeded on first init
	resp = authed(http.MethodGet, "/api/v1/models", nil)
	var models []db.ModelSpec
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(models) != 7 {
		log.Fatalf("FAIL: expected 7 seeded default models, got %d", len(models))
	}
	fmt.Println("PASS: default models seeded on first init")

	// route with no capability constraint, economy tier — cheapest of all 7 should win
	route1Body, _ := json.Marshal(api.RouteInput{Source: "smoketest", TaskType: "general", MinQualityTier: db.TierEconomy})
	resp = authed(http.MethodPost, "/api/v1/route", route1Body)
	var route1 api.RouteResult
	if err := json.NewDecoder(resp.Body).Decode(&route1); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if route1.ChosenModel == nil || route1.ChosenModel.ID != "google:gemini-2.5-flash" {
		log.Fatalf("FAIL: expected cheapest overall model google:gemini-2.5-flash, got %+v", route1.ChosenModel)
	}
	fmt.Println("PASS: unconstrained route picks the cheapest model overall (gemini-2.5-flash)")

	// route requiring long_context + standard tier — cheapest qualifying should be gemini-2.5-pro
	route2Body, _ := json.Marshal(api.RouteInput{
		Source: "smoketest", TaskType: "long-doc-summary",
		RequiredCapabilities: []string{"long_context"}, MinQualityTier: db.TierStandard,
	})
	resp = authed(http.MethodPost, "/api/v1/route", route2Body)
	var route2 api.RouteResult
	if err := json.NewDecoder(resp.Body).Decode(&route2); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if route2.ChosenModel == nil || route2.ChosenModel.ID != "google:gemini-2.5-pro" {
		log.Fatalf("FAIL: expected google:gemini-2.5-pro for long_context+standard, got %+v", route2.ChosenModel)
	}
	fmt.Println("PASS: capability+tier constraint picks the cheapest qualifying model (gemini-2.5-pro)")

	// route requiring a capability nothing has — must return nil with a clear reason, not error
	route3Body, _ := json.Marshal(api.RouteInput{
		Source: "smoketest", TaskType: "impossible",
		RequiredCapabilities: []string{"quantum_ml"}, MinQualityTier: db.TierEconomy,
	})
	resp = authed(http.MethodPost, "/api/v1/route", route3Body)
	var route3 api.RouteResult
	if err := json.NewDecoder(resp.Body).Decode(&route3); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if route3.ChosenModel != nil || route3.Log.Reasoning == "" {
		log.Fatalf("FAIL: expected no match with a reasoning message, got %+v", route3)
	}
	fmt.Println("PASS: impossible constraint returns no match with an explicit reason")

	// disable the previous winner, re-run the unconstrained route — next cheapest must win
	disableBody, _ := json.Marshal(map[string]bool{"enabled": false})
	_ = authed(http.MethodPatch, "/api/v1/models/google:gemini-2.5-flash/enabled", disableBody).Body.Close()

	resp = authed(http.MethodPost, "/api/v1/route", route1Body)
	var route4 api.RouteResult
	if err := json.NewDecoder(resp.Body).Decode(&route4); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if route4.ChosenModel == nil || route4.ChosenModel.ID != "openai:gpt-4o-mini" {
		log.Fatalf("FAIL: expected openai:gpt-4o-mini after disabling the cheapest model, got %+v", route4.ChosenModel)
	}
	fmt.Println("PASS: disabling a model excludes it from routing, next cheapest wins")

	// routing history should have 4 entries
	resp = authed(http.MethodGet, "/api/v1/routes", nil)
	var logs []db.RoutingLog
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(logs) != 4 {
		log.Fatalf("FAIL: expected 4 routing logs in history, got %d", len(logs))
	}
	fmt.Println("PASS: routing history contains all 4 decisions")

	// add a custom model via upsert
	customBody, _ := json.Marshal(api.UpsertModelInput{
		Provider: "local", ModelID: "llama-70b", QualityTier: db.TierEconomy,
		InputPricePer1M: 0.1, OutputPricePer1M: 0.1, Capabilities: []string{"function_calling"},
	})
	resp = authed(http.MethodPost, "/api/v1/models", customBody)
	var custom db.ModelSpec
	if err := json.NewDecoder(resp.Body).Decode(&custom); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if custom.ID != "local:llama-70b" {
		log.Fatalf("FAIL: expected custom model id local:llama-70b, got %s", custom.ID)
	}
	fmt.Println("PASS: custom model added via upsert")

	fmt.Println("SMOKE TEST OK")
}
