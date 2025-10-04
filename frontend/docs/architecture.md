## 全体的なアーキテクチャ

```
frontend/app/src/
├── assets/          # 静的リソース（アイコン、画像等）
├── components/      # 再利用可能なUIコンポーネント
├── hooks/          # カスタムReactフック
├── libs/           # ライブラリ・SDK・ユーティリティ
├── pages/          # ページコンポーネント（ルーティング対応）
├── router/         # ルーティング設定
├── test/           # テスト設定・ユーティリティ
├── types/          # TypeScript型定義
├── utils/          # 汎用ユーティリティ関数
├── App.tsx         # アプリケーションルートコンポーネント
└── main.tsx        # アプリケーションエントリーポイント
```

---

## 技術スタック・設計思想

### アーキテクチャパターン

#### **Presentation Layer パターン**
```
Pages → Components → Hooks → Libs → API
```

#### **状態管理戦略**
- **Local State**: useState, useReducer
- **Global State**: Zustand (認証状態)
- **Server State**: React Query (API キャッシュ)
- **URL State**: React Router

#### **依存性注入**
```tsx
// Context Provider パターン
<AuthProvider>
  <QueryClientProvider>
    {/* アプリケーション */}
  </QueryClientProvider>
</AuthProvider>
```

### セキュリティ設計

#### **多層防御アーキテクチャ**
1. **ルーティングレベル**: AuthGuard によるページアクセス制御
2. **コンポーネントレベル**: 条件付きレンダリング
3. **API レベル**: JWT トークン認証
4. **データレベル**: 入力値サニタイズ

### パフォーマンス最適化

#### **コード分割**
```tsx
// React.lazy によるページレベル分割
const HomeMonth = React.lazy(() => import('./pages/HomeMonth'));
```

#### **バンドル最適化**
- Tree Shaking による不要コード除去
- Dynamic Import によるチャンク分割
- Asset 最適化 (画像圧縮、WebP 対応)

#### **レンダリング最適化**
- React.memo による不要再レンダリング防止
- useMemo, useCallback による値・関数メモ化
- 仮想スクロール (長いリスト表示)

---

## 開発・運用考慮事項

### 開発効率

#### **開発サーバー**
```bash
# 開発サーバー起動
cd frontend/app
pnpm dev

# アクセス: http://localhost:5173/app
```

#### **ホットリロード**
- ファイル変更の即座反映
- 状態保持機能
- エラーオーバーレイ

### デバッグ・監視

#### **開発ツール**
- React DevTools
- React Query DevTools  
- Redux DevTools (認証状態)

#### **エラー監視**
- Sentry によるエラートラッキング
- Console ログによるデバッグ情報
- パフォーマンス監視

### テスト戦略

#### **テストピラミッド**
```
E2E テスト (少数)
    ↑
Integration テスト (中程度)
    ↑  
Unit テスト (多数)
```

#### **テスト対象**
- **Unit**: フック、ユーティリティ関数
- **Integration**: コンポーネント結合テスト
- **E2E**: ユーザーシナリオテスト

---

## 今後の拡張可能性

### 機能拡張

#### **新機能追加パターン**
1. `pages/` に新しいページ追加
2. `router/` にルート定義追加
3. `components/` に専用コンポーネント作成
4. `libs/query.ts` に API 通信追加

---

## 技術スタック詳細

### **Core技術**
- **React 19**: UI ライブラリ
- **TypeScript**: 型安全性
- **Vite**: ビルドツール・開発サーバー

### **状態管理**
- **Zustand**: 軽量グローバル状態管理
- **React Query**: サーバー状態管理・キャッシュ
- **React Router**: URL 状態管理

### **UI/UX**
- **Tailwind CSS**: ユーティリティファースト CSS
- **React Hot Toast**: 通知システム
- **Recharts**: データ視覚化

### **開発・ビルド**
- **pnpm**: パッケージマネージャー
- **ESLint + Prettier**: コード品質
- **Vitest**: テストフレームワーク
