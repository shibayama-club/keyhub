# データベースマイグレーション戦略

KeyHubのデータベースマイグレーション戦略と実装計画です。

## 目次

1. [マイグレーション戦略](#1-マイグレーション戦略)
2. [初期セットアップSQL](#2-初期セットアップsql)
3. [将来の拡張マイグレーション](#3-将来の拡張マイグレーション)
4. [マイグレーション管理](#4-マイグレーション管理)

---

## 1. マイグレーション戦略

### 1.1 基本方針

**段階的アプローチ**
1. Phase 1: 現在実装（単一DB、固定Organization）
2. Phase 2: 機能拡張（Groups、Audit Logs追加）
3. Phase 3: マルチDB移行（Admin DB分離）

**マイグレーションツール**
- 使用ツール: golang-migrate または sqlc migrate
- バージョン管理: 連番プレフィックス
- ロールバック可能な設計

### 1.2 命名規則

```
YYYYMMDDHHMMSS_description.up.sql
YYYYMMDDHHMMSS_description.down.sql

例:
20240101120000_create_users_table.up.sql
20240101120000_create_users_table.down.sql
```

---

## 2. 初期セットアップSQL

### 2.1 完全な初期化スクリプト

```sql
-- 20240101000001_initial_schema.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users table
CREATE TABLE users (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    icon TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- 2. User Identities table
CREATE TABLE user_identities (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    provider TEXT NOT NULL,
    provider_sub TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (user_id, provider, provider_sub),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider_sub ON user_identities(provider, provider_sub);

-- 3. Tenants table
CREATE TABLE tenants (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id TEXT NOT NULL DEFAULT 'ORG-DEFAULT-001',
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    tenant_type TEXT NOT NULL DEFAULT 'team',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (organization_id, slug),
    CONSTRAINT check_tenant_type CHECK (tenant_type IN ('team', 'department', 'project', 'laboratory'))
);

CREATE INDEX idx_tenants_organization_id ON tenants(organization_id);
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_is_default ON tenants(is_default) WHERE is_default = true;

-- 4. Tenant Memberships table
CREATE TABLE tenant_memberships (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, user_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_role CHECK (role IN ('admin', 'member'))
);

CREATE INDEX idx_tenant_memberships_user_id ON tenant_memberships(user_id);

-- 5. Sessions table
CREATE TABLE sessions (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    session_id TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    active_membership_id UUID,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    last_accessed_at TIMESTAMPTZ,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (active_membership_id) REFERENCES tenant_memberships(id) ON DELETE SET NULL
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked_expires ON sessions(revoked, expires_at) WHERE revoked = false;

-- 6. OAuth States table
CREATE TABLE oauth_states (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    state TEXT NOT NULL UNIQUE,
    code_verifier TEXT NOT NULL,
    nonce TEXT NOT NULL,
    redirect_uri TEXT,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE INDEX idx_oauth_states_state ON oauth_states(state);
CREATE INDEX idx_oauth_states_created_at ON oauth_states(created_at);
CREATE INDEX idx_oauth_states_consumed ON oauth_states(consumed_at) WHERE consumed_at IS NULL;

-- 7. Tenant Join Codes table
CREATE TABLE tenant_join_codes (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ,
    max_uses INTEGER NOT NULL DEFAULT 0,
    used_count INTEGER NOT NULL DEFAULT 0,
    role TEXT NOT NULL DEFAULT 'member',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT check_max_uses CHECK (max_uses >= 0),
    CONSTRAINT check_used_count CHECK (used_count >= 0),
    CONSTRAINT check_join_role CHECK (role IN ('member', 'viewer'))
);

CREATE INDEX idx_tenant_join_codes_code ON tenant_join_codes(code);
CREATE INDEX idx_tenant_join_codes_tenant_id ON tenant_join_codes(tenant_id);
CREATE INDEX idx_tenant_join_codes_expires_at ON tenant_join_codes(expires_at);
CREATE INDEX idx_tenant_join_codes_active ON tenant_join_codes(code, expires_at) WHERE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP;

-- 8. Console Sessions table
CREATE TABLE console_sessions (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    session_id TEXT NOT NULL UNIQUE,
    organization_id TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE INDEX idx_console_sessions_session_id ON console_sessions(session_id);
CREATE INDEX idx_console_sessions_expires_at ON console_sessions(expires_at);
CREATE INDEX idx_console_sessions_active ON console_sessions(expires_at) WHERE expires_at > CURRENT_TIMESTAMP;

-- トリガー: updated_at自動更新
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_identities_updated_at BEFORE UPDATE ON user_identities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenant_memberships_updated_at BEFORE UPDATE ON tenant_memberships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 2.2 ロールバックスクリプト

```sql
-- 20240101000001_initial_schema.down.sql

-- Drop triggers
DROP TRIGGER IF EXISTS update_tenant_memberships_updated_at ON tenant_memberships;
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;
DROP TRIGGER IF EXISTS update_user_identities_updated_at ON user_identities;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS console_sessions;
DROP TABLE IF EXISTS tenant_join_codes;
DROP TABLE IF EXISTS oauth_states;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tenant_memberships;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS user_identities;
DROP TABLE IF EXISTS users;
```

### 2.3 初期データ投入

```sql
-- 20240101000002_seed_data.up.sql

-- デフォルトTenant作成
INSERT INTO tenants (
    id,
    organization_id,
    name,
    slug,
    description,
    tenant_type,
    is_default
) VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'ORG-DEFAULT-001',
    'Default Team',
    'default',
    'デフォルトのチームです',
    'team',
    true
);

-- サンプル参加コード
INSERT INTO tenant_join_codes (
    tenant_id,
    code,
    expires_at,
    max_uses,
    role
) VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'KH-DEMO-001',
    CURRENT_TIMESTAMP + INTERVAL '30 days',
    100,
    'member'
);
```

---

## 3. 将来の拡張マイグレーション

### 3.1 Groups機能追加

```sql
-- 20240201000001_add_groups.up.sql

-- Groups table
CREATE TABLE groups (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_group_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_groups_tenant_id ON groups(tenant_id);
CREATE INDEX idx_groups_parent_group_id ON groups(parent_group_id);

-- Group Memberships table
CREATE TABLE group_memberships (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_memberships_group_id ON group_memberships(group_id);
CREATE INDEX idx_group_memberships_user_id ON group_memberships(user_id);

-- Add trigger
CREATE TRIGGER update_groups_updated_at BEFORE UPDATE ON groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.2 Audit Logs追加

```sql
-- 20240301000001_add_audit_logs.up.sql

CREATE TABLE audit_logs (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id TEXT NOT NULL,
    tenant_id UUID,
    user_id UUID,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_audit_logs_organization ON audit_logs(organization_id, created_at DESC);
CREATE INDEX idx_audit_logs_tenant ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action, created_at DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- パーティショニング（月単位）
CREATE TABLE audit_logs_2024_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE audit_logs_2024_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

### 3.3 Organizations分離

```sql
-- 20240401000001_add_organizations.up.sql

-- Organizations table (将来のマルチDB構成への準備)
CREATE TABLE organizations (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    settings JSONB,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Organization Domains
CREATE TABLE organization_domains (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id TEXT NOT NULL,
    domain TEXT NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- 既存データの移行
INSERT INTO organizations (id, name, slug, status)
VALUES ('ORG-DEFAULT-001', 'Default Organization', 'default', 'active');

-- Tenantsテーブルに外部キー追加
ALTER TABLE tenants
    ADD CONSTRAINT fk_tenants_organization
    FOREIGN KEY (organization_id) REFERENCES organizations(id);
```

---

## 4. マイグレーション管理

### 4.1 マイグレーションツール設定

**golang-migrate設定**
```yaml
# migrations.yaml
database:
  driver: postgres
  datasource: ${DATABASE_URL}

migrations:
  path: ./migrations
  table: schema_migrations
```

**実行コマンド**
```bash
# 最新バージョンまでマイグレート
migrate -path ./migrations -database ${DATABASE_URL} up

# 特定バージョンまで
migrate -path ./migrations -database ${DATABASE_URL} goto 3

# ロールバック
migrate -path ./migrations -database ${DATABASE_URL} down 1

# 現在のバージョン確認
migrate -path ./migrations -database ${DATABASE_URL} version
```

### 4.2 CI/CDパイプライン統合

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main]
    paths:
      - 'migrations/**'

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Run migrations
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
          ./migrate -path ./migrations -database ${DATABASE_URL} up
```

### 4.3 バックアップ戦略

```bash
#!/bin/bash
# backup_before_migration.sh

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_NAME="keyhub"
BACKUP_DIR="/backups"

# フルバックアップ
pg_dump $DB_NAME > $BACKUP_DIR/backup_${TIMESTAMP}.sql

# 圧縮
gzip $BACKUP_DIR/backup_${TIMESTAMP}.sql

# S3へアップロード
aws s3 cp $BACKUP_DIR/backup_${TIMESTAMP}.sql.gz s3://keyhub-backups/migrations/

# 古いバックアップ削除（30日以上）
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
```

### 4.4 ヘルスチェック

```sql
-- migration_health_check.sql

-- テーブル存在確認
SELECT COUNT(*) as table_count
FROM information_schema.tables
WHERE table_schema = 'public'
AND table_type = 'BASE TABLE';

-- 必須テーブル確認
SELECT
    'users' as table_name,
    EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'users') as exists
UNION ALL
SELECT
    'tenants',
    EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'tenants')
UNION ALL
SELECT
    'tenant_memberships',
    EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'tenant_memberships')
UNION ALL
SELECT
    'sessions',
    EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions');

-- インデックス確認
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

-- 外部キー確認
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
AND tc.table_schema = 'public';
```

## ベストプラクティス

### DO
- 各マイグレーションは独立して実行可能に
- ロールバック可能な設計
- データ移行は別マイグレーションで
- 本番環境前にステージング環境でテスト

### DON'T
- 既存マイグレーションの編集
- 大量データ更新とスキーマ変更の同時実行
- 外部依存のあるマイグレーション
- ハードコーディングされた環境固有値