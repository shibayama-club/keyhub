# 認証フロー詳細

KeyHubの認証フローの完全な仕様書です。

## 目次

1. [App側認証フロー](#1-app側認証フロー)
2. [Console側認証フロー](#2-console側認証フロー)
3. [Tenant参加フロー](#3-tenant参加フロー)

---

## 1. App側認証フロー

### 1.1 Google OAuth認証フロー全体

```mermaid
sequenceDiagram
    autonumber
    participant FE as Frontend
    participant HTTP as HTTP Auth
    participant RPC as ConnectRPC
    participant DB as PostgreSQL
    participant G as Google OAuth

    Note over FE: ユーザー「Googleでログイン」クリック
    FE->>HTTP: GET /auth/google/login
    HTTP->>HTTP: generate state and code_verifier and nonce
    HTTP->>DB: oauth_statesにINSERT
    HTTP-->>FE: 302 Google認可URLへリダイレクト

    FE-->>G: Google login and consent
    G-->>HTTP: callback with code and state
    HTTP->>DB: oauth_statesからstate取得と検証
    HTTP->>G: exchange tokens with code and code_verifier
    G-->>HTTP: id_token and access_token
    HTTP->>HTTP: id_token検証 [iss, aud, exp, nonce]

    rect rgb(230, 245, 255)
        HTTP->>DB: users を UPSERT
        HTTP->>DB: user_identities を UPSERT
        HTTP->>DB: sessions に INSERT [session_id 生成]
    end

    HTTP-->>FE: 302 /home + Set-Cookie session_id [HttpOnly]

    Note over FE: ホーム画面表示
    FE->>RPC: AuthService.GetMe() [Cookie送信]
    RPC->>DB: sessions検証、users取得
    RPC-->>FE: ユーザー情報返却

    Note over FE: テナント選択/参加画面へ
    FE->>RPC: TenantService.GetMyTenants()
    RPC->>DB: SELECT tenant_memberships WHERE user_id
    RPC-->>FE: 参加済みテナント一覧

    Note over FE: ログアウト
    FE->>RPC: AuthService.Logout()
    RPC->>DB: UPDATE sessions SET revoked=true
    RPC-->>FE: Set-Cookie session_id empty [Max-Age 0]
```

### 1.2 認証開始処理詳細

#### エンドポイント: `GET /auth/google/login`

**処理フロー**:

1. **セキュリティパラメータ生成**
   ```go
   state := generateRandomString(32)        // CSRF対策
   codeVerifier := generateCodeVerifier()   // PKCE用（43-128文字）
   codeChallenge := sha256(codeVerifier)    // S256方式
   nonce := generateRandomString(32)        // リプレイ攻撃対策
   ```

2. **一時保存**
   ```sql
   INSERT INTO oauth_states (state, code_verifier, nonce, created_at)
   VALUES ($1, $2, $3, NOW())
   ```

3. **Google認証URLパラメータ**
   ```
   https://accounts.google.com/o/oauth2/v2/auth?
     client_id={CLIENT_ID}
     &redirect_uri={REDIRECT_URI}
     &response_type=code
     &scope=openid+email+profile
     &state={state}
     &nonce={nonce}
     &code_challenge={codeChallenge}
     &code_challenge_method=S256
   ```

### 1.3 コールバック処理詳細

#### エンドポイント: `GET /auth/google/callback`

**処理フロー**:

1. **State検証**
   ```sql
   SELECT * FROM oauth_states
   WHERE state = $1
   AND consumed_at IS NULL
   AND created_at > NOW() - INTERVAL '15 minutes'
   ```

2. **トークン交換**
   ```
   POST https://oauth2.googleapis.com/token
   {
     "code": "{authorization_code}",
     "client_id": "{CLIENT_ID}",
     "client_secret": "{CLIENT_SECRET}",
     "redirect_uri": "{REDIRECT_URI}",
     "grant_type": "authorization_code",
     "code_verifier": "{code_verifier}"
   }
   ```

3. **IDトークン検証**
   - **署名検証**: Google JWKSを使用
   - **Issuer確認**: `https://accounts.google.com`
   - **Audience確認**: CLIENT_IDと一致
   - **有効期限確認**: 現在時刻より未来
   - **Nonce確認**: 保存されたnonceと一致

4. **ユーザー処理**
   ```sql
   -- users UPSERT
   INSERT INTO users (email, name, icon, created_at, updated_at)
   VALUES ($1, $2, $3, NOW(), NOW())
   ON CONFLICT (email)
   DO UPDATE SET
     name = EXCLUDED.name,
     icon = EXCLUDED.icon,
     updated_at = NOW()
   RETURNING id;

   -- user_identities UPSERT
   INSERT INTO user_identities (user_id, provider, provider_sub, created_at, updated_at)
   VALUES ($1, 'google', $2, NOW(), NOW())
   ON CONFLICT (user_id, provider, provider_sub)
   DO UPDATE SET updated_at = NOW();
   ```

5. **セッション作成**
   ```sql
   INSERT INTO sessions (session_id, user_id, expires_at, created_at)
   VALUES ($1, $2, NOW() + INTERVAL '7 days', NOW())
   ```

6. **Cookie設定**
   ```
   Set-Cookie: session_id={session_id};
     HttpOnly;
     Secure;
     SameSite=Lax;
     Max-Age=604800;
     Path=/
   ```

### 1.4 セッション管理

#### セッション検証（全APIリクエスト）

```go
func validateSession(sessionID string) (*User, error) {
    session, err := db.Query(`
        SELECT s.*, u.*
        FROM sessions s
        JOIN users u ON s.user_id = u.id
        WHERE s.session_id = $1
        AND s.revoked = false
        AND s.expires_at > NOW()
    `, sessionID)

    if err != nil {
        return nil, ErrInvalidSession
    }

    return session.User, nil
}
```

#### セッション更新（アクティビティ延長）

```sql
UPDATE sessions
SET expires_at = NOW() + INTERVAL '7 days'
WHERE session_id = $1
AND revoked = false
```

---

## 2. Console側認証フロー

### 2.1 Organization ID認証フロー

```mermaid
sequenceDiagram
    autonumber
    participant Admin as Console Admin
    participant Console as Console App
    participant DB as PostgreSQL

    Note over Admin: Console初回アクセス
    Admin->>Console: アクセス
    Console-->>Admin: Organization ID入力画面
    Admin->>Console: Organization ID/Key入力
    Console->>Console: 定数と照合

    alt 認証成功
        Console->>DB: INSERT console_sessions
        Console->>Console: JWTトークン生成
        Console-->>Admin: Console Dashboard + JWT Cookie
    else 認証失敗
        Console-->>Admin: エラー表示
    end
```

### 2.2 認証処理詳細

#### 定数管理（現在実装）

```go
const (
    ORGANIZATION_ID  = "ORG-DEFAULT-001"
    ORGANIZATION_KEY = "org_key_example_12345"
)

// 本番環境では環境変数から取得
func getOrganizationCredentials() (string, string) {
    if env == "production" {
        return os.Getenv("ORGANIZATION_ID"), os.Getenv("ORGANIZATION_KEY")
    }
    return ORGANIZATION_ID, ORGANIZATION_KEY
}
```

#### JWT生成

```go
type ConsoleClaims struct {
    OrganizationID string `json:"org_id"`
    SessionID      string `json:"session_id"`
    jwt.RegisteredClaims
}

func generateConsoleJWT(orgID, sessionID string) (string, error) {
    claims := ConsoleClaims{
        OrganizationID: orgID,
        SessionID:      sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   orgID,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
```

#### セッション管理

```sql
-- セッション作成
INSERT INTO console_sessions (session_id, organization_id, expires_at, created_at)
VALUES ($1, $2, NOW() + INTERVAL '24 hours', NOW());

-- セッション検証
SELECT * FROM console_sessions
WHERE session_id = $1
AND organization_id = $2
AND expires_at > NOW();
```

---

## 3. Tenant参加フロー（改訂版）

### 3.1 UIフローと画面遷移

#### 初回ログイン時のフロー

```
1. Google認証完了
   ↓
2. /home へリダイレクト（ホーム画面）
   ↓
3. テナント選択画面
   - 参加済みテナント: なし
   - 選択肢:
     a) 公開テナントから選択
     b) 参加コードを入力
   ↓
4. テナント参加
   ↓
5. アプリケーションメイン画面へ
```

#### 2回目以降のログイン時のフロー

```
1. Google認証完了
   ↓
2. /home へリダイレクト（ホーム画面）
   ↓
3. テナント選択画面
   - 参加済みテナント一覧を表示
   - 最後に使用したテナントがデフォルト選択
   - オプション:
     a) 既存テナントを選択
     b) 新規テナント追加
     c) テナントの切り替え
   ↓
4. アプリケーションメイン画面へ
```

### 3.2 Tenant参加完全フロー（新設計）

```mermaid
sequenceDiagram
    autonumber
    participant AFE as App Frontend
    participant ARPC as App RPC
    participant DB as Shared PostgreSQL

    Note over AFE: Google認証後、ホーム画面へリダイレクト
    AFE->>ARPC: AuthService.GetMe
    ARPC->>DB: SELECT users and sessions
    ARPC-->>AFE: user

    Note over AFE: テナント管理画面表示
    AFE->>ARPC: TenantService.GetMyTenants
    ARPC->>DB: SELECT tenant_memberships WHERE user_id AND status='active'
    ARPC-->>AFE: 参加済みTenant一覧

    alt 参加済みテナントがある場合
        Note over AFE: 参加済みテナント一覧を表示
        AFE->>AFE: テナント選択UI表示

        alt ユーザーが既存テナントを選択
            AFE->>ARPC: TenantService.SetActiveTenant with membership_id
            ARPC->>DB: UPDATE sessions SET active_membership_id WHERE session_id
            ARPC-->>AFE: OK
            AFE-->>AFE: アプリケーション画面へ遷移
        else 新規テナント追加
            AFE-->>AFE: テナント追加画面へ
        end
    else 参加済みテナントがない場合
        AFE->>ARPC: TenantService.ListAvailableTenants
        ARPC->>DB: SELECT tenants WHERE organization_id = 'ORG-DEFAULT-001'
        ARPC-->>AFE: 利用可能なTenant一覧

        AFE-->>AFE: テナント選択/参加画面表示
    end

    alt 新規テナント参加
        alt 公開テナントから選択
            AFE->>ARPC: TenantService.JoinTenant with tenant_id
            ARPC->>DB: INSERT tenant_memberships (user_id, tenant_id, role='member', status='active')
            ARPC->>DB: UPDATE sessions SET active_membership_id = membership_id WHERE session_id
            ARPC-->>AFE: membership_id
        else 参加コード使用
            AFE->>ARPC: TenantService.JoinByCode with code
            ARPC->>DB: SELECT tenant_join_codes WHERE code and valid
            ARPC->>DB: INSERT tenant_memberships
            ARPC->>DB: UPDATE tenant_join_codes SET used_count = used_count + 1
            ARPC->>DB: UPDATE sessions SET active_membership_id
            ARPC-->>AFE: membership_id
        end

        AFE->>ARPC: TenantService.GetMyTenants
        ARPC->>DB: SELECT updated tenant_memberships
        ARPC-->>AFE: 更新された参加済みTenant一覧
    end
```

### 3.2 複数テナント管理

#### ユーザーの参加済みテナント取得

```sql
SELECT
    tm.id as membership_id,
    tm.tenant_id,
    tm.role,
    tm.status,
    tm.joined_via,
    tm.created_at as joined_at,
    t.name as tenant_name,
    t.slug as tenant_slug,
    t.description,
    t.tenant_type,
    (tm.id = s.active_membership_id) as is_active
FROM tenant_memberships tm
JOIN tenants t ON tm.tenant_id = t.id
LEFT JOIN sessions s ON s.user_id = tm.user_id AND s.session_id = $1
WHERE tm.user_id = $2
AND tm.status = 'active'
ORDER BY tm.created_at DESC;
```

#### アクティブテナントの切り替え

```sql
-- 切り替え権限の確認
SELECT id FROM tenant_memberships
WHERE id = $1
AND user_id = $2
AND status = 'active';

-- セッションのアクティブテナント更新
UPDATE sessions
SET active_membership_id = $1
WHERE session_id = $2
AND user_id = $3;
```

### 3.3 Tenant選択処理

#### 利用可能Tenant取得

```sql
SELECT
    t.id,
    t.name,
    t.slug,
    t.description,
    t.tenant_type,
    COUNT(tm.id) as member_count
FROM tenants t
LEFT JOIN tenant_memberships tm ON t.id = tm.tenant_id AND tm.status = 'active'
WHERE t.organization_id = $1
GROUP BY t.id
ORDER BY t.name;
```

#### Tenant参加処理

```sql
-- メンバーシップ作成
INSERT INTO tenant_memberships (
    tenant_id,
    user_id,
    role,
    status,
    joined_at
) VALUES (
    $1,  -- tenant_id
    $2,  -- user_id
    'member',
    'active',
    NOW()
) ON CONFLICT (tenant_id, user_id) DO NOTHING
RETURNING id;

-- アクティブメンバーシップ設定
UPDATE sessions
SET active_membership_id = $1
WHERE session_id = $2
AND user_id = $3;
```

### 3.4 参加コード処理

#### コード検証

```sql
SELECT
    tjc.*,
    t.name as tenant_name
FROM tenant_join_codes tjc
JOIN tenants t ON tjc.tenant_id = t.id
WHERE tjc.code = $1
AND (tjc.expires_at IS NULL OR tjc.expires_at > NOW())
AND (tjc.max_uses = 0 OR tjc.used_count < tjc.max_uses);
```

#### 使用回数更新

```sql
-- トランザクション内で実行
BEGIN;

-- 参加コードの使用回数チェック
SELECT * FROM tenant_join_codes
WHERE code = $1
FOR UPDATE;

-- メンバーシップ作成
INSERT INTO tenant_memberships (...);

-- 使用回数インクリメント
UPDATE tenant_join_codes
SET used_count = used_count + 1
WHERE code = $1;

COMMIT;
```

### 3.5 Tenant切り替え

#### アクティブTenant変更

```sql
-- 権限確認
SELECT id FROM tenant_memberships
WHERE user_id = $1
AND tenant_id = $2
AND status = 'active';

-- 切り替え
UPDATE sessions
SET active_membership_id = $1
WHERE session_id = $2
AND user_id = $3;
```

#### コンテキスト取得（API処理）

```go
func GetTenantContext(ctx context.Context) (*TenantContext, error) {
    session := GetSessionFromContext(ctx)

    if session.ActiveMembershipID == nil {
        return nil, ErrNoActiveTenant
    }

    membership, err := db.GetMembership(session.ActiveMembershipID)
    if err != nil {
        return nil, err
    }

    return &TenantContext{
        TenantID: membership.TenantID,
        UserID:   membership.UserID,
        Role:     membership.Role,
    }, nil
}
```

## セキュリティ考慮事項

### CSRF対策
- OAuth: stateパラメータ（32文字のランダム文字列）
- フォーム: csrf_tokenをセッションに保存

### セッション固定攻撃対策
- ログイン成功時に新しいセッションIDを発行
- 古いセッションは無効化

### タイミング攻撃対策
- パスワード検証（将来）: constant-time comparison
- エラーレスポンスの統一化

### レート制限
- ログイン試行: 5回/分
- 参加コード試行: 10回/時
- API全般: 1000回/時

---

## 4. 実装ガイドライン

### 4.1 フロントエンド実装

#### 必要なページ/コンポーネント

1. **ホーム画面 (/home)**
   - Google認証後の初期ランディングページ
   - ユーザー情報の表示
   - テナント管理への導線

2. **テナント選択画面 (/tenants)**
   - 参加済みテナント一覧（カード or リスト形式）
   - 最後に使用したテナントのハイライト
   - 「新規テナント追加」ボタン
   - テナント切り替え機能

3. **テナント参加画面 (/tenants/join)**
   - 公開テナント一覧
   - 参加コード入力フォーム
   - 参加確認ダイアログ

#### 状態管理

```typescript
interface TenantState {
  myTenants: Tenant[]           // 参加済みテナント
  activeTenant: Tenant | null    // 現在アクティブなテナント
  availableTenants: Tenant[]     // 参加可能な公開テナント
}

interface Tenant {
  membershipId: string
  tenantId: string
  name: string
  slug: string
  description: string
  role: 'admin' | 'member'
  isActive: boolean
}
```

### 4.2 バックエンド実装

#### 必要なAPIエンドポイント

1. **TenantService.GetMyTenants**
   - 参加済みテナント一覧取得
   - アクティブテナントのマーキング

2. **TenantService.ListAvailableTenants**
   - 公開テナント一覧取得
   - 参加可否の判定

3. **TenantService.JoinTenant**
   - テナントへの新規参加
   - membership作成

4. **TenantService.JoinByCode**
   - 参加コードによる参加
   - コード検証と使用回数更新

5. **TenantService.SetActiveTenant**
   - アクティブテナントの切り替え
   - セッション更新

6. **TenantService.LeaveTenant**
   - テナントからの離脱
   - membership status更新

#### データベース考慮事項

- `tenant_memberships`テーブルのユニーク制約により、同一ユーザーは同じテナントに複数回参加不可
- `sessions.active_membership_id`でアクティブテナントを管理
- テナント切り替え時はセッションの更新のみ（新規セッション作成は不要）

### 4.3 セキュリティ考慮事項

1. **テナント分離**
   - APIレスポンスは現在のactive_membership_idに基づいてフィルタリング
   - テナント間のデータ漏洩防止

2. **権限管理**
   - テナント内での役割（admin/member）に基づくアクセス制御
   - テナント切り替え時の権限再検証

3. **セッション管理**
   - テナント切り替えはセッション継続（再ログイン不要）
   - active_membership_idの整合性チェック

### 4.4 UX改善ポイント

1. **スムーズな遷移**
   - 参加済みテナントが1つの場合は自動選択
   - 最後に使用したテナントの記憶と自動選択オプション

2. **視覚的フィードバック**
   - 現在のアクティブテナントの明示的表示
   - テナント切り替え時のローディング表示

3. **エラーハンドリング**
   - 無効な参加コードの適切なエラーメッセージ
   - テナント満員時の案内