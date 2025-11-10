# App API仕様

Main App（エンドユーザー向け）のAPI仕様です。ConnectRPCを使用しています。

## AuthService - 認証サービス

```proto
service AuthService {
    // 現在のユーザー情報取得
    rpc GetMe(GetMeRequest) returns (GetMeResponse);

    // ログアウト
    rpc Logout(LogoutRequest) returns (LogoutResponse);
}
```

### GetMe
現在ログイン中のユーザー情報を取得します。

**リクエスト**: なし（セッションから自動取得）

**レスポンス**:
```proto
message GetMeResponse {
    User user = 1;
}
```

**注意**: テナント情報は取得しません。テナント情報が必要な場合は `TenantService.GetMyTenants` または `TenantService.GetActiveTenant` を使用してください。

### Logout
ユーザーをログアウトさせ、セッションを無効化します。

**リクエスト**: なし

**レスポンス**: 成功/失敗のステータス

## TenantService - テナント管理サービス

```proto
service TenantService {
    // 参加可能なTenant一覧取得
    rpc ListAvailableTenants(ListAvailableTenantsRequest) returns (ListAvailableTenantsResponse);

    // 自分が所属するTenant一覧
    rpc GetMyTenants(GetMyTenantsRequest) returns (GetMyTenantsResponse);

    // 現在アクティブなTenant取得
    rpc GetActiveTenant(GetActiveTenantRequest) returns (GetActiveTenantResponse);

    // アクティブTenant切り替え
    rpc SetActiveTenant(SetActiveTenantRequest) returns (SetActiveTenantResponse);

    // Tenant参加
    rpc JoinTenant(JoinTenantRequest) returns (JoinTenantResponse);

    // 参加コードでTenant参加
    rpc JoinByCode(JoinByCodeRequest) returns (JoinTenantResponse);

    // Tenantメンバー一覧
    rpc ListTenantMembers(ListTenantMembersRequest) returns (ListTenantMembersResponse);

    // Tenantから退出
    rpc LeaveTenant(LeaveTenantRequest) returns (LeaveTenantResponse);

    // Tenant詳細取得
    rpc GetTenant(GetTenantRequest) returns (GetTenantResponse);
}
```

### ListAvailableTenants
参加可能な公開テナント一覧を取得します。ユーザーが参加可能なテナントのみが返されます。

**リクエスト**:
```proto
message ListAvailableTenantsRequest {
    int32 page_size = 1 [(buf.validate.field).int32 = {gte: 1, lte: 100}];
    string page_token = 2;
}
```

**レスポンス**:
```proto
message ListAvailableTenantsResponse {
    repeated Tenant tenants = 1;
    string next_page_token = 2;
}
```

### GetMyTenants
現在ログイン中のユーザーが所属する全テナントの一覧を取得します。

**リクエスト**: なし

**レスポンス**:
```proto
message GetMyTenantsResponse {
    repeated TenantMembership memberships = 1;
}
```

### GetActiveTenant
現在アクティブなテナント情報を取得します。テナントが設定されていない場合はエラーを返します。

**リクエスト**: なし

**レスポンス**:
```proto
message GetActiveTenantResponse {
    Tenant tenant = 1;
    TenantMembership membership = 2;
}
```

### SetActiveTenant
アクティブなテナントを切り替えます。指定したテナントに所属している必要があります。

**リクエスト**:
```proto
message SetActiveTenantRequest {
    string membership_id = 1 [(buf.validate.field).string.uuid = true];
}
```

**レスポンス**:
```proto
message SetActiveTenantResponse {
    Tenant tenant = 1;
}
```

### JoinTenant
指定したテナントに参加します。

**リクエスト**:
```proto
message JoinTenantRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
}
```

**レスポンス**:
```proto
message JoinTenantResponse {
    TenantMembership membership = 1;
}
```

### JoinByCode
参加コードを使用してテナントに参加します。

**リクエスト**:
```proto
message JoinByCodeRequest {
    string code = 1 [(buf.validate.field).string = {pattern: "^KH-[A-Z0-9]{5}-[A-Z0-9]{2}$"}];
}
```

**レスポンス**: `JoinTenantResponse`と同じ

