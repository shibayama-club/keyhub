# データベースパフォーマンス最適化

KeyHubのデータベースパフォーマンス最適化戦略とインデックス設計です。

## 目次

1. [インデックス戦略](#1-インデックス戦略)
2. [クエリ最適化](#2-クエリ最適化)
3. [パフォーマンスモニタリング](#3-パフォーマンスモニタリング)
4. [スケーリング戦略](#4-スケーリング戦略)

---

## 1. インデックス戦略

### 1.1 現在のインデックス一覧

**Users Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- 単一カラムインデックス
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- 用途:
-- idx_users_email: ログイン時のメールアドレス検索
-- idx_users_created_at: 管理画面での新規ユーザー一覧
```

**User Identities Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (user_id, provider, provider_sub)

-- インデックス
CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider_sub ON user_identities(provider, provider_sub);

-- 用途:
-- idx_user_identities_user_id: ユーザー削除時の関連データ取得
-- idx_user_identities_provider_sub: OAuth認証時のアイデンティティ検索
```

**Tenants Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (organization_id, slug)

-- インデックス
CREATE INDEX idx_tenants_organization_id ON tenants(organization_id);
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_is_default ON tenants(is_default) WHERE is_default = true;

-- 用途:
-- idx_tenants_organization_id: 組織別Tenant一覧取得
-- idx_tenants_slug: URLアクセス時のTenant解決
-- idx_tenants_is_default: デフォルトTenant取得（部分インデックス）
```

**Tenant Memberships Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (tenant_id, user_id)

-- インデックス
CREATE INDEX idx_tenant_memberships_tenant_id ON tenant_memberships(tenant_id);
CREATE INDEX idx_tenant_memberships_user_id ON tenant_memberships(user_id);
CREATE INDEX idx_tenant_memberships_status ON tenant_memberships(status);
CREATE INDEX idx_tenant_memberships_tenant_user_status
    ON tenant_memberships(tenant_id, user_id, status);

-- 用途:
-- idx_tenant_memberships_tenant_id: Tenantメンバー一覧
-- idx_tenant_memberships_user_id: ユーザーの所属Tenant一覧
-- idx_tenant_memberships_status: アクティブメンバーのフィルタリング
-- idx_tenant_memberships_tenant_user_status: 複合条件検索の最適化
```

**Sessions Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (session_id)

-- インデックス
CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked_expires
    ON sessions(revoked, expires_at) WHERE revoked = false;

-- 用途:
-- idx_sessions_session_id: セッション検証の高速化
-- idx_sessions_user_id: ユーザーのセッション管理
-- idx_sessions_expires_at: 期限切れセッションのクリーンアップ
-- idx_sessions_revoked_expires: アクティブセッションの効率的な取得
```

**OAuth States Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (state)

-- インデックス
CREATE INDEX idx_oauth_states_state ON oauth_states(state);
CREATE INDEX idx_oauth_states_created_at ON oauth_states(created_at);
CREATE INDEX idx_oauth_states_consumed
    ON oauth_states(consumed_at) WHERE consumed_at IS NULL;

-- 用途:
-- idx_oauth_states_state: State検証の高速化
-- idx_oauth_states_created_at: 古いStateのクリーンアップ
-- idx_oauth_states_consumed: 未使用Stateの効率的な取得
```

**Tenant Join Codes Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (code)

-- インデックス
CREATE INDEX idx_tenant_join_codes_code ON tenant_join_codes(code);
CREATE INDEX idx_tenant_join_codes_tenant_id ON tenant_join_codes(tenant_id);
CREATE INDEX idx_tenant_join_codes_expires_at ON tenant_join_codes(expires_at);
CREATE INDEX idx_tenant_join_codes_active
    ON tenant_join_codes(code, expires_at)
    WHERE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP;

-- 用途:
-- idx_tenant_join_codes_code: 参加コード検証
-- idx_tenant_join_codes_tenant_id: Tenant別コード管理
-- idx_tenant_join_codes_expires_at: 期限切れコードのクリーンアップ
-- idx_tenant_join_codes_active: 有効なコードの高速検索
```

**Console Sessions Table**
```sql
-- 主キー
PRIMARY KEY (id)

-- ユニーク制約
UNIQUE (session_id)

-- インデックス
CREATE INDEX idx_console_sessions_session_id ON console_sessions(session_id);
CREATE INDEX idx_console_sessions_expires_at ON console_sessions(expires_at);
CREATE INDEX idx_console_sessions_active
    ON console_sessions(expires_at) WHERE expires_at > CURRENT_TIMESTAMP;

-- 用途:
-- idx_console_sessions_session_id: セッション検証
-- idx_console_sessions_expires_at: 期限切れセッションのクリーンアップ
-- idx_console_sessions_active: アクティブセッションの効率的な取得
```

### 1.2 インデックス設計原則

**カーディナリティの考慮**
```sql
-- 高カーディナリティ（良いインデックス候補）
CREATE INDEX ON users(email);  -- ほぼユニーク

-- 低カーディナリティ（単体では非効率）
CREATE INDEX ON tenant_memberships(status);  -- 3つの値のみ

-- 複合インデックスで改善
CREATE INDEX ON tenant_memberships(status, joined_at DESC);
```

**カバリングインデックス**
```sql
-- クエリ
SELECT user_id, role, joined_at
FROM tenant_memberships
WHERE tenant_id = $1 AND status = 'active';

-- カバリングインデックス（全データをインデックスから取得）
CREATE INDEX idx_covering_membership
ON tenant_memberships(tenant_id, status)
INCLUDE (user_id, role, joined_at);
```

**部分インデックス**
```sql
-- アクティブレコードのみインデックス化
CREATE INDEX idx_active_sessions
ON sessions(user_id, expires_at)
WHERE revoked = false;

-- 特定条件のみインデックス化
CREATE INDEX idx_owner_memberships
ON tenant_memberships(tenant_id, user_id)
WHERE role = 'owner';
```

---

## 2. クエリ最適化

### 2.1 頻出クエリの最適化

**ユーザーのTenant一覧取得**
```sql
-- 最適化前
SELECT t.*, tm.role
FROM tenants t
JOIN tenant_memberships tm ON t.id = tm.tenant_id
WHERE tm.user_id = $1 AND tm.status = 'active';

-- 最適化後（インデックス活用）
SELECT
    t.id, t.name, t.slug, tm.role,
    t.tenant_type, t.is_default
FROM tenant_memberships tm
JOIN tenants t ON tm.tenant_id = t.id
WHERE tm.user_id = $1
AND tm.status = 'active'
ORDER BY t.is_default DESC, t.name;

-- 実行計画
/*
Nested Loop
  -> Index Scan using idx_tenant_memberships_user_id on tenant_memberships tm
       Filter: (status = 'active')
  -> Index Scan using tenants_pkey on tenants t
       Index Cond: (id = tm.tenant_id)
*/
```

**セッション検証**
```sql
-- 最適化前
SELECT * FROM sessions
WHERE session_id = $1;

-- 最適化後（必要カラムのみ取得）
SELECT
    s.user_id,
    s.active_membership_id,
    s.expires_at,
    s.revoked
FROM sessions s
WHERE s.session_id = $1
AND s.revoked = false
AND s.expires_at > NOW();

-- さらに最適化（JOINを含む）
SELECT
    s.user_id,
    u.email,
    u.name,
    tm.tenant_id,
    tm.role
FROM sessions s
JOIN users u ON s.user_id = u.id
LEFT JOIN tenant_memberships tm ON s.active_membership_id = tm.id
WHERE s.session_id = $1
AND s.revoked = false
AND s.expires_at > NOW();
```

### 2.2 バッチクエリ最適化

**N+1問題の回避**
```go
// 悪い例：N+1クエリ
tenants := getTenants()
for _, tenant := range tenants {
    memberCount := getMemberCount(tenant.ID)  // N回のクエリ
    tenant.MemberCount = memberCount
}

// 良い例：一括取得
query := `
    SELECT
        t.*,
        COUNT(tm.id) as member_count
    FROM tenants t
    LEFT JOIN tenant_memberships tm
        ON t.id = tm.tenant_id
        AND tm.status = 'active'
    WHERE t.organization_id = $1
    GROUP BY t.id
`
```

**バルクインサート**
```sql
-- 悪い例：個別INSERT
INSERT INTO tenant_memberships (tenant_id, user_id, role) VALUES ($1, $2, $3);
INSERT INTO tenant_memberships (tenant_id, user_id, role) VALUES ($4, $5, $6);

-- 良い例：バルクINSERT
INSERT INTO tenant_memberships (tenant_id, user_id, role) VALUES
    ($1, $2, $3),
    ($4, $5, $6),
    ($7, $8, $9)
ON CONFLICT (tenant_id, user_id) DO NOTHING;
```

### 2.3 集計クエリ最適化

**統計情報の事前計算**
```sql
-- マテリアライズドビュー作成
CREATE MATERIALIZED VIEW mv_tenant_statistics AS
SELECT
    t.id as tenant_id,
    t.organization_id,
    COUNT(DISTINCT tm.user_id) as total_members,
    COUNT(DISTINCT CASE WHEN tm.role = 'owner' THEN tm.user_id END) as owner_count,
    COUNT(DISTINCT CASE WHEN tm.role = 'admin' THEN tm.user_id END) as admin_count,
    COUNT(DISTINCT CASE WHEN tm.joined_at > NOW() - INTERVAL '7 days' THEN tm.user_id END) as new_members_week,
    MAX(tm.joined_at) as last_join_date
FROM tenants t
LEFT JOIN tenant_memberships tm ON t.id = tm.tenant_id AND tm.status = 'active'
GROUP BY t.id, t.organization_id;

-- インデックス追加
CREATE UNIQUE INDEX ON mv_tenant_statistics(tenant_id);
CREATE INDEX ON mv_tenant_statistics(organization_id);

-- 定期更新ジョブ
CREATE OR REPLACE FUNCTION refresh_tenant_statistics()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_tenant_statistics;
END;
$$ LANGUAGE plpgsql;

-- cronジョブで15分ごとに実行
SELECT cron.schedule('refresh-tenant-stats', '*/15 * * * *', 'SELECT refresh_tenant_statistics()');
```

---

## 3. パフォーマンスモニタリング

### 3.1 スロークエリ検出

**PostgreSQL設定**
```ini
# postgresql.conf
log_min_duration_statement = 100  # 100ms以上のクエリをログ記録
log_statement = 'all'            # 全SQLをログ記録（開発環境）
log_duration = on                 # 実行時間をログ記録
log_lock_waits = on              # ロック待機をログ記録
log_temp_files = 0               # 一時ファイル使用をログ記録
```

**スロークエリ分析**
```sql
-- pg_stat_statements有効化
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- スロークエリTOP 10
SELECT
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    stddev_exec_time,
    rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- 頻繁に実行される重いクエリ
SELECT
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    (total_exec_time / sum(total_exec_time) OVER ()) * 100 as percentage
FROM pg_stat_statements
WHERE calls > 100
ORDER BY total_exec_time DESC
LIMIT 20;
```

### 3.2 インデックス使用状況

```sql
-- インデックス使用率
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- 未使用インデックスの検出
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
AND indexname NOT LIKE '%_pkey'
ORDER BY pg_relation_size(indexrelid) DESC;

-- インデックスブロート検出
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
    ROUND(100 * (pg_relation_size(indexrelid) - pg_relation_size(indexrelid::regclass))
          / pg_relation_size(indexrelid)::numeric, 2) AS bloat_percentage
FROM pg_stat_user_indexes
WHERE pg_relation_size(indexrelid) > 1000000  -- 1MB以上
ORDER BY bloat_percentage DESC;
```

### 3.3 テーブル統計

```sql
-- テーブルアクセス統計
SELECT
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    idx_tup_fetch,
    n_tup_ins,
    n_tup_upd,
    n_tup_del,
    n_live_tup,
    n_dead_tup,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
ORDER BY seq_scan + idx_scan DESC;

-- デッドタプル率の高いテーブル
SELECT
    schemaname,
    tablename,
    n_live_tup,
    n_dead_tup,
    ROUND(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_tuple_percentage
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY dead_tuple_percentage DESC;
```

### 3.4 リアルタイムモニタリング

```go
// Prometheusメトリクス
type DBMetrics struct {
    QueryDuration   *prometheus.HistogramVec
    ConnectionPool  *prometheus.GaugeVec
    SlowQueries     *prometheus.CounterVec
    CacheHitRate    *prometheus.GaugeVec
}

func InitDBMetrics() *DBMetrics {
    return &DBMetrics{
        QueryDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "db_query_duration_seconds",
                Help:    "Database query duration",
                Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
            },
            []string{"query_type", "table"},
        ),
        ConnectionPool: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "db_connection_pool_status",
                Help: "Database connection pool status",
            },
            []string{"status"}, // active, idle, total
        ),
    }
}

// クエリ実行時間の記録
func (db *DB) ExecuteWithMetrics(query string, args ...interface{}) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        metrics.QueryDuration.WithLabelValues(getQueryType(query), getTable(query)).Observe(duration)

        if duration > 0.1 { // 100ms以上
            metrics.SlowQueries.WithLabelValues(getQueryType(query)).Inc()
        }
    }()

    return db.Query(query, args...)
}
```

---

## 4. スケーリング戦略

### 4.1 垂直スケーリング

**接続プール最適化**
```go
// 環境別設定
type PoolConfig struct {
    Environment string
    MaxOpenConns int
    MaxIdleConns int
    ConnMaxLifetime time.Duration
}

