# 技術スタック

KeyHub で使用している主要な技術をまとめたドキュメントです。

## バックエンド

| カテゴリ | 技術 |
| --- | --- |
| 言語 | Go |
| 通信プロトコル | Connect RPC (gRPC 互換, HTTP/2) |
| HTTP サーバ | Echo |
| CLI / 設定 | cobra + viper |
| データベース | PostgreSQL (pgx ドライバ) |
| クエリ生成 | sqlc |
| マイグレーション | goose |
| エラー処理 | cockroachdb/errors |
| モニタリング | Sentry |
| 認証 | Google OAuth, JWT |
| テスト | gomock |
| アーキテクチャ | クリーンアーキテクチャ + DDD(Domain-Driven Design) + RLS (Row-Level Security) によるマルチテナント |

## フロントエンド

| カテゴリ | 技術 |
| --- | --- |
| 言語 | TypeScript |
| フレームワーク | React 19 |
| ビルドツール | Vite |
| スタイリング | TailwindCSS v4 |
| 状態管理 | Zustand |
| データフェッチ | TanStack Query + connect-query |
| ルーティング | react-router-dom |
| パッケージ管理 | pnpm (workspace) |

## スキーマ定義 / コード生成

| カテゴリ | 技術 |
| --- | --- |
| IDL | Protocol Buffers (buf でビルド) |
| 生成物 | Proto から Go / TypeScript クライアントを自動生成 |

## インフラ / 開発ツール

| カテゴリ | 技術 |
| --- | --- |
| コンテナ | Docker (PostgreSQL) |
| イメージビルド | ko |
| オーケストレーション | Kubernetes |
| ツール管理 | aqua (全 CLI ツールをバージョン固定) |
| タスクランナー | Task (Taskfile) |
| Linter / Formatter | golangci-lint, ESLint, gofmt, Prettier |
| スキーマドキュメント | tbls |
