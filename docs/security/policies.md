# セキュリティポリシー

KeyHubのセキュリティ設計と実装ガイドラインです。

## 認証セキュリティ

### Organization ID管理

#### 現在の実装（開発環境）

```go
const ORGANIZATION_ID = "550e8400-e29b-41d4-a716-446655440000"  // UUID形式
const ORGANIZATION_KEY = "org_key_example_12345"
```

#### 本番環境

**config.yaml**:
```yaml
organization:
  id: "550e8400-e29b-41d4-a716-446655440000"  # UUID形式
  key: "your_secure_organization_key_here"

security:
  session:
    app_ttl: 168h      # 7日間
    console_ttl: 24h   # 24時間
```

**Goでの読み込み**:
```go
type Config struct {
    Organization struct {
        ID  string `yaml:"id" validate:"required,uuid"`
        Key string `yaml:"key" validate:"required,min=32"`
    } `yaml:"organization"`
    Security struct {
        Session struct {
            AppTTL     time.Duration `yaml:"app_ttl"`
            ConsoleTTL time.Duration `yaml:"console_ttl"`
        } `yaml:"session"`
    } `yaml:"security"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    // バリデーション
    if err := validator.New().Struct(cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

### セッション管理

#### App側（エンドユーザー）

| 設定項目 | 値 | 説明 |
|---------|-----|------|
| 有効期限 | 7日間 | 最終アクティビティから |
| 保存方法 | HttpOnly Cookie | XSS対策 |
| セッションID | 128ビット乱数 | 予測不可能 |
| CSRF対策 | トークン検証 | フォーム送信時 |

#### Console側（管理者）

| 設定項目 | 値 | 説明 |
|---------|-----|------|
| 有効期限 | 24時間 | より厳格な制限 |
| 保存方法 | JWT + Cookie | 署名付きトークン |
| 更新 | 不可 | 再ログイン必須 |

## データアクセス制御

### 現在の実装（アプリケーションレベル）

```go
// Tenantメンバーシップチェック
func (s *Service) CheckTenantAccess(ctx context.Context, userID, tenantID string) error {
    membership, err := s.db.GetMembership(ctx, userID, tenantID)
    if err != nil {
        return status.Error(codes.NotFound, "membership not found")
    }

    if membership.Status != "active" {
        return status.Error(codes.PermissionDenied, "membership not active")
    }

    return nil
}

// 権限レベルチェック
func (s *Service) RequireRole(ctx context.Context, requiredRole string) error {
    membership := GetMembershipFromContext(ctx)

    roleHierarchy := map[string]int{
        "member": 1,
        "admin":  2,
        "owner":  3,
    }

    if roleHierarchy[membership.Role] < roleHierarchy[requiredRole] {
        return status.Error(codes.PermissionDenied, "insufficient permissions")
    }

    return nil
}
```

### 将来的なRLS実装

```sql
-- Row Level Security用関数
CREATE FUNCTION current_tenant_id() RETURNS UUID AS $$
    SELECT (current_setting('app.tenant_id', true))::UUID
$$ LANGUAGE sql STABLE;

-- テーブルにRLS適用
ALTER TABLE tenant_specific_data ENABLE ROW LEVEL SECURITY;

-- ポリシー定義
CREATE POLICY tenant_isolation ON tenant_specific_data
    USING (tenant_id = current_tenant_id());
```

## 監査ログ

### ログ記録対象

#### 必須記録イベント

- ログイン/ログアウト
- Tenant作成/削除
- メンバー追加/削除
- 権限変更
- 参加コード発行

#### ログ形式

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "timestamp": "2024-01-15T10:30:00Z",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "tenant.created",
  "actor": {
    "type": "user",
    "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
    "email": "admin@example.com",
    "ip": "192.168.1.1"
  },
  "resource": {
    "type": "tenant",
    "id": "d4e5f6a7-b8c9-0123-defa-456789012345",
    "name": "情報学部"
  },
  "result": "success",
  "details": {
    "changes": {...}
  }
}
```

### ログ保持

- **保持期間**: 90日間（コンプライアンス要件に応じて調整）
- **保存先**: PostgreSQL（将来的にはS3等へアーカイブ）
- **アクセス**: 読み取り専用、改ざん防止

## 入力検証

### バリデーションルール

```go
// Email検証
func ValidateEmail(email string) error {
    if !emailRegex.MatchString(email) {
        return ErrInvalidEmail
    }
    return nil
}

// Tenant名検証
func ValidateTenantName(name string) error {
    if len(name) < 2 || len(name) > 100 {
        return ErrInvalidLength
    }
    if !nameRegex.MatchString(name) {
        return ErrInvalidCharacters
    }
    return nil
}

// 参加コード検証
func ValidateJoinCode(code string) error {
    if len(code) != 10 {
        return ErrInvalidCodeLength
    }
    if !codeRegex.MatchString(code) {
        return ErrInvalidCodeFormat
    }
    return nil
}
```

## SQLインジェクション対策

### パラメータバインディング

```go
// 正しい例
query := `
    SELECT * FROM tenants
    WHERE organization_id = $1 AND name = $2
`
rows, err := db.Query(query, orgID, name)

// NGな例（絶対に避ける）
query := fmt.Sprintf(
    "SELECT * FROM tenants WHERE name = '%s'",
    name,  // SQLインジェクションの危険
)
```

### SQLC使用（推奨）

```sql
-- queries/tenant.sql
-- name: GetTenant :one
SELECT * FROM tenants
WHERE id = $1 AND organization_id = $2;

-- name: CreateTenant :one
INSERT INTO tenants (name, slug, organization_id, tenant_type)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

## XSS対策

### 出力エスケープ

```go
// HTMLテンプレート
template.HTMLEscapeString(userInput)

// JSON出力
json.Marshal(data)  // 自動的にエスケープ

// React（フロントエンド）
// デフォルトでエスケープされる
<div>{userInput}</div>

// dangerouslySetInnerHTMLは避ける
```

## CSRF対策

### トークン実装

```go
// トークン生成
func GenerateCSRFToken() string {
    token := make([]byte, 32)
    rand.Read(token)
    return base64.URLEncoding.EncodeToString(token)
}

// トークン検証
func ValidateCSRFToken(sessionToken, requestToken string) error {
    expectedToken := GetCSRFToken(sessionToken)
    if !hmac.Equal([]byte(expectedToken), []byte(requestToken)) {
        return ErrInvalidCSRFToken
    }
    return nil
}
```

## 暗号化

### パスワード（将来実装）

```go
// bcryptを使用
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost,
)

// 検証
err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
```

### 機密データ

- **保存時**: AES-256-GCM暗号化
- **転送時**: TLS 1.2以上必須
- **キー管理**: `config.yaml`（開発環境）、AWS Secrets Manager / HashiCorp Vault（本番環境）

**config.yaml例**:
```yaml
encryption:
  key: "base64_encoded_32_byte_key_here"  # 開発環境のみ

# 本番環境ではキーは外部シークレット管理から読み込む
# secrets_manager:
#   aws_region: "ap-northeast-1"
#   secret_name: "keyhub/encryption-key"
```

## セキュリティヘッダー

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        next.ServeHTTP(w, r)
    })
}
```

## インシデント対応

### 検出

- ログ監視（異常なアクセスパターン）
- レート制限超過
- 認証失敗の連続

### 対応手順

1. **検出**: アラート受信
2. **評価**: 影響範囲の特定
3. **封じ込め**: 該当セッション/ユーザーの無効化
4. **根絶**: 脆弱性の修正
5. **復旧**: サービス正常化
6. **事後対応**: 再発防止策の実施