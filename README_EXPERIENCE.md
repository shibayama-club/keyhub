## プロダクト概要

- **GitHubリポジトリURL**: https://github.com/shibayama-club/keyhub
- **プロダクト名**: KeyHub
- **概要**: 組織内の鍵管理を効率化するマルチテナント型SaaSアプリケーション。部屋や施設の鍵の貸出・返却を管理し、組織ごとにテナントを作成して利用者やリソースを管理できる。

## 開発の経緯

私の所属する大学では、研究室ごとに貸し出せる鍵が1本のみで、地下1階の守衛室で紙の書類を記入しなければ借りることができない仕組みになっています。
また、大学は高層ビル構造であるため、例えば14階で授業を受けた後に研究室へ向かう場合、地下1階まで降りるか、22階の研究室まで行って初めて「鍵が借りられているかどうか」が分かります。結果的に、鍵がすでに貸し出されていた場合は再び移動が発生し、往復で5〜7分程度の無駄な時間が生じていました。
さらに、貸出・返却時の手続きはすべて紙による記入で行われており、混雑時には並ばなければならないことも多く、学生・管理者双方にとって非効率な運用となっていました。
このような日常的な不便さを解消するため、「鍵の状態を事前に確認でき」「貸出・返却をスマートフォンから完結できる」仕組みを作りたいと考え、keyhubの開発に至りました。

## 担当範囲

バックエンド・フロントエンドの設計・実装を担当。

---

# バックエンド

## 1. ポインタと値渡しの適切な使い分け

**適用箇所**:
- [internal/domain/model/tenant.go](backend/internal/domain/model/tenant.go)
- [internal/domain/model/key.go](backend/internal/domain/model/key.go)
- [internal/domain/model/room.go](backend/internal/domain/model/room.go)

**詳細**:

構造体のサイズが大きくない限り、ポインタを使用しない設計方針を採用しています。

```go
// internal/domain/model/tenant.go:127-135
type Tenant struct {
    ID             TenantID
    OrganizationID OrganizationID
    Name           TenantName
    Description    TenantDescription
    Type           TenantType
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

**こだわりポイント**:

| 観点 | 値渡しのメリット |
|------|-----------------|
| **nilの心配が不要** | ポインタを使わないことで`nil`チェックが不要になり、余計な防御的コードを削減 |
| **メモリ効率の向上** | 小〜中サイズの構造体はスタックに割り当てられ、不要なヒープ割り当てを回避 |
| **GC圧力の軽減** | ヒープ割り当てが減ることでGCの負荷が軽減され、パフォーマンスが安定 |
| **並行処理の安全性** | 値渡しによりデータがコピーされるため、複数のgoroutineから同時にアクセスしても競合が発生しない |
| **不変性の保証** | 関数に渡した値が意図せず変更される（エイリアシング問題）を防止 |

**ファクトリ関数での実装例**:
```go
// internal/domain/model/tenant.go:171-193
func NewTenant(
    organizationID OrganizationID,
    name TenantName,
    description TenantDescription,
    tenantType TenantType,
) (Tenant, error) {  // ポインタではなく値を返す
    now := time.Now()
    tenant := Tenant{
        ID:             TenantID(uuid.New()),
        OrganizationID: organizationID,
        Name:           name,
        Description:    description,
        Type:           tenantType,
        CreatedAt:      now,
        UpdatedAt:      now,
    }

    if err := tenant.Validate(); err != nil {
        return Tenant{}, err  // ゼロ値を返すことでnilの心配なし
    }

    return tenant, nil
}
```

**ポインタを使用する例外ケース**:
- SQLCの自動生成コードではnull許容カラムに対してポインタを使用（[sqlc.yaml:13](backend/sqlc.yaml#L13)で設定）
- DBから取得した値が`NULL`の可能性がある場合のみポインタで表現

---

## 2. 型安全なContext管理

**適用箇所**:
- [internal/domain/context.go](backend/internal/domain/context.go)

**詳細**:

Go 1.18のジェネリクスを活用し、型安全なContext値の管理を実現しています。

```go
// internal/domain/context.go:1-20
package domain

import "context"

type (
    ctxKey[T any] struct{}
)

func WithValue[T any](ctx context.Context, val T) context.Context {
    return context.WithValue(ctx, ctxKey[T]{}, val)
}

func RemoveValue[T any](ctx context.Context) context.Context {
    return context.WithValue(ctx, ctxKey[T]{}, nil)
}

