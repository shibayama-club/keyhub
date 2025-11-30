# 実装手順

## Phase 1: データベーススキーマとマイグレーション

### 1. マイグレーションファイルの作成
- `db/migrations/20251123000001_add_rooms_table.sql`
- `db/migrations/20251123000002_add_keys_table.sql`
- `db/migrations/20251123000003_add_room_assignments_table.sql`

### 2. SQLCクエリの作成
- `db/sqlc/queries/room.sql`
- `db/sqlc/queries/key.sql`
- `db/sqlc/queries/room_assignment.sql`

### 3. スキーマ生成
```bash
sqlc generate
```

---

## Phase 2: ドメイン層の実装

### 4. ドメインモデルの実装
- `internal/domain/model/room.go`
  - RoomID, RoomName, BuildingName, FloorNumber, RoomType, RoomDescription (Value Objects)
  - Room Entity
- `internal/domain/model/key.go`
  - KeyID, KeyNumber, KeyStatus (Value Objects)
  - Key Entity
- `internal/domain/model/room_assignment.go`
  - RoomAssignmentID (Value Object)
  - RoomAssignment Entity

### 5. リポジトリインターフェースの定義
- `internal/domain/repository/room.go`
  - CreateRoom, GetAllRooms, GetRoomByID, UpdateRoom, DeleteRoom
- `internal/domain/repository/key.go`
  - CreateKey, GetKeysByRoomID, GetKeyByID, UpdateKey, DeleteKey
- `internal/domain/repository/room_assignment.go`
  - CreateAssignment, GetAssignmentsByTenantID, UnassignRoom

---

## Phase 3: インフラストラクチャ層の実装

### 6. リポジトリの実装
- `internal/infrastructure/sqlc/room.go`
  - SQLCで生成されたコードを使用してリポジトリインターフェースを実装
- `internal/infrastructure/sqlc/key.go`
  - SQLCで生成されたコードを使用してリポジトリインターフェースを実装
- `internal/infrastructure/sqlc/room_assignment.go`
  - SQLCで生成されたコードを使用してリポジトリインターフェースを実装

---

## Phase 4: ユースケース層の実装

### 7. ユースケースの実装
- `internal/usecase/console/room.go`
  - CreateRoom, GetAllRooms, GetRoomByID, UpdateRoom, DeleteRoom
  - AssignRoomToTenant, UnassignRoomFromTenant, GetRoomsByTenantID
- `internal/usecase/console/key.go`
  - CreateKey, GetKeysByRoomID, GetKeyByID, UpdateKey, DeleteKey
- `internal/usecase/console/room_assignment.go`
  - AssignRoom, UnassignRoom, GetAssignments (必要に応じて)

### 8. DTOの定義
- `internal/usecase/console/dto/room.go`
  - RoomDTO, RoomWithAssignmentDTO
- `internal/usecase/console/dto/key.go`
  - KeyDTO
- `internal/usecase/console/dto/room_assignment.go`
  - RoomAssignmentDTO

---

## Phase 5: インターフェース層の実装

### 9. Protobufの定義と生成
- `proto/keyhub/console/v1/console.proto` (追加)
  - ConsoleRoomService, ConsoleKeyService
  - 関連するメッセージとEnumの定義
- コード生成
  ```bash
  buf generate
  ```

### 10. ハンドラーの実装
- `internal/interface/console/v1/room_handler.go`
  - CreateRoom, GetAllRooms, GetRoomById, UpdateRoom, DeleteRoom
  - AssignRoomToTenant, UnassignRoomFromTenant, GetRoomsByTenantId
- `internal/interface/console/v1/key_handler.go`
  - CreateKey, GetKeysByRoomId, GetKeyById, UpdateKey, DeleteKey
- `internal/interface/console/v1/room_assignment_handler.go` (必要に応じて)

---

## Phase 6: Composition Rootでの統合

### 11. サービスの登録
- `cmd/serve/console.go` (修正)
  - RoomRepository, KeyRepository, RoomAssignmentRepository のインスタンス生成
  - RoomUseCase, KeyUseCase のインスタンス生成
  - RoomHandler, KeyHandler のインスタンス生成
  - Connect RPCサーバーへのサービス登録

---

## Phase 7: テストの実装

### 12. 単体テスト
- `internal/domain/model/*_test.go`
  - Value Objectのバリデーションテスト
  - Entityの生成テスト
- `internal/usecase/console/*_test.go`
  - ユースケースのビジネスロジックテスト
  - モックリポジトリを使用したテスト

### 13. 統合テスト
- `internal/infrastructure/sqlc/*_test.go`
  - リポジトリの実装テスト
  - データベースとの連携テスト
  - トランザクション管理のテスト

---

## 実装の順序

1. **Phase 1**: データベーススキーマとマイグレーションから開始
2. **Phase 2-3**: ドメイン層とインフラ層を並行して実装
3. **Phase 4**: ユースケース層の実装
4. **Phase 5**: API定義とハンドラーの実装
5. **Phase 6**: Composition Rootでの統合
6. **Phase 7**: テストの実装と品質保証

各フェーズが完了したら、必ずテストを実施し、次のフェーズに移行します。
