このドキュメントは分割されました。最新の内容は以下を参照してください。

- docs/system/README.md
- 認証の詳細: docs/auth/authflow.md

主な分割先
- 全体ワークフロー: docs/system/workflow.md
- ER 図・テーブル詳細: docs/system/db/er-core.md, docs/system/db/tables.md
- API と RPC: docs/system/api.md
- 運用Tips/画面メモ: docs/system/operations.md
- セキュリティ/一意性ポリシー: docs/system/security/policies.md
- RLS 概略: docs/system/security/rls.md
- マルチテナント拡張: docs/system/multitenant/overview.md, docs/system/multitenant/flows.md, docs/system/multitenant/er.md, docs/system/multitenant/ddl.md

---

### 10 RLS（概略）

- コンテキスト: セッションの `active_membership_id` を接続単位で `SET app.membership_id = '...'` として渡す。
- 補助関数（例）:

```sql
CREATE SCHEMA IF NOT EXISTS app;

CREATE OR REPLACE FUNCTION app.current_membership_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('app.membership_id', true)::uuid
$$;

CREATE OR REPLACE FUNCTION app.current_tenant_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT tm.tenant_id
  FROM tenant_memberships tm
  WHERE tm.id = app.current_membership_id()
$$;
```

- ポリシー（例）: テナントに紐づく全テーブルを `tenant_id = app.current_tenant_id()` で制限。

```sql
ALTER TABLE tenants            ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_domains     ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_join_codes  ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_memberships ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_is_current ON tenants
  USING (id = app.current_tenant_id());

CREATE POLICY tenant_is_current ON tenant_domains
  USING (tenant_id = app.current_tenant_id());

CREATE POLICY tenant_is_current ON tenant_join_codes
  USING (tenant_id = app.current_tenant_id());

CREATE POLICY tenant_is_current ON tenant_memberships
  USING (tenant_id = app.current_tenant_id());
```

注意: 実運用では書き込み系（INSERT/UPDATE/DELETE）の `WITH CHECK` も合わせて定義し、管理系（console）とアプリ系（app）でロールを分離する。
