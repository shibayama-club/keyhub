## マルチテナント拡張（console × app） — 要件の要約

- console で `tenant`（学校/企業）を作成し、任意でドメイン（例: `kogakuin.ac.jp`）を紐づけ可能。
- app 側の Google ログイン時、ユーザーのメールドメインから候補テナントを自動提示。ただし加入は任意（スキップ可）。
- ドメインが一意に判別できない（例: `gmail.com`）場合は、console 発行の参加コード（join_code）を app で入力 → 特定テナントに加入。
- ユーザーは複数テナントに所属可。app セッションは active_membership を 1 つ保持（切替可能）。
- console ログインは JWT+Cookie（`tenant_name` と `tenant_password`）。
- app でテナント所属が確定したら、console 側にもユーザーが登録される（同DB or RPCで反映）。
- 将来拡張: `tenant_admins (tenant_id, user_id, role)` を導入可能。現状は最小構成として `tenants.password_hash` のみで運用し、拡張余地を明記。

