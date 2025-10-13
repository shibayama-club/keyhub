# データベーステーブル詳細仕様

このドキュメントは、KeyHubで使用される全テーブルとカラムの詳細な説明です。

## 目次

1. [認証・ユーザー管理テーブル](#1-認証ユーザー管理テーブル)
2. [組織・テナント管理テーブル](#2-組織テナント管理テーブル)
3. [管理用テーブル](#3-管理用テーブル)
4. [将来拡張用テーブル](#4-将来拡張用テーブル)

---

## 1. 認証・ユーザー管理テーブル

### users テーブル
**目的**: システムに登録されたユーザーの基本情報を管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | ユーザーの一意識別子。自動生成される | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| email | TEXT | NOT NULL, UNIQUE | Googleアカウントのメールアドレス。ログインに使用 | `user@example.com` |
| name | TEXT | - | ユーザーの表示名。Googleプロフィールから取得 | `山田太郎` |
| icon | TEXT | - | プロフィール画像のURL。Googleプロフィールから取得 | `https://lh3.googleusercontent.com/...` |
| created_at | TIMESTAMPTZ | NOT NULL | アカウント作成日時 | `2024-01-15 10:30:00+09` |
| updated_at | TIMESTAMPTZ | NOT NULL | 最終更新日時。プロフィール変更時に更新 | `2024-01-20 15:45:00+09` |

**インデックス**:
- `PRIMARY KEY (id)`
- `UNIQUE INDEX idx_users_email ON LOWER(email)`

**使用例**:
- ユーザーがGoogleログインすると自動作成
- プロフィール表示で参照
- メンバー一覧表示で利用

**関連テーブル**:
- `user_identities`: 1対多（認証プロバイダ情報）
- `sessions`: 1対多（ログインセッション）
- `tenant_memberships`: 1対多（所属情報）

---

### user_identities テーブル
**目的**: 外部認証プロバイダとの連携情報を管理（将来的に複数プロバイダ対応）

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | レコード識別子 | `b2c3d4e5-f6a7-8901-bcde-f23456789012` |
| user_id | UUID | FK(users.id), NOT NULL | 対応するユーザー | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| provider | TEXT | NOT NULL | 認証プロバイダ名 | `google` |
| provider_sub | TEXT | NOT NULL | プロバイダ側のユーザー識別子 | `1234567890987654321` |
| created_at | TIMESTAMPTZ | NOT NULL | 連携開始日時 | `2024-01-15 10:30:00+09` |
| updated_at | TIMESTAMPTZ | NOT NULL | 最終更新日時 | `2024-01-15 10:30:00+09` |

**インデックス**:
- `PRIMARY KEY (id)`
- `FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- `UNIQUE INDEX (user_id, provider, provider_sub)`

**使用例**:
- 同じユーザーが複数の認証方法を使う場合の管理
- Google以外の認証プロバイダ追加時の拡張

**ビジネスルール**:
- 同一プロバイダで同一ユーザーの重複登録は不可
- ユーザー削除時は自動削除（CASCADE）

---

### sessions テーブル
**目的**: ユーザーのログインセッションを管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| session_id | TEXT | PK, NOT NULL | セッション識別子。ランダム文字列 | `sess_abc123xyz789...` |
| user_id | UUID | FK(users.id), NOT NULL | ログインユーザー | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| active_membership_id | UUID | FK(tenant_memberships.id) | 現在選択中のテナント | `c3d4e5f6-a7b8-9012-cdef-345678901234` |
| created_at | TIMESTAMPTZ | NOT NULL | セッション開始日時 | `2024-01-15 10:30:00+09` |
| expires_at | TIMESTAMPTZ | NOT NULL | セッション有効期限（7日後） | `2024-01-22 10:30:00+09` |
| csrf_token | TEXT | - | CSRF対策用トークン | `csrf_token_abc123...` |
| revoked | BOOLEAN | DEFAULT FALSE | ログアウト済みフラグ | `false` |

**インデックス**:
- `PRIMARY KEY (session_id)`
- `INDEX idx_sessions_user ON (user_id)`
- `INDEX idx_sessions_expires ON (expires_at)`
- `FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- `FOREIGN KEY (active_membership_id) REFERENCES tenant_memberships(id) ON DELETE SET NULL`

**使用例**:
- ログイン状態の維持
- アクティブなテナントの記憶
- セキュリティ（CSRF対策、有効期限管理）

**ビジネスルール**:
- 有効期限切れセッションは定期削除（GC）
- ログアウト時は`revoked=true`に設定
- 1ユーザーあたり複数セッション許可

---

### oauth_states テーブル
**目的**: OAuth認証フローの一時的な状態を保存（セキュリティ対策）

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| state | TEXT | PK, NOT NULL | CSRF対策の一意な文字列 | `state_random123...` |
| code_verifier | TEXT | NOT NULL | PKCE認証用の検証コード | `verifier_abc456...` |
| nonce | TEXT | NOT NULL | IDトークン検証用のnonce | `nonce_xyz789...` |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 | `2024-01-15 10:29:00+09` |
| consumed_at | TIMESTAMPTZ | - | 使用済み日時（再利用防止） | `2024-01-15 10:30:00+09` |

**インデックス**:
- `PRIMARY KEY (state)`
- `INDEX idx_oauth_states_created ON (created_at)`

**使用例**:
- Google OAuth認証の開始から完了までの状態管理
- 10-15分で自動削除される一時データ

**ビジネスルール**:
- 作成から15分で自動削除
- 使用済み（consumed）のstateは再利用不可
- 定期的なガベージコレクション実行

---

## 2. 組織・テナント管理テーブル

### tenants テーブル
**目的**: 部門や研究室などの組織単位を管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | テナントの一意識別子 | `d4e5f6a7-b8c9-0123-defa-456789012345` |
| organization_id | TEXT | DEFAULT 'ORG-DEFAULT-001' | 所属組織ID（現在は固定値、将来FK化） | `ORG-DEFAULT-001` |
| name | TEXT | NOT NULL | テナント名 | `情報学部` |
| description | TEXT | NOT NULL DEFAULT '' | テナントの説明（未記入時は空文字列） | `情報学部の研究・教育部門` |
| tenant_type | TEXT | DEFAULT 'department' | テナントの種別 | `department`, `laboratory`, `division` |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 | `2024-01-10 09:00:00+09` |
| updated_at | TIMESTAMPTZ | NOT NULL | 更新日時 | `2024-01-10 09:00:00+09` |

**インデックス**:
- `PRIMARY KEY (id)`
- `UNIQUE INDEX (organization_id, name)`
- `INDEX idx_tenants_org ON (organization_id)`

**使用例**:
- Console Appでテナント作成
- ユーザーが参加するテナントの選択
- メンバー管理の単位

**ビジネスルール**:
- 同一組織内でテナント名の重複不可
- tenant_typeは事前定義された値のみ許可（`department`, `laboratory`, `division`）
- descriptionは空文字列を許可（未記入時のデフォルト値）

**関連テーブル**:
- `tenant_memberships`: 1対多（メンバー管理）
- `tenant_join_codes`: 1対多（参加コード）
- `groups`: 1対多（将来実装）

---

### tenant_join_codes テーブル
**目的**: テナントへの参加用招待コードを管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | レコード識別子 | `f6a7b8c9-d0e1-2345-fabc-678901234567` |
| tenant_id | UUID | FK(tenants.id), NOT NULL | 対象テナント | `d4e5f6a7-b8c9-0123-defa-456789012345` |
| code | TEXT | NOT NULL, UNIQUE | 参加コード（8-12文字） | `JOIN2024XYZ` |
| expires_at | TIMESTAMPTZ | - | 有効期限（nullは無期限） | `2024-12-31 23:59:59+09` |
| max_uses | INTEGER | DEFAULT 0 | 最大使用回数（0は無制限） | `100` |
| used_count | INTEGER | DEFAULT 0 | 現在の使用回数 | `25` |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 | `2024-01-10 10:00:00+09` |

**インデックス**:
- `PRIMARY KEY (id)`
- `UNIQUE INDEX (code)`
- `FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE`
- `INDEX idx_join_codes_tenant ON (tenant_id)`
- `INDEX idx_join_codes_expires ON (expires_at)`

**使用例**:
- メールドメインが判別できない場合の参加方法
- 期間限定の招待キャンペーン
- 人数制限付きの参加管理

**ビジネスルール**:
- 有効期限切れコードは使用不可
- max_uses超過時は使用不可
- コードは大文字英数字のみ
- 使用時にused_countをインクリメント

**セキュリティ**:
- コードは予測困難な乱数生成
- 短時間での連続試行はブロック
- 使用履歴を監査ログに記録

---

### tenant_memberships テーブル
**目的**: ユーザーとテナントの所属関係を管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | メンバーシップ識別子 | `a7b8c9d0-e1f2-3456-abcd-789012345678` |
| tenant_id | UUID | FK(tenants.id), NOT NULL | 所属テナント | `d4e5f6a7-b8c9-0123-defa-456789012345` |
| user_id | UUID | FK(users.id), NOT NULL | 所属ユーザー | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| role | TEXT | DEFAULT 'member' | 権限レベル | `member`, `admin`, `owner` |
| status | TEXT | DEFAULT 'active' | メンバーシップ状態 | `active`, `invited`, `suspended` |
| joined_at | TIMESTAMPTZ | NOT NULL | 参加日時 | `2024-01-15 11:00:00+09` |
| left_at | TIMESTAMPTZ | - | 退出日時 | null |

**インデックス**:
- `PRIMARY KEY (id)`
- `UNIQUE INDEX (tenant_id, user_id)`
- `FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE`
- `FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- `INDEX idx_memberships_user ON (user_id)`
- `INDEX idx_memberships_tenant ON (tenant_id)`

**権限レベル**:
- `owner`: 全権限（テナント削除、管理者任命）
- `admin`: 管理権限（メンバー管理、設定変更）
- `member`: 基本権限（閲覧、基本操作）

**状態遷移**:
- `invited` → `active`: 招待承認時
- `active` → `suspended`: 一時停止時
- `active` → （削除）: 退出時（left_atを記録）

**使用例**:
- ユーザーのテナント所属管理
- 権限チェック
- メンバー一覧表示

**ビジネスルール**:
- 1ユーザー1テナントに重複参加不可
- ownerは最低1名必須
- 退出時はleft_atを記録（物理削除しない）

---

## 3. 管理用テーブル

### console_sessions テーブル
**目的**: Console App用のセッション管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| session_id | TEXT | PK, NOT NULL | セッション識別子 | `console_sess_xyz123...` |
| organization_id | TEXT | NOT NULL | 組織ID（現在は固定） | `ORG-DEFAULT-001` |
| created_at | TIMESTAMPTZ | NOT NULL | セッション開始日時 | `2024-01-15 09:00:00+09` |
| expires_at | TIMESTAMPTZ | NOT NULL | 有効期限（24時間） | `2024-01-16 09:00:00+09` |

**インデックス**:
- `PRIMARY KEY (session_id)`
- `INDEX idx_console_sessions_expires ON (expires_at)`

**使用例**:
- Console管理者のログイン状態管理
- 24時間の短い有効期限でセキュリティ強化

**ビジネスルール**:
- App側より短い有効期限（24時間）
- 更新不可（期限切れ後は再ログイン）
- Organization IDの検証必須

---

## 4. 将来拡張用テーブル

### groups テーブル（将来実装）
**目的**: テナント内の小グループを管理

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | グループ識別子 | `b8c9d0e1-f2a3-4567-bcde-890123456789` |
| tenant_id | UUID | FK(tenants.id), NOT NULL | 所属テナント | `d4e5f6a7-b8c9-0123-defa-456789012345` |
| name | TEXT | NOT NULL | グループ名 | `AI研究室` |
| description | TEXT | NOT NULL DEFAULT '' | グループ説明（未記入時は空文字列） | `人工知能の研究を行う研究室` |
| parent_group_id | UUID | FK(groups.id) | 親グループ（階層構造） | null |
| is_optional | BOOLEAN | DEFAULT TRUE | オプショナルグループか | `true` |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 | `2024-02-01 10:00:00+09` |
| updated_at | TIMESTAMPTZ | NOT NULL | 更新日時 | `2024-02-01 10:00:00+09` |

**将来の使用例**:
- テナント内の細かい組織単位
- プロジェクトチーム管理
- 研究室や部署内のサブグループ

---

### audit_logs テーブル（将来実装）
**目的**: システム全体の操作履歴を記録（監査・デバッグ用）

| カラム名 | データ型 | 制約 | 説明 | 例 |
|---------|---------|------|------|-----|
| id | UUID | PK, NOT NULL | ログ識別子 | `c9d0e1f2-a3b4-5678-cdef-901234567890` |
| organization_id | TEXT | DEFAULT 'ORG-DEFAULT-001' | 組織ID | `ORG-DEFAULT-001` |
| event_type | TEXT | NOT NULL | イベント種別 | `tenant.created`, `user.joined` |
| actor_type | TEXT | NOT NULL | 実行者種別 | `user`, `console`, `system` |
| actor_id | TEXT | - | 実行者ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| resource_type | TEXT | - | 対象リソース種別 | `tenant`, `user` |
| resource_id | TEXT | - | 対象リソースID | `d4e5f6a7-b8c9-0123-defa-456789012345` |
| details | JSONB | - | 詳細情報（JSON形式） | `{"action": "create", "name": "情報学部"}` |
| created_at | TIMESTAMPTZ | NOT NULL | 記録日時 | `2024-01-15 11:30:00+09` |

**将来の使用例**:
- セキュリティ監査
- 問題発生時の調査
- 使用統計の分析
- コンプライアンス対応