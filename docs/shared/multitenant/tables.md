## テーブル詳細（マルチテナント拡張）

このドキュメントは、`./er.md` のERに基づき、マルチテナント関連テーブルと各カラムの意味・制約を整理したものです。app側で定義される基本テーブル（users, user_identities, sessions など）の詳細は `../../app/db/tables.md` も参照してください。

### tenants

テナント（組織・ワークスペース）を表すルートエンティティ。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。`uuid_generate_v4()` による自動採番。 |
| name | text | テナント表示名。必須。ユニーク。lower(name)にも一意インデックス。 |
| slug | text | 任意のURL用スラッグ。ユニーク。空値可。 |
| password_hash | text | 認証用ハッシュ値（用途は別ドキュメント参照）。必須。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |
| updated_at | timestamptz | 更新時刻。`now()` デフォルト。 |

備考 / 制約・インデックス
- `name` は UNIQUE。さらに `lower(name)` にも一意インデックス（大文字小文字を無視した一意性）。

### tenant_domains

メールドメインとテナントの紐付け。ドメインから候補テナントを推定する用途。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。`uuid_generate_v4()` 自動採番。 |
| tenant_id | uuid | 外部キー。`tenants(id)`。CASCADE削除。 |
| domain | text | メールドメイン。ユニーク。正規化（lowercase）推奨。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |

備考 / 制約・インデックス
- `domain` は UNIQUE。セキュリティポリシー上、lowercase正規化での一意性を推奨。

### tenant_join_codes

招待コードによる参加手段を管理。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。`uuid_generate_v4()` 自動採番。 |
| tenant_id | uuid | 外部キー。`tenants(id)`。CASCADE削除。 |
| code | text | 参加コード。ユニーク。例: 8〜12桁ランダム。 |
| expires_at | timestamptz | 期限切れ時刻。NULL可（無期限）。 |
| max_uses | integer | 最大利用回数。`0` は無制限。デフォルト `0`。 |
| used_count | integer | 使用済み回数。デフォルト `0`。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |

備考 / 制約・インデックス
- UNIQUE: `code`
- INDEX: `(tenant_id)`, `(expires_at)`（クリーンアップや検索の最適化）
- 運用上は、生の `code` を一度だけ表示し、保存時はハッシュ化を検討可能（セキュリティ強化）。

### tenant_memberships

ユーザーとテナントの所属関係。役割/状態や参加経路を持つ。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。`uuid_generate_v4()` 自動採番。 |
| tenant_id | uuid | 外部キー。`tenants(id)`。CASCADE削除。 |
| user_id | uuid | 外部キー。`users(id)`。CASCADE削除。 |
| role | text | 役割。`admin` / `member`。デフォルト `member`。 |
| status | text | 状態。`active` / `invited` / `suspended` / `left`。デフォルト `active`。 |
| joined_via | text | 参加経路。`domain` / `code` / `manual`。NULL可。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |
| left_at | timestamptz | 退会時刻。NULL可。 |

備考 / 制約・インデックス
- UNIQUE: `(tenant_id, user_id)`（同一テナントに二重参加を防止）
- INDEX: `(user_id)`（ユーザーから所属一覧を引く用途）
- CHECK制約: `role`, `status`, `joined_via` は定義済みの値に限定。

### users（app側）

ユーザー基本情報。詳細は `../../app/db/tables.md` 参照。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。`uuid_generate_v4()` 自動採番。 |
| email | text | メールアドレス。lowerで一意性。 |
| name | text | 表示名。 |
| icon | text | プロフィール画像URL。 |
| created_at | timestamptz | 作成時刻。 |
| updated_at | timestamptz | 更新時刻。 |

### user_identities（app側）

外部IDプロバイダとの紐付け。詳細は `../../app/db/tables.md` 参照。

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。 |
| user_id | uuid | 外部キー。`users(id)`。CASCADE削除。 |
| provider | text | プロバイダ（例: `google`）。 |
| provider_sub | text | プロバイダ側の一意ID（Googleの `sub` など）。 |
| created_at / updated_at | timestamptz | 作成・更新時刻。 |

### oauth_states（app側）

OAuth 2.0 / OIDC認証フロー中の一時的な状態を管理。CSRF攻撃とリプレイ攻撃を防止する。

| カラム | 型 | 説明 |
| --- | --- | --- |
| state | text | 主キー。OAuth CSRF対策のランダム文字列。認証リクエスト時に生成し、コールバック時に検証。 |
| code_verifier | text | PKCE（Proof Key for Code Exchange）のcode_verifier。必須。 |
| nonce | text | IDトークン検証用nonce。リプレイ攻撃防止。必須。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |
| consumed_at | timestamptz | 使用済み時刻。NULL可。コールバック処理完了時に記録。 |

備考 / 制約・インデックス
- INDEX: `(created_at)` - 古いレコードのクリーンアップに使用。
- ライフサイクル: 認証開始時に作成 → コールバック時に`consumed_at`を更新 → 定期的に古いレコードを削除（例: 作成から15分経過したもの）。
- セキュリティ: `consumed_at`が既に設定されているレコードは再利用を拒否することで、リプレイ攻撃を防ぐ。

### sessions（app側 + マルチテナント拡張）

アプリセッション。`active_membership_id` により現在のテナント文脈を保持。詳細は `../../app/db/tables.md` 参照。

| カラム | 型 | 説明 |
| --- | --- | --- |
| session_id | text | 主キー。クッキーに保存されるランダム文字列。 |
| user_id | uuid | 外部キー。`users(id)`。 |
| active_membership_id | uuid | 外部キー。`tenant_memberships(id)`。NULL可（未選択時）。 |
| created_at | timestamptz | 作成時刻。 |
| expires_at | timestamptz | 期限。GC対象。 |
| csrf_token | text | CSRF対策トークン。 |
| revoked | boolean | 手動無効化フラグ。 |

備考 / 制約・インデックス
- INDEX: `(active_membership_id)` - テナント切り替えやメンバーシップに基づくセッション検索を高速化。
- 外部キー: `fk_sessions_active_membership` - `active_membership_id`が`tenant_memberships(id)`を参照。データ整合性を強制し、存在しないメンバーシップIDの設定を防ぐ。
- マイグレーション: `active_membership_id`カラムは`sessions`テーブル作成時に定義されるが、外部キー制約は`tenant_memberships`テーブル作成後（`20251007060715_add_tenant_memberships_table.sql`）に追加される。

### console_sessions

コンソール（管理）用のセッション。テナント単位の管理操作に利用。

| カラム | 型 | 説明 |
| --- | --- | --- |
| session_id | text | 主キー。 |
| tenant_id | uuid | 外部キー。`tenants(id)`。CASCADE削除。 |
| created_at | timestamptz | 作成時刻。`now()` デフォルト。 |
| expires_at | timestamptz | 期限。必須。 |

備考 / 制約・インデックス
- テナント境界の管理操作に使用。TTLに従って定期的にクリーンアップ。

---

## リレーション要約

- tenants 1:N tenant_domains / tenant_join_codes / tenant_memberships / console_sessions
- users 1:N tenant_memberships / sessions
- sessions N:1 tenant_memberships（`active_membership_id`）

## RLS との関係（概要）

- 接続ごとに `SET app.membership_id = '...'` を設定し、`app.current_tenant_id()` を通じて現在テナントを決定。
- `tenant_id = app.current_tenant_id()` を条件とするRLSポリシーで、テナント境界を強制。
- 書き込み時は `WITH CHECK (tenant_id = app.current_tenant_id())` 等を併用する。

