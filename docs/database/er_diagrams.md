# ER図

## 現在実装予定のER図

```mermaid
erDiagram
    users ||--o{ user_identities : "has"
    users ||--o{ sessions : "has"
    users ||--o{ tenant_memberships : "belongs to"

    tenants ||--o{ tenant_memberships : "has members"
    tenants ||--o{ tenant_join_codes : "has codes"

    tenant_memberships ||--o| sessions : "active in"

    oauth_states ||--|| users : "authenticates"
    console_sessions ||--|| tenants : "manages"

    users {
        uuid id PK "ユーザーの一意識別子"
        text email UK "メールアドレス（ユニーク）"
        text name "表示名"
        text icon "プロフィール画像URL"
        timestamptz created_at "作成日時"
        timestamptz updated_at "更新日時"
    }

    user_identities {
        uuid id PK "ID連携の識別子"
        uuid user_id FK "usersテーブルへの参照"
        text provider "認証プロバイダ（google）"
        text provider_sub "プロバイダ側のユーザーID"
        timestamptz created_at "作成日時"
        timestamptz updated_at "更新日時"
    }

    sessions {
        text session_id PK "セッション識別子"
        uuid user_id FK "ユーザーID"
        uuid active_membership_id FK "現在アクティブなメンバーシップ"
        timestamptz created_at "作成日時"
        timestamptz expires_at "有効期限"
        text csrf_token "CSRF対策トークン"
        boolean revoked "無効化フラグ"
    }

    oauth_states {
        text state PK "OAuth state"
        text code_verifier "PKCE verifier"
        text nonce "リプレイ攻撃対策"
        timestamptz created_at "作成日時"
        timestamptz consumed_at "使用済み日時"
    }

    tenants {
        uuid id PK "テナント識別子"
        text organization_id "組織ID（現在は定数）"
        text name "テナント名"
        text description "説明（NOT NULL DEFAULT ''）"
        text tenant_type "種別"
        timestamptz created_at "作成日時"
        timestamptz updated_at "更新日時"
    }

    tenant_join_codes {
        uuid id PK "参加コード識別子"
        uuid tenant_id FK "テナントID"
        text code UK "参加コード"
        timestamptz expires_at "有効期限"
        integer max_uses "最大使用回数"
        integer used_count "使用済み回数"
        timestamptz created_at "作成日時"
    }

    tenant_memberships {
        uuid id PK "メンバーシップ識別子"
        uuid tenant_id FK "テナントID"
        uuid user_id FK "ユーザーID"
        text role "権限（member/admin/owner）"
        text status "状態（active/invited/suspended）"
        timestamptz joined_at "参加日時"
        timestamptz left_at "退出日時"
    }

    console_sessions {
        text session_id PK "セッション識別子"
        text organization_id "組織ID"
        timestamptz created_at "作成日時"
        timestamptz expires_at "有効期限"
    }
```

## 将来拡張時の完全ER図

将来的にAdmin LayerとIntegration Layerを追加した場合の完全な構成については、[multi_db_architecture.md](../multi_db_architecture.md)を参照してください。

### 主な追加要素

- **organizations テーブル**: 組織の最上位管理
- **organization_domains テーブル**: 組織レベルのドメイン管理
- **groups テーブル**: Tenant内の小グループ
- **audit_logs テーブル**: 監査ログ

## テーブル関係の説明

### 認証フロー
1. `oauth_states` → 一時的な認証状態
2. `users` → ユーザー作成
3. `user_identities` → Google連携
4. `sessions` → ログインセッション

### テナント管理フロー
1. `tenants` → Console経由で作成
2. `tenant_join_codes` → 参加コード発行
3. `tenant_memberships` → ユーザー参加
4. `sessions.active_membership_id` → アクティブなテナント選択