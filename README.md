# AI Task Router

「どのAIを使うか」の意思決定を自動化する — comet-taskAI ロードマップ Product J。

タスク種別・必要capability・品質下限を渡すと、登録済みモデルの中からコスト最適な候補を
決定論的ルールで選び、理由つきで返す。実際にそのモデルを呼び出すのは外部システムの仕事。

本製品自身はどのAIプロバイダーも呼び出さない。I. AI Output Validatorが「AIを呼ばずに出力を
検証する」と決めたのと同じ考え方で、本製品は「AIを呼ばずにどのモデルを使うべきか決定する」
ことに徹する。

## 現在のステータス: v0.1.0 リリース済み

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・ルーティングエンジン・CRUD API
- [x] Phase 2: CLI（`trcli`、CI/CD組み込み用）
- [x] Phase 3: Wails + Vue3 UI
- [x] Phase 4: 仕上げ・署名・配布・LP

macOSアプリ（署名・公証済み）は [GitHub Releases](https://github.com/chankei613/ai-task-router/releases) から、
ランディングページは https://ai-task-router-three.vercel.app/ から入手できる。
アプリ内のHelpタブに使い方の説明がある。

## 使い方（開発用ヘッドレスサーバー）

```bash
go mod tidy
go run ./cmd/trserve   # :8427 でAPIサーバー起動
go run ./cmd/smoketest
```

初回起動時、Claude/GPT/Geminiの主要モデルがデフォルト単価でシードされる（`internal/db/models.go`の
`DefaultModels()`）。実際の契約単価に合わせて編集・追加できる。

### ルーティングを試す

```bash
curl -X POST localhost:8427/api/v1/route \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"task_type":"long-doc-summary","required_capabilities":["long_context"],"min_quality_tier":"standard"}'
```

### ルーティングのロジック

1. 有効なモデルの中から、要求されたcapabilityを全て満たすものに絞る
2. 品質下限（economy < standard < premium）を満たすものに絞る
3. 残った候補の中で入力+出力単価合計が最も安いものを選ぶ
4. 該当なしなら理由つきで「該当なし」を返す（無言でnullにしない）

### CI/CDへの組み込み（`trcli`）

```bash
export TR_API_KEY=...
trcli -task summarization -require vision -quality economy
# 標準出力: provider/model_id
# 標準エラー: 選定理由
# 該当モデルがなければexit 1
```

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST/GET/DELETE | `/api/v1/keys` | APIキー管理 |
| POST | `/api/v1/models` | モデル登録・更新（upsert） |
| GET | `/api/v1/models` | モデル一覧 |
| PATCH | `/api/v1/models/{id}/enabled` | 有効/無効切り替え |
| DELETE | `/api/v1/models/{id}` | モデル削除 |
| POST | `/api/v1/route` | ルーティング判定（RoutingLog作成） |
| GET | `/api/v1/routes` | 判定履歴 |
| GET | `/api/v1/routes/{id}` | 単一判定の詳細 |

## ディレクトリ構成

```
internal/db/       GORMモデル（ModelSpec/RoutingLog/AgentKey）+ デフォルトモデルのシード
internal/api/       REST API（models/route/routes/keys）+ ルーティングエンジン + 認証ミドルウェア
cmd/trserve/        開発用ヘッドレスAPIサーバー
cmd/trcli/          CI/CD組み込み用CLI
cmd/smoketest/      通しスモークテスト
```