func Value[T any](ctx context.Context) (T, bool) {
    value, ok := ctx.Value(ctxKey[T]{}).(T)
    return value, ok
}
```

**こだわりポイント**:
- **ジェネリクスによる型安全性**: 従来の`context.Value`はキーと値が`any`型で型安全性が低いが、ジェネリクスを使用することで型パラメータ`T`によりコンパイル時に型チェック可能
- **キーの衝突防止**: `ctxKey[T]{}`という空構造体をキーとして使用することで、異なる型同士のキー衝突を完全に防止
- **シンプルなAPI**: `domain.WithValue`、`domain.Value`の2つの関数のみで完結

---

## 3. Contextの応用(contextへのさまざまデータの格納、各層での値の取り出し、自動ロギング)

**適用箇所**:
- [internal/interface/console/v1/interceptor/auth.go](backend/internal/interface/console/v1/interceptor/auth.go)
- [internal/infrastructure/sqlc/driver.go](backend/internal/infrastructure/sqlc/driver.go)
- [internal/domain/logger/logger.go](backend/internal/domain/logger/logger.go)

**詳細**:

Contextを通じて認証情報を各層で共有し、ロギングにも自動的に反映させています。

**認証Interceptorでの設定**:
```go
// internal/interface/console/v1/interceptor/auth.go:39-41
session, err := i.useCase.ValidateSession(ctx, token)
if err != nil {
    return ctx, connect.NewError(connect.CodeUnauthenticated, err)
}
ctx = domain.WithValue(ctx, session.OrganizationID)
```

**DB接続時のRLS設定への活用**:
```go
// internal/infrastructure/sqlc/driver.go:33-44
config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
    if orgID, ok := domain.Value[model.OrganizationID](ctx); ok {
        if _, err := conn.Exec(ctx, "SELECT set_config('keyhub.organization_id', $1, false)", orgID.UUID().String()); err != nil {
            slog.ErrorContext(ctx, "BeforeAcquire: failed to set RLS organization_id", slog.String("error", err.Error()))
            return false
        }
    } else {
        if _, err := conn.Exec(ctx, "RESET keyhub.organization_id"); err != nil {
            slog.ErrorContext(ctx, "BeforeAcquire: failed to reset RLS organization_id", slog.String("error", err.Error()))
            return false
        }
    }
    // ...
}
```

**ロガーでの自動コンテキスト抽出**:
```go
// internal/domain/logger/logger.go:39-62
func (h contextualLoggingHandler) Handle(ctx context.Context, record slog.Record) error {
    requestID, ok := domain.Value[RequestID](ctx)
    if ok {
        record.Add(slog.String("request_id", string(requestID)))
    }

    // app
    userID, ok := domain.Value[model.UserID](ctx)
    if ok {
        record.Add(slog.String("user_id", userID.String()))
    }
    tenantID, ok := domain.Value[model.TenantID](ctx)
    if ok {
        record.Add(slog.String("tenant_id", tenantID.String()))
    }

    // console
    organizationID, ok := domain.Value[model.OrganizationID](ctx)
    if ok {
        record.Add(slog.String("organization_id", organizationID.String()))
    }

    return h.Handler.Handle(ctx, record)
}
```

**こだわりポイント**:
- **横断的関心事の一元管理**: 認証情報をContextに格納することで、各層で明示的にパラメータを渡す必要がなく、コードがクリーンに
- **自動ログエンリッチメント**: カスタムslog.Handlerにより、すべてのログに自動的にユーザーID、テナントID、リクエストIDが付与され、トレーサビリティが向上
- **RLSとの連携**: ContextからDB接続時に自動的にPostgreSQLのセッション変数を設定し、Row Level Securityを有効化

---

## 4. DBのTransaction, pool, bulk, migration

**適用箇所**:
- [internal/infrastructure/sqlc/sqlc.go](backend/internal/infrastructure/sqlc/sqlc.go)
- [internal/infrastructure/sqlc/driver.go](backend/internal/infrastructure/sqlc/driver.go)
- [internal/usecase/console/tenant.go](backend/internal/usecase/console/tenant.go)
- [db/migrations/](backend/db/migrations/)

**詳細**:

**トランザクション管理**:
```go
// internal/infrastructure/sqlc/sqlc.go:31-60
func (r *SqlcRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx repository.Transaction) error) (err error) {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return errors.Wrap(err, "failed to begin transaction")
    }

    committed := false
    defer func() {
        if committed {
            return
        }
        if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
            if err != nil {
                err = errors.Wrapf(err, "failed to rollback transaction: %v", rbErr)
            }
        }
    }()

    if txErr := fn(ctx, &SqlcTransaction{queries: r.queries.WithTx(tx)}); txErr != nil {
        err = errors.Wrap(txErr, "failed to Queries")
        return err
    }

    if err = tx.Commit(ctx); err != nil {
        return errors.Wrap(err, "failed to Commit")
    }
    committed = true

    return nil
}
```

**Usecaseでのトランザクション使用例**:
```go
// internal/usecase/console/tenant.go:64-92
err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
    err = tx.CreateTenant(ctx, repository.CreateTenantArg{
        ID:             tenant.ID,
        OrganizationID: tenant.OrganizationID,
        Name:           tenant.Name,
        Description:    tenant.Description,
        Type:           tenant.Type,
    })
    if err != nil {
        return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant in repository")
    }

    err = tx.CreateTenantJoinCode(ctx, repository.CreateTenantJoinCodeArg{
        ID:        joinCodeEntity.ID,
        TenantID:  joinCodeEntity.TenantID,
        Code:      joinCodeEntity.Code,
        ExpiresAt: joinCodeEntity.ExpiresAt,
        MaxUses:   joinCodeEntity.MaxUses,
        UsedCount: joinCodeEntity.UsedCount,
    })
    if err != nil {
        return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant join code in repository")
    }

    return nil
})
```

**コネクションプールの設定**:
```go
// internal/infrastructure/sqlc/driver.go:17-80
func NewPool(ctx context.Context, cf config.DBConfig) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s",
        cf.Host, cf.Port, cf.User, cf.Password, cf.Database)

    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, errors.Wrap(err, "failed to parse pgx config")
    }

    // BeforeAcquire/AfterReleaseでRLS用セッション変数を管理
    config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool { /* ... */ }
    config.AfterRelease = func(conn *pgx.Conn) bool { /* ... */ }

    pool, err := pgxpool.NewWithConfig(ctx, config)
    // ...
}
```

**マイグレーション**: gooseを使用し、Up/Down両方向のマイグレーションをサポート。

**こだわりポイント**:
- **関数型トランザクションパターン**: `WithTransaction`に関数を渡すパターンにより、ロールバックの保証やcommit漏れを防止
- **コネクションプールフック**: `BeforeAcquire`/`AfterRelease`を活用し、各接続でRLS用のセッション変数を自動設定・クリア
- **エラー時の適切なラップ**: ロールバック失敗時も元のエラー情報を保持

---

## 5. エラーのラッピング、キャプチャー(Sentry)

**適用箇所**:
- [internal/domain/errors/errors.go](backend/internal/domain/errors/errors.go)
- [internal/interface/sentry/sentry.go](backend/internal/interface/sentry/sentry.go)

**詳細**:

**ドメインエラー定義**:
```go
// internal/domain/errors/errors.go:1-31
package errors

