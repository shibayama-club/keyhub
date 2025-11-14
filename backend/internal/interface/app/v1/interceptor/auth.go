package interceptor

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/shibayama-club/keyhub/internal/domain"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/usecase/app/iface"
)

type AuthInterceptor struct {
	useCase iface.IUseCase
}

func NewAuthInterceptor(useCase iface.IUseCase) *AuthInterceptor {
	return &AuthInterceptor{
		useCase: useCase,
	}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if strings.Contains(req.Spec().Procedure, "Health") {
			return next(ctx, req)
		}

		cookies := req.Header().Get("Cookie")
		sessionID := extractSessionID(cookies)

		if sessionID == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}

		user, err := i.useCase.GetMe(ctx, sessionID)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		ctx = domain.WithValue(ctx, user.UserId)
		appSessionID, _ := model.NewAppSessionID(sessionID)
		ctx = domain.WithValue(ctx, appSessionID)

		return next(ctx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func extractSessionID(cookies string) string {
	parts := strings.Split(cookies, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "session_id=") {
			return strings.TrimPrefix(part, "session_id=")
		}
	}
	return ""
}