### ListTenantMembers
指定したテナントのメンバー一覧を取得します。

**リクエスト**:
```proto
message ListTenantMembersRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
    int32 page_size = 2 [(buf.validate.field).int32 = {gte: 1, lte: 100}];
    string page_token = 3;
}
```

**レスポンス**:
```proto
message ListTenantMembersResponse {
    repeated TenantMember members = 1;
    string next_page_token = 2;
}
```

## データ型定義

### Enum定義

```proto
// ユーザーロール
enum Role {
    ROLE_UNSPECIFIED = 0;
    ROLE_VIEWER = 1;      // 閲覧のみ
    ROLE_MEMBER = 2;      // 通常メンバー
    ROLE_ADMIN = 3;       // 管理者
    ROLE_OWNER = 4;       // オーナー
}

// メンバーシップステータス
enum MembershipStatus {
    MEMBERSHIP_STATUS_UNSPECIFIED = 0;
    MEMBERSHIP_STATUS_ACTIVE = 1;      // アクティブ
    MEMBERSHIP_STATUS_INACTIVE = 2;    // 非アクティブ
    MEMBERSHIP_STATUS_SUSPENDED = 3;   // 停止中
    MEMBERSHIP_STATUS_INVITED = 4;     // 招待中
}

// Tenantタイプ
enum TenantType {
    TENANT_TYPE_UNSPECIFIED = 0;
    TENANT_TYPE_TEAM = 1;        // チーム
    TENANT_TYPE_DEPARTMENT = 2;  // 部署
    TENANT_TYPE_PROJECT = 3;     // プロジェクト
    TENANT_TYPE_LABORATORY = 4;  // 研究室
}
```

### メッセージ定義

```proto
message User {
    string id = 1 [(buf.validate.field).string.uuid = true];  // UUID形式
    string email = 2 [(buf.validate.field).string.email = true];
    string name = 3 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    string icon = 4 [(buf.validate.field).string.uri = true];  // URL形式
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
}

message Tenant {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string organization_id = 2 [(buf.validate.field).string.uuid = true];  // 組織ID (UUID形式)
    string name = 3 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    string slug = 4 [(buf.validate.field).string = {pattern: "^[a-z0-9-]+$", min_len: 3, max_len: 50}];  // URL用識別子
    string description = 5 [(buf.validate.field).string.max_len = 500];
    TenantType tenant_type = 6;
    bool is_default = 7;  // 組織内のデフォルトテナントフラグ
    int32 member_count = 8;
    google.protobuf.Timestamp created_at = 9;
    google.protobuf.Timestamp updated_at = 10;
}

message TenantMember {
    string user_id = 1 [(buf.validate.field).string.uuid = true];
    string email = 2 [(buf.validate.field).string.email = true];
    string name = 3;
    string icon = 4;
    Role role = 5;
    MembershipStatus status = 6;
    google.protobuf.Timestamp joined_at = 7;
    google.protobuf.Timestamp left_at = 8;
}

message TenantMembership {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string tenant_id = 2 [(buf.validate.field).string.uuid = true];
    string user_id = 3 [(buf.validate.field).string.uuid = true];
    Tenant tenant = 4;
    Role role = 5;
    MembershipStatus status = 6;
    google.protobuf.Timestamp joined_at = 7;
    google.protobuf.Timestamp left_at = 8;
    google.protobuf.Timestamp updated_at = 9;
}

message TenantJoinCode {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string tenant_id = 2 [(buf.validate.field).string.uuid = true];
    string code = 3 [(buf.validate.field).string = {pattern: "^KH-[A-Z0-9]{5}-[A-Z0-9]{2}$"}];
    google.protobuf.Timestamp expires_at = 4;
    int32 max_uses = 5 [(buf.validate.field).int32.gte = 0];
    int32 used_count = 6 [(buf.validate.field).int32.gte = 0];
    Role role = 7;  // 参加時に付与されるロール
    string created_by = 8 [(buf.validate.field).string.uuid = true];
    google.protobuf.Timestamp created_at = 9;
}
```

## HTTP エンドポイント

認証関連のHTTPエンドポイント：

