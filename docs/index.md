# KeyHub Documentation Index

## 主要ドキュメント

- [統合ドキュメント](./all_docs.md) - システム全体の統合ドキュメント（マスター）
- [マルチDB設計](./multi_db_architecture.md) - 将来のマルチデータベース設計

## アーキテクチャ

- [システム概要](./architecture/overview.md) - システムアーキテクチャ概要
- [データフロー](./architecture/data_flow.md) - データフローとユースケースシナリオ
- [将来拡張](./architecture/future_extensions.md) - 将来の拡張計画

## データベース

- [スキーマ設計](./database/schema.md) - データベーススキーマ設計
- [ER図](./database/er_diagrams.md) - Entity-Relationship図
- [テーブル仕様](./database/tables.md) - 詳細なテーブルとカラム仕様
- [マイグレーション](./database/migrations.md) - マイグレーション戦略とSQL
- [パフォーマンス](./database/performance.md) - インデックスと最適化

## API

- [App API](./api/app_api.md) - アプリケーション側API仕様
- [Console API](./api/console_api.md) - 管理コンソール側API仕様

## 認証

- [Google OAuth](./authentication/google_oauth.md) - Google OAuth実装
- [認証フロー](./authentication/flows.md) - 完全な認証フロー詳細

## Console

- [管理機能](./console/management.md) - Console管理機能仕様

## セキュリティ

- [セキュリティポリシー](./security/policies.md) - セキュリティ実装方針

## その他

- [コミット規約](./other/commit.md) - Gitコミットメッセージ規約

---

## ドキュメント構成について

すべてのドキュメントは `all_docs.md` をマスターソースとして、各トピックごとに分割されています。

### 現在実装
- 単一PostgreSQLデータベース
- 固定Organization ID（ORG-DEFAULT-001）
- Google OAuth認証
- Tenant/User管理
- Console管理画面

### 将来実装予定
- マルチデータベース構成
- 動的Organization管理
- Groups機能
- 監査ログ
- 高度な権限管理（RBAC）