var poolConfigs = map[string]PoolConfig{
    "development": {
        MaxOpenConns: 10,
        MaxIdleConns: 5,
        ConnMaxLifetime: 30 * time.Minute,
    },
    "staging": {
        MaxOpenConns: 25,
        MaxIdleConns: 10,
        ConnMaxLifetime: 1 * time.Hour,
    },
    "production": {
        MaxOpenConns: 100,
        MaxIdleConns: 25,
        ConnMaxLifetime: 1 * time.Hour,
    },
}
```

### 4.2 水平スケーリング

**読み取りレプリカ**
```go
type DBCluster struct {
    Primary  *sql.DB
    Replicas []*sql.DB
    current  int32
}

// ラウンドロビンでレプリカ選択
func (c *DBCluster) ReadDB() *sql.DB {
    if len(c.Replicas) == 0 {
        return c.Primary
    }

    index := atomic.AddInt32(&c.current, 1)
    return c.Replicas[index%int32(len(c.Replicas))]
}

// 書き込みは常にプライマリ
func (c *DBCluster) WriteDB() *sql.DB {
    return c.Primary
}

// 使用例
func GetUser(cluster *DBCluster, userID string) (*User, error) {
    db := cluster.ReadDB()  // 読み取りレプリカ使用
    return db.QueryRow("SELECT * FROM users WHERE id = $1", userID)
}
```

### 4.3 キャッシング戦略

**Redisキャッシュ層**
```go
type CacheLayer struct {
    redis *redis.Client
    db    *sql.DB
    ttl   time.Duration
}

