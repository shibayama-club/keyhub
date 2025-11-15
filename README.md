# KeyHub

## 開発環境セットアップ

```bash
# 初期化（依存関係インストール）
task init

# Docker起動
task compose:up
```

## よく使うコマンド

### データベース

```bash
# マイグレーション実行
task migrate:up

# マイグレーション1つだけ実行
task migrate:up-by-one

# マイグレーション巻き戻し
task migrate:down

# 新しいマイグレーションファイル作成
task migrate:new

# シードデータ投入
task seed:up
```

### 開発サーバー起動

```bash
# 全て起動（app + console / frontend + backend）
task run

# app のみ起動
task run:app

# console のみ起動
task run:console

# frontend のみ起動
task run:frontend

# backend のみ起動
task run:backend
```

### コード生成

```bash
# protobuf から Go/TypeScript コード生成
task proto

# sqlc でクエリコード生成
task gen:sqlc

# DBスキーマドキュメント生成
task gen:docs:schema

# 全て実行
task gen
```

### コード品質

```bash
# フォーマット
task fmt

# Lint
task lint
```

### データベース管理

```bash
# スキーマダンプ
task dump:schema

# Docker停止＆削除
task compose:down
```

## ディレクトリ構成

```
keyhub/
├── backend/           # Go バックエンド
│   ├── cmd/          # エントリーポイント (app/console)
│   ├── db/
│   │   ├── migrations/  # DBマイグレーション
│   │   └── seeds/       # シードデータ
│   └── sqlc.yaml
├── frontend/
│   ├── app/          # ユーザー向けフロントエンド
│   └── console/      # 管理画面フロントエンド
├── proto/            # Protocol Buffers定義
├── docs/             # ドキュメント
│   ├── app/         # app関連ドキュメント
│   ├── console/     # console関連ドキュメント
│   └── shared/      # 共通ドキュメント
└── Taskfile.yaml     # タスクランナー設定
```

## 開発時、RowLevelSecurityのかかっているテーブルの中身を見たい時のQuery

```
SET keyhub.organization_id = <組織ID>;
SELECT * FROM public.<テーブル名>;
```

## ドキュメント

- [認証フロー](docs/app/authflow.md)
- [ER図（app側）](docs/app/db/er-core.md)
- [テーブル詳細（app側）](docs/app/db/tables.md)
- [マルチテナントER図](docs/shared/multitenant/er.md)
- [マルチテナントテーブル詳細](docs/shared/multitenant/tables.md)
- [RLS設計](docs/shared/security/rls.md)
- [API仕様](docs/shared/api.md)
