# Row Level Security (RLS) 設計

## 概要

KeyHubのマルチテナント機能では、PostgreSQLの**Row Level Security (RLS)** を使用してテナント間のデータ分離を実現します。各データベース接続は特定のテナントコンテキストを持ち、そのテナントに属するデータのみにアクセスできます。

---

## RLSの仕組み

### 1. セッションコンテキストの設定

アプリケーションは、データベース接続確立後に現在のメンバーシップIDを設定します。

```sql
-- アプリケーション側で実行（リクエストごと）
SET app.membership_id = 'ユーザーのactive_membership_id';
```

**具体例**:
```sql
-- ユーザーAさんが「会社X」テナントに切り替えた場合
SET app.membership_id = '123e4567-e89b-12d3-a456-426614174000';
```

### 2. ヘルパー関数

RLSポリシーで使用する2つのヘルパー関数を`app`スキーマに定義します。

#### `app.current_membership_id()`

現在のセッションに設定されているメンバーシップIDを取得します。

```sql
CREATE OR REPLACE FUNCTION app.current_membership_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('app.membership_id', true)::uuid
$$;
```

**動作**:
- `current_setting('app.membership_id', true)`: PostgreSQLのセッション変数から値を取得
- `true`: 変数が未設定の場合NULLを返す（エラーにしない）
- `::uuid`: 文字列からUUID型にキャスト

#### `app.current_tenant_id()`

現在のメンバーシップIDから所属するテナントIDを取得します。

```sql
CREATE OR REPLACE FUNCTION app.current_tenant_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT tm.tenant_id
  FROM tenant_memberships tm
  WHERE tm.id = app.current_membership_id()
$$;
```

**動作例**:
```
app.current_membership_id() = '123e4567-...'
  ↓
tenant_memberships テーブルを検索
  ↓
tenant_id = '987e4567-...' (会社Xのテナント)
```

---

## 3. RLSポリシーの適用

### テーブルへのRLS有効化

```sql
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_domains ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_join_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_memberships ENABLE ROW LEVEL SECURITY;
```

### ポリシー定義

各テーブルに`tenant_is_current`という名前のポリシーを作成します。

#### `tenants`テーブル

```sql
CREATE POLICY tenant_is_current ON tenants
  USING (id = app.current_tenant_id());
```

**意味**: 現在のテナントIDと一致するテナントレコードのみ表示

#### `tenant_domains`テーブル

```sql
CREATE POLICY tenant_is_current ON tenant_domains
  USING (tenant_id = app.current_tenant_id());
```

**意味**: 現在のテナントに属するドメインのみ表示

#### `tenant_join_codes`テーブル

```sql
CREATE POLICY tenant_is_current ON tenant_join_codes
  USING (tenant_id = app.current_tenant_id());
```

**意味**: 現在のテナントの招待コードのみ表示

#### `tenant_memberships`テーブル

```sql
CREATE POLICY tenant_is_current ON tenant_memberships
  USING (tenant_id = app.current_tenant_id());
```

**意味**: 現在のテナントのメンバーシップのみ表示

---

## 動作例

### シナリオ: ユーザーAさんのアクセス

```
ユーザーA:
- 会社Xのテナント（tenant_id: aaa）にadminで所属（membership_id: 111）
- 個人プロジェクト（tenant_id: bbb）にmemberで所属（membership_id: 222）
```

#### 例1: 会社Xのコンテキストでクエリ

```sql
-- セッション設定
SET app.membership_id = '111';  -- 会社Xのメンバーシップ

-- テナント情報を取得
SELECT * FROM tenants;
-- 結果: 会社X（tenant_id: aaa）のみ表示

-- ドメイン一覧を取得
SELECT * FROM tenant_domains;
-- 結果: 会社Xに紐づくドメイン（例: @company-x.com）のみ表示

-- メンバー一覧を取得
SELECT * FROM tenant_memberships;
-- 結果: 会社Xのメンバー全員が表示
```

#### 例2: 個人プロジェクトのコンテキストでクエリ

```sql
-- セッション切り替え
SET app.membership_id = '222';  -- 個人プロジェクトのメンバーシップ

-- テナント情報を取得
SELECT * FROM tenants;
-- 結果: 個人プロジェクト（tenant_id: bbb）のみ表示

-- ドメイン一覧を取得
SELECT * FROM tenant_domains;
-- 結果: 個人プロジェクトに紐づくドメインのみ表示
```

#### 例3: コンテキスト未設定の場合

```sql
-- セッション変数なし（またはNULL）
SELECT * FROM tenants;
-- 結果: 0件（RLSポリシーにより全て除外）
```

---

## セキュリティ上の注意点

### 現在のポリシーの制限

**現在のRLSポリシーは読み取り専用（SELECT）のみに適用されます。**

```sql
CREATE POLICY tenant_is_current ON tenants
  USING (id = app.current_tenant_id());  -- SELECT時のみ有効
```

### 書き込み操作への対応（今後の実装推奨）

#### INSERT/UPDATE/DELETE用ポリシー

```sql
-- 書き込み時の制約を追加
CREATE POLICY tenant_is_current_write ON tenants
  FOR INSERT WITH CHECK (id = app.current_tenant_id());

CREATE POLICY tenant_is_current_update ON tenants
  FOR UPDATE USING (id = app.current_tenant_id())
  WITH CHECK (id = app.current_tenant_id());

CREATE POLICY tenant_is_current_delete ON tenants
  FOR DELETE USING (id = app.current_tenant_id());
```

---

## ロール分離（推奨設計）

### アプリケーション用ロール

```sql
CREATE ROLE keyhub_app;
-- RLSポリシーが適用される
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO keyhub_app;
```

### 管理用ロール（Console）

```sql
CREATE ROLE keyhub_console;
-- RLSをバイパス（全テナントにアクセス可能）
ALTER ROLE keyhub_console BYPASSRLS;
GRANT ALL ON ALL TABLES IN SCHEMA public TO keyhub_console;
```

---

## マイグレーション情報

RLS関連の設定は以下のマイグレーションファイルで実装されています:

**ファイル**: `backend/db/migrations/20251007061025_add_rls_functions_and_policies.sql`

**実装内容**:
1. `app`スキーマの作成
2. `app.current_membership_id()`関数の定義
3. `app.current_tenant_id()`関数の定義
4. テーブルへのRLS有効化
5. `tenant_is_current`ポリシーの作成

**ロールバック**: Down migrationで全て削除可能

---

## トラブルシューティング

### データが取得できない場合

```sql
-- 現在のコンテキストを確認
SELECT current_setting('app.membership_id', true);
-- NULL の場合: セッション変数が未設定

-- 関数の動作を確認
SELECT app.current_membership_id();
SELECT app.current_tenant_id();
```

### RLSをバイパスしてデバッグ

```sql
-- スーパーユーザーまたはBYPASSRLS権限を持つロールでログイン
SELECT * FROM tenants;  -- 全テナントが表示される
```

---

## 参考リンク

- PostgreSQL公式ドキュメント: [Row Security Policies](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- 関連ドキュメント: [docs/shared/security/policies.md](./policies.md)
