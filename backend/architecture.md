# KeyHub Backend Architecture

このドキュメントでは、KeyHub Backendプロジェクトのアーキテクチャ設計について詳しく説明します。

---

## 目次

1. [アーキテクチャ概要](#アーキテクチャ概要)
2. [ディレクトリ構造](#ディレクトリ構造)
3. [レイヤードアーキテクチャ](#レイヤードアーキテクチャ)
4. [依存関係の方向](#依存関係の方向)
5. [各レイヤーの詳細](#各レイヤーの詳細)
6. [ドメイン駆動設計 (DDD)](#ドメイン駆動設計-ddd)
7. [データフロー](#データフロー)
8. [設計原則](#設計原則)
9. [技術スタック](#技術スタック)

---

## アーキテクチャ概要

KeyHub Backendは**クリーンアーキテクチャ**と**Composition Rootパターン**を組み合わせたレイヤードアーキテクチャを採用しています。

### 主要な特徴

- **レイヤー分離**: Domain、UseCase、Interface、Infrastructureの4層構造
- **依存性逆転の原則 (DIP)**: 内側のレイヤーは外側のレイヤーに依存しない
- **Composition Root**: `cmd/`層で具象実装を組み立て、依存性を注入
- **Connect RPC**: gRPC互換のHTTP/2ベースのRPCフレームワーク

---

## ディレクトリ構造

```
backend/
├── cmd/                        # エントリポイント・Composition Root
│   ├── cmd.go                 # CLIルート
│   ├── config/                # 設定管理
│   │   └── config.go          # 設定の読み込み・バリデーション
│   └── serve/                 # サーバー起動
│       ├── app.go             # App APIサーバー
│       └── console.go         # Console APIサーバー
│
├── internal/                   # プライベートパッケージ
│   ├── domain/                # ドメイン層（最内層）
│   │   ├── authenticator/     # 認証器インターフェース
│   │   ├── errors/            # ドメインエラー定義
│   │   ├── healthcheck/       # ヘルスチェック抽象化
│   │   ├── logger/            # ロガーインターフェース
│   │   ├── model/             # ドメインモデル
│   │   └── repository/        # リポジトリインターフェース
│   │
│   ├── usecase/               # ユースケース層（アプリケーションロジック）
│   │   ├── app/               # App機能のユースケース
│   │   └── console/           # Console機能のユースケース
│   │
│   ├── interface/             # インターフェース層（外部との境界）
│   │   ├── app/v1/            # App API v1ハンドラー
│   │   ├── console/v1/        # Console API v1ハンドラー
│   │   │   └── interceptor/   # Connect RPCインターセプター
│   │   │       ├── auth.go    # 認証インターセプター
│   │   │       └── sentry.go  # Sentryエラー送信インターセプター
│   │   ├── gen/               # 自動生成コード（Protobuf）
│   │   └── health/            # ヘルスチェックハンドラー
│   │
│   └── infrastructure/        # インフラストラクチャ層（外部システム統合）
│       ├── auth/              # 認証実装
│       │   ├── claim/         # クレームベース認証
│       │   └── console/       # コンソール認証
│       ├── jwt/               # JWT実装
│       └── sqlc/              # データベースアクセス
│           └── gen/           # SQLC自動生成コード
│
├── config.yaml                # 設定ファイル
├── go.mod                     # Go依存関係
└── main.go                    # アプリケーションエントリポイント
```

---

## レイヤードアーキテクチャ

KeyHub Backendは4つの主要レイヤーで構成されています：

```
┌─────────────────────────────────────────────────────────────┐
│                     cmd/ (Composition Root)                  │
│  - エントリポイント                                            │
│  - 依存性の組み立て（DI）                                      │
│  - サーバー設定・起動                                          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              interface/ (Adapters - 外向き)                  │
│  - HTTPハンドラー                                             │
│  - gRPCサービス実装                                           │
│  - リクエスト・レスポンスの変換                                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  usecase/ (Application Logic)                │
│  - ビジネスロジックのオーケストレーション                       │
│  - トランザクション管理                                        │
│  - 複数ドメインの調整                                          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                domain/ (Business Rules - 最内層)              │
│  - ドメインモデル                                              │
│  - ビジネスルール                                              │
│  - インターフェース定義（Repository, Authenticatorなど）        │
└─────────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────────┐
│            infrastructure/ (Adapters - 内向き)               │
│  - データベースアクセス                                        │
│  - 外部API連携                                                │
│  - 認証実装                                                   │
│  - ミドルウェア                                                │
└─────────────────────────────────────────────────────────────┘
```

---

## 依存関係の方向

### 原則: **内側のレイヤーは外側のレイヤーに依存しない**

```
                         依存の方向
                              ↓

┌──────────────────────────────────────────────────────────┐
│  cmd/                                                     │
│  ├─→ interface/        (ハンドラーを利用)                 │
│  ├─→ usecase/          (ユースケースを利用)               │
│  ├─→ infrastructure/   (具象実装を注入)                   │
│  └─→ domain/           (型定義を利用)                     │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  interface/                                               │
│  ├─→ usecase/          (ユースケースに委譲)               │
│  └─→ domain/           (ドメインモデルを利用)             │
│  ✗ infrastructure/     (依存しない)                       │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  usecase/                                                 │
│  ├─→ domain/           (ドメインロジックを利用)           │
│  ✗ interface/          (依存しない)                       │
│  ✗ infrastructure/     (具象実装に依存しない)             │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  domain/                                                  │
│  ✗ どのレイヤーにも依存しない（最内層）                    │
│  - 標準ライブラリのみ使用                                  │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  infrastructure/                                          │
│  ├─→ domain/           (インターフェースを実装)           │
│  ✗ usecase/            (依存しない)                       │
│  ✗ interface/          (依存しない)                       │
└──────────────────────────────────────────────────────────┘
```

### 依存性逆転の原則 (DIP)

```
┌─────────────────┐          ┌─────────────────┐
│    usecase/     │          │ infrastructure/ │
│                 │          │                 │
│  LoginUseCase   │          │  PostgresRepo   │
│                 │          │                 │
└────────┬────────┘          └────────┬────────┘
         │                            │
         │ 依存                        │ 実装
         ↓                            ↓
    ┌─────────────────────────────────┐
    │         domain/                 │
    │                                 │
    │  interface Repository {         │
    │      GetUser(id) User           │
    │  }                              │
    └─────────────────────────────────┘

✅ usecaseは具象実装に依存せず、interfaceに依存
✅ infrastructureがinterfaceを実装
✅ cmd/でPostgresRepoをusecaseに注入
```

---

## 各レイヤーの詳細

### 1. cmd/ - Composition Root

**責務**: アプリケーション全体の組み立て

```go
// cmd/serve/console.go
func SetupConsole(ctx context.Context, config config.Config) (*echo.Echo, error) {
    // 1. 外部サービスの初期化
    sentry.Init(config.Sentry.DSN, ...)

    // 2. Infrastructure層の具象実装を作成
    pool := sqlc.NewPool(ctx, config.Postgres)
    repo := sqlc.NewRepository(pool)
    consoleAuth := consoleauth.NewAuthService(jwtSecret)

    // 3. UseCase層に依存性を注入
    consoleUseCase := console.NewUseCase(ctx, repo, config, consoleAuth)

    // 4. Interface層にインターセプターとハンドラーを作成
    sentryInterceptor := interceptor.NewSentryInterceptor()
    authInterceptor := interceptor.NewAuthInterceptor(consoleUseCase)
    consoleHandler := consolev1.NewHandler(consoleUseCase, jwtSecret)

    // 5. Connect RPCサービスにインターセプターを登録
    e := echo.New()
    authPath, authHandler := consolev1connect.NewConsoleAuthServiceHandler(
        consoleHandler,
        connect.WithInterceptors(sentryInterceptor, authInterceptor),
    )
    servicePath, serviceHandler := consolev1connect.NewConsoleServiceHandler(
        consoleHandler,
        connect.WithInterceptors(sentryInterceptor, authInterceptor),
    )
    e.Any(authPath+"*", echo.WrapHandler(authHandler))
    e.Any(servicePath+"*", echo.WrapHandler(serviceHandler))

    return e, nil
}
```

**特徴**:
- ✅ すべての具象実装を直接import
- ✅ 依存関係を明示的に構築
- ✅ ビジネスロジックを含まない
- ✅ Go言語の標準的なComposition Rootパターン

---

### 2. domain/ - ドメイン層（最内層）

**責務**: ビジネスロジックとルール

```go
// domain/model/user.go
type User struct {
    ID        string
    Email     string
    CreatedAt time.Time
}

// domain/repository/user.go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
}

// domain/authenticator/authenticator.go
type Authenticator interface {
    GenerateToken(userID string) (string, error)
    ValidateToken(token string) (*Claims, error)
}
```

**特徴**:
- ✅ 外部への依存なし（標準ライブラリのみ）
- ✅ インターフェース定義（具象実装は含まない）
- ✅ ドメインモデル・ビジネスルール
- ✅ 変更に強い（他のレイヤーの影響を受けにくい）

---

### 3. usecase/ - ユースケース層

**責務**: アプリケーションロジックのオーケストレーション

```go
// usecase/console/login.go
type UseCase struct {
    repo domain.Repository
    auth domain.Authenticator
}

func (u *UseCase) Login(ctx context.Context, orgID, orgKey string) (string, error) {
    // 1. ドメインロジック: 組織の検証
    org, err := u.repo.GetOrganization(ctx, orgID)
    if err != nil {
        return "", err
    }

    // 2. ドメインロジック: 認証キーの検証
    if !org.ValidateKey(orgKey) {
        return "", ErrInvalidCredentials
    }

    // 3. Infrastructure: トークン生成
    token, err := u.auth.GenerateToken(orgID)
    if err != nil {
        return "", err
    }

    return token, nil
}
```

**特徴**:
- ✅ ドメインインターフェースに依存
- ✅ 具象実装には依存しない
- ✅ トランザクション管理
- ✅ 複数リポジトリの調整

---

### 4. interface/ - インターフェース層

**責務**: 外部との境界・プロトコル変換

```go
// interface/console/v1/handler.go
type Handler struct {
    usecase *console.UseCase
}

func (h *Handler) LoginWithOrgId(
    ctx context.Context,
    req *connect.Request[consolev1.LoginWithOrgIdRequest],
) (*connect.Response[consolev1.LoginWithOrgIdResponse], error) {
    // 1. リクエストのバリデーション
    if req.Msg.OrganizationId == "" {
        return nil, connect.NewError(connect.CodeInvalidArgument, ...)
    }

    // 2. ユースケースに委譲
    token, err := h.usecase.Login(ctx,
        req.Msg.OrganizationId,
        req.Msg.OrganizationKey,
    )
    if err != nil {
        return nil, toConnectError(err)
    }

    // 3. レスポンスの構築
    return connect.NewResponse(&consolev1.LoginWithOrgIdResponse{
        SessionToken: token,
        ExpiresIn:    3600,
    }), nil
}
```

**特徴**:
- ✅ HTTPリクエスト・レスポンスの処理
- ✅ プロトコル固有のエラーハンドリング
- ✅ ユースケースへの委譲
- ✅ データ変換（DTO ↔ ドメインモデル）

---

### 5. infrastructure/ - インフラストラクチャ層

**責務**: 外部システムとの統合

#### 5-1. infrastructure/sqlc/ - データベースアクセス

```go
// infrastructure/sqlc/repository.go
type Repository struct {
    db *pgxpool.Pool
}

// ドメインのRepositoryインターフェースを実装
func (r *Repository) GetOrganization(ctx context.Context, id string) (*model.Organization, error) {
    row, err := r.queries.GetOrganization(ctx, id)
    if err != nil {
        return nil, err
    }

    // データベースモデル → ドメインモデル
    return &model.Organization{
        ID:  row.ID,
        Key: row.Key,
    }, nil
}
```

#### 5-2. infrastructure/auth/ - 認証実装

```go
// infrastructure/auth/console/auth.go
type AuthService struct {
    jwtSecret []byte
}

// ドメインのAuthenticatorインターフェースを実装
func (a *AuthService) GenerateToken(orgID string) (string, error) {
    claims := jwt.MapClaims{
        "organization_id": orgID,
        "exp": time.Now().Add(time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.jwtSecret)
}
```

#### 5-3. interface/console/v1/interceptor/ - Connect RPCインターセプター

```go
// interface/console/v1/interceptor/sentry.go
type sentryInterceptor struct{}

func NewSentryInterceptor() connect.Interceptor {
    return &sentryInterceptor{}
}

func (i *sentryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
    return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
        resp, err := next(ctx, req)
        if err != nil {
            // 5xxエラー（サーバー内部エラー）のみSentryに送信
            if connectErr, ok := err.(*connect.Error); ok {
                if connectErr.Code() == connect.CodeInternal ||
                    connectErr.Code() == connect.CodeUnknown ||
                    connectErr.Code() == connect.CodeDataLoss {

                    hub := sentry.CurrentHub().Clone()
                    hub.WithScope(func(scope *sentry.Scope) {
                        scope.SetTag("rpc_method", req.Spec().Procedure)
                        scope.SetContext("request", map[string]interface{}{
                            "procedure": req.Spec().Procedure,
                            "peer":      req.Peer(),
                        })
                        hub.CaptureException(err)
                    })
                }
            }
        }
        return resp, err
    }
}
```

**特徴**:
- ✅ Connect RPCのインターセプターインターフェースを実装
- ✅ 5xxエラー（Internal、Unknown、DataLoss）のみをSentryに送信
- ✅ RPCメソッド名やピア情報などのコンテキストを付与
- ✅ Unary（単一リクエスト）とStreaming（ストリーミング）の両方に対応

---

## ドメイン駆動設計 (DDD)

KeyHub Backendでは、**ドメイン駆動設計 (Domain-Driven Design: DDD)** の原則とパターンを活用しています。

### DDD概要

DDDは、複雑なビジネスロジックをソフトウェアで表現するための設計手法です。以下の2つの側面があります:

1. **戦略的設計 (Strategic Design)**: システム全体の構造・境界の定義
2. **戦術的設計 (Tactical Design)**: ドメインモデルの実装パターン

---

### 戦術的パターン (Tactical Patterns)

KeyHubで活用しているDDDの戦術的パターンを紹介します。

#### 1. Value Object（値オブジェクト）

**特徴**: 不変で、同一性を持たないオブジェクト

```go
// domain/model/organization.go

// OrganizationID - 値オブジェクト
type OrganizationID uuid.UUID

func (id OrganizationID) String() string {
    return uuid.UUID(id).String()
}

// バリデーションロジックを内包
func (id OrganizationID) Validate() error {
    if id == OrganizationID(uuid.Nil) {
        return errors.WithHint(
            errors.New("organization ID is required"),
            "組織IDは必須です。",
        )
    }
    return nil
}

// ファクトリメソッド - 不正な値の生成を防ぐ
func NewOrganizationID(id uuid.UUID) (OrganizationID, error) {
    orgID := OrganizationID(id)
    if err := orgID.Validate(); err != nil {
        return OrganizationID(uuid.Nil), err
    }
    return orgID, nil
}

// OrganizationKey - 値オブジェクト
type OrganizationKey string

func (k OrganizationKey) Validate() error {
    if k == "" {
        return errors.New("organization key is required")
    }

    length := utf8.RuneCountInString(string(k))
    if length < 1 || length > 20 {
        return errors.New("key must be between 1 and 20 characters")
    }

    return nil
}

func NewOrganizationKey(value string) (OrganizationKey, error) {
    k := OrganizationKey(value)
    if err := k.Validate(); err != nil {
        return "", err
    }
    return k, nil
}
```

**メリット**:
- ✅ ビジネスルールを型で表現
- ✅ 不正な値の生成を防止
- ✅ ドメイン知識のカプセル化

---

#### 2. Entity（エンティティ）

**特徴**: 同一性（ID）を持つオブジェクト

```go
// domain/model/organization.go
type Organization struct {
    ID  OrganizationID   // ← 同一性
    Key OrganizationKey
}

// domain/model/user.go
type User struct {
    ID        UserID       // ← 同一性
    Email     Email
    TenantID  TenantID
    CreatedAt time.Time
}
```

**エンティティの判定基準**:
- 同じ属性でも別のインスタンス（IDが異なる）なら異なるオブジェクト
- ライフサイクルを持つ（作成・更新・削除）

---

#### 3. Repository（リポジトリ）

**特徴**: ドメインオブジェクトの永続化を抽象化

```go
// domain/repository/user.go
type UserRepository interface {
    GetByID(ctx context.Context, id model.UserID) (*model.User, error)
    GetByEmail(ctx context.Context, email model.Email) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id model.UserID) error
}

// domain/repository/repository.go
type Repository interface {
    Transaction  // トランザクション境界

    WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type Transaction interface {
    UserRepository
    TenantRepository
    ConsoleSessionRepository
}
```

**DDDのリポジトリパターンの特徴**:
- ✅ ドメイン層にインターフェース定義（DIP: 依存性逆転）
- ✅ コレクション風のAPI（`GetByID`, `Create`, `Update`）
- ✅ 永続化の詳細を隠蔽（SQLやORMの詳細は隠す）

**実装** (`infrastructure/sqlc/`):
```go
// infrastructure/sqlc/repository.go
type Repository struct {
    db      *pgxpool.Pool
    queries *gen.Queries
}

// domain.UserRepositoryインターフェースを実装
func (r *Repository) GetByID(ctx context.Context, id model.UserID) (*model.User, error) {
    row, err := r.queries.GetUserByID(ctx, id.UUID())
    if err != nil {
        return nil, err
    }

    // データベースモデル → ドメインモデル変換
    return &model.User{
        ID:        model.UserID(row.ID),
        Email:     model.Email(row.Email),
        TenantID:  model.TenantID(row.TenantID),
        CreatedAt: row.CreatedAt,
    }, nil
}
```

---

#### 4. Domain Service（ドメインサービス）

**特徴**: エンティティや値オブジェクトに属さないドメインロジック

```go
// domain/authenticator/auth_console.go
type ConsoleAuthenticator interface {
    GenerateToken(organizationID, sessionID string, expiresIn time.Duration) (string, error)
    ValidateToken(token string) (*claim.ConsoleClaims, error)
}

// domain/authenticator/auth_user.go
type UserAuthenticator interface {
    GenerateToken(userID, tenantID string, expiresIn time.Duration) (string, error)
    ValidateToken(token string) (*claim.UserClaims, error)
}
```

**ドメインサービスの判定基準**:
- エンティティのメソッドとして定義するのが不自然
- 複数のエンティティにまたがる処理
- ステートレスな操作

**実装** (`infrastructure/auth/`):
```go
// infrastructure/auth/console/auth.go
type AuthService struct {
    jwtSecret []byte
}

func (a *AuthService) GenerateToken(orgID, sessionID string, expiresIn time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "organization_id": orgID,
        "session_id":      sessionID,
        "exp":             time.Now().Add(expiresIn).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.jwtSecret)
}
```

---

#### 5. Application Service（アプリケーションサービス）

**特徴**: ユースケースの実装・ドメインオブジェクトのオーケストレーション

```go
// usecase/console/usecase.go
type UseCase struct {
    repo repository.Repository
    auth authenticator.ConsoleAuthenticator
}

func (u *UseCase) LoginWithOrgId(
    ctx context.Context,
    orgID model.OrganizationID,
    orgKey model.OrganizationKey,
) (string, error) {
    // 1. リポジトリでエンティティ取得
    session, err := u.repo.CreateConsoleSession(ctx, orgID, orgKey)
    if err != nil {
        return "", err
    }

    // 2. ドメインサービスでトークン生成
    token, err := u.auth.GenerateToken(
        orgID.String(),
        session.ID.String(),
        time.Hour,
    )
    if err != nil {
        return "", err
    }

    return token, nil
}
```

**アプリケーションサービス（ユースケース）の責務**:
- ✅ トランザクション境界の管理
- ✅ ドメインオブジェクトの調整
- ✅ ビジネスロジックの実行順序の決定
- ❌ ビジネスルール自体は含まない（ドメイン層に委譲）

---

### 戦略的パターン (Strategic Patterns)

#### 1. Bounded Context（境界付けられたコンテキスト）

KeyHubでは、以下の境界付けられたコンテキストを識別しています:

```
┌─────────────────────────────────────────────────────────┐
│                    KeyHub System                         │
│                                                           │
│  ┌──────────────────────┐  ┌─────────────────────────┐  │
│  │  Console Context     │  │   App Context           │  │
│  │  (管理者向け)         │  │   (エンドユーザー向け)   │  │
│  │                      │  │                         │  │
│  │  - Organization      │  │  - User                 │  │
│  │  - ConsoleSession    │  │  - Tenant               │  │
│  │  - ConsoleAuth       │  │  - UserAuth             │  │
│  └──────────────────────┘  └─────────────────────────┘  │
│           ↓                          ↓                   │
│  ┌───────────────────────────────────────────────────┐  │
│  │        Shared Kernel (共有カーネル)                │  │
│  │  - Repository interfaces                          │  │
│  │  - Common domain errors                           │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**境界の分離**:
- **Console Context**: 組織管理・コンソール認証
- **App Context**: ユーザー管理・アプリ認証
- **Shared Kernel**: 共通のインターフェース・エラー定義

**実装の分離**:
```
internal/
├── usecase/
│   ├── console/          # Console Context
│   └── app/              # App Context
├── interface/
│   ├── console/v1/       # Console Context API
│   └── app/v1/           # App Context API
└── domain/
    ├── model/            # 各コンテキストのモデル
    └── repository/       # 共有インターフェース
```

---

#### 2. Ubiquitous Language（ユビキタス言語）

ドメインの専門用語をコード全体で統一的に使用しています。

| ユビキタス言語 | コード上の表現 | 説明 |
|--------------|--------------|------|
| **Organization** | `model.Organization` | 組織エンティティ |
| **Organization ID** | `model.OrganizationID` | 組織の識別子（値オブジェクト） |
| **Organization Key** | `model.OrganizationKey` | 組織の認証キー（値オブジェクト） |
| **Console Session** | `model.ConsoleSession` | コンソールセッション |
| **Tenant** | `model.Tenant` | テナント（マルチテナント） |
| **User** | `model.User` | ユーザーエンティティ |
| **Authenticator** | `authenticator.ConsoleAuthenticator` | 認証サービス |

**ユビキタス言語の利点**:
- ✅ ビジネスサイドと開発サイドの共通理解
- ✅ コードが仕様書になる
- ✅ 変更時の影響範囲が明確

---

### DDDの適用レベル

KeyHubでは、以下のようにDDDパターンを適用しています:

| パターン | 適用度 | 実装場所 | 備考 |
|---------|-------|---------|------|
| **Value Object** | ⭐⭐⭐⭐⭐ | `domain/model/*` | ID型、Key型など |
| **Entity** | ⭐⭐⭐⭐⭐ | `domain/model/*` | User, Organization, Tenant |
| **Repository** | ⭐⭐⭐⭐⭐ | `domain/repository/*` | インターフェース定義 |
| **Domain Service** | ⭐⭐⭐⭐ | `domain/authenticator/*` | 認証サービス |
| **Application Service** | ⭐⭐⭐⭐⭐ | `usecase/*` | ユースケース実装 |
| **Aggregate** | ⭐⭐⭐ | 一部適用 | トランザクション境界 |
| **Domain Event** | ⭐ | 未適用 | 将来的に検討 |
| **Specification** | ⭐ | 未適用 | 将来的に検討 |

---

### DDD適用の利点

KeyHubでDDDを適用することで、以下の利点を享受しています:

1. **ビジネスロジックの明確化**
   - ドメインモデルがビジネスルールを表現
   - コードがドメイン知識の宝庫

2. **変更への柔軟性**
   - ビジネスルールの変更がドメイン層に局所化
   - 技術的詳細の変更がインフラ層に局所化

3. **テスタビリティ**
   - ドメインロジックが純粋関数的
   - リポジトリのモック化が容易

4. **チーム間のコミュニケーション**
   - ユビキタス言語による共通理解
   - コードがドキュメント

---

### DDDのベストプラクティス

KeyHubで実践しているDDDのベストプラクティス:

#### ✅ Do（推奨）

```go
// ✅ 値オブジェクトでバリデーション
func NewOrganizationKey(value string) (OrganizationKey, error) {
    k := OrganizationKey(value)
    if err := k.Validate(); err != nil {
        return "", err
    }
    return k, nil
}

// ✅ ドメインエラーを明確に
var (
    ErrOrganizationNotFound = errors.New("organization not found")
    ErrInvalidCredentials   = errors.New("invalid credentials")
)

// ✅ リポジトリインターフェースはドメイン層
type UserRepository interface {
    GetByID(ctx context.Context, id UserID) (*User, error)
}
```

#### ❌ Don't（非推奨）

```go
// ❌ プリミティブ型の多用（Primitive Obsession）
func Login(ctx context.Context, orgID string, orgKey string) error {
    // stringのままではビジネスルールが表現できない
}

// ❌ ドメインロジックをユースケース層に
func (u *UseCase) Login(ctx context.Context, orgID, orgKey string) error {
    // バリデーションをここに書くのは NG
    if len(orgKey) < 1 || len(orgKey) > 20 {
        return errors.New("invalid key length")
    }
}

// ❌ リポジトリの実装詳細を漏らす
type UserRepository interface {
    ExecuteSQL(query string) error  // NG: SQL詳細を露出
}
```

---

### まとめ: DDDの活用

KeyHub Backendでは、DDDの以下のパターンを積極的に活用しています:

- **Value Object**: 型安全性とビジネスルールのカプセル化
- **Entity**: ドメインの中心オブジェクト
- **Repository**: 永続化の抽象化
- **Domain Service**: ビジネスロジックの適切な配置
- **Application Service**: ユースケースのオーケストレーション
- **Bounded Context**: コンテキストの明確な分離
- **Ubiquitous Language**: 共通言語によるコミュニケーション

これにより、**保守性**、**拡張性**、**テスタビリティ**の高いシステムを実現しています。

---

## データフロー

### リクエスト処理フロー

```
1. HTTPリクエスト受信
        ↓
2. Echo Middleware処理
   - Recover (パニック復旧)
   - Logging (ログ出力)
   - CORS (クロスオリジン)
        ↓
3. Connect RPC Interceptor処理
   - SentryInterceptor (5xxエラー監視)
   - AuthInterceptor (認証・認可)
        ↓
4. Interface層 (Handler)
   - リクエストバリデーション
   - DTOからドメインモデルへ変換
        ↓
5. UseCase層
   - ビジネスロジック実行
   - Repositoryでデータ取得
   - ドメインルール適用
        ↓
6. Infrastructure層
   - データベースアクセス
   - 外部API呼び出し
        ↓
7. Domain層
   - ドメインモデルの操作
   - ビジネスルール検証
        ↓
8. レスポンス返却（逆順）
   - ドメインモデル → DTO
   - HTTPレスポンス構築
   - Interceptorでエラーをキャプチャ（5xxのみSentryへ）
```

### 具体例: ログインフロー

```
[Client] POST /keyhub.console.v1.ConsoleAuthService/LoginWithOrgId
    ↓
[Echo Middleware] Recover, Logging, CORS
    ↓
[Interceptor] SentryInterceptor, AuthInterceptor
    ↓
[Interface] consolev1.Handler.LoginWithOrgId()
    ├─ リクエストバリデーション
    ├─ req.Msg → (organizationId, organizationKey)
    └─ usecase.Login() を呼び出し
        ↓
[UseCase] console.UseCase.Login()
    ├─ repo.GetOrganization() でDB取得
    │   ↓
    │  [Infrastructure] sqlc.Repository.GetOrganization()
    │       ├─ SQL実行: SELECT * FROM organizations
    │       └─ DB Row → domain.Organization
    │
    ├─ org.ValidateKey(orgKey) でビジネスルール検証
    │   ↓
    │  [Domain] Organization.ValidateKey()
    │       └─ ハッシュ比較・検証ロジック
    │
    └─ auth.GenerateToken(orgID) でトークン生成
        ↓
       [Infrastructure] console.AuthService.GenerateToken()
            └─ JWT生成・署名
    ↓
[Interface] Handler
    ├─ token → LoginWithOrgIdResponse
    └─ connect.Response構築
    ↓
[Client] レスポンス受信
```

---

## 設計原則

### 1. SOLID原則

#### Single Responsibility Principle (単一責任の原則)
- 各レイヤーは明確な責務を持つ
- 例: `UserRepository`はユーザーのデータアクセスのみ担当

#### Open/Closed Principle (オープン・クローズドの原則)
- インターフェースで拡張可能
- 例: 新しい認証方式を`Authenticator`インターフェース実装で追加

#### Liskov Substitution Principle (リスコフの置換原則)
- インターフェース実装は置換可能
- 例: `PostgresRepository` ↔ `MockRepository`

#### Interface Segregation Principle (インターフェース分離の原則)
- 小さく焦点を絞ったインターフェース
- 例: `UserRepository`, `OrganizationRepository` を分離

#### Dependency Inversion Principle (依存性逆転の原則)
- 抽象に依存、具象に依存しない
- 例: UseCaseはRepositoryインターフェースに依存

---

### 2. Composition Root パターン

```go
// cmd/serve/console.go - 唯一の配線箇所
func SetupConsole(ctx context.Context, config config.Config) (*echo.Echo, error) {
    // すべての依存関係をここで解決
    db := sqlc.NewPool(ctx, config.Postgres)          // ← 具象
    repo := sqlc.NewRepository(db)                     // ← 具象
    auth := consoleauth.NewAuthService(jwtSecret)      // ← 具象

    // 抽象に注入
    uc := console.NewUseCase(ctx, repo, config, auth)  // ← インターフェース経由
    handler := consolev1.NewHandler(uc)                // ← インターフェース経由

    return e, nil
}
```

**メリット**:
- ✅ 依存関係が一箇所に集約
- ✅ テスト時にモック注入が容易
- ✅ 実装の切り替えが簡単

---

### 3. クリーンアーキテクチャ

```
┌─────────────────────────────────────────────┐
│            外側（詳細・技術）                 │
│  ┌───────────────────────────────────────┐  │
│  │         中間（ロジック）                 │  │
│  │  ┌─────────────────────────────────┐  │  │
│  │  │    内側（ビジネスルール）         │  │  │
│  │  │                                 │  │  │
│  │  │  Domain                         │  │  │
│  │  │  - モデル                        │  │  │
│  │  │  - インターフェース              │  │  │
│  │  └─────────────────────────────────┘  │  │
│  │                                       │  │
│  │  UseCase                              │  │
│  │  - アプリケーションロジック           │  │
│  └───────────────────────────────────────┘  │
│                                             │
│  Interface / Infrastructure                 │
│  - HTTP, DB, 外部サービス                   │
└─────────────────────────────────────────────┘

依存の方向: 外 → 内
```

---

## 技術スタック

### フレームワーク・ライブラリ

| 用途 | 技術 | 説明 |
|------|------|------|
| **Web Framework** | [Echo v4](https://echo.labstack.com/) | 高速・軽量なHTTPフレームワーク |
| **RPC** | [Connect RPC](https://connectrpc.com/) | gRPC互換のHTTP/2 RPC |
| **Database** | PostgreSQL + [pgx](https://github.com/jackc/pgx) | 高性能なPostgreSQLドライバ |
| **SQL Builder** | [SQLC](https://sqlc.dev/) | 型安全なSQLコード生成 |
| **JWT** | [golang-jwt](https://github.com/golang-jwt/jwt) | JWT認証実装 |
| **Validation** | [validator](https://github.com/go-playground/validator) | 構造体バリデーション |
| **Config** | [Viper](https://github.com/spf13/viper) + [Cobra](https://github.com/spf13/cobra) | 設定管理・CLI |
| **Logging** | `log/slog` + [slog-echo](https://github.com/samber/slog-echo) | 構造化ログ |
| **Error Tracking** | [Sentry](https://sentry.io/) | エラー監視・トラッキング |
| **Error Handling** | [cockroachdb/errors](https://github.com/cockroachdb/errors) | 拡張エラー処理 |

---

### プロトコル

- **HTTP/2** (h2c - TLS無しHTTP/2)
- **Connect RPC** (gRPC互換のJSON/Protobuf)

---

## まとめ

KeyHub Backendのアーキテクチャは以下の特徴を持ちます:

### ✅ 長所

1. **明確な責務分離**: 各レイヤーが独立した責務を持つ
2. **テスタビリティ**: インターフェースによりモック化が容易
3. **保守性**: ビジネスロジックが技術的詳細から分離
4. **拡張性**: 新機能追加時の影響範囲が限定的
5. **Go言語らしさ**: Composition Rootパターンで明示的なDI

### 🎯 設計思想

- **内側は外側を知らない**: ドメイン層は外部実装に依存しない
- **インターフェースで抽象化**: 具象実装は交換可能
- **明示的な依存注入**: `cmd/`で全ての依存を組み立て
- **技術的詳細の隠蔽**: Infrastructureが外部サービスを抽象化

### 📚 参考資料

- [The Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Clean Architecture Example](https://github.com/bxcodec/go-clean-arch)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