func (c *CacheLayer) GetTenant(tenantID string) (*Tenant, error) {
    // キャッシュチェック
    cacheKey := fmt.Sprintf("tenant:%s", tenantID)
    cached, err := c.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var tenant Tenant
        json.Unmarshal([]byte(cached), &tenant)
        return &tenant, nil
    }

    // DBから取得
    tenant, err := c.db.GetTenant(tenantID)
    if err != nil {
        return nil, err
    }

    // キャッシュ保存
    data, _ := json.Marshal(tenant)
    c.redis.Set(ctx, cacheKey, data, c.ttl)

    return tenant, nil
}

// キャッシュ無効化
func (c *CacheLayer) InvalidateTenant(tenantID string) {
    cacheKey := fmt.Sprintf("tenant:%s", tenantID)
    c.redis.Del(ctx, cacheKey)

    // 関連キャッシュも削除
    pattern := fmt.Sprintf("tenant:%s:*", tenantID)
    keys, _ := c.redis.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        c.redis.Del(ctx, keys...)
    }
}
```

### 4.4 パーティショニング

**時系列データのパーティション**
```sql
-- セッションテーブルの月単位パーティション
CREATE TABLE sessions (
    id UUID NOT NULL,
    session_id TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    -- other columns
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- 月単位パーティション作成
CREATE TABLE sessions_2024_01 PARTITION OF sessions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE sessions_2024_02 PARTITION OF sessions
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- 自動パーティション作成関数
CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    start_date date;
    end_date date;
    partition_name text;
BEGIN
    start_date := DATE_TRUNC('month', CURRENT_DATE + INTERVAL '1 month');
    end_date := start_date + INTERVAL '1 month';
    partition_name := 'sessions_' || TO_CHAR(start_date, 'YYYY_MM');

    EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF sessions FOR VALUES FROM (%L) TO (%L)',
        partition_name, start_date, end_date);
END;
$$ LANGUAGE plpgsql;

-- 月次cronジョブ
SELECT cron.schedule('create-partition', '0 0 25 * *', 'SELECT create_monthly_partition()');
```

## チューニングチェックリスト

### 定期メンテナンス
- [ ] VACUUM ANALYZE実行（週次）
- [ ] インデックス再構築（月次）
- [ ] 統計情報更新（日次）
- [ ] スロークエリレビュー（週次）

### パフォーマンス監視
- [ ] クエリ実行時間の監視
- [ ] インデックス使用率の確認
- [ ] デッドタプル率の確認
- [ ] 接続プール使用率の監視

### 最適化対象
- [ ] 頻繁に実行されるクエリ
- [ ] 実行時間の長いクエリ
- [ ] フルテーブルスキャンの削減
- [ ] 不要インデックスの削除