import "github.com/cockroachdb/errors"

var (
    ErrValidation    = errors.New("Validation Error")
    ErrNotFound      = errors.New("Not Found Error")
    ErrUnAuthorized  = errors.New("Unauthorized Error")
    ErrInternal      = errors.New("Internal Error")
    ErrAlreadyExists = errors.New("Already Exists Error")
)

func IsValidationError(err error) bool {
    return errors.Is(err, ErrValidation)
}
// ...
```

**Sentry統合Interceptor**:
```go
// internal/interface/sentry/sentry.go:72-88
func (i *ErrorInterceptor) connectError(err error) *connect.Error {
    switch {
    case errors.Is(err, domainerrors.ErrValidation):
        return connect.NewError(connect.CodeInvalidArgument, err)
    case errors.Is(err, domainerrors.ErrNotFound):
        return connect.NewError(connect.CodeNotFound, err)
    case errors.Is(err, domainerrors.ErrUnAuthorized):
        return connect.NewError(connect.CodeUnauthenticated, err)
    case errors.Is(err, domainerrors.ErrAlreadyExists):
        return connect.NewError(connect.CodeAlreadyExists, err)
    case errors.Is(err, domainerrors.ErrInternal):
        return connect.NewError(connect.CodeInternal, err)
    default:
        return connect.NewError(connect.CodeUnknown, err)
    }
}
```

**Sentryへのキャプチャ（サーバーエラーのみ）**:
```go
// internal/interface/sentry/sentry.go:230-243
func shouldCaptureInSentry(code connect.Code) bool {
    switch code {
    case connect.CodeInternal,
        connect.CodeUnknown,
        connect.CodeDataLoss,
        connect.CodeUnimplemented,
        connect.CodeUnavailable:
        return true
    default:
        return false
    }
}
```

**エラーヒント機能による多言語対応**:
```go
// internal/interface/sentry/sentry.go:137-150
hints := errors.FlattenHints(err)
if len(hints) > 0 {
    detail, detailErr := connect.NewErrorDetail(&errdetails.LocalizedMessage{
        Locale:  "ja-JP",
        Message: hints,
    })
    if detailErr != nil {
        slog.ErrorContext(ctx, "failed to create error detail", slog.String("error", detailErr.Error()))
    } else {
        newError.AddDetail(detail)
    }
}
```

**こだわりポイント**:
- **cockroachdb/errorsの活用**: スタックトレース、エラーラッピング、ヒント機能を持つ高機能エラーライブラリを使用
- **エラーヒントによるユーザーフレンドリーなメッセージ**: `errors.WithHint`でユーザー向けの日本語メッセージを添付し、Connect RPCのエラー詳細として送信
- **選択的なSentry送信**: クライアントエラー（バリデーション等）はWarnログのみ、サーバーエラーのみSentryに送信

---

## 6. DDDの実践

**適用箇所**:
- [internal/domain/model/](backend/internal/domain/model/) - エンティティ・値オブジェクト
- [internal/domain/repository/](backend/internal/domain/repository/) - リポジトリインターフェース
- [internal/usecase/](backend/internal/usecase/) - ユースケース

**詳細**:

**値オブジェクトの定義とバリデーション**:
```go
// internal/domain/model/tenant.go:33-62
type TenantName string

