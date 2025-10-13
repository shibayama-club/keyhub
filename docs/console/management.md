# Console管理機能

KeyHub Console（管理者向けダッシュボード）の機能仕様です。

## 目次

1. [Console認証](#1-console認証)
2. [Tenant管理](#2-tenant管理)
3. [参加コード管理](#3-参加コード管理)
4. [メンバー管理](#4-メンバー管理)
5. [監査ログ](#5-監査ログ)

---

## 1. Console認証

### 1.1 認証方式

**Organization ID/Key認証**
- 組織に発行されたID/Keyペアで認証
- 現在実装：定数として管理
- 将来実装：Admin DBから発行

```typescript
// Console認証フォーム
interface ConsoleAuthForm {
  organizationId: string;  // UUID形式: 550e8400-e29b-41d4-a716-446655440000
  organizationKey: string; // org_key_example_12345
}
```

### 1.2 セッション管理

```sql
-- Console用セッション
CREATE TABLE console_sessions (
    session_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,  -- UUID形式
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**JWT構造**
```json
{
  "org_id": "550e8400-e29b-41d4-a716-446655440000",
  "session_id": "uuid-value",
  "exp": 1699999999,
  "iat": 1699913599,
  "sub": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 2. Tenant管理

### 2.1 Tenant一覧

**取得クエリ**
```sql
SELECT
    t.id,
    t.name,
    t.description,
    t.tenant_type,
    t.is_default,
    t.created_at,
    t.updated_at,
    COUNT(DISTINCT tm.user_id) as member_count,
    COUNT(DISTINCT tjc.id) as join_code_count
FROM tenants t
LEFT JOIN tenant_memberships tm
    ON t.id = tm.tenant_id
    AND tm.status = 'active'
LEFT JOIN tenant_join_codes tjc
    ON t.id = tjc.tenant_id
WHERE t.organization_id = $1
GROUP BY t.id
ORDER BY t.created_at DESC;
```

### 2.2 Tenant作成

**作成フォーム**
```typescript
interface CreateTenantForm {
  name: string;            // 表示名
  description?: string;    // 説明
  tenant_type: 'team' | 'department' | 'project' | 'laboratory';
  is_default?: boolean;    // デフォルトテナント
}
```

**バリデーション**
- name: 1-100文字、必須
- description: 最大500文字
- tenant_type: 定義済みの値のみ許可
- organization_idごとに1つのみis_default=true

### 2.3 Tenant編集

**更新可能フィールド**
- name
- description
- tenant_type
- is_default

**更新不可フィールド**
- id（一意識別子）
- organization_id（所属組織変更不可）
- created_at（作成日時）

### 2.4 Tenant削除

**削除前チェック**
```sql
-- アクティブメンバー確認
SELECT COUNT(*)
FROM tenant_memberships
WHERE tenant_id = $1
AND status = 'active';

-- 有効な参加コード確認
SELECT COUNT(*)
FROM tenant_join_codes
WHERE tenant_id = $1
AND (expires_at IS NULL OR expires_at > NOW());
```

**ソフトデリート推奨**
```sql
-- 実装案（将来）
UPDATE tenants
SET
    deleted_at = NOW(),
    status = 'deleted'
WHERE id = $1;
```

---

## 3. 参加コード管理

### 3.1 参加コード生成

**生成パラメータ**
```typescript
interface GenerateJoinCodeForm {
  tenant_id: string;
  expires_at?: Date;      // 有効期限
  max_uses?: number;       // 最大使用回数（0=無制限）
  role?: 'member' | 'viewer';  // 付与するロール
}
```

**コード生成ロジック**
```go
func generateJoinCode() string {
    // PREFIX-RANDOM-CHECK
    // 例: KH-X7Y9Z-A3
    prefix := "KH"
    random := generateRandomString(5, "ABCDEFGHJKLMNPQRSTUVWXYZ23456789")
    check := calculateChecksum(random)
    return fmt.Sprintf("%s-%s-%s", prefix, random, check)
}
```

### 3.2 参加コード一覧

```sql
SELECT
    tjc.id,
    tjc.code,
    tjc.expires_at,
    tjc.max_uses,
    tjc.used_count,
    tjc.role,
    tjc.created_at,
    t.name as tenant_name,
    u.name as created_by_name,
    CASE
        WHEN tjc.expires_at IS NOT NULL AND tjc.expires_at < NOW() THEN 'expired'
        WHEN tjc.max_uses > 0 AND tjc.used_count >= tjc.max_uses THEN 'exhausted'
        ELSE 'active'
    END as status
FROM tenant_join_codes tjc
JOIN tenants t ON tjc.tenant_id = t.id
LEFT JOIN users u ON tjc.created_by = u.id
WHERE t.organization_id = $1
ORDER BY tjc.created_at DESC;
```

### 3.3 使用履歴

```sql
-- 参加コード使用ログ（将来実装）
CREATE TABLE join_code_usage_logs (
    id UUID PRIMARY KEY,
    join_code_id UUID REFERENCES tenant_join_codes(id),
    user_id UUID REFERENCES users(id),
    used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT
);
```

### 3.4 参加コード無効化

```sql
-- 即座に無効化
UPDATE tenant_join_codes
SET expires_at = NOW()
WHERE id = $1;

-- または削除
DELETE FROM tenant_join_codes
WHERE id = $1;
```

---

## 4. メンバー管理

### 4.1 メンバー一覧

```sql
SELECT
    tm.id as membership_id,
    tm.tenant_id,
    tm.user_id,
    tm.role,
    tm.status,
    tm.joined_at,
    tm.left_at,
    u.email,
    u.name,
    u.icon,
    t.name as tenant_name,
    CASE
        WHEN s.active_membership_id = tm.id THEN true
        ELSE false
    END as is_active
FROM tenant_memberships tm
JOIN users u ON tm.user_id = u.id
JOIN tenants t ON tm.tenant_id = t.id
LEFT JOIN sessions s ON s.user_id = u.id AND s.revoked = false
WHERE t.organization_id = $1
ORDER BY tm.joined_at DESC;
```

### 4.2 ロール管理

**ロール更新**
```sql
UPDATE tenant_memberships
SET
    role = $1,
    updated_at = NOW()
WHERE id = $2;
```

**ロール階層**
```typescript
enum Role {
  OWNER = 'owner',      // 全権限
  ADMIN = 'admin',      // 管理権限
  MEMBER = 'member',    // 通常メンバー
  VIEWER = 'viewer'     // 閲覧のみ
}

// 権限チェック
function hasPermission(role: Role, action: string): boolean {
  const permissions = {
    owner: ['*'],
    admin: ['read', 'write', 'invite', 'manage_members'],
    member: ['read', 'write'],
    viewer: ['read']
  };

  return permissions[role].includes('*') ||
         permissions[role].includes(action);
}
```

### 4.3 メンバー削除（退会処理）

```sql
-- ソフトデリート
UPDATE tenant_memberships
SET
    status = 'inactive',
    left_at = NOW()
WHERE id = $1;

-- セッションからも削除
UPDATE sessions
SET active_membership_id = NULL
WHERE active_membership_id = $1;
```

### 4.4 一括操作

**一括招待**
```typescript
interface BulkInvite {
  emails: string[];
  tenant_id: string;
  role: Role;
  send_email?: boolean;
}
```

**CSVエクスポート**
```sql
-- メンバーリストCSV生成
COPY (
    SELECT
        u.email,
        u.name,
        t.name as tenant,
        tm.role,
        tm.joined_at
    FROM tenant_memberships tm
    JOIN users u ON tm.user_id = u.id
    JOIN tenants t ON tm.tenant_id = t.id
    WHERE t.organization_id = $1
    AND tm.status = 'active'
) TO '/tmp/members.csv' WITH CSV HEADER;
```

---

## 5. 監査ログ

### 5.1 ログ記録対象

**記録するイベント**
- Tenant作成/編集/削除
- 参加コード生成/使用/無効化
- メンバー追加/ロール変更/削除
- Console管理者ログイン

### 5.2 ログテーブル（将来実装）

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,  -- UUID形式
    tenant_id UUID,
    user_id UUID,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス
CREATE INDEX idx_audit_logs_organization ON audit_logs(organization_id, created_at DESC);
CREATE INDEX idx_audit_logs_tenant ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
```

### 5.3 ログ検索

```sql
-- 特定期間のログ取得
SELECT * FROM audit_logs
WHERE organization_id = $1
AND created_at BETWEEN $2 AND $3
AND ($4::text IS NULL OR action = $4)
AND ($5::text IS NULL OR resource_type = $5)
ORDER BY created_at DESC
LIMIT 100;
```

### 5.4 ログ保持

**保持ポリシー**
- 通常ログ：90日
- セキュリティイベント：1年
- ログイン履歴：30日

```sql
-- 定期削除バッチ
DELETE FROM audit_logs
WHERE created_at < NOW() - INTERVAL '90 days'
AND action NOT IN ('security_alert', 'permission_change');
```

## Console画面構成

### ダッシュボード
- Organization情報表示
- Tenant数、総メンバー数
- 最近のアクティビティ

### Tenant管理画面
- Tenant一覧テーブル
- 作成/編集モーダル
- メンバー数表示

### 参加コード管理画面
- コード一覧テーブル
- 生成フォーム
- 使用状況グラフ

### メンバー管理画面
- メンバー一覧テーブル
- フィルタ（Tenant、ロール、状態）
- 一括操作メニュー

### 監査ログ画面
- ログ一覧テーブル
- 検索フィルタ
- エクスポート機能

## API エンドポイント

### Console API（ConnectRPC）

```protobuf
service ConsoleService {
  // Tenant管理
  rpc ListTenants(ListTenantsRequest) returns (ListTenantsResponse);
  rpc CreateTenant(CreateTenantRequest) returns (CreateTenantResponse);
  rpc UpdateTenant(UpdateTenantRequest) returns (UpdateTenantResponse);
  rpc DeleteTenant(DeleteTenantRequest) returns (DeleteTenantResponse);

  // 参加コード管理
  rpc ListJoinCodes(ListJoinCodesRequest) returns (ListJoinCodesResponse);
  rpc GenerateJoinCode(GenerateJoinCodeRequest) returns (GenerateJoinCodeResponse);
  rpc RevokeJoinCode(RevokeJoinCodeRequest) returns (RevokeJoinCodeResponse);

  // メンバー管理
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse);
  rpc UpdateMemberRole(UpdateMemberRoleRequest) returns (UpdateMemberRoleResponse);
  rpc RemoveMember(RemoveMemberRequest) returns (RemoveMemberResponse);

  // 監査ログ
  rpc GetAuditLogs(GetAuditLogsRequest) returns (GetAuditLogsResponse);
}
```