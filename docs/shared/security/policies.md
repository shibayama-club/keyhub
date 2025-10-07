## セキュリティ＆一意性ポリシー

- `tenant_domains.domain` は lowercase 一意に（正規化）。
- 参加コードは十分なエントロピーと期限、可能ならハッシュ保存（生コードは一度だけ表示）。
- `tenant_memberships (tenant_id, user_id)` を UNIQUE にし二重参加を防止。
- `sessions.active_membership_id` の整合性は外部キーで強制（ロジック簡略化）。