func (n TenantName) String() string {
    return string(n)
}

func (n TenantName) Validate() error {
    if n == "" {
        return errors.WithHint(
            errors.New("tenant name is required"),
            "Tenant Nameは必須です。",
        )
    }

    if utf8.RuneCountInString(string(n)) > 30 {
        return errors.WithHint(
            errors.New("Please enter a tenantname within 30 characters"),
            "テナント名は30文字以内で入力してください。",
        )
    }
    return nil
}

func NewTenantName(value string) (TenantName, error) {
    n := TenantName(value)
    if err := n.Validate(); err != nil {
        return "", err
    }
    return n, nil
}
```

**リポジトリインターフェース（ドメイン層に定義）**:
```go
// internal/domain/repository/tenant.go:32-38
type TenantRepository interface {
    CreateTenant(ctx context.Context, arg CreateTenantArg) error
    GetAllTenants(ctx context.Context) ([]model.Tenant, error)
    GetTenantsByUserID(ctx context.Context, userID model.UserID) ([]TenantWithMemberCount, error)
    GetTenantByID(ctx context.Context, id model.TenantID) (TenantWithJoinCode, error)
    UpdateTenant(ctx context.Context, arg UpdateTenantArg) error
}
```

**こだわりポイント**:
- **ドメインプリミティブ**: `TenantID`、`TenantName`などプリミティブ型をラップした型で表現力と型安全性を向上
- **不変条件の保証**: 各値オブジェクトは`Validate()`を持ち、ファクトリ関数で必ずバリデーションを実行
- **リポジトリパターン**: ドメイン層にインターフェースを定義し、インフラ層で実装（依存性逆転の原則）

---

## 7. クリーンアーキテクチャの実践

**適用箇所**:
```
backend/internal/
├── domain/           # ドメイン層（ビジネスルール）
│   ├── model/        # エンティティ・値オブジェクト
│   ├── repository/   # リポジトリインターフェース
│   ├── errors/       # ドメインエラー
│   └── logger/       # ロギング
├── usecase/          # ユースケース層（アプリケーションロジック）
│   ├── app/          # appユースケース
│   └── console/      # consoleユースケース
├── interface/        # インターフェース層（入出力）
│   ├── app/          # app向けハンドラー
│   └── console/      # console向けハンドラー
└── infrastructure/   # インフラ層（外部システム連携）
    ├── sqlc/         # DB実装
    ├── jwt/          # JWT処理
    └── auth/         # 認証サービス
