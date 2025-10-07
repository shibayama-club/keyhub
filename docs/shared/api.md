## 4️⃣ API と DB対応 / RPC インターフェース

### API と DB対応

| APIエンドポイント | 操作 | 関連テーブル |
| --- | --- | --- |
| GET /auth/google/login | INSERT | oauth_states |
| GET /auth/google/callback | UPSERT/INSERT | users, user_identities, sessions, oauth_states |
| AuthService.GetMe | SELECT | sessions, users |
| AuthService.Logout | UPDATE (revoked=true) | sessions |

---

### RPCインターフェース（例）

| サービス | メソッド | 目的 |
| --- | --- | --- |
| App.AuthService | GetMe, Logout | 認証状態/ログアウト |
| App.TenantDiscovery | SuggestByEmailDomain(domain) → tenants[] | ドメインから候補取得 |
| App.Membership | JoinByTenantId, JoinByCode, ListMyTenants | 所属/参加コード/所属一覧 |
| App.Session | SetActiveMembership(membership_id) | アクティブメンバーシップ切替 |
| Console.Auth | Login(name, password) | console ログイン（JWT+Cookie） |
| Console.TenantService | CreateTenant, AddDomain, RemoveDomain, GenerateJoinCode, ListMembers | テナント管理 |
