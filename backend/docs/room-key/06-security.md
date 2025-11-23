# セキュリティ考慮事項

## 1. 認証・認可

### Console認証必須
- ✅ すべてのAPIはConsole認証を必要とする
- ✅ 認証されたユーザーのみがリソースにアクセス可能

### Organization IDの検証
- ✅ リクエスト元の組織IDとリソースの組織IDを照合
- ✅ 異なる組織のリソースへのアクセスを防止
- ✅ ハンドラー層で組織IDの一致を確認

### Row Level Security (RLS)
- ✅ PostgreSQLのRow Level Securityで組織レベルのデータ分離を保証
- ✅ データベースレベルでのマルチテナント分離
- ✅ アプリケーション層のバグによるデータ漏洩を防止

---

## 2. データバリデーション

### ドメインモデルでの検証
- ✅ すべての値オブジェクト(Value Object)でバリデーションを実施
- ✅ 不正な値の生成を防止
- ✅ 型安全性の確保

**検証項目例:**
- RoomName: 1〜100文字
- RoomLocation: 1〜200文字
- RoomCapacity: 0より大きい整数、またはNULL
- KeyNumber: 1〜50文字
- 各ENUM値の妥当性

### Protobufバリデーション
- ✅ buf validateでリクエストパラメータを検証
- ✅ API境界での入力検証
- ✅ 不正なリクエストの早期検出

**検証ルール:**
- UUID形式の検証: `[(buf.validate.field).string.uuid = true]`
- 文字列長の検証: `{min_len: 1, max_len: 100}`
- 数値範囲の検証: `{gt: 0}`

### SQL制約
- ✅ CHECK制約でデータベースレベルの整合性を保証
- ✅ NOT NULL制約による必須項目の強制
- ✅ UNIQUE制約による重複防止

**制約例:**
- `room_assignments_date_check`: expires_at IS NULL OR expires_at > assigned_at
- `idx_keys_organization_key_number`: (organization_id, key_number) UNIQUE
- `idx_rooms_organization_name`: (organization_id, name) UNIQUE

---

## 3. トランザクション管理

### ACID保証
- ✅ 複数テーブルの更新はトランザクション内で実行
- ✅ 部分的な更新の防止
- ✅ データの一貫性保証

**トランザクションが必要な操作:**
- 部屋割り当て時の競合チェックと登録
- 鍵のステータス更新とロギング
- 複数のリレーションシップの更新

### 外部キー制約
- ✅ カスケード削除を適切に設定
- ✅ 参照整合性の保証
- ✅ 孤立レコードの防止

**外部キー設定:**
- keys.room_id → rooms.id (ON DELETE CASCADE)
- room_assignments.tenant_id → tenants.id (ON DELETE CASCADE)
- room_assignments.room_id → rooms.id (ON DELETE CASCADE)

---

## 4. 監査ログ

### タイムスタンプ
- ✅ すべてのテーブルに`created_at`と`updated_at`を記録
- ✅ リソースの作成・更新日時の追跡
- ✅ 監査証跡の確保

### Sentryインテグレーション
- ✅ エラー発生時の自動通知
- ✅ スタックトレースの記録
- ✅ エラーの集約と分析

### 将来的な拡張
- 📋 操作ログテーブルの追加（誰が・いつ・何を操作したか）
- 📋 変更履歴の保持（差分ログ）
- 📋 削除操作の監査

---

## 5. データ整合性

### 期間重複チェック
- ✅ 同じ部屋への複数割り当てを防止
- ✅ アプリケーション層での検証ロジック実装
- 📋 将来的にはPostgreSQLのEXCLUDE制約の使用も検討

**チェックロジック:**
```
同じroom_idに対して、期間が重複するroom_assignmentsが存在しないことを確認
WHERE room_id = ?
  AND (expires_at IS NULL OR expires_at > NOW())
```

### ソフトデリート
- 📋 将来的には削除フラグ(`deleted_at`)の導入を検討
- 📋 物理削除ではなく論理削除で履歴を保持
- 📋 誤削除からの復旧を容易にする

### 一意性制約
- ✅ 組織内での名前の重複を防止
- ✅ (organization_id, name) の複合UNIQUE制約
- ✅ (organization_id, key_number) の複合UNIQUE制約

---

## 6. セキュアコーディング

### SQLインジェクション対策
- ✅ SQLCによる型安全なクエリ生成
- ✅ プリペアドステートメントの使用
- ✅ 動的SQLの排除

### XSS対策
- ✅ Protobufによる型安全なシリアライゼーション
- ✅ 入力値のエスケープ不要（バイナリプロトコル）

### 権限の最小化
- ✅ データベースユーザーの権限を必要最小限に設定
- ✅ RLSポリシーによる行レベルのアクセス制御
- ✅ 組織IDによるスコープの制限

---

## まとめ

この機能では、以下のセキュリティ対策を多層的に実装します:

1. **認証・認可**: Console認証 + Organization ID検証 + RLS
2. **データバリデーション**: ドメインモデル + Protobuf + SQL制約
3. **トランザクション管理**: ACID保証 + 外部キー制約
4. **監査ログ**: タイムスタンプ + Sentry
5. **データ整合性**: 重複チェック + 一意性制約

これらの対策により、セキュアで信頼性の高いシステムを実現します。