| エンドポイント | メソッド | 説明 |
|---------------|---------|------|
| `/auth/google/login` | GET | Google OAuth認証開始 |
| `/auth/google/callback` | GET | OAuth認証コールバック |

ConnectRPCエンドポイント：
- ベースURL: `/connect`
- Content-Type: `application/connect+proto`

## エラーコード

| コード | 説明 |
|--------|------|
| `UNAUTHENTICATED` | 認証が必要 |
| `PERMISSION_DENIED` | 権限不足 |
| `NOT_FOUND` | リソースが見つからない |
| `ALREADY_EXISTS` | すでに存在する（重複参加など） |
| `INVALID_ARGUMENT` | 無効なパラメータ |
| `FAILED_PRECONDITION` | 前提条件を満たしていない |

## 完全なProtobufスキーマ

```proto
syntax = "proto3";

package keyhub.app.v1;

import "google/protobuf/timestamp.proto";
import "buf/validate/validate.proto";

option go_package = "github.com/yourusername/keyhub/gen/go/keyhub/app/v1;appv1";

// ========== Enums ==========

// ユーザーロール
enum Role {
    ROLE_UNSPECIFIED = 0;
    ROLE_VIEWER = 1;      // 閲覧のみ
    ROLE_MEMBER = 2;      // 通常メンバー
    ROLE_ADMIN = 3;       // 管理者
    ROLE_OWNER = 4;       // オーナー
}

// メンバーシップステータス
enum MembershipStatus {
    MEMBERSHIP_STATUS_UNSPECIFIED = 0;
    MEMBERSHIP_STATUS_ACTIVE = 1;      // アクティブ
    MEMBERSHIP_STATUS_INACTIVE = 2;    // 非アクティブ
    MEMBERSHIP_STATUS_SUSPENDED = 3;   // 停止中
    MEMBERSHIP_STATUS_INVITED = 4;     // 招待中
}

// Tenantタイプ
enum TenantType {
    TENANT_TYPE_UNSPECIFIED = 0;
    TENANT_TYPE_TEAM = 1;        // チーム
    TENANT_TYPE_DEPARTMENT = 2;  // 部署
    TENANT_TYPE_PROJECT = 3;     // プロジェクト
    TENANT_TYPE_LABORATORY = 4;  // 研究室
}

// ========== Common Messages ==========

message User {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string email = 2 [(buf.validate.field).string.email = true];
    string name = 3 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    string icon = 4 [(buf.validate.field).string.uri = true];
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
}

message Tenant {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string organization_id = 2 [(buf.validate.field).string.uuid = true];  // 組織ID (UUID形式)
    string name = 3 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    string slug = 4 [(buf.validate.field).string = {pattern: "^[a-z0-9-]+$", min_len: 3, max_len: 50}];  // URL用識別子
    string description = 5 [(buf.validate.field).string.max_len = 500];
    TenantType tenant_type = 6;
    bool is_default = 7;  // 組織内のデフォルトテナントフラグ
    int32 member_count = 8;
    google.protobuf.Timestamp created_at = 9;
    google.protobuf.Timestamp updated_at = 10;
}

message TenantMember {
    string user_id = 1 [(buf.validate.field).string.uuid = true];
    string email = 2 [(buf.validate.field).string.email = true];
    string name = 3;
    string icon = 4;
    Role role = 5;
    MembershipStatus status = 6;
    google.protobuf.Timestamp joined_at = 7;
    google.protobuf.Timestamp left_at = 8;
}

message TenantMembership {
    string id = 1 [(buf.validate.field).string.uuid = true];
    string tenant_id = 2 [(buf.validate.field).string.uuid = true];
    string user_id = 3 [(buf.validate.field).string.uuid = true];
    Tenant tenant = 4;
    Role role = 5;
    MembershipStatus status = 6;
    google.protobuf.Timestamp joined_at = 7;
    google.protobuf.Timestamp left_at = 8;
    google.protobuf.Timestamp updated_at = 9;
}

// ========== AuthService ==========

service AuthService {
    // 現在のユーザー情報取得
    rpc GetMe(GetMeRequest) returns (GetMeResponse);

    // ログアウト
    rpc Logout(LogoutRequest) returns (LogoutResponse);

    // セッション検証
    rpc ValidateSession(ValidateSessionRequest) returns (ValidateSessionResponse);

    // アクティブTenant切り替え
    rpc SwitchTenant(SwitchTenantRequest) returns (SwitchTenantResponse);
}

message GetMeRequest {}

message GetMeResponse {
    User user = 1;
}

message LogoutRequest {}

message LogoutResponse {
    bool success = 1;
}

message ValidateSessionRequest {}

message ValidateSessionResponse {
    bool valid = 1;
    User user = 2;
}

// ========== TenantService ==========

service TenantService {
    // 参加可能なTenant一覧取得
    rpc ListAvailableTenants(ListAvailableTenantsRequest) returns (ListAvailableTenantsResponse);

    // 自分が所属するTenant一覧
    rpc GetMyTenants(GetMyTenantsRequest) returns (GetMyTenantsResponse);

    // 現在アクティブなTenant取得
    rpc GetActiveTenant(GetActiveTenantRequest) returns (GetActiveTenantResponse);

    // アクティブTenant切り替え
    rpc SetActiveTenant(SetActiveTenantRequest) returns (SetActiveTenantResponse);

    // Tenant参加
    rpc JoinTenant(JoinTenantRequest) returns (JoinTenantResponse);

    // 参加コードでTenant参加
    rpc JoinByCode(JoinByCodeRequest) returns (JoinTenantResponse);

    // Tenantメンバー一覧
    rpc ListTenantMembers(ListTenantMembersRequest) returns (ListTenantMembersResponse);

    // Tenantから退出
    rpc LeaveTenant(LeaveTenantRequest) returns (LeaveTenantResponse);

    // Tenant詳細取得
    rpc GetTenant(GetTenantRequest) returns (GetTenantResponse);
}

message ListAvailableTenantsRequest {
    int32 page_size = 1 [(buf.validate.field).int32 = {gte: 1, lte: 100}];
    string page_token = 2;
}

message ListAvailableTenantsResponse {
    repeated Tenant tenants = 1;
    string next_page_token = 2;
}

message GetMyTenantsRequest {}

message GetMyTenantsResponse {
    repeated TenantMembership memberships = 1;
}

message GetActiveTenantRequest {}

message GetActiveTenantResponse {
    Tenant tenant = 1;
    TenantMembership membership = 2;
}

message SetActiveTenantRequest {
    string membership_id = 1 [(buf.validate.field).string.uuid = true];
}

message SetActiveTenantResponse {
    Tenant tenant = 1;
}

message JoinTenantRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
}

message JoinTenantResponse {
    TenantMembership membership = 1;
}

message JoinByCodeRequest {
    string code = 1 [(buf.validate.field).string = {pattern: "^KH-[A-Z0-9]{5}-[A-Z0-9]{2}$"}];
}

message ListTenantMembersRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
    int32 page_size = 2 [(buf.validate.field).int32 = {gte: 1, lte: 100}];
    string page_token = 3;
}

message ListTenantMembersResponse {
    repeated TenantMember members = 1;
    string next_page_token = 2;
}

message LeaveTenantRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
}

message LeaveTenantResponse {
    bool success = 1;
}

message GetTenantRequest {
    string tenant_id = 1 [(buf.validate.field).string.uuid = true];
}

message GetTenantResponse {
    Tenant tenant = 1;
    TenantMembership membership = 2;
}

// ========== UserService ==========

service UserService {
    // ユーザープロフィール取得
    rpc GetUser(GetUserRequest) returns (GetUserResponse);

    // ユーザープロフィール更新
    rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);

    // ユーザー検索
    rpc SearchUsers(SearchUsersRequest) returns (SearchUsersResponse);
}

message GetUserRequest {
    string user_id = 1 [(buf.validate.field).string.uuid = true];
}

message GetUserResponse {
    User user = 1;
}

message UpdateProfileRequest {
    string name = 1 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    string icon = 2 [(buf.validate.field).string.uri = true];
}

message UpdateProfileResponse {
    User user = 1;
}

message SearchUsersRequest {
    string query = 1 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
    int32 page_size = 2 [(buf.validate.field).int32 = {gte: 1, lte: 50}];
    string page_token = 3;
}

message SearchUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
}
```