```

**詳細**:

**依存方向**:
- interface → usecase → domain ← infrastructure
- ドメイン層は外部に依存しない（リポジトリはインターフェースのみ定義）

**こだわりポイント**:
- **層の分離**: 各層が明確な責務を持ち、依存関係が一方向
- **テスト容易性**: リポジトリインターフェースによりモックが容易
- **技術選択の柔軟性**: インフラ層を差し替えてもドメイン/ユースケースは影響なし

---

## 8. テスト&mock

**適用箇所**:
- [internal/usecase/console/tenant_test.go](backend/internal/usecase/console/tenant_test.go)
- [internal/domain/repository/mock/](backend/internal/domain/repository/mock/)
- [Taskfile.yaml](Taskfile.yaml#L130-L176) - mockタスク

**詳細**:

**テーブル駆動テスト**:
```go
// internal/usecase/console/tenant_test.go:19-171
func TestUseCase_CreateTenant(t *testing.T) {
    type fields struct {
        setupMock func(*mock.MockRepository)
    }
    type args struct {
        ctx   context.Context
        input dto.CreateTenantInput
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    string
        wantErr bool
        errType error
    }{
        {
            name: "正常系: テナント作成成功",
            fields: fields{
                setupMock: func(m *mock.MockRepository) {
                    m.EXPECT().
                        WithTransaction(gomock.Any(), gomock.Any()).
                        DoAndReturn(func(ctx context.Context, fn func(context.Context, repository.Transaction) error) error {
                            mockTx := mock.NewMockTransaction(gomock.NewController(t))
                            mockTx.EXPECT().CreateTenant(gomock.Any(), gomock.Any()).Return(nil)
                            mockTx.EXPECT().CreateTenantJoinCode(gomock.Any(), gomock.Any()).Return(nil)
                            return fn(ctx, mockTx)
                        })
                },
            },
            // ...
        },
        {
            name: "異常系: 名前が空",
            // ...
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockRepo := mock.NewMockRepository(ctrl)
            tt.fields.setupMock(mockRepo)

            u := &UseCase{repo: mockRepo, config: config.Config{}}

            got, err := u.CreateTenant(tt.args.ctx, tt.args.input)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.True(t, errors.Is(err, tt.errType))
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**mockgenによる自動生成**:
```go
// internal/domain/repository/repository.go:3
//go:generate go run go.uber.org/mock/mockgen@latest -source=$GOFILE -destination=mock/mock_repository.go -package=mock
```

**こだわりポイント**:
- **テーブル駆動テスト**: 正常系・異常系を網羅的にテスト、新規ケース追加が容易
- **uber/mockを使用したモック生成**: インターフェースから自動生成、型安全なモック
- **トランザクション内の動作もテスト**: `DoAndReturn`を使用してトランザクション内のモック動作を定義

---

## 9. SQL(RowLevelSecurity, sqlc.embed)

**適用箇所**:
- [db/migrations/20251007060706_add_rls_functions.sql](backend/db/migrations/20251007060706_add_rls_functions.sql)
- [db/migrations/20251007060707_add_tenants_table.sql](backend/db/migrations/20251007060707_add_tenants_table.sql)
- [db/sqlc/queries/tenant.sql](backend/db/sqlc/queries/tenant.sql)

**詳細**:

**RLS関数の定義**:
```sql
-- db/migrations/20251007060706_add_rls_functions.sql:7-15
CREATE OR REPLACE FUNCTION current_organization_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('keyhub.organization_id', true), '')::uuid
$$;
```

**RLSポリシーの適用**:
```sql
-- db/migrations/20251007060707_add_tenants_table.sql:20-30
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenants FORCE ROW LEVEL SECURITY;

CREATE POLICY tenants_org_isolation ON tenants
    FOR ALL
    TO keyhub
    USING (
        current_organization_id() IS NULL
        OR organization_id = current_organization_id()
    );
```

**sqlc.embedによるJOIN結果の構造化**:
```sql
-- db/sqlc/queries/tenant.sql:17-24
-- name: GetTenantById :one
SELECT
    sqlc.embed(t),
    sqlc.embed(jc)
FROM tenants t
INNER JOIN tenant_join_codes jc
    ON jc.tenant_id = t.id
WHERE t.id = $1;
```

**こだわりポイント**:
- **PostgreSQL RLSによるマルチテナント分離**: アプリケーションコードに依存せずDB層でデータ分離を保証
- **セッション変数によるRLS制御**: Contextから取得したorganization_idをDB接続時に自動設定
- **sqlc.embed**: JOIN結果を個別の構造体として取得、型安全なマッピング

---

# フロントエンド

## 1. TanStack Queryの活用

**適用箇所**:
- [frontend/app/src/lib/query.ts](frontend/app/src/lib/query.ts)
- [frontend/console/src/libs/query.ts](frontend/console/src/libs/query.ts)

**詳細**:

```typescript
// frontend/app/src/lib/query.ts:1-54
import { Code, ConnectError } from '@connectrpc/connect';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';

const retry = (failureCount: number, err: unknown) => {
  if (err instanceof ConnectError) {
    if (err.code === Code.PermissionDenied || err.code === Code.Unauthenticated) {
      return false; // 認証エラーはリトライしない
    }
  }
  return failureCount < 3;
};

export const queryClient = new QueryClient({
  queryCache: new QueryCache({ onError }),
  mutationCache: new MutationCache({ onError }),
  defaultOptions: {
    queries: {
      retry,
      staleTime: 60 * 1000, // 1分間はstaleとみなさない
    },
    mutations: {
      retry: false,
    },
  },
});

// connect-queryとの統合
export const useQueryGetTenantByJoinCode = (joinCode: string) => {
  return useQuery(getTenantByJoinCode, { joinCode }, { enabled: !!joinCode });
};
```

**こだわりポイント**:
- **@connectrpc/connect-queryとの統合**: Protocol Buffersで定義したAPIをTanStack Queryのhooksとして使用
- **スマートなリトライ戦略**: 認証エラーは即座に失敗、その他は3回まで再試行
- **適切なキャッシュ設定**: `staleTime`で不要なリフェッチを防止

---

## 2. Sentryの活用

**適用箇所**:
- [frontend/app/src/lib/sentry.ts](frontend/app/src/lib/sentry.ts)
- [frontend/console/src/libs/sentry.ts](frontend/console/src/libs/sentry.ts)

**詳細**:

```typescript
// frontend/app/src/lib/sentry.ts:90-124
const SKIP_RPC_CODES = [
  Code.Unauthenticated,  // 認証エラー（ログインページへ）
  Code.PermissionDenied, // 権限なし（正常なビジネスロジック）
  Code.InvalidArgument,  // バリデーションエラー
  Code.NotFound,         // リソースなし
  Code.AlreadyExists,    // 重複
];

function shouldCaptureError(error: unknown): boolean {
  if (!isSentryEnabled) return false;

  // Connect RPCエラーの場合、ビジネスエラーはスキップ
  if (error instanceof ConnectError) {
    if (SKIP_RPC_CODES.includes(error.code)) {
      return false;
    }
  }
  return true;
}

// エラーログ関数
export function logError(error: unknown, context?: { [key: string]: unknown }) {
  console.error('Error:', error, context);

  if (!shouldCaptureError(error)) return;

  if (error instanceof ConnectError) {
    Sentry.withScope((scope) => {
      scope.setFingerprint(['{{ default }}', endpoint, String(error.code), method || 'unknown']);
      scope.setContext('connectRPC', {
        code: error.code,
        message: error.message,
        endpoint,
        method,
        metadata: error.metadata,
      });
      scope.setTag('error_type', 'connect_rpc');
      Sentry.captureException(error);
    });
  }
}
```

**こだわりポイント**:
- **ノイズの削減**: ビジネスエラー（認証、バリデーション等）はSentryに送信せずログのみ
- **Connect RPCエラーの詳細なコンテキスト**: エンドポイント、メソッド、エラーコードをSentryのタグ・コンテキストに設定
- **適切なフィンガープリンティング**: 同じエンドポイント・エラーコードでグルーピング

---

## 3. Zodの活用(バリデーション)

**適用箇所**:
- [frontend/console/src/libs/utils/schema.ts](frontend/console/src/libs/utils/schema.ts)
- [frontend/console/src/hooks/useForm.ts](frontend/console/src/hooks/useForm.ts)

**詳細**:

**スキーマ定義**:
```typescript
// frontend/console/src/libs/utils/schema.ts:6-61
export const tenantnameValidation = z.preprocess(
  (val) => (typeof val === 'string' ? val.trim() : val),
  z
    .string({ message: 'テナント名を文字列で入力してください' })
    .nonempty({ message: 'テナント名を1文字以上入力してください' })
    .max(15, { message: 'テナント名は15文字以内で入力してください' })
    .refine((value: string) => !isBlankOrInvisible(value), {
      message: 'テナント名を1文字以上で入力してください',
    }),
);

export const joinCodeValidation = z.preprocess(
  (val) => (typeof val === 'string' ? val.trim() : val),
  z
    .string({ message: '参加コードを文字列で入力してください' })
    .min(6, { message: '参加コードは6文字以上で入力してください' })
    .max(20, { message: '参加コードは20文字以内で入力してください' })
    .regex(/^[a-zA-Z0-9]+$/, { message: '参加コードは英数字のみで入力してください' }),
);

export const tenantSchema = z.object({
  name: tenantnameValidation,
  description: descriptionValidation.optional(),
  tenantType: tenanttypeValidation,
  joinCode: joinCodeValidation,
  joinCodeExpiry: joinCodeExpiryValidation,
  joinCodeMaxUse: joinCodeMaxUseValidation,
});

export type TenantFormData = z.infer<typeof tenantSchema>;
```

**カスタムフォームフック**:
```typescript
// frontend/console/src/hooks/useForm.ts:13-96
export const useForm = <T extends z.ZodObject<z.ZodRawShape>>(
  schema: T,
  options: { revalidate?: boolean; initialValues?: Partial<z.infer<T>> } = {},
): UseFormReturn<T> => {
  // ...
  const validate = () => {
    const result = schema.safeParse(state);
    if (!result.success) {
      const newErrors = {} as FormErrorsType;
      result.error.issues.forEach((issue) => {
        const field = issue.path[0] as keyof FormType;
        if (!newErrors[field]) {
          newErrors[field] = [];
        }
        newErrors[field].push(issue.message);
      });
      setErrors(newErrors);
    }
    return result;
  };

  return { state, errors, updateField, validateField, validate, setState };
};
```

**こだわりポイント**:
- **日本語エラーメッセージ**: ユーザーに分かりやすいメッセージを定義
- **preprocess**: 入力値のトリミングをスキーマレベルで統一
- **型推論**: `z.infer<typeof schema>`で型を自動導出、フォーム状態とスキーマの一貫性を保証
- **リアルタイムバリデーション**: `revalidate`オプションで入力中の即時フィードバック

---

## 4. 適切な責任分離(Presentation Layerパターン)

**適用箇所**:
- [frontend/console/src/pages/CreateTenantPage.tsx](frontend/console/src/pages/CreateTenantPage.tsx) - Pageコンポーネント
- [frontend/console/src/components/CreateTenantForm.tsx](frontend/console/src/components/CreateTenantForm.tsx) - Formコンポーネント

**詳細**:

**Page（Container）コンポーネント**:
```typescript
// frontend/console/src/pages/CreateTenantPage.tsx:10-44
export const CreateTenantPage = () => {
  const navigate = useNavigate();
  const { mutateAsync: createTenant, isPending } = useMutationCreateTenant();

  const handleSubmit = async (data: { /* ... */ }) => {
    try {
      await createTenant({
        name: data.name,
        description: data.description || '',
        // ...
      });

      await queryClient.invalidateQueries();
      toast.success('テナントを作成しました');
      navigate('/tenants', { replace: true });
    } catch (error) {
      Sentry.captureException(error);
      toast.error('テナントの作成に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8">
        <h1>新しいテナントを作成</h1>
        <CreateTenantForm onSubmit={handleSubmit} isSubmitting={isPending} />
      </div>
    </div>
  );
};
```

**Form（Presentational）コンポーネント**:
```typescript
// frontend/console/src/components/CreateTenantForm.tsx:8-52
type CreateTenantFormProps = {
  onSubmit: (data: TenantFormData) => void;
  isSubmitting?: boolean;
};

export const CreateTenantForm = ({ onSubmit, isSubmitting = false }: CreateTenantFormProps) => {
  const form = useForm(tenantSchema, { revalidate: true, initialValues });

  const nameField = useFormField(form, 'name');
  const descriptionField = useFormField(form, 'description');
  // ...

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const result = form.validate();
    if (result.success) {
      onSubmit(result.data);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input value={nameField.value} onChange={nameField.onChange} onBlur={nameField.onBlur} />
      {nameField.error.length > 0 && <p>{nameField.error[0]}</p>}
      {/* ... */}
    </form>
  );
};
```

**こだわりポイント**:
- **Pageコンポーネント**: データ取得、mutation、ナビゲーション、エラーハンドリングを担当
- **Formコンポーネント**: 表示とフォームロジック（バリデーション含む）に集中、propsで受け取るシンプルなインターフェース
- **再利用性**: Formコンポーネントは編集画面でも再利用可能な設計

---

# その他

## 1. CI

**適用箇所**:
- [.github/workflows/](https://github.com/shibayama-club/keyhub/tree/main/.github/workflows)

**詳細**:

| ワークフロー | 目的 |
|------------|------|
| **go-lint.yaml** | golangci-lintによるGoコードの静的解析 |
| **go-test.yaml** | Goのユニットテスト実行 |
| **ts-lint.yaml** | ESLint + Prettierによるフロントエンドのlint/format |
| **gen-lint.yaml** | 自動生成コードの差分チェック（sqlc, protobuf, tbls） |

**gen-lint.yamlの特徴**:
```yaml
# .github/workflows/gen-lint.yaml
# PostgreSQLを起動し、マイグレーションを実行
- run: docker compose up -d postgres
- run: goose -dir ./backend/db/migrations postgres '...' up

# 各種コード生成
- run: cd backend && tbls doc --force -c tbls.yaml
- run: cd backend && sqlc generate
- run: buf generate
  working-directory: ./proto

# 生成コードに差分がないことを確認
- run: git diff --exit-code
```

**こだわりポイント**:
- **自動生成コードの整合性保証**: CIで差分を検出することで、生成忘れを防止
- **実際のDBを使用したスキーマ生成**: マイグレーションを適用したPostgreSQLからER図やスキーマを生成

---

## 2. ER図の自動生成

**適用箇所**:
- [backend/tbls.yaml](backend/tbls.yaml)
- [backend/docs/schema/](backend/docs/schema/)

**詳細**:

[tbls](https://github.com/k1LoW/tbls)を使用し、PostgreSQLスキーマからER図とテーブルドキュメントを自動生成。

```yaml
# backend/tbls.yaml
dsn: postgresql://postgres:keyhub@localhost:5432/keyhub?sslmode=disable
docPath: docs/schema

exclude:
  - goose_db_version  # マイグレーション管理テーブルを除外

er:
  format: svg
  comment: true
  distance: 2
```

**生成成果物**:
- `schema.svg` - 全体ER図
- `public.tenants.svg` - テーブルごとのER図
- `public.tenants.md` - カラム定義、インデックス、外部キー情報

---

## 3. protobuf, connectRPC

**適用箇所**:
- [proto/](https://github.com/shibayama-club/keyhub/tree/main/proto)
- [proto/buf.gen.yaml](proto/buf.gen.yaml)

**詳細**:

[Connect](https://connectrpc.com/)を使用し、Protocol BuffersベースのRPCを実現。

```yaml
# proto/buf.gen.yaml
plugins:
  # Go向け
  - remote: buf.build/protocolbuffers/go
    out: ../backend/internal/interface/gen
  - remote: buf.build/connectrpc/go
    out: ../backend/internal/interface/gen

  # TypeScript向け
  - remote: buf.build/bufbuild/es
    out: ../frontend/gen/src
    opt: target=ts
  - remote: buf.build/connectrpc/query-es  # TanStack Query統合
    out: ../frontend/gen/src
```

**こだわりポイント**:
- **connect-query-es**: Protocol Buffersの定義からTanStack Queryのhooksを自動生成
- **型安全なAPI通信**: フロントエンド・バックエンドで同じスキーマを共有

---

## 4. コードの自動生成(開発体験の向上)

**適用箇所**:
- [Taskfile.yaml](Taskfile.yaml)

**詳細**:

```yaml
# Taskfile.yaml - 主要な生成タスク

# Protocol Buffers → Go/TypeScript
proto:
  dir: ./proto
  cmds:
    - buf dep update
    - buf generate
    - buf lint
    - buf format -w

# SQLクエリ → Go (型安全なDB操作)
gen:sqlc:
  dir: ./backend
  cmds:
    - sqlc generate -f sqlc.yaml

# DBスキーマ → ER図/ドキュメント
gen:docs:schema:
  dir: ./backend
  cmds:
    - tbls doc --rm-dist -c tbls.yaml

# モック生成
mock:
  deps:
    - mock:install
    - mock:console
    - mock:repository
```

**こだわりポイント**:
- **task一発で全生成**: `task gen`で全ての自動生成を実行
- **増分生成**: `sources`/`generates`でファイル変更時のみ再生成
- **開発ワークフローの統一**: 新規参加者も`task init && task gen`で環境構築完了
