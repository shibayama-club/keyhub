# データベーススキーマ定義

## 現在実装のテーブル一覧

### 認証・ユーザー管理
- `users` - ユーザー基本情報
- `user_identities` - 外部認証プロバイダ連携
- `sessions` - ログインセッション
- `oauth_states` - OAuth認証の一時状態

### 組織・テナント管理
- `tenants` - 部門・研究室などの組織単位
- `tenant_memberships` - ユーザーとテナントの所属関係
- `tenant_join_codes` - 参加用招待コード

### 管理
- `console_sessions` - Console App用セッション

## スキーマ定義

### users テーブル

```sql
CREATE TABLE users (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    icon TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_users_email ON users(LOWER(email));
```

### tenants テーブル

```sql
CREATE TABLE tenants (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id TEXT DEFAULT 'ORG-DEFAULT-001', -- 将来的にFK化
    name TEXT NOT NULL,
    slug TEXT UNIQUE,
    description TEXT,
    tenant_type TEXT DEFAULT 'department', -- department, laboratory, division
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (organization_id, name)
);

CREATE INDEX idx_tenants_org ON tenants(organization_id);
CREATE UNIQUE INDEX idx_tenants_org_name ON tenants(organization_id, name);
```

### tenant_memberships テーブル

```sql
CREATE TABLE tenant_memberships (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL DEFAULT 'member', -- member, admin (将来使用予定)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, user_id)
);

CREATE INDEX idx_memberships_user ON tenant_memberships(user_id);
```

### sessions テーブル

```sql
CREATE TABLE sessions (
    session_id TEXT NOT NULL,
    user_id UUID NOT NULL,
    active_membership_id UUID, -- 現在アクティブなmembership
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    csrf_token TEXT,
    revoked BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (session_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (active_membership_id) REFERENCES tenant_memberships(id) ON DELETE SET NULL
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

## トリガーとファンクション

### 自動更新トリガー

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tenants_updated_at
BEFORE UPDATE ON tenants
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## データ整合性ルール

| ルール | 説明 | 実装方法 |
|--------|------|----------|
| **一意性保証** | メールアドレス、参加コードは重複不可 | UNIQUE制約 |
| **参照整合性** | 削除時は関連データも削除 | CASCADE DELETE |
| **有効期限管理** | 期限切れセッション・コードは無効 | アプリケーション層でチェック |
| **権限階層** | owner > admin > member | アプリケーション層で制御 |
| **組織整合性** | 全テナントは組織に所属必須 | NOT NULL制約 + デフォルト値 |