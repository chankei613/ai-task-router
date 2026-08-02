// Package db はAI Task RouterのGORMモデルとSQLite初期化を提供する。
package db

import "time"

// QualityTier はモデルの品質階層。文字列の大小比較ではなく tierRank() で順序評価する。
const (
	TierEconomy  = "economy"
	TierStandard = "standard"
	TierPremium  = "premium"
)

// ModelSpec は登録済みAIモデル1件。本製品はこのモデルを実際には呼び出さない — あくまで
// 「どれを使うべきか」を単価とcapabilityで判定するためのカタログ。
type ModelSpec struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Provider string `json:"provider"` // 例: anthropic, openai, google
	ModelID  string `json:"model_id"` // 例: claude-sonnet-5, gpt-4o, gemini-2.5-pro

	QualityTier      string   `json:"quality_tier"` // economy/standard/premium
	InputPricePer1M  float64  `json:"input_price_per_1m"`
	OutputPricePer1M float64  `json:"output_price_per_1m"`
	Capabilities     []string `gorm:"serializer:json" json:"capabilities"`
	Enabled          bool     `json:"enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutingLog は1回のルーティング判断の記録（追記専用、監査・集計用）。
type RoutingLog struct {
	ID     string `gorm:"primaryKey" json:"id"`
	Source string `json:"source"` // 呼び出し元の自由記述

	TaskType             string   `json:"task_type"`
	RequiredCapabilities []string `gorm:"serializer:json" json:"required_capabilities"`
	MinQualityTier       string   `json:"min_quality_tier"`

	ChosenModelID string `json:"chosen_model_id"` // 該当なしなら空文字
	Reasoning     string `json:"reasoning"`

	RequestedAt time.Time `json:"requested_at"`
}

// AgentKey — CRUD/ルーティングAPIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

// DefaultModels は初回起動時にシードする主要モデルの一覧（2026年中頃時点の目安単価）。
// ユーザーは実際の契約に合わせて編集・追加できる。
func DefaultModels() []ModelSpec {
	now := time.Now()
	mk := func(provider, modelID, tier string, inPrice, outPrice float64, caps []string) ModelSpec {
		return ModelSpec{
			ID: provider + ":" + modelID, Provider: provider, ModelID: modelID,
			QualityTier: tier, InputPricePer1M: inPrice, OutputPricePer1M: outPrice,
			Capabilities: caps, Enabled: true, CreatedAt: now, UpdatedAt: now,
		}
	}
	return []ModelSpec{
		mk("anthropic", "claude-opus-5", TierPremium, 15, 75, []string{"vision", "function_calling", "long_context"}),
		mk("anthropic", "claude-sonnet-5", TierStandard, 3, 15, []string{"vision", "function_calling", "long_context"}),
		mk("anthropic", "claude-haiku-4-5", TierEconomy, 0.8, 4, []string{"vision", "function_calling"}),
		mk("openai", "gpt-4o", TierStandard, 2.5, 10, []string{"vision", "function_calling", "long_context"}),
		mk("openai", "gpt-4o-mini", TierEconomy, 0.15, 0.6, []string{"vision", "function_calling"}),
		mk("google", "gemini-2.5-pro", TierStandard, 1.25, 5, []string{"vision", "function_calling", "long_context"}),
		mk("google", "gemini-2.5-flash", TierEconomy, 0.075, 0.3, []string{"vision", "function_calling"}),
	}
}
