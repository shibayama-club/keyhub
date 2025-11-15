package console

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/shibayama-club/keyhub/cmd/config"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/domain/repository/mock"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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
					expectedTenant := model.Tenant{
						ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")),
						OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
						Name:           "テストテナント",
						Description:    "テスト説明",
						Type:           model.TenantTypeTeam,
					}

					m.EXPECT().
						WithTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context, repository.Transaction) error) error {
							// トランザクション内の処理をシミュレート
							mockTx := mock.NewMockTransaction(gomock.NewController(t))
							mockTx.EXPECT().
								CreateTenant(gomock.Any(), gomock.Any()).
								Return(expectedTenant, nil)
							mockTx.EXPECT().
								CreateTenantJoinCode(gomock.Any(), gomock.Any()).
								Return(model.TenantJoinCodeEntity{}, nil)
							return fn(ctx, mockTx)
						})
				},
			},
			args: args{
				ctx: context.Background(),
				input: dto.CreateTenantInput{
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           "テストテナント",
					Description:    "テスト説明",
					TenantType:     model.TenantTypeTeam.String(),
					JoinCode:       "testcode123",
					JoinCodeExpiry: nil,
					JoinCodeMaxUse: 0,
				},
			},
			want:    "550e8400-e29b-41d4-a716-446655440001",
			wantErr: false,
		},
		{
			name: "異常系: 名前が空",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					// モックは呼ばれない
				},
			},
			args: args{
				ctx: context.Background(),
				input: dto.CreateTenantInput{
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           "",
					Description:    "テスト説明",
					TenantType:     model.TenantTypeTeam.String(),
					JoinCode:       "testcode123",
					JoinCodeExpiry: nil,
					JoinCodeMaxUse: 0,
				},
			},
			want:    "",
			wantErr: true,
			errType: domainerrors.ErrValidation,
		},
		{
			name: "異常系: 不正なテナントタイプ",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					// モックは呼ばれない
				},
			},
			args: args{
				ctx: context.Background(),
				input: dto.CreateTenantInput{
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           "テストテナント",
					Description:    "テスト説明",
					TenantType:     "invalid_type",
					JoinCode:       "testcode123",
					JoinCodeExpiry: nil,
					JoinCodeMaxUse: 0,
				},
			},
			want:    "",
			wantErr: true,
			errType: domainerrors.ErrValidation,
		},
		{
			name: "異常系: 名前が30文字を超える",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					// モックは呼ばれない
				},
			},
			args: args{
				ctx: context.Background(),
				input: dto.CreateTenantInput{
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           string(make([]byte, 31)), // 31文字
					Description:    "テスト説明",
					TenantType:     model.TenantTypeTeam.String(),
					JoinCode:       "testcode123",
					JoinCodeExpiry: nil,
					JoinCodeMaxUse: 0,
				},
			},
			want:    "",
			wantErr: true,
			errType: domainerrors.ErrValidation,
		},
		{
			name: "正常系: 説明が空でも成功",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					expectedTenant := model.Tenant{
						ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")),
						OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
						Name:           "テストテナント2",
						Description:    "",
						Type:           model.TenantTypeDepartment,
					}

					m.EXPECT().
						WithTransaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(context.Context, repository.Transaction) error) error {
							mockTx := mock.NewMockTransaction(gomock.NewController(t))
							mockTx.EXPECT().
								CreateTenant(gomock.Any(), gomock.Any()).
								Return(expectedTenant, nil)
							mockTx.EXPECT().
								CreateTenantJoinCode(gomock.Any(), gomock.Any()).
								Return(model.TenantJoinCodeEntity{}, nil)
							return fn(ctx, mockTx)
						})
				},
			},
			args: args{
				ctx: context.Background(),
				input: dto.CreateTenantInput{
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           "テストテナント2",
					Description:    "",
					TenantType:     model.TenantTypeDepartment.String(),
					JoinCode:       "testcode456",
					JoinCodeExpiry: nil,
					JoinCodeMaxUse: 0,
				},
			},
			want:    "550e8400-e29b-41d4-a716-446655440002",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tt.fields.setupMock(mockRepo)

			u := &UseCase{
				repo:   mockRepo,
				config: config.Config{},
			}

			// Act
			got, err := u.CreateTenant(tt.args.ctx, tt.args.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.Is(err, tt.errType), "expected error type %v, got %v", tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUseCase_GetAllTenants(t *testing.T) {
	type fields struct {
		setupMock func(*mock.MockRepository)
	}
	type args struct {
		ctx            context.Context
		organizationID model.OrganizationID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Tenant
		wantErr bool
	}{
		{
			name: "正常系: テナント一覧取得成功",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					expectedTenants := []model.Tenant{
						{
							ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")),
							OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
							Name:           model.TenantName("テナント1"),
							Description:    model.TenantDescription("説明1"),
							Type:           model.TenantTypeTeam,
						},
						{
							ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")),
							OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
							Name:           model.TenantName("テナント2"),
							Description:    model.TenantDescription("説明2"),
							Type:           model.TenantTypeDepartment,
						},
					}

					m.EXPECT().
						GetAllTenants(gomock.Any(), gomock.Any()).
						Return(expectedTenants, nil)
				},
			},
			args: args{
				ctx:            context.Background(),
				organizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
			},
			want: []model.Tenant{
				{
					ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")),
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           model.TenantName("テナント1"),
					Description:    model.TenantDescription("説明1"),
					Type:           model.TenantTypeTeam,
				},
				{
					ID:             model.TenantID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")),
					OrganizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
					Name:           model.TenantName("テナント2"),
					Description:    model.TenantDescription("説明2"),
					Type:           model.TenantTypeDepartment,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系: テナントが0件でも成功",
			fields: fields{
				setupMock: func(m *mock.MockRepository) {
					m.EXPECT().
						GetAllTenants(gomock.Any(), gomock.Any()).
						Return([]model.Tenant{}, nil)
				},
			},
			args: args{
				ctx:            context.Background(),
				organizationID: model.OrganizationID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
			},
			want:    []model.Tenant{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tt.fields.setupMock(mockRepo)

			u := &UseCase{
				repo:   mockRepo,
				config: config.Config{},
			}

			// Act
			got, err := u.GetAllTenants(tt.args.ctx, tt.args.organizationID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
