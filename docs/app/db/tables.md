## テーブル詳細とカラム解説

### users

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | ユーザーの主キー。`uuid_generate_v4()`で自動生成。 |
| email | text | Googleから取得するメールアドレス。lowerでユニーク制約。 |
| name | text | Googleプロフィールの表示名。 |
| icon | text | Googleプロフィール画像URL。 |
| created_at | timestamptz | レコード作成時刻。 |
| updated_at | timestamptz | 最終更新時刻。 |

### user_identities

| カラム | 型 | 説明 |
| --- | --- | --- |
| id | uuid | 主キー。 |
| user_id | uuid | users.idへの外部キー。削除時はCASCADE。 |
| provider | text | 連携元プロバイダ。Google固定。 |
| provider_sub | text | Googleの`sub`（一意のユーザーID）。 |
| created_at / updated_at | timestamptz | 作成・更新時刻。 |

備考: 現状は Google API を呼ばないため refresh_token は保存しない。

### sessions

| カラム | 型 | 説明 |
| --- | --- | --- |
| session_id | text | セッションを識別するランダム文字列（クッキーに保存）。 |
| user_id | uuid | users.idへの外部キー。 |
| active_membership_id | uuid | tenant_memberships.idへの外部キー。NULL許容：テナント未選択。 |
| created_at | timestamptz | セッション開始時刻。 |
| expires_at | timestamptz | 期限切れ時刻。GC対象。 |
| csrf_token | text | フォームPOST等のCSRF対策。 |
| revoked | boolean | 手動ログアウトなどで無効化された場合true。 |

### oauth_states

| カラム | 型 | 説明 |
| --- | --- | --- |
| state | text | OAuth CSRF対策のランダム文字列。 |
| code_verifier | text | PKCEのcode_verifier。 |
| nonce | text | IDトークン検証用nonce。 |
| created_at | timestamptz | 作成時刻。通常10〜15分で削除。 |
| consumed_at | timestamptz | 使用済み時刻。再利用防止。 |

