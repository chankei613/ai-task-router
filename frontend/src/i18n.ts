import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'AI Task Router',
    'lang.toggle': 'JA',
    'nav.models': 'Models',
    'nav.route': 'Route',
    'nav.help': 'Help',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'models.title': 'Model Registry',
    'models.empty': 'No models registered.',
    'models.new': 'Register a model',
    'models.new.provider': 'Provider (e.g. anthropic)',
    'models.new.modelId': 'Model ID (e.g. claude-sonnet-5)',
    'models.new.tier': 'Quality tier',
    'models.new.inputPrice': 'Input $/1M tokens',
    'models.new.outputPrice': 'Output $/1M tokens',
    'models.new.capabilities': 'Capabilities (comma-separated)',
    'models.new.save': 'Save',
    'models.enabled': 'Enabled',
    'models.disabled': 'Disabled',
    'models.delete': 'Delete',
    'models.delete.confirm': 'Delete this model from the registry?',

    'tier.economy': 'economy',
    'tier.standard': 'standard',
    'tier.premium': 'premium',

    'route.title': 'Try a route',
    'route.form.task': 'Task type (free text)',
    'route.form.capabilities': 'Required capabilities (comma-separated)',
    'route.form.tier': 'Minimum quality tier',
    'route.form.source': 'Source (optional)',
    'route.form.submit': 'Route',
    'route.result.chosen': 'Chosen',
    'route.result.none': 'No model matched',
    'route.result.reasoning': 'Reasoning',

    'history.title': 'Routing history',
    'history.empty': 'No routing decisions yet.',
    'history.open': 'View',
    'history.none': 'no match',

    'detail.title': 'Routing decision',
    'detail.back': 'Back to route',

    'help.title': 'Help',
    'help.intro': 'How the model registry and routing engine fit together.',
    'help.what.title': 'What this app does',
    'help.what.body': 'This app decides which AI model to use for a task — it never calls an AI provider itself. Register models with a provider, price, quality tier, and capability tags. Ask for a route by describing required capabilities and a minimum quality tier, and it returns the cheapest registered model that qualifies, with the reasoning spelled out. Something else — your own code, a CI job, another comet-taskAI product — is responsible for actually calling the chosen model.',
    'help.start.title': 'Getting started',
    'help.start.1': 'Check the Models tab — Claude, GPT, and Gemini models are pre-seeded with rough pricing. Edit them to match your actual contract, or add your own.',
    'help.start.2': 'Go to Route, describe a task, and set any required capabilities and a minimum quality tier.',
    'help.start.3': 'The result shows the chosen model and exactly why — or why nothing qualified.',
    'help.start.4': 'Every decision is kept in history, so you can see which models get picked most often.',
    'help.stuck.title': 'Common snags',
    'help.stuck.1': 'Want to call this from a script or CI? Use the `trcli` command bundled in the repository — it prints the chosen provider/model to stdout and exits 1 if nothing matched.',
    'help.stuck.2': 'No match? Check whether the model you expected is disabled, or missing one of the required capability tags.',
    'help.stuck.3': 'This app never calls an AI provider. It only decides — actually running the request against the chosen model is up to whatever asked for the route.',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'External systems (CI, scripts) can request routing decisions here.',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app?',
  },
  ja: {
    'app.subtitle': 'AI Task Router',
    'lang.toggle': 'EN',
    'nav.models': 'モデル',
    'nav.route': 'ルーティング',
    'nav.help': 'ヘルプ',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'models.title': 'モデル登録',
    'models.empty': 'まだモデルが登録されていません。',
    'models.new': 'モデルを登録',
    'models.new.provider': 'プロバイダー（例: anthropic）',
    'models.new.modelId': 'モデルID（例: claude-sonnet-5）',
    'models.new.tier': '品質階層',
    'models.new.inputPrice': '入力 $/1Mトークン',
    'models.new.outputPrice': '出力 $/1Mトークン',
    'models.new.capabilities': 'capability（カンマ区切り）',
    'models.new.save': '保存',
    'models.enabled': '有効',
    'models.disabled': '無効',
    'models.delete': '削除',
    'models.delete.confirm': 'このモデルを登録から削除しますか？',

    'tier.economy': 'economy',
    'tier.standard': 'standard',
    'tier.premium': 'premium',

    'route.title': 'ルーティングを試す',
    'route.form.task': 'タスク種別（自由記述）',
    'route.form.capabilities': '必要なcapability（カンマ区切り）',
    'route.form.tier': '品質下限',
    'route.form.source': '実行元（任意）',
    'route.form.submit': '実行',
    'route.result.chosen': '選定結果',
    'route.result.none': '該当モデルなし',
    'route.result.reasoning': '理由',

    'history.title': 'ルーティング履歴',
    'history.empty': 'まだルーティング履歴がありません。',
    'history.open': '詳細',
    'history.none': '該当なし',

    'detail.title': 'ルーティング判定',
    'detail.back': 'ルーティングへ戻る',

    'help.title': 'ヘルプ',
    'help.intro': 'モデル登録とルーティングエンジンがどう連動するかをまとめました。',
    'help.what.title': 'このアプリでできること',
    'help.what.body': 'タスクにどのAIモデルを使うべきかを決定するアプリです — 本アプリ自身はAIプロバイダーを一切呼び出しません。プロバイダー・単価・品質階層・capabilityタグ付きでモデルを登録し、必要なcapabilityと品質下限を指定してルーティングを依頼すると、条件を満たす最安のモデルが理由つきで返ります。実際にそのモデルを呼び出すのは、あなた自身のコード・CIジョブ・他のcomet-taskAI製品など、依頼した側の仕事です。',
    'help.start.title': 'はじめに',
    'help.start.1': 'モデルタブを確認 — Claude・GPT・Geminiの主要モデルが目安単価であらかじめ登録されています。実際の契約に合わせて編集するか、独自のモデルを追加してください。',
    'help.start.2': 'ルーティングタブでタスクを記述し、必要なcapabilityと品質下限を設定します。',
    'help.start.3': '結果には選ばれたモデルと、なぜそれが選ばれたか（あるいは、なぜ何も該当しなかったか）が表示されます。',
    'help.start.4': '全ての判定は履歴に残るので、どのモデルがよく選ばれているか確認できます。',
    'help.stuck.title': 'よくある詰まりどころ',
    'help.stuck.1': 'スクリプトやCIから呼びたい → リポジトリ同梱の`trcli`コマンドを使ってください。選ばれたprovider/model_idを標準出力に表示し、該当なしならexit 1で終了します。',
    'help.stuck.2': '該当なしになる → 期待していたモデルが無効化されていないか、必要なcapabilityタグが不足していないか確認してください。',
    'help.stuck.3': '本アプリはAIプロバイダーを一切呼び出しません。決定するだけで、実際に選ばれたモデルへリクエストを送るのは依頼した側の仕事です。',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': '外部システム（CI・スクリプト等）はここでルーティング判定を依頼できます。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
