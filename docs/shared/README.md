# System Design Docs (Index)

このインデックスは、docs を App / Console / Shared に整理した上での導線です。

- App
  - workflow: ../app/workflow.md
  - 認証の詳細: ../app/authflow.md
  - ER 図（コア）: ../app/db/er-core.md
  - テーブル詳細（コア）: ../app/db/tables.md
  - 運用Tips・画面/体験メモ: ../app/operations.md

- API（共有）
  - API と DB 対応・RPC 一覧: ./api.md

- セキュリティ（共有）
  - セキュリティ＆一意性ポリシー: ./security/policies.md
  - RLS（概略・DDLサンプル）: ./security/rls.md

- マルチテナント（共有/コンソール連携）
  - 要件の要約（Console 視点）: ../console/overview.md
  - フロー（Console / App / 参加コード・切替）: ../console/flows.md
  - ER 図（マルチテナント）: ./multitenant/er.md
  - テーブル詳細（マルチテナント）: ./multitenant/tables.md
  - 追加 DDL（抜粋）: ../console/ddl.md

補足:
- `sessions.active_membership_id` に統一。
- Logout は `revoked=true` の UPDATE 運用（監査・多端末管理の都合）。
- Google API は呼ばない前提のため、refresh_token は取得・保存しない。

