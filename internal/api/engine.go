package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chankei613/ai-task-router/internal/db"
)

var tierRank = map[string]int{
	db.TierEconomy:  0,
	db.TierStandard: 1,
	db.TierPremium:  2,
}

func hasAllCapabilities(model db.ModelSpec, required []string) bool {
	have := make(map[string]bool, len(model.Capabilities))
	for _, c := range model.Capabilities {
		have[c] = true
	}
	for _, r := range required {
		if !have[r] {
			return false
		}
	}
	return true
}

func costScore(m db.ModelSpec) float64 {
	return m.InputPricePer1M + m.OutputPricePer1M
}

// decideRoute はcandidatesの中からrequiredCapabilitiesを全て満たし、minQualityTier以上の
// ものだけに絞り込み、その中でコスト（入力+出力単価合計）最安のモデルを選ぶ。
// 該当がなければ chosen は nil になり、reasoning に除外理由の要約が入る。
func decideRoute(candidates []db.ModelSpec, requiredCapabilities []string, minQualityTier string) (*db.ModelSpec, string) {
	minRank, ok := tierRank[minQualityTier]
	if !ok {
		minRank = tierRank[db.TierEconomy]
	}

	var qualifying []db.ModelSpec
	excludedByCapability := 0
	excludedByTier := 0
	excludedByDisabled := 0

	for _, m := range candidates {
		if !m.Enabled {
			excludedByDisabled++
			continue
		}
		if !hasAllCapabilities(m, requiredCapabilities) {
			excludedByCapability++
			continue
		}
		if tierRank[m.QualityTier] < minRank {
			excludedByTier++
			continue
		}
		qualifying = append(qualifying, m)
	}

	if len(qualifying) == 0 {
		reasons := []string{}
		if excludedByDisabled > 0 {
			reasons = append(reasons, fmt.Sprintf("%d disabled", excludedByDisabled))
		}
		if excludedByCapability > 0 {
			reasons = append(reasons, fmt.Sprintf("%d missing a required capability", excludedByCapability))
		}
		if excludedByTier > 0 {
			reasons = append(reasons, fmt.Sprintf("%d below minimum quality tier %q", excludedByTier, minQualityTier))
		}
		reasoning := "no model satisfies the constraints"
		if len(reasons) > 0 {
			reasoning += " (" + strings.Join(reasons, ", ") + ")"
		}
		return nil, reasoning
	}

	sort.Slice(qualifying, func(i, j int) bool {
		return costScore(qualifying[i]) < costScore(qualifying[j])
	})
	chosen := qualifying[0]

	reasoning := fmt.Sprintf(
		"chosen %s/%s: cheapest of %d qualifying model(s) (cost score %.2f per 1M tokens)",
		chosen.Provider, chosen.ModelID, len(qualifying), costScore(chosen),
	)
	return &chosen, reasoning
}
