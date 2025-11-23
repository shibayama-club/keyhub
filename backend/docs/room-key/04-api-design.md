# API設計

## Protobuf定義（Console API）

`proto/keyhub/console/v1/console.proto` に以下を追加:

```protobuf
// ============================================================
// Room Management
// ============================================================

service ConsoleRoomService {
  // 部屋を作成
  rpc CreateRoom(CreateRoomRequest) returns (CreateRoomResponse);

  // 組織の全部屋を取得
  rpc GetAllRooms(GetAllRoomsRequest) returns (GetAllRoomsResponse);

  // IDから部屋を取得
  rpc GetRoomById(GetRoomByIdRequest) returns (GetRoomByIdResponse);

  // 部屋情報を更新
  rpc UpdateRoom(UpdateRoomRequest) returns (UpdateRoomResponse);

  // 部屋を削除
  rpc DeleteRoom(DeleteRoomRequest) returns (DeleteRoomResponse);

  // テナントに部屋を割り当て
  rpc AssignRoomToTenant(AssignRoomToTenantRequest) returns (AssignRoomToTenantResponse);

  // テナントの部屋割り当てを解除
  rpc UnassignRoomFromTenant(UnassignRoomFromTenantRequest) returns (UnassignRoomFromTenantResponse);

  // テナントに割り当てられた部屋一覧を取得
  rpc GetRoomsByTenantId(GetRoomsByTenantIdRequest) returns (GetRoomsByTenantIdResponse);
}

// ============================================================
// Key Management
// ============================================================

service ConsoleKeyService {
  // 鍵を作成
  rpc CreateKey(CreateKeyRequest) returns (CreateKeyResponse);

  // 部屋の鍵一覧を取得
  rpc GetKeysByRoomId(GetKeysByRoomIdRequest) returns (GetKeysByRoomIdResponse);

  // IDから鍵を取得
  rpc GetKeyById(GetKeyByIdRequest) returns (GetKeyByIdResponse);

  // 鍵情報を更新
  rpc UpdateKey(UpdateKeyRequest) returns (UpdateKeyResponse);

  // 鍵を削除
  rpc DeleteKey(DeleteKeyRequest) returns (DeleteKeyResponse);
}

// ============================================================
// Messages - Room
// ============================================================

enum RoomType {
  ROOM_TYPE_UNSPECIFIED = 0;
  ROOM_TYPE_CLASSROOM = 1;      // 教室
  ROOM_TYPE_MEETING_ROOM = 2;   // 会議室
  ROOM_TYPE_LABORATORY = 3;     // 実験室
  ROOM_TYPE_OFFICE = 4;         // オフィス
  ROOM_TYPE_WORKSHOP = 5;       // 作業室
  ROOM_TYPE_STORAGE = 6;        // 倉庫
}

message Room {
  string id = 1 [(buf.validate.field).string.uuid = true];
  string organization_id = 2 [(buf.validate.field).string.uuid = true];
  string name = 3;
  string location = 4;
  RoomType room_type = 5;
  optional int32 capacity = 6;
  string description = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}

message CreateRoomRequest {
  string name = 1 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
  string location = 2 [(buf.validate.field).string = {min_len: 1, max_len: 200}];
  RoomType room_type = 3;
  optional int32 capacity = 4 [(buf.validate.field).int32 = {gt: 0}];
  string description = 5 [(buf.validate.field).string.max_len = 500];
}

message CreateRoomResponse {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message GetAllRoomsRequest {}

message GetAllRoomsResponse {
  repeated Room rooms = 1;
}

message GetRoomByIdRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message GetRoomByIdResponse {
  Room room = 1;
  repeated Key keys = 2; // この部屋に関連する鍵
}

message UpdateRoomRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
  string name = 2 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
  string location = 3 [(buf.validate.field).string = {min_len: 1, max_len: 200}];
  RoomType room_type = 4;
  optional int32 capacity = 5 [(buf.validate.field).int32 = {gt: 0}];
  string description = 6 [(buf.validate.field).string.max_len = 500];
}

message UpdateRoomResponse {}

message DeleteRoomRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message DeleteRoomResponse {}

// ============================================================
// Messages - Room Assignment
// ============================================================

message RoomAssignment {
  string id = 1 [(buf.validate.field).string.uuid = true];
  string tenant_id = 2 [(buf.validate.field).string.uuid = true];
  string room_id = 3 [(buf.validate.field).string.uuid = true];
  google.protobuf.Timestamp assigned_at = 4;
  optional google.protobuf.Timestamp expires_at = 5;
}

message AssignRoomToTenantRequest {
  string tenant_id = 1 [(buf.validate.field).string.uuid = true];
  string room_id = 2 [(buf.validate.field).string.uuid = true];
  optional google.protobuf.Timestamp expires_at = 3;
}

message AssignRoomToTenantResponse {
  string assignment_id = 1 [(buf.validate.field).string.uuid = true];
}

message UnassignRoomFromTenantRequest {
  string assignment_id = 1 [(buf.validate.field).string.uuid = true];
}

message UnassignRoomFromTenantResponse {}

message GetRoomsByTenantIdRequest {
  string tenant_id = 1 [(buf.validate.field).string.uuid = true];
}

message RoomWithAssignment {
  Room room = 1;
  RoomAssignment assignment = 2;
  repeated Key keys = 3; // この部屋に関連する鍵
}

message GetRoomsByTenantIdResponse {
  repeated RoomWithAssignment rooms = 1;
}

// ============================================================
// Messages - Key
// ============================================================

enum KeyType {
  KEY_TYPE_UNSPECIFIED = 0;
  KEY_TYPE_PHYSICAL = 1;  // 物理鍵
  KEY_TYPE_CARD = 2;      // カードキー
  KEY_TYPE_DIGITAL = 3;   // デジタルキー
}

enum KeyStatus {
  KEY_STATUS_UNSPECIFIED = 0;
  KEY_STATUS_AVAILABLE = 1;  // 利用可能
  KEY_STATUS_IN_USE = 2;     // 貸出中
  KEY_STATUS_LOST = 3;       // 紛失
  KEY_STATUS_DAMAGED = 4;    // 破損
}

message Key {
  string id = 1 [(buf.validate.field).string.uuid = true];
  string room_id = 2 [(buf.validate.field).string.uuid = true];
  string organization_id = 3 [(buf.validate.field).string.uuid = true];
  string key_number = 4;
  KeyType key_type = 5;
  KeyStatus status = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
}

message CreateKeyRequest {
  string room_id = 1 [(buf.validate.field).string.uuid = true];
  string key_number = 2 [(buf.validate.field).string = {min_len: 1, max_len: 50}];
  KeyType key_type = 3;
}

message CreateKeyResponse {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message GetKeysByRoomIdRequest {
  string room_id = 1 [(buf.validate.field).string.uuid = true];
}

message GetKeysByRoomIdResponse {
  repeated Key keys = 1;
}

message GetKeyByIdRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message GetKeyByIdResponse {
  Key key = 1;
}

message UpdateKeyRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
  string key_number = 2 [(buf.validate.field).string = {min_len: 1, max_len: 50}];
  KeyType key_type = 3;
  KeyStatus status = 4;
}

message UpdateKeyResponse {}

message DeleteKeyRequest {
  string id = 1 [(buf.validate.field).string.uuid = true];
}

message DeleteKeyResponse {}
```
