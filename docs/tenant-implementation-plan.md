# Tenant参加フロー実装計画

## 概要
Google認証後にホーム画面へ遷移し、ユーザーがテナントを選択・管理できる新しいフローを実装します。

## 変更内容サマリー

### 1. 認証フローの変更
- **現在**: Google認証 → `/app`へリダイレクト
- **変更後**: Google認証 → `/home`へリダイレクト → テナント選択 → アプリケーション

### 2. マルチテナント対応
- ユーザーは複数のテナントに参加可能
- テナント間の切り替えが可能
- 最後に使用したテナントの記憶

## フロントエンド実装タスク

### Phase 1: 基本ルーティング設定
- [ ] `/home`ルートの作成
- [ ] `/tenants`ルートの作成
- [ ] `/tenants/join`ルートの作成
- [ ] 認証後のリダイレクト先を`/home`に変更

### Phase 2: ホーム画面実装
- [ ] ユーザー情報表示コンポーネント
- [ ] テナント管理への導線ボタン
- [ ] 初回ユーザー向けのウェルカムメッセージ

### Phase 3: テナント選択画面
- [ ] 参加済みテナント一覧コンポーネント
  - [ ] カード形式のテナント表示
  - [ ] アクティブテナントのハイライト
  - [ ] テナント情報（名前、説明、メンバー数）
- [ ] テナント切り替え機能
- [ ] 新規テナント追加ボタン

### Phase 4: テナント参加画面
- [ ] 公開テナント一覧
  - [ ] 検索・フィルター機能
  - [ ] テナント詳細表示
- [ ] 参加コード入力フォーム
  - [ ] バリデーション
  - [ ] エラー表示
- [ ] 参加確認ダイアログ

### Phase 5: 状態管理
- [ ] TenantContextの実装
- [ ] テナント関連のstate管理
- [ ] APIクライアントの実装

## バックエンド実装タスク

### Phase 1: データモデル確認
- [x] `tenant_memberships`テーブル構造確認
- [x] `sessions.active_membership_id`の動作確認
- [ ] マイグレーションファイルの準備（必要な場合）

### Phase 2: Repository層
- [ ] `GetUserTenants(userID)` - 参加済みテナント取得
- [ ] `GetAvailableTenants(orgID)` - 公開テナント取得
- [ ] `CreateTenantMembership(userID, tenantID)` - メンバーシップ作成
- [ ] `UpdateActiveMembership(sessionID, membershipID)` - アクティブテナント更新
- [ ] `ValidateJoinCode(code)` - 参加コード検証

### Phase 3: UseCase層
- [ ] `GetMyTenants(ctx, sessionID)` - ユーザーのテナント一覧
- [ ] `ListAvailableTenants(ctx)` - 参加可能テナント一覧
- [ ] `JoinTenant(ctx, sessionID, tenantID)` - テナント参加
- [ ] `JoinByCode(ctx, sessionID, code)` - コードでテナント参加
- [ ] `SetActiveTenant(ctx, sessionID, membershipID)` - アクティブテナント切替
- [ ] `LeaveTenant(ctx, sessionID, membershipID)` - テナント離脱

### Phase 4: RPC Handler実装
- [ ] Proto定義の作成/更新
  ```proto
  service TenantService {
    rpc GetMyTenants(GetMyTenantsRequest) returns (GetMyTenantsResponse);
    rpc ListAvailableTenants(ListAvailableTenantsRequest) returns (ListAvailableTenantsResponse);
    rpc JoinTenant(JoinTenantRequest) returns (JoinTenantResponse);
    rpc JoinByCode(JoinByCodeRequest) returns (JoinByCodeResponse);
    rpc SetActiveTenant(SetActiveTenantRequest) returns (SetActiveTenantResponse);
    rpc LeaveTenant(LeaveTenantRequest) returns (LeaveTenantResponse);
  }
  ```
- [ ] Handler実装
- [ ] インターセプターでのテナントコンテキスト処理

### Phase 5: セキュリティ実装
- [ ] テナント間のデータ分離確認
- [ ] 権限チェックミドルウェア
- [ ] レート制限の実装

## テスト計画

### ユニットテスト
- [ ] Repository層のテスト
- [ ] UseCase層のテスト
- [ ] Handler層のテスト

### 統合テスト
- [ ] 認証フロー全体のE2Eテスト
- [ ] テナント参加フローのテスト
- [ ] テナント切り替えのテスト

### シナリオテスト
1. 新規ユーザーの初回ログインとテナント参加
2. 既存ユーザーのログインと複数テナント管理
3. 参加コードによるテナント参加
4. テナント切り替えとコンテキスト保持

## リリース計画

### Stage 1: Backend API実装（1-2週間）
- Repository/UseCase層の実装
- RPC Handlerの実装
- テスト実装

### Stage 2: Frontend基本実装（1週間）
- ルーティング設定
- 基本画面の実装
- API連携

### Stage 3: UI/UX改善（1週間）
- デザイン適用
- アニメーション追加
- エラーハンドリング改善

### Stage 4: 本番リリース準備
- パフォーマンステスト
- セキュリティレビュー
- ドキュメント更新

## 注意事項

1. **後方互換性**: 既存のセッションは維持し、段階的に移行
2. **データ移行**: 既存ユーザーのデフォルトテナント設定
3. **ロールバック計画**: 問題発生時の切り戻し手順を準備

## 成功指標

- ユーザーがスムーズにテナントを選択・切り替えできる
- 複数テナント参加時のパフォーマンス低下がない
- セキュリティ面でテナント間のデータ分離が完全
- 既存機能への影響がない